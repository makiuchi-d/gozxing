package detector

import (
	"reflect"
	"sort"
	"testing"
	"unsafe"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode/detector"
)

var qrstr = "" +
	"##############      ##  ##  ##############        ##############      ##  ##  ##############\n" +
	"##          ##          ##  ##          ##        ##          ##          ##  ##          ##\n" +
	"##  ######  ##  ##  ##      ##  ######  ##        ##  ######  ##  ##  ##      ##  ######  ##\n" +
	"##  ######  ##          ##  ##  ######  ##        ##  ######  ##          ##  ##  ######  ##\n" +
	"##  ######  ##    ##  ####  ##  ######  ##        ##  ######  ##    ##  ####  ##  ######  ##\n" +
	"##          ##    ######    ##          ##        ##          ##    ######    ##          ##\n" +
	"##############  ##  ##  ##  ##############        ##############  ##  ##  ##  ##############\n" +
	"                ##  ##                                            ##  ##                    \n" +
	"######  ##########  ##  ######      ##            ######  ##########  ##  ######      ##    \n" +
	"  ##  ##        ########  ##  ##      ####          ##  ##        ########  ##  ##      ####\n" +
	"##    ####  ##  ########  ######  ########        ##    ####  ##  ########  ######  ########\n" +
	"    ####  ##  ####    ######  ####    ##              ####  ##  ####    ######  ####    ##  \n" +
	"        ##  ##    ##  ##  ######                          ##  ##    ##  ##  ######          \n" +
	"                ##  ##      ####    ######                        ##  ##      ####    ######\n" +
	"##############  ##  ##  ##      ##  ######        ##############  ##  ##  ##      ##  ######\n" +
	"##          ##  ######      ######    ####        ##          ##  ######      ######    ####\n" +
	"##  ######  ##  ####    ##  ##        ####        ##  ######  ##  ####    ##  ##        ####\n" +
	"##  ######  ##    ######  ##  ##    ####          ##  ######  ##    ######  ##  ##    ####  \n" +
	"##  ######  ##  ########  ####  ##  ##  ##        ##  ######  ##  ########  ####  ##  ##  ##\n" +
	"##          ##  ##  ########    ##    ##          ##          ##  ##  ########    ##    ##  \n" +
	"##############  ########  ######      ####        ##############  ########  ######      ####\n" +
	"                                                                                            \n" +
	"                                                                                            \n" +
	"                                                                                            \n" +
	"                                                                                            \n" +
	"##############      ##  ##  ##############        ##############      ##  ##  ##############\n" +
	"##          ##          ##  ##          ##        ##          ##          ##  ##          ##\n" +
	"##  ######  ##  ##  ##      ##  ######  ##        ##  ######  ##  ##  ##      ##  ######  ##\n" +
	"##  ######  ##          ##  ##  ######  ##        ##  ######  ##          ##  ##  ######  ##\n" +
	"##  ######  ##    ##  ####  ##  ######  ##        ##  ######  ##    ##  ####  ##  ######  ##\n" +
	"##          ##    ######    ##          ##        ##          ##    ######    ##          ##\n" +
	"##############  ##  ##  ##  ##############        ##############  ##  ##  ##  ##############\n" +
	"                ##  ##                                            ##  ##                    \n" +
	"######  ##########  ##  ######      ##            ######  ##########  ##  ######      ##    \n" +
	"  ##  ##        ########  ##  ##      ####          ##  ##        ########  ##  ##      ####\n" +
	"##    ####  ##  ########  ######  ########        ##    ####  ##  ########  ######  ########\n" +
	"    ####  ##  ####    ######  ####    ##              ####  ##  ####    ######  ####    ##  \n" +
	"        ##  ##    ##  ##  ######                          ##  ##    ##  ##  ######          \n" +
	"                ##  ##      ####    ######                        ##  ##      ####    ######\n" +
	"##############  ##  ##  ##      ##  ######        ##############  ##  ##  ##      ##  ######\n" +
	"##          ##  ######      ######    ####        ##          ##  ######      ######    ####\n" +
	"##  ######  ##  ####    ##  ##        ####        ##  ######  ##  ####    ##  ##        ####\n" +
	"##  ######  ##    ######  ##  ##    ####          ##  ######  ##    ######  ##  ##    ####  \n" +
	"##  ######  ##  ########  ####  ##  ##  ##        ##  ######  ##  ########  ####  ##  ##  ##\n" +
	"##          ##  ##  ########    ##    ##          ##          ##  ##  ########    ##    ##  \n" +
	"##############  ########  ######      ####        ##############  ########  ######      ####\n"

