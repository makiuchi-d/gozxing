package encoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
)

func stringsToByteMatrix(str string) *ByteMatrix {
	bytes := make([][]int8, 1)
	bytes[0] = make([]int8, 0)
	y := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '\n' {
			bytes = append(bytes, make([]int8, 0))
			y++
			continue
		}
		i++
		var b int8
		switch str[i] {
		case '1':
			b = 1
		case '0':
			b = 0
		default:
			b = -1
		}
		bytes[y] = append(bytes[y], b)
	}

	return &ByteMatrix{bytes, len(bytes[0]), len(bytes)}
}

func TestEmbedPositionDetectionPatternsAndSeparators(t *testing.T) {
	var matrix *ByteMatrix
	var e error

	matrix = NewByteMatrix(17, 17)
	clearMatrix(matrix)
	matrix.SetBool(0, 7, true)
	e = embedPositionDetectionPatternsAndSeparators(matrix)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("embed must be WriterException, %T", e)
	}

	clearMatrix(matrix)
	matrix.SetBool(7, 0, true)
	e = embedPositionDetectionPatternsAndSeparators(matrix)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("embed must be WriterException, %T", e)
	}

	clearMatrix(matrix)
	e = embedPositionDetectionPatternsAndSeparators(matrix)
	if e != nil {
		t.Fatalf("embed returns error, %v", e)
	}
	expect := stringsToByteMatrix("" +
		" 1 1 1 1 1 1 1 0   0 1 1 1 1 1 1 1\n" +
		" 1 0 0 0 0 0 1 0   0 1 0 0 0 0 0 1\n" +
		" 1 0 1 1 1 0 1 0   0 1 0 1 1 1 0 1\n" +
		" 1 0 1 1 1 0 1 0   0 1 0 1 1 1 0 1\n" +
		" 1 0 1 1 1 0 1 0   0 1 0 1 1 1 0 1\n" +
		" 1 0 0 0 0 0 1 0   0 1 0 0 0 0 0 1\n" +
		" 1 1 1 1 1 1 1 0   0 1 1 1 1 1 1 1\n" +
		" 0 0 0 0 0 0 0 0   0 0 0 0 0 0 0 0\n" +
		"                                  \n" +
		" 0 0 0 0 0 0 0 0                  \n" +
		" 1 1 1 1 1 1 1 0                  \n" +
		" 1 0 0 0 0 0 1 0                  \n" +
		" 1 0 1 1 1 0 1 0                  \n" +
		" 1 0 1 1 1 0 1 0                  \n" +
		" 1 0 1 1 1 0 1 0                  \n" +
		" 1 0 0 0 0 0 1 0                  \n" +
		" 1 1 1 1 1 1 1 0                  ")
	if !reflect.DeepEqual(matrix.GetArray(), expect.GetArray()) {
		t.Fatalf("matrix bits:\n%vexpect:\n%v", matrix, expect)
	}
}

