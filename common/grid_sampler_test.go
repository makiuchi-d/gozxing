package common_test

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestGridSampler_GetSetInstance(t *testing.T) {
	dummySampler := testutil.DummyGridSampler{}

	common.GridSampler_SetGridSampler(dummySampler)

	if s := common.GridSampler_GetInstance(); s != dummySampler {
		t.Fatalf("sampler is not DummyGridSampler")
	}
}

func TestGridSampler_checkAndNudgePoints(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)
	var points []float64
	var e error

	points = []float64{-2, 0}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{11, 0}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{0, -2}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{0, 11}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{0, 0, -2, 0}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{0, 0, 11, 0}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{0, 0, 0, -2}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{0, 0, 0, 11}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e == nil {
		t.Fatalf("return must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("return must be NotFoundException, %v", reflect.TypeOf(e))
	}

	points = []float64{-1, -1, 10, 10, 0, 0, -1, -1, 10, 10}
	e = common.GridSampler_checkAndNudgePoints(image, points)
	if e != nil {
		t.Fatalf("return must not be error")
	}
}