func TestModuleSizeComparator(t *testing.T) {
	patterns := []*detector.FinderPattern{
		detector.NewFinderPattern1(1, 1, 5),
		detector.NewFinderPattern1(1, 1, 3),
		detector.NewFinderPattern1(1, 1, 1),
		detector.NewFinderPattern1(1, 1, 4),
		detector.NewFinderPattern1(1, 1, 2),
	}
	sort.Slice(patterns, ModuleSizeComparator(patterns))

	expects := []float64{5, 4, 3, 2, 1}
	for i := 0; i < 5; i++ {
		if ems := patterns[i].GetEstimatedModuleSize(); ems != expects[i] {
			t.Fatalf("patterns[%v].GetEstimatedModuleSize = %v, expected %v", i, ems, expects[i])
		}
	}
}

func injectPossibleCenters(finder *MultiFinderPatternFinder, patterns []*detector.FinderPattern) {
	v := reflect.ValueOf(finder).Elem().FieldByName("possibleCenters")
	p := (*[]*detector.FinderPattern)(unsafe.Pointer(v.UnsafeAddr()))
	*p = patterns
}

func compResultPoint(p gozxing.ResultPoint, x, y float64) bool {
	return p.GetX() == x && p.GetY() == y
}

func TestMultiFinderPatternFinder_selectMultipleBestPatterns(t *testing.T) {

	finder := &MultiFinderPatternFinder{&detector.FinderPatternFinder{}}
	injectPossibleCenters(finder, []*detector.FinderPattern{
		detector.NewFinderPattern1(0, 0, 1),
		detector.NewFinderPattern1(0, 0, 2),
	})

	_, e := finder.selectMultipleBestPatterns()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("selectMultipleBestPatterns must be NotFoundException, %T", e)
	}

	ptns := []*detector.FinderPattern{
		detector.NewFinderPattern1(10, 10, 10),
		detector.NewFinderPattern1(50, 50, 11),
		detector.NewFinderPattern1(10, 50, 12),
	}
	injectPossibleCenters(finder, ptns)

	r, e := finder.selectMultipleBestPatterns()
	if e != nil {
		t.Fatalf("selectMultipleBestPatterns returns error: %v", e)
	}
	if !reflect.DeepEqual(r, [][]*detector.FinderPattern{ptns}) {
		t.Fatalf("bestPatterns = %v, expect %v", r, ptns)
	}

	injectPossibleCenters(finder, []*detector.FinderPattern{
		// estimatedModuleSize missmatch.
		detector.NewFinderPattern1(100, 100, 5.59),
		detector.NewFinderPattern1(100, 200, 5.58),
		detector.NewFinderPattern1(200, 100, 15.57),
		detector.NewFinderPattern1(200, 200, 15.56),
	})
	_, e = finder.selectMultipleBestPatterns()
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("selectMultipleBestPatterns must be NotFoundException, %T", e)
	}

	injectPossibleCenters(finder, []*detector.FinderPattern{
		detector.NewFinderPattern1(100, 100, 5.59),
		detector.NewFinderPattern1(100, 130, 5.58),
		detector.NewFinderPattern1(130, 100, 5.57), // 0-1-2: failed on check the size
		detector.NewFinderPattern1(100, 200, 5.56), // 0-1-3: failed on difference of the edge lengths
		detector.NewFinderPattern1(197, 75, 5.55),  // 0-3-4: failed on angle at topleft
		detector.NewFinderPattern1(200, 100, 5.54), // 0-3-5: OK
		detector.NewFinderPattern1(170, 10, 5.53),  // 2-4-6: OK
	})
	// expected results: [(100,200),(100,100),(200,100)] and [(130,100),(197,75),(170,10)]
	r, e = finder.selectMultipleBestPatterns()
	if e != nil {
		t.Fatalf("selectMultipleBestPatterns returns error: %v", e)
	}
	if len(r) != 2 {
		t.Fatalf("patterns length = %v, expect 2", len(r))
	}
	bl1, tl1, tr1 := gozxing.ResultPoint_OrderBestPatterns(r[0][0], r[0][1], r[0][2])
	bl2, tl2, tr2 := gozxing.ResultPoint_OrderBestPatterns(r[1][0], r[1][1], r[1][2])
	if tl1.GetX() == 197 {
		bl1, tl1, tr1, bl2, tl2, tr2 = bl2, tl2, tr2, bl1, tl1, tr1
	}
	if p := bl1; !compResultPoint(p, 100, 200) {
		t.Fatalf("result[0] BottomLeft = (%v,%v), expect (100,200)", p.GetX(), p.GetY())
	}
	if p := tl1; !compResultPoint(p, 100, 100) {
		t.Fatalf("result[0] BottomLeft = (%v,%v), expect (100,100)", p.GetX(), p.GetY())
	}
	if p := tr1; !compResultPoint(p, 200, 100) {
		t.Fatalf("result[0] BottomLeft = (%v,%v), expect (200,100)", p.GetX(), p.GetY())
	}
	if p := bl2; !compResultPoint(p, 170, 10) {
		t.Fatalf("result[0] BottomLeft = (%v,%v), expect (170,10)", p.GetX(), p.GetY())
	}
	if p := tl2; !compResultPoint(p, 197, 75) {
		t.Fatalf("result[0] BottomLeft = (%v,%v), expect (197,75)", p.GetX(), p.GetY())
	}
	if p := tr2; !compResultPoint(p, 130, 100) {
		t.Fatalf("result[0] BottomLeft = (%v,%v), expect (130,100)", p.GetX(), p.GetY())
	}
}

