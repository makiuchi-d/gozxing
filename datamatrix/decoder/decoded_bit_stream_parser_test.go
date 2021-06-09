package decoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

func testUnrandomize255State(t testing.TB, randomized, codewordPosition, expect int) {
	t.Helper()
	c := unrandomize255State(randomized, codewordPosition)
	if c != expect {
		t.Fatalf("unrandomize255State(%v,%v) = %v, expect %v", randomized, codewordPosition, c, expect)
	}
}

func TestUnrandomize255State(t *testing.T) {
	testUnrandomize255State(t, 0, 0, 255)
	testUnrandomize255State(t, 10, 0, 9)
	testUnrandomize255State(t, 255, 0, 254)
	testUnrandomize255State(t, 0, 10, 40)
	testUnrandomize255State(t, 10, 10, 50)
	testUnrandomize255State(t, 255, 10, 39)
	testUnrandomize255State(t, 0, 255, 255)
	testUnrandomize255State(t, 10, 255, 9)
	testUnrandomize255State(t, 255, 255, 254)
}

func TestDecodeBase256Segment(t *testing.T) {
	result := make([]byte, 0)
	byteSegments := make([][]byte, 0)

	// d1=0, count=0
	bits := common.NewBitSource([]byte{150})
	r, _, e := decodeBase256Segment(bits, result, byteSegments)
	if e != nil {
		t.Fatalf("decodeBase256Segment returns error, %v", e)
	}
	if len(r) != 0 {
		t.Fatalf("decodeBase256Segment = %v, expect []", r)
	}

	// d1=1, count=1, but length=0
	bits = common.NewBitSource([]byte{151})
	_, _, e = decodeBase256Segment(bits, result, byteSegments)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeBase256Segment must be FormatException, %T", e)
	}

	// båd, d1=0
	bits = common.NewBitSource([]byte{150, 142, 166, 187})
	expect := "båd"
	expectSeg := [][]byte{{0x62, 0xe5, 0x64}}
	r, s, e := decodeBase256Segment(bits, result, byteSegments)
	if e != nil {
		t.Fatalf("decodeBase256Segment returns error, %v", e)
	}
	if string(r) != expect {
		t.Fatalf("decodeBase256Segment = \"%v\", expect \"%v\"", string(r), expect)
	}
	if !reflect.DeepEqual(s, expectSeg) {
		t.Fatalf("decodeBase256Segment segments = %v, expect = %v", s, expectSeg)
	}

	// count=256, d1=250, d2=6
	bits = common.NewBitSource([]byte{144, 50,
		192, 86, 235, 129, 23, 172, 66, 215, 109, 3, 152, 46, 195, 89, 238, 132, 26,
		175, 69, 218, 112, 6, 155, 49, 198, 92, 241, 135, 29, 178, 72, 221, 115, 9,
		158, 52, 201, 95, 244, 138, 32, 181, 75, 224, 118, 12, 161, 55, 204, 98, 247,
		141, 35, 184, 78, 227, 121, 15, 164, 58, 207, 101, 250, 144, 38, 187, 81,
		230, 124, 18, 167, 61, 210, 104, 253, 147, 41, 190, 84, 233, 127, 21, 170,
		64, 213, 107, 1, 150, 44, 193, 87, 236, 130, 24, 173, 67, 216, 110, 4, 153,
		47, 196, 90, 239, 133, 27, 176, 70, 219, 113, 7, 156, 50, 199, 93, 242, 136,
		30, 179, 73, 222, 116, 10, 159, 53, 202, 96, 245, 139, 33, 182, 76, 225, 119,
		13, 162, 56, 205, 99, 248, 142, 36, 185, 79, 228, 122, 16, 165, 59, 208, 102,
		251, 145, 39, 188, 82, 231, 125, 19, 168, 62, 211, 105, 254, 148, 42, 191,
		85, 234, 128, 22, 171, 65, 214, 108, 2, 151, 45, 194, 88, 237, 131, 25, 174,
		68, 217, 111, 5, 154, 48, 197, 91, 240, 134, 28, 177, 71, 220, 114, 8, 157,
		51, 200, 94, 243, 137, 31, 180, 74, 223, 117, 11, 160, 54, 203, 97, 246, 140,
		34, 183, 77, 226, 120, 14, 163, 57, 206, 100, 249, 143, 37, 186, 80, 229,
		123, 17, 166, 60, 209, 103, 252, 146, 40, 189, 83, 232, 126, 20, 169, 63,
		212, 106, 0, 149, 43, 192,
	})
	expect = "" +
		"ÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿ" +
		"ÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿ" +
		"ÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿ" +
		"ÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿ"
	expectSeg = [][]byte{{
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	}}
	r, s, e = decodeBase256Segment(bits, result, byteSegments)
	if e != nil {
		t.Fatalf("decodeBase256Segment returns error, %v", e)
	}
	if string(r) != expect {
		t.Fatalf("decodeBase256Segment = \"%v\", expect \"%v\"", string(r), expect)
	}
	if !reflect.DeepEqual(s, expectSeg) {
		t.Fatalf("decodeBase256Segment segments = %v, expect = %v", s, expectSeg)
	}
}

