package gozxing

import (
	"testing"
)

func TestResultPoint_GetXY(t *testing.T) {
	rp := NewResultPoint(0.5, 2.0)
	if x := rp.GetX(); x != 0.5 {
		t.Fatalf("GetX() = %v, expect %v", x, 0.5)
	}
	if y := rp.GetY(); y != 2.0 {
		t.Fatalf("GetY() = %v, expect %v", y, 2.0)
	}
}

func TestResultPoint_OrderBestPatterns(t *testing.T) {
	var arr []ResultPoint
	a := NewResultPoint(1.5, -1)
	b := NewResultPoint(1, 2)
	c := NewResultPoint(-1, 1)

	arr = []ResultPoint{a, b, c}
	ResultPoint_OrderBestPatterns(arr)
	if arr[0] != a || arr[1] != b || arr[2] != c {
		t.Fatalf("not best pattern %v, expect %v", arr, []ResultPoint{a, b, c})
	}

	arr = []ResultPoint{c, b, a}
	ResultPoint_OrderBestPatterns(arr)
	if arr[0] != a || arr[1] != b || arr[2] != c {
		t.Fatalf("not best pattern %v, expect %v", arr, []ResultPoint{a, b, c})
	}

	arr = []ResultPoint{b, c, a}
	ResultPoint_OrderBestPatterns(arr)
	if arr[0] != a || arr[1] != b || arr[2] != c {
		t.Fatalf("not best pattern %v, expect %v", arr, []ResultPoint{a, b, c})
	}

	arr = []ResultPoint{c, a, b}
	ResultPoint_OrderBestPatterns(arr)
	if arr[0] != a || arr[1] != b || arr[2] != c {
		t.Fatalf("not best pattern %v, expect %v", arr, []ResultPoint{a, b, c})
	}
}
