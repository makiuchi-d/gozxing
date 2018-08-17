package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestUPCAWriter(t *testing.T) {
	writer := NewUPCAWriter()

	_, e := writer.EncodeWithoutHint("123456789012", gozxing.BarcodeFormat_EAN_13, 1, 1)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"    # #  ##  #  #  ## #### # #   ## ##   # # #### # # #   #  #  #   ### #  ###  # ##  ## ## ##  # #     \n"+
		"    # #  ##  #  #  ## #### # #   ## ##   # # #### # # #   #  #  #   ### #  ###  # ##  ## ## ##  # #     \n",
		"#", " ")

	matrix, e := writer.Encode("12345678901", gozxing.BarcodeFormat_UPC_A, 1, 2, nil)
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
