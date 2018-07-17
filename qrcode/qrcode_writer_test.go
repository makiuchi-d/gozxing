package qrcode

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
	"github.com/makiuchi-d/gozxing/qrcode/detector"
	"github.com/makiuchi-d/gozxing/qrcode/encoder"
)

func TestQRCodeWriter_renderResult(t *testing.T) {
	code := encoder.NewQRCode()

	_, e := renderResult(code, 10, 10, 5)
	if e == nil {
		t.Fatalf("renderResult must be error")
	}

	code, _ = encoder.Encoder_encode("test", decoder.ErrorCorrectionLevel_M, nil)

	_, e = renderResult(code, 0, 1, -11)
	if e == nil {
		t.Fatalf("renderResult must be error")
	}

	matrix, e := renderResult(code, 60, 64, 4)
	if e != nil {
		t.Fatalf("renderResult returns error, %v", e)
	}
	if w, h := matrix.GetWidth(), matrix.GetHeight(); w != 60 || h != 64 {
		t.Fatalf("renderResult matrix size = %vx%v, expect 60x64", w, h)
	}
	expect := "" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                  X X X X X X X X X X X X X X             X X             X X X X X X X X X X X X X X                   \n" +
		"                  X X X X X X X X X X X X X X             X X             X X X X X X X X X X X X X X                   \n" +
		"                  X X                     X X     X X     X X     X X     X X                     X X                   \n" +
		"                  X X                     X X     X X     X X     X X     X X                     X X                   \n" +
		"                  X X     X X X X X X     X X                 X X         X X     X X X X X X     X X                   \n" +
		"                  X X     X X X X X X     X X                 X X         X X     X X X X X X     X X                   \n" +
		"                  X X     X X X X X X     X X         X X X X             X X     X X X X X X     X X                   \n" +
		"                  X X     X X X X X X     X X         X X X X             X X     X X X X X X     X X                   \n" +
		"                  X X     X X X X X X     X X     X X X X     X X X X     X X     X X X X X X     X X                   \n" +
		"                  X X     X X X X X X     X X     X X X X     X X X X     X X     X X X X X X     X X                   \n" +
		"                  X X                     X X         X X     X X         X X                     X X                   \n" +
		"                  X X                     X X         X X     X X         X X                     X X                   \n" +
		"                  X X X X X X X X X X X X X X     X X     X X     X X     X X X X X X X X X X X X X X                   \n" +
		"                  X X X X X X X X X X X X X X     X X     X X     X X     X X X X X X X X X X X X X X                   \n" +
		"                                                      X X X X                                                           \n" +
		"                                                      X X X X                                                           \n" +
		"                  X X     X X     X X     X X             X X     X X             X X         X X                       \n" +
		"                  X X     X X     X X     X X             X X     X X             X X         X X                       \n" +
		"                          X X X X X X X X     X X         X X X X     X X     X X             X X X X                   \n" +
		"                          X X X X X X X X     X X         X X X X     X X     X X             X X X X                   \n" +
		"                  X X     X X             X X X X X X X X X X X X     X X X X X X     X X X X X X X X                   \n" +
		"                  X X     X X             X X X X X X X X X X X X     X X X X X X     X X X X X X X X                   \n" +
		"                              X X X X X X         X X X X X X X X X X X X     X X X X         X X                       \n" +
		"                              X X X X X X         X X X X X X X X X X X X     X X X X         X X                       \n" +
		"                  X X     X X     X X     X X         X X X X X X     X X X X X X     X X     X X X X                   \n" +
		"                  X X     X X     X X     X X         X X X X X X     X X X X X X     X X     X X X X                   \n" +
		"                                                  X X                     X X         X X         X X                   \n" +
		"                                                  X X                     X X         X X         X X                   \n" +
		"                  X X X X X X X X X X X X X X         X X         X X             X X X X     X X X X                   \n" +
		"                  X X X X X X X X X X X X X X         X X         X X             X X X X     X X X X                   \n" +
		"                  X X                     X X             X X             X X                 X X                       \n" +
		"                  X X                     X X             X X             X X                 X X                       \n" +
		"                  X X     X X X X X X     X X     X X X X X X     X X     X X     X X X X     X X X X                   \n" +
		"                  X X     X X X X X X     X X     X X X X X X     X X     X X     X X X X     X X X X                   \n" +
		"                  X X     X X X X X X     X X             X X X X     X X     X X             X X                       \n" +
		"                  X X     X X X X X X     X X             X X X X     X X     X X             X X                       \n" +
		"                  X X     X X X X X X     X X     X X X X X X X X     X X X X X X         X X     X X                   \n" +
		"                  X X     X X X X X X     X X     X X X X X X X X     X X X X X X         X X     X X                   \n" +
		"                  X X                     X X             X X X X X X X X     X X X X X X     X X                       \n" +
		"                  X X                     X X             X X X X X X X X     X X X X X X     X X                       \n" +
		"                  X X X X X X X X X X X X X X     X X X X X X X X     X X X X X X         X X X X X X                   \n" +
		"                  X X X X X X X X X X X X X X     X X X X X X X X     X X X X X X         X X X X X X                   \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n" +
		"                                                                                                                        \n"
	if r := matrix.String(); r != expect {
		t.Fatalf("renderResult matrix:\n%v\nexpect:\n%v", r, expect)
	}
}

