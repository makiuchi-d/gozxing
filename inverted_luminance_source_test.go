package gozxing

import (
	"reflect"
	"testing"

	errors "golang.org/x/xerrors"
)

func TestInvertedLuminanceSource(t *testing.T) {
	s := newTestLuminanceSource(16).Invert()

	expect := []byte{
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
		255, 238, 221, 204, 187, 170, 153, 136, 119, 102, 85, 68, 51, 34, 17, 0,
	}

	_, e := s.GetRow(16, make([]byte, 16))
	if e == nil {
		t.Fatalf("GetRow must be error")
	}

	row, _ := s.GetRow(1, make([]byte, 16))
	if !reflect.DeepEqual(row, expect[:16]) {
		t.Fatalf("GetRow = %v, expect %v", row, expect[:16])
	}

	matrix := s.GetMatrix()
	if !reflect.DeepEqual(matrix, expect) {
		t.Fatalf("GetMatrix = %v, expect %v", matrix, expect)
	}

	if s.IsCropSupported() {
		t.Fatalf("IsCropped is not false")
	}

	if _, e := s.Crop(1, 1, 10, 10); e == nil {
		t.Fatalf("Crop must be error")
	}

	if s.IsRotateSupported() {
		t.Fatalf("IsRotateSupported is not false")
	}

	if _, e := s.RotateCounterClockwise(); e == nil {
		t.Fatalf("RotateCounterClockwise must be error")
	}

	if _, e := s.RotateCounterClockwise45(); e == nil {
		t.Fatalf("RotateCounterClockwise45 must be error")
	}

	inv := s.Invert()
	if _, ok := inv.(*testLuminanceSource); !ok {
		t.Fatalf("Invert returns %T, expect *testLuminanceSource", inv)
	}

	strexpect := "" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n" +
		"    ....++++####\n"
	if str := s.String(); str != strexpect {
		t.Fatalf("s.String:\n%v\nexpect:\n%v", str, strexpect)
	}
}

type croppableLS struct {
	LuminanceSource
	left, top, width, height int
}

func newCroppableLS(size int) LuminanceSource {
	return &croppableLS{newTestLuminanceSource(size), 0, 0, size, size}
}
func (this *croppableLS) GetWidth() int {
	return this.width
}
func (this *croppableLS) GetHeight() int {
	return this.height
}
func (this *croppableLS) IsCropSupported() bool {
	return true
}
func (this *croppableLS) Crop(left, top, width, height int) (LuminanceSource, error) {
	if left+width > this.width || top+height > this.height {
		return nil, errors.New("invalid argument")
	}
	return &croppableLS{
		this.LuminanceSource,
		this.left + left,
		this.top + top,
		width,
		height,
	}, nil
}
func (this *croppableLS) GetRow(y int, row []byte) ([]byte, error) {
	if len(row) < this.LuminanceSource.GetWidth() {
		row = make([]byte, this.LuminanceSource.GetWidth())
	}
	row, e := this.LuminanceSource.GetRow(this.top+y, row)
	if e != nil {
		return row, e
	}
	return row[this.left : this.left+this.width], nil
}
func (this *croppableLS) GetMatrix() []byte {
	matrix := make([]byte, 0, this.width*this.height)
	row := make([]byte, this.width)
	for y := 0; y < this.height; y++ {
		row, _ = this.GetRow(y, row)
		matrix = append(matrix, row...)
	}
	return matrix
}
func (this *croppableLS) String() string {
	return LuminanceSourceString(this)
}

func TestInvertedLuminanceSourceCrop(t *testing.T) {
	s := NewInvertedLuminanceSource(newCroppableLS(16))

	if !s.IsCropSupported() {
		t.Fatalf("IsCropSupported must be true")
	}

	c, e := s.Crop(5, 5, 8, 8)
	if e != nil {
		t.Fatalf("Crop returns error, %v", e)
	}

	if w, h := c.GetWidth(), c.GetHeight(); w != 8 || h != 8 {
		t.Fatalf("Croped size = %v,%v, expect 8,8", w, h)
	}

	row, _ := c.GetRow(1, make([]byte, 16))
	expect := []byte{170, 153, 136, 119, 102, 85, 68, 51}
	if !reflect.DeepEqual(row, expect) {
		t.Fatalf("Cropped row = %v, expect %v", row, expect)
	}
}

type dummyRotateLS90 struct {
	LuminanceSource
}
type dummyRotateLS45 struct {
	LuminanceSource
}

func (this *dummyRotateLS90) IsRotateSupported() bool {
	return true
}
func (this *dummyRotateLS45) IsRotateSupported() bool {
	return true
}
func (this *dummyRotateLS90) RotateCounterClockwise() (LuminanceSource, error) {
	return &dummyRotateLS90{this}, nil
}
func (this *dummyRotateLS90) RotateCounterClockwise45() (LuminanceSource, error) {
	return &dummyRotateLS45{this}, nil
}
func (this *dummyRotateLS45) RotateCounterClockwise() (LuminanceSource, error) {
	return &dummyRotateLS90{this}, nil
}
func (this *dummyRotateLS45) RotateCounterClockwise45() (LuminanceSource, error) {
	return &dummyRotateLS45{this}, nil
}

func TestInvertedLuminanceSourceRotate(t *testing.T) {
	s := NewInvertedLuminanceSource(&dummyRotateLS45{newTestLuminanceSource(16)})

	if !s.IsRotateSupported() {
		t.Fatalf("IsRotatedSupported must be true")
	}

	is90, e := s.RotateCounterClockwise()
	if e != nil {
		t.Fatalf("RotateCounterClockwise returns error, %v", e)
	}
	s90 := is90.Invert()
	if _, ok := s90.(*dummyRotateLS90); !ok {
		t.Fatalf("RotateCounterClockwise return type Inverted{%T}, expect Inverted{*dummyRotate90}", s90)
	}

	is45, e := s.RotateCounterClockwise45()
	if e != nil {
		t.Fatalf("RotateCounterClockwise45 returns error, %v", e)
	}
	s45 := is45.Invert()
	if _, ok := s45.(*dummyRotateLS45); !ok {
		t.Fatalf("RotateCounterClockwise return type Inverted{%T}, expect Inverted{*dummyRotate45}", s45)
	}
}
