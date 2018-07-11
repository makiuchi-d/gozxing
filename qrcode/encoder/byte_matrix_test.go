package encoder

import (
	"reflect"
	"testing"
)

func TestNewByteMatrix(t *testing.T) {
	bm := NewByteMatrix(3, 5)

	if r := bm.GetHeight(); r != 5 {
		t.Fatalf("GetHeight = %v, expect 5", r)
	}
	if r := bm.GetWidth(); r != 3 {
		t.Fatalf("GetWidth = %v, expect 3", r)
	}
	arr := bm.GetArray()
	expect := [][]int8{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	if !reflect.DeepEqual(arr, expect) {
		t.Fatalf("GetArr = %v, expect %v", arr, expect)
	}
}

func TestByteMatrix_SetGet(t *testing.T) {
	bm := NewByteMatrix(3, 5)

	bm.Set(1, 1, 1)
	bm.Set(2, 2, 0)
	bm.SetBool(0, 3, true)
	bm.SetBool(1, 4, false)
	if r := bm.Get(0, 0); r != 0 {
		t.Fatalf("Get(0,0) = %v, expect 0", r)
	}
	if r := bm.Get(1, 1); r != 1 {
		t.Fatalf("Get(1,1) = %v, expect 1", r)
	}
	if r := bm.Get(2, 2); r != 0 {
		t.Fatalf("Get(2,2) = %v, expect 0", r)
	}
	if r := bm.Get(0, 3); r != 1 {
		t.Fatalf("Get(0,3) = %v, expect 1", r)
	}
	if r := bm.Get(1, 4); r != 0 {
		t.Fatalf("Get(1,4) = %v, expect 0", r)
	}

	arr := bm.GetArray()
	expect := [][]int8{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}, {1, 0, 0}, {0, 0, 0}}
	if !reflect.DeepEqual(arr, expect) {
		t.Fatalf("GetArr = %v, expect %v", arr, expect)
	}
}

func TestByteMatrix_Clear(t *testing.T) {
	bm := NewByteMatrix(3, 2)

	bm.Clear(1)
	arr := bm.GetArray()
	expect := [][]int8{{1, 1, 1}, {1, 1, 1}}
	if !reflect.DeepEqual(arr, expect) {
		t.Fatalf("GetArr = %v, expect %v", arr, expect)
	}
}

func TestByteMatrix_String(t *testing.T) {
	bm := NewByteMatrix(3, 3)
	bm.Set(1, 0, 1)
	bm.Set(1, 1, 2)
	bm.Set(2, 1, 1)
	bm.Set(0, 2, 1)
	bm.Set(1, 2, 1)
	bm.Set(2, 2, 1)
	expect := " 0 1 0\n 0   1\n 1 1 1\n"
	if s := bm.String(); s != expect {
		t.Fatalf("String:\n%v\nexpect:\n%v", s, expect)
	}
}
