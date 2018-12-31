package detector

import (
	"math"
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func makeAlignPattern(image *gozxing.BitMatrix, x, y int) {
	image.SetRegion(x-2, y-2, 5, 5)
	unsetRegion(image, x-1, y-1, 3, 3)
	image.Set(x, y)
}

func makeAlignPattern3(image *gozxing.BitMatrix, x, y int) {
	image.SetRegion(x-7, y-7, 15, 15)
	unsetRegion(image, x-4, y-4, 9, 9)
	image.SetRegion(x-1, y-1, 3, 3)
}
func makePattern3(image *gozxing.BitMatrix, x, y int) {
	image.SetRegion(x-10, y-10, 21, 21)
	unsetRegion(image, x-7, y-7, 15, 15)
	image.SetRegion(x-4, y-4, 9, 9)
}

func TestNewDetector(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)
	d := NewDetector(image)

	if r := d.GetImage(); r != image {
		t.Fatalf("GetImage returns %p, expect %p", r, image)
	}

	if r := d.GetResultPointCallback(); r != nil {
		t.Fatalf("GetResultPointCallback must be nil")
	}
}

func TestDetector_findAlingmentInReagion(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)
	d := NewDetector(image)

	_, e := d.findAlignmentInRegion(5, 10, 10, 3)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("findAlignmentInRegion must be NotFoundException")
	}

	image, _ = gozxing.NewBitMatrix(20, 10)
	d = NewDetector(image)
	_, e = d.findAlignmentInRegion(5, 10, 10, 3)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("findAlignmentInRegion must be NotFoundException")
	}

	image, _ = gozxing.NewBitMatrix(40, 40)
	image.SetRegion(5, 5, 30, 30)

	unsetRegion(image, 7, 7, 6, 6)
	image.SetRegion(9, 9, 2, 2)
	unsetRegion(image, 28, 20, 6, 6)
	image.SetRegion(30, 22, 2, 2)
	unsetRegion(image, 15, 28, 6, 6)
	image.SetRegion(17, 30, 2, 2)

	d = NewDetector(image)
	ap, e := d.findAlignmentInRegion(2, 28, 30, 8)

	if e != nil {
		t.Fatalf("findAlignmentInRegion returns error: %v", e)
	}

	ex := NewAlignmentPattern(31, 23, 2)
	if ap.GetX() != ex.GetX() || ap.GetY() != ex.GetY() || ap.estimatedModuleSize != ex.estimatedModuleSize {
		t.Fatalf("findAlignmentInRegion returns %v, expect %v", ap, ex)
	}
}

func TestDetector_sizeOfBlackWhiteBlackRun(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(20, 20)
	d := NewDetector(image)

	s := d.sizeOfBlackWhiteBlackRun(5, 3, 18, 10)
	if !math.IsNaN(s) {
		t.Fatalf("size must be NaN, %v", s)
	}

	image.SetRegion(9, 0, 4, 20)
	s = d.sizeOfBlackWhiteBlackRun(5, 3, 17, 12)
	if s != 10 {
		t.Fatalf("size must be 10, %v", s)
	}

	image.Clear()
	image.SetRegion(0, 10, 20, 4)
	s = d.sizeOfBlackWhiteBlackRun(12, 17, 3, 5)
	if s != 10 {
		t.Fatalf("size must be 10, %v", s)
	}

	image.Clear()
	image.SetRegion(5, 3, 1, 1)
	image.SetRegion(16, 11, 2, 2)
	s = d.sizeOfBlackWhiteBlackRun(5, 3, 16, 12)
	if s != 15 {
		t.Fatalf("size must be 15, %v", s)
	}
}

