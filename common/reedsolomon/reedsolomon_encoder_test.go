package reedsolomon

import (
	"reflect"
	"testing"
)

func TestReedSolomonEncoder_buildGenerator(t *testing.T) {
	field := GenericGF_QR_CODE_FIELD_256
	enc := NewReedSolomonEncoder(field)

	d := 20
	expect := []int{
		1, 152, 185, 240, 5, 111, 99, 6, 220, 112, 150,
		69, 36, 187, 22, 228, 198, 121, 121, 165, 174,
	}
	r := enc.buildGenerator(d)
	if r.field != field {
		t.Fatalf("buildGenerator(%d) field is %v, expect %v", d, r.field, field)
	}
	if !reflect.DeepEqual(r.coefficients, expect) {
		t.Fatalf("buildGenerator(%d) coefficients = %v, expect %v", 2, r.coefficients, expect)
	}

	d = 2
	expect = []int{1, 3, 2}
	r = enc.buildGenerator(d)
	if r.field != field {
		t.Fatalf("buildGenerator(%d) field is %v, expect %v", d, r.field, field)
	}
	if !reflect.DeepEqual(r.coefficients, expect) {
		t.Fatalf("buildGenerator(%d) coefficients = %v, expect %v", 2, r.coefficients, expect)
	}
}

func TestReedSolomonEncoder_Encode(t *testing.T) {
	enc := NewReedSolomonEncoder(GenericGF_QR_CODE_FIELD_256)
	var e error

	e = enc.Encode([]int{1, 2, 3}, 0)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	e = enc.Encode([]int{1, 2, 3}, 3)
	if e == nil {
		t.Fatalf("Encode must be error")
	}

	toEncode := []int{0, 0, 0, 0, 0, 0}
	expect := []int{0, 0, 0, 0, 0, 0}
	e = enc.Encode(toEncode, 5)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	if !reflect.DeepEqual(toEncode, expect) {
		t.Fatalf("Encode result %v, expect %v", toEncode, expect)
	}

	toEncode = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0}
	expect = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 52, 151, 95, 170, 87}
	e = enc.Encode(toEncode, 5)
	if e != nil {
		t.Fatalf("Encode returns error, %v", e)
	}
	if !reflect.DeepEqual(toEncode, expect) {
		t.Fatalf("Encode result %v, expect %v", toEncode, expect)
	}
}
