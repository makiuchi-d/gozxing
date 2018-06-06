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
	a := NewResultPoint(1.5, -1)
	b := NewResultPoint(1, 2)
	c := NewResultPoint(-1, 1)

	r0, r1, r2 := ResultPoint_OrderBestPatterns(a, b, c)
	if r0 != a || r1 != b || r2 != c {
		t.Fatalf("not best pattern [%v, %v, %v], expect [%v, %v, %v]", r0, r1, r2, a, b, c)
	}

	r0, r1, r2 = ResultPoint_OrderBestPatterns(c, b, a)
	if r0 != a || r1 != b || r2 != c {
		t.Fatalf("not best pattern [%v, %v, %v], expect [%v, %v, %v]", r0, r1, r2, a, b, c)
	}

	r0, r1, r2 = ResultPoint_OrderBestPatterns(b, c, a)
	if r0 != a || r1 != b || r2 != c {
		t.Fatalf("not best pattern [%v, %v, %v], expect [%v, %v, %v]", r0, r1, r2, a, b, c)
	}

	r0, r1, r2 = ResultPoint_OrderBestPatterns(c, a, b)
	if r0 != a || r1 != b || r2 != c {
		t.Fatalf("not best pattern [%v, %v, %v], expect [%v, %v, %v]", r0, r1, r2, a, b, c)
	}
}
