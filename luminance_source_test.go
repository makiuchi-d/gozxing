package gozxing

import (
	"testing"

	errors "golang.org/x/xerrors"
)

type testLuminanceSource struct {
	LuminanceSourceBase
	get func(x, y int) byte
}

func newTestLuminanceSource(size int) LuminanceSource {
	return &testLuminanceSource{
		LuminanceSourceBase{size, size},
		func(x, y int) byte { return byte(255 * x / (size - 1)) },
	}
}

func (this *testLuminanceSource) GetRow(y int, row []byte) ([]byte, error) {
	if y >= this.GetHeight() {
		return row, errors.Errorf("y(%d) >= height(%d)", y, this.GetHeight())
	}
	width := this.GetWidth()
	for i := 0; i < width; i++ {
		row[i] = this.get(i, y)
	}
	return row, nil
}

func (this *testLuminanceSource) GetMatrix() []byte {
	width := this.GetWidth()
	height := this.GetHeight()
	matrix := make([]byte, width*height)
	for y := 0; y < height; y++ {
		this.GetRow(y, matrix[width*y:])
	}
	return matrix
}

func (this *testLuminanceSource) Invert() LuminanceSource {
	return LuminanceSourceInvert(this)
}

func (this *testLuminanceSource) String() string {
	return LuminanceSourceString(this)
}

func TestLuminanceSource(t *testing.T) {
	s := newTestLuminanceSource(16)

	if w, h := s.GetWidth(), s.GetHeight(); w != 16 || h != 16 {
		t.Fatalf("TestLuminanceSource size = %v,%v, expect (16,16)", w, h)
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
	if _, ok := inv.(*InvertedLuminanceSource); !ok {
		t.Fatalf("Invert returns %T, expect *InvertedLuminanceSource", inv)
	}

	expect := "" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n" +
		"####++++....    \n"
	if str := s.String(); str != expect {
		t.Fatalf("s.String:\n%v\nexpect:\n%v", str, expect)
	}
}
