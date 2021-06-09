package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestCode39DecodeExtended(t *testing.T) {
	for _, str := range []string{"+0", "$0", "%0", "/0"} {
		_, e := code39DecodeExtended([]byte(str))
		if _, ok := e.(gozxing.FormatException); !ok {
			t.Fatalf("code39DecodeExtended(\"%v\") must be FormatException, %T", str, e)
		}
	}

	str := "0+A$J/A/O/Z%A%F%K%P%U%V%W%X%Y%Z"
	expect := "0a\n!/:\x1b;[{\x00@`\x7f\x7f\x7f"
	r, e := code39DecodeExtended([]byte(str))
	if e != nil {
		t.Fatalf("code39DecodeExtended returns error: %v", e)
	}
	if r != expect {
		t.Fatalf("code39DecodeExtended = \"%v\", expect = \"%v\"", r, expect)
	}
}

func TestCode39PatternToChar(t *testing.T) {
	_, e := code39PatternToChar(0x30)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("code39PatternToChar must be NotFoundException, %T", e)
	}

	c, e := code39PatternToChar(0x94)
	if e != nil {
		t.Fatalf("code39PatternToChar returns error: %v", e)
	}
	if c != '*' {
		t.Fatalf("code39PatternToChar = '%c', expect '%c'", c, '*')
	}

	c, e = code39PatternToChar(0x0D0)
	if e != nil {
		t.Fatalf("code39PatternToChar returns error: %v", e)
	}
	if c != 'Z' {
		t.Fatalf("code39PatternToChar = '%c', expect '%c'", c, 'Z')
	}
}

func TestCode39ToNarrowWidePattern(t *testing.T) {

	counters := []int{1, 2, 2, 2, 2, 1, 1, 1, 1}
	r := code39ToNarrowWidePattern(counters)
	if r != -1 {
		t.Fatalf("code39ToNarrowWidePattern = %v expect -1", r)
	}

	counters = []int{4, 2, 2, 1, 1, 1, 1, 1, 1}
	r = code39ToNarrowWidePattern(counters)
	if r != -1 {
		t.Fatalf("code39ToNarrowWidePattern = %v expect -1", r)
	}

	counters = []int{3, 7, 3, 2, 6, 4, 6, 3, 3} // 1,2,1,1,2,1,2,1,1
	r = code39ToNarrowWidePattern(counters)
	if r != code39AsteriskEncoding {
		t.Fatalf("code39ToNarrowWidePattern = %v, expect %v", r, code39AsteriskEncoding)
	}
}

func TestCode39FindAsteriskPattern(t *testing.T) {
	counters := make([]int, 9)

	src := testutil.NewBitArrayFromString("000110100000010110110110100000011111")
	_, _, e := code39FindAsteriskPattern(src, counters)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("code39FindAsteriskPattern must be NotFoundException, %T", e)
	}

	for i := range counters {
		counters[i] = 0
	}

	src = testutil.NewBitArrayFromString("0001101000000100101101101010100011111")
	l, r, e := code39FindAsteriskPattern(src, counters)
	if e != nil {
		t.Fatalf("code39FindAsteriskPattern returns error: %v", e)
	}
	if l != 13 || r != 25 {
		t.Fatalf("code39FindAsteriskPattern = %v-%v, expect 13-25", l, r)
	}
}

