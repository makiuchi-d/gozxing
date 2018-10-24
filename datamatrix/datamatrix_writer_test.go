package datamatrix

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix/encoder"
	qrencoder "github.com/makiuchi-d/gozxing/qrcode/encoder"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestConvertByteMatrixToBitMatrix(t *testing.T) {
	src, _ := gozxing.ParseStringToBitMatrix(""+
		"##  ##  ##  ##  ##  ##  ##  ##  \n"+
		"##  ####  ####  ##  ####      ##\n"+
		"####      ####  ####            \n"+
		"######  ####  ##    ####      ##\n"+
		"####                ##########  \n"+
		"######      ##    ##        ####\n"+
		"####    ####          ##        \n"+
		"######  ####  ##    ##        ##\n"+
		"##    ##  ##############        \n"+
		"##  ####  ##    ######  ####  ##\n"+
		"######  ##    ##      ######    \n"+
		"##  ####  ######        ####  ##\n"+
		"##    ##  ######  ##  ######    \n"+
		"########    ##  ##        ##  ##\n"+
		"############  ##  ######    ##  \n"+
		"################################\n", "##", "  ")
	bm := qrencoder.NewByteMatrix(src.GetWidth(), src.GetHeight())
	for j := 0; j < bm.GetHeight(); j++ {
		for i := 0; i < bm.GetWidth(); i++ {
			bm.SetBool(i, j, src.Get(i, j))
		}
	}

	b := convertByteMatrixToBitMatrix(bm, 48, 35)
	expect := "" +
		"                                                                                                \n" +
		"                X X     X X     X X     X X     X X     X X     X X     X X                     \n" +
		"                X X     X X     X X     X X     X X     X X     X X     X X                     \n" +
		"                X X     X X X X     X X X X     X X     X X X X             X X                 \n" +
		"                X X     X X X X     X X X X     X X     X X X X             X X                 \n" +
		"                X X X X             X X X X     X X X X                                         \n" +
		"                X X X X             X X X X     X X X X                                         \n" +
		"                X X X X X X     X X X X     X X         X X X X             X X                 \n" +
		"                X X X X X X     X X X X     X X         X X X X             X X                 \n" +
		"                X X X X                                 X X X X X X X X X X                     \n" +
		"                X X X X                                 X X X X X X X X X X                     \n" +
		"                X X X X X X             X X         X X                 X X X X                 \n" +
		"                X X X X X X             X X         X X                 X X X X                 \n" +
		"                X X X X         X X X X                     X X                                 \n" +
		"                X X X X         X X X X                     X X                                 \n" +
		"                X X X X X X     X X X X     X X         X X                 X X                 \n" +
		"                X X X X X X     X X X X     X X         X X                 X X                 \n" +
		"                X X         X X     X X X X X X X X X X X X X X                                 \n" +
		"                X X         X X     X X X X X X X X X X X X X X                                 \n" +
		"                X X     X X X X     X X         X X X X X X     X X X X     X X                 \n" +
		"                X X     X X X X     X X         X X X X X X     X X X X     X X                 \n" +
		"                X X X X X X     X X         X X             X X X X X X                         \n" +
		"                X X X X X X     X X         X X             X X X X X X                         \n" +
		"                X X     X X X X     X X X X X X                 X X X X     X X                 \n" +
		"                X X     X X X X     X X X X X X                 X X X X     X X                 \n" +
		"                X X         X X     X X X X X X     X X     X X X X X X                         \n" +
		"                X X         X X     X X X X X X     X X     X X X X X X                         \n" +
		"                X X X X X X X X         X X     X X                 X X     X X                 \n" +
		"                X X X X X X X X         X X     X X                 X X     X X                 \n" +
		"                X X X X X X X X X X X X     X X     X X X X X X         X X                     \n" +
		"                X X X X X X X X X X X X     X X     X X X X X X         X X                     \n" +
		"                X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X                 \n" +
		"                X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X X                 \n" +
		"                                                                                                \n" +
		"                                                                                                \n"
	if str := b.String(); str != expect {
		t.Fatalf("convertByteMatrixToBitMatrix:\n%vexpect:\n%v", str, expect)
	}

	b = convertByteMatrixToBitMatrix(bm, 10, 10)
	expect = "" +
		"X   X   X   X   X   X   X   X   \n" +
		"X   X X   X X   X   X X       X \n" +
		"X X       X X   X X             \n" +
		"X X X   X X   X     X X       X \n" +
		"X X                 X X X X X   \n" +
		"X X X       X     X         X X \n" +
		"X X     X X           X         \n" +
		"X X X   X X   X     X         X \n" +
		"X     X   X X X X X X X         \n" +
		"X   X X   X     X X X   X X   X \n" +
		"X X X   X     X       X X X     \n" +
		"X   X X   X X X         X X   X \n" +
		"X     X   X X X   X   X X X     \n" +
		"X X X X     X   X         X   X \n" +
		"X X X X X X   X   X X X     X   \n" +
		"X X X X X X X X X X X X X X X X \n"
	if str := b.String(); str != expect {
		t.Fatalf("convertByteMatrixToBitMatrix:\n%vexpect:\n%v", str, expect)
	}
}

