package decoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

const (
	dm4str = "" +
		"##  ##  ##  ##  ##  ##  ##  ##  \n" +
		"##  ####  ####  ##  ####      ##\n" +
		"####      ####  ####            \n" +
		"######  ####  ##    ####      ##\n" +
		"####                ##########  \n" +
		"######      ##    ##        ####\n" +
		"####    ####          ##        \n" +
		"######  ####  ##    ##        ##\n" +
		"##    ##  ##############        \n" +
		"##  ####  ##    ######  ####  ##\n" +
		"######  ##    ##      ######    \n" +
		"##  ####  ######        ####  ##\n" +
		"##    ##  ######  ##  ######    \n" +
		"########    ##  ##        ##  ##\n" +
		"############  ##  ######    ##  \n" +
		"################################\n"
	dm5str = "" +
		"##  ##  ##  ##  ##  ##  ##  ##  ##  \n" +
		"##########      ####  ####        ##\n" +
		"##    ##  ##########    ####    ##  \n" +
		"######  ######        ####  ########\n" +
		"##    ##  ########  ##  ##  ####    \n" +
		"##  ######  ##  ##        ##########\n" +
		"######  ##    ##    ##  ####        \n" +
		"######        ##  ####            ##\n" +
		"####  ##  ####  ##    ##  ##        \n" +
		"####  ##              ##    ########\n" +
		"######    ##        ####            \n" +
		"##  ####  ######  ##########      ##\n" +
		"####  ######    ##      ####        \n" +
		"##  ##########        ##############\n" +
		"##############  ##          ##  ##  \n" +
		"########      ##    ######  ##    ##\n" +
		"##  ##  ##        ##  ########      \n" +
		"####################################\n"
	dm7str = "" +
		"##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  \n" +
		"##        ##    ####          ##  ##    ####\n" +
		"######                  ##########  ######  \n" +
		"######    ####        ####  ####  ##    ####\n" +
		"##  ##  ####  ########  ##  ############    \n" +
		"########  ##  ##  ##              ##########\n" +
		"######  ####        ##  ####  ##      ##    \n" +
		"##  ##    ####    ######  ##  ##  ##########\n" +
		"##        ####  ##      ##        ##  ##    \n" +
		"##  ############  ########    ##  ##  ######\n" +
		"########  ##  ##  ######################    \n" +
		"##    ##    ##  ######    ####            ##\n" +
		"##  ####    ####  ##  ######      ####  ##  \n" +
		"##    ########      ######        ####  ####\n" +
		"##########  ##    ##  ##                    \n" +
		"####          ##      ##      ##############\n" +
		"####  ##############  ####  ##  ##  ####    \n" +
		"####    ######  ##  ##  ##              ####\n" +
		"####      ##  ######    ####  ####    ##    \n" +
		"####    ##################    ######  ######\n" +
		"##      ##  ######  ##      ##        ####  \n" +
		"############################################\n"
	dm25str = "" +
		"##  ##  ##  ##  ##  ##  ##  ##  ##  \n" +
		"##  ####    ##    ##    ##  ####  ##\n" +
		"####      ####  ##    ##  ####      \n" +
		"##  ######    ##    ####  ####  ####\n" +
		"##    ######  ##################    \n" +
		"##  ####    ##    ########      ####\n" +
		"####  ######      ##########  ####  \n" +
		"####################################\n"
	dm26str = "" +
		"##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  ##  \n" +
		"##  ####    ##        ####  ######          ##  ##  ##  ##    ##\n" +
		"####      ####  ##############  ##  ##########                  \n" +
		"##  ######    ##            ########      ##        ######  ####\n" +
		"##    ######              ##    ####  ####    ##############    \n" +
		"######      ##    ##    ############      ####        ######  ##\n" +
		"############      ##    ##      ####  ####  ##    ########  ##  \n" +
		"################################################################\n"
)

func TestReadVersion(t *testing.T) {
	dm, _ := gozxing.NewBitMatrix(14, 14)
	v, e := readVersion(dm)
	if e != nil {
		t.Fatalf("readVersion returns error, %v", e)
	}
	if n := v.getVersionNumber(); n != 3 {
		t.Fatalf("readVersion verNum = %v, expect 3", n)
	}
}

