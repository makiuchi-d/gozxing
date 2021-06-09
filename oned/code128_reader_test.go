package oned

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestCode128FindStartPattern(t *testing.T) {
	row := testutil.NewBitArrayFromString("00010000000110100100100100")
	_, e := code128FindStartPattern(row)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("code128FindStartPattern must be NotFoundException, %T", e)
	}

	// Start Code A
	row = testutil.NewBitArrayFromString("00010000011010000100100100100")
	r, e := code128FindStartPattern(row)
	if e != nil {
		t.Fatalf("code128FindStartPattern returns error: %v", e)
	}
	if r[0] != 9 {
		t.Fatalf("pattern start = %v, expect 9", r[0])
	}
	if r[1] != 20 {
		t.Fatalf("pattern end = %v, epxpect 20", r[1])
	}
	if r[2] != code128CODE_START_A {
		t.Fatalf("start code = %v, expect %v", r[2], code128CODE_CODE_A)
	}

	// Start Code B
	row = testutil.NewBitArrayFromString("0010000000000011110011000011000000001100111001111")
	r, e = code128FindStartPattern(row)
	if e != nil {
		t.Fatalf("code128FindStartPattern returns error: %v", e)
	}
	if r[0] != 14 {
		t.Fatalf("pattern start = %v, expect 14", r[0])
	}
	if r[1] != 36 {
		t.Fatalf("pattern end = %v, epxpect 36", r[1])
	}
	if r[2] != code128CODE_START_B {
		t.Fatalf("start code = %v, expect %v", r[2], code128CODE_CODE_B)
	}

	// Start Code C
	row = testutil.NewBitArrayFromString("011010011100100")
	r, e = code128FindStartPattern(row)
	if e != nil {
		t.Fatalf("code128FindStartPattern returns error: %v", e)
	}
	if r[0] != 1 {
		t.Fatalf("pattern start = %v, expect 1", r[0])
	}
	if r[1] != 12 {
		t.Fatalf("pattern end = %v, epxpect 12", r[1])
	}
	if r[2] != code128CODE_START_C {
		t.Fatalf("start code = %v, expect %v", r[2], code128CODE_CODE_C)
	}
}

func TestCode128DecodeCode(t *testing.T) {
	counters := make([]int, 6)

	row := testutil.NewBitArrayFromString("10101")
	_, e := code128DecodeCode(row, counters, 0)
	if e == nil {
		t.Fatalf("code128DecodeCode must be error")
	}

	row = testutil.NewBitArrayFromString("1010101")
	_, e = code128DecodeCode(row, counters, 0)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("code128DecodeCode must be NotFoundException, %T", e)
	}

	row = testutil.NewBitArrayFromString("110011011001")
	r, e := code128DecodeCode(row, counters, 0)
	if e != nil {
		t.Fatalf("code128DecodeCode returns error: %v", e)
	}
	if r != 1 {
		t.Fatalf("code128DecodeCode = %v, expect 1", r)
	}
}

