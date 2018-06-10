package decoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestQRCodeDecoderMetaData(t *testing.T) {
	points := []gozxing.ResultPoint{
		gozxing.NewResultPoint(10, 50),
		gozxing.NewResultPoint(10, 10),
		gozxing.NewResultPoint(50, 10),
	}

	md := NewQRCodeDecoderMetaData(false)
	if md.IsMirrored() {
		t.Fatalf("IsMirrored must be false")
	}
	md.ApplyMirroredCorrection(points)
	if points[0].GetX() != 10 || points[0].GetY() != 50 {
		t.Fatalf("points[0] = (%v,%v), expect (10,50)", points[0].GetX(), points[0].GetY())
	}
	if points[1].GetX() != 10 || points[1].GetY() != 10 {
		t.Fatalf("points[1] = (%v,%v), expect (10,10)", points[1].GetX(), points[1].GetY())
	}
	if points[2].GetX() != 50 || points[2].GetY() != 10 {
		t.Fatalf("points[2] = (%v,%v), expect (50,10)", points[2].GetX(), points[2].GetY())
	}

	md = NewQRCodeDecoderMetaData(true)
	if !md.IsMirrored() {
		t.Fatalf("IsMirrored must be true")
	}
	md.ApplyMirroredCorrection(points)
	if points[0].GetX() != 50 || points[0].GetY() != 10 {
		t.Fatalf("points[0] = (%v,%v), expect (50,10)", points[0].GetX(), points[0].GetY())
	}
	if points[1].GetX() != 10 || points[1].GetY() != 10 {
		t.Fatalf("points[1] = (%v,%v), expect (10,10)", points[1].GetX(), points[1].GetY())
	}
	if points[2].GetX() != 10 || points[2].GetY() != 50 {
		t.Fatalf("points[2] = (%v,%v), expect (10,50)", points[2].GetX(), points[2].GetY())
	}

}
