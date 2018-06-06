package detector

import (
	"testing"
)

func TestFinderPatternInfo(t *testing.T) {
	bl := NewFinderPattern(1, 2, 3, 4)
	tl := NewFinderPattern(2, 3, 4, 5)
	tr := NewFinderPattern(3, 4, 5, 6)

	fpi := NewFinderPatternInfo(bl, tl, tr)
	if fpi == nil {
		t.Fatal("NewFinderPatternInfo returns nil")
	}

	if r := fpi.GetBottomLeft(); r != bl {
		t.Fatalf("BottomLeft is not match, %v, %v", r, bl)
	}
	if r := fpi.GetTopLeft(); r != tl {
		t.Fatalf("TopLeft is not match, %v, %v", r, tl)
	}
	if r := fpi.GetTopRight(); r != tr {
		t.Fatalf("TopRight is not match, %v, %v", r, tr)
	}
}