func TestEncodeLowLevel(t *testing.T) {
	cw := make([]byte, 24)
	for i := 0; i < len(cw); i++ {
		cw[i] = 0xaa
	}

	symbol := encoder.NewSymbolInfo(false, 5, 7, 10, 10, 1)
	placement := encoder.NewDefaultPlacement(cw[:12], 10, 10)
	placement.Place()

	bm := encodeLowLevel(placement, symbol, 14, 16)
	expect := "" +
		"                            \n" +
		"                            \n" +
		"  X   X   X   X   X   X     \n" +
		"  X X     X   X   X X   X   \n" +
		"  X X   X X     X   X       \n" +
		"  X   X   X   X X     X X   \n" +
		"  X X X     X   X   X X     \n" +
		"  X   X   X X     X   X X   \n" +
		"  X     X   X   X X         \n" +
		"  X   X X     X   X   X X   \n" +
		"  X X   X   X X     X       \n" +
		"  X X     X   X   X X   X   \n" +
		"  X X   X X     X     X     \n" +
		"  X X X X X X X X X X X X   \n" +
		"                            \n" +
		"                            \n"
	if str := bm.String(); str != expect {
		t.Fatalf("encodeLowLevel:\n%v\nexpect:\n%v", str, expect)
	}
}

func TestDataMatrixWriter_Encode(t *testing.T) {
	writer := NewDataMatrixWriter()

	_, e := writer.EncodeWithoutHint("", gozxing.BarcodeFormat_DATA_MATRIX, 10, 10)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	_, e = writer.EncodeWithoutHint("123456", gozxing.BarcodeFormat_QR_CODE, 10, 10)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	_, e = writer.EncodeWithoutHint("123456", gozxing.BarcodeFormat_DATA_MATRIX, -1, 10)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	_, e = writer.EncodeWithoutHint("123456", gozxing.BarcodeFormat_DATA_MATRIX, 10, -1)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	hint := make(map[gozxing.EncodeHintType]interface{})

	hint[gozxing.EncodeHintType_MAX_SIZE], _ = gozxing.NewDimension(5, 5)
	_, e = writer.Encode("123456", gozxing.BarcodeFormat_DATA_MATRIX, 10, 10, hint)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	hint[gozxing.EncodeHintType_DATA_MATRIX_SHAPE] = encoder.SymbolShapeHint_FORCE_RECTANGLE
	hint[gozxing.EncodeHintType_MIN_SIZE], _ = gozxing.NewDimension(5, 5)
	hint[gozxing.EncodeHintType_MAX_SIZE], _ = gozxing.NewDimension(20, 20)
	b, e := writer.Encode("123456", gozxing.BarcodeFormat_DATA_MATRIX, 20, 20, hint)
	expect := "" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n" +
		"  X   X   X   X   X   X   X   X   X     \n" +
		"  X X     X           X X           X   \n" +
		"  X X       X     X X   X X X X   X     \n" +
		"  X X     X X       X       X X X   X   \n" +
		"  X X X X   X X     X X X     X         \n" +
		"  X   X X X X       X       X   X X X   \n" +
		"  X         X X X X   X X   X X   X     \n" +
		"  X X X X X X X X X X X X X X X X X X   \n" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n" +
		"                                        \n"
	if e != nil {
		t.Fatalf("Encode returns error: %v", e)
	}
	if str := b.String(); str != expect {
		t.Fatalf("Encode:\n%vexpect:\n%v", str, expect)
	}

	contents := "Hello, world!"
	b, e = writer.Encode(contents, gozxing.BarcodeFormat_DATA_MATRIX, 100, 100, nil)
	if e != nil {
		t.Fatalf("Encode returns error: %v", e)
	}
	bmp := testutil.NewBinaryBitmapFromBitMatrix(b)
	reader := NewDataMatrixReader()
	result, e := reader.DecodeWithoutHints(bmp)
	if e != nil {
		t.Fatalf("Decode returns error: %v", e)
	}
	if txt := result.GetText(); txt != contents {
		t.Fatalf("result = \"%v\", expect \"%v\"", txt, contents)
	}
}
