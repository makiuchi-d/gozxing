package gozxing

import (
	"testing"
)

func testDecodeHintType_String(t testing.TB, d DecodeHintType, e string) {
	t.Helper()
	if s := d.String(); s != e {
		t.Fatalf("DecodeHintType(%d) stringified \"%s\", expect \"%s\"", d, s, e)
	}
}

func TestDecodeHintType_String(t *testing.T) {
	testDecodeHintType_String(t, DecodeHintType_OTHER, "OTHER")
	testDecodeHintType_String(t, DecodeHintType_PURE_BARCODE, "PURE_BARCODE")
	testDecodeHintType_String(t, DecodeHintType_POSSIBLE_FORMATS, "POSSIBLE_FORMATS")
	testDecodeHintType_String(t, DecodeHintType_TRY_HARDER, "TRY_HARDER")
	testDecodeHintType_String(t, DecodeHintType_CHARACTER_SET, "CHARACTER_SET")
	testDecodeHintType_String(t, DecodeHintType_ALLOWED_LENGTHS, "ALLOWED_LENGTHS")
	testDecodeHintType_String(t, DecodeHintType_ASSUME_CODE_39_CHECK_DIGIT, "ASSUME_CODE_39_CHECK_DIGIT")
	testDecodeHintType_String(t, DecodeHintType_ASSUME_GS1, "ASSUME_GS1")
	testDecodeHintType_String(t, DecodeHintType_RETURN_CODABAR_START_END, "RETURN_CODABAR_START_END")
	testDecodeHintType_String(t, DecodeHintType_NEED_RESULT_POINT_CALLBACK, "NEED_RESULT_POINT_CALLBACK")
	testDecodeHintType_String(t, DecodeHintType_ALLOWED_EAN_EXTENSIONS, "ALLOWED_EAN_EXTENSIONS")
	testDecodeHintType_String(t, DecodeHintType_ALSO_INVERTED, "ALSO_INVERTED")
	testDecodeHintType_String(t, DecodeHintType(-1), "Unknown DecodeHintType")
}
