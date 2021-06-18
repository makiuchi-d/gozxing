package rss

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestRSS14Reader_addOrTally(t *testing.T) {
	reader := NewRSS14Reader().(*rss14Reader)
	pairs := []*Pair{}

	pairs = reader.addOrTally(pairs, nil)
	if l := len(pairs); l != 0 {
		t.Fatalf("addOrTally length = %v, wants 0", l)
	}

	p1 := NewPair(60, 5000, NewFinderPattern(2, []int{178, 247}, 178, 247, 154))
	p2 := NewPair(97250, 19572, NewFinderPattern(5, []int{178, 247}, 178, 247, 154))
	p3 := NewPair(60, 5000, NewFinderPattern(2, []int{178, 247}, 178, 247, 154))

	pairs = reader.addOrTally(pairs, p1)
	if l := len(pairs); l != 1 {
		t.Fatalf("addOrTally length = %v, wants 1", l)
	}
	if pairs[0] != p1 {
		t.Fatalf("addOrTally [0] = %v, wants %v", pairs[0], p1)
	}
	if c := pairs[0].GetCount(); c != 0 {
		t.Fatalf("addOrTally [0].GetCount = %v, wants 0", c)
	}

	pairs = reader.addOrTally(pairs, p2)
	if l := len(pairs); l != 2 {
		t.Fatalf("addOrTally length = %v, wants 2", l)
	}
	if pairs[0] != p1 {
		t.Fatalf("addOrTally [0] = %v, wants %v", pairs[0], p1)
	}
	if pairs[1] != p2 {
		t.Fatalf("addOrTally [1] = %v, wants %v", pairs[1], p2)
	}

	pairs = reader.addOrTally(pairs, p3)
	if l := len(pairs); l != 2 {
		t.Fatalf("addOrTally length = %v, wants 2", l)
	}
	if c := pairs[0].GetCount(); c != 1 {
		t.Fatalf("addOrTally [0].GetCount = %v, wants 1", c)
	}
}

func TestRSS14Reader_findFinderPattern(t *testing.T) {
	reader := NewRSS14Reader().(*rss14Reader)

	tests := []struct {
		row   string
		right bool
		wants []int
	}{
		{"00011000", false, nil},
		{"00100011100000001011", false, []int{6, 18}},
		{"00100000111111101000", true, []int{3, 17}},
	}

	for _, test := range tests {
		row := testutil.NewBitArrayFromString(test.row)
		r, e := reader.findFinderPattern(row, test.right)

		if test.wants == nil {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("findFinderPattern(%v,%v) must NotFoundException, %T", test.row, test.right, e)
			}
			continue
		}

		if !reflect.DeepEqual(r, test.wants) {
			t.Fatalf("findFinderPattern(%v,%v) = %v, wants %v", test.row, test.right, r, test.wants)
		}
	}
}

func TestRSS14Reader_parseFoundFinderPattern(t *testing.T) {
	reader := NewRSS14Reader().(*rss14Reader)

	tests := []struct {
		row   string
		right bool
		wants *FinderPattern
	}{
		{"0101111100000101", false, nil},
		{"00100011100000001011", false, NewFinderPattern(2, []int{3, 18}, 3, 18, 0)},
		{"0010000011111110100000000", true, NewFinderPattern(7, []int{2, 17}, 22, 7, 0)},
	}

	for _, test := range tests {
		row := testutil.NewBitArrayFromString(test.row)
		startEnd, e := reader.findFinderPattern(row, test.right)
		if e != nil {
			t.Fatalf("findFinderPattern(%v) error: %v", test.row, e)
		}

		r, e := reader.parseFoundFinderPattern(row, 0, test.right, startEnd)
		if test.wants == nil {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("parseFoundFinderPattern(%v) must NotFoundException: %T", test.row, e)
			}
			continue
		}

		if !reflect.DeepEqual(*r, *test.wants) {
			t.Fatalf("parseFoundFinderPattern(%v) = %v, wants %v", test.row, r, test.wants)
		}
	}
}

