package gozxing

import (
	"image"
	"image/color"
	"image/draw"
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

func TestNewBinaryBitmapFromImage(t *testing.T) {
	img := newTestImage(15, 10)
	bmp, e := NewBinaryBitmapFromImage(img)
	if e != nil {
		t.Fatalf("NewBinaryBitmapFromImage returns error: %v", e)
	}
	if w, h := bmp.GetWidth(), bmp.GetHeight(); w != 15 || h != 10 {
		t.Fatalf("NewBinaryBitmapFromImage = %vx%v, expect 15x10", w, h)
	}
}

func TestNewLuminanceSourceFromImage(t *testing.T) {
	img := newTestImage(10, 10)
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	imgs := []image.Image{img, rgba}
	for _, img := range imgs {

		src := NewLuminanceSourceFromImage(img)

		if _, ok := src.(*GoImageLuminanceSource); !ok {
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
}

func TestNewLuminanceSourceFromGray(t *testing.T) {
	img := newTestImage(10, 10)
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Src)
	src := NewLuminanceSourceFromImage(gray)

	if _, ok := src.(*GoImageLuminanceSource); !ok {
		t.Fatalf("NewLuminanceSourceFromImage must return *RGBLuminanceSource, %T", src)
	}

	matrix := src.GetMatrix()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			expect := gray.GrayAt(x, y).Y
			lumina := matrix[y*10+x]
			if lumina != expect {
				t.Fatalf("matrix[%v,%v] = %v, expect %v", x, y, lumina, expect)
			}
		}
	}
}

func TestGoImageLuminanceSource_Crop(t *testing.T) {
	img := newTestImage(20, 20)
	src := NewLuminanceSourceFromImage(img)

	_, e := src.Crop(10, 10, 20, 20)
	if e == nil {
		t.Fatalf("Crop must be error")
	}

	cropped, e := src.Crop(5, 5, 10, 10)
	if e != nil {
		t.Fatalf("Crop returns error, %v", e)
	}

	matrix := cropped.GetMatrix()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			expect := rgb2lumina(xy2rgb(x+5, y+5, 20, 20))
			lumina := matrix[y*10+x]
			if lumina != expect {
				t.Fatalf("matrix[%v,%v] = %v, expect %v", x, y, lumina, expect)
			}
		}
	}
}

func TestGoImageLuminanceSource_Invert(t *testing.T) {
	img := newTestImage(10, 10)
	src := NewLuminanceSourceFromImage(img)
	inverted := src.Invert()

	matrix := inverted.GetMatrix()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			expect := 255 - rgb2lumina(xy2rgb(x, y, 10, 10))
			lumina := matrix[y*10+x]
			if lumina != expect {
				t.Fatalf("matrix[%v,%v] = %v, expect %v", x, y, lumina, expect)
			}
		}
	}
}

func TestGoImageLuminanceSource_RotateCounterClockwise(t *testing.T) {
	img := newTestImage(20, 20)
	src, _ := NewLuminanceSourceFromImage(img).Crop(5, 3, 10, 15)

	if !src.IsRotateSupported() {
		t.Fatalf("IsRotateSupported must be true")
	}

	rotated, e := src.RotateCounterClockwise()
	if e != nil {
		t.Fatalf("RotateCounterClockwise returns error, %v", e)
	}
	if w, h := rotated.GetWidth(), rotated.GetHeight(); w != 15 || h != 10 {
		t.Fatalf("roated size = %v,%v, expect 15,10", w, h)
	}

	matrix := rotated.GetMatrix()
	for y := 0; y < 10; y++ {
		for x := 0; x < 15; x++ {
			oldx := 5 + 10 - 1 - y
			oldy := 3 + x
			expect := rgb2lumina(xy2rgb(oldx, oldy, 20, 20))
			lumina := matrix[y*15+x]
			if lumina != expect {
				t.Fatalf("matrix[%v,%v] = %v, expect %v", x, y, lumina, expect)
			}
		}
	}
}

func TestGoImageLuminanceSource_RotateCounterClockwise45(t *testing.T) {
	img := newTestImage(10, 10)
	src := NewLuminanceSourceFromImage(img)

	_, e := src.RotateCounterClockwise45()
	if e == nil {
		t.Fatalf("RotateCounterClockwise45 must be error")
	}
}