func TestEmbedBasicPatterns(t *testing.T) {
	var e error
	ver, _ := decoder.Version_GetVersionForNumber(1)
	d := ver.GetDimensionForVersion()
	matrix := NewByteMatrix(d, d)

	clearMatrix(matrix)
	matrix.SetBool(7, 0, true)
	e = embedBasicPatterns(ver, matrix)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("embed must be WriterException, %T", e)
	}

	clearMatrix(matrix)
	matrix.SetBool(8, 13, false)
	e = embedBasicPatterns(ver, matrix)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("embed must be WriterException, %T", e)
	}

	clearMatrix(matrix)
	e = embedBasicPatterns(ver, matrix)
	if e != nil {
		t.Fatalf("embed returns error, %v", e)
	}
	expect := stringsToByteMatrix("" +
		" 1 1 1 1 1 1 1 0           0 1 1 1 1 1 1 1\n" +
		" 1 0 0 0 0 0 1 0           0 1 0 0 0 0 0 1\n" +
		" 1 0 1 1 1 0 1 0           0 1 0 1 1 1 0 1\n" +
		" 1 0 1 1 1 0 1 0           0 1 0 1 1 1 0 1\n" +
		" 1 0 1 1 1 0 1 0           0 1 0 1 1 1 0 1\n" +
		" 1 0 0 0 0 0 1 0           0 1 0 0 0 0 0 1\n" +
		" 1 1 1 1 1 1 1 0 1 0 1 0 1 0 1 1 1 1 1 1 1\n" +
		" 0 0 0 0 0 0 0 0           0 0 0 0 0 0 0 0\n" +
		"             1                            \n" +
		"             0                            \n" +
		"             1                            \n" +
		"             0                            \n" +
		"             1                            \n" +
		" 0 0 0 0 0 0 0 0 1                        \n" +
		" 1 1 1 1 1 1 1 0                          \n" +
		" 1 0 0 0 0 0 1 0                          \n" +
		" 1 0 1 1 1 0 1 0                          \n" +
		" 1 0 1 1 1 0 1 0                          \n" +
		" 1 0 1 1 1 0 1 0                          \n" +
		" 1 0 0 0 0 0 1 0                          \n" +
		" 1 1 1 1 1 1 1 0                          ")
	if !reflect.DeepEqual(matrix.GetArray(), expect.GetArray()) {
		t.Fatalf("matrix bits:\n%vexpect:\n%v", matrix, expect)
	}

	ver, _ = decoder.Version_GetVersionForNumber(7)
	d = ver.GetDimensionForVersion()
	matrix = NewByteMatrix(d, d)
	clearMatrix(matrix)
	e = embedBasicPatterns(ver, matrix)
	if e != nil {
		t.Fatalf("embed returns error, %v", e)
	}
	expect = stringsToByteMatrix("" +
		" 1 1 1 1 1 1 1 0                                                           0 1 1 1 1 1 1 1\n" +
		" 1 0 0 0 0 0 1 0                                                           0 1 0 0 0 0 0 1\n" +
		" 1 0 1 1 1 0 1 0                                                           0 1 0 1 1 1 0 1\n" +
		" 1 0 1 1 1 0 1 0                                                           0 1 0 1 1 1 0 1\n" +
		" 1 0 1 1 1 0 1 0                         1 1 1 1 1                         0 1 0 1 1 1 0 1\n" +
		" 1 0 0 0 0 0 1 0                         1 0 0 0 1                         0 1 0 0 0 0 0 1\n" +
		" 1 1 1 1 1 1 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 1 1 1 1 1 1\n" +
		" 0 0 0 0 0 0 0 0                         1 0 0 0 1                         0 0 0 0 0 0 0 0\n" +
		"             1                           1 1 1 1 1                                        \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"         1 1 1 1 1                       1 1 1 1 1                       1 1 1 1 1        \n" +
		"         1 0 0 0 1                       1 0 0 0 1                       1 0 0 0 1        \n" +
		"         1 0 1 0 1                       1 0 1 0 1                       1 0 1 0 1        \n" +
		"         1 0 0 0 1                       1 0 0 0 1                       1 0 0 0 1        \n" +
		"         1 1 1 1 1                       1 1 1 1 1                       1 1 1 1 1        \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                                                                            \n" +
		"             0                                                                            \n" +
		"             1                           1 1 1 1 1                       1 1 1 1 1        \n" +
		" 0 0 0 0 0 0 0 0 1                       1 0 0 0 1                       1 0 0 0 1        \n" +
		" 1 1 1 1 1 1 1 0                         1 0 1 0 1                       1 0 1 0 1        \n" +
		" 1 0 0 0 0 0 1 0                         1 0 0 0 1                       1 0 0 0 1        \n" +
		" 1 0 1 1 1 0 1 0                         1 1 1 1 1                       1 1 1 1 1        \n" +
		" 1 0 1 1 1 0 1 0                                                                          \n" +
		" 1 0 1 1 1 0 1 0                                                                          \n" +
		" 1 0 0 0 0 0 1 0                                                                          \n" +
		" 1 1 1 1 1 1 1 0                                                                          ")
	if !reflect.DeepEqual(matrix.GetArray(), expect.GetArray()) {
		t.Fatalf("matrix bits:\n%vexpect:\n%v", matrix, expect)
	}
}

func TestFindMsbSet(t *testing.T) {
	if r := findMSBSet(0); r != 0 {
		t.Fatalf("findMsbSet(0) = %d, expect 0", r)
	}
	if r := findMSBSet(1); r != 1 {
		t.Fatalf("findMsbSet(1) = %d, expect 1", r)
	}
	if r := findMSBSet(255); r != 8 {
		t.Fatalf("findMsbSet(255) = %d, expect 8", r)
	}
}