func TestMultiFinderPatternFinder_FindMulti(t *testing.T) {
	img, _ := gozxing.NewBitMatrix(100, 100)
	finder := NewMultiFinderPatternFinder(img, nil)
	_, e := finder.FindMulti(nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("FindMulti must be NotFoundException, %T", e)
	}

	img, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")

	finder = NewMultiFinderPatternFinder(img, nil)

	hint := make(map[gozxing.DecodeHintType]interface{})
	hint[gozxing.DecodeHintType_TRY_HARDER] = true

	r, e := finder.FindMulti(hint)
	if e != nil {
		t.Fatalf("FindMulti returns error: %v", e)
	}

	testsPoints := []struct{ tlx, tly, blx, bly, trx, try float64 }{
		{3.5, 3.5, 3.5, 17.5, 17.5, 3.5},     // TopLeft
		{28.5, 3.5, 28.5, 17.5, 42.5, 3.5},   // TopRight
		{28.5, 3.5, 28.5, 17.5, 42.5, 3.5},   // BottomLeft
		{28.5, 28.5, 28.5, 42.5, 42.5, 28.5}, // BottomLeft
	}
FORTESTPOINTS:
	for _, test := range testsPoints {
		for _, p := range r {
			if compResultPoint(p.GetTopLeft(), test.tlx, test.tly) &&
				compResultPoint(p.GetBottomLeft(), test.blx, test.bly) &&
				compResultPoint(p.GetTopRight(), test.trx, test.try) {
				continue FORTESTPOINTS
			}
		}
		t.Fatalf("result must contain: {(%v,%v), (%v,%v), (%v,%v)}",
			test.tlx, test.tly, test.blx, test.bly, test.trx, test.try)
	}
}