func TestExtractDataRegion(t *testing.T) {
	// version mismatch
	dm, _ := gozxing.NewBitMatrix(10, 10)
	v, _ := getVersionForDimensions(14, 14)
	_, e := extractDataRegion(v, dm)
	if e == nil {
		t.Fatalf("extractDataRegion must be error")
	}

	dm, _ = gozxing.ParseStringToBitMatrix(dm5str, "##", "  ")
	v, _ = readVersion(dm)
	b, e := extractDataRegion(v, dm)
	if e != nil {
		t.Fatalf("extractDataRegion returns error, %v", e)
	}
	expect, _ := gozxing.ParseStringToBitMatrix(""+
		"########      ####  ####        \n"+
		"    ##  ##########    ####    ##\n"+
		"####  ######        ####  ######\n"+
		"    ##  ########  ##  ##  ####  \n"+
		"  ######  ##  ##        ########\n"+
		"####  ##    ##    ##  ####      \n"+
		"####        ##  ####            \n"+
		"##  ##  ####  ##    ##  ##      \n"+
		"##  ##              ##    ######\n"+
		"####    ##        ####          \n"+
		"  ####  ######  ##########      \n"+
		"##  ######    ##      ####      \n"+
		"  ##########        ############\n"+
		"############  ##          ##  ##\n"+
		"######      ##    ######  ##    \n"+
		"  ##  ##        ##  ########    \n", "##", "  ")
	if w, h, ew, eh := b.GetWidth(), b.GetHeight(), expect.GetWidth(), expect.GetHeight(); w != ew || h != eh {
		t.Fatalf("extractDataRegion size=%vx%v, expect %vx%v", w, h, ew, eh)
	}
	for j := 0; j < b.GetHeight(); j++ {
		for i := 0; i < b.GetWidth(); i++ {
			if bb, eb := b.Get(i, j), expect.Get(i, j); bb != eb {
				t.Fatalf("extractDataRegion [%v,%v] = %v, expect %v", i, j, bb, eb)
			}
		}
	}

	dm, _ = gozxing.ParseStringToBitMatrix(dm26str, "##", "  ")
	v, _ = readVersion(dm)
	b, e = extractDataRegion(v, dm)
	if e != nil {
		t.Fatalf("extractDataRegion returns error, %v", e)
	}
	expect, _ = gozxing.ParseStringToBitMatrix(""+
		"  ####    ##        ####  ##          ##  ##  ##  ##    \n"+
		"##      ####  ##############  ##########                \n"+
		"  ######    ##            ####      ##        ######  ##\n"+
		"    ######              ##  ##  ####    ##############  \n"+
		"####      ##    ##    ########      ####        ######  \n"+
		"##########      ##    ##    ##  ####  ##    ########  ##\n", "##", "  ")

	if w, h, ew, eh := b.GetWidth(), b.GetHeight(), expect.GetWidth(), expect.GetHeight(); w != ew || h != eh {
		t.Fatalf("extractDataRegion size=%vx%v, expect %vx%v", w, h, ew, eh)
	}
	for j := 0; j < b.GetHeight(); j++ {
		for i := 0; i < b.GetWidth(); i++ {
			if bb, eb := b.Get(i, j), expect.Get(i, j); bb != eb {
				t.Fatalf("extractDataRegion [%v,%v] = %v, expect %v", i, j, bb, eb)
			}
		}
	}
}

func TestReadCorner1(t *testing.T) {
	// datamatrix ver 3
	dm, _ := gozxing.NewBitMatrix(14, 14)
	parser, _ := NewBitMatrixParser(dm)
	if n := parser.GetVersion().getVersionNumber(); n != 3 {
		t.Fatalf("DataMatrix version = %v, expect 3", n)
	}
	parser.mappingBitMatrix.Set(10, 0)
	parser.mappingBitMatrix.Set(11, 0)
	parser.mappingBitMatrix.Set(11, 1)
	parser.mappingBitMatrix.Set(11, 2)
	parser.mappingBitMatrix.Set(11, 3)
	parser.mappingBitMatrix.Set(0, 11)
	parser.mappingBitMatrix.Set(1, 11)
	parser.mappingBitMatrix.Set(2, 11)

	if b := parser.readCorner1(12, 12); b != 255 {
		t.Fatalf("readCorner1 = %v, expect 255", b)
	}
}

func TestReadCorner2(t *testing.T) {
	// datamatrix ver 4
	dm, _ := gozxing.NewBitMatrix(16, 16)
	parser, _ := NewBitMatrixParser(dm)
	if n := parser.GetVersion().getVersionNumber(); n != 4 {
		t.Fatalf("DataMatrix version = %v, expect 4", n)
	}
	parser.mappingBitMatrix.Set(10, 0)
	parser.mappingBitMatrix.Set(11, 0)
	parser.mappingBitMatrix.Set(12, 0)
	parser.mappingBitMatrix.Set(13, 0)
	parser.mappingBitMatrix.Set(13, 1)
	parser.mappingBitMatrix.Set(0, 11)
	parser.mappingBitMatrix.Set(0, 12)
	parser.mappingBitMatrix.Set(0, 13)

	if b := parser.readCorner2(14, 14); b != 255 {
		t.Fatalf("readCorner2 = %v, expect 255", b)
	}
}

func TestReadCorner3(t *testing.T) {
	// datamatrix ver 25
	dm, _ := gozxing.NewBitMatrix(18, 8)
	parser, _ := NewBitMatrixParser(dm)
	if n := parser.GetVersion().getVersionNumber(); n != 25 {
		t.Fatalf("DataMatrix version = %v, expect 25", n)
	}
	parser.mappingBitMatrix.Set(13, 0)
	parser.mappingBitMatrix.Set(14, 0)
	parser.mappingBitMatrix.Set(15, 0)
	parser.mappingBitMatrix.Set(13, 1)
	parser.mappingBitMatrix.Set(14, 1)
	parser.mappingBitMatrix.Set(15, 1)
	parser.mappingBitMatrix.Set(0, 5)
	parser.mappingBitMatrix.Set(15, 5)

	if b := parser.readCorner3(6, 16); b != 255 {
		t.Fatalf("readCorner3 = %v, expect 255", b)
	}
}