func TestCalculateBCHCode(t *testing.T) {
	if _, e := calculateBCHCode(0, 0); e == nil {
		t.Fatalf("calculateBCHCode(0,0) must be error")
	}

	c, e := calculateBCHCode(7, matrixUtil_VERSION_INFO_POLY)
	if e != nil {
		t.Fatalf("calculateBCHCode returns error, %v", e)
	}
	if c != 0xc94 {
		t.Fatalf("calculateBCHCode(7) = 0x%x, expect 0xc94", c)
	}

}

func TestMakeTypeInfoBits(t *testing.T) {
	var e error
	ecl := decoder.ErrorCorrectionLevel_L
	bits := gozxing.NewEmptyBitArray()

	// invalid mask pattern
	e = makeTypeInfoBits(ecl, 8, bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("makeTypeInfoBits must be WriterException, %T", e)
	}

	// non empty bitarray
	bits.AppendBit(true)
	e = makeTypeInfoBits(ecl, 7, bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("makeTypeInfoBits must be WriterException, %T", e)
	}

	bits = gozxing.NewEmptyBitArray()
	e = makeTypeInfoBits(ecl, 5, bits)
	if e != nil {
		t.Fatalf("makeTypeInfoBits returns error, %v", e)
	}
	// typeInfo(L): 01
	// maskPattern(5): 101
	// bchCode(poly(011010000000000) mod poly(10100110111)) = 1100001010
	// 011011100001010 xor 101010000010010 = 110001100011000
	expect := []bool{
		true, true, false, false, false, true, true, false,
		false, false, true, true, false, false, false,
	}
	for i, b := range expect {
		if r := bits.Get(i); r != b {
			t.Fatalf("TypeInfo bits[%v] = %v, expect %v", i, r, b)
		}
	}
}

func TestEmbedTypeInfo(t *testing.T) {
	var e error
	ecl := decoder.ErrorCorrectionLevel_H
	matrix := NewByteMatrix(21, 21)
	clearMatrix(matrix)

	e = embedTypeInfo(ecl, 9, matrix)
	if e == nil {
		t.Fatalf("embedTypeInfo must be error")
	}

	clearMatrix(matrix)
	e = embedTypeInfo(ecl, 5, matrix)
	if e != nil {
		t.Fatalf("embedTypeInfo returns error, %v", e)
	}
	// typeInfo : 000001001010101
	expects := [][]int{
		{0, 8, 8, 20, 0},
		{1, 8, 8, 19, 0},
		{2, 8, 8, 18, 0},
		{3, 8, 8, 17, 0},
		{4, 8, 8, 16, 0},
		{5, 8, 8, 15, 1},
		{7, 8, 8, 14, 0},
		{8, 8, 13, 8, 0},
		{8, 7, 14, 8, 1},
		{8, 5, 15, 8, 0},
		{8, 4, 16, 8, 1},
		{8, 3, 17, 8, 0},
		{8, 2, 18, 8, 1},
		{8, 1, 19, 8, 0},
		{8, 0, 20, 8, 1},
	}
	for _, exp := range expects {
		if r := matrix.Get(exp[0], exp[1]); int(r) != exp[4] {
			t.Fatalf("embeded typeInfo matrix(%v,%v) = %v, expect %v", exp[0], exp[1], r, exp[4])
		}
		if r := matrix.Get(exp[2], exp[3]); int(r) != exp[4] {
			t.Fatalf("embeded typeInfo matrix(%v,%v) = %v, expect %v", exp[2], exp[3], r, exp[4])
		}
	}
}

func TestMakeVersionInfoBits(t *testing.T) {
	ver, _ := decoder.Version_GetVersionForNumber(10)

	bits := gozxing.NewEmptyBitArray()
	bits.AppendBit(true)
	e := makeVersionInfoBits(ver, bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("makeVersionInfoBits must be WriterException, %T", e)
	}

	bits = gozxing.NewEmptyBitArray()
	e = makeVersionInfoBits(ver, bits)
	if e != nil {
		t.Fatalf("makeVersionInfoBits returns error, %v", e)
	}
	// ver: 001010
	// poly: 1111100100101
	// bch: 010011010011
	// versionInfo: 001010 010011 010011
	expects := []bool{
		false, false, true, false, true, false,
		false, true, false, false, true, true,
		false, true, false, false, true, true,
	}
	for i, expect := range expects {
		if r := bits.Get(i); r != expect {
			t.Fatalf("makeVersionInfoBits [%v] = %v, expect %v", i, r, expect)
		}
	}
}

