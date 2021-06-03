package gozxing

import (
	"reflect"
	"testing"
)

func TestNewSquareBitMatrix(t *testing.T) {
	b, e := NewSquareBitMatrix(3)

	if e != nil {
		t.Fatalf("NewSquareBitMatrix returns error: %v", e)
	}
	if b.width != 3 || b.height != 3 {
		t.Fatalf("width/height = (%d,%d), expect (3,3)", b.width, b.height)
	}
}

func TestNewBitMatrix(t *testing.T) {
	if _, e := NewBitMatrix(0, 3); e == nil {
		t.Fatalf("NewBitMatrix(0, 3) must be error")
	}
	if _, e := NewBitMatrix(3, -1); e == nil {
		t.Fatalf("NewBitMatrix(3, -1) must be error")
	}

	b, e := NewBitMatrix(33, 3)
	if e != nil {
		t.Fatalf("NewBitMatrix returns error: %v", e)
	}
	if b.width != 33 || b.height != 3 {
		t.Fatalf("width/height = (%d,%d), expect (33,3)", b.width, b.height)
	}
	if b.rowSize < b.width/32 {
		t.Fatalf("b.rowSize too small, rowSize=%v (width=%v)", b.rowSize, b.width)
	}
	if len(b.bits) < b.rowSize*b.height {
		t.Fatalf("b.bits too small, len(b.bits)=%v", len(b.bits))
	}
}

func testBitMatrixGet(t testing.TB, b *BitMatrix, x, y int, expect bool) {
	t.Helper()
	if r := b.Get(x, y); r != expect {
		t.Fatalf("[%v,%v] = %v, expect %v", x, y, r, expect)
	}
}

func TestParseBoolMapToBitMatrix(t *testing.T) {
	if _, e := ParseBoolMapToBitMatrix([][]bool{}); e == nil {
		t.Fatalf("Parse from empty image must be error")
	}

	image := [][]bool{
		{true, true, false},
		{false, true, false},
	}
	b, e := ParseBoolMapToBitMatrix(image)

	if e != nil {
		t.Fatalf("ParseBoolMapToBitMatrix returns error, %v", e)
	}
	if b.width != 3 || b.height != 2 {
		t.Fatalf("BitMatrix size is (%v,%v), expect (3,2)", b.width, b.height)
	}
	testBitMatrixGet(t, b, 2, 0, false)
	testBitMatrixGet(t, b, 1, 1, true)
}

func TestParseStringToBitMatrix(t *testing.T) {
	if _, e := ParseStringToBitMatrix("", "", ""); e == nil {
		t.Fatalf("Parse from empty string must be error")
	}
	if _, e := ParseStringToBitMatrix("\n\n", "X", "."); e == nil {
		t.Fatalf("Parse from no width string must be error")
	}
	if _, e := ParseStringToBitMatrix("XXX\n..\n", "X", "."); e == nil {
		t.Fatalf("Parse from different width string must be error")
	}
	if _, e := ParseStringToBitMatrix("XY\nZZ\n", "X", "."); e == nil {
		t.Fatalf("Parse from string containing unknown chars must be error")
	}
	if _, e := ParseStringToBitMatrix("XX\n.", "X", "."); e == nil {
		t.Fatalf("Parse from incomplete string must be error")
	}

	var b *BitMatrix
	var e error
	b, e = ParseStringToBitMatrix("XX...", "XX", "...")
	if e != nil {
		t.Fatalf("ParseStringToBitMatrix returns error: %v", e)
	}
	if b.width != 2 || b.height != 1 {
		t.Fatalf("BitMatrix size = (%v,%v), expect (2,1)", b.width, b.height)
	}

	b, e = ParseStringToBitMatrix("XX.\n.X.\n", "X", ".")
	if e != nil {
		t.Fatalf("ParseStringToBitMatrix returns error: %v", e)
	}
	if b.width != 3 || b.height != 2 {
		t.Fatalf("BitMatrix size = (%v,%v), expect (3,2)", b.width, b.height)
	}
	if b.Get(1, 1) != true {
		t.Fatalf("[1,1] is not true")
	}
	if b.Get(2, 1) != false {
		t.Fatalf("[0, 2] is not false")
	}
}

