package decoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	"github.com/makiuchi-d/gozxing/testutil"
)

func unsetRegion(image *gozxing.BitMatrix, x, y, w, h int) {
	for i := y; i < y+h; i++ {
		for j := x; j < x+w; j++ {
			image.Unset(j, i)
		}
	}
}

func TestDecoder_Decode(t *testing.T) {
	decoder := NewDecoder()
	var result *common.DecoderResult
	var e error

	bits, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	rbits := testutil.MirrorBitMatrix(bits)

	// normal qrcode
	result, e = decoder.Decode(bits, nil)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if r := result.GetText(); r != "hello" {
		t.Fatalf("Decoder result text=\"%v\", expect \"hello\"", r)
	}
	// mirrored qrcode
	result, e = decoder.Decode(rbits, nil)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if r := result.GetText(); r != "hello" {
		t.Fatalf("Decoder result text=\"%v\", expect \"hello\"", r)
	}

	bits, _ = gozxing.NewSquareBitMatrix(1)
	_, e = decoder.Decode(bits, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	bits, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	for i := 0; i < 21; i++ {
		bits.Unset(i, 8)
	}
	_, e = decoder.Decode(bits, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	bits, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	bits.SetRegion(10, 10, 10, 10)
	_, e = decoder.Decode(bits, nil)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("Decode must be ChecksumException, %T", e)
	}

	bits, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	result, e = decoder.DecodeWithoutHint(bits)
	if e != nil {
		t.Fatalf("DecodeWithoutHint returns error, %v", e)
	}
	if r := result.GetText(); r != "hello" {
		t.Fatalf("DecodeWithoutHint result text=\"%v\", expect \"hello\"", r)
	}
}

func TestDecoder_DecodeBoolMap(t *testing.T) {
	decoder := NewDecoder()
	bits, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	w := bits.GetWidth()
	h := bits.GetHeight()
	boolmap := make([][]bool, h)
	for j := 0; j < h; j++ {
		boolmap[j] = make([]bool, w)
		for i := 0; i < bits.GetWidth(); i++ {
			boolmap[j][i] = bits.Get(i, j)
		}
	}

	_, e := decoder.DecodeBoolMap([][]bool{}, nil)
	if e == nil {
		t.Fatalf("DecodeBoolMap must be error")
	}

	result, e := decoder.DecodeBoolMap(boolmap, nil)
	if e != nil {
		t.Fatalf("DecodeBoolMap returns error, %v", e)
	}
	if r := result.GetText(); r != "hello" {
		t.Fatalf("DecodeBoolMap result text=\"%v\", expect \"hello\"", r)
	}

	result, e = decoder.DecodeBoolMapWithoutHint(boolmap)
	if e != nil {
		t.Fatalf("DecodeBoolMapWithoutHint returns error, %v", e)
	}
	if r := result.GetText(); r != "hello" {
		t.Fatalf("DecodeBoolMapWithoutHint result text=\"%v\", expect \"hello\"", r)
	}
}

func TestDecoder_decode(t *testing.T) {
	var bits *gozxing.BitMatrix
	var parser *BitMatrixParser
	var e error
	decoder := NewDecoder()

	// no version bits
	bits, _ = gozxing.NewSquareBitMatrix(45)
	parser, _ = NewBitMatrixParser(bits)
	_, e = decoder.decode(parser, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decode must be FormatException, %T", e)
	}

	// invalid format info
	bits, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	for i := 0; i < 21; i++ {
		bits.Unset(i, 8)
	}
	parser, _ = NewBitMatrixParser(bits)
	_, e = decoder.decode(parser, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decode must be FormatException, %T", e)
	}

	// too many errors
	bits, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	bits.SetRegion(10, 10, 10, 10)
	parser, _ = NewBitMatrixParser(bits)
	_, e = decoder.decode(parser, nil)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("decode must be ChecksumException, %T", e)
	}

	bits, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	parser, _ = NewBitMatrixParser(bits)
	result, e := decoder.decode(parser, nil)
	if e != nil {
		t.Fatalf("decode returns error, %v", e)
	}
	if r := result.GetText(); r != "hello" {
		t.Fatalf("decoder result text=\"%v\", expect \"hello\"", r)
	}
}