func TestCode128Reader_DecodeRowFail(t *testing.T) {
	dec := NewCode128Reader().(*code128Reader)

	row := testutil.NewBitArrayFromString("00010000000110100100100100")
	_, e := dec.DecodeRow(10, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// error in code128decodeCode
	row = testutil.NewBitArrayFromString("00010000011010000100101010101")
	_, e = dec.DecodeRow(10, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// start code after started
	row = testutil.NewBitArrayFromString("00010000011010010000110100001001")
	_, e = dec.DecodeRow(10, row, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("DecodeRow must be FormatException, %T", e)
	}

	// less tailing space
	row = testutil.NewBitArrayFromString("0000000" + "11010000100" +
		"1100011101011" + "000001")
	_, e = dec.DecodeRow(10, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// checksum error
	row = testutil.NewBitArrayFromString("0000000" + "11010000100" +
		"11001101100" + // 1
		"1100011101011" + "00000001")
	_, e = dec.DecodeRow(10, row, nil)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("DecodeRow must be ChecksumException, %T", e)
	}

	// empty result
	row = testutil.NewBitArrayFromString("0000000" + "11010000100" +
		"1100011101011" + "00000001")
	_, e = dec.DecodeRow(10, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}
}

func TestCode128Reader_DecodeRowCodeA(t *testing.T) {
	dec := NewCode128Reader().(*code128Reader)
	hint := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_ASSUME_GS1: true,
	}

	row := testutil.NewBitArrayFromString("0000000" + "11010000100" + // StartA
		"11110101110" + "11001101100" + "10000110010" + //FNC1 ! \n
		"11110101110" + "11110101000" + "10111100010" + // FNC1 FNC2 FNC3
		"11101011110" + "11101011110" + "11001101100" + "10010110000" + // FNC4 FNC4 160+1 64+65
		"11101011110" + "11101011110" + "11001101100" + // FNC4 FNC4 !
		"11110100010" + "10010110000" + // Shift a
		"10111101110" + "10010110000" + "11101011110" + //CodeB a CodeA
		"10111011110" + "11011001100" + "11101011110" + //CodeC 00 CodeA
		"10001111010" + // Checksum=79
		"1100011101011" + "00000001")
	r, e := dec.DecodeRow(10, row, hint)
	expTxt := "]C1!\n\x1d\xa1\x81!aa00"
	expRaw := []byte{
		103, 102, 1, 74, 102, 97, 96, 101, 101, 1, 65, 101, 101, 1,
		98, 65, 100, 65, 101, 99, 0, 101, 79, 106,
	}
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if format := r.GetBarcodeFormat(); format != gozxing.BarcodeFormat_CODE_128 {
		t.Fatalf("format = %v, expect %v", format, gozxing.BarcodeFormat_CODE_128)
	}
	if txt := r.GetText(); txt != expTxt {
		t.Fatalf("text = \"%v\", expect \"%v\"", txt, expTxt)
	}
	if raw := r.GetRawBytes(); !reflect.DeepEqual(raw, expRaw) {
		t.Fatalf("rawBytes = %v, expect %v", raw, expRaw)
	}
	rps := r.GetResultPoints()
	if x, y := rps[0].GetX(), rps[0].GetY(); x != 12.5 || y != 10 {
		t.Fatalf("resultPoint[0] = (%v,%v), expect(12.5,10)", x, y)
	}
	rps = r.GetResultPoints()
	if x, y := rps[1].GetX(), rps[1].GetY(); x != 265.5 || y != 10 {
		t.Fatalf("resultPoint[0] = (%v,%v), expect(265.5,10)", x, y)
	}
}

func TestCode128Reader_DecodeRowSymbologyIdentifier(t *testing.T) {
	fmt.Println("====DecodeRowSymbologyIdentifier====")

	dec := NewCode128Reader().(*code128Reader)
	hint := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_ASSUME_GS1: true,
	}

	tests := []struct {
		symbology string
		bits      string
	}{
		{
			"]C1", "0000000" + "11010000100" + // StartA
				"11110101110" + //FNC1
				"11110101110" + // Checksum=102
				"1100011101011" + "00000001",
		},
		{
			"]C2", "0000000" + "11010000100" + // StartA
				"10100011000" + "11110101110" + // A, FNC1
				"11011000110" + // Checksum=31
				"1100011101011" + "00000001",
		},
		{
			"]C4", "0000000" + "11010000100" + // StartA
				"10100011000" + "11110101000" + // A, FNC2
				"11011100100" + // Checksum=21
				"1100011101011" + "00000001",
		},
		{
			"]C1", "0000000" + "11010010000" + // StartB
				"11110101110" + //FNC1
				"11011001100" + // Checksum=0
				"1100011101011" + "00000001",
		},
		{
			"]C2", "0000000" + "11010010000" + // StartB
				"10100011000" + "11110101110" + // A FNC1
				"11000110110" + // Checksum=32
				"1100011101011" + "00000001",
		},
		{
			"]C4", "0000000" + "11010010000" + // StartB
				"10100011000" + "11110101000" + // A, FNC2
				"11001110100" + // Checksum=22
				"1100011101011" + "00000001",
		},
		{
			"]C1", "0000000" + "11010011100" + // StartC
				"11110101110" + //FNC1
				"11001101100" + // Checksum=1
				"1100011101011" + "00000001",
		},
		{
			"]C2", "0000000" + "11010000100" + // StartA
				"10100011000" + "10111011110" + "11110101110" + // A CodeC FNC1
				"11001110100" + // Checksum=22
				"1100011101011" + "00000001",
		},
	}
	for _, test := range tests {
		row := testutil.NewBitArrayFromString(test.bits)
		r, e := dec.DecodeRow(10, row, hint)
		if e != nil {
			t.Fatalf("DecodeRow (symbology=%v) error: %v", test.symbology, e)
		}
		symbology, ok := r.GetResultMetadata()[gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER].(string)
		if !ok {
			t.Fatalf("DecodeRow (symbology=%v) metadata not contains SIMBOLOGY_IDENTIFIER", test.symbology)
		}
		if symbology != test.symbology {
			t.Fatalf(
				"DecodeRow (symbology=%v) metadata[SIMBOLOGY_IDENTIFIER] = %v, wants %v",
				test.symbology, symbology, test.symbology)
		}
	}
}

