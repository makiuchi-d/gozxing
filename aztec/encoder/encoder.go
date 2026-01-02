package encoder

import (
	"fmt"
	"golang.org/x/text/encoding"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common/reedsolomon"
)

const (
	DEFAULT_EC_PERCENT      = 33 // default minimal percentage of error check words
	DEFAULT_AZTEC_LAYERS     = 0
	MAX_NB_BITS              = 32
	MAX_NB_BITS_COMPACT      = 4
)

var WORD_SIZE = []int{
	4, 6, 6, 8, 8, 8, 8, 8, 8, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
}

// EncodeString Encodes the given string content as an Aztec symbol (without ECI code)
func EncodeString(data string) (*AztecCode, error) {
	return EncodeStringWithParams(data, DEFAULT_EC_PERCENT, DEFAULT_AZTEC_LAYERS, nil)
}

// EncodeStringWithParams Encodes the given string content as an Aztec symbol
func EncodeStringWithParams(data string, minECCPercent, userSpecifiedLayers int, charset encoding.Encoding) (*AztecCode, error) {
	var bytes []byte
	var err error
	if charset != nil {
		enc := charset.NewEncoder()
		bytes, err = enc.Bytes([]byte(data))
		if err != nil {
			return nil, gozxing.WrapWriterException(err)
		}
	} else {
		// Default to ISO-8859-1 (Latin-1)
		// ISO-8859-1 is a single-byte encoding, so we can directly convert
		bytes = []byte(data)
	}
	return EncodeBytesWithParams(bytes, minECCPercent, userSpecifiedLayers, charset)
}

// EncodeBytes Encodes the given binary content as an Aztec symbol (without ECI code)
func EncodeBytes(data []byte) (*AztecCode, error) {
	return EncodeBytesWithParams(data, DEFAULT_EC_PERCENT, DEFAULT_AZTEC_LAYERS, nil)
}

