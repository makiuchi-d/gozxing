package rss

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestAbstractRSSReader(t *testing.T) {
	rdr := NewAbstractRSSReader()

	if l, wants := len(rdr.GetDecodeFinderCounters()), 4; l != wants {
		t.Fatalf("len(decodeFinderCounters) = %v, wants %v", l, wants)
	}
	if l, wants := len(rdr.GetDataCharacterCounters()), 8; l != wants {
		t.Fatalf("len(dataCharacterCounters) = %v, wants %v", l, wants)
	}
	if l, wants := len(rdr.GetOddRoundingErrors()), 4; l != wants {
		t.Fatalf("len(oddRoundingErrors) = %v, wants %v", l, wants)
	}
	if l, wants := len(rdr.GetEvenRoundingErrors()), 4; l != wants {
		t.Fatalf("len(evenRoundingErrors) = %v, wants %v", l, wants)
	}
	if l, wants := len(rdr.GetOddCounts()), 4; l != wants {
		t.Fatalf("len(oddCounts) = %v, wants %v", l, wants)
	}
	if l, wants := len(rdr.GetEvenCounts()), 4; l != wants {
		t.Fatalf("len(evenCounts) = %v, wants %v", l, wants)
	}
}

func TestRSSReader_parseFinderValue(t *testing.T) {
	counters := []int{10, 10, 10, 10}
	_, e := RSSReader_parseFinderValue(counters, rss14_FINDER_PATTERNS)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("parseFinderValue must NotFoundException: %T", e)
	}

	tests := []struct {
		counters []int
		wants    int
	}{
		{[]int{10, 25, 7, 3}, 0},
		{[]int{7, 16, 20, 3}, 5},
	}
	for _, test := range tests {
		r, e := RSSReader_parseFinderValue(test.counters, rss14_FINDER_PATTERNS)
		if e != nil {
			t.Fatalf("parseFinderValue(%v) returns error: %v", test, e)
		}
		if r != test.wants {
			t.Fatalf("parseFinderValue(%v) = %v, wants %v", test, r, test.wants)
		}
	}
}

func TestRSSReader_increment_decrement(t *testing.T) {
	array := []int{1, 2, 3, 4}

	RSSReader_increment(array, []float64{1.5, 1.8, 0.3, 0.9})
	wants := []int{1, 3, 3, 4}
	if !reflect.DeepEqual(array, wants) {
		t.Fatalf("increment result = %v, wants %v", array, wants)
	}

	RSSReader_increment(array, []float64{1.5, 1.5, 1.5, 1.5})
	wants = []int{2, 3, 3, 4}
	if !reflect.DeepEqual(array, wants) {
		t.Fatalf("increment result = %v, wants %v", array, wants)
	}

	RSSReader_decrement(array, []float64{1.5, 1.5, 1.5, 1.5})
	wants = []int{1, 3, 3, 4}
	if !reflect.DeepEqual(array, wants) {
		t.Fatalf("increment result = %v, wants %v", array, wants)
	}

	RSSReader_decrement(array, []float64{1.5, 0.8, 0.2, 1.8})
	wants = []int{1, 3, 2, 4}
	if !reflect.DeepEqual(array, wants) {
		t.Fatalf("increment result = %v, wants %v", array, wants)
	}
}

func TestRSSReader_isFinderPattern(t *testing.T) {

	counters := []int{3, 6, 1, 2} // 9/12
	if RSSReader_isFinderPattern(counters) {
		t.Fatalf("isFinderPattern = true, wants false")
	}

	counters = []int{13, 13, 1, 1} // 13/14
	if RSSReader_isFinderPattern(counters) {
		t.Fatalf("isFinderPattern = true, wants false")
	}

	counters = []int{11, 1, 2, 1}
	if RSSReader_isFinderPattern(counters) {
		t.Fatalf("isFinderPattern = true, wants false")
	}

	counters = []int{8, 2, 1, 1}
	if !RSSReader_isFinderPattern(counters) {
		t.Fatalf("isFinderPattern = false, wants true")
	}

}