func TestCode128Reader_DecodeRowCodeB(t *testing.T) {
	dec := NewCode128Reader().(*code128Reader)
	hint := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_ASSUME_GS1: true,
	}

	row := testutil.NewBitArrayFromString("0000000" + "11010010000" + // StartB
		"11110101110" + "11001101100" + "10010110000" + //FNC1 ! a
		"11110101110" + "11110101000" + "10111100010" + // FNC1 FNC2 FNC3
		"10111101110" + "10111101110" + "11001101100" + // FNC4 FNC4 160+1
		"10111101110" + "10111101110" + "11001101100" + // FNC4 FNC4 1
		"11110100010" + "10000110010" + // Shift \n
		"11101011110" + "10000110010" + "10111101110" + // CodeA \n CodeB
		"10111011110" + "11011001100" + "10111101110" + // CodeC 00 CodeB
		"11001011100" + // Checksum=19
		"1100011101011" + "00000001")
	r, e := dec.DecodeRow(10, row, hint)
	expTxt := "]C1!a\x1d\xa1!\n\n00"
	expRaw := []byte{
		104, 102, 1, 65, 102, 97, 96, 100, 100, 1, 100, 100, 1,
		98, 74, 101, 74, 100, 99, 0, 100, 19, 106,
	}
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if format := r.GetBarcodeFormat(); format != gozxing.BarcodeFormat_CODE_128 {
		t.Fatalf("format = %v, expect %v", format, gozxing.BarcodeFormat_CODE_128)
	}
	if txt := r.GetText(); txt != expTxt {
		t.Fatalf("text = \"%v\", expect \"%v\"", txt, expTxt)
	}
	if raw := r.GetRawBytes(); !reflect.DeepEqual(raw, expRaw) {
		t.Fatalf("rawBytes = %v, expect %v", raw, expRaw)
	}
	rps := r.GetResultPoints()
	if x, y := rps[0].GetX(), rps[0].GetY(); x != 12.5 || y != 10 {
		t.Fatalf("resultPoint[0] = (%v,%v), expect(12.5,10)", x, y)
	}
	rps = r.GetResultPoints()
	if x, y := rps[1].GetX(), rps[1].GetY(); x != 254.5 || y != 10 {
		t.Fatalf("resultPoint[0] = (%v,%v), expect(254.5,10)", x, y)
	}
}

func TestCode128Reader_DecodeRowCodeC(t *testing.T) {
	dec := NewCode128Reader().(*code128Reader)
	hint := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_ASSUME_GS1: true,
	}

	row := testutil.NewBitArrayFromString("0000000" + "11010011100" + // StartC
		"11110101110" + "10111011110" + "11110101110" + // FNC1 99 FNC1
		"11101011110" + "10000110010" + "10111011110" + // CodeA \n CodeC
		"10111101110" + "10010110000" + "10111011110" + // CodeB a CodeC
		"10110000100" + // Checksum=70
		"1100011101011" + "00000001")
	r, e := dec.DecodeRow(10, row, hint)
	expTxt := "]C199\x1d\na"
	expRaw := []byte{
		105, 102, 99, 102, 101, 74, 99, 100, 65, 99, 70, 106,
	}
	if e != nil {
		t.Fatalf("DecodeRow returns error: %v", e)
	}
	if format := r.GetBarcodeFormat(); format != gozxing.BarcodeFormat_CODE_128 {
		t.Fatalf("format = %v, expect %v", format, gozxing.BarcodeFormat_CODE_128)
	}
	if txt := r.GetText(); txt != expTxt {
		t.Fatalf("text = \"%v\", expect \"%v\"", txt, expTxt)
	}
	if raw := r.GetRawBytes(); !reflect.DeepEqual(raw, expRaw) {
		t.Fatalf("rawBytes = %v, expect %v", raw, expRaw)
	}
	rps := r.GetResultPoints()
	if x, y := rps[0].GetX(), rps[0].GetY(); x != 12.5 || y != 10 {
		t.Fatalf("resultPoint[0] = (%v,%v), expect(12.5,10)", x, y)
	}
	rps = r.GetResultPoints()
	if x, y := rps[1].GetX(), rps[1].GetY(); x != 133.5 || y != 10 {
		t.Fatalf("resultPoint[0] = (%v,%v), expect(133.5,10)", x, y)
	}
}

