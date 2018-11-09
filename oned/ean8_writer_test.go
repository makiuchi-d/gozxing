package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestEAN8Writer_encode(t *testing.T) {
	enc := ean8Encoder{}

	_, e := enc.encode("")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("123456a")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("123456ab")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("12345678")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	expect := []bool{
		true, false, true, // start
		false, false, true, true, false, false, true, // 1
		false, false, true, false, false, true, true, // 2
		false, true, true, true, true, false, true, // 3
		false, true, false, false, false, true, true, // 4
		false, true, false, true, false, // middle
		true, false, false, true, true, true, false, // 5
		true, false, true, false, false, false, false, // 6
		true, false, false, false, true, false, false, // 7
		true, true, true, false, false, true, false, // 0
		true, false, true, // end
	}

	r, e := enc.encode("1234567")
	if e != nil {
		t.Fatalf("encode returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode:\n%v\nexpect:\n%v", r, expect)
	}

	r, e = enc.encode("12345670")
	if e != nil {
		t.Fatalf("encode returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode:\n%v\nexpect:\n%v", r, expect)
	}
}

func TestEAN8Writer(t *testing.T) {
	writer := NewEAN8Writer()

	if _, e := writer.Encode("96385abc", gozxing.BarcodeFormat_EAN_8, 0, 0, nil); e == nil {
		t.Fatalf("Encode must be error")
	}

	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"    # #  ##  #  #  ## #### # #   ## # # #  ### # #    #   #  ###  # # #     \n"+
		"    # #  ##  #  #  ## #### # #   ## # # #  ### # #    #   #  ###  # # #     \n"+
		"    # #  ##  #  #  ## #### # #   ## # # #  ### # #    #   #  ###  # # #     \n",
		"#", " ")

	matrix, e := writer.Encode("1234567", gozxing.BarcodeFormat_EAN_8, 1, 3, nil)
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
