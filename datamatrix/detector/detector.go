package detector

import (
	"fmt"
	"sort"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	cdetector "github.com/makiuchi-d/gozxing/common/detector"
	"github.com/makiuchi-d/gozxing/common/util"
)

type Detector struct {
	image             *gozxing.BitMatrix
	rectangleDetector *cdetector.WhiteRectangleDetector
}

func NewDetector(image *gozxing.BitMatrix) (*Detector, error) {
	rectangleDetector, e := cdetector.NewWhiteRectangleDetectorFromImage(image)
	if e != nil {
		return nil, e
	}
	return &Detector{image, rectangleDetector}, nil
}

// Detect Detects a Data Matrix Code in an image.
//
// @return {@link DetectorResult} encapsulating results of detecting a Data Matrix Code
// @throws NotFoundException if no Data Matrix Code can be found
//
func (this *Detector) Detect() (*common.DetectorResult, error) {

	cornerPoints, e := this.rectangleDetector.Detect()
	if e != nil {
		return nil, e
	}
	pointA := cornerPoints[0]
	pointB := cornerPoints[1]
	pointC := cornerPoints[2]
	pointD := cornerPoints[3]

	// Point A and D are across the diagonal from one another,
	// as are B and C. Figure out which are the solid black lines
	// by counting transitions
	transitions := make([]*ResultPointsAndTransitions, 4)
	transitions[0] = this.transitionsBetween(pointA, pointB)
	transitions[1] = this.transitionsBetween(pointA, pointC)
	transitions[2] = this.transitionsBetween(pointB, pointD)
	transitions[3] = this.transitionsBetween(pointC, pointD)
	sort.Slice(transitions, func(a, b int) bool {
		return ResultPointsAndTransitionsComparator(transitions[a], transitions[b])
	})

	// Sort by number of transitions. First two will be the two solid sides; last two
	// will be the two alternating black/white sides
	lSideOne := transitions[0]
	lSideTwo := transitions[1]

	// Figure out which point is their intersection by tallying up the number of times we see the
	// endpoints in the four endpoints. One will show up twice.
	pointCount := make(map[gozxing.ResultPoint]int)
	pointCount = increment(pointCount, lSideOne.getFrom())
	pointCount = increment(pointCount, lSideOne.getTo())
	pointCount = increment(pointCount, lSideTwo.getFrom())
	pointCount = increment(pointCount, lSideTwo.getTo())

	var maybeTopLeft gozxing.ResultPoint
	var bottomLeft gozxing.ResultPoint
	var maybeBottomRight gozxing.ResultPoint
	for point, value := range pointCount {
		if value == 2 {
			bottomLeft = point // this is definitely the bottom left, then -- end of two L sides
		} else {
			// Otherwise it's either top left or bottom right -- just assign the two arbitrarily now
			if maybeTopLeft == nil {
				maybeTopLeft = point
			} else {
				maybeBottomRight = point
			}
		}
	}

	if maybeTopLeft == nil || bottomLeft == nil || maybeBottomRight == nil {
		return nil, gozxing.GetNotFoundExceptionInstance()
	}

	// Bottom left is correct but top left and bottom right might be switched

	// Use the dot product trick to sort them out

	// Now we know which is which:
	bottomRight, bottomLeft, topLeft := gozxing.ResultPoint_OrderBestPatterns(maybeTopLeft, bottomLeft, maybeBottomRight)
	// Which point didn't we find in relation to the "L" sides? that's the top right corner
	var topRight gozxing.ResultPoint
	if _, ok := pointCount[pointA]; !ok {
		topRight = pointA
	} else if _, ok := pointCount[pointB]; !ok {
		topRight = pointB
	} else if _, ok := pointCount[pointC]; !ok {
		topRight = pointC
	} else {
		topRight = pointD
	}

	// Next determine the dimension by tracing along the top or right side and counting black/white
	// transitions. Since we start inside a black module, we should see a number of transitions
	// equal to 1 less than the code dimension. Well, actually 2 less, because we are going to
	// end on a black module:

	// The top right point is actually the corner of a module, which is one of the two black modules
	// adjacent to the white module at the top right. Tracing to that corner from either the top left
	// or bottom right should work here.

	dimensionTop := this.transitionsBetween(topLeft, topRight).getTransitions()
	dimensionRight := this.transitionsBetween(bottomRight, topRight).getTransitions()

	if (dimensionTop & 0x01) == 1 {
		// it can't be odd, so, round... up?
		dimensionTop++
	}
	dimensionTop += 2

	if (dimensionRight & 0x01) == 1 {
		// it can't be odd, so, round... up?
		dimensionRight++
	}
	dimensionRight += 2

	var bits *gozxing.BitMatrix
	var correctedTopRight gozxing.ResultPoint

	// Rectangular symbols are 6x16, 6x28, 10x24, 10x32, 14x32, or 14x44. If one dimension is more
	// than twice the other, it's certainly rectangular, but to cut a bit more slack we accept it as
	// rectangular if the bigger side is at least 7/4 times the other:
	if 4*dimensionTop >= 7*dimensionRight || 4*dimensionRight >= 7*dimensionTop {
		// The matrix is rectangular

		correctedTopRight = this.correctTopRightRectangular(
			bottomLeft, bottomRight, topLeft, topRight, dimensionTop, dimensionRight)
		if correctedTopRight == nil {
			correctedTopRight = topRight
		}

		dimensionTop = this.transitionsBetween(topLeft, correctedTopRight).getTransitions()
		dimensionRight = this.transitionsBetween(bottomRight, correctedTopRight).getTransitions()

		if (dimensionTop & 0x01) == 1 {
			// it can't be odd, so, round... up?
			dimensionTop++
		}

		if (dimensionRight & 0x01) == 1 {
			// it can't be odd, so, round... up?
			dimensionRight++
		}

		bits, e = sampleGrid(
			this.image, topLeft, bottomLeft, bottomRight, correctedTopRight, dimensionTop, dimensionRight)
		if e != nil {
			return nil, e
		}

	} else {
		// The matrix is square

		dimension := min(dimensionRight, dimensionTop)
		// correct top right point to match the white module
		correctedTopRight = this.correctTopRight(bottomLeft, bottomRight, topLeft, topRight, dimension)
		if correctedTopRight == nil {
			correctedTopRight = topRight
		}

		// Redetermine the dimension using the corrected top right point
		dimensionCorrected := max(
			this.transitionsBetween(topLeft, correctedTopRight).getTransitions(),
			this.transitionsBetween(bottomRight, correctedTopRight).getTransitions())
		dimensionCorrected++
		if (dimensionCorrected & 0x01) == 1 {
			dimensionCorrected++
		}

		bits, e = sampleGrid(
			this.image,
			topLeft,
			bottomLeft,
			bottomRight,
			correctedTopRight,
			dimensionCorrected,
			dimensionCorrected)
		if e != nil {
			return nil, e
		}
	}

	return common.NewDetectorResult(bits, []gozxing.ResultPoint{topLeft, bottomLeft, bottomRight, correctedTopRight}), nil
}