func TestEmbedVersionInfo(t *testing.T) {
	ver, _ := decoder.Version_GetVersionForNumber(6)
	d := ver.GetDimensionForVersion()
	matrix := NewByteMatrix(d, d)
	clearMatrix(matrix)

	e := maybeEmbedVersionInfo(ver, matrix)
	if e != nil {
		t.Fatalf("maybeEmbedVersionInfo(ver=6) returns error, %v", e)
	}
	for y := 0; y < d; y++ {
		for x := 0; x < d; x++ {
			if !isEmpty(matrix.Get(x, y)) {
				t.Fatalf("matrix (%v,%v) must be empty", x, y)
			}
		}
	}

	ver, _ = decoder.Version_GetVersionForNumber(10)
	d = ver.GetDimensionForVersion()
	matrix = NewByteMatrix(d, d)
	e = maybeEmbedVersionInfo(ver, matrix)
	if e != nil {
		t.Fatalf("maybeEmbedVersionInfo(ver=10) returns error, %v", e)
	}
	expectBits := []int8{0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1}
	expectPos := [][]int{
		{5, 48}, {5, 47}, {5, 46}, {4, 48}, {4, 47}, {4, 46},
		{3, 48}, {3, 47}, {3, 46}, {2, 48}, {2, 47}, {2, 46},
		{1, 48}, {1, 47}, {1, 46}, {0, 48}, {0, 47}, {0, 46},
	}
	for i, p := range expectPos {
		if r := matrix.Get(p[0], p[1]); r != expectBits[i] {
			t.Fatalf("matrix (%v,%v) = %v, expect %v", p[0], p[1], r, expectBits[i])
		}
		if r := matrix.Get(p[1], p[0]); r != expectBits[i] {
			t.Fatalf("matrix (%v,%v) = %v, expect %v", p[1], p[0], r, expectBits[i])
		}
	}
}

func TestEmbedDataBits(t *testing.T) {
	var e error
	ver, _ := decoder.Version_GetVersionForNumber(7)
	d := ver.GetDimensionForVersion()
	ecLevel := decoder.ErrorCorrectionLevel_M
	matrix := NewByteMatrix(d, d)

	// invalid maskPattern
	clearMatrix(matrix)
	embedBasicPatterns(ver, matrix)
	embedTypeInfo(ecLevel, 0, matrix)
	maybeEmbedVersionInfo(ver, matrix)
	e = embedDataBits(gozxing.NewEmptyBitArray(), 8, matrix)
	if e == nil {
		t.Fatalf("embedDataBits must be error")
	}

	clearMatrix(matrix)
	embedBasicPatterns(ver, matrix)
	embedTypeInfo(ecLevel, 0, matrix)
	maybeEmbedVersionInfo(ver, matrix)
	e = embedDataBits(gozxing.NewBitArray((124+72+1)*8), 0, matrix)
	if e == nil {
		t.Fatalf("embedDataBits must be error")
	}

	// valid
	dataBits := gozxing.NewEmptyBitArray()
	for i := 0; i < (124 + 72 - 1); i++ {
		dataBits.AppendBits(i, 8)
	}
	// the 125th byte position
	pos := [][]int{
		{18, 30}, {17, 30},
		{18, 31}, {17, 31},
		{18, 32}, {17, 32},
		{18, 33}, {17, 33},
	}

	// without mask
	clearMatrix(matrix)
	embedBasicPatterns(ver, matrix)
	embedTypeInfo(ecLevel, 0, matrix)
	maybeEmbedVersionInfo(ver, matrix)
	e = embedDataBits(dataBits, -1, matrix)
	if e != nil {
		t.Fatalf("embedDataBits returns error, %v", e)
	}
	v := 0
	for _, p := range pos {
		v = (v * 2) + int(matrix.Get(p[0], p[1]))
	}
	if v != 125 {
		t.Fatalf("embedDataBits 125th byte = %v, expect 125", v)
	}

	// with mask
	matrix = NewByteMatrix(d, d)
	clearMatrix(matrix)
	embedBasicPatterns(ver, matrix)
	embedTypeInfo(ecLevel, 0, matrix)
	maybeEmbedVersionInfo(ver, matrix)
	e = embedDataBits(dataBits, 0, matrix)
	if e != nil {
		t.Fatalf("embedDataBits returns error, %v", e)
	}
	v = 0
	for _, p := range pos {
		v <<= 1
		b, _ := MaskUtil_getDataMaskBit(0, p[0], p[1])
		i := matrix.Get(p[0], p[1])
		if (b && i == 0) || (!b && i == 1) {
			v |= 1
		}
	}
	if v != 125 {
		t.Fatalf("embedDataBits 125th byte = %v, expect 125", v)
	}
}

