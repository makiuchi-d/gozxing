package encoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestDefaultPlacement_module(t *testing.T) {
	cw := make([]byte, 8)
	cw[0] = 255
	cw[1] = 0xaa
	cw[2] = 255
	p := NewDefaultPlacement(cw, 8, 8)

	p.module(0, 0, 1, 1)
	if p.bits[0] != 1 {
		t.Fatalf("module(0,0) bits[0]=%v, expect 1", p.bits[0])
	}
	p.module(0, 1, 1, 2)
	if p.bits[1] != 0 {
		t.Fatalf("module(0,0) bits[1]=%v, expect 0", p.bits[1])
	}

	p.module(2, -2, 0, 1)
	if p.bits[8*2+8-2] != 1 {
		t.Fatalf("module(0,0) bits[1]=%v, expect 1", p.bits[1])
	}

	p.module(-2, 2, 2, 1)
	if p.bits[8*(8-2)+2] != 1 {
		t.Fatalf("module(0,0) bits[1]=%v, expect 1", p.bits[1])
	}
}

func testBits(t testing.TB, p *DefaultPlacement, x, y int, expect int8) {
	t.Helper()
	pos := p.numcols*y + x
	bits := p.getBits()
	if r := bits[pos]; r != expect {
		t.Fatalf("bits(%v,%v) = %v, expect %v", x, y, r, expect)
	}
}

func TestDefaultPlacement_utah(t *testing.T) {
	cw := make([]byte, 12)
	cw[0] = 0x55
	cw[1] = 0xaa
	cw[2] = 0x55
	cw[3] = 0xaa
	p := NewDefaultPlacement(cw, 10, 10)

	p.utah(4, 0, 0)
	testBits(t, p, 8, 0, 0)
	testBits(t, p, 9, 0, 1)
	testBits(t, p, 8, 1, 0)
	testBits(t, p, 9, 1, 1)
	testBits(t, p, 0, 3, 0)
	testBits(t, p, 8, 2, 1)
	testBits(t, p, 9, 2, 0)
	testBits(t, p, 0, 4, 1)

	p.utah(2, 2, 1)
	testBits(t, p, 0, 0, 1)
	testBits(t, p, 1, 0, 0)
	testBits(t, p, 0, 1, 1)
	testBits(t, p, 1, 1, 0)
	testBits(t, p, 2, 1, 1)
	testBits(t, p, 0, 2, 0)
	testBits(t, p, 1, 2, 1)
	testBits(t, p, 2, 2, 0)

	p.utah(0, 4, 2)
	testBits(t, p, 0, 8, 0)
	testBits(t, p, 1, 8, 1)
	testBits(t, p, 0, 9, 0)
	testBits(t, p, 1, 9, 1)
	testBits(t, p, 2, 9, 0)
	testBits(t, p, 2, 0, 1)
	testBits(t, p, 3, 0, 0)
	testBits(t, p, 4, 0, 1)

	p.utah(1, 7, 3)
	testBits(t, p, 3, 9, 1)
	testBits(t, p, 4, 9, 0)
	testBits(t, p, 5, 0, 1)
	testBits(t, p, 6, 0, 0)
	testBits(t, p, 7, 0, 1)
	testBits(t, p, 5, 1, 0)
	testBits(t, p, 6, 1, 1)
	testBits(t, p, 7, 1, 0)

	p.utah(3, 5, 4)
	testBits(t, p, 3, 1, 0)
	testBits(t, p, 4, 1, 0)
	testBits(t, p, 3, 2, 0)
	testBits(t, p, 4, 2, 0)
	testBits(t, p, 5, 2, 0)
	testBits(t, p, 3, 3, 0)
	testBits(t, p, 4, 3, 0)
	testBits(t, p, 5, 3, 0)
}

func TestDefaultPlacement_corner1(t *testing.T) {
	cw := make([]byte, 18)
	cw[7] = 0x55
	p := NewDefaultPlacement(cw, 12, 12)
	p.corner1(7)
	testBits(t, p, 0, 11, 0)
	testBits(t, p, 1, 11, 1)
	testBits(t, p, 2, 11, 0)
	testBits(t, p, 10, 0, 1)
	testBits(t, p, 11, 0, 0)
	testBits(t, p, 11, 1, 1)
	testBits(t, p, 11, 2, 0)
	testBits(t, p, 11, 3, 1)
}

func TestDefaultPlacement_corner2(t *testing.T) {
	cw := make([]byte, 24)
	cw[7] = 0x55
	p := NewDefaultPlacement(cw, 14, 14)
	p.corner2(7)
	testBits(t, p, 0, 11, 0)
	testBits(t, p, 0, 12, 1)
	testBits(t, p, 0, 13, 0)
	testBits(t, p, 10, 0, 1)
	testBits(t, p, 11, 0, 0)
	testBits(t, p, 12, 0, 1)
	testBits(t, p, 13, 0, 0)
	testBits(t, p, 13, 1, 1)
}

func TestDefaultPlacement_corner3(t *testing.T) {
	cw := make([]byte, 21)
	cw[0] = 0x55
	p := NewDefaultPlacement(cw, 28, 6)
	p.corner3(0)
	testBits(t, p, 0, 5, 0)
	testBits(t, p, 0, 4, 1)
	testBits(t, p, 0, 3, 0)
	testBits(t, p, 26, 0, 1)
	testBits(t, p, 27, 0, 0)
	testBits(t, p, 27, 1, 1)
	testBits(t, p, 27, 2, 0)
	testBits(t, p, 27, 3, 1)
}

