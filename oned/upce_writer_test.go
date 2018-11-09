package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestUPCEWriter_encode(t *testing.T) {
	enc := upcEEncoder{}

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

	_, e = enc.encode("2345678")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	// contents=1234567, check=0, parity[1,0]=0x7=LLLGGG
	expect := []bool{
		true, false, true, // start
		false, false, true, false, false, true, true, // 2(L)
		false, true, true, true, true, false, true, // 3(L)
		false, true, false, false, false, true, true, // 4(L)
		false, true, true, true, false, false, true, // 5(G)
		false, false, false, false, true, false, true, // 6(G)
		false, false, true, false, false, false, true, // 7(G)
		false, true, false, true, false, true, // end
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

func TestUPCEWriter(t *testing.T) {
	writer := NewUPCEWriter()

	if _, e := writer.Encode("05096abc", gozxing.BarcodeFormat_UPC_E, 0, 0, nil); e == nil {
		t.Fatalf("Encode must be error")
	}

	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"    # #  #  ## #### # #   ## ###  #    # #  #   # # # #     \n"+
		"    # #  #  ## #### # #   ## ###  #    # #  #   # # # #     \n"+
		"    # #  #  ## #### # #   ## ###  #    # #  #   # # # #     \n",
		"#", " ")

	matrix, e := writer.Encode("1234567", gozxing.BarcodeFormat_UPC_E, 1, 3, nil)
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
