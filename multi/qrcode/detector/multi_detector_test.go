package detector

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

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

	testsPoints := []struct{ tlx, tly, blx, bly, trx, try float64 }{
		{3.5, 3.5, 3.5, 17.5, 17.5, 3.5},     // TopLeft
		{28.5, 3.5, 28.5, 17.5, 42.5, 3.5},   // TopRight
		{28.5, 3.5, 28.5, 17.5, 42.5, 3.5},   // BottomLeft
		{28.5, 28.5, 28.5, 42.5, 42.5, 28.5}, // BottomLeft
	}
FORTESTPOINTS:
	for _, test := range testsPoints {
		for _, r := range results {
			ps := r.GetPoints()
			bl, tl, tr := ps[0], ps[1], ps[2]

			if bl.GetX() == test.blx && bl.GetY() == test.bly &&
				tl.GetX() == test.tlx && tl.GetY() == test.tly &&
				tr.GetX() == test.trx && tr.GetY() == test.try {
				continue FORTESTPOINTS
			}
		}

		t.Fatalf("result must contain: {(%v,%v), (%v,%v), (%v,%v)}",
			test.tlx, test.tly, test.blx, test.bly, test.trx, test.try)
	}

	testCallbacks := []struct{ x, y float64 }{
		{3.5, 3.5}, {3.5, 17.5}, {17.5, 3.5}, // TopLeft
		{28.5, 3.5}, {28.5, 17.5}, {42.5, 3.5}, // TopRight
		{3.5, 28.5}, {3.5, 42.5}, {17.5, 28.5}, // BottomLeft
		{28.5, 28.5}, {28.5, 42.5}, {42.5, 28.5}, // BottomRight
	}
FORTESTCALLBACKS:
	for _, test := range testCallbacks {
		for _, p := range points {
			if p.GetX() == test.x && p.GetY() == test.y {
				continue FORTESTCALLBACKS
			}
		}
		t.Fatalf("callbacked points must contain %v", test)
	}
}
