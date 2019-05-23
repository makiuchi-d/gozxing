package gozxing

import (
	"reflect"
	"testing"
)

func makeYUVSource(w, h int) *PlanarYUVLuminanceSource {
	pixels := make([]byte, w*h*2)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			pixels[j*w+i] = byte((i + j) * 255 / (w + h - 2))
		}
	}

	src, _ := NewPlanarYUVLuminanceSource(pixels, w, h, 0, 0, w, h, false)
	return src.(*PlanarYUVLuminanceSource)
}

func TestNewPlanarYUVLuminanceSource(t *testing.T) {
	_, e := NewPlanarYUVLuminanceSource(nil, 10, 10, 5, 5, 10, 10, false)
	if e == nil {
		t.Fatalf("NewPlanarYUVLuminanceSource must be error")
	}

	src, e := NewPlanarYUVLuminanceSource(
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, 5, 2, 0, 0, 5, 2, true)
	if e != nil {
		t.Fatalf("NewPlanarYUVLuminanceSource returns error, %v", e)
	}
	yuvsrc, ok := src.(*PlanarYUVLuminanceSource)
	if !ok {
		t.Fatalf("src is not *PlanarYUVLuminanceSource, %T", src)
	}
	expect := []byte{5, 4, 3, 2, 1, 0, 9, 8, 7, 6}
	if !reflect.DeepEqual(yuvsrc.yuvData, expect) {
		t.Fatalf("yuvData = %v, expect %v", yuvsrc.yuvData, expect)
	}
}

func TestPlanarYUVLuminanceSource_GetRow(t *testing.T) {
	src := makeYUVSource(10, 10)

	if _, e := src.GetRow(-1, nil); e == nil {
		t.Fatalf("GetRow must be error")
	}
	if _, e := src.GetRow(10, nil); e == nil {
		t.Fatalf("GetRow must be error")
	}

	row, e := src.GetRow(5, nil)
	if e != nil {
		t.Fatalf("GetRow returns error, %v", e)
	}
	expect := make([]byte, 10)
	for i := 0; i < 10; i++ {
		expect[i] = byte((5 + i) * 255 / 18)
	}
	if !reflect.DeepEqual(row, expect) {
		t.Fatalf("GetRow = %v, expect %v", row, expect)
	}
}

func TestPlanarYUVLuminanceSource_IsCropSupported(t *testing.T) {
	src := makeYUVSource(10, 10)
	if !src.IsCropSupported() {
		t.Fatalf("IsCropSupported must be true")
	}
}

func TestPlanarYUVLuminanceSource_GetCroppedMatrix(t *testing.T) {
	src := makeYUVSource(10, 10)

	matrix := src.GetMatrix()
	if !reflect.DeepEqual(matrix, src.yuvData) {
		t.Fatalf("GetMatrix:\n%v\nexpect:\n%v", matrix, src.yuvData)
	}

	_, e := src.Crop(8, 8, 8, 8)
	if e == nil {
		t.Fatalf("Crop(8,8,8,8) must be error")
	}

	crop, e := src.Crop(0, 3, 10, 5)
	if e != nil {
		t.Fatalf("Crop returns error, %v", e)
	}
	matrix = crop.GetMatrix()
	expect := make([]byte, 10*5)
	for j := 0; j < 5; j++ {
		for i := 0; i < 10; i++ {
			expect[j*10+i] = byte((i + j + 3) * 255 / 18)
		}
	}
	if !reflect.DeepEqual(matrix, expect) {
		t.Fatalf("GetMatrix:\n%v\nexpect:\n%v", matrix, expect)
	}

	crop2, e := crop.Crop(3, 0, 5, 5)
	if e != nil {
		t.Fatalf("Crop returns error, %v", e)
	}
	matrix = crop2.GetMatrix()
	expect = make([]byte, 5*5)
	for j := 0; j < 5; j++ {
		for i := 0; i < 5; i++ {
			expect[j*5+i] = byte((i + j + 6) * 255 / 18)
		}
	}
	if !reflect.DeepEqual(matrix, expect) {
		t.Fatalf("GetMatrix:\n%v\nexpect:\n%v", matrix, expect)
	}
}

func TestPlanarYUVLuminanceSource_RenderThumbnail(t *testing.T) {
	src := makeYUVSource(10, 10)

	thumb := src.RenderThumbnail()
	expect := make([]uint, 5*5)
	for j := uint(0); j < 5; j++ {
		for i := uint(0); i < 5; i++ {
			y := (i + j) * 2 * 255 / 18
			expect[j*5+i] = 0xff000000 | y*0x00010101
		}
	}
	if !reflect.DeepEqual(thumb, expect) {
		t.Fatalf("\n%v\n%v", thumb, expect)
	}
}

func TestPlanarYUVLuminanceSource_Invert(t *testing.T) {
	src := makeYUVSource(10, 10)
	inv := src.Invert()

	smat := src.GetMatrix()
	imat := inv.GetMatrix()

	if len(smat) < 100 || len(imat) < 100 {
		t.Fatalf("matrix length is not enough, smat=%v, imat=%v", len(smat), len(imat))
	}
	for i := 0; i < 100; i++ {
		if imat[i] != 255-smat[i] {
			t.Fatalf("inverted matrix[%v] = %v, expect %v", i, imat[i], 255-smat[i])
		}
	}
}

func TestPlanarYUVLuminanceSource_String(t *testing.T) {
	src := makeYUVSource(10, 10)
	str := src.String()
	expect := "" +
		"#####+++++\n" +
		"####+++++.\n" +
		"###+++++..\n" +
		"##+++++...\n" +
		"#+++++....\n" +
		"+++++.... \n" +
		"++++....  \n" +
		"+++....   \n" +
		"++....    \n" +
		"+....     \n"
	if str != expect {
		t.Fatalf("String:\n%v\nexpect:\n%v", str, expect)
	}
}
