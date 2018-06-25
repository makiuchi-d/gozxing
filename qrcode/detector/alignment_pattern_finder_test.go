package detector

import (
	"math"
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func unsetRegion(image *gozxing.BitMatrix, x, y, w, h int) {
	for i := y; i < y+h; i++ {
		for j := x; j < x+w; j++ {
			image.Unset(j, i)
		}
	}
}

func TestAlignmentPatternFinder_centerFromEnd(t *testing.T) {
	stateCount := []int{1, 2, 3}
	end := 10
	expect := float64(6)
	if r := AlignmentPatternFinder_centerFromEnd(stateCount, end); r != expect {
		t.Fatalf("centerFromEnd(%v, %v) = %v, expect %v", stateCount, end, r, expect)
	}
}

func TestAlignmentPatternFinder_foundPatternCross(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)
	f := NewAlignmentPatternFinder(image, 0, 0, 10, 10, 2, nil)

	stateCount := []int{2, 2, 2}
	if !f.foundPatternCross(stateCount) {
		t.Fatalf("foundPatternCross(%v) returns false, expect true", stateCount)
	}

	f.moduleSize = 1
	if f.foundPatternCross(stateCount) {
		t.Fatalf("foundPatternCross(%v) returns true, expect false", stateCount)
	}
}

func TestAlignmentPatternFinder_crossCheckVertical(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(6, 6)
	image.SetRegion(0, 0, 6, 6)

	f := NewAlignmentPatternFinder(image, 0, 0, 5, 5, 1, nil)

	if r := f.crossCheckVertical(2, 2, 1, 3); !math.IsNaN(r) {
		t.Fatalf("crossCheckVertical() must return NaN, %v", r)
	}

	image.Unset(2, 0)
	image.Unset(2, 1)
	if r := f.crossCheckVertical(2, 2, 1, 3); !math.IsNaN(r) {
		t.Fatalf("crossCheckVertical() must return NaN, %v", r)
	}

	image.Set(2, 0)
	if r := f.crossCheckVertical(2, 2, 1, 3); !math.IsNaN(r) {
		t.Fatalf("crossCheckVertical() must return NaN, %v", r)
	}

	image.Unset(2, 3)
	image.Unset(2, 4)
	image.Unset(2, 5)
	if r := f.crossCheckVertical(2, 2, 1, 3); !math.IsNaN(r) {
		t.Fatalf("crossCheckVertical() must return NaN, %v", r)
	}

	image.Set(2, 5)
	if r := f.crossCheckVertical(2, 2, 2, 10); !math.IsNaN(r) {
		t.Fatalf("crossCheckVertical() must return NaN, %v", r)
	}

	if r := f.crossCheckVertical(2, 2, 2, 3); !math.IsNaN(r) {
		t.Fatalf("crossCheckVertical() must return NaN, %v", r)
	}

	image.Set(2, 4)
	if r := f.crossCheckVertical(2, 2, 2, 3); r != 2.5 {
		t.Fatalf("crossCheckVertical() = %v, expect 2.5", r)
	}
}

func TestAlignmentPattern_handlePossibleCenter(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)
	image.SetRegion(0, 0, 10, 10)

	var foundPoint *AlignmentPattern
	rpCallback := func(p gozxing.ResultPoint) {
		foundPoint = p.(*AlignmentPattern)
	}

	f := NewAlignmentPatternFinder(image, 0, 0, 10, 10, 1, rpCallback)

	stateCount := []int{1, 1, 1}

	r := f.handlePossibleCenter(stateCount, 3, 3)
	if r != nil {
		t.Fatalf("handlePossibleCenter() must be nil")
	}
	if foundPoint != nil {
		t.Fatalf("RsultPointCallback called, %v", foundPoint)
	}

	unsetRegion(image, 1, 1, 3, 3)
	image.Set(2, 2)
	foundPoint = nil
	expoint := NewAlignmentPattern(2.5, 2.5, 1)

	r = f.handlePossibleCenter(stateCount, 2, 4)
	if r != nil {
		t.Fatalf("handlePossibleCenter() must be nil")
	}
	if !reflect.DeepEqual(foundPoint, expoint) {
		t.Fatalf("foundPoint = %v, expect %v", foundPoint, expoint)
	}

	foundPoint = nil

	r = f.handlePossibleCenter(stateCount, 2, 4)
	if !reflect.DeepEqual(r, expoint) {
		t.Fatalf("handlePossibleCenter() = %v, expect %v", r, expoint)
	}
	if foundPoint != nil {
		t.Fatalf("RsultPointCallback called, %v", foundPoint)
	}

	unsetRegion(image, 6, 6, 3, 3)
	image.Set(7, 7)
	foundPoint = nil
	expoint = NewAlignmentPattern(7.5, 7.5, 1)

	r = f.handlePossibleCenter(stateCount, 7, 9)
	if r != nil {
		t.Fatalf("handlePossibleCenter() must be nil")
	}
	if !reflect.DeepEqual(foundPoint, expoint) {
		t.Fatalf("foundPoint = %v, expect %v", foundPoint, expoint)
	}
}

func TestAlignmentPatternFinder_Find(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(30, 30)

	f := NewAlignmentPatternFinder(image, 0, 0, 0, 0, 1, nil)

	r, e := f.Find()
	if r != nil || e == nil {
		t.Fatalf("Find must be error, %v, %v", r, e)
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("error must be NotFoundException, %T", e)
	}

	image.SetRegion(1, 1, 5, 5)
	unsetRegion(image, 2, 2, 3, 3)
	image.Set(3, 3)
	expect := NewAlignmentPattern(3.5, 3.5, 1)

	image.SetRegion(9, 9, 12, 12)
	unsetRegion(image, 10, 10, 9, 10)
	image.SetRegion(13, 13, 3, 4)
	expect2 := NewAlignmentPattern(14.5, 15, 3)

	f.width = 10
	f.height = 10
	r, e = f.Find()
	if e != nil {
		t.Fatalf("find failed, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("find returns %v, expect %v", r, expect)
	}

	f.startX = 8
	f.startY = 8
	f.width = 15
	f.height = 15
	f.moduleSize = 3
	r, e = f.Find()
	if e != nil {
		t.Fatalf("find failed, %v", e)
	}
	if !reflect.DeepEqual(r, expect2) {
		t.Fatalf("find returns %v, expect %v", r, expect2)
	}

	// found on edge
	f.startX = 0
	f.startY = 0
	f.width = 5
	f.height = 5
	f.moduleSize = 1
	r, e = f.Find()
	if e != nil {
		t.Fatalf("find failed, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("find returns %v, expect %v", r, expect)
	}
}
