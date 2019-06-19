package gozxing

import (
	"reflect"
	"testing"
)

func TestBitArray_GetSize(t *testing.T) {
	b := NewEmptyBitArray()
	if b.GetSize() != 0 || b.GetSizeInBytes() != 0 {
		t.Fatalf("Empty BitArray size is not 0")
	}

	b = NewBitArray(1)
	if b.GetSize() != 1 {
		t.Fatalf("BitArray(1) size is %v, expect 1", b.GetSize())
	}
	if b.GetSizeInBytes() != 1 {
		t.Fatalf("BitArray(1) size in bytes is %v, expect 1", b.GetSizeInBytes())
	}

	b = NewBitArray(8)
	if b.GetSize() != 8 {
		t.Fatalf("BitArray(8) size is %v, expect 8", b.GetSize())
	}
	if b.GetSizeInBytes() != 1 {
		t.Fatalf("BitArray(8) size in bytes is %v, expect 1", b.GetSizeInBytes())
	}

	b = NewBitArray(9)
	if b.GetSize() != 9 {
		t.Fatalf("BitArray(9) size is %v, expect 9", b.GetSize())
	}
	if b.GetSizeInBytes() != 2 {
		t.Fatalf("BitArray(9) size in bytes is %v, expect 2", b.GetSizeInBytes())
	}
}

func TestBitArray_EnsureCapacity(t *testing.T) {
	b := NewBitArray(31)
	b.ensureCapacity(30)
	if b.size != 31 {
		t.Fatalf("Size is changed to %v, expect 31", b.GetSize())
	}

	b.bits[0] = 10
	b.ensureCapacity(33)
	if len(b.bits) != 2 {
		t.Fatalf("bits is not expanded")
	}
	if b.bits[0] != 10 {
		t.Fatalf("bits is not copied")
	}
}

func TestBitArray_GetSetFlip(t *testing.T) {
	b := NewBitArray(5)

	b.Set(3)
	if r := b.Get(3); r != true {
		t.Fatalf("b[3] = %v, expect true", r)
	}

	b.Set(3)
	if r := b.Get(3); r != true {
		t.Fatalf("b[3] = %v, expect true", r)
	}

	b.Flip(3)
	if r := b.Get(3); r != false {
		t.Fatalf("b[3] = %v, expect false", r)
	}

	b.Flip(3)
	if r := b.Get(3); r != true {
		t.Fatalf("b[3] = %v, expect true", r)
	}
}

func TestBitArray_GetNextSet(t *testing.T) {
	b := NewBitArray(65)
	b.Set(10)
	b.Set(33)

	if r := b.GetNextSet(70); r != 65 {
		t.Fatalf("Next set from 70 is %v, expect 65", r)
	}
	if r := b.GetNextSet(3); r != 10 {
		t.Fatalf("Next set from 3 is %v, expect 10", r)
	}
	if r := b.GetNextSet(10); r != 10 {
		t.Fatalf("Next set from 10 is %v, expect 10", r)
	}
	if r := b.GetNextSet(11); r != 33 {
		t.Fatalf("Next set from 11 is %v, expect 33", r)
	}
	if r := b.GetNextSet(34); r != 65 {
		t.Fatalf("Next set from 53 is %v, expect 65", r)
	}

	b.bits[2] = 0x100 // set bit out of range
	if r := b.GetNextSet(34); r != 65 {
		t.Fatalf("Next set from 53 is %v, expect 65", r)
	}
}

func TestBitArray_GetNextUnset(t *testing.T) {
	b := NewBitArray(65)
	b.bits[0] = 0xffffffff
	b.bits[1] = 0xffffffff
	b.bits[2] = 0xffffffff
	b.Flip(10)
	b.Flip(33)

	if r := b.GetNextUnset(70); r != 65 {
		t.Fatalf("Next unset from 70 is %v, expect 65", r)
	}
	if r := b.GetNextUnset(0); r != 10 {
		t.Fatalf("Next unset from 0 is %v, expect 10", r)
	}
	if r := b.GetNextUnset(10); r != 10 {
		t.Fatalf("Next unset from 10 is %v, expect 10", r)
	}
	if r := b.GetNextUnset(11); r != 33 {
		t.Fatalf("Next set from 11 is %v, expect 33", r)
	}
	if r := b.GetNextUnset(34); r != 65 {
		t.Fatalf("Next set from 53 is %v, expect 65", r)
	}

	b.bits[2] = 0x7ff // unset bit out of range
	if r := b.GetNextUnset(34); r != 65 {
		t.Fatalf("Next set from 53 is %v, expect 65", r)
	}
}