func TestDetector_sizeOfBlackWhiteBlackRunBothWays(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(30, 30)
	d := NewDetector(image)

	image.SetRegion(16, 0, 3, 30)
	image.SetRegion(12, 0, 3, 30)
	s := d.sizeOfBlackWhiteBlackRunBothWays(15, 15, 27, 24) // otherTo=(3,6)
	if s != 9 {
		t.Fatalf("size = %v, expect 9", s)
	}

	image.Clear()
	image.SetRegion(16-10, 0, 3, 30)
	image.SetRegion(12-10, 0, 3, 30)
	s = d.sizeOfBlackWhiteBlackRunBothWays(15-10, 15-10, 27-10, 24-10)
	if s != 9 {
		t.Fatalf("size = %v, expect 9", s)
	}

	image.Clear()
	image.SetRegion(16+10, 0, 3, 30)
	image.SetRegion(12+10, 0, 3, 30)
	s = d.sizeOfBlackWhiteBlackRunBothWays(15+10, 15+10, 3+10, 6+10)
	if s != 9 {
		t.Fatalf("size = %v, expect 9", s)
	}

	image.Clear()
	image.SetRegion(0, 16-5, 30, 3)
	image.SetRegion(0, 12-5, 30, 3)
	s = d.sizeOfBlackWhiteBlackRunBothWays(15-5, 15-5, 24-5, 27-5)
	if s != 9 {
		t.Fatalf("size = %v, expect 9", s)
	}

	image.Clear()
	image.SetRegion(0, 16+5, 30, 3)
	image.SetRegion(0, 12+5, 30, 3)
	s = d.sizeOfBlackWhiteBlackRunBothWays(15+5, 15+5, 6+5, 3+5)
	if s != 9 {
		t.Fatalf("size = %v, expect 9", s)
	}
}

func TestDetector_calculateModuleSizeOneWay(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(100, 100)
	makePattern(image, 20, 20, 3)
	makePattern(image, 80, 20, 3)
	d := NewDetector(image)

	p1 := gozxing.NewResultPoint(20, 20)
	p2 := gozxing.NewResultPoint(80, 20)

	s := d.calculateModuleSizeOneWay(p1, p2)
	if s != 3 {
		t.Fatalf("size = %v, expect 3", s)
	}

	image.Clear()
	makePattern(image, 20, 20, 4)
	s = d.calculateModuleSizeOneWay(p1, p2)
	if s != 4 {
		t.Fatalf("size = %v, expect 4", s)
	}

	image.Clear()
	makePattern(image, 20, 80, 2)
	p2 = gozxing.NewResultPoint(20, 80)
	s = d.calculateModuleSizeOneWay(p1, p2)
	if s != 2 {
		t.Fatalf("size = %v, expect 2", s)
	}
}

func TestDetector_calculateModuleSize(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(100, 100)
	makePattern(image, 20, 20, 4)
	makePattern(image, 20, 80, 5)
	makePattern(image, 80, 20, 3)

	d := NewDetector(image)
	topLeft := gozxing.NewResultPoint(20, 20)
	topRight := gozxing.NewResultPoint(80, 20)
	bottomLeft := gozxing.NewResultPoint(20, 80)

	s := d.calculateModuleSize(topLeft, topRight, bottomLeft)
	if s != 4 {
		t.Fatalf("size = %v, expect 4", s)
	}
}

