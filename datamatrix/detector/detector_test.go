package detector

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestNewDetector(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(1, 1)

	_, e := NewDetector(img)
	if e == nil {
		t.Fatalf("NewDetector must be error")
	}

	img, _ = gozxing.NewBitMatrix(50, 50)
	d, e := NewDetector(img)
	if e != nil {
		t.Fatalf("NewDetector returns error, %v", e)
	}
	if d == nil {
		t.Fatalf("detector is nil")
	}
	if d.image != img {
		t.Fatalf("idetector.image != img")
	}
	if d.rectangleDetector == nil {
		t.Fatalf("detecotr.rectangleDetector is nil")
	}
}

func TestAbs(t *testing.T) {
	if r := abs(0); r != 0 {
		t.Fatalf("abs(0) = %v, expect 0", r)
	}
	if r := abs(10); r != 10 {
		t.Fatalf("abs(0) = %v, expect 10", r)
	}
	if r := abs(-20); r != 20 {
		t.Fatalf("abs(0) = %v, expect 20", r)
	}
}

func TestMax(t *testing.T) {
	if r := max(-1, 10); r != 10 {
		t.Fatalf("max(-1, 10) = %v, expect 10", r)
	}
	if r := max(20, 10); r != 20 {
		t.Fatalf("max(20, 10) = %v, expect 20", r)
	}
}

func TestMin(t *testing.T) {
	if r := min(-1, 10); r != -1 {
		t.Fatalf("min(-1, 10) = %v, expect -1", r)
	}
	if r := min(20, 10); r != 10 {
		t.Fatalf("min(-1, 10) = %v, expect 10", r)
	}
}

func checkResultPoint(r gozxing.ResultPoint, x, y float64) bool {
	return r != nil && r.GetX() == x && r.GetY() == y
}

func TestShiftPoint(t *testing.T) {
	point := gozxing.NewResultPoint(10, 10)
	to := gozxing.NewResultPoint(30, 50)
	r := shiftPoint(point, to, 4)
	if !checkResultPoint(r, 14, 18) {
		t.Fatalf("shiftPoint = %v, expect {14 18}", r)
	}
}

func TestMoveAway(t *testing.T) {
	point := gozxing.NewResultPoint(10, 10)

	r := moveAway(point, 15, 15)
	if !checkResultPoint(r, 9, 9) {
		t.Fatalf("shiftPoint = %v, expect {9 9}", r)
	}
	r = moveAway(point, 5, 15)
	if !checkResultPoint(r, 11, 9) {
		t.Fatalf("shiftPoint = %v, expect {11 9}", r)
	}
	r = moveAway(point, 15, 5)
	if !checkResultPoint(r, 9, 11) {
		t.Fatalf("shiftPoint = %v, expect {9 11}", r)
	}
	r = moveAway(point, 5, 5)
	if !checkResultPoint(r, 11, 11) {
		t.Fatalf("shiftPoint = %v, expect {11 11}", r)
	}
}

func TestDetectSolid1(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(50, 50)
	det, _ := NewDetector(img)
	points := make([]gozxing.ResultPoint, 4)
	points[0] = gozxing.NewResultPoint(10, 10)
	points[1] = gozxing.NewResultPoint(10, 39)
	points[2] = gozxing.NewResultPoint(39, 10)
	points[3] = gozxing.NewResultPoint(39, 39)

	img.SetRegion(10, 10, 30, 30)

	// solid1 = left
	img.Unset(25, 39)
	img.Unset(39, 25)
	img.Unset(25, 10)
	r := det.detectSolid1(points)
	expect := []gozxing.ResultPoint{
		points[2], points[0], points[1], points[3],
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("detectSolid1 = %v, expect %v", r, expect)
	}

	// solid1 = bottom
	img.Unset(10, 25)
	img.Set(25, 39)
	r = det.detectSolid1(points)
	expect = []gozxing.ResultPoint{
		points[0], points[1], points[3], points[2],
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("detectSolid1 = %v, expect %v", r, expect)
	}

	// solid1 = right
	img.Unset(25, 39)
	img.Set(39, 25)
	r = det.detectSolid1(points)
	expect = []gozxing.ResultPoint{
		points[1], points[3], points[2], points[0],
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("detectSolid1 = %v, expect %v", r, expect)
	}

	// solid1 = top
	img.Unset(39, 25)
	img.Set(25, 10)
	r = det.detectSolid1(points)
	expect = []gozxing.ResultPoint{
		points[3], points[2], points[0], points[1],
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("detectSolid1 = %v, expect %v", r, expect)
	}
}

