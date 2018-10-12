package gozxing

import (
	"testing"
)

func TestDimension(t *testing.T) {
	_, e := NewDimension(-1, 1)
	if e == nil {
		t.Fatalf("NewDimension(0,1) must be error")
	}
	_, e = NewDimension(1, -1)
	if e == nil {
		t.Fatalf("NewDimension(1,0) must be error")
	}
	d, e := NewDimension(10, 20)
	if e != nil {
		t.Fatalf("NewDimension(10, 20) returns error: %v", e)
	}
	if w := d.GetWidth(); w != 10 {
		t.Fatalf("GetWidth = %v, expect 10", w)
	}
	if h := d.GetHeight(); h != 20 {
		t.Fatalf("GetHeight = %v, expect 20", h)
	}
}

func TestDimension_Equals(t *testing.T) {
	d1, _ := NewDimension(10, 20)
	d2, _ := NewDimension(11, 20)
	if d1.Equals(d2) {
		t.Fatalf("%v not equals %v", d1, d2)
	}

	d2, _ = NewDimension(10, 19)
	if d1.Equals(d2) {
		t.Fatalf("%v not equals %v", d1, d2)
	}

	d2, _ = NewDimension(10, 20)
	if !d1.Equals(d2) {
		t.Fatalf("%v equals %v", d1, d2)
	}
}

func TestDimension_HashCode(t *testing.T) {
	d, _ := NewDimension(10, 20)
	if h := d.HashCode(); h != 327150 {
		t.Fatalf("HashCode = %v, expect 327150", h)
	}

	d, _ = NewDimension(20, 15)
	if h := d.HashCode(); h != 654275 {
		t.Fatalf("HashCode = %v, expect 654275", h)
	}
}

func TestDimension_String(t *testing.T) {
	d, _ := NewDimension(10, 20)
	if s := d.String(); s != "10x20" {
		t.Fatalf("String = \"%v\", expect \"10x20\"", s)
	}
}
