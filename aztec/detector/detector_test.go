package detector

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func writeFrame(img *gozxing.BitMatrix, x, y, w, h, l int) {
	img.SetRegion(x, y, w, l)
	img.SetRegion(x, y, l, h)
	img.SetRegion(x+w-l, y, l, h)
	img.SetRegion(x, y+h-l, w, l)
}

func TestDetector_Detect(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(60, 60)
	det := NewDetector(img)

	// bulls eye not found
	_, e := det.DetectNoMirror()
	if e == nil {
		t.Fatalf("Detect must be error")
	}

	// extract pattern error (reedsolomon error)
	img, _ = gozxing.ParseStringToBitMatrix(""+
		"                                  \n"+
		"                                  \n"+
		"      ######          ##  ##      \n"+
		"      ######################      \n"+
		"        ##              ##        \n"+
		"        ##  ##########  ##        \n"+
		"        ##  ##      ##  ##        \n"+
		"        ##  ##  ##  ##  ##        \n"+
		"        ##  ##      ##  ##        \n"+
		"        ##  ##########  ##        \n"+
		"        ##              ##        \n"+
		"        ####################      \n"+
		"                    ##            \n"+
		"                                  \n"+
		"                                  \n",
		"##", "  ")
	det = NewDetector(testutil.ExpandBitMatrix(img, 3))
	_, e = det.Detect(false)
	if e == nil {
		t.Fatalf("Detect must be error")
	}

	img, _ = gozxing.ParseStringToBitMatrix(""+
		"                                  \n"+
		"                                  \n"+
		"      ######    ########  ##      \n"+
		"      ######################      \n"+
		"      ####              ####      \n"+
		"      ####  ##########  ####      \n"+
		"        ##  ##      ##  ####      \n"+
		"        ##  ##  ##  ##  ##        \n"+
		"      ####  ##      ##  ##        \n"+
		"      ####  ##########  ####      \n"+
		"      ####              ####      \n"+
		"        ####################      \n"+
		"          ######  ######          \n"+
		"                                  \n"+
		"                                  \n",
		"##", "  ")
	det = NewDetector(testutil.ExpandBitMatrix(img, 3))
	_, e = det.Detect(false)
	if e == nil {
		t.Fatalf("Detect must be error")
	}

	// correct data (compact)
	img, _ = gozxing.ParseStringToBitMatrix(""+
		"    ##    ##  ####        ##  \n"+
		"  ######    ##  ######      ##\n"+
		"    ####        ##  ##  ##    \n"+
		"##########################    \n"+
		"####  ##              ##      \n"+
		"    ####  ##########  ##  ##  \n"+
		"  ##  ##  ##      ##  ##      \n"+
		"  ######  ##  ##  ##  ########\n"+
		"  ######  ##      ##  ##      \n"+
		"  ######  ##########  ####    \n"+
		"    ####              ######  \n"+
		"##    ####################  ##\n"+
		"##        ##    ##  ##        \n"+
		"####      ######  ##  ##    ##\n"+
		"########    ####  ####  ##  ##\n",
		"##", "  ")
	det = NewDetector(testutil.ExpandBitMatrix(img, 3))
	r, e := det.Detect(false)
	if e != nil {
		t.Fatalf("detect error: %v", e)
	}
	if l := r.GetNbLayers(); l != 1 {
		t.Fatalf("NbLayers = %v, wants 1", l)
	}
	if n := r.GetNbDatablocks(); n != 11 {
		t.Fatalf("NbDatablocks = %v, wants 11", n)
	}
	if c := r.IsCompact(); !c {
		t.Fatalf("IsCompact = %v, wants true", c)
	}
	if b := r.GetBits(); !reflect.DeepEqual(b, img) {
		t.Fatalf("detected img:\n%v\nwants:\n%v", b, img)
	}

	// mirrored
	det = NewDetector(testutil.ExpandBitMatrix(testutil.MirrorBitMatrix(img), 3))
	r, e = det.Detect(true)
	if e != nil {
		t.Fatalf("detect error: %v", e)
	}
	if l := r.GetNbLayers(); l != 1 {
		t.Fatalf("NbLayers = %v, wants 1", l)
	}
	if n := r.GetNbDatablocks(); n != 11 {
		t.Fatalf("NbDatablocks = %v, wants 11", n)
	}
	if c := r.IsCompact(); !c {
		t.Fatalf("IsCompact = %v, wants true", c)
	}
	if b := r.GetBits(); !reflect.DeepEqual(b, img) {
		t.Fatalf("detected img:\n%v\nwants:\n%v", b, img)
	}

	// full size
	img, _ = gozxing.ParseStringToBitMatrix(""+
		"          ####  ##    ##  ##    ######\n"+
		"      ####        ##    ##            \n"+
		"##  ####                        ####  \n"+
		"  ##################################  \n"+
		"####  ##                      ##    ##\n"+
		"    ####  ##################  ##    ##\n"+
		"##  ####  ##              ##  ####    \n"+
		"      ##  ##  ##########  ##  ##  ##  \n"+
		"    ####  ##  ##      ##  ##  ##  ####\n"+
		"  ##  ##  ##  ##  ##  ##  ##  ##  ##  \n"+
		"  ##  ##  ##  ##      ##  ##  ####    \n"+
		"##  ####  ##  ##########  ##  ######  \n"+
		"##    ##  ##              ##  ##  ####\n"+
		"  ##  ##  ##################  ####    \n"+
		"##  ####                      ##    ##\n"+
		"####  ############################    \n"+
		"####    ##          ####  ####        \n"+
		"        ####  ######    ####  ##      \n"+
		"    ####  ####              ##########\n",
		"##", "  ")
	det = NewDetector(testutil.ExpandBitMatrix(img, 3))
	r, e = det.Detect(false)
	if e != nil {
		t.Fatalf("detect error: %v", e)
	}
	if l := r.GetNbLayers(); l != 1 {
		t.Fatalf("NbLayers = %v, wants 1", l)
	}
	if n := r.GetNbDatablocks(); n != 10 {
		t.Fatalf("NbDatablocks = %v, wants 10", n)
	}
	if c := r.IsCompact(); c {
		t.Fatalf("IsCompact = %v, wants false", c)
	}
	if b := r.GetBits(); !reflect.DeepEqual(b, img) {
		t.Fatalf("detected img:\n%v\nwants:\n%v", b, img)
	}
}

