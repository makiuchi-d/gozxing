package detector

import (
	"github.com/makiuchi-d/gozxing"
)

type FinderPatternInfo struct {
	bottomLeft FinderPattern
	topLeft    FinderPattern
	topRight   FinderPattern
}

func NewFinderPatternInfo(patternCenters []gozxing.ResultPoint) *FinderPatternInfo {
	return &FinderPatternInfo{
		bottomLeft: patternCenters[0].(FinderPattern),
		topLeft:    patternCenters[1].(FinderPattern),
		topRight:   patternCenters[2].(FinderPattern),
	}
}

func (f *FinderPatternInfo) GetBottomLeft() FinderPattern {
	return f.bottomLeft
}

func (f *FinderPatternInfo) GetTopLeft() FinderPattern {
	return f.topLeft
}
func (f *FinderPatternInfo) GetTopRight() FinderPattern {
	return f.topRight
}
