package util

import (
	"testing"
)

func testRound(t testing.TB, d float64, e int) {
	t.Helper()
	if r := MathUtils_Round(d); r != e {
		t.Fatalf("Round(%f) = %d, expect %d", d, r, e)
	}
}

func TestMathUtils_Round(t *testing.T) {
	testRound(t, 5.5, 6)
	testRound(t, 10.4, 10)
	testRound(t, -1.5, -2)
	testRound(t, -0.4, 0)
}

func testDistanceFloat(t testing.TB, aX, aY, bX, bY, e float64) {
	t.Helper()
	if r := MathUtils_DistanceFloat(aX, aY, bX, bY); r != e {
		t.Fatalf("Distance(%f, %f, %f, %f) = %f, expect %f", aX, aY, bX, bY, r, e)
	}
}

func TestMathUtils_DistanceFloat(t *testing.T) {
	// 3^2 + 4^2 = 5^2
	testDistanceFloat(t, 0, 0, 3, 4, 5)
	testDistanceFloat(t, 1, 0, -2, 4, 5)
	testDistanceFloat(t, 1, 1, -2, -3, 5)
}

func testDistanceInt(t testing.TB, aX, aY, bX, bY int, e float64) {
	t.Helper()
	if r := MathUtils_DistanceInt(aX, aY, bX, bY); r != e {
		t.Fatalf("Distance(%d, %d, %d, %d) = %f, expect %f", aX, aY, bX, bY, r, e)
	}
}

func TestMathUtils_DistanceInt(t *testing.T) {
	testDistanceInt(t, 0, 0, 3, 4, 5)
	testDistanceInt(t, 1, 0, -2, 4, 5)
	testDistanceInt(t, 1, 1, -2, -3, 5)
}

func testSum(t testing.TB, arr []int, e int) {
	t.Helper()
	if r := MathUtils_Sum(arr); r != e {
		t.Fatalf("Sum(%v) = %d, expect %d", arr, r, e)
	}
}

func TestMathUtils_Sum(t *testing.T) {
	testSum(t, []int{}, 0)
	testSum(t, []int{1, 2, 3, 4}, 10)
	testSum(t, []int{1, 2, -3}, 0)
}