func TestDetector_extractParameters(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(60, 60)
	det := NewDetector(img)

	e := det.extractParameters([]gozxing.ResultPoint{
		gozxing.NewResultPoint(0, 0), gozxing.NewResultPoint(0, 0),
		gozxing.NewResultPoint(0, 0), gozxing.NewResultPoint(80, 80),
	})
	if e == nil {
		t.Fatalf("extractParameters must be error")
	}

	// layer5
	writeFrame(img, 18, 18, 27, 27, 3)
	det.compact = true
	det.nbCenterLayers = 5
	points := []gozxing.ResultPoint{
		gozxing.NewResultPoint(46, 16),
		gozxing.NewResultPoint(46, 46),
		gozxing.NewResultPoint(16, 46),
		gozxing.NewResultPoint(16, 16),
	}

	// no orientation mark
	e = det.extractParameters(points)
	if e == nil {
		t.Fatalf("extractParameters must be error")
	}

	// orientation mark
	img.SetRegion(42, 15, 6, 6)
	img.SetRegion(42, 45, 6, 3)
	img.SetRegion(18, 45, 3, 3)
	// reedsolomon decode error
	img.SetRegion(45, 21, 3, 3)
	img.SetRegion(45, 24, 3, 3)
	img.SetRegion(45, 39, 3, 3)
	img.SetRegion(15, 36, 3, 3)
	e = det.extractParameters(points)
	if e == nil {
		t.Fatalf("extractParameters must be error")
	}

	// correct param
	img.SetRegion(45, 30, 3, 9)
	img.SetRegion(33, 45, 9, 3)
	img.SetRegion(21, 45, 6, 3)
	img.SetRegion(15, 33, 3, 9)
	img.SetRegion(15, 21, 3, 9)
	img.SetRegion(21, 15, 9, 3)
	img.SetRegion(36, 15, 6, 3)
	e = det.extractParameters(points)
	if e != nil {
		t.Fatalf("extractParameters error: %v", e)
	}
	if det.nbLayers != 3 {
		t.Fatalf("nbLayers = %v, wants 3", det.nbLayers)
	}
	if det.nbDataBlocks != 32 {
		t.Fatalf("nbDataBlocks = %v, wants 32", det.nbDataBlocks)
	}
	if det.shift != 0 {
		t.Fatalf("shift = %v, wants 0", det.shift)
	}

	// layer7
	writeFrame(img, 12, 12, 39, 39, 3)
	img.SetRegion(9, 9, 6, 6)
	img.SetRegion(51, 9, 3, 6)
	img.SetRegion(51, 48, 3, 3)
	img.SetRegion(18, 9, 3, 3)
	img.SetRegion(42, 9, 3, 3)
	img.SetRegion(51, 18, 3, 3)
	img.SetRegion(51, 24, 3, 3)
	img.SetRegion(51, 33, 3, 3)
	img.SetRegion(51, 45, 3, 3)
	img.SetRegion(42, 51, 6, 3)
	img.SetRegion(18, 51, 3, 3)
	img.SetRegion(9, 45, 3, 3)
	img.SetRegion(9, 33, 3, 3)
	img.SetRegion(9, 27, 3, 3)
	img.SetRegion(9, 18, 3, 3)
	det.compact = false
	det.nbCenterLayers = 7
	points = []gozxing.ResultPoint{
		gozxing.NewResultPoint(52, 10),
		gozxing.NewResultPoint(52, 52),
		gozxing.NewResultPoint(10, 52),
		gozxing.NewResultPoint(10, 10),
	}
	e = det.extractParameters(points)
	if e != nil {
		t.Fatalf("extractParameters error: %v", e)
	}
	if det.nbLayers != 9 {
		t.Fatalf("nbLayers = %v, wants 9", det.nbLayers)
	}
	if det.nbDataBlocks != 150 {
		t.Fatalf("nbDataBlocks = %v, wants 150", det.nbDataBlocks)
	}
	if det.shift != 3 {
		t.Fatalf("shift = %v, wants 3", det.shift)
	}

}