func TestDetectSolid2(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(50, 50)
	det, _ := NewDetector(img)

	img.SetRegion(10, 10, 5, 30)
	img.SetRegion(20, 10, 5, 5)
	img.SetRegion(30, 10, 5, 5)
	img.SetRegion(35, 15, 5, 5)
	img.SetRegion(35, 25, 5, 5)
	img.SetRegion(10, 35, 30, 5)
	pointA := gozxing.NewResultPoint(10, 10)
	pointB := gozxing.NewResultPoint(10, 39)
	pointC := gozxing.NewResultPoint(39, 39)
	pointD := gozxing.NewResultPoint(34, 10)
	expect := []gozxing.ResultPoint{pointA, pointB, pointC, pointD}

	r := det.detectSolid2([]gozxing.ResultPoint{pointA, pointB, pointC, pointD})
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("detectSolid2 = %v, expect %v", r, expect)
	}

	r = det.detectSolid2([]gozxing.ResultPoint{pointD, pointA, pointB, pointC})
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("detectSolid2 = %v, expect %v", r, expect)
	}
}

func TestCorrectTopRight(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(50, 50)
	det, _ := NewDetector(img)
	points := make([]gozxing.ResultPoint, 4)

	// invalid both
	img.SetRegion(20, 0, 5, 30)
	img.SetRegion(30, 0, 5, 5)
	img.SetRegion(40, 0, 6, 5)
	img.SetRegion(45, 5, 5, 5)
	img.SetRegion(45, 15, 5, 5)
	img.SetRegion(20, 25, 30, 5)
	points[0] = gozxing.NewResultPoint(20, 0)
	points[1] = gozxing.NewResultPoint(20, 29)
	points[2] = gozxing.NewResultPoint(49, 29)
	points[3] = gozxing.NewResultPoint(45, 0)
	r := det.correctTopRight(points)
	if r != nil {
		t.Fatalf("correctTopRight must be nil, %v", r)
	}

	// invalid candidate2
	img.Clear()
	img.SetRegion(18, 3, 5, 30)
	img.SetRegion(28, 3, 5, 5)
	img.SetRegion(38, 3, 5, 5)
	img.SetRegion(43, 8, 5, 5)
	img.SetRegion(43, 18, 5, 5)
	img.SetRegion(18, 28, 30, 5)
	points[0] = gozxing.NewResultPoint(18, 3)
	points[1] = gozxing.NewResultPoint(18, 32)
	points[2] = gozxing.NewResultPoint(47, 32)
	points[3] = gozxing.NewResultPoint(42, 3)
	r = det.correctTopRight(points)
	if !checkResultPoint(r, 47.8, 3) {
		t.Fatalf("correctTopRight = %v, expect {47.8 3}", r)
	}

	// invalide candidate1
	points[3] = gozxing.NewResultPoint(47, 8)
	r = det.correctTopRight(points)
	if !checkResultPoint(r, 47, 2.2) {
		t.Fatalf("correctTopRight = %v, expect {47 2.2}", r)
	}

	// candidate1
	img.Clear()
	img.SetRegion(10, 10, 5, 30)
	img.SetRegion(20, 10, 5, 5)
	img.SetRegion(30, 10, 5, 5)
	img.SetRegion(35, 15, 5, 5)
	img.SetRegion(35, 25, 5, 5)
	img.SetRegion(10, 35, 30, 5)
	points[0] = gozxing.NewResultPoint(10, 10)
	points[1] = gozxing.NewResultPoint(10, 39)
	points[2] = gozxing.NewResultPoint(39, 39)
	points[3] = gozxing.NewResultPoint(34, 10)
	r = det.correctTopRight(points)
	if !checkResultPoint(r, 39.8, 10) {
		t.Fatalf("correctTopRight = %v, expect {39.8 10}", r)
	}

	// candidate2
	points[3] = gozxing.NewResultPoint(39, 15)
	r = det.correctTopRight(points)
	if !checkResultPoint(r, 39, 9.2) {
		t.Fatalf("correctTopRight = %v, expect {39 9.2}", r)
	}
}

