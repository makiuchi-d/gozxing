package detector

import (
	"math"
	"sort"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode/detector"
)

// This class attempts to find finder patterns in a QR Code.
// Finder patterns are the square markers at three corners of a QR Code.
//
// This class is thread-safe but not reentrant. Each thread must allocate its own object.
//
// In contrast to {@link FinderPatternFinder}, this class will return an array of all possible
// QR code locations in the image.
//
// Use the TRY_HARDER hint to ask for a more thorough detection.
//
type MultiFinderPatternFinder struct {
	*detector.FinderPatternFinder
}

// private static final FinderPatternInfo[] EMPTY_RESULT_ARRAY = new FinderPatternInfo[0];
// private static final FinderPattern[][] EMPTY_FP_2D_ARRAY = new FinderPattern[0][];

const (
	// TODO MIN_MODULE_COUNT and MAX_MODULE_COUNT would be great hints to ask the user for
	// since it limits the number of regions to decode

	// max. legal count of modules per QR code edge (177)
	MAX_MODULE_COUNT_PER_EDGE = 180
	// min. legal count per modules per QR code edge (11)
	MIN_MODULE_COUNT_PER_EDGE = 9

	// More or less arbitrary cutoff point for determining if two finder patterns might belong
	// to the same code if they differ less than DIFF_MODSIZE_CUTOFF_PERCENT percent in their
	// estimated modules sizes.
	DIFF_MODSIZE_CUTOFF_PERCENT = 0.05

	// More or less arbitrary cutoff point for determining if two finder patterns might belong
	// to the same code if they differ less than DIFF_MODSIZE_CUTOFF pixels/module in their
	// estimated modules sizes.
	DIFF_MODSIZE_CUTOFF = 0.5
)

// ModuleSizeComparator A comparator that orders FinderPatterns by their estimated module size.
func ModuleSizeComparator(possibleCenters []*detector.FinderPattern) func(int, int) bool {
	return func(i, j int) bool {
		center1 := possibleCenters[i]
		center2 := possibleCenters[j]
		value := center2.GetEstimatedModuleSize() - center1.GetEstimatedModuleSize()
		return value < 0
	}
}

// NewMultiFinderPatternFinder Creates a finder that will search the image for three finder patterns.
//
// @param image image to search
//
func NewMultiFinderPatternFinder(image *gozxing.BitMatrix, resultPointCallback gozxing.ResultPointCallback) *MultiFinderPatternFinder {
	return &MultiFinderPatternFinder{
		detector.NewFinderPatternFinder(image, resultPointCallback),
	}
}

// selectMultipleBestPatterns select the best patterns.
// @return the 3 best {@link FinderPattern}s from our list of candidates. The "best" are
//         those that have been detected at least 2 times, and whose module
//         size differs from the average among those patterns the least
// @throws NotFoundException if 3 such finder patterns do not exist
func (this *MultiFinderPatternFinder) selectMultipleBestPatterns() ([][]*detector.FinderPattern, error) {
	possibleCenters := this.GetPossibleCenters()
	size := len(possibleCenters)

	if size < 3 {
		// Couldn't find enough finder patterns
		return nil, gozxing.NewNotFoundException("Couldn't find enough finder patterns (%d)", size)
	}

	// Begin HE modifications to safely detect multiple codes of equal size
	if size == 3 {
		return [][]*detector.FinderPattern{
			{
				possibleCenters[0],
				possibleCenters[1],
				possibleCenters[2],
			},
		}, nil
	}

	// Sort by estimated module size to speed up the upcoming checks
	sort.Slice(possibleCenters, ModuleSizeComparator(possibleCenters))

	// Now lets start: build a list of tuples of three finder locations that
	//  - feature similar module sizes
	//  - are placed in a distance so the estimated module count is within the QR specification
	//  - have similar distance between upper left/right and left top/bottom finder patterns
	//  - form a triangle with 90° angle (checked by comparing top right/bottom left distance
	//    with pythagoras)
	//
	// Note: we allow each point to be used for more than one code region: this might seem
	// counterintuitive at first, but the performance penalty is not that big. At this point,
	// we cannot make a good quality decision whether the three finders actually represent
	// a QR code, or are just by chance laid out so it looks like there might be a QR code there.
	// So, if the layout seems right, lets have the decoder try to decode.

	results := make([][]*detector.FinderPattern, 0) // holder for the results

	for i1 := 0; i1 < (size - 2); i1++ {
		p1 := possibleCenters[i1]
		if p1 == nil {
			continue
		}

		for i2 := i1 + 1; i2 < (size - 1); i2++ {
			p2 := possibleCenters[i2]
			if p2 == nil {
				continue
			}

			// Compare the expected module sizes; if they are really off, skip
			vModSize12 := (p1.GetEstimatedModuleSize() - p2.GetEstimatedModuleSize()) /
				math.Min(p1.GetEstimatedModuleSize(), p2.GetEstimatedModuleSize())
			vModSize12A := math.Abs(p1.GetEstimatedModuleSize() - p2.GetEstimatedModuleSize())
			if vModSize12A > DIFF_MODSIZE_CUTOFF && vModSize12 >= DIFF_MODSIZE_CUTOFF_PERCENT {
				// break, since elements are ordered by the module size deviation there cannot be
				// any more interesting elements for the given p1.
				break
			}

			for i3 := i2 + 1; i3 < size; i3++ {
				p3 := possibleCenters[i3]
				if p3 == nil {
					continue
				}

				// Compare the expected module sizes; if they are really off, skip
				vModSize23 := (p2.GetEstimatedModuleSize() - p3.GetEstimatedModuleSize()) /
					math.Min(p2.GetEstimatedModuleSize(), p3.GetEstimatedModuleSize())
				vModSize23A := math.Abs(p2.GetEstimatedModuleSize() - p3.GetEstimatedModuleSize())
				if vModSize23A > DIFF_MODSIZE_CUTOFF && vModSize23 >= DIFF_MODSIZE_CUTOFF_PERCENT {
					// break, since elements are ordered by the module size deviation there cannot be
					// any more interesting elements for the given p1.
					break
				}

				bl, tl, tr := gozxing.ResultPoint_OrderBestPatterns(p1, p2, p3)
				test := []*detector.FinderPattern{
					bl.(*detector.FinderPattern), tl.(*detector.FinderPattern), tr.(*detector.FinderPattern),
				}

				// Calculate the distances: a = topleft-bottomleft, b=topleft-topright, c = diagonal
				info := detector.NewFinderPatternInfo(test[0], test[1], test[2])
				dA := gozxing.ResultPoint_Distance(info.GetTopLeft(), info.GetBottomLeft())
				dC := gozxing.ResultPoint_Distance(info.GetTopRight(), info.GetBottomLeft())
				dB := gozxing.ResultPoint_Distance(info.GetTopLeft(), info.GetTopRight())

				// Check the sizes
				estimatedModuleCount := (dA + dB) / (p1.GetEstimatedModuleSize() * 2.0)
				if estimatedModuleCount > MAX_MODULE_COUNT_PER_EDGE ||
					estimatedModuleCount < MIN_MODULE_COUNT_PER_EDGE {
					continue
				}

				// Calculate the difference of the edge lengths in percent
				vABBC := math.Abs((dA - dB) / math.Min(dA, dB))
				if vABBC >= 0.1 {
					continue
				}

				// Calculate the diagonal length by assuming a 90° angle at topleft
				dCpy := math.Sqrt(dA*dA + dB*dB)
				// Compare to the real distance in %
				vPyC := math.Abs((dC - dCpy) / math.Min(dC, dCpy))

				if vPyC >= 0.1 {
					continue
				}

				// All tests passed!
				results = append(results, test)
			}
		}
	}

	if len(results) > 0 {
		return results, nil
	}

	// Nothing found!
	return nil, gozxing.NewNotFoundException()
}