func TestBitArray_SetBulk(t *testing.T) {
	b := NewBitArray(64)
	b.SetBulk(0, 0xff00ff00)
	if b.bits[0] != 0xff00ff00 {
		t.Fatalf("bits[0] = %#x, expect 0xff00ff00", b.bits[0])
	}

	b.SetBulk(32, 0x070707)
	if b.bits[1] != 0x070707 {
		t.Fatalf("bits[1] = %#x, expect 0x070707", b.bits[1])
	}
}

func TestBitArray_SetRange(t *testing.T) {
	b := NewBitArray(50)

	if e := b.SetRange(20, 10); e == nil {
		t.Fatalf("SetRange(20, 10) must be error")
	}
	if e := b.SetRange(-1, 10); e == nil {
		t.Fatalf("SetRange(-1, 10) must be error")
	}
	if e := b.SetRange(10, 60); e == nil {
		t.Fatalf("SetRange(10, 60) must be error")
	}
	if e := b.SetRange(10, 10); e != nil {
		t.Fatalf("SetRange(10, 10) must not be error, %v", e)
	}

	e := b.SetRange(10, 40)
	if e != nil {
		t.Fatalf("SetRange(10, 40) must not be error, %v", e)
	}
	if r := b.GetNextSet(0); r != 10 {
		t.Fatalf("Next set is %v, expect 10", r)
	}
	if r := b.GetNextUnset(10); r != 40 {
		t.Fatalf("Next unset is %v, expect 40", r)
	}
}

func TestBitArray_Clear(t *testing.T) {
	b := NewBitArray(64)
	b.SetRange(10, 40)

	b.Clear()
	if b.bits[0] != 0 || b.bits[1] != 0 {
		t.Fatalf("Not cleared, [0]=%x, [1]=%x", b.bits[0], b.bits[1])
	}
}

func testIsRange(t testing.TB, b *BitArray, start, end int, value, expect bool) {
	t.Helper()
	if v, e := b.IsRange(start, end, value); v != expect || e != nil {
		t.Fatalf("IsRange(%v, %v, %v) must %v (%v, %v)", start, end, value, expect, v, e)
	}
}

func TestBitArray_IsRange(t *testing.T) {
	b := NewBitArray(128)

	if _, e := b.IsRange(100, 10, true); e == nil {
		t.Fatalf("IsRange(100, 10) must be error")
	}
	if _, e := b.IsRange(-1, 10, true); e == nil {
		t.Fatalf("IsRange(-1, 10) must be error")
	}
	if _, e := b.IsRange(10, 200, true); e == nil {
		t.Fatalf("IsRange(10, 200) must be error")
	}

	// start==end
	testIsRange(t, b, 10, 10, true, true)

	// all zero
	testIsRange(t, b, 0, 128, false, true)

	b.SetRange(10, 100)

	testIsRange(t, b, 0, 10, true, false)
	testIsRange(t, b, 0, 10, false, true)
	testIsRange(t, b, 0, 11, false, false)
	testIsRange(t, b, 9, 100, true, false)
	testIsRange(t, b, 10, 100, true, true)
	testIsRange(t, b, 10, 100, false, false)
	testIsRange(t, b, 10, 101, true, false)
	testIsRange(t, b, 99, 128, false, false)
	testIsRange(t, b, 100, 128, false, true)
}

func testBit(t testing.TB, b *BitArray, pos int, expect bool) {
	t.Helper()
	if v := b.Get(pos); v != expect {
		t.Fatalf("[%v] = %v, expect %v", pos, v, expect)
	}
}

func TestBitArray_AppendBit(t *testing.T) {
	b := NewBitArray(31)

	b.AppendBit(true)
	if s := b.GetSize(); s != 32 {
		t.Fatalf("size is %v, expect 32", s)
	}
	testBit(t, b, 31, true)

	b.AppendBit(false)
	if s := b.GetSize(); s != 33 {
		t.Fatalf("size is %v, expect 33", s)
	}
	testBit(t, b, 32, false)
}

