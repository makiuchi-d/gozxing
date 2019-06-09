package rss

import (
	"testing"
)

func TestCombines(t *testing.T) {
	tests := []struct {
		n     int
		r     int
		wants int
	}{
		{1, 1, 1},
		{5, 5, 1},
		{5, 2, 10},
		{5, 3, 10},
		{10, 3, 120},
		{10, 4, 210},
		{10, 5, 252},
		{10, 6, 210},
	}
	for _, test := range tests {
		c := combins(test.n, test.r)
		if c != test.wants {
			t.Fatalf("combins(%v, %v) = %v, wants %v", test.n, test.r, c, test.wants)
		}
	}
}

func TestRSSUtils_getRSSvalue(t *testing.T) {
	tests := []struct {
		widths   []int
		maxWidth int
		noNarrow bool
		wants    int
	}{
		{[]int{2, 3, 4, 3}, 8, false, 60},
		{[]int{1, 1, 1, 1}, 1, true, 0},
		{[]int{1, 2, 5, 1}, 6, true, 10},
		{[]int{2, 2, 1, 1}, 3, false, 8},
		{[]int{1, 3, 1, 1}, 3, false, 5},
		{[]int{2, 3, 1, 4}, 6, true, 33},
		{[]int{1, 2, 1, 1}, 2, true, 2},
		{[]int{1, 1, 2, 8}, 8, false, 0},
	}
	for _, test := range tests {
		r := RSSUtils_getRSSvalue(test.widths, test.maxWidth, test.noNarrow)
		if r != test.wants {
			t.Fatalf("getRSSValue(%v) = %v", test, r)
		}
	}
}
