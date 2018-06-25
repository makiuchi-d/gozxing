package common

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestDetectorResult(t *testing.T) {
	bits, _ := gozxing.NewBitMatrix(10, 10)
	points := make([]gozxing.ResultPoint, 5)

	dr := NewDetectorResult(bits, points)

	if dr == nil {
		t.Fatalf("NewDetectorResult returns nil")
	}

	if r := dr.GetBits(); r != bits {
		t.Fatalf("GetBits returns %p, expect %p", r, bits)
	}

	p := dr.GetPoints()
	if len(p) != len(points) {
		t.Fatalf("GetPoints slice length is %v, expect %v", len(p), len(points))
	}
	if &p[0] != &points[0] {
		t.Fatalf("GetPoints slice top address is %p, expect %p", &p[0], &points[0])
	}

}
