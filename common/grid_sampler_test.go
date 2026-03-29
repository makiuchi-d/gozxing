package common_test

import (
	"math"
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

func TestGridSampler_checkAndNudgePoints_NaN(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(10, 10)

	points := []float64{math.NaN(), 5}
	if err := common.GridSampler_checkAndNudgePoints(image, points); err == nil {
		t.Fatal("expected error for NaN x")
	}

	points = []float64{5, math.NaN()}
	if err := common.GridSampler_checkAndNudgePoints(image, points); err == nil {
		t.Fatal("expected error for NaN y")
	}

	points = []float64{math.Inf(1), 5}
	if err := common.GridSampler_checkAndNudgePoints(image, points); err == nil {
		t.Fatal("expected error for +Inf x")
	}

	points = []float64{5, math.Inf(-1)}
	if err := common.GridSampler_checkAndNudgePoints(image, points); err == nil {
		t.Fatal("expected error for -Inf y")
	}
}

func TestGridSampler_checkAndNudgePoints_Clamp(t *testing.T) {
	// Use a 200x200 image so the 5% tolerance is 11 pixels,
	// giving us room to test clamping of moderate overshoot.
	image, _ := gozxing.NewBitMatrix(200, 200)

	tests := []struct {
		name  string
		in    []float64
		wantX float64
		wantY float64
	}{
		{"negative x clamped to 0", []float64{-5, 100}, 0, 100},
		{"negative y clamped to 0", []float64{100, -5}, 100, 0},
		{"x past width clamped to width-1", []float64{204, 100}, 199, 100},
		{"y past height clamped to height-1", []float64{100, 204}, 100, 199},
		{"both clamped", []float64{-1, 205}, 0, 199},
		{"in bounds unchanged", []float64{100, 100}, 100, 100},
		{"at edge unchanged", []float64{0, 199}, 0, 199},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := make([]float64, len(tt.in))
			copy(points, tt.in)
			if err := common.GridSampler_checkAndNudgePoints(image, points); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if points[0] != tt.wantX {
				t.Errorf("x = %v, want %v", points[0], tt.wantX)
			}
			if points[1] != tt.wantY {
				t.Errorf("y = %v, want %v", points[1], tt.wantY)
			}
		})
	}
}

func TestGridSampler_checkAndNudgePoints_FarOutOfBounds(t *testing.T) {
	// Points far beyond the tolerance should still be rejected.
	image, _ := gozxing.NewBitMatrix(200, 200)

	// 5% of 200 = 10, +1 = 11 tolerance. 212 is outside that range.
	points := []float64{212, 100}
	if err := common.GridSampler_checkAndNudgePoints(image, points); err == nil {
		t.Fatal("expected error for x far beyond image width")
	}

	points = []float64{100, -15}
	if err := common.GridSampler_checkAndNudgePoints(image, points); err == nil {
		t.Fatal("expected error for y far below 0")
	}
}

func TestGridSampler_checkAndNudgePoints_AllPointsChecked(t *testing.T) {
	image, _ := gozxing.NewBitMatrix(200, 200)

	// All points should be checked, not just endpoints.
	// Interior point at x=204 (within tolerance) should be clamped.
	points := []float64{100, 100, 204, 100, 100, 100}
	if err := common.GridSampler_checkAndNudgePoints(image, points); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if points[2] != 199 {
		t.Errorf("interior point x = %v, want 199 (width-1)", points[2])
	}
}