func TestDefaultPlacement_corner4(t *testing.T) {
	cw := make([]byte, 12)
	cw[6] = 0x55
	p := NewDefaultPlacement(cw, 16, 6)
	p.corner4(6)
	testBits(t, p, 0, 5, 0)
	testBits(t, p, 15, 5, 1)
	testBits(t, p, 13, 0, 0)
	testBits(t, p, 14, 0, 1)
	testBits(t, p, 15, 0, 0)
	testBits(t, p, 13, 1, 1)
	testBits(t, p, 14, 1, 0)
	testBits(t, p, 15, 1, 1)
}

func testPlace(t testing.TB, p *DefaultPlacement, m *gozxing.BitMatrix) {
	t.Helper()
	col := p.getNumcols()
	row := p.getNumrows()

	if w, h := m.GetWidth(), m.GetHeight(); col != w || row != h {
		t.Fatalf("DefaultPlacement(%vx%v) numcols/numrows expect %vx%v", col, row, w, h)
	}

	p.Place()

	for j := 0; j < row; j++ {
		for i := 0; i < col; i++ {
			if a, b := p.GetBit(i, j), m.Get(i, j); a != b {
				t.Fatalf("DefaultPlacement(%vx%v) bits[%v,%v] =%v, expect %v", col, row, i, j, a, b)
			}
		}
	}
}

func TestDefaultPlacement_Place(t *testing.T) {
	cw := make([]byte, 24)
	for i := 0; i < len(cw); i++ {
		cw[i] = 0xaa
	}

	p := NewDefaultPlacement(cw[:12], 10, 10)
	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"##    ##  ##  ####  \n"+
		"##  ####    ##  ##  \n"+
		"  ##  ##  ####    ##\n"+
		"####    ##  ##  ####\n"+
		"  ##  ####    ##  ##\n"+
		"    ##  ##  ####    \n"+
		"  ####    ##  ##  ##\n"+
		"##  ##  ####    ##  \n"+
		"##    ##  ##  ####  \n"+
		"##  ####    ##    ##\n"+
		"", "##", "  ")
	testPlace(t, p, expect)

	// corner1
	p = NewDefaultPlacement(cw[:18], 12, 12)
	expect, _ = gozxing.ParseStringToBitMatrix(""+
		"##    ##  ##  ####    ##\n"+
		"##  ####    ##  ##  ##  \n"+
		"  ##  ##  ####    ##  ##\n"+
		"####    ##  ##  ####    \n"+
		"  ##  ####    ##  ##  ##\n"+
		"    ##  ##  ####    ##  \n"+
		"  ####    ##  ##  ####  \n"+
		"##  ##  ####    ##  ##  \n"+
		"##    ##  ##  ####    ##\n"+
		"##  ####    ##  ##  ####\n"+
		"  ##  ##  ####    ##  ##\n"+
		"##  ##  ##  ##  ####    \n", "##", "  ")
	testPlace(t, p, expect)

	// corner2
	p = NewDefaultPlacement(cw[:24], 14, 14)
	expect, _ = gozxing.ParseStringToBitMatrix(""+
		"##    ##  ##  ####    ##  ##\n"+
		"##  ####    ##  ##  ####    \n"+
		"  ##  ##  ####    ##  ##  ##\n"+
		"####    ##  ##  ####    ##  \n"+
		"  ##  ####    ##  ##  ####  \n"+
		"    ##  ##  ####    ##  ##  \n"+
		"  ####    ##  ##  ####    ##\n"+
		"##  ##  ####    ##  ##  ####\n"+
		"##    ##  ##  ####    ##  ##\n"+
		"##  ####    ##  ##  ####    \n"+
		"  ##  ##  ####    ##  ##  ##\n"+
		"####    ##  ##  ####    ##  \n"+
		"  ##  ####    ##  ##  ####  \n"+
		"##  ##  ##  ####    ##    ##\n"+
		"", "##", "  ")
	testPlace(t, p, expect)

	// corner3
	p = NewDefaultPlacement(cw[:21], 28, 6)
	expect, _ = gozxing.ParseStringToBitMatrix(""+
		"##    ##  ##  ####    ##  ##  ####    ##  ##  ####    ##\n"+
		"##  ####    ##  ##  ####    ##  ##  ####    ##  ##  ##  \n"+
		"  ##  ##  ####    ##  ##  ####    ##  ##  ####    ##  ##\n"+
		"####    ##  ##  ####    ##  ##  ####    ##  ##  ####    \n"+
		"  ##  ####    ##  ##  ####    ##  ##  ####    ##  ##  ##\n"+
		"##  ##  ##  ####    ##  ##  ####    ##  ##  ####    ##  \n"+
		"", "##", "  ")
	testPlace(t, p, expect)

	// corner4
	p = NewDefaultPlacement(cw[:12], 16, 6)
	expect, _ = gozxing.ParseStringToBitMatrix(""+
		"##    ##  ##  ####    ##  ##  ##\n"+
		"##  ####    ##  ##  ####    ##  \n"+
		"  ##  ##  ####    ##  ##  ####  \n"+
		"####    ##  ##  ####    ##  ##  \n"+
		"  ##  ####    ##  ##  ####    ##\n"+
		"##  ##  ##  ####    ##  ##  ##  \n"+
		"", "##", "  ")
	testPlace(t, p, expect)

}
