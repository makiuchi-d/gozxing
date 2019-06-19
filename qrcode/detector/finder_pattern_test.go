package detector

import (
	"testing"
)

func TestNewFinderPattern1(t *testing.T) {
	f := NewFinderPattern1(10, 20, 5)
	if r := f.GetX(); r != 10 {
		t.Fatalf("GetX() = %v, expect 10", r)
	}
	if r := f.GetY(); r != 20 {
		t.Fatalf("GetY() = %v, expect 20", r)
	}
	if r := f.GetEstimatedModuleSize(); r != 5 {
		t.Fatalf("GetEstimatedModuleSize() = %v, expect 5", r)
	}
	if r := f.GetCount(); r != 1 {
		t.Fatalf("GetCount() = %v, expect 1", r)
	}
}

func testFinderPattern_AboutEquals(t testing.TB, f *FinderPattern, moduleSize, i, j float64, expect bool) {
	t.Helper()
	if r := f.AboutEquals(moduleSize, i, j); r != expect {
		t.Fatalf("AboutEquals(%v,%v,%v) = %v, expect %v", moduleSize, i, j, r, expect)
	}
}

func TestFinderPattern_AboutEquals(t *testing.T) {
	f := NewFinderPattern(10, 20, 4.5, 2)
	testFinderPattern_AboutEquals(t, f, 4, 21, 11, true)
	testFinderPattern_AboutEquals(t, f, 3, 21, 11, true)
	testFinderPattern_AboutEquals(t, f, 10, 21, 11, false)
	testFinderPattern_AboutEquals(t, f, 4, 20, 40, false)
	testFinderPattern_AboutEquals(t, f, 4, 40, 10, false)
}

func TestFinderPattern_CombineEstimate(t *testing.T) {
	f := NewFinderPattern(10, 20, 4, 3)
	f2 := f.CombineEstimate(5, 8, 2)
	expect := NewFinderPattern(9.5, 16.25, 3.5, 4)
	if *f2 != *expect {
		t.Fatalf("CombinedEstimate = %v, expect %v", f2, expect)
	}
}
