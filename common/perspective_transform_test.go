package common

import (
	"reflect"
	"testing"
)

func TestPerspectiveTransform_QuadrilateralToQuadrilateral(t *testing.T) {
	p := PerspectiveTransform_QuadrilateralToQuadrilateral(
		3, 2, 7, 5, 2, 5, 3, 7,
		2, 3, 5, 7, 3, 2, 7, 5)

	expect := PerspectiveTransform{
		-10.500000, -1.000000, 3.500000,
		-18.000000, 4.000000, 1.000000,
		-1.500000, 2.000000, -14.500000,
	}
	if *p != expect {
		t.Fatalf("result is %v, expect %v", *p, expect)
	}
}

func TestPerspectiveTransform_TransformPoints(t *testing.T) {
	p := PerspectiveTransform_SquareToQuadrilateral(3, 2, 4, 5, 3, 4, 5, 2)
	points := []float64{1, 2, 5, 2}
	expect := []float64{2.5, 3.5, 3.3, 3.5}
	p.TransformPoints(points)
	if !reflect.DeepEqual(points, expect) {
		t.Fatalf("points is %v, expect %v", points, expect)
	}
}

func TestPerspectiveTransform_TransformPointsXY(t *testing.T) {
	p := PerspectiveTransform_SquareToQuadrilateral(3, 2, 4, 5, 3, 4, 5, 2)
	xs := []float64{1, 2, 5, 2}
	ys := []float64{2, 2, 2, 4}
	expectx := []float64{2.500000, 3.000000, 3.300000, 2.600000}
	expecty := []float64{3.500000, 3.500000, 3.500000, 3.200000}
	p.TransformPointsXY(xs, ys)
	if !reflect.DeepEqual(xs, expectx) {
		t.Fatalf("xValues is %v, expect %v", xs, expectx)
	}
	if !reflect.DeepEqual(ys, expecty) {
		t.Fatalf("yValues is %v, expect %v", ys, expecty)
	}
}

func TestPerspectiveTransform_SquareToQuadrilateral(t *testing.T) {
	p := PerspectiveTransform_SquareToQuadrilateral(2, 3, 5, 7, 7, 5, 3, 2)
	expect := PerspectiveTransform{
		0.500000, 1.000000, 2.000000,
		0.500000, -1.000000, 3.000000,
		-0.500000, 0.000000, 1.000000,
	}
	if *p != expect {
		t.Fatalf("result is %v, expect %v", *p, expect)
	}

	p = PerspectiveTransform_SquareToQuadrilateral(1, 5, 2, 4, 3, 5, 2, 6)
	expect = PerspectiveTransform{
		1.000000, 1.000000, 1.000000,
		-1.000000, 1.000000, 5.000000,
		0.000000, 0.000000, 1.000000,
	}
	if *p != expect {
		t.Fatalf("result is %v, expect %v", *p, expect)
	}
}

func TestPerspectiveTransform_QuadrilateralToSquare(t *testing.T) {
	p := PerspectiveTransform_QuadrilateralToSquare(2, 3, 5, 7, 3, 2, 7, 5)
	expect := PerspectiveTransform{
		-1.000000, 2.500000, -5.500000,
		-2.000000, 1.500000, -0.500000,
		-3.500000, 3.500000, 0.000000,
	}
	if *p != expect {
		t.Fatalf("result is %v, expect %v", *p, expect)
	}
}