func TestReadCorner4(t *testing.T) {
	// datamatrix ver 26
	dm, _ := gozxing.NewBitMatrix(32, 8)
	parser, _ := NewBitMatrixParser(dm)
	if n := parser.GetVersion().getVersionNumber(); n != 26 {
		t.Fatalf("DataMatrix version = %v, expect 26", n)
	}
	parser.mappingBitMatrix.Set(26, 0)
	parser.mappingBitMatrix.Set(27, 0)
	parser.mappingBitMatrix.Set(27, 1)
	parser.mappingBitMatrix.Set(27, 2)
	parser.mappingBitMatrix.Set(27, 3)
	parser.mappingBitMatrix.Set(0, 3)
	parser.mappingBitMatrix.Set(0, 4)
	parser.mappingBitMatrix.Set(0, 5)

	if b := parser.readCorner4(6, 28); b != 255 {
		t.Fatalf("readCorner4 = %v, expect 255", b)
	}
}

func TestReadCodewords(t *testing.T) {
	// with corner1
	dm, _ := gozxing.ParseStringToBitMatrix(dm7str, "##", "  ")
	expect := []byte{
		231, 54, 241, 136, 30, 181, 76, 226, 121, 15, 166, 61, 129, 192, 87, 237,
		133, 28, 178, 73, 223, 118, 14, 164, 59, 209, 104, 254, 150, 45, 185, 156,
		250, 26, 23, 222, 122, 70, 235, 125, 103, 178, 1, 9, 184, 168, 61, 131, 225, 13,
	}
	parser, _ := NewBitMatrixParser(dm)
	b, e := parser.readCodewords()
	if e != nil {
		t.Fatalf("readCodewords returns error, %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("readCodewords =\n%v,\nexpect\n%v", b, expect)
	}

	// with corner2
	dm, _ = gozxing.ParseStringToBitMatrix(dm4str, "##", "  ")
	expect = []byte{
		73, 102, 109, 109, 112, 33, 88, 112, 115, 109, 101, 129, 99, 7, 194, 15,
		202, 155, 119, 170, 200, 35, 246, 70,
	}
	parser, _ = NewBitMatrixParser(dm)
	b, e = parser.readCodewords()
	if e != nil {
		t.Fatalf("readCodewords returns error, %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("readCodewords =\n%v,\nexpect\n%v", b, expect)
	}

	// with corner3
	dm, _ = gozxing.ParseStringToBitMatrix(dm25str, "##", "  ")
	expect = []byte{
		98, 99, 100, 101, 102, 115, 244, 185, 139, 42, 95, 255,
	}
	parser, _ = NewBitMatrixParser(dm)
	b, e = parser.readCodewords()
	if e != nil {
		t.Fatalf("readCodewords returns error, %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("readCodewords =\n%v,\nexpect\n%v", b, expect)
	}

	// with corner4
	dm, _ = gozxing.ParseStringToBitMatrix(dm26str, "##", "  ")
	expect = []byte{
		98, 99, 100, 101, 102, 103, 129, 56, 206, 101, 202, 9, 172, 57, 10, 232,
		131, 157, 121, 70, 245,
	}
	parser, _ = NewBitMatrixParser(dm)
	b, e = parser.readCodewords()
	if e != nil {
		t.Fatalf("readCodewords returns error, %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("readCodewords =\n%v,\nexpect\n%v", b, expect)
	}

	// invalid mapping bitmatrix
	parser.mappingBitMatrix, _ = gozxing.NewBitMatrix(6, 6)
	_, e = parser.readCodewords()
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("readCodewords must be FormatException, %T", e)
	}
}

func TestNewBitMatrixParser(t *testing.T) {
	bm, _ := gozxing.NewBitMatrix(5, 5)
	_, e := NewBitMatrixParser(bm)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("NewBitMatrixParser must be FormatException, %T", e)
	}

	bm, _ = gozxing.NewBitMatrix(1, 8)
	_, e = NewBitMatrixParser(bm)
	if e == nil {
		t.Fatalf("NewBitMatrixParser must be error")
	}
}

func TestBitMatrixParser_readModule(t *testing.T) {
	bm, _ := gozxing.NewBitMatrix(10, 10)
	bm.Set(4, 3)
	p, _ := NewBitMatrixParser(bm)

	if r := p.readModule(2, 3, 8, 8); !r {
		t.Fatalf("readModule(2,3) should be true")
	}

	if r := p.readModule(-6, 3, 8, 8); !r {
		t.Fatalf("readModule(-6,3) should be true")
	}

	if r := p.readModule(2, -5, 8, 8); !r {
		t.Fatalf("readModule(2,-5) should be true")
	}

	if r := p.readModule(10, 3, 8, 8); !r {
		t.Fatalf("readModule(2,-5) should be true")
	}
}