func TestMatrixUtil_buildMatrix(t *testing.T) {
	var e error
	var matrix *ByteMatrix
	ver, _ := decoder.Version_GetVersionForNumber(1)

	// fail on basic pattern
	matrix = NewByteMatrix(15, 15)
	e = MatrixUtil_buildMatrix(nil, decoder.ErrorCorrectionLevel_M, ver, 0, matrix)
	if e == nil {
		t.Fatalf("buildMatrix must be error")
	}

	// fail on type info
	matrix = NewByteMatrix(21, 21)
	e = MatrixUtil_buildMatrix(nil, decoder.ErrorCorrectionLevel_M, ver, 8, matrix)
	if e == nil {
		t.Fatalf("buildMatrix must be error")
	}

	qrstr := "" +
		"##############    ##  ####  ##############\n" +
		"##          ##  ####  ##    ##          ##\n" +
		"##  ######  ##  ####    ##  ##  ######  ##\n" +
		"##  ######  ##    ##  ##    ##  ######  ##\n" +
		"##  ######  ##  ##      ##  ##  ######  ##\n" +
		"##          ##  ##    ####  ##          ##\n" +
		"##############  ##  ##  ##  ##############\n" +
		"                ##########                \n" +
		"####  ##    ####  ####      ######  ####  \n" +
		"  ##########  ######        ##        ####\n" +
		"    ####  ########  ##  ####      ####  ##\n" +
		"      ##  ##    ##    ##          ##  ####\n" +
		"        ##  ####  ####  ##  ##  ##        \n" +
		"                ########      ####  ##  ##\n" +
		"##############  ######    ##  ##  ######  \n" +
		"##          ##    ##########  ####        \n" +
		"##  ######  ##    ##  ##    ######      ##\n" +
		"##  ######  ##  ##  ####      ##  ########\n" +
		"##  ######  ##    ####  ##      ##  ##  ##\n" +
		"##          ##  ######    ####            \n" +
		"##############  ##  ######    ##  ##  ##  "
	img, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	hellodata := []int{
		0x40, 0x56, 0x86, 0x56, 0xc6, 0xc6, 0xf0, 0xec, 0x11, 0xec,
		0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0x25,
		0x19, 0xd0, 0xd2, 0x68, 0x59, 0x39,
	}
	bits := gozxing.NewEmptyBitArray()
	for _, c := range hellodata {
		bits.AppendBits(c, 8)
	}
	e = MatrixUtil_buildMatrix(bits, decoder.ErrorCorrectionLevel_L, ver, 7, matrix)

	for y := 0; y < matrix.GetHeight(); y++ {
		for x := 0; x < matrix.GetWidth(); x++ {
			b := matrix.Get(x, y) == 1
			exp := img.Get(x, y)
			if b != exp {
				t.Fatalf("bits (%v,%v) = %v, expect %v", x, y, b, exp)
			}
		}
	}
}

func TestEmbedTimingPatterns(t *testing.T) {
	version, _ := decoder.Version_GetVersionForNumber(40)
	dimension := version.GetDimensionForVersion()
	matrix := NewByteMatrix(dimension, dimension)
	matrix.Clear(-1)

	embedTimingPatterns(matrix)

	for i := 8; i < matrix.GetWidth()-8; i++ {
		bit := matrix.Get(i, 6)
		wants := (i + 1) % 2
		if bit != int8(wants) {
			t.Fatalf("matrix(%d,6) = %v, wants %v", i, bit, wants)
		}
	}

	for i := 8; i < matrix.GetHeight()-8; i++ {
		bit := matrix.Get(6, i)
		wants := (i + 1) % 2
		if bit != int8(wants) {
			t.Fatalf("matrix(6, %d) = %v, wants %v", i, bit, wants)
		}
	}
}