func TestEecodeEdifactSegment(t *testing.T) {
	result := make([]byte, 0)

	bits := common.NewBitSource([]byte{0})
	r := decodeEdifactSegment(bits, result)
	if len(r) > 0 {
		t.Fatalf("decodeEdifactSegment = %v, expect []", r)
	}

	// 101010 010101 011111 00
	// -> 00101010, 01010101
	bits = common.NewBitSource([]byte{0xa9, 0x57, 0xc0})
	expect := []byte{0x2a, 0x55}
	r = decodeEdifactSegment(bits, result)
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeEdifactSegment = %v, expect %v", r, expect)
	}
	if r := bits.GetBitOffset(); r != 0 {
		t.Fatalf("decodeEdifactSegment result.GetBitOffset must be 0, %v", r)
	}

	// 101010 010101 101010 010101
	bits = common.NewBitSource([]byte{0xa9, 0x5a, 0x95})
	expect = []byte{0x2a, 0x55, 0x2a, 0x55}
	r = decodeEdifactSegment(bits, result)
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeEdifactSegment = %v, expect %v", r, expect)
	}
	if r := bits.GetBitOffset(); r != 0 {
		t.Fatalf("decodeEdifactSegment result.GetBitOffset must be 0, %v", r)
	}
}

func TestDecodeAnsiX12Segment(t *testing.T) {
	result := make([]byte, 0)

	bits := common.NewBitSource([]byte{})
	r, e := decodeAnsiX12Segment(bits, result)
	if e != nil {
		t.Fatalf("decodeAnsiX12Segment returns error, %v", e)
	}
	if len(r) != 0 {
		t.Fatalf("decodeAnsiX12Segment = %v, expect []", r)
	}

	bits = common.NewBitSource([]byte{251, 0})
	_, e = decodeAnsiX12Segment(bits, result)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeAnsiX12Segment must be FormatException, %T", e)
	}

	bits = common.NewBitSource([]byte{33, 116, 12, 169})
	expect := []byte("1A >*\r")
	r, e = decodeAnsiX12Segment(bits, result)
	if e != nil {
		t.Fatalf("decodeAnsiX12Segment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeAnsiX12Segment = %v, expect %v", r, expect)
	}

	bits = common.NewBitSource([]byte{33, 116, 254, 0})
	expect = []byte("1A ")
	r, e = decodeAnsiX12Segment(bits, result)
	if e != nil {
		t.Fatalf("decodeAnsiX12Segment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeAnsiX12Segment = %v, expect %v", r, expect)
	}
	if r := bits.Available(); r != 8 {
		t.Fatalf("decodeAnsiX12Segment bits.Available = %v, expect 8", r)
	}

	bits = common.NewBitSource([]byte{87, 90, 64})
	expect = []byte("9Z*")
	r, e = decodeAnsiX12Segment(bits, result)
	if e != nil {
		t.Fatalf("decodeAnsiX12Segment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeAnsiX12Segment = %v, expect %v", r, expect)
	}
	if r := bits.Available(); r != 8 {
		t.Fatalf("decodeAnsiX12Segment bits.Available = %v, expect 8", r)
	}
}

