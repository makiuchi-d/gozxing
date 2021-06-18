package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestITFReader_decodeDigit(t *testing.T) {
	tests := []struct {
		wants    int
		counters []int
	}{
		{0, []int{1, 1, 2, 2, 1}},
		{7, []int{1, 1, 1, 2, 2}},
		{9, []int{1, 3, 1, 3, 1}},
		{-1, []int{16, 16, 16, 39, 39}},
		{-1, []int{4, 4, 1, 4, 1}},
	}

	for _, test := range tests {
		r, e := itfReader_decodeDigit(test.counters)

		if test.wants < 0 {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("decodeDigit(%v) must NotFoundException, %T, %v", test.counters, e, r)
			}
		} else {
			if e != nil {
				t.Fatalf("decodeDigit(%v) returns error: %v", test.counters, e)
			}
			if r != test.wants {
				t.Fatalf("decodeDigit(%v) = %v, wants %v", test.counters, r, test.wants)
			}
		}
	}
}

func TestITFReader_findGuardPattern(t *testing.T) {
	row := testutil.NewBitArrayFromString("0000000000111111100000001100000011100")
	_, e := itfReader_findGuardPattern(row, 0, itfReader_START_PATTERN)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("findGuardPattern must be NotFoundException, %T", e)
	}

	row = testutil.NewBitArrayFromString("0000110011000111000010101000")
	r, e := itfReader_findGuardPattern(row, 0, itfReader_START_PATTERN)
	if e != nil {
		t.Fatalf("findGuardPattern returns error: %v", e)
	}
	wants := []int{4, 13}
	if !reflect.DeepEqual(r, wants) {
		t.Fatalf("findGuardPattern = %v, wants %v", r, wants)
	}

	r, e = itfReader_findGuardPattern(row, 13, itfReader_START_PATTERN)
	if e != nil {
		t.Fatalf("findGuardPattern returns error: %v", e)
	}
	wants = []int{20, 24}
	if !reflect.DeepEqual(r, wants) {
		t.Fatalf("findGuardPattern = %v, wants %v", r, wants)
	}
}

func TestITFReader_skipWhiteSpace(t *testing.T) {
	row := testutil.NewBitArrayFromString("0000000000")
	_, e := itfReader_skipWhiteSpace(row)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("skipWhiteSpace must NotFoundException, %T", e)
	}

	row = testutil.NewBitArrayFromString("0000000100")
	r, e := itfReader_skipWhiteSpace(row)
	if e != nil {
		t.Fatalf("skipWhiteSpace returns error: %v", e)
	}
	wants := 7
	if r != wants {
		t.Fatalf("skipWhiteSpace = %v, wants %v", r, wants)
	}
}

func TestITFReader_validateQuietZone(t *testing.T) {
	reader := NewITFReader().(*itfReader)

	tests := []struct {
		row    string
		narrow int
		start  int
		valid  bool
	}{
		{"10001111", 1, 4, false},
		{"11100000000000000000001111", 2, 22, false},
		{"11000000000000000000001111", 2, 22, true},
	}
	for _, test := range tests {
		row := testutil.NewBitArrayFromString(test.row)
		reader.narrowLineWidth = test.narrow
		e := reader.validateQuietZone(row, test.start)
		if test.valid {
			if e != nil {
				t.Fatalf("validateQuietZone returns error: %v", e)
			}
		} else {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("validateQuietZone must NotFoundException: %T", e)
			}
		}
	}
}

func TestITFReader_decodeEnd(t *testing.T) {
	reader := NewITFReader().(*itfReader)

	tests := []struct {
		row    string
		narrow int
		wants  []int
	}{
		{"00000000", 1, nil},
		{"00001111", 1, nil},
		{"0110001111001100010", 2, nil},
		{"0110001111001100000", 2, []int{6, 14}},
		{"0110001111111001100000", 2, []int{6, 17}},
	}
	for _, test := range tests {
		reader.narrowLineWidth = test.narrow
		row := testutil.NewBitArrayFromString(test.row)
		r, e := reader.decodeEnd(row)
		if test.wants == nil {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("decodeEnd(%v) must NotFoundException, %T", test.row, e)
			}
		} else {
			if e != nil {
				t.Fatalf("decodeEnd(%v) returns error: %v", test.row, e)
			}
			if !reflect.DeepEqual(r, test.wants) {
				t.Fatalf("decodeEnd(%v) = %v, wants %v", test.row, r, test.wants)
			}
		}
		for i, b := range test.row {
			if row.Get(i) != (b == '1') {
				t.Fatalf("row bits = \"%v\", wants \"%v\"", row.String(), test.row)
			}
		}
	}
}

func TestITFReader_decodeStart(t *testing.T) {
	reader := NewITFReader().(*itfReader)

	tests := []struct {
		row    string
		wants  []int
		narrow int
	}{
		{"000000", nil, 0},
		{"001100", nil, 0},
		{"01000000110011001111", nil, 0},
		{"0000010101111", []int{5, 9}, 1},
		{"00000000110011001111", []int{8, 16}, 2},
	}
	for _, test := range tests {
		row := testutil.NewBitArrayFromString(test.row)
		r, e := reader.decodeStart(row)
		if test.wants == nil {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("decodeStart(%v) must NotFoundException, %T", test.row, e)
			}
		} else {
			if e != nil {
				t.Fatalf("decodeStart(%v) returns error: %v", test.row, e)
			}
			if !reflect.DeepEqual(r, test.wants) {
				t.Fatalf("decodeStart(%v) = %v, wants %v", test.row, r, test.wants)
			}
			if reader.narrowLineWidth != test.narrow {
				t.Fatalf("narrowLineWidth = %v, wants %v", reader.narrowLineWidth, test.narrow)
			}
		}
	}
}

