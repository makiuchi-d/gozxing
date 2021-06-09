package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestEan8Reader_getBarcodeFormat(t *testing.T) {
	reader := NewEAN8Reader().(*ean8Reader)
	f := reader.getBarcodeFormat()
	if f != gozxing.BarcodeFormat_EAN_8 {
		t.Fatalf("getBarcodeFormat = %v, expect EAN_8", f)
	}
}

func TestEan8Reader_decodeMiddle(t *testing.T) {
	reader := NewEAN8Reader().(*ean8Reader)
	result := make([]byte, 0, 8)
	row := gozxing.NewBitArray(70)
	startRange := []int{3, 6}

	_, _, e := reader.decodeMiddle(row, startRange, result[:0])
	if e == nil {
		t.Fatalf("decodeMiddle must be error")
	}

	// [6-12] 1(L): 0011001
	row.Set(8)
	row.Set(9)
	row.Set(12)
	// [13-19] 2(L): 0010011
	row.Set(15)
	row.Set(18)
	row.Set(19)
	// [20-26] 3(L): 0111101
	row.Set(21)
	row.Set(22)
	row.Set(23)
	row.Set(24)
	row.Set(26)
	// [27-33] 4(L) 0100011
	row.Set(28)
	row.Set(32)
	row.Set(33)
	_, _, e = reader.decodeMiddle(row, startRange, result[:0])
	if e == nil {
		t.Fatalf("decodeMiddle must be error")
	}

	// center guard
	row.Set(35)
	row.Set(37)
	row.Set(39)
	_, _, e = reader.decodeMiddle(row, startRange, result[:0])
	if e == nil {
		t.Fatalf("decodeMiddle must be error")
	}

	// [39-45] 5(R): 1001110
	row.Set(42)
	row.Set(43)
	row.Set(44)
	// [46-52] 6(R): 1010000
	row.Set(46)
	row.Set(48)
	// [53-59] 7(R): 1000100
	row.Set(53)
	row.Set(57)
	// [60-66] 8(R): 1001000
	row.Set(60)
	row.Set(63)
	// end
	row.Set(67)
	offset, result, e := reader.decodeMiddle(row, startRange, result[:0])
	if e != nil {
		t.Fatalf("decodeMiddle returns error, %v", e)
	}
	if offset != 67 {
		t.Fatalf("decodeMiddle offset = %v, expect 67", offset)
	}
	if s := string(result); s != "12345678" {
		t.Fatalf("decodeMiddle result = \"%v\", expect \"12345678\"", s)
	}
}

func TestEAN8Reader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/ean8-1/
	reader := NewEAN8Reader()
	format := gozxing.BarcodeFormat_EAN_8

	tests := []struct {
		file     string
		wants    string
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		{
			"testdata/ean8/1.png", "48512343",
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]E4",
			},
		},
		{"testdata/ean8/2.png", "12345670", nil},
		{"testdata/ean8/3.png", "12345670", nil},
		{"testdata/ean8/4.png", "67678983", nil},
		{"testdata/ean8/5.png", "80674313", nil},
		{"testdata/ean8/6.png", "59001270", nil},
		{"testdata/ean8/7.png", "50487066", nil},
		{"testdata/ean8/8.png", "55123457", nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, nil, test.metadata)
	}
}
