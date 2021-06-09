package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestCode93CheckChecksums(t *testing.T) {
	s := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0.")
	e := code93CheckChecksums(s)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("code93CheckChecksums must be ChecksumException, %T", e)
	}

	s = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ00")
	e = code93CheckChecksums(s)
	if e != nil {
		t.Fatalf("code93CheckChecksums returns error: %v", e)
	}
}

func TestCode93DecodeExtended(t *testing.T) {
	s := []byte("a")
	_, e := code93DecodeExtended(s)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("code39DecodeExtended must be FormatException, %T", e)
	}

	for _, s := range []string{"a0", "b0", "c0", "d0"} {
		_, e := code93DecodeExtended([]byte(s))
		if _, ok := e.(gozxing.FormatException); !ok {
			t.Fatalf("code39DecodeExtended must be FormatException, %T", e)
		}
	}

	s = []byte("CdOdDdE93cAaJ$%/+ cZbAbFbKbPbTbUbVbWbXbYbZ")
	expect := "Code93!\n$%/+ :\x1b;[{\x7f\x00@`\x7f\x7f\x7f"
	r, e := code93DecodeExtended(s)
	if e != nil {
		t.Fatalf("code93DecodeExtended returns error: %v", e)
	}
	if r != expect {
		t.Fatalf("code93DecodeExtended = \"%v\", expect \"%v\"", r, expect)
	}
}

func TestCode93PatternToChar(t *testing.T) {
	_, e := code93PatternToChar(0x1ea) // unused
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("code93DecodeExtended must be NotFoundException, %T", e)
	}

	c, e := code93PatternToChar(0x114)
	if e != nil {
		t.Fatalf("code93PatternToChar returns error: %v", e)
	}
	if c != '0' {
		t.Fatalf("code93PatternToChar = '%c', expect '0'", c)
	}

	c, e = code93PatternToChar(0x15e)
	if e != nil {
		t.Fatalf("code93PatternToChar returns error: %v", e)
	}
	if c != '*' {
		t.Fatalf("code93PatternToChar = '%c', expect '*'", c)
	}
}

func TestCode93ToPattern(t *testing.T) {
	// invalid pattern
	r := code93ToPattern([]int{1, 1, 1, 1, 1, 5})
	if r != -1 {
		t.Fatalf("code93ToPattern = %v, expect %v", r, -1)
	}

	// * (start/stop)
	r = code93ToPattern([]int{1, 1, 1, 1, 4, 1})
	if r != 0x15e {
		t.Fatalf("code93ToPattern = %v, expect %v", r, 0x15e)
	}

	// 0 x3
	r = code93ToPattern([]int{3, 9, 3, 3, 3, 6})
	if r != 0x114 {
		t.Fatalf("code93ToPattern = %v, expect %v", r, 0x114)
	}
}

func TestCode93Reader_findAsteriskPattern(t *testing.T) {
	dec := NewCode93Reader().(*code93Reader)

	row := testutil.NewBitArrayFromString("00001000101101101111111")
	_, _, e := dec.findAsteriskPattern(row)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("findAsteriskPattern must be NotFoundException, %T", e)
	}

	row = testutil.NewBitArrayFromString("000010001100110011111111001111111")
	l, r, e := dec.findAsteriskPattern(row)
	if e != nil {
		t.Fatalf("findAsteriskPattern returns error: %v", e)
	}
	if l != 8 || r != 26 {
		t.Fatalf("findAsteriskPattern = %v-%v, expect 8-26", l, r)
	}
}

func TestCode93Reader_DecodeRow(t *testing.T) {
	dec := NewCode93Reader().(*code93Reader)

	// no start asterisk
	row := testutil.NewBitArrayFromString("00001010100000")
	_, e := dec.DecodeRow(1, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// error in recordPattern
	row = testutil.NewBitArrayFromString("0001010111101000")
	_, e = dec.DecodeRow(1, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// error in toPattern
	row = testutil.NewBitArrayFromString("000101011110101011111010111")
	_, e = dec.DecodeRow(1, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// error in patternToChar
	row = testutil.NewBitArrayFromString("000101011110101111010")
	_, e = dec.DecodeRow(1, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// no guard black module at end
	row = testutil.NewBitArrayFromString("000101011110" + "100010100" + "101011110")
	_, e = dec.DecodeRow(1, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// short result
	row = testutil.NewBitArrayFromString("000101011110" + "100010100" + "101011110" + "1111")
	_, e = dec.DecodeRow(1, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// invalid checksum
	row = testutil.NewBitArrayFromString("000101011110" +
		"110100010" + // C
		"100010100" + "100010100" + // checksum 0,0
		"1010111101111")
	_, e = dec.DecodeRow(1, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// invalid extented code
	row = testutil.NewBitArrayFromString("000101011110" +
		"100100110" + "100010100" + // ($)0
		"111001010" + "110110010" + // checksum $,R
		"1010111101111")
	_, e = dec.DecodeRow(1, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	row = testutil.NewBitArrayFromString("000101011110" +
		"110100010" + // C
		"100110010" + "100101100" + // o = (+)O
		"100110010" + "110010100" + // d = (+)d
		"100110010" + "110010010" + // e = (+)e
		"100001010" + "101000010" + // 93
		"110101000" + "111010110" + //Checksum A (/)
		"1010111101111")
	r, e := dec.DecodeRow(1, row, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if txt := r.GetText(); txt != "Code93" {
		t.Fatalf("text = \"%v\", expect \"Code93\"", txt)
	}
	if format := r.GetBarcodeFormat(); format != gozxing.BarcodeFormat_CODE_93 {
		t.Fatalf("format = %v, expect %v", format, gozxing.BarcodeFormat_CODE_93)
	}
	rps := r.GetResultPoints()
	if x, y := rps[0].GetX(), rps[0].GetY(); x != 7.5 || y != 1 {
		t.Fatalf("ResultPoint[0] = (%v,%v), expect (7.5,1)", x, y)
	}
	if x, y := rps[1].GetX(), rps[1].GetY(); x != 115.5 || y != 1 {
		t.Fatalf("ResultPoint[0] = (%v,%v), expect (115.5,1)", x, y)
	}
}

func TestCode93Reader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/code93-1/
	reader := NewCode93Reader()
	format := gozxing.BarcodeFormat_CODE_93

	tests := []struct {
		file     string
		wants    string
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		{
			"testdata/code93/1.png", "1234567890",
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]G0",
			},
		},
		{"testdata/code93/2.png", "CODE 93", nil},
		{"testdata/code93/3.png", "DATA", nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, nil, test.metadata)
	}
}
