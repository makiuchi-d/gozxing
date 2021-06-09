package gozxing

import (
	"testing"
)

func testResultMetadataTypeString(t testing.TB, metadataType ResultMetadataType, expect string) {
	t.Helper()
	if r := metadataType.String(); r != expect {
		t.Fatalf("String = %v, expect %v", r, expect)
	}
}

func TestResultMetadataTypeString(t *testing.T) {
	testResultMetadataTypeString(t, ResultMetadataType_OTHER, "OTHER")
	testResultMetadataTypeString(t, ResultMetadataType_ORIENTATION, "ORIENTATION")
	testResultMetadataTypeString(t, ResultMetadataType_BYTE_SEGMENTS, "BYTE_SEGMENTS")
	testResultMetadataTypeString(t, ResultMetadataType_ERROR_CORRECTION_LEVEL, "ERROR_CORRECTION_LEVEL")
	testResultMetadataTypeString(t, ResultMetadataType_ISSUE_NUMBER, "ISSUE_NUMBER")
	testResultMetadataTypeString(t, ResultMetadataType_SUGGESTED_PRICE, "SUGGESTED_PRICE")
	testResultMetadataTypeString(t, ResultMetadataType_POSSIBLE_COUNTRY, "POSSIBLE_COUNTRY")
	testResultMetadataTypeString(t, ResultMetadataType_UPC_EAN_EXTENSION, "UPC_EAN_EXTENSION")
	testResultMetadataTypeString(t, ResultMetadataType_PDF417_EXTRA_METADATA, "PDF417_EXTRA_METADATA")
	testResultMetadataTypeString(t, ResultMetadataType_STRUCTURED_APPEND_SEQUENCE, "STRUCTURED_APPEND_SEQUENCE")
	testResultMetadataTypeString(t, ResultMetadataType_STRUCTURED_APPEND_PARITY, "STRUCTURED_APPEND_PARITY")
	testResultMetadataTypeString(t, ResultMetadataType_SYMBOLOGY_IDENTIFIER, "SYMBOLOGY_IDENTIFIER")

	testResultMetadataTypeString(t, -1, "unknown metadata type")
}