func TestBitMatrix_GetSetFlip(t *testing.T) {
	b, _ := NewBitMatrix(7, 7)

	testBitMatrixGet(t, b, 1, 1, false)
	b.Set(1, 1)
	testBitMatrixGet(t, b, 1, 1, true)
	testBitMatrixGet(t, b, 1, 0, false)

	b.Unset(1, 1)
	testBitMatrixGet(t, b, 1, 1, false)

	b.Flip(6, 6)
	testBitMatrixGet(t, b, 6, 6, true)
	b.Flip(6, 6)
	testBitMatrixGet(t, b, 6, 6, false)

	testBitMatrixGet(t, b, -1, 0, false)
	testBitMatrixGet(t, b, 0, -1, false)
	testBitMatrixGet(t, b, 7, 0, false)
	testBitMatrixGet(t, b, 0, 7, false)
}

func TestBitMatrix_FlipAll(t *testing.T) {
	b, _ := ParseStringToBitMatrix("XX.X\n..XX\nX.X.", "X", ".")
	w, _ := ParseStringToBitMatrix("XX.X\n..XX\nX.X.", ".", "X")
	b.FlipAll()

	if b.width != w.width || b.height != w.height {
		t.Fatalf("size mismatch: %vx%v, %vx%v", b.width, b.height, w.width, w.height)
	}
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			bb := b.Get(x, y)
			wb := w.Get(x, y)
			if bb != wb {
				t.Fatalf("[%v,%v] = %v, expect %v", x, y, bb, wb)
			}
		}
	}
}

func TestBitMatrix_Xor(t *testing.T) {
	var b, m *BitMatrix

	b, _ = NewBitMatrix(1, 1)
	m, _ = NewBitMatrix(2, 1)
	if e := b.Xor(m); e == nil {
		t.Fatalf("dimension missmatch xor must be error")
	}
	m, _ = NewBitMatrix(1, 2)
	if e := b.Xor(m); e == nil {
		t.Fatalf("dimension missmatch xor must be error")
	}
	m = &BitMatrix{1, 1, 2, nil}
	if e := b.Xor(m); e == nil {
		t.Fatalf("rowSize missmatch xor must be error")
	}

	// XX.       X.X     .XX
	// ..X  xor  .XX  =  .X.
	b, _ = ParseStringToBitMatrix("XX.\n..X", "X", ".")
	m, _ = ParseStringToBitMatrix("X.X\n.XX", "X", ".")

	e := b.Xor(m)
	if e != nil {
		t.Fatalf("Xor returns error: %v", e)
	}
	testBitMatrixGet(t, b, 0, 0, false)
	testBitMatrixGet(t, b, 1, 0, true)
	testBitMatrixGet(t, b, 2, 0, true)
	testBitMatrixGet(t, b, 0, 1, false)
	testBitMatrixGet(t, b, 1, 1, true)
	testBitMatrixGet(t, b, 2, 1, false)
}

func TestBitMatrix_Clear(t *testing.T) {
	b, _ := ParseStringToBitMatrix("X.XX..\nXX..XX\n..XX.X", "X", ".")

	b.Clear()
	for i := 0; i < len(b.bits); i++ {
		if b.bits[i] != 0 {
			t.Fatalf("bits[%v] is not cleared", i)
		}
	}
}

