package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func testEan13Reader_determineFirstDigit(t testing.TB, p int, expect byte) {
	t.Helper()
	d, e := ean13Reader_determineFirstDigit(p)
	if e != nil {
		t.Fatalf("determineFirstDigit(%v) returns errror, %v", p, e)
	}
	if d != expect {
		t.Fatalf("determineFirstDigit(%v) = %v, expect %v", p, d, expect)
	}
}

func TestEan13Reader_determineFirstDigit(t *testing.T) {
	testEan13Reader_determineFirstDigit(t, 0x00, 0)
	testEan13Reader_determineFirstDigit(t, 0x0b, 1)
	testEan13Reader_determineFirstDigit(t, 0x0d, 2)
	testEan13Reader_determineFirstDigit(t, 0x0e, 3)
	testEan13Reader_determineFirstDigit(t, 0x13, 4)
	testEan13Reader_determineFirstDigit(t, 0x19, 5)
	testEan13Reader_determineFirstDigit(t, 0x1c, 6)
	testEan13Reader_determineFirstDigit(t, 0x15, 7)
	testEan13Reader_determineFirstDigit(t, 0x16, 8)
	testEan13Reader_determineFirstDigit(t, 0x1a, 9)

	_, e := ean13Reader_determineFirstDigit(0x01)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("determineFirstDigit(0x01) must be NotFoundException, %T", e)
	}
}

func TestEan13Reader_getBarcodeFormat(t *testing.T) {
	r := &ean13Reader{}
	f := r.getBarcodeFormat()
	if f != gozxing.BarcodeFormat_EAN_13 {
		t.Fatalf("getBarcodeFormat = %v, expect EAN_13", f)
	}
}

func TestEan13Reader_decodeMiddle(t *testing.T) {
	r := &ean13Reader{decodeMiddleCounters: make([]int, 4)}
	result := make([]byte, 13)
	row := gozxing.NewBitArray(120)
	startRange := []int{3, 6}

	_, _, e := r.decodeMiddle(row, startRange, result[:0])
	if e == nil {
		t.Fatalf("decodeMiddle must be error")
	}

	// 1: LLGLGG (error=LLGLGL)
	// [6-12] 2(L): 0010011
	row.Set(8)
	row.Set(11)
	row.Set(12)
	// [13-19] 3(L): 0111101
	row.Set(14)
	row.Set(15)
	row.Set(16)
	row.Set(17)
	row.Set(19)
	// [20-26] 4(G): 0011101
	row.Set(22)
	row.Set(23)
	row.Set(24)
	row.Set(26)
	// [27-33] 5(L): 0110001
	row.Set(28)
	row.Set(29)
	row.Set(33)
	// [34-40] 6(G): 0000101
	row.Set(38)
	row.Set(40)
	// [41-47] 7(L): 0111011
	row.Set(42)
	row.Set(43)
	row.Set(44)
	row.Set(46)
	row.Set(47)
	_, _, e = r.decodeMiddle(row, startRange, result[:0])
	if e == nil {
		t.Fatalf("decodeMiddle must be error")
	}

	// [41-47] 7(G): 0010001
	row.Flip(42)
	row.Flip(44)
	row.Flip(46)
	_, _, e = r.decodeMiddle(row, startRange, result[:0])
	if e == nil {
		t.Fatalf("decodeMiddle must be error")
	}

	// center guard pattern
	row.Set(49)
	row.Set(51)
	row.Set(53)
	_, _, e = r.decodeMiddle(row, startRange, result[:0])
	if e == nil {
		t.Fatalf("decodeMiddle must be error")
	}

	// [53-59] 8(R): 1001000
	row.Set(56)
	// [60-66] 9(R): 1110100
	row.Set(60)
	row.Set(61)
	row.Set(62)
	row.Set(64)
	// [67-73] 0(R): 1110010
	row.Set(67)
	row.Set(68)
	row.Set(69)
	row.Set(72)
	// [74-80] 1(R): 1100110
	row.Set(74)
	row.Set(75)
	row.Set(78)
	row.Set(79)
	// [81-87] 2(R): 1101100
	row.Set(81)
	row.Set(82)
	row.Set(84)
	row.Set(85)
	// [88-94] 3(R): 1000010
	row.Set(88)
	row.Set(93)
	// end
	row.Set(95)

	offset, result, e := r.decodeMiddle(row, startRange, result[:0])
	if e != nil {
		t.Fatalf("decodeMiddle retuns error, %v", e)
	}
	if offset != 95 {
		t.Fatalf("decodeMiddle offset = %v, expect 95", offset)
	}
	if s := string(result); s != "1234567890123" {
		t.Fatalf("decodeMiddle result = \"%v\", expect \"1234567890123\"", s)
	}
}

