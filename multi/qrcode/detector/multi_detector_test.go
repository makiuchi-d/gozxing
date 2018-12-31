package detector

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

func containsPoints(points []*common.DetectorResult, tlx, tly, blx, bly, trx, try float64) bool {
	for _, r := range points {
		ps := r.GetPoints()
		bl, tl, tr := ps[0], ps[1], ps[2]

		if bl.GetX() == blx && bl.GetY() == bly &&
			tl.GetX() == tlx && tl.GetY() == tly &&
			tr.GetX() == trx && tr.GetY() == try {
			return true
		}
	}

	return false
}

func containsPoint(points []gozxing.ResultPoint, x, y float64) bool {
	for _, p := range points {
		if p.GetX() == x && p.GetY() == y {
			return true
		}
	}
	return false
}

func TestMultiDetector_DetectMulti(t *testing.T) {
	hints := make(map[gozxing.DecodeHintType]interface{})

	img, _ := gozxing.NewBitMatrix(10, 10)
	det := NewMultiDetector(img)
	_, e := det.DetectMulti(hints)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DetectMulti must be NotFoundException: %T", e)
	}

	img, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	det = NewMultiDetector(img)
	points := make([]gozxing.ResultPoint, 0)
	hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK] = gozxing.ResultPointCallback(
		func(p gozxing.ResultPoint) {
			points = append(points, p)
		})
	results, e := det.DetectMulti(hints)
	if e != nil {
		t.Fatalf("DetectMulti returns error: %v", e)
	}

	// top left
	if !containsPoints(results, 3.5, 3.5, 3.5, 17.5, 17.5, 3.5) {
		t.Fatalf("result must contain: {(3.5,3.5), (3.5,17.5), (17.5,3.5)}")
	}
	// top right
	if !containsPoints(results, 28.5, 3.5, 28.5, 17.5, 42.5, 3.5) {
		t.Fatalf("result must contain: {(28.5,3.5), (28.5,17.5), (42.5,3.5)}")
	}
	// bottom left
	if !containsPoints(results, 3.5, 28.5, 3.5, 42.5, 17.5, 28.5) {
		t.Fatalf("result must contain: {(3.5,28.5), (3.5,42.5), (17.5,28.5)}")
	}
	// bottom right
	if !containsPoints(results, 28.5, 28.5, 28.5, 42.5, 42.5, 28.5) {
		t.Fatalf("result must contain: {(28.5,28.5), (28.5,42.5), (42.5,28.5)}")
	}

	if !containsPoint(points, 3.5, 3.5) {
		t.Fatalf("callback point must contain (3.5,3.5)")
	}
	if !containsPoint(points, 3.5, 17.5) {
		t.Fatalf("callback point must contain (3.5,17.5)")
	}
	if !containsPoint(points, 17.5, 3.5) {
		t.Fatalf("callback point must contain (17.5,3.5)")
	}

	if !containsPoint(points, 28.5, 3.5) {
		t.Fatalf("callback point must contain (28.5,3.5)")
	}
	if !containsPoint(points, 28.5, 17.5) {
		t.Fatalf("callback point must contain (28.5,3.5)")
	}
	if !containsPoint(points, 42.5, 3.5) {
		t.Fatalf("callback point must contain (28.5,3.5)")
	}

	if !containsPoint(points, 3.5, 28.5) {
		t.Fatalf("callback point must contain (3.5,28.5)")
	}
	if !containsPoint(points, 3.5, 42.5) {
		t.Fatalf("callback point must contain (3.5,42.5)")
	}
	if !containsPoint(points, 17.5, 28.5) {
		t.Fatalf("callback point must contain (17.5,28.5)")
	}

	if !containsPoint(points, 28.5, 28.5) {
		t.Fatalf("callback point must contain (28.5,28.5)")
	}
	if !containsPoint(points, 28.5, 42.5) {
		t.Fatalf("callback point must contain (28.5,42.5)")
	}
	if !containsPoint(points, 42.5, 28.5) {
		t.Fatalf("callback point must contain (42.5,28.5)")
	}
}
