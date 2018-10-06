package testutil

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

func ExpandBitMatrix(src *gozxing.BitMatrix, factor int) *gozxing.BitMatrix {
	dst, _ := gozxing.NewBitMatrix(src.GetWidth() * factor, src.GetHeight() * factor)
	for j := 0; j < src.GetHeight(); j++ {
		y := j * factor
		for i := 0; i < src.GetWidth(); i++ {
			x := i * factor
			if src.Get(i, j) {
				dst.SetRegion(x, y, factor, factor)
			}
		}
	}
	return dst
}

func NewBinaryBitmapFromBitMatrix(matrix *gozxing.BitMatrix) *gozxing.BinaryBitmap {
	src := newTestBitMatrixSource(matrix)
	binarizer := common.NewHybridBinarizer(src)
	bmp, _ := gozxing.NewBinaryBitmap(binarizer)
	return bmp
}

func NewBinaryBitmapFromFile(filename string) *gozxing.BinaryBitmap {
	file, _ := os.Open(filename)
	img, _, _ := image.Decode(file)
	src := gozxing.NewLuminanceSourceFromImage(img)
	binarizer := common.NewHybridBinarizer(src)
	bmp, _ := gozxing.NewBinaryBitmap(binarizer)
	return bmp
}

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
