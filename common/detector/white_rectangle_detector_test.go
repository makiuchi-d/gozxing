package detector

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestNewWhiteRectangleDetector(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(20, 20)

	_, e := NewWhiteRectangleDetector(img, 10, 3, 3)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("NewWhiteRectangleDetector must be NotFoundException, %T", e)
	}

	detector, e := NewWhiteRectangleDetectorFromImage(img)
	if e != nil {
		t.Fatalf("NewWhiteRectangleDetectorFromImage returns error, %v", e)
	}
	if detector.image != img {
		t.Fatalf("image != img, %p, %p", detector.image, img)
	}
	if detector.width != 20 || detector.height != 20 {
		t.Fatalf("width,height = %v,%v, expect 20, 20", detector.width, detector.height)
	}
	if detector.leftInit != 5 || detector.rightInit != 15 || detector.upInit != 5 || detector.downInit != 15 {
		t.Fatalf("init = %v,%v,%v,%v, expect 5,15,5,15",
			detector.leftInit, detector.rightInit, detector.upInit, detector.downInit)
	}
}

func TestWhiterectangledetector_getBlackPointOnSegment(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(20, 20)
	img.Set(10, 10)
	detector, _ := NewWhiteRectangleDetectorFromImage(img)

	p := detector.getBlackPointOnSegment(5, 8, 15, 12)
	if p == nil {
		t.Fatalf("getBlackPointOnSegment returns nil")
	}
	if x, y := p.GetX(), p.GetY(); x != 10 || y != 10 {
		t.Fatalf("getBlackPointOnSegment = %v,%v, expect 10,10", x, y)
	}

	p = detector.getBlackPointOnSegment(5, 8, 15, 10)
	if p != nil {
		t.Fatalf("getBlackPointOnSegment must be nil, %v", p)
	}
}

func TestWhiterectangledetector_centerEdges(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(25, 25)
	decoder, _ := NewWhiteRectangleDetectorFromImage(img)

	t0 := gozxing.NewResultPoint(15, 3)
	z0 := gozxing.NewResultPoint(3, 5)
	x0 := gozxing.NewResultPoint(20, 10)
	y0 := gozxing.NewResultPoint(5, 20)

	edges := decoder.centerEdges(y0, z0, x0, t0)

	if p := edges[0]; p.GetX() != 14 || p.GetY() != 4 {
		t.Fatalf("edges[0] = (%v,%v), expect (14,4)", p.GetX(), p.GetY())
	}
	if p := edges[1]; p.GetX() != 4 || p.GetY() != 6 {
		t.Fatalf("edges[1] = (%v,%v), expect (4,6)", p.GetX(), p.GetY())
	}
	if p := edges[2]; p.GetX() != 19 || p.GetY() != 9 {
		t.Fatalf("edges[2] = (%v,%v), expect (19,9)", p.GetX(), p.GetY())
	}
	if p := edges[3]; p.GetX() != 6 || p.GetY() != 19 {
		t.Fatalf("edges[3] = (%v,%v), expect (6,19)", p.GetX(), p.GetY())
	}

	t0 = gozxing.NewResultPoint(5, 3)
	z0 = gozxing.NewResultPoint(3, 15)
	x0 = gozxing.NewResultPoint(20, 5)
	y0 = gozxing.NewResultPoint(15, 20)

	edges = decoder.centerEdges(y0, z0, x0, t0)

	if p := edges[0]; p.GetX() != 6 || p.GetY() != 4 {
		t.Fatalf("edges[0] = (%v,%v), expect (6,4)", p.GetX(), p.GetY())
	}
	if p := edges[1]; p.GetX() != 4 || p.GetY() != 14 {
		t.Fatalf("edges[1] = (%v,%v), expect (4,14)", p.GetX(), p.GetY())
	}
	if p := edges[2]; p.GetX() != 19 || p.GetY() != 6 {
		t.Fatalf("edges[2] = (%v,%v), expect (19,6)", p.GetX(), p.GetY())
	}
	if p := edges[3]; p.GetX() != 14 || p.GetY() != 19 {
		t.Fatalf("edges[3] = (%v,%v), expect (14,19)", p.GetX(), p.GetY())
	}
}

func TestWhiterectangledetector_containsBlackPoint(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(20, 20)
	img.Set(10, 10)
	detector, _ := NewWhiteRectangleDetectorFromImage(img)

	if !detector.containsBlackPoint(5, 15, 10, true) {
		t.Fatalf("containsBlackPoint(5, 15, 10, true) must be true")
	}

	if !detector.containsBlackPoint(5, 15, 10, false) {
		t.Fatalf("containsBlackPoint(5, 15, 10, false) must be true")
	}

	if detector.containsBlackPoint(5, 15, 5, true) {
		t.Fatalf("containsBlackPoint(5, 15, 5, true) must be false")
	}
}

func TestWhiterectangledetector_Detect(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(20, 20)
	detector, _ := NewWhiteRectangleDetectorFromImage(img)

	_, e := detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}

	// right
	img.Set(18, 10)

	_, e = detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}

	// bottom
	img.Set(14, 18)

	_, e = detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}

	// left
	img.Set(3, 15)

	_, e = detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}

	// up
	img.Set(8, 3)

	points, e := detector.Detect()
	if e != nil {
		t.Fatalf("Detect returns error, %v", e)
	}

	// t
	if p := points[0]; p.GetX() != 9 || p.GetY() != 4 {
		t.Fatalf("Detect points[0] = (%v,%v), expect (9,4)", p.GetX(), p.GetY())
	}
	// z
	if p := points[1]; p.GetX() != 4 || p.GetY() != 14 {
		t.Fatalf("Detect points[1] = (%v,%v), expect (4,14)", p.GetX(), p.GetY())
	}
	// x
	if p := points[2]; p.GetX() != 17 || p.GetY() != 11 {
		t.Fatalf("Detect points[2] = (%v,%v), expect (17,11)", p.GetX(), p.GetY())
	}
	// y
	if p := points[3]; p.GetX() != 13 || p.GetY() != 17 {
		t.Fatalf("Detect points[3] = (%v,%v), expect (13,17)", p.GetX(), p.GetY())
	}
}

func TestWhiterectangledetector_DetectFail(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(20, 20)
	detector, _ := NewWhiteRectangleDetectorFromImage(img)

	img.Set(5, 5)
	img.Set(15, 15)

	// z not found
	_, e := detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}

	img.Set(7, 12)

	// x not found
	_, e = detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}

	img.Clear()
	img.Set(5, 15)
	img.Set(15, 5)

	// t not found
	_, e = detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}

	// y not found
	img.Set(6, 6)
	_, e = detector.Detect()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Detect must be NotFoundException, %T", e)
	}
}
