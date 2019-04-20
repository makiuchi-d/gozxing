package detector

import (
	"math"
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func makePattern(image *gozxing.BitMatrix, x, y, m int) {
	image.SetRegion(x-(3*m), y-(3*m), 7*m, 7*m)
	unsetRegion(image, x-(2*m), y-(2*m), 5*m, 5*m)
	image.SetRegion(x-m, y-m, 3*m, 3*m)
}

func TestNewFinderPatternFinder(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)
	f := NewFinderPatternFinder(image, nil)
	if r := f.GetImage(); r != image {
		t.Fatalf("GetImage ptr=%p, expect %p", r, image)
	}
	if r := f.GetPossibleCenters(); r == nil {
		t.Fatalf("PossibleCenters is nil")
	}
}

func TestFinderPatternFinder_centerFromEnd(t *testing.T) {
	sc := []int{2, 5, 11, 7, 3}
	end := 23
	e := 7.5
	if r := FinderPatternFinder_centerFromEnd(sc, end); r != e {
		t.Fatalf("centerFromEnd = %v, expect %v", r, e)
	}
}

func TestFinderPatternFinder_foundPatternCross(t *testing.T) {
	sc := []int{0, 0, 0, 0, 0}
	if FinderPatternFinder_foundPatternCross(sc) {
		t.Fatalf("foundPatternCross(%v) must be false", sc)
	}

	sc = []int{1, 1, 2, 1, 1}
	if FinderPatternFinder_foundPatternCross(sc) {
		t.Fatalf("foundPatternCross(%v) must be false", sc)
	}

	sc = []int{4, 13, 30, 13, 10} // moduleSize = 10.0, maxVariance = 5
	if FinderPatternFinder_foundPatternCross(sc) {
		t.Fatalf("foundPatternCross(%v) must be false", sc)
	}

	sc = []int{8, 7, 20, 9, 6} // almost good
	if !FinderPatternFinder_foundPatternCross(sc) {
		t.Fatalf("foundPatternCross(%v) must be true", sc)
	}
}

func TestFinderPatternFinder_foundPatternDiagonal(t *testing.T) {
	stateCount := []int{0, 0, 0, 0, 0}
	if FinderPatternFinder_foundPatternDiagonal(stateCount) {
		t.Fatalf("foundPatternDiagonal must return false")
	}

	stateCount = []int{1, 1, 1, 1, 2}
	if FinderPatternFinder_foundPatternDiagonal(stateCount) {
		t.Fatalf("foundPatternDiagonal must return false")
	}

	stateCount = []int{1, 1, 3, 1, 1}
	if !FinderPatternFinder_foundPatternDiagonal(stateCount) {
		t.Fatalf("foundPatternDiagonal must return true")
	}

	stateCount = []int{11, 10, 29, 9, 11}
	if !FinderPatternFinder_foundPatternDiagonal(stateCount) {
		t.Fatalf("foundPatternDiagonal must return true")
	}
}

func TestFinderPatternFinder_ShiftCounts2(t *testing.T) {
	sc := []int{2, 3, 5, 7, 11}
	FinderPatternFinder_ShiftCounts2(sc)
	e := []int{5, 7, 11, 1, 0}

	if !reflect.DeepEqual(sc, e) {
		t.Fatalf("ShiftPattern2 result %v, expect %v", sc, e)
	}
}

func TestFinderPatternFinder_crossCheckDiagonal(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(7, 7)
	f := NewFinderPatternFinder(image, nil)
	if f.crossCheckDiagonal(3, 3) {
		t.Fatalf("crossCheckDiagonal must return false")
	}

	image.Set(3, 3)
	if f.crossCheckDiagonal(3, 3) {
		t.Fatalf("crossCheckDiagonal must return false")
	}

	image.Set(0, 0)
	image.Set(1, 1)
	image.Set(2, 2)
	if f.crossCheckDiagonal(3, 3) {
		t.Fatalf("crossCheckDiagonal must return false")
	}

	image.Unset(1, 1)
	if f.crossCheckDiagonal(3, 3) {
		t.Fatalf("crossCheckDiagonal must return false")
	}

	image.Set(4, 4)
	image.Set(5, 5)
	image.Set(6, 6)
	if f.crossCheckDiagonal(3, 3) {
		t.Fatalf("crossCheckDiagonal must return false")
	}

	image.Unset(5, 5)
	if !f.crossCheckDiagonal(3, 3) {
		t.Fatalf("crossCheckDiagonal must return true")
	}
}

