package gozxing

import (
	"reflect"
	"testing"
)

func makeRGBLSource(width, height int) *RGBLuminanceSource {
	pixels := make([]int, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels[y*width+x] = xy2rgb(x, y, width, height)
		}
	}
	return NewRGBLuminanceSource(width, height, pixels).(*RGBLuminanceSource)
}

func xy2rgb(x, y, width, height int) int {
	r := (255 * x * 2 / (width - 1)) & 0xff
	g := (255 * (x + y) / (width + height - 2)) & 0xff
	b := (255 * y * 2 / (height - 1)) & 0xff
	return (r << 16) + (g << 8) + b
}

func rgb2lumina(rgb int) byte {
	r := (rgb >> 16) & 0xff
	g := (rgb >> 8) & 0xff
	b := rgb & 0xff
	return byte((r + 2*g + b) / 4)
}

func TestRGBLuminanceSource_GetRow(t *testing.T) {
	src := makeRGBLSource(10, 10)

	if _, e := src.GetRow(-1, nil); e == nil {
		t.Fatalf("GetRow(-1) must be error")
	}
	if _, e := src.GetRow(10, nil); e == nil {
		t.Fatalf("GetRow(-1) must be error")
	}

	row, e := src.GetRow(0, nil)
	if e != nil {
		t.Fatalf("GetRow(0) returns error, %v", e)
	}
	// r = (255 * w * 2 / 9) & 0xff => 0, 56,113,170,226, 27, 84,140,197,254
	// g = (255 * w / 18) & 0xff    => 0, 14, 28, 42, 56, 70, 85, 99,113,127
	// b = 0
	// (r + 2*g + b) 4
	expect := []byte{0, 21, 42, 63, 84, 41, 63, 84, 105, 127}
	if !reflect.DeepEqual(row, expect) {
		t.Fatalf("GetRow(0) = %v, expect %v", row, expect)
	}

	row, e = src.GetRow(9, row)
	if e != nil {
		t.Fatalf("GetRow(0) returns error, %v", e)
	}
	// r = (255 * w * 2 / 9) & 0xff  =>   0, 56,113,170,226, 27, 84,140,197,254
	// g = (255 * (w+9) / 18) & 0xff => 127,141,155,170,184,198,212,226,240,255
	// b = 254
	expect = []byte{127, 148, 169, 191, 212, 169, 190, 211, 232, 254}
	if !reflect.DeepEqual(row, expect) {
		t.Fatalf("GetRow(0) = %v, expect %v", row, expect)
	}
}

func TestRGBLuminanceSource_GetMatrix(t *testing.T) {
	src := makeRGBLSource(10, 10)
	matrix := src.GetMatrix()
	if len(matrix) != len(src.luminances) {
		t.Fatalf("GetMatrix len=%d, expect %d", len(matrix), len(src.luminances))
	}
	for i := range matrix {
		x := i % src.dataWidth
		y := i / src.dataWidth
		if exp := rgb2lumina(xy2rgb(x, y, 10, 10)); matrix[i] != exp {
			t.Fatalf("GetMatrix matrix[%v] = %v, expect %v", i, matrix[i], exp)
		}
	}

	if !src.IsCropSupported() {
		t.Fatalf("IsCropSupported must be true")
	}

	if _, e := src.Crop(1, 0, 10, 10); e == nil {
		t.Fatalf("Crop must be error")
	}

	if _, e := src.Crop(0, 0, 10, 11); e == nil {
		t.Fatalf("Crop must be error")
	}

	crop, e := src.Crop(0, 2, 10, 6)
	if e != nil {
		t.Fatalf("Crop returns error, %v", e)
	}
	cmatrix := crop.GetMatrix()
	if !reflect.DeepEqual(cmatrix, matrix[20:len(matrix)-20]) {
		t.Fatalf("%v", crop)
	}

	crop, e = crop.Crop(2, 0, 6, 6)
	if e != nil {
		t.Fatalf("Crop returns error, %v", e)
	}
	cmatrix = crop.GetMatrix()
	expect := []byte{
		84, 105, 127, 84, 105, 126,
		105, 127, 148, 105, 127, 148,
		127, 148, 169, 126, 148, 169,
		84, 105, 126, 84, 105, 126,
		105, 127, 148, 105, 127, 148,
		126, 148, 169, 126, 148, 169,
	}
	if !reflect.DeepEqual(cmatrix, expect) {
		t.Fatalf("Cropped matrix = %v, expect %v", cmatrix, expect)
	}
}

func TestRGBLuminanceSource_Invert(t *testing.T) {
	src := makeRGBLSource(10, 10)
	inv := src.Invert()

	srcm := src.GetMatrix()
	invm := inv.GetMatrix()

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			i := y*10 + x
			exp := byte(255 - int(srcm[i]))
			if invm[i] != exp {
				t.Fatalf("Inverted matrix[%v] = %v, expect %v", i, invm[i], exp)
			}
		}
	}
}

func TestRGBLuminanceSource_String(t *testing.T) {
	src := makeRGBLSource(10, 10)
	str := src.String()
	expect := "" +
		"####+##+++\n" +
		"###++#+++.\n" +
		"##++++++..\n" +
		"#+++.++...\n" +
		"+++..+... \n" +
		"##++++++..\n" +
		"#+++.++...\n" +
		"+++..+... \n" +
		"++......  \n" +
		"+... ..   \n"
	if str != expect {
		t.Fatalf("String:\n%vexpect:\n%v", str, expect)
	}
}