func (this *MultiFinderPatternFinder) FindMulti(hints map[gozxing.DecodeHintType]interface{}) ([]*detector.FinderPatternInfo, error) {
	_, tryHarder := hints[gozxing.DecodeHintType_TRY_HARDER]
	image := this.GetImage()
	maxI := image.GetHeight()
	maxJ := image.GetWidth()
	// We are looking for black/white/black/white/black modules in
	// 1:1:3:1:1 ratio; this tracks the number of such modules seen so far

	// Let's assume that the maximum version QR Code we support takes up 1/4 the height of the
	// image, and then account for the center being 3 modules in size. This gives the smallest
	// number of pixels the center could be, so skip this often. When trying harder, look for all
	// QR versions regardless of how dense they are.
	iSkip := (3 * maxI) / (4 * detector.FinderPatternFinder_MAX_MODULES)
	if iSkip < detector.FinderPatternFinder_MIN_SKIP || tryHarder {
		iSkip = detector.FinderPatternFinder_MIN_SKIP
	}

	stateCount := make([]int, 5)
	for i := iSkip - 1; i < maxI; i += iSkip {
		// Get a row of black/white values
		detector.FinderPatternFinder_doClearCounts(stateCount)
		currentState := 0
		for j := 0; j < maxJ; j++ {
			if image.Get(j, i) {
				// Black pixel
				if (currentState & 1) == 1 { // Counting white pixels
					currentState++
				}
				stateCount[currentState]++
			} else { // White pixel
				if (currentState & 1) == 0 { // Counting black pixels
					if currentState == 4 { // A winner?
						if detector.FinderPatternFinder_foundPatternCross(stateCount) &&
							this.HandlePossibleCenter(stateCount, i, j) { // Yes
							// Clear state to start looking again
							currentState = 0
							detector.FinderPatternFinder_doClearCounts(stateCount)
						} else { // No, shift counts back by two
							detector.FinderPatternFinder_doShiftCounts2(stateCount)
							currentState = 3
						}
					} else {
						currentState++
						stateCount[currentState]++
					}
				} else { // Counting white pixels
					stateCount[currentState]++
				}
			}
		} // for j=...

		if detector.FinderPatternFinder_foundPatternCross(stateCount) {
			this.HandlePossibleCenter(stateCount, i, maxJ)
		}
	} // for i=iSkip-1 ...
	patternInfo, e := this.selectMultipleBestPatterns()
	if e != nil {
		return nil, e
	}
	result := make([]*detector.FinderPatternInfo, 0)
	for _, pattern := range patternInfo {
		bl, tl, tr := gozxing.ResultPoint_OrderBestPatterns(pattern[0], pattern[1], pattern[2])
		result = append(result,
			detector.NewFinderPatternInfo(
				bl.(*detector.FinderPattern), tl.(*detector.FinderPattern), tr.(*detector.FinderPattern)))
	}

	return result, nil
}
