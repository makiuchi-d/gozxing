package decoder

import (
	"testing"
)

var (
	MASKED_TEST_FORMAT_INFO   = uint(0x2BED)
	UNMASKED_TEST_FORMAT_INFO = uint(MASKED_TEST_FORMAT_INFO ^ 0x5412)
)

func testFormatInformation_NumBitsDiffering(t testing.TB, a, b uint, expect int) {
	t.Helper()
	if r := FormatInformation_NumBitsDiffering(a, b); r != expect {
		t.Fatalf("numBitsDiffering(%v,%v) = %v, expect %v", a, b, r, expect)
	}
}

func TestFormatInformation_NumBitsDiffering(t *testing.T) {
	testFormatInformation_NumBitsDiffering(t, 1, 1, 0)
	testFormatInformation_NumBitsDiffering(t, 0, 2, 1)
	testFormatInformation_NumBitsDiffering(t, 1, 2, 2)
	testFormatInformation_NumBitsDiffering(t, 0xffffffff, 0, 32)
}

func TestFormatInformation_Decode(t *testing.T) {
	f1 := FormatInformation_DecodeFormatInformation(MASKED_TEST_FORMAT_INFO, MASKED_TEST_FORMAT_INFO)
	if f1 == nil {
		t.Fatalf("FormatInformation is nil")
	}
	if r := f1.GetDataMask(); r != 0x07 {
		t.Fatalf("DataMask is %v, expect 7", r)
	}
	if r := f1.GetErrorCorrectionLevel(); r != ErrorCorrectionLevel_Q {
		t.Fatalf("ErrorCorrectionLevel is %v, expect Q", r)
	}

	f2 := FormatInformation_DecodeFormatInformation(UNMASKED_TEST_FORMAT_INFO, MASKED_TEST_FORMAT_INFO)

	if *f1 != *f2 {
		t.Fatalf("not same f1 and f2, f1=%v, f2=%v", *f1, *f2)
	}
}

func TestFormatInformation_DecodeWithBitDifference(t *testing.T) {
	expected := FormatInformation_DecodeFormatInformation(MASKED_TEST_FORMAT_INFO, MASKED_TEST_FORMAT_INFO)

	f := FormatInformation_DecodeFormatInformation(
		MASKED_TEST_FORMAT_INFO^0x01, MASKED_TEST_FORMAT_INFO^0x01)
	if *f != *expected {
		t.Fatalf("not same with masked 0x01, %v, expect %v", *f, *expected)
	}

	f = FormatInformation_DecodeFormatInformation(
		MASKED_TEST_FORMAT_INFO^0x03, MASKED_TEST_FORMAT_INFO^0x03)
	if *f != *expected {
		t.Fatalf("not same with masked 0x03, %v, expect %v", *f, *expected)
	}

	f = FormatInformation_DecodeFormatInformation(
		MASKED_TEST_FORMAT_INFO^0x07, MASKED_TEST_FORMAT_INFO^0x07)
	if *f != *expected {
		t.Fatalf("not same with masked 0x07, %v, expect %v", *f, *expected)
	}

	f = FormatInformation_DecodeFormatInformation(
		MASKED_TEST_FORMAT_INFO^0x0f, MASKED_TEST_FORMAT_INFO^0x0f)
	if f != nil {
		t.Fatalf("Decode with invalid formatInfo must return nil")
	}
}

func TestFormatInformation_DecodeWithMisread(t *testing.T) {
	expected := FormatInformation_DecodeFormatInformation(
		MASKED_TEST_FORMAT_INFO, MASKED_TEST_FORMAT_INFO)
	f := FormatInformation_DecodeFormatInformation(
		MASKED_TEST_FORMAT_INFO^0x03, MASKED_TEST_FORMAT_INFO^0x0F)

	if *f != *expected {
		t.Fatalf("not same with missread, %v, expect %v", *f, *expected)
	}

}
