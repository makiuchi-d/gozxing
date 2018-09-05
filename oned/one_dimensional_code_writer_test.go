package oned

import (
	"errors"
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

type dummyEncoder struct{}

func (dummyEncoder) encodeContents(contents string) ([]bool, error) {
	code := make([]bool, 0)
	for _, c := range contents {
		if c != '1' && c != '0' {
			return nil, errors.New("dummy encoder error")
		}
		code = append(code, c != '0')
	}
	return code, nil
}

func (dummyEncoder) getFormat() gozxing.BarcodeFormat {
	return gozxing.BarcodeFormat_UPC_A
}

func TestOnedWriter_renderResult(t *testing.T) {
	_, e := onedWriter_renderResult([]bool{}, 0, 1, -3)
	if e == nil {
		t.Fatalf("renderResult must return error")
	}

	code := []bool{true, false, true, true, false, false, false, true}
	r, e := onedWriter_renderResult(code, 5, 0, 4)
	if e != nil {
		t.Fatalf("renderResult returns error, %v", e)
	}
	if w, h := r.GetWidth(), r.GetHeight(); w != len(code)+4 || h != 1 {
		t.Fatalf("renderResult (%vx%v), expect (%vx%v)", w, h, len(code)+4, 1)
	}
	for i := 0; i < r.GetWidth(); i++ {
		if i < 2 || i >= r.GetWidth()-2 {
			if v := r.Get(i, 0); v != false {
				t.Fatalf("renderResult matrix[%v] = %v, expect %v", i, v, false)
			}
		} else {
			if v := r.Get(i, 0); v != code[i-2] {
				t.Fatalf("renderResult matrix[%v] = %v, expect %v", i, v, code[i-2])
			}
		}
	}
}

func TestOnedWriter_checkNumeric(t *testing.T) {
	e := onedWriter_checkNumeric("1234567890")
	if e != nil {
		t.Fatalf("onedWriter_checkNumeric returns error, %v", e)
	}

	e = onedWriter_checkNumeric("1234a56789")
	if e == nil {
		t.Fatalf("onedWriter_checkNumeric must be error")
	}
}

func TestOnedWriter_appendPattern(t *testing.T) {
	target := make([]bool, 10)

	numAdded := onedWriter_appendPattern(target, 3, []int{1, 2, 3}, true)
	if numAdded != 1+2+3 {
		t.Fatalf("appendPattern numAdded = %v, expect %v", numAdded, 1+2+3)
	}
	expect := []bool{
		false, false, false, true, false, false, true, true, true, false,
	}
	if !reflect.DeepEqual(target, expect) {
		t.Fatalf("appendPattern target = %v, expect %v", target, expect)
	}
}

func TestOnedWriter_encode(t *testing.T) {
	writer := NewOneDimensionalCodeWriter(dummyEncoder{})

	_, e := writer.EncodeWithoutHint("", gozxing.BarcodeFormat_EAN_13, 10, 1)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	_, e = writer.EncodeWithoutHint("", gozxing.BarcodeFormat_UPC_A, 10, 1)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	_, e = writer.EncodeWithoutHint("10110001", gozxing.BarcodeFormat_UPC_A, -1, -1)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	_, e = writer.EncodeWithoutHint("a", gozxing.BarcodeFormat_UPC_A, 1, 1)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	hints := make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_MARGIN] = 1.5
	_, e = writer.Encode("10110001", gozxing.BarcodeFormat_UPC_A, 1, 1, hints)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	hints[gozxing.EncodeHintType_MARGIN] = "abc"
	_, e = writer.Encode("10110001", gozxing.BarcodeFormat_UPC_A, 1, 1, hints)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	hints[gozxing.EncodeHintType_MARGIN] = "0"
	r, e := writer.Encode("10110001", gozxing.BarcodeFormat_UPC_A, 1, 2, hints)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	if w, h := r.GetWidth(), r.GetHeight(); w != 8 || h != 2 {
		t.Fatalf("matrix = (%vx%v), expect (8x2)", w, h)
	}
	for i, b := range "10110001" {
		expect := b != '0'
		if v := r.Get(i, 0); v != expect {
			t.Fatalf("matrix[%v] = %v, expect %v", i, v, expect)
		}
	}

	hints[gozxing.EncodeHintType_MARGIN] = 4
	r, e = writer.Encode("10110001", gozxing.BarcodeFormat_UPC_A, 1, 2, hints)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	if w, h := r.GetWidth(), r.GetHeight(); w != 12 || h != 2 {
		t.Fatalf("matrix = (%vx%v), expect (12x2)", w, h)
	}
	for i, b := range "001011000100" {
		expect := b != '0'
		if v := r.Get(i, 0); v != expect {
			t.Fatalf("matrix[%v] = %v, expect %v", i, v, expect)
		}
	}
}