func TestBitArray_AppendBits(t *testing.T) {
	b := NewBitArray(30)

	if e := b.AppendBits(0, -1); e == nil {
		t.Fatalf("AppendBits(0, -1) must be error, %v", e)
	}
	if e := b.AppendBits(0, 33); e == nil {
		t.Fatalf("AppendBits(0, 33) must be error, %v", e)
	}

	e := b.AppendBits(0x1e, 6)
	if e != nil {
		t.Fatalf("AppendBits failed, %v", e)
	}
	if b.GetSize() != 36 {
		t.Fatalf("size is not expanded, size = %v", b.GetSize())
	}
	testBit(t, b, 30, false)
	testBit(t, b, 31, true)
	testBit(t, b, 32, true)
	testBit(t, b, 33, true)
	testBit(t, b, 34, true)
	testBit(t, b, 35, false)
}

func TestBitArray_AppendBitArray(t *testing.T) {
	b := NewBitArray(30)
	c := NewEmptyBitArray()
	c.AppendBits(0x1e, 6)

	b.AppendBitArray(c)
	if b.GetSize() != 36 {
		t.Fatalf("size is not expanded, size = %v", b.GetSize())
	}
	testBit(t, b, 30, false)
	testBit(t, b, 31, true)
	testBit(t, b, 32, true)
	testBit(t, b, 33, true)
	testBit(t, b, 34, true)
	testBit(t, b, 35, false)
}

func TestBitArray_Xor(t *testing.T) {
	b := NewBitArray(30)
	b.AppendBits(0x1e, 6)

	b1 := NewBitArray(40)
	if b.Xor(b1) == nil {
		t.Fatalf("b xor b1 must be error")
	}

	b2 := NewBitArray(30)
	b2.AppendBits(0x3f, 6)
	if e := b.Xor(b2); e != nil {
		t.Fatalf("b xor b2 must not be error, %v", e)
	}
	testBit(t, b, 29, false)
	testBit(t, b, 30, true)
	testBit(t, b, 31, false)
	testBit(t, b, 32, false)
	testBit(t, b, 33, false)
	testBit(t, b, 34, false)
	testBit(t, b, 35, true)
}

func TestBitArray_ToByte(t *testing.T) {
	// 0f f8 01 f0 00 00
	b := NewBitArray(32 * 6)
	b.SetRange(4, 13)
	b.SetRange(23, 28)

	array := make([]byte, 5)

	b.ToBytes(4, array, 1, 4)
	// [00 ff 80 1f 00]
	if array[0] != 0x00 {
		t.Fatalf("array[0] = %#02x, expect 0x00", array[0])
	}
	if array[1] != 0xff {
		t.Fatalf("array[1] = %#02x, expect 0xff", array[1])
	}
	if array[2] != 0x80 {
		t.Fatalf("array[2] = %#02x, expect 0x80", array[2])
	}
	if array[3] != 0x1f {
		t.Fatalf("array[3] = %#02x, expect 0x1f", array[3])
	}
	if array[4] != 0x00 {
		t.Fatalf("array[4] = %#02x, expect 0x00", array[4])
	}
}

func TestBitArray_GetBitArray(t *testing.T) {
	b := NewEmptyBitArray()
	b.AppendBits(0x1e, 6)

	array := b.GetBitArray()
	if !reflect.DeepEqual(array, b.bits) {
		t.Fatalf("GetBitArray() not equals BitArray.bits")
	}
}

func TestBitArray_Reverse(t *testing.T) {
	b := NewBitArray(80)
	expect := make([]bool, 80)

	for i := 0; i < 80; i++ {
		if i%3 == 0 || i%7 == 0 {
			b.Set(i)
			expect[80-i-1] = true
		}
	}

	b.Reverse()
	for i := 0; i < 80; i++ {
		if b.Get(i) != expect[i] {
			t.Fatalf("[%v] is %v, expect %v", i, b.Get(i), expect[i])
		}
	}
}

func TestBitArray_String(t *testing.T) {
	b := NewBitArray(24)
	expect := " XX....X. XXX..X.. ..XX...."
	b.Set(0)
	b.Set(1)
	b.Set(6)
	b.Set(8)
	b.Set(9)
	b.Set(10)
	b.Set(13)
	b.Set(18)
	b.Set(19)

	if s := b.String(); s != expect {
		t.Fatalf("String is \"%s\", expect \"%s\"", s, expect)
	}
}