func TestDecodeTextSegment(t *testing.T) {
	result := make([]byte, 0)

	// basic set
	// 40,0,0 (out of range)
	bits := common.NewBitSource([]byte{250, 1})
	fnc1poss := intSet{}
	_, e := decodeTextSegment(bits, result, fnc1poss)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeTextSegment must be FormatException, %T", e)
	}
	// 14,1,30, 3,15,16 ('a',uppser-shift,128+32,'b','c')
	bits = common.NewBitSource([]byte{87, 199, 21, 41})
	expect := []byte{'a', 128 + 32, 'b', 'c'}
	r, e := decodeTextSegment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}

	// shift-1 set
	// 0,13,1, 30,0,0 ('\r',upper-shift,128+0)
	bits = common.NewBitSource([]byte{2, 10, 187, 129})
	expect = []byte{'\r', 128}
	r, e = decodeTextSegment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}

	// shift-2
	// 1,28,0 (out of range)
	bits = common.NewBitSource([]byte{10, 161})
	_, e = decodeTextSegment(bits, result, fnc1poss)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeTextSegment must be FormatException, %T", e)
	}
	// 1,0,1, 27,1,30, 1,26,3 ('!',FNC1,upper-shift,128+95)
	bits = common.NewBitSource([]byte{6, 66, 169, 7, 10, 84})
	expect = []byte{'!', 29, 128 + 95, ' '}
	r, e = decodeTextSegment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}

	// shift-3
	// 2,32,0 (out of range)
	bits = common.NewBitSource([]byte{17, 129})
	_, e = decodeTextSegment(bits, result, fnc1poss)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeTextSegment must be FormatException, %T", e)
	}
	// 2,0,1,30,2,31 ('`', upper-shift, 128+127)
	bits = common.NewBitSource([]byte{12, 130, 187, 240})
	expect = []byte{'`', 255}
	r, e = decodeTextSegment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}

	// rest bits handling
	// 4,5,6, <8bit>
	bits = common.NewBitSource([]byte{25, 207, 0})
	expect = []byte("012")
	r, e = decodeTextSegment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}
	if r := bits.Available(); r != 8 {
		t.Fatalf("decodeTextSegment rest bits = %v, expect 8", r)
	}

	// 4,5,6, <254><8bit>
	bits = common.NewBitSource([]byte{25, 207, 254, 0})
	expect = []byte("012")
	r, e = decodeTextSegment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}
	if r := bits.Available(); r != 8 {
		t.Fatalf("decodeTextSegment rest bits = %v, expect 8", r)
	}
}

