package common

import (
	"math"

	"github.com/makiuchi-d/gozxing"
)

type GridSampler interface {
	SampleGrid(image *gozxing.BitMatrix, dimensionX, dimensionY int,
		p1ToX, p1ToY, p2ToX, p2ToY, p3ToX, p3ToY, p4ToX, p4ToY float64,
		p1FromX, p1FromY, p2FromX, p2FromY, p3FromX, p3FromY, p4FromX, p4FromY float64) (*gozxing.BitMatrix, error)

	SampleGridWithTransform(image *gozxing.BitMatrix,
		dimensionX, dimensionY int, transform *PerspectiveTransform) (*gozxing.BitMatrix, error)
}

var gridSampler GridSampler = NewDefaultGridSampler()

func GridSampler_SetGridSampler(newGridSampler GridSampler) {
	gridSampler = newGridSampler
}

func GridSampler_GetInstance() GridSampler {
	return gridSampler
}

// GridSampler_checkAndNudgePoints validates and adjusts transformed coordinates.
// Points that are slightly outside the image (within a tolerance proportional
// to image size) are clamped to the nearest edge. This accommodates the small
// overshoot that perspective transforms commonly produce at image borders.
// Points that are NaN, Inf, or far outside the image indicate a degenerate
// transform and cause an immediate error.
func GridSampler_checkAndNudgePoints(image *gozxing.BitMatrix, points []float64) error {
	width := image.GetWidth()
	height := image.GetHeight()

	// Allow overshoot up to ~5% of the image dimension. A perspective
	// transform on a 200px image can legitimately overshoot by ~10px at
	// the edges; rejecting that kills valid decodes.
	maxOvershootX := width/20 + 1
	maxOvershootY := height/20 + 1

	for offset := 0; offset < len(points)-1; offset += 2 {
		px := points[offset]
		py := points[offset+1]

		if math.IsNaN(px) || math.IsNaN(py) || math.IsInf(px, 0) || math.IsInf(py, 0) {
			return gozxing.NewNotFoundException(
				"(w, h) = (%v, %v),  (x, y) = (%v, %v)", width, height, px, py)
		}

		x := int(px)
		y := int(py)

		// Reject points that are far outside the image — the transform
		// is too distorted for a valid decode.
		if x < -maxOvershootX || x >= width+maxOvershootX ||
			y < -maxOvershootY || y >= height+maxOvershootY {
			return gozxing.NewNotFoundException(
				"(w, h) = (%v, %v),  (x, y) = (%v, %v)", width, height, x, y)
		}

		// Clamp to image bounds. BitMatrix.Get returns false (white) for
		// out-of-bounds access, so moderate overshoot reads as white
		// rather than aborting the decode.
		if x < 0 {
			points[offset] = 0
		} else if x >= width {
			points[offset] = float64(width - 1)
		}
		if y < 0 {
			points[offset+1] = 0
		} else if y >= height {
			points[offset+1] = float64(height - 1)
		}
	}
	return nil
}
