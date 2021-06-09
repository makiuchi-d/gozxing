package qrcode

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestNewQRCodeReader(t *testing.T) {
	reader := NewQRCodeReader().(*QRCodeReader)
	if reader.GetDecoder() == nil {
		t.Fatalf("decoder must not be nil")
	}
}

func TestQRCodeReader_DecodeBitMatrixSource(t *testing.T) {
	reader := NewQRCodeReader()
	var matrix *gozxing.BitMatrix
	var bmp *gozxing.BinaryBitmap
	var e error

	qrstr := "" +
		"                                                          \n" +
		"                                                          \n" +
		"                                                          \n" +
		"                                                          \n" +
		"        ##############      ##  ##  ##############        \n" +
		"        ##          ##          ##  ##          ##        \n" +
		"        ##  ######  ##  ##  ##      ##  ######  ##        \n" +
		"        ##  ######  ##          ##  ##  ######  ##        \n" +
		"        ##  ######  ##    ##  ####  ##  ######  ##        \n" +
		"        ##          ##    ######    ##          ##        \n" +
		"        ##############  ##  ##  ##  ##############        \n" +
		"                        ##  ##                            \n" +
		"        ######  ##########  ##  ######      ##            \n" +
		"          ##  ##        ########  ##  ##      ####        \n" +
		"        ##    ####  ##  ########  ######  ########        \n" +
		"            ####  ##  ####    ######  ####    ##          \n" +
		"                ##  ##    ##  ##  ######                  \n" +
		"                        ##  ##      ####    ######        \n" +
		"        ##############  ##  ##  ##      ##  ######        \n" +
		"        ##          ##  ######      ######    ####        \n" +
		"        ##  ######  ##  ####    ##  ##        ####        \n" +
		"        ##  ######  ##    ######  ##  ##    ####          \n" +
		"        ##  ######  ##  ########  ####  ##  ##  ##        \n" +
		"        ##          ##  ##  ########    ##    ##          \n" +
		"        ##############  ########  ######      ####        \n" +
		"                                                          \n" +
		"                                                          \n" +
		"                                                          \n" +
		"                                                          \n"
	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_PURE_BARCODE] = true

	// fail GetBitMatrixMatrix
	matrix, _ = gozxing.NewSquareBitMatrix(1)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(matrix)
	_, e = reader.DecodeWithoutHints(bmp)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Decode must be NotFoundException, %T", e)
	}

	// fail Detect
	matrix, _ = gozxing.NewSquareBitMatrix(50)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(matrix)
	_, e = reader.DecodeWithoutHints(bmp)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Decode must be NotFoundException, %T", e)
	}
	_, e = reader.Decode(bmp, hints)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Decode must be NotFoundException, %T", e)
	}

	// fail checksum
	matrix, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	matrix.SetRegion(15, 15, 10, 10)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(matrix)
	_, e = reader.DecodeWithoutHints(bmp)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("Decode must be ChecksumException, %T", e)
	}
	_, e = reader.Decode(bmp, hints)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("Decode must be ChecksumException, %T", e)
	}

	// success
	matrix, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	bmp = testutil.NewBinaryBitmapFromBitMatrix(matrix)
	textexpect := "hello\n"
	r, e := reader.DecodeWithoutHints(bmp)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if txt := r.GetText(); txt != textexpect {
		t.Fatalf("Decode text = \"%v\", expect \"%v\"", txt, textexpect)
	}
	r, e = reader.Decode(bmp, hints)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if txt := r.GetText(); txt != textexpect {
		t.Fatalf("Decode text = \"%v\", expect \"%v\"", txt, textexpect)
	}
}

func TestQRCodeReader_moduleSize(t *testing.T) {
	reader := NewQRCodeReader().(*QRCodeReader)
	matrix, _ := gozxing.ParseStringToBitMatrix(""+
		"                 \n"+
		" ##############  \n"+
		" ##############  \n"+
		" ##          ##  \n"+
		" ##          ##  \n"+
		" ##  ######  ##  \n"+
		" ##  ######  ##  \n"+
		" ##  ######  ##  \n"+
		" ##  ######  ##  \n"+
		" ##  ######  ##  \n"+
		" ##  ######  ##  \n"+
		" ##          ##  \n"+
		" ##          ##  \n"+
		" ##############  \n"+
		" ##############  \n"+
		"                 \n"+
		"                 \n", "#", " ")
	s, e := reader.moduleSize([]int{1, 1}, matrix)
	if e != nil {
		t.Fatalf("moduleSize returns error, %v", e)
	}
	if s != 2 {
		t.Fatalf("moduleSize = %v, expect 2", s)
	}

	matrix.SetRegion(0, 0, matrix.GetWidth(), matrix.GetHeight())
	s, e = reader.moduleSize([]int{1, 1}, matrix)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("moduleSize must be NotFoundException, %T", e)
	}
}