func TestDecodeC40Segment(t *testing.T) {
	result := make([]byte, 0)

	// basic set
	// 40,0,0 (out of range)
	bits := common.NewBitSource([]byte{250, 1})
	fnc1poss := intSet{}
	_, e := decodeC40Segment(bits, result, fnc1poss)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeC40Segment must be FormatException, %T", e)
	}
	// 14,1,30, 3,15,16 ('A',uppser-shift,128+32,'B','C')
	bits = common.NewBitSource([]byte{87, 199, 21, 41})
	expect := []byte{'A', 128 + 32, 'B', 'C'}
	r, e := decodeC40Segment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeC40Segment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeC40Segment = %v, expect %v", r, expect)
	}

	// shift-1 set
	// 0,13,1, 30,0,0 ('\r',upper-shift,128+0)
	bits = common.NewBitSource([]byte{2, 10, 187, 129})
	expect = []byte{'\r', 128}
	r, e = decodeC40Segment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeC40Segment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeC40Segment = %v, expect %v", r, expect)
	}

	// shift-2
	// 1,28,0 (out of range)
	bits = common.NewBitSource([]byte{10, 161})
	_, e = decodeC40Segment(bits, result, fnc1poss)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeC40Segment must be FormatException, %T", e)
	}
	// 1,0,1, 27,1,30, 1,26,3 ('!',FNC1,upper-shift,128+95)
	bits = common.NewBitSource([]byte{6, 66, 169, 7, 10, 84})
	expect = []byte{'!', 29, 128 + 95, ' '}
	r, e = decodeC40Segment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeC40Segment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeC40Segment = %v, expect %v", r, expect)
	}

	// shift-3
	// 2,0,1,30,2,31 ('`', upper-shift, 128+127)
	bits = common.NewBitSource([]byte{12, 130, 187, 240})
	expect = []byte{'`', 255}
	r, e = decodeC40Segment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeC40Segment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeC40Segment = %v, expect %v", r, expect)
	}

	// rest bits handling
	// 4,5,6, <8bit>
	bits = common.NewBitSource([]byte{25, 207, 0})
	expect = []byte("012")
	r, e = decodeC40Segment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}
	if r := bits.Available(); r != 8 {
		t.Fatalf("decodeTextSegment rest bits = %v, expect 8", r)
	}

	// 4,5,6, <254><8bit>
	bits = common.NewBitSource([]byte{25, 207, 254, 0})
	expect = []byte("012")
	r, e = decodeC40Segment(bits, result, fnc1poss)
	if e != nil {
		t.Fatalf("decodeTextSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeTextSegment = %v, expect %v", r, expect)
	}
	if r := bits.Available(); r != 8 {
		t.Fatalf("decodeTextSegment rest bits = %v, expect 8", r)
	}
}

func testDecodeAsciiSegment(t testing.TB, bits *common.BitSource, mode Mode, result, trailer []byte) {
	t.Helper()
	r := make([]byte, 0)
	rt := make([]byte, 0)
	fnc1poss := intSet{}
	m, r, rt, e := decodeAsciiSegment(bits, r, rt, fnc1poss)
	if e != nil {
		t.Fatalf("decodeAsciiSegment(%v) return error, %v", bits, e)
	}
	if m != mode {
		t.Fatalf("decodeAsciiSegment(%v) mode = %v, expect %v", bits, m, mode)
	}
	if !reflect.DeepEqual(r, result) {
		t.Fatalf("decodeAsciiSegment(%v) = %v, expect %v", bits, r, result)
	}
	if !reflect.DeepEqual(rt, trailer) {
		t.Fatalf("decodeAsciiSegment(%v) trailer = %v, expect %v", bits, r, trailer)
	}
}

