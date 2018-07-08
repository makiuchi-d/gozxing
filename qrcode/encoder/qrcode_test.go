package encoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing/qrcode/decoder"
)

func TestQRCode(t *testing.T) {
	qr := NewQRCode()

	str := "<<\n" +
		" mode: \n" +
		" ecLevel: M\n" +
		" version: \n" +
		" maskPattern: -1\n" +
		" matrix: nil\n" +
		">>\n"
	if r := qr.String(); r != str {
		t.Fatalf("String:%vexpect:%v", r, str)
	}

	if r := qr.GetMode(); r != nil {
		t.Fatalf("GetMode must be nil, %v", r)
	}
	qr.SetMode(decoder.Mode_ALPHANUMERIC)
	if r := qr.GetMode(); r != decoder.Mode_ALPHANUMERIC {
		t.Fatalf("GetMode = %v, expect %v", r, decoder.Mode_ALPHANUMERIC)
	}

	// default of ErrorCorrectionLevel (0x00) is "M"
	if r := qr.GetECLevel(); r != decoder.ErrorCorrectionLevel_M {
		t.Fatalf("GetECLevel = %v, expect %v", r, decoder.ErrorCorrectionLevel_M)
	}
	qr.SetECLevel(decoder.ErrorCorrectionLevel_H)
	if r := qr.GetECLevel(); r != decoder.ErrorCorrectionLevel_H {
		t.Fatalf("GetECLevel = %v, expect %v", r, decoder.ErrorCorrectionLevel_H)
	}

	if r := qr.GetVersion(); r != nil {
		t.Fatalf("GetVersion must be nil, %v", r)
	}
	version, _ := decoder.Version_GetVersionForNumber(10)
	qr.SetVersion(version)
	if r := qr.GetVersion(); r != version {
		t.Fatalf("GetVersion = %v, expect %v", r, version)
	}

	if r := qr.GetMaskPattern(); QRCode_IsValidMaskPattern(r) {
		t.Fatalf("GetMaskPattern must be invalid mask pattern, %v", r)
	}
	qr.SetMaskPattern(4)
	if r := qr.GetMaskPattern(); r != 4 {
		t.Fatalf("GetMaskPattern = %v, expect %v", r, 4)
	}

	if r := qr.GetMatrix(); r != nil {
		t.Fatalf("GetMatrix must be nil, %v", r)
	}
	matrix := NewByteMatrix(3, 3)
	qr.SetMatrix(matrix)
	if r := qr.GetMatrix(); r != matrix {
		t.Fatalf("GetMatrix = %v, expect %v", r, matrix)
	}

	str = "<<\n" +
		" mode: ALPHANUMERIC\n" +
		" ecLevel: H\n" +
		" version: 10\n" +
		" maskPattern: 4\n" +
		" matrix:\n" +
		" 0 0 0\n" +
		" 0 0 0\n" +
		" 0 0 0\n" +
		">>\n"
	if r := qr.String(); r != str {
		t.Fatalf("String:%vexpect:%v", r, str)
	}
}

func TestQRCode_IsValidMaskPattern(t *testing.T) {
	if QRCode_IsValidMaskPattern(-1) {
		t.Fatalf("IsValidMaskPattern(-1) must be false")
	}
	if !QRCode_IsValidMaskPattern(0) {
		t.Fatalf("IsValidMaskPattern(0) must be true")
	}
	if !QRCode_IsValidMaskPattern(QRCode_NUM_MASK_PATERNS - 1) {
		t.Fatalf("IsValidMaskPattern(QRCode_NUM_MASK_PATERNS - 1) must be true")
	}
	if QRCode_IsValidMaskPattern(QRCode_NUM_MASK_PATERNS) {
		t.Fatalf("IsValidMaskPattern(QRCode_NUM_MASK_PATERNS) must be false")
	}
}
