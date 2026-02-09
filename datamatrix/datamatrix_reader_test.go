package datamatrix

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

const (
	dm4str = "" +
		"                                        \n" +
		"                                        \n" +
		"    ##  ##  ##  ##  ##  ##  ##  ##      \n" +
		"    ##  ####  ####  ##  ####      ##    \n" +
		"    ####      ####  ####                \n" +
		"    ######  ####  ##    ####      ##    \n" +
		"    ####                ##########      \n" +
		"    ######      ##    ##        ####    \n" +
		"    ####    ####          ##            \n" +
		"    ######  ####  ##    ##        ##    \n" +
		"    ##    ##  ##############            \n" +
		"    ##  ####  ##    ######  ####  ##    \n" +
		"    ######  ##    ##      ######        \n" +
		"    ##  ####  ######        ####  ##    \n" +
		"    ##    ##  ######  ##  ######        \n" +
		"    ########    ##  ##        ##  ##    \n" +
		"    ############  ##  ######    ##      \n" +
		"    ################################    \n" +
		"                                        \n" +
		"                                        \n"

	dmstr = "" +
		"                                                \n" +
		"                                                \n" +
		"      ##  ##  ##  ##  ##  ##  ##  ##  ##        \n" +
		"      ##########      ####  ####        ##      \n" +
		"      ##    ##  ##########    ####    ##        \n" +
		"      ######  ######        ####  ########      \n" +
		"      ##    ##  ########  ##  ##  ####          \n" +
		"      ##  ######  ##  ##        ##########      \n" +
		"      ######  ##    ##    ##  ####              \n" +
		"      ######        ##  ####            ##      \n" +
		"      ####  ##  ####  ##    ##  ##              \n" +
		"      ####  ##              ##    ########      \n" +
		"      ######    ##        ####                  \n" +
		"      ##  ####  ######  ##########      ##      \n" +
		"      ####  ######    ##      ####              \n" +
		"      ##  ##########        ##############      \n" +
		"      ##############  ##          ##  ##        \n" +
		"      ########      ##    ######  ##    ##      \n" +
		"      ##  ##  ##        ##  ########            \n" +
		"      ####################################      \n" +
		"                                                \n" +
		"                                                \n"
)

func TestModuleSize(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)
	leftTopBlack := []int{3, 3}

	_, e := moduleSize(leftTopBlack, image)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("moduleSize must be NotFoundException, %T", e)
	}

	image.SetRegion(3, 3, 7, 7)
	_, e = moduleSize(leftTopBlack, image)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("moduleSize must be NotFoundException, %T", e)
	}

	image.Clear()
	image.SetRegion(3, 3, 3, 3)
	m, e := moduleSize(leftTopBlack, image)
	if e != nil {
		t.Fatalf("moduleSize returns error, %v", e)
	}
	if m != 3 {
		t.Fatalf("moduleSize = %v, expect 3", m)
	}
}

func TestExtractPureBits(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 8)

	_, e := extractPureBits(image)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("extractPureBits must be NotFoundException, %T", e)
	}

	image.SetRegion(3, 3, 7, 5)
	_, e = extractPureBits(image)
	if e == nil {
		t.Fatalf("extractPureBits must be error")
	}

	image.Clear()
	image.SetRegion(0, 0, 9, 8)
	_, e = extractPureBits(image)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("extractPureBits must be NotFoundException, %T", e)
	}

	image, _ = gozxing.ParseStringToBitMatrix(dmstr, "##", "  ")
	image = testutil.ExpandBitMatrix(image, 2)
	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"##  ##  ##  ##  ##  ##  ##  ##  ##  \n"+
		"##########      ####  ####        ##\n"+
		"##    ##  ##########    ####    ##  \n"+
		"######  ######        ####  ########\n"+
		"##    ##  ########  ##  ##  ####    \n"+
		"##  ######  ##  ##        ##########\n"+
		"######  ##    ##    ##  ####        \n"+
		"######        ##  ####            ##\n"+
		"####  ##  ####  ##    ##  ##        \n"+
		"####  ##              ##    ########\n"+
		"######    ##        ####            \n"+
		"##  ####  ######  ##########      ##\n"+
		"####  ######    ##      ####        \n"+
		"##  ##########        ##############\n"+
		"##############  ##          ##  ##  \n"+
		"########      ##    ######  ##    ##\n"+
		"##  ##  ##        ##  ########      \n"+
		"####################################\n", "##", "  ")

	b, e := extractPureBits(image)
	if e != nil {
		t.Fatalf("extractPureBits returns error, %v", e)
	}
	for j := 0; j < expect.GetHeight(); j++ {
		for i := 0; i < expect.GetWidth(); i++ {
			if b.Get(i, j) != expect.Get(i, j) {
				t.Fatalf("bits(%v,%v) = %v, expect %v", i, j, b.Get(i, j), expect.Get(i, j))
			}
		}
	}
}