func TestDecodeAsciiSegment(t *testing.T) {
	bits := common.NewBitSource([]byte{0})
	fnc1poss := intSet{}
	_, _, _, e := decodeAsciiSegment(bits, []byte{}, []byte{}, fnc1poss)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeAsciiSegment must be FormatException, %T", e)
	}

	// ASCII
	bits = common.NewBitSource([]byte{'A' + 1})
	mode := Mode_ASCII_ENCODE
	expect := []byte{'A'}
	trailer := []byte{}

	// padding
	bits = common.NewBitSource([]byte{129})
	mode = Mode_PDA_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// mumeric
	bits = common.NewBitSource([]byte{135})
	mode = Mode_ASCII_ENCODE
	expect = []byte("05")
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	bits = common.NewBitSource([]byte{229})
	mode = Mode_ASCII_ENCODE
	expect = []byte("99")
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// launch to C40 encodation
	bits = common.NewBitSource([]byte{230})
	mode = Mode_C40_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// launch to base 256 encodation
	bits = common.NewBitSource([]byte{231})
	mode = Mode_BASE256_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// FNC1
	bits = common.NewBitSource([]byte{232})
	mode = Mode_ASCII_ENCODE
	expect = []byte{29}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// structured append (not implemented)
	bits = common.NewBitSource([]byte{233})
	mode = Mode_ASCII_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)
	bits = common.NewBitSource([]byte{234})
	mode = Mode_ASCII_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// upper-shift + extended ascii 200
	bits = common.NewBitSource([]byte{235, 200 - 128 + 1})
	mode = Mode_ASCII_ENCODE
	expect = []byte{200}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// 05 Macro
	bits = common.NewBitSource([]byte{236})
	mode = Mode_ASCII_ENCODE
	expect = []byte("[)>\u001E05\u001D")
	trailer = []byte("\u001E\u0004")
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// 06 Macro
	bits = common.NewBitSource([]byte{237})
	mode = Mode_ASCII_ENCODE
	expect = []byte("[)>\u001E06\u001D")
	trailer = []byte("\u001E\u0004")
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// launch ANSI X12 encodation
	bits = common.NewBitSource([]byte{238})
	mode = Mode_ANSIX12_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// launch text encodation
	bits = common.NewBitSource([]byte{239})
	mode = Mode_TEXT_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// launch EDIFACT encodation
	bits = common.NewBitSource([]byte{240})
	mode = Mode_EDIFACT_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// ECI encodation
	bits = common.NewBitSource([]byte{241})
	mode = Mode_ECI_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// latch back to ASCII
	bits = common.NewBitSource([]byte{254})
	mode = Mode_ASCII_ENCODE
	expect = []byte{}
	trailer = []byte{}
	testDecodeAsciiSegment(t, bits, mode, expect, trailer)

	// invalid code
	for i := 250; i < 256; i++ {
		bits = common.NewBitSource([]byte{byte(i), 'A' + 1})
		_, _, _, e := decodeAsciiSegment(bits, []byte{}, []byte{}, fnc1poss)
		if _, ok := e.(gozxing.FormatException); !ok {
			t.Fatalf("decodeAsciiSegment must be FormatException, %T", e)
		}
	}
}

func testDecodedBitStreamParser_decode(t testing.TB, bytes []byte, str string, segments [][]byte) {
	t.Helper()
	r, e := DecodedBitStreamParser_decode(bytes)
	if e != nil {
		t.Fatalf("DecodedBitStreamParser_decode(%v) returns error, %v", bytes, e)
	}
	if s := r.GetText(); s != str {
		t.Fatalf("DecodedBitStreamParser_decode(%v) = \"%v\", expect \"%v\"", bytes, s, str)
	}
	if r := r.GetByteSegments(); !reflect.DeepEqual(r, segments) {
		t.Fatalf("DecodedBitStreamParser_decode(%v) byteSegments = %v, expect %v", bytes, r, segments)
	}
}

func TestDecodedBitStreamParser_decode(t *testing.T) {
	bytes := []byte{'a' + 1, 0}
	_, e := DecodedBitStreamParser_decode(bytes)
	if e == nil {
		t.Fatalf("DecodedBitStreamParser_decode must be error")
	}

	// ASCII
	bytes = []byte{'a' + 1, 'b' + 1, 'c' + 1, 131}
	str := "abc01"
	segments := [][]byte(nil)
	testDecodedBitStreamParser_decode(t, bytes, str, segments)

	// ASCII with macro
	bytes = []byte{'a' + 1, 236, 'b' + 1, 237, 'c' + 1}
	str = "a[)>\u001E05\u001Db[)>\u001E06\u001Dc\u001E\u0004\u001E\u0004"
	segments = [][]byte(nil)
	testDecodedBitStreamParser_decode(t, bytes, str, segments)

	// C40, Text segments
	bytes = []byte{230, 87, 169, 254, 239, 87, 169}
	str = "A!a!"
	segments = [][]byte(nil)
	testDecodedBitStreamParser_decode(t, bytes, str, segments)

	// ANSIX12, EDIFACT, BASE256, PAD
	bytes = []byte{
		238, 7, 48,
		254, 240, 0xa9, 0x57, 0xc0,
		231, 219, 208, 233, 253,
		34, 129, 34, 34,
	}
	str = "*1Z*Ubåd!"
	segments = [][]byte{{0x62, 0xe5, 0x64}}
	testDecodedBitStreamParser_decode(t, bytes, str, segments)
}

