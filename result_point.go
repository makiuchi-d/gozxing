package gozxing

import (
	"github.com/makiuchi-d/gozxing/common/detector"
)

type ResultPoint struct {
	x float64
	y float64
}

func NewResultPoint(x, y float64) ResultPoint {
	return ResultPoint{x, y}
}

func (rp ResultPoint) GetX() float64 {
	return rp.x
}

func (rp ResultPoint) GetY() float64 {
	return rp.y
}

// Orders an array of three ResultPoints in an order [A,B,C] such that AB is less than AC
// and BC is less than AC, and the angle between BC and BA is less than 180 degrees.
// @param patterns array of three {@code ResultPoint} to order
func ResultPoint_OrderBestPatterns(patterns []ResultPoint) {
	zeroOneDistance := distance(patterns[0], patterns[1])
	oneTwoDistance := distance(patterns[1], patterns[2])
	zeroTwoDistance := distance(patterns[0], patterns[2])

	var pointA, pointB, pointC ResultPoint

	if oneTwoDistance >= zeroOneDistance && oneTwoDistance >= zeroTwoDistance {
		pointB = patterns[0]
		pointA = patterns[1]
		pointC = patterns[2]
	} else if zeroTwoDistance >= oneTwoDistance && zeroTwoDistance >= zeroOneDistance {
		pointB = patterns[1]
		pointA = patterns[0]
		pointC = patterns[2]
	} else {
		pointB = patterns[2]
		pointA = patterns[0]
		pointC = patterns[1]
	}

	if crossProductZ(pointA, pointB, pointC) < 0.0 {
		pointA, pointC = pointC, pointA
	}

	patterns[0] = pointA
	patterns[1] = pointB
	patterns[2] = pointC
}

func distance(pattern1, pattern2 ResultPoint) float64 {
	return detector.MathUtils_DistanceFloat(pattern1.x, pattern1.y, pattern2.x, pattern2.y)
}

func crossProductZ(pointA, pointB, pointC ResultPoint) float64 {
	bX := pointB.x
	bY := pointB.y
	return ((pointC.x - bX) * (pointA.y - bY)) - ((pointC.y - bY) * (pointA.x - bX))
}