func TestITFReader_decodeMiddle(t *testing.T) {
	reader := NewITFReader().(*itfReader)

	tests := []struct {
		row    string
		start  int
		end    int
		narrow int
		wants  []byte
	}{
		{"001010101010", 6, 13, 1, nil},
		{"001010101010101011010", 6, 13, 1, nil},
		{"00101010101011011011010", 6, 18, 1, nil},
		{"0010101001011011010011010", 6, 20, 1, []byte{'0', '1'}},
		{"0010101001100101011010010110010110110100", 6, 34, 1, []byte{'2', '3', '4', '5'}},
	}
	for _, test := range tests {
		buf := []byte{}
		row := testutil.NewBitArrayFromString(test.row)
		reader.narrowLineWidth = test.narrow
		r, e := reader.decodeMiddle(row, test.start, test.end, buf)

		if test.wants == nil {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("decodeMiddle(%v) must NotFoundException, %T", test.row, e)
			}
		} else {
			if e != nil {
				t.Fatalf("decodeMiddle(%v) returns error: %v", test.row, e)
			}
			if !reflect.DeepEqual(r, test.wants) {
				t.Fatalf("decodeMiddle(%v) = %v, wants %v", test.row, r, test.wants)
			}
		}
	}
}

func TestITFReader_DecodeRow(t *testing.T) {
	reader := NewITFReader().(*itfReader)

	failtests := []string{
		"00000000",                   // decodeStart failed
		"00001010111100",             // decodeEnd failed
		"00001010101010101011010000", // decodeMiddle failed
	}
	for _, str := range failtests {
		row := testutil.NewBitArrayFromString(str)
		_, e := reader.DecodeRow(10, row, nil)
		if _, ok := e.(gozxing.NotFoundException); !ok {
			t.Fatalf("DecodeRow(%v) must NotFoundException, %T", str, e)
		}
	}

	hint := make(map[gozxing.DecodeHintType]interface{})
	hint[gozxing.DecodeHintType_ALLOWED_LENGTHS] = []int{2, 4, 6}

	rowstr := "000010101011011010010011010000"
	wants := "67"
	start := float64(8)
	end := float64(22)
	row := testutil.NewBitArrayFromString(rowstr)

	_, e := reader.DecodeRow(10, row, nil) // length (=2) is not allowed
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("DecodeRow must FormatException, %T", e)
	}

	r, e := reader.DecodeRow(10, row, hint)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if txt := r.GetText(); txt != wants {
		t.Fatalf("DecodeRow = \"%v\", wants \"%v\"", txt, wants)
	}
	pts := r.GetResultPoints()
	if x, y := pts[0].GetX(), pts[0].GetY(); x != start || y != 10 {
		t.Fatalf("start point = (%v,%v), wants (%v,%v)", x, y, start, 10)
	}
	if x, y := pts[1].GetX(), pts[1].GetY(); x != end || y != 10 {
		t.Fatalf("end point = (%v,%v), wants (%v,%v)", x, y, end, 10)
	}

	rowstr = "00001010100101101101001001100101011010010110010110101101101001001101001011001011010000"
	wants = "0123456789"
	start = float64(8)
	end = float64(78)
	row = testutil.NewBitArrayFromString(rowstr)

	r, e = reader.DecodeRow(10, row, hint)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if txt := r.GetText(); txt != wants {
		t.Fatalf("DecodeRow = \"%v\", wants \"%v\"", txt, wants)
	}
	pts = r.GetResultPoints()
	if x, y := pts[0].GetX(), pts[0].GetY(); x != start || y != 10 {
		t.Fatalf("start point = (%v,%v), wants (%v,%v)", x, y, start, 10)
	}
	if x, y := pts[1].GetX(), pts[1].GetY(); x != end || y != 10 {
		t.Fatalf("end point = (%v,%v), wants (%v,%v)", x, y, end, 10)
	}
}

func TestITFReader(t *testing.T) {
	reader := NewITFReader()

	tests := []struct {
		file     string
		wants    string
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		{
			"testdata/itf/1.png", "30712345000010",
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]I0",
			},
		},
		{"testdata/itf/2.png", "00012345678905", nil},
		{"testdata/itf/3.png", "0053611912", nil},
		{"testdata/itf/5.png", "0829220875", nil},
		{"testdata/itf/6.png", "0829220874", nil},
		{"testdata/itf/7.png", "0817605453", nil},
		{"testdata/itf/8.png", "0829220874", nil},
		{"testdata/itf/9.png", "0053611912", nil},
		{"testdata/itf/10.png", "0053611912", nil},
		{"testdata/itf/13.png", "0829220875", nil},
		{"testdata/itf/14.png", "0829220875", nil},
		{"testdata/itf/15.png", "0829220875", nil},
		{"testdata/itf/16.png", "0829220874", nil},
		{"testdata/itf/17.png", "3018108390", nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, gozxing.BarcodeFormat_ITF, nil, test.metadata)
	}
}