func TestDecodedBitStreamParser_decode_SymbologyModifire(t *testing.T) {
	// no fnc1
	bytes := []byte{'0' + 1}
	sm := 1
	r, e := DecodedBitStreamParser_decode(bytes)
	if e != nil {
		t.Fatalf("decode error: %v", e)
	}
	if s := r.GetSymbologyModifier(); s != sm {
		t.Fatalf("SymbologModifier = %v, expect %v", s, sm)
	}

	// fnc1 at 0
	bytes = []byte{232}
	sm = 2
	r, e = DecodedBitStreamParser_decode(bytes)
	if e != nil {
		t.Fatalf("decode error: %v", e)
	}
	if s := r.GetSymbologyModifier(); s != sm {
		t.Fatalf("SymbologModifier = %v, expect %v", s, sm)
	}

	// fnc1 at 1
	bytes = []byte{'0' + 1, 232}
	sm = 3
	r, e = DecodedBitStreamParser_decode(bytes)
	if e != nil {
		t.Fatalf("decode error: %v", e)
	}
	if s := r.GetSymbologyModifier(); s != sm {
		t.Fatalf("SymbologModifier = %v, expect %v", s, sm)
	}

	// no fnc1 with ECI encoded
	bytes = []byte{241, '0' + 1}
	sm = 4
	r, e = DecodedBitStreamParser_decode(bytes)
	if e != nil {
		t.Fatalf("decode error: %v", e)
	}
	if s := r.GetSymbologyModifier(); s != sm {
		t.Fatalf("SymbologModifier = %v, expect %v", s, sm)
	}

	// fnc1 at 0 with ECI encoded
	bytes = []byte{241, 232, '0' + 1}
	sm = 5
	r, e = DecodedBitStreamParser_decode(bytes)
	if e != nil {
		t.Fatalf("decode error: %v", e)
	}
	if s := r.GetSymbologyModifier(); s != sm {
		t.Fatalf("SymbologModifier = %v, expect %v", s, sm)
	}

	// fnc1 at 1 with ECI encoded
	bytes = []byte{241, '0' + 1, 232, '0' + 1}
	sm = 6
	r, e = DecodedBitStreamParser_decode(bytes)
	if e != nil {
		t.Fatalf("decode error: %v", e)
	}
	if s := r.GetSymbologyModifier(); s != sm {
		t.Fatalf("SymbologModifier = %v, expect %v", s, sm)
	}
}

func TestMode_String(t *testing.T) {
	if s := Mode_PDA_ENCODE.String(); s != "PAD_ENCODE" {
		t.Fatalf("PAD_ENCODE string = %v", s)
	}
	if s := Mode_ASCII_ENCODE.String(); s != "ASCII_ENCODE" {
		t.Fatalf("ASCII_ENCODE string = %v", s)
	}
	if s := Mode_C40_ENCODE.String(); s != "C40_ENCODE" {
		t.Fatalf("C40_ENCODE string = %v", s)
	}
	if s := Mode_TEXT_ENCODE.String(); s != "TEXT_ENCODE" {
		t.Fatalf("TEXT_ENCODE string = %v", s)
	}
	if s := Mode_ANSIX12_ENCODE.String(); s != "ANSIX12_ENCODE" {
		t.Fatalf("ANSIX12_ENCODE string = %v", s)
	}
	if s := Mode_EDIFACT_ENCODE.String(); s != "EDIFACT_ENCODE" {
		t.Fatalf("EDIFACT_ENCODE string = %v", s)
	}
	if s := Mode_BASE256_ENCODE.String(); s != "BASE256_ENCODE" {
		t.Fatalf("BASE256_ENCODE string = %v", s)
	}
	if s := Mode_ECI_ENCODE.String(); s != "ECI_ENCODE" {
		t.Fatalf("ECI_ENCODE string = %v", s)
	}
	if s := Mode(-1).String(); s != "" {
		t.Fatalf("unknown mode string = \"%v\"", s)
	}
}