// EncodeBytesWithParams Encodes the given binary content as an Aztec symbol
func EncodeBytesWithParams(data []byte, minECCPercent, userSpecifiedLayers int, charset encoding.Encoding) (*AztecCode, error) {
	// High-level encode
	var bits *gozxing.BitArray
	if charset != nil {
		bits = NewHighLevelEncoderWithCharset(data, charset).Encode()
	} else {
		bits = NewHighLevelEncoder(data).Encode()
	}

	// stuff bits and choose symbol size
	eccBits := bits.GetSize()*minECCPercent/100 + 11
	totalSizeBits := bits.GetSize() + eccBits
	var compact bool
	var layers int
	var totalBitsInLayer int
	var wordSize int
	var stuffedBits *gozxing.BitArray

	if userSpecifiedLayers != DEFAULT_AZTEC_LAYERS {
		compact = userSpecifiedLayers < 0
		layers = userSpecifiedLayers
		if layers < 0 {
			layers = -layers
		}
		maxLayers := MAX_NB_BITS
		if compact {
			maxLayers = MAX_NB_BITS_COMPACT
		}
		if layers > maxLayers {
			return nil, gozxing.NewWriterException(
				"IllegalArgumentException: Illegal value %d for layers", userSpecifiedLayers)
		}
		totalBitsInLayer = calculateTotalBitsInLayer(layers, compact)
		wordSize = WORD_SIZE[layers]
		usableBitsInLayers := totalBitsInLayer - (totalBitsInLayer % wordSize)
		stuffedBits = stuffBits(bits, wordSize)
		if stuffedBits.GetSize()+eccBits > usableBitsInLayers {
			return nil, gozxing.NewWriterException("Data too large for user specified layer")
		}
		if compact && stuffedBits.GetSize() > wordSize*64 {
			// Compact format only allows 64 data words, though C4 can hold more words than that
			return nil, gozxing.NewWriterException("Data too large for user specified layer")
		}
	} else {
		wordSize = 0
		stuffedBits = nil
		// We look at the possible table sizes in the order Compact1, Compact2, Compact3,
		// Compact4, Normal4,...  Normal(i) for i < 4 isn't typically used since Compact(i+1)
		// is the same size, but has more data.
		for i := 0; ; i++ {
			if i > MAX_NB_BITS {
				return nil, gozxing.NewWriterException("Data too large for an Aztec code")
			}
			compact = i <= 3
			if compact {
				layers = i + 1
			} else {
				layers = i
			}
			totalBitsInLayer = calculateTotalBitsInLayer(layers, compact)
			if totalSizeBits > totalBitsInLayer {
				continue
			}
			// [Re]stuff the bits if this is the first opportunity, or if the
			// wordSize has changed
			if stuffedBits == nil || wordSize != WORD_SIZE[layers] {
				wordSize = WORD_SIZE[layers]
				stuffedBits = stuffBits(bits, wordSize)
			}
			usableBitsInLayers := totalBitsInLayer - (totalBitsInLayer % wordSize)
			if compact && stuffedBits.GetSize() > wordSize*64 {
				// Compact format only allows 64 data words, though C4 can hold more words than that
				continue
			}
			if stuffedBits.GetSize()+eccBits <= usableBitsInLayers {
				break
			}
		}
	}
	messageBits := generateCheckWords(stuffedBits, totalBitsInLayer, wordSize)

	// generate mode message
	messageSizeInWords := stuffedBits.GetSize() / wordSize
	modeMessage := generateModeMessage(compact, layers, messageSizeInWords)

	// allocate symbol
	baseMatrixSize := 11
	if !compact {
		baseMatrixSize = 14
	}
	baseMatrixSize += layers * 4 // not including alignment lines
	alignmentMap := make([]int, baseMatrixSize)
	var matrixSize int
	if compact {
		// no alignment marks in compact mode, alignmentMap is a no-op
		matrixSize = baseMatrixSize
		for i := 0; i < len(alignmentMap); i++ {
			alignmentMap[i] = i
		}
	} else {
		matrixSize = baseMatrixSize + 1 + 2*((baseMatrixSize/2-1)/15)
		origCenter := baseMatrixSize / 2
		center := matrixSize / 2
		for i := 0; i < origCenter; i++ {
			newOffset := i + i/15
			alignmentMap[origCenter-i-1] = center - newOffset - 1
			alignmentMap[origCenter+i] = center + newOffset + 1
		}
	}
	matrix, err := gozxing.NewBitMatrix(matrixSize, matrixSize)
	if err != nil {
		return nil, gozxing.WrapWriterException(err)
	}

	// draw data bits
	rowOffset := 0
	for i := 0; i < layers; i++ {
		rowSize := (layers-i)*4 + 9
		if !compact {
			rowSize = (layers-i)*4 + 12
		}
		low := i * 2
		high := baseMatrixSize - 1 - low
		for j := 0; j < rowSize; j++ {
			columnOffset := j * 2
			for k := 0; k < 2; k++ {
				// left column
				if messageBits.Get(rowOffset + columnOffset + k) {
					matrix.Set(alignmentMap[low+k], alignmentMap[low+j])
				}
				// bottom row
				if messageBits.Get(rowOffset + rowSize*2 + columnOffset + k) {
					matrix.Set(alignmentMap[low+j], alignmentMap[high-k])
				}
				// right column
				if messageBits.Get(rowOffset + rowSize*4 + columnOffset + k) {
					matrix.Set(alignmentMap[high-k], alignmentMap[high-j])
				}
				// top row
				if messageBits.Get(rowOffset + rowSize*6 + columnOffset + k) {
					matrix.Set(alignmentMap[high-j], alignmentMap[low+k])
				}
			}
		}
		rowOffset += rowSize * 8
	}

	// draw mode message
	drawModeMessage(matrix, compact, matrixSize, modeMessage)

	// draw alignment marks (bulls eye drawn last so it overwrites any overlapping data)
	if compact {
		drawBullsEye(matrix, matrixSize/2, 5)
	} else {
		drawBullsEye(matrix, matrixSize/2, 7)
		for i, j := 0, 0; i < baseMatrixSize/2-1; i, j = i+15, j+16 {
			for k := matrixSize / 2 & 1; k < matrixSize; k += 2 {
				matrix.Set(matrixSize/2-j, k)
				matrix.Set(matrixSize/2+j, k)
				matrix.Set(k, matrixSize/2-j)
				matrix.Set(k, matrixSize/2+j)
			}
		}
	}

	aztec := NewAztecCode()
	aztec.SetCompact(compact)
	aztec.SetSize(matrixSize)
	aztec.SetLayers(layers)
	aztec.SetCodeWords(messageSizeInWords)
	aztec.SetMatrix(matrix)
	return aztec, nil
}

func drawBullsEye(matrix *gozxing.BitMatrix, center, size int) {
	for i := 0; i < size; i += 2 {
		for j := center - i; j <= center+i; j++ {
			matrix.Set(j, center-i)
			matrix.Set(j, center+i)
			matrix.Set(center-i, j)
			matrix.Set(center+i, j)
		}
	}
	matrix.Set(center-size, center-size)
	matrix.Set(center-size+1, center-size)
	matrix.Set(center-size, center-size+1)
	matrix.Set(center+size, center-size)
	matrix.Set(center+size, center-size+1)
	matrix.Set(center+size, center+size-1)
}

