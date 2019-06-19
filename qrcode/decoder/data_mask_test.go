package decoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func testMask(t testing.TB, mask DataMask, dimension int, condition func(int, int) bool) {
	t.Helper()
	bits, _ := gozxing.NewSquareBitMatrix(dimension)
	mask.UnmaskBitMatrix(bits, dimension)
	for i := 0; i < dimension; i++ {
		for j := 0; j < dimension; j++ {
			if condition(i, j) != bits.Get(j, i) {
				t.Fatalf("(%v, %v)", i, j)
			}
		}
	}
}

func testMaskAcrossDimensions(t testing.TB, reference int, condition func(int, int) bool) {
	t.Helper()
	mask := DataMaskValues[reference]
	for version := 1; version <= 40; version++ {
		dimension := 17 + 4*version
		testMask(t, mask, dimension, condition)
	}
}

func TestMask0(t *testing.T) {
	testMaskAcrossDimensions(t, 0, func(i, j int) bool {
		return (i+j)%2 == 0
	})
}

func TestMask1(t *testing.T) {
	testMaskAcrossDimensions(t, 1, func(i, j int) bool {
		return i%2 == 0
	})
}

func TestMask2(t *testing.T) {
	testMaskAcrossDimensions(t, 2, func(i, j int) bool {
		return j%3 == 0
	})
}

func TestMask3(t *testing.T) {
	testMaskAcrossDimensions(t, 3, func(i, j int) bool {
		return (i+j)%3 == 0
	})
}

func TestMask4(t *testing.T) {
	testMaskAcrossDimensions(t, 4, func(i, j int) bool {
		return (i/2+j/3)%2 == 0
	})
}

func TestMask5(t *testing.T) {
	testMaskAcrossDimensions(t, 5, func(i, j int) bool {
		return (i*j)%2+(i*j)%3 == 0
	})
}

func TestMask6(t *testing.T) {
	testMaskAcrossDimensions(t, 6, func(i, j int) bool {
		return ((i*j)%2+(i*j)%3)%2 == 0
	})
}

func TestMask7(t *testing.T) {
	testMaskAcrossDimensions(t, 7, func(i, j int) bool {
		return ((i+j)%2+(i*j)%3)%2 == 0
	})
}