func TestBitMatrix_SetRegion(t *testing.T) {
	b, _ := NewBitMatrix(40, 40)

	if e := b.SetRegion(-1, 10, 15, 20); e == nil {
		t.Fatalf("SetRegion must be error")
	}
	if e := b.SetRegion(15, -1, 15, 20); e == nil {
		t.Fatalf("SetRegion must be error")
	}
	if e := b.SetRegion(15, 10, 0, 20); e == nil {
		t.Fatalf("SetRegion must be error")
	}
	if e := b.SetRegion(15, 10, 15, 0); e == nil {
		t.Fatalf("SetRegion must be error")
	}
	if e := b.SetRegion(15, 10, 26, 20); e == nil {
		t.Fatalf("SetRegion must be error")
	}
	if e := b.SetRegion(15, 10, 15, 31); e == nil {
		t.Fatalf("SetRegion must be error")
	}

	e := b.SetRegion(15, 10, 15, 20)
	if e != nil {
		t.Fatalf("SetRegion returns error, %v", e)
	}
	testBitMatrixGet(t, b, 14, 10, false)
	testBitMatrixGet(t, b, 15, 9, false)
	testBitMatrixGet(t, b, 15, 10, true)
	testBitMatrixGet(t, b, 14, 29, false)
	testBitMatrixGet(t, b, 15, 30, false)
	testBitMatrixGet(t, b, 15, 29, true)
	testBitMatrixGet(t, b, 29, 9, false)
	testBitMatrixGet(t, b, 30, 10, false)
	testBitMatrixGet(t, b, 29, 10, true)
	testBitMatrixGet(t, b, 29, 30, false)
	testBitMatrixGet(t, b, 30, 29, false)
	testBitMatrixGet(t, b, 29, 29, true)
	testBitMatrixGet(t, b, 20, 29, true)

	e = b.SetRegion(20, 20, 20, 20)
	if e != nil {
		t.Fatalf("SetRegion returns error, %v", e)
	}
	testBitMatrixGet(t, b, 20, 39, true)
	testBitMatrixGet(t, b, 39, 20, true)
	testBitMatrixGet(t, b, 39, 39, true)
}

func TestBitMatrix_GetRow(t *testing.T) {
	b, _ := ParseStringToBitMatrix("X.X.\nX..X\n.XX.", "X", ".")

	ba := b.GetRow(0, nil)

	testBit(t, ba, 0, true)
	testBit(t, ba, 1, false)
	testBit(t, ba, 2, true)
	testBit(t, ba, 3, false)

	ba = b.GetRow(1, NewBitArray(2))
	testBit(t, ba, 0, true)
	testBit(t, ba, 1, false)
	testBit(t, ba, 2, false)
	testBit(t, ba, 3, true)

	ba = b.GetRow(2, ba)
	testBit(t, ba, 0, false)
	testBit(t, ba, 1, true)
	testBit(t, ba, 2, true)
	testBit(t, ba, 3, false)
}

func TestBitMatrix_SetRow(t *testing.T) {
	b, _ := ParseStringToBitMatrix("X.X.XX.\n..XX..X", "X", ".")
	a := NewBitArray(b.width)
	a.SetRange(1, 4) // .XXX...

	b.SetRow(1, a)

	testBitMatrixGet(t, b, 0, 0, true)
	testBitMatrixGet(t, b, 1, 0, false)
	testBitMatrixGet(t, b, 2, 0, true)
	testBitMatrixGet(t, b, 3, 0, false)
	testBitMatrixGet(t, b, 0, 1, false)
	testBitMatrixGet(t, b, 1, 1, true)
	testBitMatrixGet(t, b, 2, 1, true)
	testBitMatrixGet(t, b, 3, 1, true)
	testBitMatrixGet(t, b, 4, 1, false)
	testBitMatrixGet(t, b, 5, 1, false)
	testBitMatrixGet(t, b, 6, 1, false)
}

func TestBitMatrix_Rotate180(t *testing.T) {
	b, _ := ParseStringToBitMatrix("...X.\n.XX..\n..XXX", "X", ".")
	b.Rotate180()

	testBitMatrixGet(t, b, 0, 0, true)
	testBitMatrixGet(t, b, 1, 0, true)
	testBitMatrixGet(t, b, 2, 0, true)
	testBitMatrixGet(t, b, 3, 0, false)
	testBitMatrixGet(t, b, 4, 0, false)

	testBitMatrixGet(t, b, 0, 1, false)
	testBitMatrixGet(t, b, 1, 1, false)
	testBitMatrixGet(t, b, 2, 1, true)
	testBitMatrixGet(t, b, 3, 1, true)
	testBitMatrixGet(t, b, 4, 1, false)

	testBitMatrixGet(t, b, 0, 2, false)
	testBitMatrixGet(t, b, 1, 2, true)
	testBitMatrixGet(t, b, 2, 2, false)
	testBitMatrixGet(t, b, 3, 2, false)
	testBitMatrixGet(t, b, 4, 2, false)
}