func generateModeMessage(compact bool, layers, messageSizeInWords int) *gozxing.BitArray {
	modeMessage := gozxing.NewEmptyBitArray()
	if compact {
		_ = modeMessage.AppendBits(layers-1, 2)
		_ = modeMessage.AppendBits(messageSizeInWords-1, 6)
		modeMessage = generateCheckWords(modeMessage, 28, 4)
	} else {
		_ = modeMessage.AppendBits(layers-1, 5)
		_ = modeMessage.AppendBits(messageSizeInWords-1, 11)
		modeMessage = generateCheckWords(modeMessage, 40, 4)
	}
	return modeMessage
}

func drawModeMessage(matrix *gozxing.BitMatrix, compact bool, matrixSize int, modeMessage *gozxing.BitArray) {
	center := matrixSize / 2
	if compact {
		for i := 0; i < 7; i++ {
			offset := center - 3 + i
			if modeMessage.Get(i) {
				matrix.Set(offset, center-5)
			}
			if modeMessage.Get(i + 7) {
				matrix.Set(center+5, offset)
			}
			if modeMessage.Get(20 - i) {
				matrix.Set(offset, center+5)
			}
			if modeMessage.Get(27 - i) {
				matrix.Set(center-5, offset)
			}
		}
	} else {
		for i := 0; i < 10; i++ {
			offset := center - 5 + i + i/5
			if modeMessage.Get(i) {
				matrix.Set(offset, center-7)
			}
			if modeMessage.Get(i + 10) {
				matrix.Set(center+7, offset)
			}
			if modeMessage.Get(29 - i) {
				matrix.Set(offset, center+7)
			}
			if modeMessage.Get(39 - i) {
				matrix.Set(center-7, offset)
			}
		}
	}
}

func generateCheckWords(bitArray *gozxing.BitArray, totalBits, wordSize int) *gozxing.BitArray {
	// bitArray is guaranteed to be a multiple of the wordSize, so no padding needed
	messageSizeInWords := bitArray.GetSize() / wordSize
	rs := reedsolomon.NewReedSolomonEncoder(getGF(wordSize))
	totalWords := totalBits / wordSize
	messageWords := bitsToWords(bitArray, wordSize, totalWords)
	err := rs.Encode(messageWords, totalWords-messageSizeInWords)
	if err != nil {
		panic(fmt.Sprintf("Reed-Solomon encoding error: %v", err))
	}
	startPad := totalBits % wordSize
	messageBits := gozxing.NewEmptyBitArray()
	_ = messageBits.AppendBits(0, startPad)
	for _, messageWord := range messageWords {
		_ = messageBits.AppendBits(messageWord, wordSize)
	}
	return messageBits
}

func bitsToWords(stuffedBits *gozxing.BitArray, wordSize, totalWords int) []int {
	message := make([]int, totalWords)
	n := stuffedBits.GetSize() / wordSize
	for i := 0; i < n; i++ {
		value := 0
		for j := 0; j < wordSize; j++ {
			if stuffedBits.Get(i*wordSize + j) {
				value |= 1 << uint(wordSize-j-1)
			}
		}
		message[i] = value
	}
	return message
}

func getGF(wordSize int) *reedsolomon.GenericGF {
	switch wordSize {
	case 4:
		return reedsolomon.GenericGF_AZTEC_PARAM
	case 6:
		return reedsolomon.GenericGF_AZTEC_DATA_6
	case 8:
		return reedsolomon.GenericGF_AZTEC_DATA_8
	case 10:
		return reedsolomon.GenericGF_AZTEC_DATA_10
	case 12:
		return reedsolomon.GenericGF_AZTEC_DATA_12
	default:
		panic(fmt.Sprintf("Unsupported word size %d", wordSize))
	}
}

func stuffBits(bits *gozxing.BitArray, wordSize int) *gozxing.BitArray {
	out := gozxing.NewEmptyBitArray()

	n := bits.GetSize()
	mask := (1 << uint(wordSize)) - 2
	for i := 0; i < n; i += wordSize {
		word := 0
		for j := 0; j < wordSize; j++ {
			if i+j >= n || bits.Get(i+j) {
				word |= 1 << uint(wordSize-1-j)
			}
		}
		if (word & mask) == mask {
			_ = out.AppendBits(word&mask, wordSize)
			i--
		} else if (word & mask) == 0 {
			_ = out.AppendBits(word|1, wordSize)
			i--
		} else {
			_ = out.AppendBits(word, wordSize)
		}
	}
	return out
}

func calculateTotalBitsInLayer(layers int, compact bool) int {
	base := 88
	if !compact {
		base = 112
	}
	return (base + 16*layers) * layers
}

