package detector

import (
	"math"
	"sort"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

const (
	FinderPatternFinder_CENTER_QUORUM = 2
	FinderPatternFinder_MIN_SKIP      = 3
	FinderPatternFinder_MAX_MODULES   = 97
)

type FinderPatternFinder struct {
	image                *common.BitMatrix
	possibleCenters      []*FinderPattern
	hasSkipped           bool
	crossCheckStateCount []int
	resultPointCallback  gozxing.ResultPointCallback
}

func NewFinderPatternFinder(image *common.BitMatrix, resultPointCallback gozxing.ResultPointCallback) *FinderPatternFinder {
	return &FinderPatternFinder{
		image:                image,
		possibleCenters:      make([]*FinderPattern, 0),
		crossCheckStateCount: make([]int, 5),
		resultPointCallback:  resultPointCallback,
	}
}

func (f *FinderPatternFinder) GetImage() *common.BitMatrix {
	return f.image
}

func (f *FinderPatternFinder) GetPossibleCenters() []*FinderPattern {
	return f.possibleCenters
}

func (f *FinderPatternFinder) Find(hints map[gozxing.DecodeHintType]interface{}) (*FinderPatternInfo, error) {
	_, tryHarder := hints[gozxing.DecodeHintType_TRY_HARDER]
	maxI := f.image.GetHeight()
	maxJ := f.image.GetWidth()

	iSkip := (3 * maxI) / (4 * FinderPatternFinder_MAX_MODULES)
	if iSkip < FinderPatternFinder_MIN_SKIP || tryHarder {
		iSkip = FinderPatternFinder_MIN_SKIP
	}

	done := false
	stateCount := make([]int, 5)
	for i := iSkip - 1; i < maxI && !done; i += iSkip {
		FinderPatternFinder_ClearCounts(stateCount)
		currentState := 0
		for j := 0; j < maxJ; j++ {
			if f.image.Get(j, i) {
				if (currentState & 1) == 1 {
					currentState++
				}
				stateCount[currentState]++
			} else {
				if (currentState & 1) == 0 {
					if currentState == 4 {
						if FinderPatternFinder_foundPatternCross(stateCount) {
							confirmed := f.HandlePossibleCenter(stateCount, i, j)
							if confirmed {
								iSkip = 2
								if f.hasSkipped {
									done = f.HaveMultiplyConfirmedCenters()
								} else {
									rowSkip := f.FindRowSkip()
									if rowSkip > stateCount[2] {
										i += rowSkip - stateCount[2] - iSkip
										j = maxJ - 1
									}
								}
							} else {
								FinderPatternFinder_ShiftCounts2(stateCount)
								currentState = 3
								continue
							}
							currentState = 0
							FinderPatternFinder_ClearCounts(stateCount)
						} else {
							FinderPatternFinder_ShiftCounts2(stateCount)
							currentState = 3
						}
					} else {
						currentState++
						stateCount[currentState]++
					}
				} else {
					stateCount[currentState]++
				}
			}
		}
		if FinderPatternFinder_foundPatternCross(stateCount) {
			confirmed := f.HandlePossibleCenter(stateCount, i, maxJ)
			if confirmed {
				iSkip = stateCount[0]
				if f.hasSkipped {
					done = f.HaveMultiplyConfirmedCenters()
				}
			}
		}
	}

	patternInfo, e := f.SelectBestPatterns()
	if e != nil {
		return nil, e
	}

	gozxing.ResultPoint_OrderBestPatterns(patternInfo)

	return NewFinderPatternInfo(patternInfo), nil
}

func FinderPatternFinder_centerFromEnd(stateCount []int, end int) float64 {
	return float64(end-stateCount[4]-stateCount[3]) - float64(stateCount[2])/2
}

func FinderPatternFinder_foundPatternCross(stateCount []int) bool {
	totalModuleSize := 0
	for i := 0; i < 5; i++ {
		count := stateCount[i]
		if count == 0 {
			return false
		}
		totalModuleSize += count
	}
	if totalModuleSize < 7 {
		return false
	}
	moduleSize := float64(totalModuleSize) / 7.0
	maxVariance := moduleSize / 2.0
	return math.Abs(moduleSize-float64(stateCount[0])) < maxVariance &&
		math.Abs(moduleSize-float64(stateCount[1])) < maxVariance &&
		math.Abs(3.0*moduleSize-float64(stateCount[2])) < 3*maxVariance &&
		math.Abs(moduleSize-float64(stateCount[3])) < maxVariance &&
		math.Abs(moduleSize-float64(stateCount[4])) < maxVariance
}

func FinderPatternFinder_foundPatternDiagonal(stateCount []int) bool {
	totalModuleSize := 0
	for i := 0; i < 5; i++ {
		count := stateCount[i]
		if count == 0 {
			return false
		}
		totalModuleSize += count
	}
	if totalModuleSize < 7 {
		return false
	}
	moduleSize := float64(totalModuleSize) / 7.0
	maxVariance := moduleSize / 1.333
	return math.Abs(moduleSize-float64(stateCount[0])) < maxVariance &&
		math.Abs(moduleSize-float64(stateCount[1])) < maxVariance &&
		math.Abs(3.0*moduleSize-float64(stateCount[2])) < 3*maxVariance &&
		math.Abs(moduleSize-float64(stateCount[3])) < maxVariance &&
		math.Abs(moduleSize-float64(stateCount[4])) < maxVariance
}

func (f *FinderPatternFinder) GetCrossCheckStateCount() []int {
	FinderPatternFinder_ClearCounts(f.crossCheckStateCount)
	return f.crossCheckStateCount
}

func FinderPatternFinder_ClearCounts(counts []int) {
	for x := 0; x < len(counts); x++ {
		counts[x] = 0
	}
}

func FinderPatternFinder_ShiftCounts2(stateCount []int) {
	stateCount[0] = stateCount[2]
	stateCount[1] = stateCount[3]
	stateCount[2] = stateCount[4]
	stateCount[3] = 1
	stateCount[4] = 0
}

func (f *FinderPatternFinder) crossCheckDiagonal(centerI, centerJ int) bool {
	stateCount := f.GetCrossCheckStateCount()

	i := 0
	for centerI >= i && centerJ >= i && f.image.Get(centerJ-i, centerI-i) {
		stateCount[2]++
		i++
	}
	if stateCount[2] == 0 {
		return false
	}

	for centerI >= i && centerJ >= i && !f.image.Get(centerJ-i, centerI-i) {
		stateCount[1]++
		i++
	}
	if stateCount[1] == 0 {
		return false
	}

	for centerI >= i && centerJ >= i && f.image.Get(centerJ-i, centerI-i) {
		stateCount[0]++
		i++
	}
	if stateCount[0] == 0 {
		return false
	}

	maxI := f.image.GetHeight()
	maxJ := f.image.GetWidth()

	i = 1
	for centerI+i < maxI && centerJ+i < maxJ && f.image.Get(centerJ+i, centerI+i) {
		stateCount[2]++
		i++
	}
	for centerI+i < maxI && centerJ+i < maxJ && !f.image.Get(centerJ+i, centerI+i) {
		stateCount[3]++
		i++
	}
	if stateCount[3] == 0 {
		return false
	}

	for centerI+i < maxI && centerJ+i < maxJ && f.image.Get(centerJ+i, centerI+i) {
		stateCount[4]++
		i++
	}
	if stateCount[4] == 0 {
		return false
	}

	return FinderPatternFinder_foundPatternDiagonal(stateCount)
}

func (f *FinderPatternFinder) CrossCheckVertical(startI, centerJ, maxCount, originalStateCountTotal int) float64 {
	image := f.image

	maxI := image.GetHeight()
	stateCount := f.GetCrossCheckStateCount()

	// Start counting up from center
	i := startI

	for i >= 0 && image.Get(centerJ, i) {
		stateCount[2]++
		i--
	}
	if i < 0 {
		return math.NaN()
	}
	for i >= 0 && !image.Get(centerJ, i) && stateCount[1] <= maxCount {
		stateCount[1]++
		i--
	}

	if i < 0 || stateCount[1] > maxCount {
		return math.NaN()
	}
	for i >= 0 && image.Get(centerJ, i) && stateCount[0] <= maxCount {
		stateCount[0]++
		i--
	}
	if stateCount[0] > maxCount {
		return math.NaN()
	}

	// Now also count down from center
	i = startI + 1
	for i < maxI && image.Get(centerJ, i) {
		stateCount[2]++
		i++
	}
	if i == maxI {
		return math.NaN()
	}
	for i < maxI && !image.Get(centerJ, i) && stateCount[3] < maxCount {
		stateCount[3]++
		i++
	}
	if i == maxI || stateCount[3] >= maxCount {
		return math.NaN()
	}
	for i < maxI && image.Get(centerJ, i) && stateCount[4] < maxCount {
		stateCount[4]++
		i++
	}
	if stateCount[4] >= maxCount {
		return math.NaN()
	}

	stateCountTotal := stateCount[0] + stateCount[1] + stateCount[2] + stateCount[3] + stateCount[4]
	if 5*math.Abs(float64(stateCountTotal-originalStateCountTotal)) >= float64(2*originalStateCountTotal) {
		return math.NaN()
	}

	if FinderPatternFinder_foundPatternCross(stateCount) {
		return FinderPatternFinder_centerFromEnd(stateCount, i)
	}
	return math.NaN()
}

func (f *FinderPatternFinder) CrossCheckHorizontal(startJ, centerI, maxCount, originalStateCountTotal int) float64 {
	image := f.image

	maxJ := image.GetWidth()
	stateCount := f.GetCrossCheckStateCount()

	j := startJ
	for j >= 0 && image.Get(j, centerI) {
		stateCount[2]++
		j--
	}
	if j < 0 {
		return math.NaN()
	}
	for j >= 0 && !image.Get(j, centerI) && stateCount[1] <= maxCount {
		stateCount[1]++
		j--
	}
	if j < 0 || stateCount[1] > maxCount {
		return math.NaN()
	}
	for j >= 0 && image.Get(j, centerI) && stateCount[0] <= maxCount {
		stateCount[0]++
		j--
	}
	if stateCount[0] > maxCount {
		return math.NaN()
	}

	j = startJ + 1
	for j < maxJ && image.Get(j, centerI) {
		stateCount[2]++
		j++
	}
	if j == maxJ {
		return math.NaN()
	}
	for j < maxJ && !image.Get(j, centerI) && stateCount[3] < maxCount {
		stateCount[3]++
		j++
	}
	if j == maxJ || stateCount[3] >= maxCount {
		return math.NaN()
	}
	for j < maxJ && image.Get(j, centerI) && stateCount[4] < maxCount {
		stateCount[4]++
		j++
	}
	if stateCount[4] >= maxCount {
		return math.NaN()
	}

	stateCountTotal := stateCount[0] + stateCount[1] + stateCount[2] + stateCount[3] + stateCount[4]
	if 5*math.Abs(float64(stateCountTotal-originalStateCountTotal)) >= float64(originalStateCountTotal) {
		return math.NaN()
	}

	if FinderPatternFinder_foundPatternCross(stateCount) {
		return FinderPatternFinder_centerFromEnd(stateCount, j)
	}
	return math.NaN()
}

func (f *FinderPatternFinder) HandlePossibleCenterWithPureBarcode(stateCount []int, i, j int, pureBarcode bool) bool {
	return f.HandlePossibleCenter(stateCount, i, j)
}

func (f *FinderPatternFinder) HandlePossibleCenter(stateCount []int, i, j int) bool {
	stateCountTotal := stateCount[0] + stateCount[1] + stateCount[2] + stateCount[3] + stateCount[4]
	centerJ := FinderPatternFinder_centerFromEnd(stateCount, j)

	centerI := f.CrossCheckVertical(i, int(centerJ), stateCount[2], stateCountTotal)
	if !math.IsNaN(centerI) {
		centerJ = f.CrossCheckHorizontal(int(centerJ), int(centerI), stateCount[2], stateCountTotal)
		if !math.IsNaN(centerJ) && f.crossCheckDiagonal(int(centerI), int(centerJ)) {
			estimatedModuleSize := float64(stateCountTotal) / 7.0
			found := false
			for index := 0; index < len(f.possibleCenters); index++ {
				center := f.possibleCenters[index]

				if center.AboutEquals(estimatedModuleSize, centerI, centerJ) {
					f.possibleCenters[index] = center.CombineEstimate(centerI, centerJ, estimatedModuleSize)
					found = true
					break
				}
			}
			if !found {
				point := NewFinderPattern1(centerJ, centerI, estimatedModuleSize)
				f.possibleCenters = append(f.possibleCenters, point)
				if f.resultPointCallback != nil {
					f.resultPointCallback(point)
				}
			}
			return true
		}
	}
	return false
}

func (f *FinderPatternFinder) FindRowSkip() int {
	if len(f.possibleCenters) <= 1 {
		return 0
	}
	var firstConfirmedCenter *FinderPattern
	for _, center := range f.possibleCenters {
		if center.GetCount() >= FinderPatternFinder_CENTER_QUORUM {
			if firstConfirmedCenter == nil {
				firstConfirmedCenter = center
			} else {
				f.hasSkipped = true
				return int((math.Abs(firstConfirmedCenter.GetX()-center.GetX()) -
					math.Abs(firstConfirmedCenter.GetY()-center.GetY())) / 2)
			}
		}
	}
	return 0
}

func (f *FinderPatternFinder) HaveMultiplyConfirmedCenters() bool {
	confirmedCount := 0
	totalModuleSize := 0.0
	max := len(f.possibleCenters)
	for _, pattern := range f.possibleCenters {
		if pattern.GetCount() >= FinderPatternFinder_CENTER_QUORUM {
			confirmedCount++
			totalModuleSize += pattern.GetEstimatedModuleSize()
		}
	}
	if confirmedCount < 3 {
		return false
	}

	average := totalModuleSize / float64(max)
	totalDeviation := 0.0
	for _, pattern := range f.possibleCenters {
		totalDeviation += math.Abs(pattern.GetEstimatedModuleSize() - average)
	}
	return totalDeviation <= 0.05*totalModuleSize
}

func (f *FinderPatternFinder) SelectBestPatterns() ([]gozxing.ResultPoint, error) {
	startSize := float64(len(f.possibleCenters))
	if startSize < 3 {
		return nil, gozxing.GetNotFoundExceptionInstance()
	}

	if startSize > 3 {
		totalModuleSize := 0.0
		square := 0.0
		for _, center := range f.possibleCenters {
			size := center.GetEstimatedModuleSize()
			totalModuleSize += size
			square += size * size
		}
		average := totalModuleSize / startSize
		stdDev := math.Sqrt(square/startSize - average*average)

		sort.Slice(f.possibleCenters, func(i, j int) bool {
			a := math.Abs(f.possibleCenters[i].GetEstimatedModuleSize() - average)
			b := math.Abs(f.possibleCenters[j].GetEstimatedModuleSize() - average)
			return b < a
		})

		limit := math.Max(0.2*average, stdDev)

		for i := 0; i < len(f.possibleCenters) && len(f.possibleCenters) > 3; i++ {
			pattern := f.possibleCenters[i]
			if math.Abs(pattern.GetEstimatedModuleSize()-average) > limit {
				f.possibleCenters = append(f.possibleCenters[:i], f.possibleCenters[i+1:]...)
				i--
			}
		}
	}

	if len(f.possibleCenters) > 3 {
		totalModuleSize := 0.0
		for _, possibleCenter := range f.possibleCenters {
			totalModuleSize += possibleCenter.GetEstimatedModuleSize()
		}

		average := totalModuleSize / float64(len(f.possibleCenters))

		sort.Slice(f.possibleCenters, func(i, j int) bool {
			center1 := f.possibleCenters[i]
			center2 := f.possibleCenters[j]

			countCompare := center2.GetCount() - center1.GetCount()
			if countCompare == 0 {
				a := math.Abs(f.possibleCenters[i].GetEstimatedModuleSize() - average)
				b := math.Abs(f.possibleCenters[j].GetEstimatedModuleSize() - average)
				return a < b
			}
			return countCompare < 0
		})

		f.possibleCenters = f.possibleCenters[:3]
	}

	return []gozxing.ResultPoint{
		f.possibleCenters[0],
		f.possibleCenters[1],
		f.possibleCenters[2]}, nil
}