func TestBitMatrix_Rotate180WideWidth(t *testing.T) {
	b, _ := NewBitMatrix(130, 3)
	for i := 0; i < b.width; i++ {
		if i%3 == 0 {
			b.Set(i, 0)
		}
		if i%5 == 0 {
			b.Set(i, 1)
		}
		if i%7 == 0 {
			b.Set(i, 2)
		}
	}

	b.Rotate180()

	for i := 0; i < b.width; i++ {
		x := b.width - 1 - i
		if b.Get(x, 2) != (i%3 == 0) {
			t.Fatalf("[%v,0] must be %v", x, (i%3 == 0))
		}
		if b.Get(x, 1) != (i%5 == 0) {
			t.Fatalf("[%v,1] must be %v", x, (i%5 == 0))
		}
		if b.Get(x, 0) != (i%7 == 0) {
			t.Fatalf("[%v,2] must be %v", x, (i%7 == 0))
		}
	}
}

func TestBitMatrix_Rotate90(t *testing.T) {
	b, _ := NewBitMatrix(3, 3)
	b.Set(0, 0)
	b.Set(0, 1)
	b.Set(1, 2)
	b.Set(2, 1)

	b.Rotate90()

	testBitMatrixGet(t, b, 0, 2, true)
	testBitMatrixGet(t, b, 1, 2, true)
	testBitMatrixGet(t, b, 2, 1, true)
	testBitMatrixGet(t, b, 1, 0, true)
}

func TestGetEnclosingRectangle(t *testing.T) {

}

func TestBitMatrix_GetEnclosingRecangle(t *testing.T) {
	b, _ := NewBitMatrix(35, 35)

	n := b.GetEnclosingRectangle()
	if n != nil {
		t.Fatalf("enclosing rectangle must be nil")
	}

	b.Set(5, 8)
	b.Set(3, 10)
	b.Set(33, 20)
	r := b.GetEnclosingRectangle()
	e := []int{3, 8, 31, 13}
	if !reflect.DeepEqual(r, e) {
		t.Fatalf("rectangle is %v, expect %v", r, e)
	}
}

func TestBitMatrix_GetTopLeftOnBit(t *testing.T) {
	b, _ := NewBitMatrix(35, 35)

	if b.GetTopLeftOnBit() != nil {
		t.Fatalf("TopLeftOnBit must be nil")
	}

	b.Set(34, 34)
	r := b.GetTopLeftOnBit()
	e := []int{34, 34}
	if !reflect.DeepEqual(r, e) {
		t.Fatalf("TopLeft is %v, expect %v", r, e)
	}

	b.Set(3, 20)
	b.Set(8, 13)
	r = b.GetTopLeftOnBit()
	e = []int{8, 13}
	if !reflect.DeepEqual(r, e) {
		t.Fatalf("TopLeft is %v, expect %v", r, e)
	}

	b.Set(0, 0)
	r = b.GetTopLeftOnBit()
	e = []int{0, 0}
	if !reflect.DeepEqual(r, e) {
		t.Fatalf("TopLeft is %v, expect %v", r, e)
	}
}

func TestBitMatrix_GetBottomRightOnBit(t *testing.T) {
	b, _ := NewBitMatrix(35, 35)

	if b.GetBottomRightOnBit() != nil {
		t.Fatalf("BottomRightOnBit must be nil")
	}

	b.Set(0, 0)
	r := b.GetBottomRightOnBit()
	e := []int{0, 0}
	if !reflect.DeepEqual(r, e) {
		t.Fatalf("BottomRight is %v, expect %v", r, e)
	}
	b.Set(3, 20)
	b.Set(8, 13)
	r = b.GetBottomRightOnBit()
	e = []int{3, 20}
	if !reflect.DeepEqual(r, e) {
		t.Fatalf("BottomRight is %v, expect %v", r, e)
	}

	b.Set(34, 34)
	r = b.GetBottomRightOnBit()
	e = []int{34, 34}
	if !reflect.DeepEqual(r, e) {
		t.Fatalf("BottomRight is %v, expect %v", r, e)
	}
}

func TestBitMatrix_String(t *testing.T) {
	s := "X   X   \n  X X   \n  X   X \n"
	b, _ := ParseStringToBitMatrix(s, "X ", "  ")

	if r := b.String(); r != s {
		t.Fatalf("String is\n%s\nexpect:\n%s", r, s)
	}

	s2 := "X .X .\n.X X .\n.X .X \n"
	if r := b.ToString("X ", "."); r != s2 {
		t.Fatalf("String is\n%s\nexpect:\n%s", r, s2)
	}
}
