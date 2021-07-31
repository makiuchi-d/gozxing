package gozxing

import (
	"testing"
)

func testEncodeHintType_String(t testing.TB, h EncodeHintType, e string) {
	t.Helper()
	if s := h.String(); s != e {
		t.Fatalf("DecodeHintType(%d) stringified \"%s\", expect \"%s\"", h, s, e)
	}
}

func TestEncodeHintType_String(t *testing.T) {
	testEncodeHintType_String(t, EncodeHintType_ERROR_CORRECTION, "ERROR_CORRECTION")
	testEncodeHintType_String(t, EncodeHintType_CHARACTER_SET, "CHARACTER_SET")
	testEncodeHintType_String(t, EncodeHintType_DATA_MATRIX_SHAPE, "DATA_MATRIX_SHAPE")
	testEncodeHintType_String(t, EncodeHintType_MIN_SIZE, "MIN_SIZE")
	testEncodeHintType_String(t, EncodeHintType_MAX_SIZE, "MAX_SIZE")
	testEncodeHintType_String(t, EncodeHintType_MARGIN, "MARGIN")
	testEncodeHintType_String(t, EncodeHintType_PDF417_COMPACT, "PDF417_COMPACT")
	testEncodeHintType_String(t, EncodeHintType_PDF417_COMPACTION, "PDF417_COMPACTION")
	testEncodeHintType_String(t, EncodeHintType_PDF417_DIMENSIONS, "PDF417_DIMENSIONS")
	testEncodeHintType_String(t, EncodeHintType_AZTEC_LAYERS, "AZTEC_LAYERS")
	testEncodeHintType_String(t, EncodeHintType_QR_VERSION, "QR_VERSION")
	testEncodeHintType_String(t, EncodeHintType_QR_MASK_PATTERN, "QR_MASK_PATTERN")
	testEncodeHintType_String(t, EncodeHintType_GS1_FORMAT, "GS1_FORMAT")
	testEncodeHintType_String(t, EncodeHintType_FORCE_CODE_SET, "FORCE_CODE_SET")
	testEncodeHintType_String(t, EncodeHintType(-1), "")
}