func TestEan13Reader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/ean13-1/
	reader := NewEAN13Reader()
	format := gozxing.BarcodeFormat_EAN_13
	harder := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}

	tests := []struct {
		file     string
		wants    string
		hints    map[gozxing.DecodeHintType]interface{}
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		{
			"testdata/ean13/1.png", "8413000065504", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]E0",
			},
		},
		{"testdata/ean13/2.png", "8480010092271", nil, nil},
		{"testdata/ean13/3.png", "8480000823274", nil, nil},
		{"testdata/ean13/4.png", "5449000039231", nil, nil},
		{"testdata/ean13/5.png", "8410054010412", nil, nil},
		{"testdata/ean13/6.png", "8480010045062", nil, nil},
		{"testdata/ean13/7.png", "9788430532674", nil, nil},
		{"testdata/ean13/8.png", "8480017507990", nil, nil},
		{"testdata/ean13/9.png", "3166298099809", nil, nil},
		{"testdata/ean13/10.png", "8480010001136", nil, nil},
		{"testdata/ean13/12.png", "5201815331227", nil, nil},
		{"testdata/ean13/13.png", "8413600298517", nil, nil},
		{"testdata/ean13/14.png", "3560070169443", nil, nil},
		{"testdata/ean13/15.png", "4045787034318", nil, nil},
		{"testdata/ean13/18.png", "3086126100326", nil, nil},
		{"testdata/ean13/19.png", "4820024790635", nil, nil},
		{"testdata/ean13/20.png", "4000539017100", harder, nil},
		{"testdata/ean13/21.png", "7622200008018", nil, nil},
		{"testdata/ean13/22.png", "5603667020517", nil, nil},
		{"testdata/ean13/23.png", "7622400791949", nil, nil},
		{"testdata/ean13/24.png", "5709262942503", nil, nil},
		{"testdata/ean13/25.png", "9780140013993", nil, nil},
		{"testdata/ean13/26.png", "4901780188352", nil, nil},
		{"testdata/ean13/28.png", "9771699057002", nil, nil},
		{"testdata/ean13/29.png", "4007817327098", nil, nil},
		{"testdata/ean13/30.png", "5025121072311", nil, nil},
		{"testdata/ean13/31.png", "9780393058673", nil, nil},
		{"testdata/ean13/32.png", "9780393058673", nil, nil},
		{"testdata/ean13/33.png", "9781558604971", nil, nil},
		// original zxing could't read too
		// {"testdata/ean13/34.png", "9781558604971", harder, nil},
		{"testdata/ean13/35.png", "5030159003930", harder, nil},
		// original zxing couldn't read too
		// {"testdata/ean13/36.png", "5000213101025", harder, nil},
		{"testdata/ean13/37.png", "5000213002834", harder, nil},
		{"testdata/ean13/38.png", "9780201752847", harder, nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, test.hints, test.metadata)
	}
}

func TestEan13ReaderWithExtension(t *testing.T) {
	// testdata from zxing core/src/test/resources/benchmark/android-2/ean13-1.png
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_ALLOWED_EAN_EXTENSIONS: []int{2, 5},
	}

	reader := NewEAN13Reader()
	bmp := testutil.NewBinaryBitmapFromFile("testdata/ean13/ean13-1.png")
	result, e := reader.Decode(bmp, hints)
	if e != nil {
		t.Fatalf("read file failed, %v", e)
	}

	expect := "9780201310054"
	if txt := result.GetText(); txt != expect {
		t.Fatalf("result = \"%v\", expect \"%v\"", txt, expect)
	}

	meta := result.GetResultMetadata()
	ext, ok := meta[gozxing.ResultMetadataType_UPC_EAN_EXTENSION]
	if !ok {
		t.Fatalf("metadata must contain key UPC_EAN_EXTENSION")
	}
	if ext != "54999" {
		t.Fatalf("metadata[UPC_EAN_EXTENSION] = \"%v\", expect \"54999\"", ext)
	}
	price, ok := meta[gozxing.ResultMetadataType_SUGGESTED_PRICE]
	if !ok {
		t.Fatalf("metadata must contain key SUGGESTED_PRICE")
	}
	if price != "$49.99" {
		t.Fatalf("metadata[SUGGESTED_PRICE] = \"%v\", expect \"$49.99\"", price)
	}
}
