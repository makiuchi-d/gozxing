package gozxing

import (
	"testing"
)

func testDecoderHintType_String(t *testing.T, d DecodeHintType, e string) {
	if s := d.String(); s != e {
		t.Fatalf("DecodeHintType(%d) stringified \"%s\", expect \"%s\"", d, s, e)
	}
}

func TestDecoderHintType_String(t *testing.T) {
	testDecoderHintType_String(t, DecodeHintType_OTHER, "OTHER")
	testDecoderHintType_String(t, DecodeHintType_PURE_BARCODE, "PURE_BARCODE")
	testDecoderHintType_String(t, DecodeHintType_POSSIBLE_FORMATS, "POSSIBLE_FORMATS")
	testDecoderHintType_String(t, DecodeHintType_TRY_HARDER, "TRY_HARDER")
	testDecoderHintType_String(t, DecodeHintType_CHARACTER_SET, "CHARACTER_SET")
	testDecoderHintType_String(t, DecodeHintType_ALLOWED_LENGTHS, "ALLOWED_LENGTHS")
	testDecoderHintType_String(t, DecodeHintType_ASSUME_CODE_39_CHECK_DIGIT, "ASSUME_CODE_39_CHECK_DIGIT")
	testDecoderHintType_String(t, DecodeHintType_ASSUME_GS1, "ASSUME_GS1")
	testDecoderHintType_String(t, DecodeHintType_RETURN_CODABAR_START_END, "RETURN_CODABAR_START_END")
	testDecoderHintType_String(t, DecodeHintType_NEED_RESULT_POINT_CALLBACK, "NEED_RESULT_POINT_CALLBACK")
	testDecoderHintType_String(t, DecodeHintType_ALLOWED_EAN_EXTENSIONS, "ALLOWED_EAN_EXTENSIONS")
	testDecoderHintType_String(t, DecodeHintType(-1), "Unknown DecodeHintType")
}
