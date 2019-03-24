package gozxing

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestBitMatrix_ColorModel(t *testing.T) {
	img, _ := NewBitMatrix(10, 10)
	c := img.ColorModel()
	if c != color.GrayModel {
		t.Fatalf("ColorModel = %v, wants %v", c, color.GrayModel)
	}
}

func TestBitMatrix_Bounds(t *testing.T) {
	img, _ := NewBitMatrix(10, 20)
	rect := img.Bounds()
	expect := image.Rect(0, 0, 10, 20)
	if !reflect.DeepEqual(rect, expect) {
		t.Fatalf("Bounds = %v, wants %v", rect, expect)
	}
}

func TestBitMatrix_At(t *testing.T) {
	str := "" +
		"# # #\n" +
		" # # \n" +
		"# # #\n"
	img, _ := ParseStringToBitMatrix(str, "#", " ")

	for j := 0; j < 3; j++ {
		for i := 0; i < 5; i++ {
			c := uint32(0)
			if !img.Get(i, j) {
				c = 0xffff // white
			}
			r, g, b, a := img.At(i, j).RGBA()
			if r != c || g != c || b != c || a != 0xffff {
				t.Fatalf(
					"At(%v,%v) = RGBA(%v,%v,%v,%v), wants (%v,%v,%v,%v)",
					i, j, r, g, b, a, c, c, c, 0xffff)
			}
		}
	}
}
