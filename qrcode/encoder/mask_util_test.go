package encoder

import (
	"testing"
)

func makeCheckerMatrix(size int) *ByteMatrix {
	matrix := NewByteMatrix(size, size)
	for j := 0; j < size; j++ {
		for i := j % 2; i < size; i += 2 {
			matrix.SetBool(i, j, true)
		}
	}
	return matrix
}

func TestMaskUtil_applyMaskPenaltyRule1(t *testing.T) {
	matrix := makeCheckerMatrix(15)

	for i := 0; i < 5; i++ {
		matrix.SetBool(i+3, 3, true)
		matrix.SetBool(i+2, 8, true)
		matrix.SetBool(3, i+3, true)
		matrix.SetBool(8, i+1, true)

		matrix.SetBool(i+8, 5, false)
		matrix.SetBool(i+1, 6, false)
		matrix.SetBool(10, i+8, false)
		matrix.SetBool(12, i+6, false)
		matrix.SetBool(13, i+10, false)
	}

	p1 := applyMaskPenaltyRule1Internal(matrix, false)
	if p1 != 23 {
		t.Fatalf("PenaltyRule1(vertical) = %v, expect 23", p1)
	}

	p2 := applyMaskPenaltyRule1Internal(matrix, true)
	if p2 != 14 {
		t.Fatalf("PenaltyRule1(vertical) = %v, expect 23", p2)
	}

	if p := MaskUtil_applyMaskPenaltyRule1(matrix); p != (23 + 14) {
		t.Fatalf("PenaltyRule1 = %v, expect %v", p, 23+14)
	}
}

func TestMaskUtil_applyMaskPenaltyRule2(t *testing.T) {
	matrix := makeCheckerMatrix(15)

	matrix.SetBool(3, 3, false)
	matrix.SetBool(4, 4, false)
	matrix.SetBool(10, 4, false)
	matrix.SetBool(11, 5, false)

	matrix.SetBool(8, 9, true)
	matrix.SetBool(9, 10, true)
	matrix.SetBool(4, 7, true)
	matrix.SetBool(5, 8, true)
	matrix.SetBool(4, 9, true)

	if p := MaskUtil_applyMaskPenaltyRule2(matrix); p != 5*3 {
		t.Fatalf("PenaltyRule2 = %v, expect %v", p, 5*3)
	}
}

func TestMaskUtil_applyMaskPenaltyRule3(t *testing.T) {
	matrix := NewByteMatrix(15, 15)

	matrix.SetBool(4, 0, true)
	matrix.SetBool(6, 0, true)
	matrix.SetBool(7, 0, true)
	matrix.SetBool(8, 0, true)
	matrix.SetBool(10, 0, true)

	matrix.SetBool(8, 14, true)
	matrix.SetBool(10, 14, true)
	matrix.SetBool(11, 14, true)
	matrix.SetBool(12, 14, true)
	matrix.SetBool(14, 14, true)

	matrix.SetBool(4, 2, true)
	matrix.SetBool(4, 3, true)
	matrix.SetBool(4, 4, true)
	matrix.SetBool(4, 6, true)

	matrix.SetBool(14, 8, true)
	matrix.SetBool(14, 10, true)
	matrix.SetBool(14, 11, true)
	matrix.SetBool(14, 12, true)
	matrix.SetBool(14, 14, true)

	if p := MaskUtil_applyMaskPenaltyRule3(matrix); p != 4*40 {
		t.Fatalf("PenaltyRule3 = %v, expect %v", p, 4*40)
	}
}

func TestMaskUtil_isWhiteHorizontal(t *testing.T) {
	arr := []int8{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}

	if r := isWhiteHorizontal(arr, -1, 4); r != true {
		t.Fatalf("isWhiteHorizontal(-1,4) must be %v", !r)
	}

	if r := isWhiteHorizontal(arr, 0, 8); r != false {
		t.Fatalf("isWhiteHorizontal(0,8) must be %v", !r)
	}

	if r := isWhiteHorizontal(arr, 8, 14); r != false {
		t.Fatalf("isWhiteHorizontal(8,14) must be %v", !r)
	}

	if r := isWhiteHorizontal(arr, 15, 20); r != true {
		t.Fatalf("isWhiteHorizontal(15,20) must be %v", !r)
	}
}

