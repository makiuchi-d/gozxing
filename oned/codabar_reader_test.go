package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestCodabarReader_arrayContains(t *testing.T) {
	tests := []struct {
		key   byte
		wants bool
	}{
		{'A', true},
		{'B', true},
		{'C', true},
		{'D', true},
		{'0', false},
		{'9', false},
		{'-', false},
		{':', false},
		{'+', false},
	}
	for _, test := range tests {
		r := codabarReader_arrayContains(codabarReader_STARTEND_ENCODING, test.key)
		if r != test.wants {
			t.Fatalf("arrayContains(%v) = %v, wants %v", test.key, r, test.wants)
		}
	}
}

func TestCodabarReader_toNarrowWidePattern(t *testing.T) {
	reader := NewCodaBarReader().(*codabarReader)
	reader.counters = []int{
		3, 3, 5, 6, 2, 5, 2, // A
		9, 1, 5, 1, 6, 3, 5, // 5
		2, 1, 2, 1, 2, 1, 2, // not a pattern
		1, 2, 3,
	}
	reader.counterLength = len(reader.counters)

	tests := []struct {
		pos   int
		wants int
	}{
		{0, 16},
		{7, 5},
		{14, -1},
		{21, -1},
	}
	for _, test := range tests {
		r := reader.toNarrowWidePattern(test.pos)
		if r != test.wants {
			t.Fatalf("toNarrowWidePattern = %v, wants %v", r, test.wants)
		}
	}
}

func TestCodabarReader_setCounters(t *testing.T) {
	reader := NewCodaBarReader().(*codabarReader)
	e := reader.setCounters(testutil.NewBitArrayFromString("11111111111111111111"))
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("setCounters must be error, %T", e)
	}

	reader = NewCodaBarReader().(*codabarReader)
	e = reader.setCounters(testutil.NewBitArrayFromString("00001001110000100"))
	if e != nil {
		t.Fatalf("setCounters returns error: %v", e)
	}
	wants := []int{4, 1, 2, 3, 4, 1, 2}
	if !reflect.DeepEqual(reader.counters, wants) {
		t.Fatalf("counters = %v, wants %v", reader.counters, wants)
	}
}

func TestCodabarReader_findStartPattern(t *testing.T) {
	tests := []struct {
		label    string
		counters []int
		found    bool
		wants    int
	}{
		{
			"no pattern",
			[]int{1, 1, 1, 1, 1, 1, 1, 2, 2, 1, 1, 1, 1, 1},
			false, 0,
		},
		{
			"no quiet zone",
			[]int{1, 1, 1, 1, 1, 2, 2, 1, 2, 1, 1, 1, 1, 1},
			false, 0,
		},
		{
			"start A",
			[]int{1, 1, 5, 1, 1, 2, 2, 1, 2, 1, 1, 1, 1, 1},
			true, 3,
		},
		{
			"start B",
			[]int{1, 1, 1, 1, 5, 1, 2, 1, 2, 1, 1, 2, 1, 1},
			true, 5,
		},
		{
			"start C",
			[]int{5, 1, 1, 1, 2, 1, 2, 2, 1, 1, 1, 1, 1, 1},
			true, 1,
		},
		{
			"start D",
			[]int{5, 1, 1, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1},
			true, 1,
		},
		{
			"start 0",
			[]int{5, 1, 1, 1, 1, 1, 2, 2, 1, 1, 1, 1, 1, 1},
			false, 0,
		},
	}

	for _, test := range tests {
		reader := NewCodaBarReader().(*codabarReader)
		reader.counters = test.counters
		reader.counterLength = len(test.counters)

		r, e := reader.findStartPattern()
		if test.found {
			if e != nil {
				t.Fatalf("findStartPattern[%v] returns error: %v", test.label, e)
			}
			if r != test.wants {
				t.Fatalf("findStartPattern[%v] = %v, wants %v", test.label, r, test.wants)
			}
		} else {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("findStartPattern[%v] must be NotFoundException, %T", test.label, e)
			}
		}
	}
}

