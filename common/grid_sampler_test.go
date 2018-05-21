package common

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

type DummyGridSampler struct{}
type Dummy struct{}

func (s DummyGridSampler) SampleGrid(image *BitMatrix, dimensionX, dimensionY int,
	p1ToX, p1ToY, p2ToX, p2ToY, p3ToX, p3ToY, p4ToX, p4ToY float64,
	p1FromX, p1FromY, p2FromX, p2FromY, p3FromX, p3FromY, p4FromX, p4FromY float64) (*BitMatrix, error) {
	return nil, nil
}

func (s DummyGridSampler) SampleGridFromTransform(image *BitMatrix,
	dimensionX, dimensionY int, transform *PerspectiveTransform) (*BitMatrix, error) {
	return nil, nil
}

func TestGridSampler_GetSetInstance(t *testing.T) {
	dummySampler := DummyGridSampler{}

	GridSampler_SetGridSampler(dummySampler)

	if s := GridSampler_GetInstance(); s != dummySampler {
		t.Fatalf("sampler is not DummyGridSampler")
	}
}

func TestGridSampler_checkAndNudgePoints(t *testing.T) {
	image, _ := NewBitMatrix(10, 10)
	var points []float64
	var e error

	points = []float64{-2, 0}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{11, 0}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{0, -2}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{0, 11}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{0, 0, -2, 0}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{0, 0, 11, 0}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{0, 0, 0, -2}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{0, 0, 0, 11}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException")
	}

	points = []float64{-1, -1, 10, 10, 0, 0, -1, -1, 10, 10}
	e = GridSampler_checkAndNudgePoints(image, points)
	if e != nil {
		t.Fatalf("return must not be error")
	}
}