func TestDataMatrixReader_Reset(t *testing.T) {
	d := NewDataMatrixReader()
	d.Reset() // do nothing
}

func TestDataMatrixReader_Metadata(t *testing.T) {
	// Tests ERRORS_CORRECTED and SYMBOLOGY_IDENTIFIER metadata (GS1 reader support)
	reader := NewDataMatrixReader()
	pure := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_PURE_BARCODE: true,
	}
	format := gozxing.BarcodeFormat_DATA_MATRIX

	tests := []struct {
		file            string
		wants           string
		hints           map[gozxing.DecodeHintType]interface{}
		errorsCorrected int
		symbologyID     string
	}{
		{"testdata/0123456789.png", "0123456789", nil, 0, "]d1"},
		{"testdata/HelloWorld_Text_L_Kaywa.png", "Hello World", pure, 0, "]d1"},
		{"testdata/HelloWorld_Text_L_Kaywa_1_error_byte.png", "Hello World", pure, 1, "]d1"},
		{"testdata/HelloWorld_Text_L_Kaywa_2_error_byte.png", "Hello World", pure, 2, "]d1"},
		{"testdata/HelloWorld_Text_L_Kaywa_3_error_byte.png", "Hello World", pure, 3, "]d1"},
		{"testdata/HelloWorld_Text_L_Kaywa_4_error_byte.png", "Hello World", pure, 5, "]d1"},
		// GS1 format: FNC1 in first position, symbology modifier ]d2
		{"testdata/gs1.png", "\x1d0105391538880109172611301028967\x1d211049931544602459836", nil, 0, "]d2"},
	}
	for _, tt := range tests {
		metadata := map[gozxing.ResultMetadataType]interface{}{
			gozxing.ResultMetadataType_ERRORS_CORRECTED: tt.errorsCorrected,
			gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: tt.symbologyID,
		}
		testutil.TestFile(t, reader, tt.file, tt.wants, format, tt.hints, metadata)
	}
}

func TestDataMatrixReader_DecodePureBarcode(t *testing.T) {
	reader := NewDataMatrixReader()
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_PURE_BARCODE: true,
	}

	img, _ := gozxing.NewBitMatrix(10, 10)
	bmp := testutil.NewBinaryBitmapFromBitMatrix(img)
	_, e := reader.Decode(bmp, hints)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	img.SetRegion(0, 0, 10, 10)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(img)
	_, e = reader.Decode(bmp, hints)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	img, _ = gozxing.ParseStringToBitMatrix(dm4str, "##", "  ")
	img.SetRegion(5, 5, 10, 10)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 2))
	_, e = reader.Decode(bmp, hints)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	img, _ = gozxing.ParseStringToBitMatrix(dm4str, "##", "  ")
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 2))
	result, e := reader.Decode(bmp, hints)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	expect := "Hello World"
	if r := result.GetText(); r != expect {
		t.Fatalf("Decode result=\"%v\", expect \"%v\"", r, expect)
	}
}

func TestDataMatrixReader_DecodeWithoutHints(t *testing.T) {
	reader := NewDataMatrixReader()

	img, _ := gozxing.NewBitMatrix(10, 10)
	bmp := testutil.NewBinaryBitmapFromBitMatrix(img)
	_, e := reader.DecodeWithoutHints(bmp)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	img.SetRegion(0, 0, 10, 10)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(img)
	_, e = reader.DecodeWithoutHints(bmp)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	img, _ = gozxing.NewBitMatrix(40, 40)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(img)
	_, e = reader.DecodeWithoutHints(bmp)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	img, _ = gozxing.ParseStringToBitMatrix(dm4str, "##", "  ")
	img.SetRegion(5, 5, 10, 10)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 2))
	_, e = reader.DecodeWithoutHints(bmp)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	img, _ = gozxing.ParseStringToBitMatrix(dmstr, "##", "  ")
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 4))
	result, e := reader.DecodeWithoutHints(bmp)
	if e != nil {
		t.Fatalf("Decode returns error, %v\n", e)
	}
	expect := "Testing C40"
	if r := result.GetText(); r != expect {
		t.Fatalf("Decode result=\"%v\", expect \"%v\"", r, expect)
	}
}