func TestCodabarReader_validatePattern(t *testing.T) {
	reader := NewCodaBarReader().(*codabarReader)

	reader.decodeRowResult = []byte{16, 1, 15} // "A1+"

	reader.counters = []int{
		10,
		1, 1, 2, 2, 1, 2, 1, 1, // A
		3, 1, 3, 1, 6, 2, 3, 1, // 1
		1, 1, 2, 1, 2, 1, 2, 1, // +
	}
	reader.counterLength = len(reader.counters)
	e := reader.validatePattern(1)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("validatePattern must be NotFoundException, %T", e)
	}

	reader.counters = []int{
		10,
		1, 1, 2, 2, 1, 2, 1, 1, // A
		1, 1, 1, 1, 2, 2, 1, 1, // 1
		1, 1, 2, 1, 2, 1, 2, 1, // +
	}
	reader.counterLength = len(reader.counters)
	e = reader.validatePattern(1)
	if e != nil {
		t.Fatalf("validatePattern returns error: %v", e)
	}
}

func TestCodabarReader_DecodeRow(t *testing.T) {
	reader := NewCodaBarReader().(*codabarReader)

	tests := []struct {
		label string
		row   string
	}{
		{"no white", "11111111111111111111"},
		{"no start", "00000101011001000000"},
		{"invalid pattern", "000001011001001010101010"},
		{"no quiet zone", "010110010010100100101101"},
		{"validate fail", "0000010110010010111011101111110011101001001011000001"},
		{"no end pattern", "000001011001001010101100100000"},
		{"start-end only", "0000010110010010101001100100000"},
	}
	for _, test := range tests {
		row := testutil.NewBitArrayFromString(test.row)
		_, e := reader.DecodeRow(10, row, nil)
		if _, ok := e.(gozxing.NotFoundException); !ok {
			t.Fatalf("DecodeRow[%v] must be NotFoundException, %T", test.label, e)
		}
	}

	// D3$+C
	row := testutil.NewBitArrayFromString(
		"00000" + "10100110010" + "1100101010" + "1011001010" + "10110110110" + "10100100110" + "00000")
	points := []gozxing.ResultPoint{gozxing.NewResultPoint(5, 10), gozxing.NewResultPoint(57, 10)}

	r, e := reader.DecodeRow(10, row, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if txt, wants := r.GetText(), "3$+"; txt != wants {
		t.Fatalf("text = \"%v\", wants \"%v\"", txt, wants)
	}
	if rps := r.GetResultPoints(); !reflect.DeepEqual(rps, points) {
		t.Fatalf("resultPoints = %v, wants %v", rps, points)
	}

	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_RETURN_CODABAR_START_END: true,
	}
	r, e = reader.DecodeRow(10, row, hints)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if txt, wants := r.GetText(), "D3$+C"; txt != wants {
		t.Fatalf("text = \"%v\", wants \"%v\"", txt, wants)
	}
	if rps := r.GetResultPoints(); !reflect.DeepEqual(rps, points) {
		t.Fatalf("resultPoints = %v, wants %v", rps, points)
	}
}

func TestCodaBarReader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/codabar-1/
	reader := NewCodaBarReader()
	tests := []struct {
		file     string
		wants    string
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		{
			"testdata/codabar/01.png", "1234567890",
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]F0",
			},
		},
		{"testdata/codabar/02.png", "1234567890", nil},
		{"testdata/codabar/03.png", "294/586", nil},
		{"testdata/codabar/04.png", "123455", nil},
		{"testdata/codabar/09.png", "12345", nil},
		{"testdata/codabar/10.png", "123456", nil},
		{"testdata/codabar/11.png", "3419500", nil},
		{"testdata/codabar/12.png", "31117013206375", nil},
		{"testdata/codabar/13.png", "12345", nil},
		{"testdata/codabar/14.png", "31117013206375", nil},
		{"testdata/codabar/15.png", "123456789012", nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, gozxing.BarcodeFormat_CODABAR, nil, test.metadata)
	}
}
