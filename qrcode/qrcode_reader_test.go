package qrcode

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

type testBitMatrixSource struct {
	gozxing.LuminanceSourceBase
	matrix *gozxing.BitMatrix
}

func newTestBitMatrixSource(matrix *gozxing.BitMatrix) gozxing.LuminanceSource {
	return &testBitMatrixSource{
		gozxing.LuminanceSourceBase{matrix.GetWidth(), matrix.GetHeight()},
		matrix,
	}
}
func (this *testBitMatrixSource) GetRow(y int, row []byte) ([]byte, error) {
	for x := 0; x < this.matrix.GetWidth(); x++ {
		if this.matrix.Get(x, y) {
			row[x] = 0
		} else {
			row[x] = 255
		}
	}
	return row, nil
}
func (this *testBitMatrixSource) GetMatrix() []byte {
	width := this.GetWidth()
	height := this.GetHeight()
	matrix := make([]byte, width*height)
	for y := 0; y < height; y++ {
		offset := y * width
		for x := 0; x < width; x++ {
			if !this.matrix.Get(x, y) {
				matrix[offset+x] = 255
			}
		}
	}
	return matrix
}
func (this *testBitMatrixSource) Invert() gozxing.LuminanceSource {
	return gozxing.LuminanceSourceInvert(this)
}
func (this *testBitMatrixSource) String() string {
	return gozxing.LuminanceSourceString(this)
}

func TestNewQRCodeReader(t *testing.T) {
	reader := NewQRCodeReader().(*QRCodeReader)
	if reader.GetDecoder() == nil {
		t.Fatalf("decoder must not be nil")
	}
}

func TestQRCodeReader_DecodeBitMatrixSource(t *testing.T) {
	reader := NewQRCodeReader()
	var matrix *gozxing.BitMatrix
	var src gozxing.LuminanceSource
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
	src = newTestBitMatrixSource(matrix)
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
	_, e = reader.DecodeWithoutHints(bmp)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Decode must be NotFoundException, %T", e)
	}

	// fail Detect
	matrix, _ = gozxing.NewSquareBitMatrix(50)
	src = newTestBitMatrixSource(matrix)
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
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
	src = newTestBitMatrixSource(matrix)
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
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
	src = newTestBitMatrixSource(matrix)
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
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

type testImageSource struct {
	gozxing.LuminanceSourceBase
	img       image.Image
	top, left int
}

func newTestImageSource(filename string) gozxing.LuminanceSource {
	file, _ := os.Open(filename)
	defer file.Close()
	img, _, _ := image.Decode(file)
	rect := img.Bounds()
	top := rect.Min.Y
	left := rect.Min.X
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	return &testImageSource{
		gozxing.LuminanceSourceBase{width, height},
		img,
		top,
		left,
	}
}
func (this *testImageSource) GetRow(y int, row []byte) ([]byte, error) {
	if y >= this.GetHeight() {
		return row, fmt.Errorf("y(%d) >= height(%d)", y, this.GetHeight())
	}
	for x := 0; x < this.GetWidth(); x++ {
		r, g, b, _ := this.img.At(this.left+x, this.top+y).RGBA()
		row[x] = byte((r + 2*g + b) * 255 / (4 * 0xffff))
	}
	return row, nil
}
func (this *testImageSource) GetMatrix() []byte {
	width := this.GetWidth()
	height := this.GetHeight()
	matrix := make([]byte, width*height)
	for y := 0; y < height; y++ {
		offset := y * width
		for x := 0; x < width; x++ {
			r, g, b, _ := this.img.At(this.left+x, this.top+y).RGBA()
			matrix[offset+x] = byte((r + 2*g + b) * 255 / (4 * 0xffff))
		}
	}
	return matrix
}
func (this *testImageSource) Invert() gozxing.LuminanceSource {
	return gozxing.LuminanceSourceInvert(this)
}
func (this *testImageSource) String() string {
	return gozxing.LuminanceSourceString(this)
}

func testDecodeImage(t *testing.T, file, expect string) *gozxing.Result {
	src := newTestImageSource(file)
	bmp, _ := gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
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

func testStructuredAppend(t *testing.T, file string, metadata map[gozxing.ResultMetadataType]interface{}, seq, parity int) {
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
