package rss

import (
	"reflect"
	"testing"
)

func TestFinderPattern(t *testing.T) {
	fp := NewFinderPattern(2, []int{178, 247}, 178, 247, 154)

	if r, wants := fp.GetValue(), 2; r != wants {
		t.Fatalf("GetValue() = %v, wants %v", r, wants)
	}
	if r, wants := fp.GetStartEnd(), []int{178, 247}; !reflect.DeepEqual(r, wants) {
		t.Fatalf("GetValue() = %v, wants %v", r, wants)
	}
	rp := fp.GetResultPoints()
	if r, wants := len(rp), 2; r != wants {
		t.Fatalf("GetResultPoints length = %v, wants %v", r, wants)
	}
	if x, y, wx, wy := rp[0].GetX(), rp[0].GetY(), float64(178), float64(154); x != wx || y != wy {
		t.Fatalf("resultPoints[0] = (%v,%v), wants (%v,%v)", x, y, wx, wy)
	}
	if x, y, wx, wy := rp[1].GetX(), rp[1].GetY(), float64(247), float64(154); x != wx || y != wy {
		t.Fatalf("resultPoints[1] = (%v,%v), wants (%v,%v)", x, y, wx, wy)
	}
	if r, wants := fp.HashCode(), 2; r != wants {
		t.Fatalf("HashCode() = %v, wants %v", r, wants)
	}

	tests := []struct {
		rp    interface{}
		wants bool
	}{
		{NewFinderPattern(5, []int{170, 240}, 457, 387, 154), false},
		{NewFinderPattern(2, []int{170, 240}, 457, 387, 154), true},
		{struct{}{}, false},
	}
	for _, test := range tests {
		if r := fp.Equals(test.rp); r != test.wants {
			t.Fatalf("Equals(%v) = %v, wants %v", test.rp, r, test.wants)
		}
	}
}
