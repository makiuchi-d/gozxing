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

func testDataMatrixReader_Decode(t testing.TB, file, expect string, pure bool) {
	t.Helper()
	reader := NewDataMatrixReader()
	bmp := testutil.NewBinaryBitmapFromFile(file)
	var help map[gozxing.DecodeHintType]interface{}
	if pure {
		help = make(map[gozxing.DecodeHintType]interface{})
		help[gozxing.DecodeHintType_PURE_BARCODE] = true
	}
	result, e := reader.Decode(bmp, help)
	if e != nil {
		t.Fatalf("Decode(%v) returns error, %v", file, e)
	}
	if txt := result.GetText(); txt != expect {
		t.Fatalf("Decode(%v) = \"%v\", expect \"%v\"", file, txt, expect)
	}
}

func TestDataMatrixReader_Decode(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/datamatrix-1/
	testDataMatrixReader_Decode(t, "testdata/0123456789.png", "0123456789", false)
	testDataMatrixReader_Decode(t, "testdata/C40.png", "Testing C40", true)
	testDataMatrixReader_Decode(t, "testdata/EDIFACT.png", "EDIFACTEDIFACT", true)
	testDataMatrixReader_Decode(t, "testdata/GUID.png", "10f27ce-acb7-4e4e-a7ae-a0b98da6ed4a", false)
	testDataMatrixReader_Decode(t, "testdata/HelloWorld_Text_L_Kaywa.png", "Hello World", true)
	testDataMatrixReader_Decode(t, "testdata/HelloWorld_Text_L_Kaywa_1_error_byte.png", "Hello World", true)
	testDataMatrixReader_Decode(t, "testdata/HelloWorld_Text_L_Kaywa_2_error_byte.png", "Hello World", true)
	testDataMatrixReader_Decode(t, "testdata/HelloWorld_Text_L_Kaywa_3_error_byte.png", "Hello World", true)
	testDataMatrixReader_Decode(t, "testdata/HelloWorld_Text_L_Kaywa_4_error_byte.png", "Hello World", true)
	testDataMatrixReader_Decode(t, "testdata/X12.png", "X12X12X12X12", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-18x8.png", "abcde", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-26x12.png", "abcdefghijklm", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-32x8.png", "abcdef", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-36x12.png", "abcdefghijklmnopq", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-36x16.png", "abcdefghijklmnopqrstuvwxyz", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-48x16.png", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVW", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-52x52-IDAutomation.png",
		"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd", true)
	testDataMatrixReader_Decode(t, "testdata/abcd-52x52.png",
		"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd", true)
	testDataMatrixReader_Decode(t, "testdata/abcdefg-64x64.png", ""+
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\"+
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\"+
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\", true)
	testDataMatrixReader_Decode(t, "testdata/abcdefg.png", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*(),./\\", true)
	testDataMatrixReader_Decode(t, "testdata/zxing_URL_L_Kayway.png", "http://code.google.com/p/zxing/", true)

	// testdata from zxing core/src/test/resources/blackbox/datamatrix-2/
	testDataMatrixReader_Decode(t, "testdata/01.png", "http://google.com/m", false)
	testDataMatrixReader_Decode(t, "testdata/02.png", "http://google.com/m", false)
	testDataMatrixReader_Decode(t, "testdata/04.png", "http://google.com/m", false)
	testDataMatrixReader_Decode(t, "testdata/05.png", "http://google.com/m", false)
	testDataMatrixReader_Decode(t, "testdata/06.png", "http://google.com/m", false)
	testDataMatrixReader_Decode(t, "testdata/07.png", "http://google.com/m", false)
	testDataMatrixReader_Decode(t, "testdata/08.png", "http://google.com/m", false)
	testDataMatrixReader_Decode(t, "testdata/09.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	testDataMatrixReader_Decode(t, "testdata/10.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	testDataMatrixReader_Decode(t, "testdata/12.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	testDataMatrixReader_Decode(t, "testdata/13.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	testDataMatrixReader_Decode(t, "testdata/14.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	testDataMatrixReader_Decode(t, "testdata/15.png", "This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)

	// original zxing cannot read these too.
	//testDataMatrixReader_Decode(t, "testdata/03.png", "http://google.com/m", false)
	//testDataMatrixReader_Decode(t, "testdata/11.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	//testDataMatrixReader_Decode(t, "testdata/16.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	//testDataMatrixReader_Decode(t, "testdata/17.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)
	//testDataMatrixReader_Decode(t, "testdata/18.png","This is a test of our DataMatrix support using a longer piece of text, and therefore a more dense barcode.", false)

	// testdata from zxing core/src/test/resources/blackbox/datamatrix-3/
	testDataMatrixReader_Decode(t, "testdata/3/abcd-36x20.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl", false)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-40x26.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxy", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-44x20.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcd", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-48x8.png", "abcdefghijklmnopqrstuvwxy", false)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-48x22.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab", false)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-48x24.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmn", false)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-48x26.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabc", false)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-64x8.png", "abcdefghijklmnopqrstuvwxyzabcdefgh", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-64x12.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-64x16.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-64x20.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst", false)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-64x24.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcd", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-64x26.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrs", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-80x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-88x12.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnop", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-96x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabc", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-120x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst", true)
	testDataMatrixReader_Decode(t, "testdata/3/abcd-144x8.png", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmno", true)
}
