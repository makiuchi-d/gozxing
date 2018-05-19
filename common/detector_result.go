package common

import (
	"github.com/makiuchi-d/gozxing"
)

type DetectorResult struct {
	bits   *BitMatrix
	points []gozxing.ResultPoint
}

func NewDetectorResult(bits *BitMatrix, points []gozxing.ResultPoint) *DetectorResult {
	return &DetectorResult{bits, points}
}

func (d *DetectorResult) GetBits() *BitMatrix {
	return d.bits
}

func (d *DetectorResult) GetPoints() []gozxing.ResultPoint {
	return d.points
}