func TestRSS14Reader_adjustOddEvenCounts(t *testing.T) {
	reader := NewRSS14Reader().(*rss14Reader)

	tests := []struct {
		outside    bool
		oddCounts  []int
		oddErrors  []float64
		oddWants   []int
		evenCounts []int
		evenErrors []float64
		evenWants  []int
	}{
		{ // mismatch=-2
			false,
			[]int{2, 2, 2, 1}, []float64{0, 0, 0, 0}, nil,
			[]int{2, 2, 1, 1}, []float64{0, 0, 0, 0}, nil,
		},
		{ // inc&decOdd error: incOdd & decEven, mismatch=1 oddParityBad & !evenParityBad
			false,
			[]int{1, 1, 1, 1}, []float64{0, 0, 0, 0}, nil,
			[]int{3, 3, 3, 3}, []float64{0, 0, 0, 0}, nil,
		},
		{ // inc&decEven error: decOdd & incEven, mismatch=1 !oddParityBad & evenParityBad
			false,
			[]int{3, 3, 3, 4}, []float64{0, 0, 0, 0}, nil,
			[]int{1, 1, 1, 0}, []float64{0, 0, 0, 0}, nil,
		},
		{ // decOdd & incEven: mismatch=0
			true,
			[]int{3, 3, 3, 5}, []float64{0.3, 0.1, 0, -0.1}, []int{3, 3, 3, 4},
			[]int{1, 1, 0, 0}, []float64{0, -0.1, 0.2, 0.1}, []int{1, 1, 1, 0},
		},
		{ // incOdd & decEven: mismatch=0
			true,
			[]int{1, 1, 0, 0}, []float64{0, -0.1, 0.2, 0.1}, []int{1, 1, 1, 0},
			[]int{3, 3, 3, 5}, []float64{0.3, 0.1, 0, -0.1}, []int{3, 3, 3, 4},
		},
		{ //decOdd: mismatch=1 oddParityBad
			false,
			[]int{2, 2, 2, 2}, []float64{0.2, 0.4, -0.2, 0}, []int{2, 2, 1, 2},
			[]int{2, 2, 2, 2}, []float64{0, 0, 0, 0}, []int{2, 2, 2, 2},
		},
		{ //decEven: mismatch=1 evenParityBad
			true,
			[]int{2, 2, 2, 2}, []float64{0, 0, 0, 0}, []int{2, 2, 2, 2},
			[]int{2, 2, 2, 3}, []float64{0.2, 0.4, -0.2, 0}, []int{2, 2, 1, 3},
		},
		{ //decOdd: mismatch=-1 oddParityBad
			false,
			[]int{2, 2, 1, 1}, []float64{0.2, 0.4, -0.2, 0}, []int{2, 3, 1, 1},
			[]int{2, 2, 2, 2}, []float64{0, 0, 0, 0}, []int{2, 2, 2, 2},
		},
		{ //incEven: mismatch=-1 evenParityBad
			true,
			[]int{2, 2, 1, 1}, []float64{0, 0, 0, 0}, []int{2, 2, 1, 1},
			[]int{2, 2, 2, 3}, []float64{0.2, 0.4, -0.2, 0}, []int{2, 3, 2, 3},
		},
		{ // mismatch=0 bothBad, odd>even
			false,
			[]int{2, 2, 2, 2}, []float64{-0.1, 0, 0.1, 0.2}, []int{1, 2, 2, 2},
			[]int{2, 2, 2, 1}, []float64{0.2, 0.4, -0.2, 0}, []int{2, 3, 2, 1},
		},
		{ // mismatch=0 bothBad, even>odd
			false,
			[]int{2, 2, 1, 1}, []float64{-0.1, 0, 0.1, 0.2}, []int{2, 2, 1, 2},
			[]int{2, 2, 2, 3}, []float64{0.2, 0.4, -0.2, 0}, []int{2, 2, 1, 3},
		},
	}

	for i, test := range tests {
		reader.oddCounts = test.oddCounts
		reader.evenCounts = test.evenCounts
		reader.oddRoundingErrors = test.oddErrors
		reader.evenRoundingErrors = test.evenErrors
		numModules := 15
		if test.outside {
			numModules = 16
		}

		e := reader.adjustOddEvenCounts(test.outside, numModules)

		if test.oddWants == nil && test.evenWants == nil {
			if _, ok := e.(gozxing.NotFoundException); !ok {
				t.Fatalf("adjustOddEvenCounts[%v] must NotFoundException, %T", i, e)
			}
			continue
		}
		if e != nil {
			t.Fatalf("adjustOddEvenCounts[%v] error: %v", i, e)
		}
		if !reflect.DeepEqual(reader.oddCounts, test.oddWants) {
			t.Fatalf("adjustOddEvenCounts[%v] odd: %v, wants %v", i, reader.oddCounts, test.oddWants)
		}
		if !reflect.DeepEqual(reader.evenCounts, test.evenWants) {
			t.Fatalf("adjustOddEvenCounts[%v] even: %v, wants %v", i, reader.evenCounts, test.evenWants)
		}
	}
}