func TestCode39Reader_DecodeRow(t *testing.T) {
	dec := NewCode39Reader().(*code39Reader)

	// no start asterisk
	src := testutil.NewBitArrayFromString("0000000")
	_, e := dec.DecodeRow(1, src, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// error on recordPattern
	src = testutil.NewBitArrayFromString("000" + "1001011011010")
	_, e = dec.DecodeRow(1, src, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// error on code39ToNarrowWidePattern
	src = testutil.NewBitArrayFromString("000" + "1001011011010" + "1010101010")
	_, e = dec.DecodeRow(1, src, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// error on code39PatternToChar
	src = testutil.NewBitArrayFromString("000" + "1001011011010" + "1101101101010")
	_, e = dec.DecodeRow(1, src, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// less whitespace after end asterisk
	src = testutil.NewBitArrayFromString("000" + "1001011011010" + "1001011011010" + "001")
	_, e = dec.DecodeRow(1, src, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// empty result
	src = testutil.NewBitArrayFromString("000" + "1001011011010" + "1001011011010" + "00")
	_, e = dec.DecodeRow(1, src, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// *A*
	src = testutil.NewBitArrayFromString("000" + "1001011011010" + "1101010010110" + "1001011011010" + "000")
	r, e := dec.DecodeRow(1, src, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if txt := r.GetText(); txt != "A" {
		t.Fatalf("text = \"%v\", expect \"A\"", txt)
	}
	if format := r.GetBarcodeFormat(); format != gozxing.BarcodeFormat_CODE_39 {
		t.Fatalf("format = %v, expect %v", format, gozxing.BarcodeFormat_CODE_39)
	}
	rps := r.GetResultPoints()
	if x, y := rps[0].GetX(), rps[0].GetY(); x != 9 || y != 1 {
		t.Fatalf("ResultPoint[0] = (%v,%v), expect (9,1)", x, y)
	}
	if x, y := rps[1].GetX(), rps[1].GetY(); x != 35 || y != 1 {
		t.Fatalf("ResultPoint[0] = (%v,%v), expect (35,1)", x, y)
	}
}

func TestCode39Reader_DecodeRowWithExtendedModeCheckDigit(t *testing.T) {
	dec := NewCode39ReaderWithFlags(true, true).(*code39Reader)

	// *AA*, checkdigit=20 => 'K'
	src := testutil.NewBitArrayFromString("000" + "1001011011010" +
		"1101010010110" + "1101010010110" + "1101010010110" + // AAA
		"1001011011010" + "000")
	_, e := dec.DecodeRow(1, src, nil)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("DecodeRow must be ChecksumException, %T", e)
	}
	// *+0* (invalid extended string), checkdigit=41 =>'+'
	src = testutil.NewBitArrayFromString("000" + "1001011011010" +
		"1001010010010" + "1010011011010" + "1001010010010" + // +0+
		"1001011011010" + "000")
	_, e = dec.DecodeRow(1, src, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("DecodeRow must be FormatException, %T", e)
	}

	// 1+A%V => "1a@", checkdigit=39 => '$'
	expect := "1a@"
	src = testutil.NewBitArrayFromString("000" + "1001011011010" +
		"1101001010110" + "1001010010010" + "1101010010110" + // 1+A
		"1010010010010" + "1001101010110" + "1001001001010" + // %V$
		"1001011011010" + "000")
	r, e := dec.DecodeRow(1, src, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if txt := r.GetText(); txt != expect {
		t.Fatalf("text = \"%v\", expect \"%v\"", txt, expect)
	}
}

func TestCode39Reader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/code39-*/
	format := gozxing.BarcodeFormat_CODE_39

	reader := NewCode39Reader()
	tests := []struct {
		file     string
		wants    string
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		{
			"testdata/code39/1-1.png", "TEST-SHEET",
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]A0",
			},
		},
		{"testdata/code39/1-2.png", " WWW.CITRONSOFT.COM ", nil},
		{"testdata/code39/1-3.png", "MOROVIA", nil},
		{"testdata/code39/1-4.png", "ABC123", nil},
		{"testdata/code39/01.png", "165627", nil},
		{"testdata/code39/02.png", "165627", nil},
		{"testdata/code39/03.png", "001EC947D49B", nil},
		{"testdata/code39/04.png", "001EC947D49B", nil},
		{"testdata/code39/05.png", "001EC947D49B", nil},
		{"testdata/code39/06.png", "165340", nil},
		{"testdata/code39/07.png", "165340", nil},
		{"testdata/code39/08.png", "165340", nil},
		{"testdata/code39/09.png", "165340", nil},
		{"testdata/code39/10.png", "001EC94767E0", nil},
		{"testdata/code39/11.png", "001EC94767E0", nil},
		{"testdata/code39/12.png", "001EC94767E0", nil},
		{"testdata/code39/13.png", "001EC94767E0", nil},
		{"testdata/code39/14.png", "404785", nil},
		{"testdata/code39/15.png", "404785", nil},
		{"testdata/code39/16.png", "404785", nil},
		{"testdata/code39/17.png", "404785", nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, nil, test.metadata)
	}

	// extended mode
	reader = NewCode39ReaderWithFlags(false, true)
	testutil.TestFile(t, reader, "testdata/code39/2-1.png", "Extended !?*#", format, nil, nil)
	testutil.TestFile(t, reader, "testdata/code39/2-2.png", "12ab", format, nil, nil)
}
