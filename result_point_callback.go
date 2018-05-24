package gozxing

type ResultPointCallback interface {
	FoundPossibleResultPoint(point ResultPoint)
}