func TestDataMatrixReader_Decode(t *testing.T) {
	reader := NewDataMatrixReader()
	format := gozxing.BarcodeFormat_DATA_MATRIX
	pure := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_PURE_BARCODE: true,
	}

	tests := []struct {
		file     string
		wants    string
		hints    map[gozxing.DecodeHintType]interface{}
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		// testdata from zxing core/src/test/resources/blackbox/datamatrix-1/
		{
			"testdata/0123456789.png", "0123456789", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]d1",
				gozxing.ResultMetadataType_ERRORS_CORRECTED:     0,
			},
		},
		{"testdata/C40.png", "Testing C40", pure, nil},
		{"testdata/EDIFACT.png", "EDIFACTEDIFACT", pure, nil},
		{"testdata/GUID.png", "10f27ce-acb7-4e4e-a7ae-a0b98da6ed4a", nil, nil},
		{"testdata/HelloWorld_Text_L_Kaywa.png", "Hello World", pure, nil},
		{
			"testdata/HelloWorld_Text_L_Kaywa_1_error_byte.png", "Hello World", pure,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_ERRORS_CORRECTED: 1,
			},
		},
		{
			"testdata/HelloWorld_Text_L_Kaywa_2_error_byte.png", "Hello World", pure,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_ERRORS_CORRECTED: 2,
			},
		},
		{
			"testdata/HelloWorld_Text_L_Kaywa_3_error_byte.png", "Hello World", pure,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_ERRORS_CORRECTED: 3,
			},
		},
		{
			"testdata/HelloWorld_Text_L_Kaywa_4_error_byte.png", "Hello World", pure,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_ERRORS_CORRECTED: 5,
			},
		},
		{"testdata/X12.png", "X12X12X12X12", pure, nil},
		{"testdata/abcd-18x8.png", "abcde", pure, nil},
		{"testdata/abcd-26x12.png", "abcdefghijklm", pure, nil},
		{"testdata/abcd-32x8.png", "abcdef", pure, nil},
		{"testdata/abcd-36x12.png", "abcdefghijklmnopq", pure, nil},
		{"testdata/abcd-36x16.png", "abcdefghijklmnopqrstuvwxyz", pure, nil},
		{"testdata/abcd-48x16.png", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVW", pure, nil},
		{"testdata/abcd-52x52-IDAutomation.png", "abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd", pure, nil},
		{"testdata/abcd-52x52.png", "abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd", pure, nil},
		{"testdata/abcdefg-64x64.png", "" +
			"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\" +
			"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\" +
			"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\", pure, nil},
		{"testdata/abcdefg.png", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\", pure, nil},
		{"testdata/zxing_URL_L_Kayway.png", "http://code.google.com/p/zxing/", pure, nil},
		// GS1 format: FNC1 in first position
		{
			"testdata/gs1.png", "\x1d0105391538880109172611301028967\x1d211049931544602459836", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]d2",
				gozxing.ResultMetadataType_ERRORS_CORRECTED:     0,
			},
		},

		// testdata from zxing core/src/test/resources/blackbox/datamatrix-2/
		{"testdata/01.png", "http://google.com/m", nil, nil},
		{"testdata/02.png", "http://google.com/m", nil, nil},
		{"testdata/04.png", "http://google.com/m", nil, nil},
		{"testdata/05.png", "http://google.com/m", nil, nil},
		{"testdata/06.png", "http://google.com/m", nil, nil},
		{"testdata/07.png", "http://google.com/m", nil, nil},
		{"testdata/08.png", "http://google.com/m", nil, nil},
		{"testdata/09.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		{"testdata/10.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		{"testdata/12.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		{"testdata/13.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		{"testdata/14.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		{"testdata/15.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		// original zxing cannot read these too.
		// {"testdata/03.png", "http://google.com/m", nil, nil},
		// {"testdata/11.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		// {"testdata/16.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		// {"testdata/17.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},
		// {"testdata/18.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", nil, nil},

		// testdata from zxing core/src/test/resources/blackbox/datamatrix-3/
		{"testdata/3/abcd-36x20.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl", nil, nil},
		{"testdata/3/abcd-40x26.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy", pure, nil},
		{"testdata/3/abcd-44x20.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcd", pure, nil},
		{"testdata/3/abcd-48x8.png", "abcdefghijklmnopqrstuvwxy", nil, nil},
		{"testdata/3/abcd-48x22.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab", nil, nil},
		{"testdata/3/abcd-48x24.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmn", nil, nil},
		{"testdata/3/abcd-48x26.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabc", nil, nil},
		{"testdata/3/abcd-64x8.png", "abcdefghijklmnopqrstuvwxyzabcdefgh", pure, nil},
		{"testdata/3/abcd-64x12.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk", pure, nil},
		{"testdata/3/abcd-64x16.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm", pure, nil},
		{"testdata/3/abcd-64x20.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst", nil, nil},
		{"testdata/3/abcd-64x24.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcd", pure, nil},
		{"testdata/3/abcd-64x26.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrs", pure, nil},
		{"testdata/3/abcd-80x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst", pure, nil},
		{"testdata/3/abcd-88x12.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnop", pure, nil},
		{"testdata/3/abcd-96x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabc", pure, nil},
		{"testdata/3/abcd-120x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst", pure, nil},
		{"testdata/3/abcd-144x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmno", pure, nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, test.hints, test.metadata)
	}
}
