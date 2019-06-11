package rss

import (
	"github.com/makiuchi-d/gozxing"
)

// Encapsulates an RSS barcode finder pattern, including its start/end position and row.

type FinderPattern struct {
	value        int
	startEnd     []int
	resultPoints []gozxing.ResultPoint
}

func NewFinderPattern(value int, startEnd []int, start, end, rowNumber int) *FinderPattern {
	return &FinderPattern{
		value:    value,
		startEnd: startEnd,
		resultPoints: []gozxing.ResultPoint{
			gozxing.NewResultPoint(float64(start), float64(rowNumber)),
			gozxing.NewResultPoint(float64(end), float64(rowNumber)),
		},
	}
}

func (this *FinderPattern) GetValue() int {
	return this.value
}

func (this *FinderPattern) GetStartEnd() []int {
	return this.startEnd
}

func (this *FinderPattern) GetResultPoints() []gozxing.ResultPoint {
	return this.resultPoints
}

func (this *FinderPattern) Equals(o interface{}) bool {
	that, ok := o.(*FinderPattern)
	if !ok {
		return false
	}
	return this.value == that.value
}

func (this *FinderPattern) HashCode() int {
	return this.value
}
