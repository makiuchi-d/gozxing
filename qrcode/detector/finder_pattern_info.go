package detector

type FinderPatternInfo struct {
	bottomLeft FinderPattern
	topLeft    FinderPattern
	topRight   FinderPattern
}

func NewFinderPatternInfo(patternCenters []FinderPattern) *FinderPatternInfo {
	return &FinderPatternInfo{
		bottomLeft: patternCenters[0],
		topLeft:    patternCenters[1],
		topRight:   patternCenters[2],
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