func TestCode128Reader(t *testing.T) {
	reader := NewCode128Reader()
	format := gozxing.BarcodeFormat_CODE_128
	harder := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}

	tests := []struct {
		file     string
		wants    string
		hints    map[gozxing.DecodeHintType]interface{}
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		// testdata from zxing core/src/test/resources/blackbox/code128-1/
		{
			"testdata/code128/1.png", "168901", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]C1",
			},
		},
		{
			"testdata/code128/2.png", "Code 128", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]C0",
			},
		},
		// testdata from zxing core/src/test/resources/blackbox/code128-2/
		{"testdata/code128/01.png", "005-3379497200006", nil, nil},
		{"testdata/code128/02.png", "005-3379497200006", nil, nil},
		{"testdata/code128/03.png", "005-3379497200006", nil, nil},
		{"testdata/code128/04.png", "005-3379497200006", nil, nil},
		{"testdata/code128/05.png", "15182881", nil, nil},
		{"testdata/code128/06.png", "15182881", nil, nil},
		{"testdata/code128/07.png", "15182881", nil, nil},
		{"testdata/code128/08.png", "15182881", nil, nil},
		{"testdata/code128/09.png", "CNK8181G2C", harder, nil},
		{"testdata/code128/10.png", "CNK8181G2C", nil, nil},
		{"testdata/code128/11.png", "CNK8181G2C", nil, nil},
		{"testdata/code128/12.png", "CNK8181G2C", harder, nil},
		{"testdata/code128/13.png", "1PEF224A4", nil, nil},
		{"testdata/code128/14.png", "1PEF224A4", nil, nil},
		{"testdata/code128/15.png", "1PEF224A4", nil, nil},
		{"testdata/code128/16.png", "1PEF224A4", nil, nil},
		{"testdata/code128/17.png", "FW727", nil, nil},
		{"testdata/code128/18.png", "FW727", nil, nil},
		{"testdata/code128/19.png", "FW727", nil, nil},
		{"testdata/code128/20.png", "FW727", nil, nil},
		{"testdata/code128/21.png", "005-3354174500018", nil, nil},
		{"testdata/code128/22.png", "005-3354174500018", nil, nil},
		{"testdata/code128/23.png", "005-3354174500018", nil, nil},
		{"testdata/code128/24.png", "005-3354174500018", nil, nil},
		{"testdata/code128/25.png", "31001171800000017989625355702636", nil, nil},
		{"testdata/code128/26.png", "31001171800000017989625355702636", nil, nil},
		{"testdata/code128/27.png", "31001171800000017989625355702636", nil, nil},
		{"testdata/code128/28.png", "31001171800000017989625355702636", nil, nil},
		{"testdata/code128/29.png", "42094043", nil, nil},
		// original zxing could not read too.
		// {"testdata/code128/30.png", "42094043", harder},
		{"testdata/code128/31.png", "42094043", nil, nil},
		{"testdata/code128/32.png", "42094043", nil, nil},
		{"testdata/code128/33.png", "30885909173823", nil, nil},
		{"testdata/code128/34.png", "30885909173823", nil, nil},
		{"testdata/code128/35.png", "30885909173823", nil, nil},
		{"testdata/code128/36.png", "30885909173823", nil, nil},
		{"testdata/code128/37.png", "FGGQ6D1", harder, nil},
		{"testdata/code128/38.png", "FGGQ6D1", nil, nil},
		{"testdata/code128/39.png", "FGGQ6D1", nil, nil},
		{"testdata/code128/40.png", "FGGQ6D1", nil, nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, test.hints, test.metadata)
	}
}
