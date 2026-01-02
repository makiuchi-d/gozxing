package aztec

import (
	"strconv"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/aztec/encoder"
)

// AztecWriter Renders an Aztec code as a BitMatrix.
type AztecWriter struct{}

func NewAztecWriter() gozxing.Writer {
	return &AztecWriter{}
}

func (this *AztecWriter) EncodeWithoutHint(
	contents string, format gozxing.BarcodeFormat, width, height int) (*gozxing.BitMatrix, error) {
	return this.Encode(contents, format, width, height, nil)
}

func (this *AztecWriter) Encode(contents string, format gozxing.BarcodeFormat,
	width, height int, hints map[gozxing.EncodeHintType]interface{}) (*gozxing.BitMatrix, error) {

	if len(contents) == 0 {
		return nil, gozxing.NewWriterException("IllegalArgumentException: Found empty contents")
	}

	if format != gozxing.BarcodeFormat_AZTEC {
		return nil, gozxing.NewWriterException(
			"IllegalArgumentException: Can only encode AZTEC, but got %v", format)
	}

	if width < 0 || height < 0 {
		return nil, gozxing.NewWriterException(
			"IllegalArgumentException: Requested dimensions are too small: %vx%v", width, height)
	}

	var charset encoding.Encoding
	eccPercent := encoder.DEFAULT_EC_PERCENT
	layers := encoder.DEFAULT_AZTEC_LAYERS

	if hints != nil {
		if cs, ok := hints[gozxing.EncodeHintType_CHARACTER_SET]; ok {
			if str, ok := cs.(string); ok {
				enc, err := ianaindex.IANA.Encoding(str)
				if err == nil {
					charset = enc
				}
			} else if enc, ok := cs.(encoding.Encoding); ok {
				charset = enc
			}
		}
		if ec, ok := hints[gozxing.EncodeHintType_ERROR_CORRECTION]; ok {
			if ecInt, ok := ec.(int); ok {
				eccPercent = ecInt
			} else if str, ok := ec.(string); ok {
				ecInt, err := strconv.Atoi(str)
				if err != nil {
					return nil, gozxing.NewWriterException("EncodeHintType_ERROR_CORRECTION = \"%v\": %w", ec, err)
				}
				eccPercent = ecInt
			} else {
				return nil, gozxing.NewWriterException(
					"IllegalArgumentException: EncodeHintType_ERROR_CORRECTION %v", ec)
			}
		}
		if l, ok := hints[gozxing.EncodeHintType_AZTEC_LAYERS]; ok {
			if layersInt, ok := l.(int); ok {
				layers = layersInt
			} else if str, ok := l.(string); ok {
				layersInt, err := strconv.Atoi(str)
				if err != nil {
					return nil, gozxing.NewWriterException("EncodeHintType_AZTEC_LAYERS = \"%v\": %w", l, err)
				}
				layers = layersInt
			} else {
				return nil, gozxing.NewWriterException(
					"IllegalArgumentException: EncodeHintType_AZTEC_LAYERS %v", l)
			}
		}
	}

	aztec, err := encoder.EncodeStringWithParams(contents, eccPercent, layers, charset)
	if err != nil {
		return nil, err
	}
	return renderResult(aztec, width, height)
}

func renderResult(code *encoder.AztecCode, width, height int) (*gozxing.BitMatrix, error) {
	input := code.GetMatrix()
	if input == nil {
		return nil, gozxing.NewWriterException("IllegalStateException")
	}
	inputWidth := input.GetWidth()
	inputHeight := input.GetHeight()
	outputWidth := width
	if outputWidth < inputWidth {
		outputWidth = inputWidth
	}
	outputHeight := height
	if outputHeight < inputHeight {
		outputHeight = inputHeight
	}

	multiple := outputWidth / inputWidth
	if h := outputHeight / inputHeight; multiple > h {
		multiple = h
	}
	leftPadding := (outputWidth - (inputWidth * multiple)) / 2
	topPadding := (outputHeight - (inputHeight * multiple)) / 2

	output, err := gozxing.NewBitMatrix(outputWidth, outputHeight)
	if err != nil {
		return nil, gozxing.WrapWriterException(err)
	}

	for inputY, outputY := 0, topPadding; inputY < inputHeight; inputY, outputY = inputY+1, outputY+multiple {
		// Write the contents of this row of the barcode
		for inputX, outputX := 0, leftPadding; inputX < inputWidth; inputX, outputX = inputX+1, outputX+multiple {
			if input.Get(inputX, inputY) {
				output.SetRegion(outputX, outputY, multiple, multiple)
			}
		}
	}
	return output, nil
}