func TestRSS14Reader_decodeDataCharacter(t *testing.T) {
	reader := NewRSS14Reader().(*rss14Reader)

	// error in adjustOddEvenCounts (odd[1,1,0,0]even[9,4,0,0])
	row := testutil.NewBitArrayFromString("00001" + "000111111110010" + "1000000000100001")
	pattern := NewFinderPattern(0, []int{5, 20}, 5, 20, 0)
	_, e := reader.decodeDataCharacter(row, pattern, false)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeDataCharacter must NotFoundException: %T", e)
	}

	// outside: odd[1,2,2,1] even[1,5,1,3]
	// inside: odd[1,2,3,1] even[3,1,1,3]
	row = testutil.NewBitArrayFromString(
		"10100111110010111" + "000111111110010" + "1110100010011101")
	pattern = NewFinderPattern(0, []int{17, 32}, 17, 32, 0)

	// outside
	r, e := reader.decodeDataCharacter(row, pattern, true)
	if e != nil {
		t.Fatalf("decodeDataCharacter error: %v", e)
	}
	if v, wants := r.GetValue(), 2315; v != wants {
		t.Fatalf("decodeDataCharacter(outside) value = %v, wants %v", v, wants)
	}
	if c, wants := r.GetChecksumPortion(), 7852; c != wants {
		t.Fatalf("decodeDataCharacter(outside) checksumPotion = %v, wants %v", c, wants)
	}

	// inside
	r, e = reader.decodeDataCharacter(row, pattern, false)
	if e != nil {
		t.Fatalf("decodeDataCharacter error: %v", e)
	}
	if v, wants := r.GetValue(), 842; v != wants {
		t.Fatalf("decodeDataCharacter(inside) value = %v, wants %v", v, wants)
	}
	if c, wants := r.GetChecksumPortion(), 7831; c != wants {
		t.Fatalf("decodeDataCharacter(inside) checksumPotion = %v, wants %v", c, wants)
	}
}

func TestRSS14Reader_decodePair(t *testing.T) {
	reader := NewRSS14Reader().(*rss14Reader)

	// error on findFinderPattern
	row := testutil.NewBitArrayFromString("00011000")
	if p := reader.decodePair(row, false, 0, nil); p != nil {
		t.Fatalf("decodePair must be nil, %v", p)
	}

	// error on parseFoundFinderPattern
	row = testutil.NewBitArrayFromString("0101111100000101")
	if p := reader.decodePair(row, false, 0, nil); p != nil {
		t.Fatalf("decodePair must be nil, %v", p)
	}

	// error on decodeDataCharacter(outside)
	row = testutil.NewBitArrayFromString("" +
		"01011101111111111" + "000111111110010" + "1001100110011001")
	if p := reader.decodePair(row, false, 0, nil); p != nil {
		t.Fatalf("decodePair must be nil, %v", p)
	}

	// error on decodeDataCharacter(inside)
	row = testutil.NewBitArrayFromString(
		"10011001100110011" + "000111111110010" + "1000000000100001")
	if p := reader.decodePair(row, false, 0, nil); p != nil {
		t.Fatalf("decodePair must be nil, %v", p)
	}

	rps := make([]gozxing.ResultPoint, 0)
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK: gozxing.ResultPointCallback(
			func(rp gozxing.ResultPoint) {
				rps = append(rps, rp)
			}),
	}

	// JIS X 0509 Annex F.1
	row = testutil.NewBitArrayFromString("" +
		"01" +
		"0001010111000111" + // 1: 31111333
		"011100000000010" + // check8: 13911
		"111010001001110" + // 2: 31131231 (rev)
		"101101111001100" + // 4: 11214222
		"101111100000111" + // check1: 11553 (rev)
		"0010011101110111" + // 3: 21231313 (rev)
		"010")

	// left
	p := reader.decodePair(row, false, 0, hints)
	if p == nil {
		t.Fatalf("decodePair(left) = nil")
	}
	if v, wants := p.GetValue(), 2733309; v != wants {
		t.Fatalf("left value = %v, wants %v", v, wants)
	}
	if l, wants := len(rps), 1; l != wants {
		t.Fatalf("len(rps) = %v, wants %v", l, wants)
	}
	if x, y, wx, wy := rps[0].GetX(), rps[0].GetY(), 25.0, 0.0; x != wx || y != wy {
		t.Fatalf("rps[0] = (%v,%v), wants (%v,%v)", x, y, wx, wy)
	}

	// right
	row.Reverse()
	p = reader.decodePair(row, true, 0, hints)
	if p == nil {
		t.Fatalf("decodePair(right) = nil")
	}
	if v, wants := p.GetValue(), 1170097; v != wants {
		t.Fatalf("right value = %v, wants %v", v, wants)
	}
	if l, wants := len(rps), 2; l != wants {
		t.Fatalf("len(rps) = %v, wants %v", l, wants)
	}
	if x, y, wx, wy := rps[1].GetX(), rps[1].GetY(), 70.0, 0.0; x != wx || y != wy {
		t.Fatalf("rps[0] = (%v,%v), wants (%v,%v)", x, y, wx, wy)
	}
}

