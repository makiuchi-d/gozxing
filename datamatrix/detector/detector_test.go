package detector

import (
	"errors"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
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

func TestMin(t *testing.T) {
	if r := min(-1, 10); r != -1 {
		t.Fatalf("min(-1, 10) = %v, expect -1", r)
	}
	if r := min(20, 10); r != 10 {
		t.Fatalf("min(20, 10) = %v, expect 10", r)
	}
}

func TestMax(t *testing.T) {
	if r := max(-1, 10); r != 10 {
		t.Fatalf("mzx(-1, 10) = %v, expect 10", r)
	}
	if r := max(20, 10); r != 20 {
		t.Fatalf("mzx(20, 10) = %v, expect 20", r)
	}
}

func TestResultPointsAndTransitions(t *testing.T) {
	from := gozxing.NewResultPoint(1, 2)
	to := gozxing.NewResultPoint(5, 6)
	transitions := 1

	rpt := NewResultPointsAndTransitions(from, to, transitions)

	if r := rpt.getFrom(); r != from {
		t.Fatalf("getFrom = %v, expect %v", r, from)
	}

	if r := rpt.getTo(); r != to {
		t.Fatalf("getTo = %v, expect %v", r, to)
	}

	if r := rpt.getTransitions(); r != transitions {
		t.Fatalf("getTransitions = %v, expect %v", r, transitions)
	}

	str := "{1 2}/{5 6}/1"
	if r := rpt.String(); r != str {
		t.Fatalf("String = \"%v\", expect \"%v\"", r, str)
	}

	o1 := NewResultPointsAndTransitions(from, to, 1)
	o2 := NewResultPointsAndTransitions(from, to, 5)

	if !ResultPointsAndTransitionsComparator(o1, o2) {
		t.Fatalf("Comparator(%v,%v) must be false", o1, o2)
	}
	if ResultPointsAndTransitionsComparator(o2, o1) {
		t.Fatalf("Comparator(%v,%v) must be true", o2, o1)
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

	r := det.transitionsBetween(from, to)
	if p := r.getFrom(); p.GetX() != from.GetX() || p.GetY() != from.GetY() {
		t.Fatalf("transitionsBetween from = %v, expect %v", p, from)
	}
	if p := r.getTo(); p.GetX() != to.GetX() || p.GetY() != to.GetY() {
		t.Fatalf("transitionsBetween to = %v, expect %v", p, to)
	}
	if tr := r.getTransitions(); tr != 7 {
		t.Fatalf("transitionsBetween transitions = %v, expect 7", tr)
	}

	to = from
	from = gozxing.NewResultPoint(8, 25)

	img.SetRegion(5, 10, 5, 5)
	img.SetRegion(5, 25, 5, 5)

	r = det.transitionsBetween(from, to)
	if p := r.getFrom(); p.GetX() != from.GetX() || p.GetY() != from.GetY() {
		t.Fatalf("transitionsBetween from = %v, expect %v", p, from)
	}
	if p := r.getTo(); p.GetX() != to.GetX() || p.GetY() != to.GetY() {
		t.Fatalf("transitionsBetween to = %v, expect %v", p, to)
	}
	if tr := r.getTransitions(); tr != 3 {
		t.Fatalf("transitionsBetween transitions = %v, expect 3", tr)
	}
}

func TestCorrectTopRight(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(40, 40)
	det, _ := NewDetector(img)

	// invalid both: c1={41,0}, c2={39,-2}
	topLeft := gozxing.NewResultPoint(9, 0)
	topRight := gozxing.NewResultPoint(39, 0)
	bottomLeft := gozxing.NewResultPoint(9, 30)
	bottomRight := gozxing.NewResultPoint(39, 30)
	r := det.correctTopRight(bottomLeft, bottomRight, topLeft, topRight, 15)
	if r != nil {
		t.Fatalf("correctTopRight must be nil, %v", r)
	}

	// c2 valid: c1={41,5}, c2={39,3}
	topLeft = gozxing.NewResultPoint(9, 5)
	topRight = gozxing.NewResultPoint(39, 5)
	bottomLeft = gozxing.NewResultPoint(9, 35)
	bottomRight = gozxing.NewResultPoint(39, 35)
	r = det.correctTopRight(bottomLeft, bottomRight, topLeft, topRight, 15)
	if r.GetX() != 39 || r.GetY() != 3 {
		t.Fatalf("correctTopRight = %v, expect {39 3}", r)
	}

	// c1 valid: c1={37,1}, c2={35,-1}
	topLeft = gozxing.NewResultPoint(5, 1)
	topRight = gozxing.NewResultPoint(35, 1)
	bottomLeft = gozxing.NewResultPoint(5, 31)
	bottomRight = gozxing.NewResultPoint(35, 31)
	r = det.correctTopRight(bottomLeft, bottomRight, topLeft, topRight, 15)
	if r.GetX() != 37 || r.GetY() != 1 {
		t.Fatalf("correctTopRight = %v, expect {37 1}", r)
	}

	// c1={38,5}, c2={35,2}
	topLeft = gozxing.NewResultPoint(5, 5)
	topRight = gozxing.NewResultPoint(35, 5)
	bottomLeft = gozxing.NewResultPoint(5, 35)
	bottomRight = gozxing.NewResultPoint(35, 35)

	// l1:10-8, l2:2-8; r=c1
	img.SetRegion(8, 5, 3, 3)
	img.SetRegion(14, 5, 3, 3)
	img.SetRegion(20, 5, 3, 3)
	img.SetRegion(26, 5, 3, 3)
	img.SetRegion(33, 5, 3, 3)
	img.SetRegion(35, 11, 3, 3)
	img.SetRegion(35, 17, 3, 3)
	img.SetRegion(35, 23, 3, 3)
	img.SetRegion(35, 29, 3, 3)
	r = det.correctTopRight(bottomLeft, bottomRight, topLeft, topRight, 10)
	if r.GetX() != 38 || r.GetY() != 5 {
		t.Fatalf("correctTopRight = %v, expect {38 5}", r)
	}

	// l1:10-2, l2:8-10; r=c2
	img.Clear()
	img.SetRegion(8, 3, 3, 3)
	img.SetRegion(14, 3, 3, 3)
	img.SetRegion(20, 3, 3, 3)
	img.SetRegion(26, 3, 3, 3)
	img.SetRegion(33, 5, 3, 3)
	img.SetRegion(33, 12, 3, 3)
	img.SetRegion(33, 18, 3, 3)
	img.SetRegion(33, 24, 3, 3)
	img.SetRegion(33, 30, 3, 3)
	r = det.correctTopRight(bottomLeft, bottomRight, topLeft, topRight, 10)
	if r.GetX() != 35 || r.GetY() != 2 {
		t.Fatalf("correctTopRight = %v, expect {35 2}", r)
	}
}

func TestCorrectTopRightRectanglar(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(50, 30)
	det, _ := NewDetector(img)

	// invalid both: c1={50,0}, c2={49,-1}
	topLeft := gozxing.NewResultPoint(9, 0)
	topRight := gozxing.NewResultPoint(49, 0)
	bottomLeft := gozxing.NewResultPoint(9, 20)
	bottomRight := gozxing.NewResultPoint(49, 20)
	r := det.correctTopRightRectangular(bottomLeft, bottomRight, topLeft, topRight, 40, 20)
	if r != nil {
		t.Fatalf("correctTopRight must be nil, %v", r)
	}

	// c2 valid : c1={50,5}, c2={49,4}
	topLeft = gozxing.NewResultPoint(9, 5)
	topRight = gozxing.NewResultPoint(49, 5)
	bottomLeft = gozxing.NewResultPoint(9, 25)
	bottomRight = gozxing.NewResultPoint(49, 25)
	r = det.correctTopRightRectangular(bottomLeft, bottomRight, topLeft, topRight, 40, 20)
	if r.GetX() != 49 || r.GetY() != 4 {
		t.Fatalf("correctTopRight = %v, expect {49 4}", r)
	}

	// c1 valid : c1={46,1}, c2={45,0}
	topLeft = gozxing.NewResultPoint(5, 1)
	topRight = gozxing.NewResultPoint(45, 1)
	bottomLeft = gozxing.NewResultPoint(5, 21)
	bottomRight = gozxing.NewResultPoint(45, 21)
	r = det.correctTopRightRectangular(bottomLeft, bottomRight, topLeft, topRight, 40, 20)
	if r.GetX() != 46 || r.GetY() != 1 {
		t.Fatalf("correctTopRight = %v, expect {46 1}", r)
	}

	// c1={46 5}, c2={45 4}
	topLeft = gozxing.NewResultPoint(5, 5)
	topRight = gozxing.NewResultPoint(45, 5)
	bottomLeft = gozxing.NewResultPoint(5, 25)
	bottomRight = gozxing.NewResultPoint(45, 25)

	// l1=|40-40|+|20-19|, l2=|40-21|+|20-20|; r=c1
	for i := 5; i <= 45; i += 2 {
		img.SetRegion(i, 5, 1, 1)
	}
	for i := 5; i <= 25; i += 2 {
		img.SetRegion(45, i, 2, 1)
	}
	r = det.correctTopRightRectangular(bottomLeft, bottomRight, topLeft, topRight, 40, 20)
	if r.GetX() != 46 || r.GetY() != 5 {
		t.Fatalf("correctTopRight = %v, expect {46 5}", r)
	}

	// l1=|40-40|+|20-11|, l2=|40-39|+|20-20|; r=c2
	img.Clear()
	for i := 5; i <= 45; i += 2 {
		img.SetRegion(i, 4, 1, 2)
	}
	for i := 5; i <= 25; i += 2 {
		img.SetRegion(45, i, 1, 1)
	}
	r = det.correctTopRightRectangular(bottomLeft, bottomRight, topLeft, topRight, 40, 20)
	if r.GetX() != 45 || r.GetY() != 4 {
		t.Fatalf("correctTopRight = %v, expect {45, 4}", r)
	}
}

func checkResultPoint(r gozxing.ResultPoint, x, y float64) bool {
	return r.GetX() == x && r.GetY() == y
}

type dummySampler struct{}

func (dummySampler) SampleGrid(image *gozxing.BitMatrix, dimensionX, dimensionY int,
	p1ToX, p1ToY, p2ToX, p2ToY, p3ToX, p3ToY, p4ToX, p4ToY float64,
	p1FromX, p1FromY, p2FromX, p2FromY, p3FromX, p3FromY, p4FromX, p4FromY float64) (*gozxing.BitMatrix, error) {
	return nil, errors.New("dummy sampler")
}

func (dummySampler) SampleGridWithTransform(image *gozxing.BitMatrix,
	dimensionX, dimensionY int, transform *common.PerspectiveTransform) (*gozxing.BitMatrix, error) {
	return nil, errors.New("dummy sampler")
}

func TestDetector_DetectSquare(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(40, 40)
	det, _ := NewDetector(img)

	if _, e := det.Detect(); e == nil {
		t.Fatalf("Detect must be error")
	}

	// no bottomLeft
	img.SetRegion(15, 15, 20, 5)
	img.Set(25, 30)
	if _, e := det.Detect(); e == nil {
		t.Fatalf("Detect must be error")
	}

	// topRight = A
	img.Clear()
	img.SetRegion(10, 10, 18, 18)
	for i := 0; i < 20; i += 4 {
		img.SetRegion(10+i, 8, 2, 2)
		img.SetRegion(8, 10+i, 2, 2)
	}
	r, e := det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if p := r.GetPoints()[0]; !checkResultPoint(p, 9, 26) {
		t.Fatalf("topLeft = %v, expect {9 26}", p)
	}
	if p := r.GetPoints()[1]; !checkResultPoint(p, 26, 26) {
		t.Fatalf("bottomLeft = %v, expect {26 26}", p)
	}
	if p := r.GetPoints()[2]; !checkResultPoint(p, 26, 9) {
		t.Fatalf("topLeft = %v, expect {26 9}", p)
	}

	// topRight = B
	img.Clear()
	img.SetRegion(10, 10, 18, 18)
	for i := 0; i < 20; i += 4 {
		img.SetRegion(10+i, 28, 2, 2)
		img.SetRegion(8, 10+i, 2, 2)
	}
	r, e = det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if p := r.GetPoints()[0]; !checkResultPoint(p, 26, 28) {
		t.Fatalf("topLeft = %v, expect {26 28}", p)
	}
	if p := r.GetPoints()[1]; !checkResultPoint(p, 26, 11) {
		t.Fatalf("bottomLeft = %v, expect {26 11}", p)
	}
	if p := r.GetPoints()[2]; !checkResultPoint(p, 9, 11) {
		t.Fatalf("topLeft = %v, expect {9 11}", p)
	}

	// topRight = C
	img.Clear()
	img.SetRegion(10, 10, 18, 18)
	for i := 0; i < 20; i += 4 {
		img.SetRegion(10+i, 8, 2, 2)
		img.SetRegion(28, 10+i, 2, 2)
	}
	r, e = det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if p := r.GetPoints()[0]; !checkResultPoint(p, 11, 9) {
		t.Fatalf("topLeft = %v, expect {11 9}", p)
	}
	r, e = det.Detect()
	if p := r.GetPoints()[1]; !checkResultPoint(p, 11, 26) {
		t.Fatalf("bottomLeft = %v, expect {11 26}", p)
	}
	if p := r.GetPoints()[2]; !checkResultPoint(p, 28, 26) {
		t.Fatalf("topLeft = %v, expect {28 26}", p)
	}

	// topRight = D
	img.Clear()
	img.SetRegion(10, 10, 18, 18)
	for i := 0; i < 20; i += 4 {
		img.SetRegion(10+i, 28, 2, 2)
		img.SetRegion(28, 10+i, 2, 2)
	}
	r, e = det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if p := r.GetPoints()[0]; !checkResultPoint(p, 28, 11) {
		t.Fatalf("topLeft = %v, expect {28 11}", p)
	}
	r, e = det.Detect()
	if p := r.GetPoints()[1]; !checkResultPoint(p, 11, 11) {
		t.Fatalf("bottomLeft = %v, expect {11 11}", p)
	}
	if p := r.GetPoints()[2]; !checkResultPoint(p, 11, 28) {
		t.Fatalf("topLeft = %v, expect {11 28}", p)
	}

	img.Clear()
	img.SetRegion(4, 5, 31, 31)
	img.SetRegion(4, 1, 5, 5)
	img.SetRegion(14, 1, 5, 5)
	img.SetRegion(24, 1, 5, 5)
	img.SetRegion(34, 1, 5, 5)
	img.SetRegion(34, 11, 5, 5)
	img.SetRegion(34, 21, 5, 5)
	img.SetRegion(34, 31, 5, 5)
	r, e = det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if p := r.GetPoints()[0]; !checkResultPoint(p, 5, 2) {
		t.Fatalf("topLeft = %v, expect {5 2}", p)
	}
	r, e = det.Detect()
	if p := r.GetPoints()[1]; !checkResultPoint(p, 5, 34) {
		t.Fatalf("bottomLeft = %v, expect {5 34}", p)
	}
	if p := r.GetPoints()[2]; !checkResultPoint(p, 37, 34) {
		t.Fatalf("topLeft = %v, expect {37 34}", p)
	}
	if p := r.GetPoints()[3]; !checkResultPoint(p, 37, 2) {
		t.Fatalf("topLeft = %v, expect {37 2}", p)
	}

	// sample error
	sampler := common.GridSampler_GetInstance()
	common.GridSampler_SetGridSampler(dummySampler{})
	_, e = det.Detect()
	common.GridSampler_SetGridSampler(sampler)
	if e == nil {
		t.Fatalf("Detect must be error")
	}
}

func TestDetector_DetectRectangler(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(70, 40)
	det, _ := NewDetector(img)

	img.SetRegion(4, 5, 61, 31)
	img.SetRegion(4, 1, 5, 5)
	img.SetRegion(14, 1, 5, 5)
	img.SetRegion(24, 1, 5, 5)
	img.SetRegion(34, 1, 5, 5)
	img.SetRegion(44, 1, 5, 5)
	img.SetRegion(54, 1, 5, 5)
	img.SetRegion(64, 1, 5, 5)
	img.SetRegion(64, 11, 5, 5)
	img.SetRegion(64, 21, 5, 5)
	img.SetRegion(64, 31, 5, 5)
	r, e := det.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if p := r.GetPoints()[0]; !checkResultPoint(p, 5, 2) {
		t.Fatalf("topLeft = %v, expect {5 2}", p)
	}
	r, e = det.Detect()
	if p := r.GetPoints()[1]; !checkResultPoint(p, 5, 34) {
		t.Fatalf("bottomLeft = %v, expect {5 34}", p)
	}
	if p := r.GetPoints()[2]; !checkResultPoint(p, 67, 34) {
		t.Fatalf("topLeft = %v, expect {67 34}", p)
	}
	if p := r.GetPoints()[3]; !checkResultPoint(p, 67, 2) {
		t.Fatalf("topLeft = %v, expect {67 2}", p)
	}

	// sample error
	sampler := common.GridSampler_GetInstance()
	common.GridSampler_SetGridSampler(dummySampler{})
	_, e = det.Detect()
	common.GridSampler_SetGridSampler(sampler)
	if e == nil {
		t.Fatalf("Detect must be error")
	}

	img, _ = gozxing.NewBitMatrix(75, 40)
	det, _ = NewDetector(img)
	img.SetRegion(4, 8, 61, 31)
	img.SetRegion(4, 4, 5, 5)
	img.SetRegion(14, 4, 5, 5)
	img.SetRegion(24, 4, 5, 5)
	img.SetRegion(34, 4, 5, 5)
	img.SetRegion(44, 4, 5, 5)
	img.SetRegion(54, 4, 5, 5)
	img.SetRegion(64, 4, 5, 5)
	img.SetRegion(64, 14, 5, 5)
	img.SetRegion(64, 24, 5, 5)
	img.SetRegion(64, 34, 5, 5)
	r, e = det.Detect()

	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if p := r.GetPoints()[0]; !checkResultPoint(p, 5, 5) {
		t.Fatalf("topLeft = %v, expect {5 5}", p)
	}
	r, e = det.Detect()
	if p := r.GetPoints()[1]; !checkResultPoint(p, 5, 37) {
		t.Fatalf("bottomLeft = %v, expect {5 37}", p)
	}
	if p := r.GetPoints()[2]; !checkResultPoint(p, 67, 37) {
		t.Fatalf("topLeft = %v, expect {67 37}", p)
	}
	if p := r.GetPoints()[3]; int(p.GetX()+0.5) != 71 || p.GetY() != 5 {
		t.Fatalf("topLeft = %v, expect {71 5}", p)
	}
}

func TestDetector_Detect(t *testing.T) {
	img, _ := gozxing.ParseStringToBitMatrix(""+
		"                                                \n"+
		"                                                \n"+
		"                                                \n"+
		"                                                \n"+
		"      ##  ##  ##  ##  ##  ##  ##  ##  ##        \n"+
		"      ##  ##  ##  ##  ##  ##  ##  ##  ##        \n"+
		"      ##########      ####  ####        ##      \n"+
		"      ##########      ####  ####        ##      \n"+
		"      ##    ##  ##########    ####    ##        \n"+
		"      ##    ##  ##########    ####    ##        \n"+
		"      ######  ######        ####  ########      \n"+
		"      ######  ######        ####  ########      \n"+
		"      ##    ##  ########  ##  ##  ####          \n"+
		"      ##    ##  ########  ##  ##  ####          \n"+
		"      ##  ######  ##  ##        ##########      \n"+
		"      ##  ######  ##  ##        ##########      \n"+
		"      ######  ##    ##    ##  ####              \n"+
		"      ######  ##    ##    ##  ####              \n"+
		"      ######        ##  ####            ##      \n"+
		"      ######        ##  ####            ##      \n"+
		"      ####  ##  ####  ##    ##  ##              \n"+
		"      ####  ##  ####  ##    ##  ##              \n"+
		"      ####  ##              ##    ########      \n"+
		"      ####  ##              ##    ########      \n"+
		"      ######    ##        ####                  \n"+
		"      ######    ##        ####                  \n"+
		"      ##  ####  ######  ##########      ##      \n"+
		"      ##  ####  ######  ##########      ##      \n"+
		"      ####  ######    ##      ####              \n"+
		"      ####  ######    ##      ####              \n"+
		"      ##  ##########        ##############      \n"+
		"      ##  ##########        ##############      \n"+
		"      ##############  ##          ##  ##        \n"+
		"      ##############  ##          ##  ##        \n"+
		"      ########      ##    ######  ##    ##      \n"+
		"      ########      ##    ######  ##    ##      \n"+
		"      ##  ##  ##        ##  ########            \n"+
		"      ##  ##  ##        ##  ########            \n"+
		"      ####################################      \n"+
		"      ####################################      \n"+
		"                                                \n"+
		"                                                \n"+
		"                                                \n"+
		"                                                \n", "#", " ")
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

	det, _ := NewDetector(img)
	r, e := det.Detect()
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
