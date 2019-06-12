package rss

import (
	"testing"
)

func TestPair(t *testing.T) {
	fp := NewFinderPattern(2, []int{178, 247}, 178, 247, 154)
	p := NewPair(60, 5000, fp)

	if r := p.GetFinderPattern(); r != fp {
		t.Fatalf("GetFinderPattern() = %v, wants %v", r, fp)
	}

	if r := p.GetCount(); r != 0 {
		t.Fatalf("GetCount() = %v, wants 0", r)
	}
	p.IncrementCount()
	p.IncrementCount()
	if r := p.GetCount(); r != 2 {
		t.Fatalf("GetCount() = %v, wants 2", r)
	}

}