// correctTopRightRectangular  Calculates the position of the white top right module
// using the output of the rectangle detector for a rectangular matrix
func (this *Detector) correctTopRightRectangular(
	bottomLeft, bottomRight, topLeft, topRight gozxing.ResultPoint,
	dimensionTop, dimensionRight int) gozxing.ResultPoint {

	corr := float64(distance(bottomLeft, bottomRight)) / float64(dimensionTop)
	norm := float64(distance(topLeft, topRight))
	cos := (topRight.GetX() - topLeft.GetX()) / norm
	sin := (topRight.GetY() - topLeft.GetY()) / norm

	c1 := gozxing.NewResultPoint(topRight.GetX()+corr*cos, topRight.GetY()+corr*sin)

	corr = float64(distance(bottomLeft, topLeft)) / float64(dimensionRight)
	norm = float64(distance(bottomRight, topRight))
	cos = (topRight.GetX() - bottomRight.GetX()) / norm
	sin = (topRight.GetY() - bottomRight.GetY()) / norm

	c2 := gozxing.NewResultPoint(topRight.GetX()+corr*cos, topRight.GetY()+corr*sin)

	if !this.isValid(c1) {
		if this.isValid(c2) {
			return c2
		}
		return nil
	}
	if !this.isValid(c2) {
		return c1
	}

	l1 := abs(dimensionTop-this.transitionsBetween(topLeft, c1).getTransitions()) +
		abs(dimensionRight-this.transitionsBetween(bottomRight, c1).getTransitions())
	l2 := abs(dimensionTop-this.transitionsBetween(topLeft, c2).getTransitions()) +
		abs(dimensionRight-this.transitionsBetween(bottomRight, c2).getTransitions())

	if l1 <= l2 {
		return c1
	}

	return c2
}

// correctTopRight Calculates the position of the white top right module
// using the output of the rectangle detector for a square matrix
func (this *Detector) correctTopRight(bottomLeft, bottomRight, topLeft, topRight gozxing.ResultPoint, dimension int) gozxing.ResultPoint {

	corr := float64(distance(bottomLeft, bottomRight)) / float64(dimension)
	norm := float64(distance(topLeft, topRight))
	cos := (topRight.GetX() - topLeft.GetX()) / norm
	sin := (topRight.GetY() - topLeft.GetY()) / norm

	c1 := gozxing.NewResultPoint(topRight.GetX()+corr*cos, topRight.GetY()+corr*sin)

	corr = float64(distance(bottomLeft, topLeft)) / float64(dimension)
	norm = float64(distance(bottomRight, topRight))
	cos = (topRight.GetX() - bottomRight.GetX()) / norm
	sin = (topRight.GetY() - bottomRight.GetY()) / norm

	c2 := gozxing.NewResultPoint(topRight.GetX()+corr*cos, topRight.GetY()+corr*sin)

	if !this.isValid(c1) {
		if this.isValid(c2) {
			return c2
		}
		return nil
	}
	if !this.isValid(c2) {
		return c1
	}

	l1 := abs(this.transitionsBetween(topLeft, c1).getTransitions() -
		this.transitionsBetween(bottomRight, c1).getTransitions())
	l2 := abs(this.transitionsBetween(topLeft, c2).getTransitions() -
		this.transitionsBetween(bottomRight, c2).getTransitions())

	if l1 <= l2 {
		return c1
	}
	return c2
}

