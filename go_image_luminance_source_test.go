package gozxing

import (
	"image"
	"image/color"
	"testing"
)

type testImage struct {
	rect image.Rectangle
}

func (this *testImage) ColorModel() color.Model {
	return color.RGBAModel
}
func (this *testImage) Bounds() image.Rectangle {
	return this.rect
}
func (this *testImage) At(x, y int) color.Color {
	width := this.rect.Max.X - this.rect.Min.X
	height := this.rect.Max.Y - this.rect.Min.Y
	rgb := xy2rgb(x, y, width, height)
	r := byte((rgb >> 16) & 0xff)
	g := byte((rgb >> 8) & 0xff)
	b := byte(rgb & 0xff)
	return color.RGBA{r, g, b, 0xff}
}

func newTestImage(w, h int) *testImage {
	return &testImage{image.Rect(0, 0, w, h)}
}

func TestNewLuminanceSourceFromImage(t *testing.T) {
	img := newTestImage(10, 10)
	src := NewLuminanceSourceFromImage(img)

	if _, ok := src.(*RGBLuminanceSource); !ok {
		t.Fatalf("NewLuminanceSourceFromImage must return *RGBLuminanceSource, %T", src)
	}

	matrix := src.GetMatrix()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			expect := rgb2lumina(xy2rgb(x, y, 10, 10))
			lumina := matrix[y*10+x]
			if lumina != expect {
				t.Fatalf("matrix[%v,%v] = %v, expect %v", x, y, lumina, expect)
			}
		}
	}
}