func TestCheckChecksum(t *testing.T) {
	left := NewPair(2733309, 40924, NewFinderPattern(8, []int{18, 33}, 18, 33, 0))
	right := NewPair(1170097, 32656, NewFinderPattern(1, []int{19, 34}, 77, 62, 0))

	if !checkChecksum(left, right) {
		t.Fatalf("checkChecksum must true")
	}

	left.checksumPortion++
	if checkChecksum(left, right) {
		t.Fatalf("checkChecksum must false")
	}

	left.checksumPortion--
	right.checksumPortion++
	if checkChecksum(left, right) {
		t.Fatalf("checkChecksum must false")
	}

	right.checksumPortion--
	left.finderPattern.value = 2
	if checkChecksum(left, right) {
		t.Fatalf("checkChecksum must false")
	}

	left.finderPattern.value = 8
	right.finderPattern.value = 2
	if checkChecksum(left, right) {
		t.Fatalf("checkChecksum must false")
	}
}

func TestConstructResult(t *testing.T) {
	// 0123456789050
	// checkdigit: 0+1+6+3+12+5+18+7+24+9+0+5+0 mod 10 = 90 mod 10 = 0
	// Vleft= 27210(39476), Vright=2923880(30676)
	// checksum=44+1, Cleft=5, Cright=0
	left := NewPair(27210, 39476, NewFinderPattern(5, []int{1, 3}, 13, 17, 10))
	right := NewPair(2923880, 30676, NewFinderPattern(0, []int{5, 7}, 31, 37, 20))

	result := constructResult(left, right)

	if f, wants := result.GetBarcodeFormat(), gozxing.BarcodeFormat_RSS_14; f != wants {
		t.Fatalf("result format = %v, wants %v", f, wants)
	}
	if txt, wants := result.GetText(), "01234567890500"; txt != wants {
		t.Fatalf("result text = \"%v\", wants \"%v\"", txt, wants)
	}
	expectPoints := []gozxing.ResultPoint{
		gozxing.NewResultPoint(13, 10),
		gozxing.NewResultPoint(17, 10),
		gozxing.NewResultPoint(31, 20),
		gozxing.NewResultPoint(37, 20),
	}
	if l, wants := len(result.GetResultPoints()), 4; l != wants {
		t.Fatalf("len(ResultPoints) = %v, wants %v", l, wants)
	}
CHECKRESULTPOINTS:
	for _, ep := range expectPoints {
		ex := ep.GetX()
		ey := ep.GetY()
		for _, rp := range result.GetResultPoints() {
			if x, y := rp.GetX(), rp.GetY(); x == ex && y == ey {
				continue CHECKRESULTPOINTS
			}
		}
		t.Fatalf("ResultPoint must contains %v", ep)
	}
}

