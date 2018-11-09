package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestEAN13Writer_encode(t *testing.T) {
	enc := ean13Encoder{}

	_, e := enc.encode("")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("12345678901a")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("12345678901ab")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("1234567890123")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	expect := []bool{
		true, false, true, // start
		// 1: LLGLGG
		false, false, true, false, false, true, true, // 2(L)
		false, true, true, true, true, false, true, // 3(L)
		false, false, true, true, true, false, true, // 4(G)
		false, true, true, false, false, false, true, // 5(L)
		false, false, false, false, true, false, true, // 6(G)
		false, false, true, false, false, false, true, // 7(G)
		false, true, false, true, false, // middle
		true, false, false, true, false, false, false, // 8(R)
		true, true, true, false, true, false, false, // 9(R)

		true, true, true, false, false, true, false, // 0(R)
		true, true, false, false, true, true, false, // 1(R)

		true, true, false, true, true, false, false, // 2(R)
		true, false, false, true, false, false, false, // 8(R)
		true, false, true, // false
	}

	r, e := enc.encode("123456789012")
	if e != nil {
		t.Fatalf("encode must be error")
	}
	if !reflect.DeepEqual(r, expect) {
		for i := range r {
			if r[i] != expect[i] {
				t.Fatalf("differ %v, %v", i, r[i])
			}
		}

		t.Fatalf("encode:\n%v\nexpect:\n%v", r, expect)
	}

	r, e = enc.encode("1234567890128")
	if e != nil {
		t.Fatalf("encode must be error")
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode:\n%v\nexpect:\n%v", r, expect)
	}
}

func TestEAN13Writer(t *testing.T) {
	writer := NewEAN13Writer()

	if _, e := writer.Encode("5901234123abc", gozxing.BarcodeFormat_EAN_13, 0, 0, nil); e == nil {
		t.Fatalf("Encode must be error")
	}

	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"    # #  #  ## #### #  ### # ##   #    # #  #   # # # #  #   ### #  ###  # ##  ## ## ##  #  #   # #     \n"+
		"    # #  #  ## #### #  ### # ##   #    # #  #   # # # #  #   ### #  ###  # ##  ## ## ##  #  #   # #     \n",
		"#", " ")

	matrix, e := writer.Encode("123456789012", gozxing.BarcodeFormat_EAN_13, 1, 2, nil)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}

	width := matrix.GetWidth()
	height := matrix.GetHeight()
	if w, h := expect.GetWidth(), expect.GetHeight(); width != w || height != h {
		t.Fatalf("Encode matrix = %vx%v, expect %vx%v", width, height, w, h)
	}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			if m, e := matrix.Get(i, j), expect.Get(i, j); m != e {
				t.Fatalf("Encode matrix[%v,%v] = %v, expect %v", i, j, m, e)
			}
		}
	}
}