func TestDetector_computeDimension(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(100, 100)
	d := NewDetector(image)
	topLeft := gozxing.NewResultPoint(20, 20)
	topRight := gozxing.NewResultPoint(80, 20)
	bottomLeft := gozxing.NewResultPoint(20, 80) // dimension %4 = 1

	r, e := d.computeDimension(topLeft, topRight, bottomLeft, 2)
	if e != nil {
		t.Fatalf("computeDimension returns error, %v", e)
	}
	if r != 37 {
		t.Fatalf("computeDimension returns %v, expect 37", r)
	}

	bottomLeft = gozxing.NewResultPoint(20, 75) // dimension %4 = 0
	r, e = d.computeDimension(topLeft, topRight, bottomLeft, 2)
	if e != nil {
		t.Fatalf("computeDimension returns error, %v", e)
	}
	if r != 37 {
		t.Fatalf("computeDimension returns %v, expect 37", r)
	}

	bottomLeft = gozxing.NewResultPoint(20, 85) // dimension %4 = 2
	r, e = d.computeDimension(topLeft, topRight, bottomLeft, 2)
	if e != nil {
		t.Fatalf("computeDimension returns error, %v", e)
	}
	if r != 37 {
		t.Fatalf("computeDimension returns %v, expect 37", r)
	}

	bottomLeft = gozxing.NewResultPoint(20, 72) // dimension %4 = 3
	_, e = d.computeDimension(topLeft, topRight, bottomLeft, 2)
	if e == nil {
		t.Fatalf("computeDimension must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("computeDimension must be NotFoundException, %v", e)
	}
}

func TestDetector_createTransform(t *testing.T) {

	topLeft := gozxing.NewResultPoint(20, 20)
	topRight := gozxing.NewResultPoint(20, 80)
	bottomLeft := gozxing.NewResultPoint(80, 20)
	alignmentPattern := NewAlignmentPattern(75, 75, 1)

	transform := Detector_createTransform(topLeft, topRight, bottomLeft, alignmentPattern, 1)
	v := reflect.ValueOf(transform).Elem()
	if r := v.FieldByName("a11").Float(); r != -31.5 {
		t.Fatalf("transform.a11 = %v, expect -31.5", r)
	}
	if r := v.FieldByName("a21").Float(); r != -328.5 {
		t.Fatalf("transform.a21 = %v, expect -328.5", r)
	}
	if r := v.FieldByName("a31").Float(); r != 1665 {
		t.Fatalf("transform.a31 = %v, expect 1665", r)
	}
	if r := v.FieldByName("a12").Float(); r != -328.5 {
		t.Fatalf("transform.a12 = %v, expect -328.5", r)
	}
	if r := v.FieldByName("a22").Float(); r != -31.5 {
		t.Fatalf("transform.a22 = %v, expect -31.5", r)
	}
	if r := v.FieldByName("a32").Float(); r != 1665 {
		t.Fatalf("transform.a32 = %v, expect 1665", r)
	}
	if r := v.FieldByName("a13").Float(); r != -1.575 {
		t.Fatalf("transform.a13 = %v, expect -1.575", r)
	}
	if r := v.FieldByName("a23").Float(); r != -1.575 {
		t.Fatalf("transform.a23 = %v, expect -1.575", r)
	}
	if r := v.FieldByName("a33").Float(); r != 31.275 {
		t.Fatalf("transform.a33 = %v, expect 31.275", r)
	}

	topLeft = gozxing.NewResultPoint(20, 50)
	topRight = gozxing.NewResultPoint(50, 20)
	bottomLeft = gozxing.NewResultPoint(50, 80)
	transform = Detector_createTransform(topLeft, topRight, bottomLeft, nil, 2)

	v = reflect.ValueOf(transform).Elem()
	if r := v.FieldByName("a11").Float(); r != -150 {
		t.Fatalf("transform.a11 = %v, expect -150", r)
	}
	if r := v.FieldByName("a21").Float(); r != -150 {
		t.Fatalf("transform.a21 = %v, expect -150", r)
	}
	if r := v.FieldByName("a31").Float(); r != 1550 {
		t.Fatalf("transform.a31 = %v, expect 1550", r)
	}
	if r := v.FieldByName("a12").Float(); r != 150 {
		t.Fatalf("transform.a12 = %v, expect 150", r)
	}
	if r := v.FieldByName("a22").Float(); r != -150 {
		t.Fatalf("transform.a22 = %v, expect -150", r)
	}
	if r := v.FieldByName("a32").Float(); r != 1250 {
		t.Fatalf("transform.a32 = %v, expect 1250", r)
	}
	if r := v.FieldByName("a13").Float(); r != 0 {
		t.Fatalf("transform.a13 = %v, expect 0", r)
	}
	if r := v.FieldByName("a23").Float(); r != 0 {
		t.Fatalf("transform.a23 = %v, expect 0", r)
	}
	if r := v.FieldByName("a33").Float(); r != 25 {
		t.Fatalf("transform.a33 = %v, expect 25", r)
	}
}

func TestDetector_ProcessFinderPatternInfo(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(300, 300)
	d := NewDetector(image)

	// too small patterns
	image.SetRegion(18, 18, 5, 5)
	unsetRegion(image, 19, 19, 3, 3)
	image.Set(20, 20)
	image.SetRegion(28, 18, 5, 5)
	unsetRegion(image, 29, 19, 3, 3)
	image.Set(30, 20)
	image.SetRegion(18, 28, 5, 5)
	unsetRegion(image, 19, 29, 3, 3)
	image.Set(20, 30)
	info := FinderPatternInfo{
		topLeft:    NewFinderPattern1(20, 20, 1),
		topRight:   NewFinderPattern1(30, 20, 1),
		bottomLeft: NewFinderPattern1(20, 30, 1),
	}
	_, e := d.ProcessFinderPatternInfo(&info)
	if e == nil {
		t.Fatalf("ProcessFinderPatternInfo must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("ProcessFinderPatternInfo must be NotFoundException, %v", e)
	}
	// dimension error
	image.Clear()
	makePattern(image, 20, 20, 1)
	makePattern(image, 20, 80, 1)
	makePattern(image, 72, 20, 1)
	info.topRight = NewFinderPattern1(20, 80, 1)
	info.bottomLeft = NewFinderPattern1(72, 20, 1)
	_, e = d.ProcessFinderPatternInfo(&info)
	if e == nil {
		t.Fatalf("ProcessFinderPatternInfo must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("ProcessFinderPatternInfo must be NotFoundException, %v", e)
	}

	// version error
	image.Clear()
	makePattern(image, 20, 20, 1)
	makePattern(image, 20, 270, 1)
	makePattern(image, 270, 20, 1)
	info.topRight = NewFinderPattern1(20, 270, 1)
	info.bottomLeft = NewFinderPattern1(270, 20, 1)
	_, e = d.ProcessFinderPatternInfo(&info)
	if e == nil {
		t.Fatalf("ProcessFinderPatternInfo must be error")
	}
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("ProcessFinderPatternInfo must be FormatException, %v", e)
	}

	// no alignment patterns
	image.Clear()
	makePattern(image, 20, 20, 1)
	makePattern(image, 35, 20, 1)
	makePattern(image, 20, 35, 1)
	info.topLeft = NewFinderPattern1(20, 20, 1)
	info.topRight = NewFinderPattern1(35, 20, 1)
	info.bottomLeft = NewFinderPattern1(20, 35, 1)
	r, e := d.ProcessFinderPatternInfo(&info)
	if e != nil {
		t.Fatalf("ProcessFinderPatternInfo returns error, %v", e)
	}
	if c := len(r.GetPoints()); c != 3 {
		t.Fatalf("points has %v members, expect 3", c)
	}
	if p := r.GetPoints()[0]; p != info.bottomLeft {
		t.Fatalf("points[0] is not bottomLeft (%p, %p)", p, info.bottomLeft)
	}
	if p := r.GetPoints()[1]; p != info.topLeft {
		t.Fatalf("points[0] is not topLeft (%p, %p)", p, info.topLeft)
	}
	if p := r.GetPoints()[2]; p != info.topRight {
		t.Fatalf("points[0] is not topRight (%p, %p)", p, info.topRight)
	}

	// with alignment patterns
	image.Clear()
	// version 7 pattern
	makePattern(image, 13, 13, 1)
	makePattern(image, 13+38, 13, 1)
	makePattern(image, 13, 13+38, 1)
	info.topLeft = NewFinderPattern1(13, 13, 1)
	info.topRight = NewFinderPattern1(13+38, 13, 1)
	info.bottomLeft = NewFinderPattern1(13, 13+38, 1)
	makeAlignPattern(image, 10+6, 10+22)
	makeAlignPattern(image, 10+22, 10+6)
	makeAlignPattern(image, 10+22, 10+22)
	makeAlignPattern(image, 10+22, 10+38)
	makeAlignPattern(image, 10+38, 10+22)
	makeAlignPattern(image, 10+38, 10+38)
	r, e = d.ProcessFinderPatternInfo(&info)
	if e != nil {
		t.Fatalf("ProcessFinderPatternInfo returns error, %v", e)
	}
	if c := len(r.GetPoints()); c != 4 {
		t.Fatalf("points has %v members, expect 4", c)
	}
	if p := r.GetPoints()[0]; p != info.bottomLeft {
		t.Fatalf("points[0] is not bottomLeft (%p, %p)", p, info.bottomLeft)
	}
	if p := r.GetPoints()[1]; p != info.topLeft {
		t.Fatalf("points[1] is not topLeft (%p, %p)", p, info.topLeft)
	}
	if p := r.GetPoints()[2]; p != info.topRight {
		t.Fatalf("points[2] is not topRight (%p, %p)", p, info.topRight)
	}
	if p := r.GetPoints()[3]; p.GetX() != 48.5 || p.GetY() != 48.5 {
		t.Fatalf("alignmentPattern is (%v,%v), expect (48.5,48.5)", p.GetX(), p.GetY())
	}

	// grid sampler fail
	image.Clear()
	makePattern(image, 290, 20, 1)
	makePattern(image, 250, 60, 1)
	makePattern(image, 290, 100, 1)
	info.topLeft = NewFinderPattern1(250, 60, 1)
	info.topRight = NewFinderPattern1(290, 20, 1)
	info.bottomLeft = NewFinderPattern1(290, 100, 1)
	r, e = d.ProcessFinderPatternInfo(&info)
	if e == nil {
		t.Fatalf("ProcessFinderPatternInfo must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("ProcessFinderPatternInfo must be NotFoundException, %v", e)
	}
}

func TestDetector_DetectWithHint(t *testing.T) {
	cbpoints := make([]gozxing.ResultPoint, 0)
	callback := gozxing.ResultPointCallback(func(p gozxing.ResultPoint) {
		cbpoints = append(cbpoints, p)
	})

	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK] = callback

	// version3 pattern
	image, _ := gozxing.NewBitMatrix(45, 45)
	makePattern(image, 10+3, 10+3, 1)
	makePattern(image, 10+3+22, 10+3, 1)
	makePattern(image, 10+3, 10+3+22, 1)
	makeAlignPattern(image, 10+22, 10+22)

	d := NewDetector(image)
	r, e := d.Detect(hints)
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}
	if len(cbpoints) != 4 {
		t.Fatalf("Callback count is %v, expect 4", len(cbpoints))
	}

	// bottomLeft
	if p := r.GetPoints()[0]; p.GetX() != 13.5 || p.GetY() != 35.5 {
		t.Fatalf("points[0] = (%v,%v), expect (13.5,35.5)", p.GetX(), p.GetY())
	}
	// topLeft
	if p := r.GetPoints()[1]; p.GetX() != 13.5 || p.GetY() != 13.5 {
		t.Fatalf("points[1] = (%v,%v), expect (13.5,13.5)", p.GetX(), p.GetY())
	}
	// topRight
	if p := r.GetPoints()[2]; p.GetX() != 35.5 || p.GetY() != 13.5 {
		t.Fatalf("points[2] = (%v,%v), expect (35.5,13.5)", p.GetX(), p.GetY())
	}
	// alignment pattern
	if p := r.GetPoints()[3]; p.GetX() != 32.5 || p.GetY() != 32.5 {
		t.Fatalf("points[3] = (%v,%v), expect (32.5,32.5)", p.GetX(), p.GetY())
	}

	for _, p := range cbpoints {
		if !(p.GetX() == 13.5 && p.GetY() == 35.5) && // bottomLeft
			!(p.GetX() == 13.5 && p.GetY() == 13.5) && // topLeft
			!(p.GetX() == 35.5 && p.GetY() == 13.5) && // topRight
			!(p.GetX() == 32.5 && p.GetY() == 32.5) { // alignmentPattern
			t.Fatalf("invalid callbacked point (%v,%v)", p.GetX(), p.GetY())
		}
	}
}

func TestDetector_DetectWithoutHint(t *testing.T) {

	image, _ := gozxing.NewBitMatrix(120, 120)

	d := NewDetector(image)
	_, e := d.DetectWithoutHints()
	if e == nil {
		t.Fatalf("DetectWithoutHints must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DetectWithoutHints must be NotFoundException, %v", e)
	}

	// version3 pattern, transformed, x3
	makePattern3(image, 20, 20)
	makePattern3(image, 96, 21)
	makePattern3(image, 22, 97)
	makeAlignPattern3(image, 89, 88)
	r, e := d.DetectWithoutHints()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}

	// bottomLeft
	if p := r.GetPoints()[0]; p.GetX() != 22.5 || p.GetY() != 97.5 {
		t.Fatalf("points[0] = (%v,%v), expect (22.5,97.5)", p.GetX(), p.GetY())
	}
	// topLeft
	if p := r.GetPoints()[1]; p.GetX() != 20.5 || p.GetY() != 20.5 {
		t.Fatalf("points[1] = (%v,%v), expect (20.5,20.5)", p.GetX(), p.GetY())
	}
	// topRight
	if p := r.GetPoints()[2]; p.GetX() != 96.5 || p.GetY() != 21.5 {
		t.Fatalf("points[2] = (%v,%v), expect (96.5,21.5)", p.GetX(), p.GetY())
	}
	// alignment pattern
	if p := r.GetPoints()[3]; p.GetX() != 89.5 || p.GetY() != 88.5 {
		t.Fatalf("points[3] = (%v,%v), expect (32.5,32.5)", p.GetX(), p.GetY())
	}

	bits := r.GetBits()
	normalizedstr := "" +
		"X X X X X X X                                       X X X X X X X \n" +
		"X           X                                       X           X \n" +
		"X   X X X   X                                       X   X X X   X \n" +
		"X   X X X   X                                       X   X X X   X \n" +
		"X   X X X   X                                       X   X X X   X \n" +
		"X           X                                       X           X \n" +
		"X X X X X X X                                       X X X X X X X \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                                  \n" +
		"                                                X X X X X         \n" +
		"                                                X       X         \n" +
		"X X X X X X X                                   X   X   X         \n" +
		"X           X                                   X       X         \n" +
		"X   X X X   X                                   X X X X X         \n" +
		"X   X X X   X                                                     \n" +
		"X   X X X   X                                                     \n" +
		"X           X                                                     \n" +
		"X X X X X X X                                                     \n"

	normalized, _ := gozxing.ParseStringToBitMatrix(normalizedstr, "X ", "  ")

	if normalized.GetWidth() != bits.GetWidth() || normalized.GetHeight() != bits.GetHeight() {
		t.Fatalf("bits = %vx%v, normalized = %vx%v",
			bits.GetWidth(), bits.GetHeight(), normalized.GetWidth(), normalized.GetHeight())
	}

	for j := 0; j < bits.GetHeight(); j++ {
		for i := 0; i < bits.GetWidth(); i++ {
			if bits.Get(i, j) != normalized.Get(i, j) {
				t.Fatalf("bits[%v,%v] = %v, expect %v", i, j, bits.Get(i, j), normalized.Get(i, j))
			}
		}
	}
}