func TestQRCodeWriter_EncodeFailed(t *testing.T) {
	writer := NewQRCodeWriter()
	var e error

	_, e = writer.EncodeWithoutHint("", gozxing.BarcodeFormat_QR_CODE, 30, 30)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}

	_, e = writer.EncodeWithoutHint("test", gozxing.BarcodeFormat_AZTEC, 30, 30)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}

	_, e = writer.EncodeWithoutHint("test", gozxing.BarcodeFormat_QR_CODE, 30, -1)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}

	hints := map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_CHARACTER_SET: "ISO-8859-1"}
	_, e = writer.Encode("エラー", gozxing.BarcodeFormat_QR_CODE, 30, 4, hints)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}
}

func parseMatrix(matrix *gozxing.BitMatrix) (*decoder.Version, *decoder.FormatInformation, error) {
	detected, e := detector.NewDetector(matrix).Detect(nil)
	if e != nil {
		return nil, nil, e
	}
	parser, e := decoder.NewBitMatrixParser(detected.GetBits())
	if e != nil {
		return nil, nil, e
	}
	version, e := parser.ReadVersion()
	if e != nil {
		return nil, nil, e
	}
	format, e := parser.ReadFormatInformation()
	if e != nil {
		return nil, nil, e
	}

	return version, format, nil
}

func TestQRCodeWriter_EncodeWithECLevelHint(t *testing.T) {
	writer := NewQRCodeWriter()
	formatQR := gozxing.BarcodeFormat_QR_CODE
	var e error
	var matrix *gozxing.BitMatrix

	hints := make(map[gozxing.EncodeHintType]interface{})

	hints[gozxing.EncodeHintType_ERROR_CORRECTION] = 1
	_, e = writer.Encode("test", formatQR, 30, 2, hints)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}

	hints[gozxing.EncodeHintType_ERROR_CORRECTION] = "A"
	_, e = writer.Encode("test", formatQR, 30, 2, hints)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}

	hints[gozxing.EncodeHintType_ERROR_CORRECTION] = decoder.ErrorCorrectionLevel_M
	matrix, e = writer.Encode("test", formatQR, 30, 2, hints)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	_, format, e := parseMatrix(matrix)
	if e != nil {
		t.Fatalf("Encode result cannot parse, %v", e)
	}
	if r := format.GetErrorCorrectionLevel(); r != decoder.ErrorCorrectionLevel_M {
		t.Fatalf("Encoder result ECLevel = %v, expect M\n", r)
	}

	hints[gozxing.EncodeHintType_ERROR_CORRECTION] = "Q"
	matrix, e = writer.Encode("test", formatQR, 30, 2, hints)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	_, format, e = parseMatrix(matrix)
	if e != nil {
		t.Fatalf("Encode result cannot parse, %v", e)
	}
	if r := format.GetErrorCorrectionLevel(); r != decoder.ErrorCorrectionLevel_Q {
		t.Fatalf("Encoder result ECLevel = %v, expect Q\n", r)
	}
}

func TestQRCodeWriter_EncodeWithMarginHint(t *testing.T) {
	writer := NewQRCodeWriter()
	formatQR := gozxing.BarcodeFormat_QR_CODE
	var e error
	var matrix *gozxing.BitMatrix

	hints := make(map[gozxing.EncodeHintType]interface{})

	hints[gozxing.EncodeHintType_MARGIN] = nil
	_, e = writer.Encode("test", formatQR, 10, 10, hints)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}

	hints[gozxing.EncodeHintType_MARGIN] = "a"
	_, e = writer.Encode("test", formatQR, 10, 10, hints)
	if e == nil {
		t.Fatalf("Encode must be Error")
	}

	hints[gozxing.EncodeHintType_MARGIN] = 10
	matrix, e = writer.Encode("test", formatQR, 10, 60, hints)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	version, _, e := parseMatrix(matrix)
	if e != nil {
		t.Fatalf("Encode result cannot parse, %v", e)
	}
	expect := version.GetDimensionForVersion() + 10*2
	if w, h := matrix.GetWidth(), matrix.GetHeight(); w != expect || h != 60 {
		t.Fatalf("Encode result size = %vx%v, expect %vx%v", w, h, expect, 60)
	}

	hints[gozxing.EncodeHintType_MARGIN] = "20"
	matrix, e = writer.Encode("test", formatQR, 10, 60, hints)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	version, _, e = parseMatrix(matrix)
	if e != nil {
		t.Fatalf("Encode result cannot parse, %v", e)
	}
	expect = version.GetDimensionForVersion() + 20*2
	if w, h := matrix.GetWidth(), matrix.GetHeight(); w != expect || h != expect {
		t.Fatalf("Encode result size = %vx%v, expect %vx%v", w, h, expect, expect)
	}
}