func TestDetector_getBullsEyeCorners(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(60, 60)
	det := NewDetector(img)
	_, e := det.getBullsEyeCorners(newPoint(25, 25))
	if e == nil {
		t.Fatalf("getBullsEyeCorners must be error")
	}

	img.SetRegion(30, 30, 3, 3)
	writeFrame(img, 24, 24, 15, 15, 3)
	writeFrame(img, 18, 18, 27, 27, 3)
	writeFrame(img, 12, 12, 39, 39, 3)

	// layer=7
	img.SetRegion(9, 9, 6, 6)
	img.SetRegion(51, 9, 3, 6)
	img.SetRegion(51, 48, 3, 3)
	ps, e := det.getBullsEyeCorners(newPoint(31, 31))
	if e != nil {
		t.Fatalf("getBullsEyeCorners error: %v", e)
	}
	if det.nbCenterLayers != 7 {
		t.Fatalf("nbCenterLayers = %v, wants 7", det.nbCenterLayers)
	}
	if det.compact {
		t.Fatalf("compact must false")
	}
	wants := []gozxing.ResultPoint{
		gozxing.NewResultPoint(52, 10),
		gozxing.NewResultPoint(52, 52),
		gozxing.NewResultPoint(10, 52),
		gozxing.NewResultPoint(10, 10),
	}
	if !reflect.DeepEqual(ps, wants) {
		t.Fatalf("corners: %v, wants %v", ps, wants)
	}

	// layer=5
	img.SetRegion(15, 15, 6, 6)
	img.SetRegion(45, 15, 3, 6)
	img.SetRegion(45, 42, 3, 3)
	ps, e = det.getBullsEyeCorners(newPoint(31, 31))
	if e != nil {
		t.Fatalf("getBullsEyeCorners error: %v", e)
	}
	if det.nbCenterLayers != 5 {
		t.Fatalf("nbCenterLayers = %v, wants 5", det.nbCenterLayers)
	}
	if !det.compact {
		t.Fatalf("compact must true")
	}
	wants = []gozxing.ResultPoint{
		gozxing.NewResultPoint(46, 16),
		gozxing.NewResultPoint(46, 46),
		gozxing.NewResultPoint(16, 46),
		gozxing.NewResultPoint(16, 16),
	}
	if !reflect.DeepEqual(ps, wants) {
		t.Fatalf("corners: %v, wants %v", ps, wants)
	}
}

