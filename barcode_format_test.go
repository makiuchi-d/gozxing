package gozxing

import (
	"testing"
)

func testBarcodeFormatString(t testing.TB, format BarcodeFormat, expect string) {
	t.Helper()
	str := format.String()
	if str != expect {
		t.Fatalf("String = \"%v\", expect \"%v\"", str, expect)
	}
}

func TestBarcodeFormatStringer(t *testing.T) {
	testBarcodeFormatString(t, BarcodeFormat_AZTEC, "AZTEC")
	testBarcodeFormatString(t, BarcodeFormat_CODABAR, "CODABAR")
	testBarcodeFormatString(t, BarcodeFormat_CODE_39, "CODE_39")
	testBarcodeFormatString(t, BarcodeFormat_CODE_93, "CODE_93")
	testBarcodeFormatString(t, BarcodeFormat_CODE_128, "CODE_128")
	testBarcodeFormatString(t, BarcodeFormat_DATA_MATRIX, "DATA_MATRIX")
	testBarcodeFormatString(t, BarcodeFormat_EAN_8, "EAN_8")
	testBarcodeFormatString(t, BarcodeFormat_EAN_13, "EAN_13")
	testBarcodeFormatString(t, BarcodeFormat_ITF, "ITF")
	testBarcodeFormatString(t, BarcodeFormat_MAXICODE, "MAXICODE")
	testBarcodeFormatString(t, BarcodeFormat_PDF_417, "PDF_417")
	testBarcodeFormatString(t, BarcodeFormat_QR_CODE, "QR_CODE")
	testBarcodeFormatString(t, BarcodeFormat_RSS_14, "RSS_14")
	testBarcodeFormatString(t, BarcodeFormat_RSS_EXPANDED, "RSS_EXPANDED")
	testBarcodeFormatString(t, BarcodeFormat_UPC_A, "UPC_A")
	testBarcodeFormatString(t, BarcodeFormat_UPC_E, "UPC_E")
	testBarcodeFormatString(t, BarcodeFormat_UPC_EAN_EXTENSION, "UPC_EAN_EXTENSION")

	testBarcodeFormatString(t, -1, "unknown format")
}

func TestBarcodeFormats_Contains(t *testing.T) {
	formats := BarcodeFormats{
		BarcodeFormat_AZTEC,
		BarcodeFormat_CODABAR,
	}

	if formats.Contains(BarcodeFormat_CODE_39) {
		t.Fatalf("CODE_39 is not contained")
	}

	if !formats.Contains(BarcodeFormat_CODABAR) {
		t.Fatalf("CODE_CODABAR is contained")
	}
}