func TestFinderPatternFinder_CrossCheckVertical(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(7, 7)
	f := NewFinderPatternFinder(image, nil)

	for i := 0; i < 7; i++ {
		image.Set(3, i)
	}
	r := f.CrossCheckVertical(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Unset(3, 0)
	r = f.CrossCheckVertical(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Set(3, 0)
	image.Set(3, 1)
	image.Unset(3, 2)
	r = f.CrossCheckVertical(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Set(3, 0)
	image.Unset(3, 1)
	image.Set(3, 2)
	r = f.CrossCheckVertical(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Unset(3, 6)
	r = f.CrossCheckVertical(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Set(3, 6)
	image.Set(3, 5)
	image.Unset(3, 4)
	r = f.CrossCheckVertical(3, 3, 2, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Unset(3, 5)
	image.Set(3, 4)
	r = f.CrossCheckVertical(3, 3, 2, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Unset(3, 6)
	image.Set(3, 5)
	image.Unset(3, 4)
	r = f.CrossCheckVertical(3, 3, 2, 7)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns not NaN, %v", r)
	}

	image.Set(3, 6)
	image.Unset(3, 5)
	image.Set(3, 4)
	r = f.CrossCheckVertical(3, 3, 2, 7)
	if math.IsNaN(r) {
		t.Fatalf("CrossCheckVertical returns NaN")
	}
}

func TestFinderPatternFinder_CrossCheckHorizontal(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(7, 7)
	f := NewFinderPatternFinder(image, nil)

	for i := 0; i < 7; i++ {
		image.Set(i, 3)
	}
	r := f.CrossCheckHorizontal(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Unset(0, 3)
	r = f.CrossCheckHorizontal(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Set(0, 3)
	image.Set(1, 3)
	image.Unset(2, 3)
	r = f.CrossCheckHorizontal(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Set(0, 3)
	image.Unset(1, 3)
	image.Set(2, 3)
	r = f.CrossCheckHorizontal(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Unset(6, 3)
	r = f.CrossCheckHorizontal(3, 3, 1, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Set(6, 3)
	image.Set(5, 3)
	image.Unset(4, 3)
	r = f.CrossCheckHorizontal(3, 3, 2, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Unset(5, 3)
	image.Set(4, 3)
	r = f.CrossCheckHorizontal(3, 3, 2, 1)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Unset(6, 3)
	image.Set(5, 3)
	image.Unset(4, 3)
	r = f.CrossCheckHorizontal(3, 3, 2, 7)
	if !math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns not NaN, %v", r)
	}

	image.Set(6, 3)
	image.Unset(5, 3)
	image.Set(4, 3)
	r = f.CrossCheckHorizontal(3, 3, 2, 7)
	if math.IsNaN(r) {
		t.Fatalf("CrossCheckHorizontal returns NaN")
	}
}

func TestFinderPatternFinder_HandlePossibleCenter(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(20, 20)

	var calledPoint gozxing.ResultPoint
	callback := func(point gozxing.ResultPoint) {
		calledPoint = point
	}
	f := NewFinderPatternFinder(image, callback)

	// vertical ok, horizontal ng
	makePattern(image, 5, 5, 1)
	image.SetRegion(2, 4, 7, 3)
	sc := []int{1, 1, 3, 1, 1}
	i := 5
	j := 9
	if f.HandlePossibleCenter(sc, i, j) {
		t.Fatalf("HandlePossibleCenter(%v,%v,%v) must be false", sc, i, j)
	}

	makePattern(image, 5, 5, 1)
	if !f.HandlePossibleCenter(sc, i, j) {
		t.Fatalf("HandlePossibleCenter(%v,%v,%v) must be true", sc, i, j)
	}
	if calledPoint == nil || calledPoint.GetX() != 5.5 || calledPoint.GetY() != 5.5 {
		t.Fatalf("calledPoint must be {5.5,5.5}, %v", calledPoint)
	}

	calledPoint = nil
	unsetRegion(image, 2, 2, 7, 7)
	makePattern(image, 5, 6, 1) // about equal point
	if !f.HandlePossibleCenter(sc, i, j) {
		t.Fatalf("HandlePossibleCenter(%v,%v,%v) must be true", sc, i, j)
	}
	if calledPoint != nil {
		t.Fatalf("calledPoint must be nil, %v", calledPoint)
	}
}

func TestFinderPatternFinder_HandlePossibleCenterWithPureBarcode(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(20, 20)
	f := NewFinderPatternFinder(image, nil)

	makePattern(image, 5, 5, 1)
	sc := []int{1, 1, 3, 1, 1}
	i := 5
	j := 9
	if !f.HandlePossibleCenterWithPureBarcode(sc, i, j, true) {
		t.Fatalf("HandlePossibleCenter(%v,%v,%v) must be true", sc, i, j)
	}
}

func TestFinderPatternFinder_FindRowSkip(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(20, 20)
	f := NewFinderPatternFinder(image, nil)

	f.possibleCenters = append(f.possibleCenters, NewFinderPattern1(5, 5, 1))
	if r := f.FindRowSkip(); r != 0 {
		t.Fatalf("FindRowSkip returns %v, expect 0", r)
	}

	f.possibleCenters = append(f.possibleCenters, NewFinderPattern(5, 10, 1, 2))
	if r := f.FindRowSkip(); r != 0 {
		t.Fatalf("FindRowSkip returns %v, expect 0", r)
	}
	if f.hasSkipped {
		t.Fatalf("hasSkipped must be false")
	}

	f.possibleCenters = append(f.possibleCenters, NewFinderPattern(10, 10, 1, 2))
	if r := f.FindRowSkip(); r != 2 {
		t.Fatalf("FindRowSkip returns %v, expect 0", 2)
	}
	if !f.hasSkipped {
		t.Fatalf("hasSkipped must be true")
	}
}

func TestNewFinderPatternFinder_HaveMultiplyConfirmedCenters(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(20, 20)
	f := NewFinderPatternFinder(image, nil)
	if f.HaveMultiplyConfirmedCenters() {
		t.Fatalf("HaveMultiplyConfirmedCenters must be false")
	}

	f.possibleCenters = []*FinderPattern{
		NewFinderPattern(10, 10, 2, 2),
		NewFinderPattern(10, 20, 2, 2),
		NewFinderPattern(20, 10, 2, 2),
	}

	if !f.HaveMultiplyConfirmedCenters() {
		t.Fatalf("HaveMultiplyConfirmedCenters must be true")
	}
}

func TestFinderPatternFinder_SelectBestPatterns(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(50, 50)
	f := NewFinderPatternFinder(image, nil)

	_, e := f.SelectBestPatterns()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("SelectBestPatterns must be NotFoundException, %v", e)
	}

	f.possibleCenters = []*FinderPattern{
		NewFinderPattern(20, 10, 4, 2),
		NewFinderPattern(10, 10, 2, 2),
		NewFinderPattern(10, 20, 2, 2),
		NewFinderPattern(20, 20, 4, 2),
	}
	_, e = f.SelectBestPatterns()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("SelectBestPatterns must be NotFoundException, %v", e)
	}

	f.possibleCenters = []*FinderPattern{
		NewFinderPattern(20, 5, 2, 2),
		NewFinderPattern(10, 10, 2, 2),
		NewFinderPattern(10, 20, 2, 2),
		NewFinderPattern(20, 10, 2, 1),
		NewFinderPattern(20, 5, 1, 2),
		NewFinderPattern(30, 30, 2, 2),
	}
	r, e := f.SelectBestPatterns()
	if e != nil {
		t.Fatalf("SelectBestPatterns returns error, %v", e)
	}

	expect1 := NewFinderPattern(20, 10, 2, 1)
	expect2 := NewFinderPattern(10, 10, 2, 2)
	expect3 := NewFinderPattern(10, 20, 2, 2)

	for _, fp := range r {
		if *fp != *expect1 && *fp != *expect2 && *fp != *expect3 {
			t.Fatalf("%v is not contained, [%v, %v, %v]", fp, expect1, expect2, expect3)
		}
	}
}

func TestFinderPatternFinder_SelectBestPatterns2(t *testing.T) {

	tests := []struct {
		input []*FinderPattern
		wants []*FinderPattern
	}{
		{
			input: []*FinderPattern{
				NewFinderPattern(60, 60, 8, 2),
				NewFinderPattern(1420, 60, 8, 2),
				NewFinderPattern(304, 1352, 9.142857, 4),
				NewFinderPattern(60, 1420, 8, 12),
			},
			wants: []*FinderPattern{
				NewFinderPattern(60, 60, 8, 2),
				NewFinderPattern(1420, 60, 8, 2),
				NewFinderPattern(60, 1420, 8, 12),
			},
		},
		{
			input: []*FinderPattern{
				NewFinderPattern(60, 63.5, 4.857142857142857, 2),
				NewFinderPattern(203, 61.5, 3.6428571428571432, 2),
				NewFinderPattern(74, 157, 7, 1),
				NewFinderPattern(54.5, 216, 4.214285714285714, 6),
			},
			wants: []*FinderPattern{
				NewFinderPattern(203, 61.5, 3.6428571428571432, 2),
				NewFinderPattern(54.5, 216, 4.214285714285714, 6),
				NewFinderPattern(60, 63.5, 4.857142857142857, 2),
			},
		},
	}

	for i, test := range tests {
		finder := &FinderPatternFinder{
			possibleCenters: test.input,
		}

		ptns, e := finder.SelectBestPatterns()
		if e != nil {
			t.Fatalf("[%d] SelectBestPatterns returns error: %v", i, e)
		}
		if len(ptns) != 3 {
			t.Fatalf("[%d] result count = %v, wants 3", i, len(ptns))
		}

		for _, p := range test.wants {
			if *ptns[0] != *p && *ptns[1] != *p && *ptns[2] != *p {
				t.Fatalf("[%d] result = [%v, %v, %v], must contains %v",
					i, ptns[0], ptns[1], ptns[2], p)
			}
		}
	}
}

func TestFinderPatternFinder_Find(t *testing.T) {
	s := "                                                          \n" +
		"                                                          \n" +
		"                                                          \n" +
		"        ##############      ##  ##  ##############        \n" +
		"        ##          ##          ##  ##          ##        \n" +
		"        ##  ######  ##  ##  ##      ##  ######  ##        \n" +
		"        ##  ######  ##          ##  ##  ######  ##        \n" +
		"        ##  ######  ##    ##  ####  ##  ######  ##        \n" +
		"        ##          ##    ######    ##          ##        \n" +
		"        ##############  ##  ##  ##  ##############        \n" +
		"                        ##  ##                            \n" +
		"        ######  ##########  ##  ######      ##            \n" +
		"          ##  ##        ########  ##  ##      ####        \n" +
		"        ##    ####  ##  ########  ######  ########        \n" +
		"            ####  ##  ####    ######  ####    ##          \n" +
		"                ##  ##    ##  ##  ######                  \n" +
		"                        ##  ##      ####    ######        \n" +
		"        ##############  ##  ##  ##      ##  ######        \n" +
		"        ##          ##  ######      ######    ####        \n" +
		"        ##  ######  ##  ####    ##  ##        ####        \n" +
		"        ##  ######  ##    ######  ##  ##    ####          \n" +
		"        ##  ######  ##  ########  ####  ##  ##  ##        \n" +
		"        ##          ##  ##  ########    ##    ##          \n" +
		"        ##############  ########  ######      ####        \n" +
		"                                                          \n" +
		"                                                          \n" +
		"                                                          \n"
	image, _ := gozxing.ParseStringToBitMatrix(s, "##", "  ")

	expect := FinderPatternInfo{
		bottomLeft: NewFinderPattern(7.5, 20.5, 1, 2),
		topLeft:    NewFinderPattern(7.5, 6.5, 1, 2),
		topRight:   NewFinderPattern(21.5, 6.5, 1, 2),
	}

	f := NewFinderPatternFinder(image, nil)
	fi, e := f.Find(nil)

	if e != nil {
		t.Fatalf("Find failed: %v", e)
	}

	if *fi.bottomLeft != *expect.bottomLeft {
		t.Fatalf("bottomLeft is %v, expect %v", fi.bottomLeft, expect.bottomLeft)
	}
	if *fi.topLeft != *expect.topLeft {
		t.Fatalf("topLeft is %v, expect %v", fi.topLeft, expect.topLeft)
	}
	if *fi.topRight != *expect.topRight {
		t.Fatalf("topRight is %v, expect %v", fi.topRight, expect.topRight)
	}
}

func TestFinderPatternFinder_Find2(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(40, 40)
	f := NewFinderPatternFinder(image, nil)

	// nested
	image.SetRegion(0, 0, 22, 22)
	unsetRegion(image, 2, 2, 18, 18)
	makePattern(image, 10, 10, 2)

	fi, e := f.Find(nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Find must be NotFoundException, %v", e)
	}

	// touch to edge
	makePattern(image, 10, 32, 2)
	makePattern(image, 32, 32, 2)

	fi, e = f.Find(nil)
	if e != nil {
		t.Fatalf("Find failed: %v", e)
	}

	expect := FinderPatternInfo{
		bottomLeft: NewFinderPattern(33, 33, 2, 2),
		topLeft:    NewFinderPattern(11, 33, 2, 2),
		topRight:   NewFinderPattern(11, 11, 2, 6),
	}

	if *fi.bottomLeft != *expect.bottomLeft {
		t.Fatalf("bottomLeft is %v, expect %v", fi.bottomLeft, expect.bottomLeft)
	}
	if *fi.topLeft != *expect.topLeft {
		t.Fatalf("topLeft is %v, expect %v", fi.topLeft, expect.topLeft)
	}
	if *fi.topRight != *expect.topRight {
		t.Fatalf("topRight is %v, expect %v", fi.topRight, expect.topRight)
	}
}