func TestMaskUtil_isWhiteVertical(t *testing.T) {
	arr := [][]int8{
		{0, 0, 1},
		{0, 0, 1},
		{0, 0, 0},
		{0, 1, 0},
		{0, 0, 0},
	}

	if r := isWhiteVertical(arr, 0, -1, 6); r != true {
		t.Fatalf("isWhiteVertical(0, -1, 6) must be %v", !r)
	}

	if r := isWhiteVertical(arr, 1, 0, 3); r != true {
		t.Fatalf("isWhiteVertical(0, 0, 3) must be %v", !r)
	}

	if r := isWhiteVertical(arr, 1, 0, 4); r != false {
		t.Fatalf("isWhiteVertical(0, 0, 3) must be %v", !r)
	}

	if r := isWhiteVertical(arr, 2, 2, 5); r != true {
		t.Fatalf("isWhiteVertical(0, 0, 3) must be %v", !r)
	}
}

func TestMaskUtil_applyMaskPenaltyRule4(t *testing.T) {
	matrix := NewByteMatrix(3, 3)

	if p := MaskUtil_applyMaskPenaltyRule4(matrix); p != 100 {
		t.Fatalf("PenaltyRule4 = %v, expect 100", p)
	}

	matrix.SetBool(0, 0, true)
	matrix.SetBool(0, 1, true)
	matrix.SetBool(0, 2, true)
	if p := MaskUtil_applyMaskPenaltyRule4(matrix); p != 30 {
		t.Fatalf("PenaltyRule4 = %v, expect 30", p)
	}

	matrix.SetBool(1, 0, true)
	matrix.SetBool(1, 1, true)
	if p := MaskUtil_applyMaskPenaltyRule4(matrix); p != 10 {
		t.Fatalf("PenaltyRule4 = %v, expect 10", p)
	}
}

func testMaskPattern(t testing.TB, maskPattern int, isMasked func(int, int) bool) {
	t.Helper()
	for x := 0; x < 30; x++ {
		for y := 0; y < 30; y++ {
			bit, e := MaskUtil_getDataMaskBit(maskPattern, x, y)
			if e != nil {
				t.Fatalf("getDataMaskBit(%v,%v,%v) returns error, %v", maskPattern, x, y, e)
			}
			if expect := isMasked(x, y); bit != expect {
				t.Fatalf("getDataMaskBit(%v,%v,%v) = %v, expect %v", maskPattern, x, y, bit, expect)
			}
		}
	}
}

func TestMaskUtil_getDataMaskBit(t *testing.T) {
	testMaskPattern(t, 0, func(j, i int) bool {
		return (i+j)%2 == 0
	})
	testMaskPattern(t, 1, func(j, i int) bool {
		return i%2 == 0
	})
	testMaskPattern(t, 2, func(j, i int) bool {
		return j%3 == 0
	})
	testMaskPattern(t, 3, func(j, i int) bool {
		return (i+j)%3 == 0
	})
	testMaskPattern(t, 4, func(j, i int) bool {
		return (i/2+j/3)%2 == 0
	})
	testMaskPattern(t, 5, func(j, i int) bool {
		return (i*j)%2+(i*j)%3 == 0
	})

	testMaskPattern(t, 6, func(j, i int) bool {
		return ((i*j)%2+(i*j)%3)%2 == 0
	})
	testMaskPattern(t, 7, func(j, i int) bool {
		return ((i+j)%2+(i*j)%3)%2 == 0
	})

	if _, e := MaskUtil_getDataMaskBit(-1, 0, 0); e == nil {
		t.Fatalf("getDataMaskBit(-1,0,0) must be error")
	}
	if _, e := MaskUtil_getDataMaskBit(8, 0, 0); e == nil {
		t.Fatalf("getDataMaskBit(8,0,0) must be error")
	}
}
