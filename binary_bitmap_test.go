package gozxing

import (
	"testing"

	errors "golang.org/x/xerrors"
)

type testBinarizer struct {
	source LuminanceSource
}

func (this *testBinarizer) GetLuminanceSource() LuminanceSource {
	return this.source
}
func (this *testBinarizer) GetBlackRow(y int, row *BitArray) (*BitArray, error) {
	width := this.GetWidth()
	if row.GetSize() < width {
		row = NewBitArray(width)
	} else {
		row.Clear()
	}
	rawrow, e := this.source.GetRow(y, make([]byte, width))
	if e != nil {
		return row, e
	}
	for i, v := range rawrow {
		if v < 128 {
			row.Set(i)
		}
	}
	return row, nil
}
func (this *testBinarizer) GetBlackMatrix() (*BitMatrix, error) {
	width := this.GetWidth()
	height := this.GetHeight()
	matrix, _ := NewBitMatrix(width, height)
	row := NewBitArray(width)
	for y := 0; y < height; y++ {
		row, _ = this.GetBlackRow(y, row)
		for x := 0; x < width; x++ {
			if row.Get(x) {
				matrix.Set(x, y)
			}
		}
	}
	return matrix, nil
}
func (this *testBinarizer) CreateBinarizer(source LuminanceSource) Binarizer {
	return &testBinarizer{source}
}
func (this *testBinarizer) GetWidth() int {
	return this.source.GetWidth()
}
func (this *testBinarizer) GetHeight() int {
	return this.source.GetHeight()
}

func TestBinaryBitmap(t *testing.T) {
	if _, e := NewBinaryBitmap(nil); e == nil {
		t.Fatalf("NewBinaryBitmap must be error")
	}

	binarizer := &testBinarizer{newTestLuminanceSource(16)}
	bmp, e := NewBinaryBitmap(binarizer)
	if e != nil {
		t.Fatalf("NewBinaryBitmap return error, %v", e)
	}
	if w, h := bmp.GetWidth(), bmp.GetHeight(); w != 16 || h != 16 {
		t.Fatalf("width,height = %v,%v, expect 16,16", w, h)
	}

	arr, e := bmp.GetBlackRow(0, NewBitArray(16))
	if e != nil {
		t.Fatalf("GetBlackRow returns error, %v", e)
	}
	for i := 0; i < arr.GetSize(); i++ {
		if arr.Get(i) != (i < 8) {
			t.Fatalf("BlackRow.Get(%v) = %v, expect %v", i, arr.Get(i), (i < 8))
		}
	}

	matrix, e := bmp.GetBlackMatrix()
	if e != nil {
		t.Fatalf("GetBlackRow returns error, %v", e)
	}
	for y := 0; y < matrix.GetHeight(); y++ {
		for x := 0; x < matrix.GetWidth(); x++ {
			if matrix.Get(x, y) != (x < 8) {
				t.Fatalf("BlackMatrix.Get(%v,%v) = %v, expect %v", x, y, matrix.Get(x, y), (x < 8))
			}
		}
	}

	if bmp.IsCropSupported() {
		t.Fatalf("IsCropSupported must false")
	}
	if _, e := bmp.Crop(1, 1, 3, 3); e == nil {
		t.Fatalf("Crop must be error")
	}
	if bmp.IsRotateSupported() {
		t.Fatalf("IsRotateSupported must false")
	}
	if _, e := bmp.RotateCounterClockwise(); e == nil {
		t.Fatalf("RotateCounterClockwise must be error")
	}
	if _, e := bmp.RotateCounterClockwise45(); e == nil {
		t.Fatalf("RotateCounterClockwise45 must be error")
	}

	expectStr := "" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n" +
		"X X X X X X X X                 \n"
	if bmp.String() != expectStr {
		t.Fatalf("\n%v", bmp)
	}
}

type notFoundBinarizer struct {
	testBinarizer
}

func (this *notFoundBinarizer) GetBlackMatrix() (*BitMatrix, error) {
	return nil, NewNotFoundException()
}

type illegalBinarizer struct {
	testBinarizer
}

func (this *illegalBinarizer) GetBlackMatrix() (*BitMatrix, error) {
	return nil, errors.New("IllegalException")
}

func TestBinaryBitmap_NotFound(t *testing.T) {
	var binarizer Binarizer = &notFoundBinarizer{testBinarizer{newTestLuminanceSource(16)}}
	bmp, _ := NewBinaryBitmap(binarizer)

	if _, e := bmp.GetBlackMatrix(); e == nil {
		t.Fatalf("GetBlackMatrix() must be error")
	}
	if s := bmp.String(); s != "" {
		t.Fatalf("Bitmap string = \"%v\", expect \"\"", s)
	}

	bmp, _ = NewBinaryBitmap(&testBinarizer{newCroppableLS(16)})
	if !bmp.IsCropSupported() {
		t.Fatalf("IsCropSupported must true")
	}
	bmp2, e := bmp.Crop(3, 3, 8, 8)
	if e != nil {
		t.Fatalf("Crop returns error, %v", e)
	}
	if w, h := bmp2.GetWidth(), bmp2.GetHeight(); w != 8 || h != 8 {
		t.Fatalf("cropped size = %v,%v, expect 8,8", w, h)
	}

	bmp, _ = NewBinaryBitmap(&testBinarizer{&dummyRotateLS90{newTestLuminanceSource(16)}})
	if !bmp.IsRotateSupported() {
		t.Fatalf("IsRotateSupported must true")
	}
	bmp2, e = bmp.RotateCounterClockwise()
	if e != nil {
		t.Fatalf("RotateCounterClockwise returns error, %v", e)
	}
	src := bmp2.binarizer.GetLuminanceSource()
	if _, ok := src.(*dummyRotateLS90); !ok {
		t.Fatalf("rotated source type must be *dummyRotateLS90, %T", src)
	}
	bmp2, e = bmp.RotateCounterClockwise45()
	if e != nil {
		t.Fatalf("RotateCounterClockwise45 returns error, %v", e)
	}
	src = bmp2.binarizer.GetLuminanceSource()
	if _, ok := src.(*dummyRotateLS45); !ok {
		t.Fatalf("rotated source type must be *dummyRotateLS45, %T", src)
	}
}

func TestBinaryBitmap_Illegal(t *testing.T) {
	var binarizer Binarizer = &illegalBinarizer{testBinarizer{newTestLuminanceSource(16)}}
	bmp, _ := NewBinaryBitmap(binarizer)
	expect := "IllegalException"
	if s := bmp.String(); s != expect {
		t.Fatalf("string = \"%v\", expect \"%v\"", s, expect)
	}
}