func TestQRCodeReader_extractPureBits(t *testing.T) {
	reader := NewQRCodeReader().(*QRCodeReader)
	matrix, _ := gozxing.NewBitMatrix(32, 32)
	var r *gozxing.BitMatrix
	var e error

	_, e = reader.extractPureBits(matrix)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("moduleSize must be NotFoundException, %T", e)
	}

	matrix.SetRegion(5, 5, 20, 20)
	_, e = reader.extractPureBits(matrix)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("moduleSize must be NotFoundException, %T", e)
	}

	matrix.Clear()
	matrix.SetRegion(20, 5, 7, 1)
	matrix.SetRegion(26, 6, 1, 6)
	matrix.SetRegion(5, 7, 20, 14)
	matrix.Set(15, 21)

	_, e = reader.extractPureBits(matrix)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("moduleSize must be NotFoundException, %T", e)
	}

	matrix.Unset(15, 21)
	_, e = reader.extractPureBits(matrix)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("moduleSize must be NotFoundException, %T", e)
	}

	matrix.Clear()
	matrix.Set(5, 5)
	matrix.SetRegion(7, 7, 3, 3)
	matrix.Set(11, 11)
	r, e = reader.extractPureBits(matrix)
	if e != nil {
		t.Fatalf("extractPureBits returns error, %v", e)
	}
	if w, h := r.GetWidth(), r.GetHeight(); w != 7 || h != 7 {
		t.Fatalf("extracted size = %v,%v, expect 7,7", w, h)
	}
	expect := "" +
		"X      \n" +
		"       \n" +
		"  XXX  \n" +
		"  XXX  \n" +
		"  XXX  \n" +
		"       \n" +
		"      X\n"
	if str := r.ToString("X", " "); str != expect {
		t.Fatalf("extracted:\n%vexpect:\n%v", str, expect)
	}
}

func TestQRCodeReader_Reset(t *testing.T) {
	reader := NewQRCodeReader()
	reader.Reset() // this method do nothing
}

func testDecodeImage(t testing.TB, file, expect string) *gozxing.Result {
	t.Helper()
	bmp := testutil.NewBinaryBitmapFromFile(file)
	r, e := NewQRCodeReader().Decode(bmp, nil)
	if e != nil {
		t.Fatalf("Decode(%s) returns error, %v", file, e)
	}

	if f := r.GetBarcodeFormat(); f != gozxing.BarcodeFormat_QR_CODE {
		t.Fatalf("Decode(%s) format is not QR_COODE, %v", file, f)
	}

	if txt := r.GetText(); txt != expect {
		t.Fatalf("Decode(%s) text = \n\"%v\", expect \n\"%v\"", file, txt, expect)
	}
	return r
}

func testStructuredAppend(t testing.TB, file string, metadata map[gozxing.ResultMetadataType]interface{}, seq, parity int) {
	t.Helper()
	s, ok := metadata[gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE]
	if !ok {
		t.Fatalf("Decode(%s) ResultMetadata must have STRUCTURED_APPEND_SEQUENCE", file)
	}
	if s != seq {
		t.Fatalf("Decode(%s) STRUCTURED_APPEND_SEQUENCE = 0x%02x, expect 0x%02x", file, s, seq)
	}

	p, ok := metadata[gozxing.ResultMetadataType_STRUCTURED_APPEND_PARITY]
	if !ok {
		t.Fatalf("Decode(%s) ResultMetadata must have STRUCTURED_APPEND_PARITY", file)
	}
	if p != parity {
		t.Fatalf("Decode(%s) STRUCTURED_APPEND_PARITY = 0x%02x, expect 0x%02x", file, p, parity)
	}
}