func TestShiftToModuleCenter(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(100, 100)
	det, _ := NewDetector(img)
	img.SetRegion(10, 10, 5, 80)
	img.SetRegion(10, 85, 80, 5)
	for i := 1; i < 8; i++ {
		img.SetRegion(10+i*10, 10, 5, 5)
		img.SetRegion(85, 5+i*10, 5, 5)
	}
	img.SetRegion(85, 10, 5, 5) // make dimension odd

	points := []gozxing.ResultPoint{
		gozxing.NewResultPoint(10, 10),
		gozxing.NewResultPoint(10, 89),
		gozxing.NewResultPoint(89, 89),
		gozxing.NewResultPoint(89, 10),
	}

	r := det.shiftToModuleCenter(points)

	if x, y := r[0].GetX(), r[0].GetY(); x <= 10 || x >= 12 || y <= 10 || y >= 12 {
		t.Fatalf("points[0] = %v, expect 10<x<12, 10<y<12", r[0])
	}
	if x, y := r[1].GetX(), r[1].GetY(); x <= 10 || x >= 13 || y <= 87 || y >= 89 {
		t.Fatalf("points[1] = %v, expect 10<x<12, 87<y<89", r[1])
	}
	if x, y := r[2].GetX(), r[2].GetY(); x <= 87 || x >= 89 || y <= 87 || y >= 89 {
		t.Fatalf("points[2] = %v, expect 87<x<89, 87<y<89", r[2])
	}
	if x, y := r[3].GetX(), r[3].GetY(); x <= 87 || x >= 89 || y <= 10 || y >= 13 {
		t.Fatalf("points[3] = %v, expect 10<x<13, 10<y<13", r[3])
	}
}

func TestTransitionsBetween(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(30, 30)
	det, _ := NewDetector(img)

	from := gozxing.NewResultPoint(5, 5)
	to := gozxing.NewResultPoint(25, 8)

	img.SetRegion(8, 5, 3, 4)
	img.SetRegion(13, 5, 3, 4)
	img.SetRegion(18, 5, 3, 4)
	img.SetRegion(23, 5, 3, 4)

	tr := det.transitionsBetween(from, to)
	if tr != 7 {
		t.Fatalf("transitionsBetween transitions = %v, expect 7", tr)
	}

	to = from
	from = gozxing.NewResultPoint(8, 25)

	img.SetRegion(5, 10, 5, 5)
	img.SetRegion(5, 25, 5, 5)

	tr = det.transitionsBetween(from, to)
	if tr != 3 {
		t.Fatalf("transitionsBetween transitions = %v, expect 3", tr)
	}
}