func TestRSS14Reader_DecodeRow(t *testing.T) {
	// 0123456789050
	// left: 123456789050 / 4537077 = 27210
	// right: 123456789050 % 4537077 = 2923880
	// left-out(1) = 27210 / 1597 = 17
	//    group=1, odd=12,8, even=4,1, Teven=1
	//    Vodd = 17 / 1 = 17, Veven = 17 % 1 = 0
	//    odd:[1,3,3,5], even[1,1,1,1]
	// left-in (2) = 27210 / 1597 = 61
	//    group=1, odd=5,2 even=10,7 Todd=4
	//    Vodd=61%4=1, Veven=61/4=15,
	//    odd[1,1,2,1], even[1,3,3,3]
	// right-out(3) = 2923880 / 1597 = 1830
	//    group=3, Gsum=961 odd=8,4 even=8,5, Teven=34
	//    Vodd=(1830-961)/34=25, Veven=(1830-961)%34=19
	//    odd[3,2,1,2], even[2,2,1,3]
	// right-in (4) = 2923880 % 1597 = 1370
	//    group=3, Gsum=1036, odd=9,6 even=6,3 Todd=48
	//    Vodd=(1370-1036)%48=46 Veven=(1370-1036)/48=6
	//    odd[4,2,2,1] even[2,1,1,2]
	rowstr := "01" +
		"0100010001000001" + // Cout=17(6376)
		"001111100000010" + // pattern5
		"111011100111010" + //  Cin=61(8275)
		"111100110110100" + //  Cin=1370(5563)
		"101100000000111" + // pattern0
		"0001101001100111" + // Cout=1830(8424)
		"01"

	reader := NewRSS14Reader().(*rss14Reader)

	row := testutil.NewBitArrayFromString(rowstr[:50])
	_, e := reader.DecodeRow(0, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must NotFoundException, %T", e)
	}

	row = testutil.NewBitArrayFromString(rowstr)
	reader.DecodeRow(0, row, nil)
	reader.DecodeRow(0, row, nil)

	row = testutil.NewBitArrayFromString(rowstr[30:])
	r, e := reader.DecodeRow(2, row, nil)
	if e != nil {
		t.Fatalf("DecodeRow error: %v", e)
	}
	if txt, wants := r.GetText(), "01234567890500"; txt != wants {
		t.Fatalf("DecodeRow = \"%v\", wants \"%v\"", txt, wants)
	}

	reader.Reset()
	if left := reader.possibleLeftPairs; len(left) != 0 {
		t.Fatalf("Reset() failed: left=%v", left)
	}
	if right := reader.possibleRightPairs; len(right) != 0 {
		t.Fatalf("Reset() failed: right=%v", right)
	}

	_, e = reader.DecodeRow(0, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must NotFoundException, %T", e)
	}
}

func TestRSS14Reader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/rss14-[12]/
	reader := NewRSS14Reader()
	format := gozxing.BarcodeFormat_RSS_14
	harder := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}

	tests := []struct {
		file     string
		wants    string
		hints    map[gozxing.DecodeHintType]interface{}
		metadata map[gozxing.ResultMetadataType]interface{}
	}{
		// blackbox/rss14-1
		{
			"testdata/1_1.png", "04412345678909", nil,
			map[gozxing.ResultMetadataType]interface{}{
				gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER: "]e0",
			},
		},
		{"testdata/1_2.png", "00821935106427", nil, nil},
		{"testdata/1_3.png", "00075678164125", nil, nil},
		{"testdata/1_4.png", "20012345678909", nil, nil},
		{"testdata/1_5.png", "00034567890125", nil, nil},
		{"testdata/1_6.png", "00012345678905", nil, nil},

		// blackbox/rss14-2
		{"testdata/2_6.png", "02001234567893", harder, nil},
		{"testdata/2_7.png", "02001234567893", nil, nil},
		{"testdata/2_8.png", "02001234567893", nil, nil},
		{"testdata/2_13.png", "02001234567893", nil, nil},
		{"testdata/2_14.png", "02001234567893", harder, nil},
		{"testdata/2_20.png", "00012345678905", harder, nil},
		{"testdata/2_23.png", "00012345678905", harder, nil},
		{"testdata/2_24.png", "00012345678905", harder, nil},
		// original zxing could not read.
		// {"testdata/2_1.png", "04412345678909", harder, nil},
		// {"testdata/2_2.png", "04412345678909", harder, nil},
		// {"testdata/2_3.png", "04412345678909", harder, nil},
		// {"testdata/2_4.png", "04412345678909", harder, nil},
		// {"testdata/2_5.png", "02001234567893", harder, nil},
		// {"testdata/2_9.png", "02001234567893", harder, nil},
		// {"testdata/2_10.png", "02001234567893", harder, nil},
		// {"testdata/2_11.png", "02001234567893", harder, nil},
		// {"testdata/2_12.png", "02001234567893", harder, nil},
		// {"testdata/2_15.png", "02001234567893", harder, nil},
		// {"testdata/2_16.png", "02001234567893", harder, nil},
		// {"testdata/2_17.png", "02001234567893", harder, nil},
		// {"testdata/2_18.png", "02001234567893", harder, nil},
		// {"testdata/2_19.png", "02001234567893", harder, nil},
		// {"testdata/2_21.png", "00012345678905", harder, nil},
		// {"testdata/2_22.png", "00012345678905", harder, nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, test.hints, test.metadata)
		reader.Reset()
	}
}