func TestDetector_getMatrixCenter(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(60, 60)
	det := NewDetector(img)

	p := det.getMatrixCenter()
	ex, ey := 30, 30
	if p.x != ex || p.y != ey {
		t.Fatalf("getMatrixCenter = %v, wants <%v, %v>", p, ex, ey)
	}

	writeFrame(img, 10, 10, 30, 30, 3)
	p = det.getMatrixCenter()
	ex, ey = 25, 25
	if p.x != ex || p.y != ey {
		t.Fatalf("getMatrixCenter = %v, wants <%v, %v>", p, ex, ey)
	}
}

func TestDetector_sampleLine(t *testing.T) {
	img, _ := gozxing.ParseStringToBitMatrix("## ### ####  ##  ", "#", " ")
	e := 0b1101110111100110

	det := NewDetector(img)

	p1 := gozxing.NewResultPoint(0, 0)
	p2 := gozxing.NewResultPoint(16, 0)

	r := det.sampleLine(p1, p2, 16)

	if r != e {
		t.Fatalf("sampleLine = %v, wants %v", r, e)
	}
}

func TestDetector_isWhiteOrBlackRectangle(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(25, 25)
	img.SetRegion(0, 0, 15, 10)
	img.SetRegion(5, 10, 10, 5)
	det := NewDetector(img)

	tests := []struct {
		p1, p2, p3, p4 Point
		wants          bool
	}{
		//d a
		//c b
		{newPoint(20, 0), newPoint(20, 15), newPoint(0, 15), newPoint(0, 0), false}, // p4-p1: 0
		{newPoint(15, 0), newPoint(15, 20), newPoint(0, 20), newPoint(0, 0), false}, // p1-p2: 0
		{newPoint(15, 0), newPoint(15, 15), newPoint(0, 15), newPoint(0, 0), false}, // p2-p3: 0
		{newPoint(17, 0), newPoint(17, 15), newPoint(1, 15), newPoint(1, 0), false}, // p3-p4: 0
		{newPoint(12, 0), newPoint(12, 12), newPoint(0, 12), newPoint(0, 0), true},
	}

	for _, test := range tests {
		r := det.isWhiteOrBlackRectangle(test.p1, test.p2, test.p3, test.p4)
		if r != test.wants {
			t.Fatalf("isWhiteOrBlackRectangle(%v,%v,%v,%v) = %v, wants %v",
				test.p1, test.p2, test.p3, test.p4, r, test.wants)
		}
	}
}

func TestMinMax(t *testing.T) {
	a, b, e := 1, 2, 1
	if r := min(a, b); r != e {
		t.Fatalf("min(%v,%v) = %v, wants %v", a, b, r, e)
	}
	a, b, e = 2, 1, 1
	if r := min(a, b); r != e {
		t.Fatalf("min(%v,%v) = %v, wants %v", a, b, r, e)
	}
	a, b, e = 1, 2, 2
	if r := max(a, b); r != e {
		t.Fatalf("max(%v,%v) = %v, wants %v", a, b, r, e)
	}
	a, b, e = 2, 1, 2
	if r := max(a, b); r != e {
		t.Fatalf("max(%v,%v) = %v, wants %v", a, b, r, e)
	}
}