func TestDetector_Detect(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(50, 50)
	det, _ := NewDetector(img)

	_, e := det.Detect()
	if e == nil {
		t.Fatalf("Detect must be error")
	}

	img.SetRegion(18, 2, 5, 30)
	img.SetRegion(28, 2, 5, 5)
	img.SetRegion(38, 2, 10, 5)
	img.SetRegion(43, 7, 5, 5)
	img.SetRegion(43, 17, 5, 5)
	img.SetRegion(18, 27, 30, 5)
	_, e = det.Detect()
	if e == nil {
		t.Fatalf("Detect must be error")
	}

	img, _ = gozxing.NewBitMatrix(100, 100)
	img.SetRegion(10, 10, 5, 80)
	img.SetRegion(10, 85, 80, 5)
	for i := 1; i < 8; i++ {
		img.SetRegion(10+i*10, 10, 5, 5)
		img.SetRegion(85, 5+i*10, 5, 5)
	}
	img.Unset(11, 11) // make dimensionTop odd
	img.Unset(87, 87) // make dimensionRight odd
	det, _ = NewDetector(img)

	r, e := det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	rp := r.GetPoints()
	if x, y := rp[0].GetX(), rp[0].GetY(); x <= 10 || x >= 12 || y <= 10 || y >= 12 {
		t.Fatalf("points[0] = %v, expect 10<x<12, 10<y<12", rp[0])
	}
	if x, y := rp[1].GetX(), rp[1].GetY(); x <= 10 || x >= 13 || y <= 87 || y >= 89 {
		t.Fatalf("points[1] = %v, expect 10<x<12, 87<y<89", rp[1])
	}
	if x, y := rp[2].GetX(), rp[2].GetY(); x <= 87 || x >= 89 || y <= 87 || y >= 89 {
		t.Fatalf("points[2] = %v, expect 87<x<89, 87<y<89", rp[2])
	}
	if x, y := rp[3].GetX(), rp[3].GetY(); x <= 87 || x >= 89 || y <= 10 || y >= 13 {
		t.Fatalf("points[3] = %v, expect 10<x<13, 10<y<13", rp[3])
	}

	sampler := common.GridSampler_GetInstance()
	common.GridSampler_SetGridSampler(testutil.DummyGridSampler{})
	_, e = det.Detect()
	common.GridSampler_SetGridSampler(sampler)
	if e == nil {
		t.Fatalf("Detect must be error")
	}

	img, _ = gozxing.ParseStringToBitMatrix(""+
		"                                                \n"+
		"                                                \n"+
		"      ##  ##  ##  ##  ##  ##  ##  ##  ##        \n"+
		"      ##########      ####  ####        ##      \n"+
		"      ##    ##  ##########    ####    ##        \n"+
		"      ######  ######        ####  ########      \n"+
		"      ##    ##  ########  ##  ##  ####          \n"+
		"      ##  ######  ##  ##        ##########      \n"+
		"      ######  ##    ##    ##  ####              \n"+
		"      ######        ##  ####            ##      \n"+
		"      ####  ##  ####  ##    ##  ##              \n"+
		"      ####  ##              ##    ########      \n"+
		"      ######    ##        ####                  \n"+
		"      ##  ####  ######  ##########      ##      \n"+
		"      ####  ######    ##      ####              \n"+
		"      ##  ##########        ##############      \n"+
		"      ##############  ##          ##  ##        \n"+
		"      ########      ##    ######  ##    ##      \n"+
		"      ##  ##  ##        ##  ########            \n"+
		"      ####################################      \n"+
		"                                                \n"+
		"                                                \n", "##", "  ")
	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"##  ##  ##  ##  ##  ##  ##  ##  ##  \n"+
		"##########      ####  ####        ##\n"+
		"##    ##  ##########    ####    ##  \n"+
		"######  ######        ####  ########\n"+
		"##    ##  ########  ##  ##  ####    \n"+
		"##  ######  ##  ##        ##########\n"+
		"######  ##    ##    ##  ####        \n"+
		"######        ##  ####            ##\n"+
		"####  ##  ####  ##    ##  ##        \n"+
		"####  ##              ##    ########\n"+
		"######    ##        ####            \n"+
		"##  ####  ######  ##########      ##\n"+
		"####  ######    ##      ####        \n"+
		"##  ##########        ##############\n"+
		"##############  ##          ##  ##  \n"+
		"########      ##    ######  ##    ##\n"+
		"##  ##  ##        ##  ########      \n"+
		"####################################\n", "##", "  ")

	det, _ = NewDetector(testutil.ExpandBitMatrix(img, 4))
	r, e = det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	bits := r.GetBits()
	for j := 0; j < expect.GetHeight(); j++ {
		for i := 0; i < expect.GetWidth(); i++ {
			if bits.Get(i, j) != expect.Get(i, j) {
				t.Fatalf("bits(%v,%v) = %v, expect %v", i, j, bits.Get(i, j), expect.Get(i, j))
			}
		}
	}
}