func TestQRCodeReader_DecodeImage(t *testing.T) {
	var file string
	var result *gozxing.Result

	// ISO/IEC 18004:2000 Figure 1
	file = "testdata/version1.png"
	result = testDecodeImage(t, file, "QR Code Symbol")
	points := result.GetResultPoints()
	// [0]: bottom-left; [1]: top-left; [2]: top-right
	if points[0].GetX() < points[1].GetX()-10 || points[1].GetX()+10 < points[0].GetX() {
		t.Fatalf("ResultPoint BottomLeft.X != TopLeft.X")
	}
	if points[2].GetY() < points[1].GetY()-10 || points[1].GetY()+10 < points[2].GetY() {
		t.Fatalf("ResultPoint TopRight.Y != TopLeft.Y")
	}
	file = "testdata/version1_mirrored.png"
	result = testDecodeImage(t, file, "QR Code Symbol")
	points = result.GetResultPoints()
	if points[0].GetX() < points[1].GetX()-10 || points[1].GetX()+10 < points[0].GetX() {
		t.Fatalf("ResultPoint BottomLeft.X != TopLeft.X")
	}
	if points[2].GetY() < points[1].GetY()-10 || points[1].GetY()+10 < points[2].GetY() {
		t.Fatalf("ResultPoint TopRight.Y != TopLeft.Y")
	}

	// https://github.com/zxing/zxing/tree/master/core/src/test/resources/benchmark/android-1
	file = "testdata/qrcode-1.jpg"
	testDecodeImage(t, file, "MECARD:N:Google 411,;TEL:18665881077;;")

	file = "testdata/qrcode-2.jpg"
	testDecodeImage(t, file,
		"UI office hours signup\r\nhttp://www.corp.google.com/sparrow/ui_office_hours/ \r\n")

	file = "testdata/qrcode-3.jpg"
	testDecodeImage(t, file, "MECARD:N:Sean Owen;TEL:+12125658770;EMAIL:srowen@google.com;;")

	file = "testdata/qrcode-4.jpg"
	testDecodeImage(t, file, "MECARD:N:Sean Owen;TEL:+12125658770;EMAIL:srowen@google.com;;")

	// structured append
	// ISO/IEC 18004:2000 Figure 22
	file = "testdata/structured_append_1.png"
	result = testDecodeImage(t, file, "ABCDEFGHIJKLMN")
	testStructuredAppend(t, file, result.GetResultMetadata(), 0x03, 1) // 1 of 4

	file = "testdata/structured_append_2.png"
	result = testDecodeImage(t, file, "OPQRSTUVWXYZ0123")
	testStructuredAppend(t, file, result.GetResultMetadata(), 0x13, 1) // 2 of 4

	file = "testdata/structured_append_3.png"
	result = testDecodeImage(t, file, "456789ABCDEFGHIJ")
	testStructuredAppend(t, file, result.GetResultMetadata(), 0x23, 1) // 3 of 4

	file = "testdata/structured_append_4.png"
	result = testDecodeImage(t, file, "KLMNOPQRSTUVWXYZ")
	testStructuredAppend(t, file, result.GetResultMetadata(), 0x33, 1) // 4 of 4
}

func TestQRCodeReaderBlackbox(t *testing.T) {
	reader := NewQRCodeReader()
	format := gozxing.BarcodeFormat_QR_CODE

	tests := []struct {
		file     string
		wants    string
		hints    map[gozxing.DecodeHintType]interface{}
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		// testdata from zxing core/src/test/resources/blackbox/qrcode-1
		{
			"testdata/qrcode-1/1.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]Q1",
			},
		},
		{"testdata/qrcode-1/2.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/3.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/4.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/5.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/6.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/7.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/8.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/9.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/10.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/11.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/12.png", "MEBKM:URL:http\\://en.wikipedia.org/wiki/Main_Page;;", nil, nil},
		{"testdata/qrcode-1/13.png", "http://google.com/gwt/n?u=bluenile.com", nil, nil},
		//{"testdata/qrcode-1/14.png", "http://google.com/gwt/n?u=bluenile.com", nil, nil},
		{"testdata/qrcode-1/15.png", "http://google.com/gwt/n?u=bluenile.com", nil, nil},
		{"testdata/qrcode-1/16.png", "Sean Owen\r\nsrowen@google.com\r\n917-364-2918\r\nhttp://awesome-thoughts.com", nil, nil},
		{"testdata/qrcode-1/17.png", "Sean Owen\r\nsrowen@google.com\r\n917-364-2918\r\nhttp://awesome-thoughts.com", nil, nil},
		{"testdata/qrcode-1/18.png", "Sean Owen\r\nsrowen@google.com\r\n917-364-2918\r\nhttp://awesome-thoughts.com", nil, nil},
		{"testdata/qrcode-1/19.png", "Sean Owen\r\nsrowen@google.com\r\n917-364-2918\r\nhttp://awesome-thoughts.com", nil, nil},
		{"testdata/qrcode-1/20.png", "Sean Owen\r\nsrowen@google.com\r\n917-364-2918\r\nhttp://awesome-thoughts.com", nil, nil},

		// testdata from zxing core/src/test/resources/blackbox/qrcode-6
		{
			"testdata/qrcode-6/1.png", "1234567890", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]Q1",
			},
		},
		{"testdata/qrcode-6/2.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/3.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/4.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/5.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/6.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/7.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/8.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/9.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/10.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/11.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/12.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/13.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/14.png", "1234567890", nil, nil},
		{"testdata/qrcode-6/15.png", "TEST", nil, nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, test.hints, test.metadata)
	}
}
