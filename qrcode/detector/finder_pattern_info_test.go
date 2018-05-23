package detector

import (
	"testing"
)

func TestFinderPatternInfo(t *testing.T) {
	pts := []FinderPattern{
		NewFinderPattern(1, 2, 3, 4),
		NewFinderPattern(2, 3, 4, 5),
		NewFinderPattern(3, 4, 5, 6),
	}
	fpi := NewFinderPatternInfo(pts)
	if fpi == nil {
		t.Fatal("NewFinderPatternInfo returns nil")
	}

	if r := fpi.GetBottomLeft(); r != pts[0] {
		t.Fatalf("BottomLeft is not pts[0], %v, %v", r, pts[0])
	}
	if r := fpi.GetTopLeft(); r != pts[1] {
		t.Fatalf("TopLeft is not pts[1], %v, %v", r, pts[1])
	}
	if r := fpi.GetTopRight(); r != pts[2] {
		t.Fatalf("TopRight is not pts[2], %v, %v", r, pts[2])
	}
}