func TestDetector_getColor(t *testing.T) {
	img, _ := gozxing.ParseStringToBitMatrix(""+
		"################            \n"+
		"###############             \n"+
		"###############             \n"+
		"###############             \n"+
		"###########                 \n"+
		"###########                 \n"+
		"                            \n",
		"#", " ")
	det := NewDetector(img)

	tests := []struct {
		x1, y1, x2, y2 int
		wants          int
	}{
		{0, 0, 0, 0, 0},
		{0, 0, 13, 5, 1},
		{0, 0, 14, 5, 0},
		{15, 0, 26, 5, -1},
	}

	for _, test := range tests {
		p1 := newPoint(test.x1, test.y1)
		p2 := newPoint(test.x2, test.y2)
		r := det.getColor(p1, p2)
		if r != test.wants {
			t.Fatalf("getColor(%v,%v) = %v, wants %v", p1, p2, r, test.wants)
		}
	}
}

func TestDetector_getFirstDifferent(t *testing.T) {
	img, _ := gozxing.ParseStringToBitMatrix(""+
		"          \n"+
		" ##### #  \n"+
		" ##### #  \n"+
		" ##### #  \n"+
		" ##### #  \n"+
		" #######  \n"+
		"       #  \n"+
		" #######  \n"+
		"          \n"+
		"          \n",
		"#", " ")
	det := NewDetector(img)

	p := det.getFirstDifferent(newPoint(1, 1), true, 1, 1)
	ex, ey := 7, 7
	if p.x != ex || p.y != ey {
		t.Fatalf("getFirstDifferent(1,1) = {%v,%v}, wants {%v,%v}", p.x, p.y, ex, ey)
	}
	p = det.getFirstDifferent(newPoint(2, 2), true, 2, 2)
	ex, ey = 4, 4
	if p.x != ex || p.y != ey {
		t.Fatalf("getFirstDifferent(2,2) = {%v,%v}, wants {%v,%v}", p.x, p.y, ex, ey)
	}
}

func TestExpandSquare(t *testing.T) {
	ps := []gozxing.ResultPoint{
		gozxing.NewResultPoint(10, -10),
		gozxing.NewResultPoint(10, 10),
		gozxing.NewResultPoint(-10, 10),
		gozxing.NewResultPoint(-10, -10),
	}

	rps := expandSquare(ps, 2, 3)

	ratio := 1.5

	for i := range rps {
		rx, ry := rps[i].GetX(), rps[i].GetY()
		ex, ey := ps[i].GetX()*ratio, ps[i].GetY()*ratio

		if rx != ex || ry != ey {
			t.Fatalf("point[%v] = {%v, %v}, wants {%v, %v}", i, rx, ry, ex, ey)
		}
	}
}

func TestDetector_isValidPoint(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(10, 10)
	det := NewDetector(img)

	tests := []struct {
		x, y   float64
		expect bool
	}{
		{-1, 0, false},
		{0, -1, false},
		{0, 0, true},
		{9.4, 9.4, true},
		{9.5, 9.5, false},
	}

	for _, test := range tests {
		p := gozxing.NewResultPoint(test.x, test.y)
		r := det.isValidPoint(p)
		if r != test.expect {
			t.Fatalf("Detecotr.isValidPoint({%v,%v}) = %v, wants %v", test.x, test.y, r, test.expect)
		}
	}
}

func TestPoint(t *testing.T) {
	p := newPoint(3, 5)
	s := p.String()
	e := "<3 5>"
	if s != e {
		t.Fatalf("Point(3,5).String() = %q, expect %q", s, e)
	}
}