func (this *Detector) isValid(p gozxing.ResultPoint) bool {
	return p.GetX() >= 0 && p.GetX() < float64(this.image.GetWidth()) &&
		p.GetY() > 0 && p.GetY() < float64(this.image.GetHeight())
}

func distance(a, b gozxing.ResultPoint) int {
	return util.MathUtils_Round(gozxing.ResultPoint_Distance(a, b))
}

// increment Increments the Integer associated with a key by one.
func increment(table map[gozxing.ResultPoint]int, key gozxing.ResultPoint) map[gozxing.ResultPoint]int {
	value, ok := table[key]
	if !ok {
		table[key] = 1
	} else {
		table[key] = value + 1
	}
	return table
}

func sampleGrid(image *gozxing.BitMatrix,
	topLeft, bottomLeft, bottomRight, topRight gozxing.ResultPoint,
	dimensionX, dimensionY int) (*gozxing.BitMatrix, error) {

	sampler := common.GridSampler_GetInstance()

	return sampler.SampleGrid(
		image,
		dimensionX,
		dimensionY,
		0.5,
		0.5,
		float64(dimensionX)-0.5,
		0.5,
		float64(dimensionX)-0.5,
		float64(dimensionY)-0.5,
		0.5,
		float64(dimensionY)-0.5,
		topLeft.GetX(),
		topLeft.GetY(),
		topRight.GetX(),
		topRight.GetY(),
		bottomRight.GetX(),
		bottomRight.GetY(),
		bottomLeft.GetX(),
		bottomLeft.GetY())
}

// transitionsBetween Counts the number of black/white transitions between two points,
// using something like Bresenham's algorithm.
func (this *Detector) transitionsBetween(from, to gozxing.ResultPoint) *ResultPointsAndTransitions {
	// See QR Code Detector, sizeOfBlackWhiteBlackRun()
	fromX := int(from.GetX())
	fromY := int(from.GetY())
	toX := int(to.GetX())
	toY := int(to.GetY())
	steep := abs(toY-fromY) > abs(toX-fromX)
	if steep {
		fromX, fromY = fromY, fromX
		toX, toY = toY, toX
	}

	dx := abs(toX - fromX)
	dy := abs(toY - fromY)
	error := -dx / 2
	ystep := 1
	if !(fromY < toY) {
		ystep = -1
	}
	xstep := 1
	if !(fromX < toX) {
		xstep = -1
	}
	transitions := 0
	var inBlack bool
	if steep {
		inBlack = this.image.Get(fromY, fromX)
	} else {
		inBlack = this.image.Get(fromX, fromY)
	}
	for x, y := fromX, fromY; x != toX; x += xstep {
		var isBlack bool
		if steep {
			isBlack = this.image.Get(y, x)
		} else {
			isBlack = this.image.Get(x, y)
		}
		if isBlack != inBlack {
			transitions++
			inBlack = isBlack
		}
		error += dy
		if error > 0 {
			if y == toY {
				break
			}
			y += ystep
			error -= dx
		}
	}
	return NewResultPointsAndTransitions(from, to, transitions)
}

// ResultPointsAndTransitions Simply encapsulates two points and a number of transitions between them.
type ResultPointsAndTransitions struct {
	from        gozxing.ResultPoint
	to          gozxing.ResultPoint
	transitions int
}

func NewResultPointsAndTransitions(from, to gozxing.ResultPoint, transitions int) *ResultPointsAndTransitions {
	return &ResultPointsAndTransitions{from, to, transitions}
}

func (this *ResultPointsAndTransitions) getFrom() gozxing.ResultPoint {
	return this.from
}

func (this *ResultPointsAndTransitions) getTo() gozxing.ResultPoint {
	return this.to
}

func (this *ResultPointsAndTransitions) getTransitions() int {
	return this.transitions
}

func (this *ResultPointsAndTransitions) String() string {
	return fmt.Sprintf("%v/%v/%v", this.from, this.to, this.transitions)
}

// ResultPointsAndTransitionsComparator Orders ResultPointsAndTransitions by number of transitions, ascending.
func ResultPointsAndTransitionsComparator(o1, o2 *ResultPointsAndTransitions) bool {
	return o1.getTransitions() < o2.getTransitions()
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
