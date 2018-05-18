package gozxing

type ResultPointCallback interface {
	foundPossibleResultPoint(point ResultPoint)
}
