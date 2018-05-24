package detector

import (
	"testing"
)

func TestNewAlignmentPattern(t *testing.T) {
	a := NewAlignmentPattern(1.5, 2.5, 3.5)
	if a.GetX() != 1.5 && a.GetY() != 2.5 {
		t.Fatalf("x,y = %v,%v, expect 1.5,2.5", a.GetX(), a.GetY())
	}
	if a.estimatedModuleSize != 3.5 {
		t.Fatalf("estimatedModuleSize = %v, expect 3.5", a.estimatedModuleSize)
	}
}

func TestAlignmentPattern_AboutEquals(t *testing.T) {
	a := NewAlignmentPattern(10.2, 15.3, 1.0)
	if a.AboutEquals(1.0, 15, 12) {
		t.Fatalf("AboutEquals(1.0,15,11) must be false")
	}
	if a.AboutEquals(1.0, 14, 10) {
		t.Fatalf("AboutEquals(1.0,14,10) must be false")
	}
	if a.AboutEquals(0.1, 15, 10) {
		t.Fatalf("AboutEquals(0.1,15,10) must be false")
	}
	if a.AboutEquals(2.1, 15, 10) {
		t.Fatalf("AboutEquals(2.1,15,10) must be false")
	}

	if !a.AboutEquals(1.1, 15, 10) {
		t.Fatalf("AboutEquals(1.1,15,10) must be true")
	}
}

func TestAlignmentPattern_CombineEstimate(t *testing.T) {
	a := NewAlignmentPattern(10.5, 15.5, 1.0)

	a2 := a.CombineEstimate(12.0, 12.0, 2.5)

	if r := a2.GetX(); r != 11.25 {
		t.Fatalf("a2.x = %v, expect 11.25", r)
	}
	if r := a2.GetY(); r != 13.75 {
		t.Fatalf("a2.y = %v, expect 13.75", r)
	}
	if r := a2.estimatedModuleSize; r != 1.75 {
		t.Fatalf("a2.estimatedModuleSize = %v, expect 1.75", r)
	}
}
