package encoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestRandomize253State(t *testing.T) {
	if r := randomize253State(1); r != 25 {
		t.Fatalf("randomize253State(1) = %v, expect 25", r)
	}
	if r := randomize253State(2); r != 175 {
		t.Fatalf("randomize253State(2) = %v, expect 175", r)
	}
	if r := randomize253State(10); r != 101 {
		t.Fatalf("randomize253State(10) = %v, expect 101", r)
	}
}

func TestFindMinimums(t *testing.T) {
	mins := make([]byte, 6)
	intCharCounts := make([]int, 6)
	min := findMinimums(
		[]float64{7.2, 1.5, 6, 2.2, 5, 2}, intCharCounts, 100, mins)

	if min != 2 {
		t.Fatalf("findMinimums min = %v, expect 2", min)
	}
	if !reflect.DeepEqual(intCharCounts, []int{8, 2, 6, 3, 5, 2}) {
		t.Fatalf("findMinimums intCharCounts = %v, expect [8 2 6 3 5 2]", intCharCounts)
	}
	if !reflect.DeepEqual(mins, []byte{0, 1, 0, 0, 0, 1}) {
		t.Fatalf("findMinimums intCharCounts = %v, expect [0 1 0 0 0 1]", mins)
	}
}

func TestGetMinimumCount(t *testing.T) {
	if r := getMinimumCount([]byte{0, 0, 0, 0, 0, 0}); r != 0 {
		t.Fatalf("getMinimumCount = %v, expect 0", r)
	}
	if r := getMinimumCount([]byte{0, 1, 0, 1, 1, 0}); r != 3 {
		t.Fatalf("getMinimumCount = %v, expect 3", r)
	}
	if r := getMinimumCount([]byte{1, 1, 1, 1, 1, 1}); r != 6 {
		t.Fatalf("getMinimumCount = %v, expect 6", r)
	}
}

func TestHighLevelEncoder_isDigit(t *testing.T) {
	if HighLevelEncoder_isDigit('/') != false {
		t.Fatalf("isDigit('/') must false")
	}
	if HighLevelEncoder_isDigit('0') != true {
		t.Fatalf("isDigit('0') must true")
	}
	if HighLevelEncoder_isDigit('9') != true {
		t.Fatalf("isDigit('9') must true")
	}
	if HighLevelEncoder_isDigit(':') != false {
		t.Fatalf("isDigit(':') must false")
	}
}

func TestHighLevelEncoder_isExtendedASCII(t *testing.T) {
	if HighLevelEncoder_isExtendedASCII(' ') != false {
		t.Fatalf("isExtendedASCII(' ') must false")
	}
	if HighLevelEncoder_isExtendedASCII(0x7f) != false {
		t.Fatalf("isExtendedASCII(0x7f) must false")
	}
	if HighLevelEncoder_isExtendedASCII(0x80) != true {
		t.Fatalf("isExtendedASCII(0x80) must true")
	}
	if HighLevelEncoder_isExtendedASCII(0xff) != true {
		t.Fatalf("isExtendedASCII(0xff) must true")
	}
}

func TestIsNativeC40(t *testing.T) {
	tcs := []byte{' ', '0', '9', 'A', 'Z'}
	fcs := []byte{'!', '@', '^', 'a', '~'}
	for _, c := range tcs {
		if isNativeC40(c) != true {
			t.Fatalf("isNativeC40(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeC40(c) != false {
			t.Fatalf("isNativeC40(%v) must false", c)
		}
	}
}

func TestIsNativeText(t *testing.T) {
	tcs := []byte{' ', '0', '9', 'a', 'z'}
	fcs := []byte{'!', '@', '^', 'A', '~'}
	for _, c := range tcs {
		if isNativeText(c) != true {
			t.Fatalf("isNativeText(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeText(c) != false {
			t.Fatalf("isNativeText(%v) must false", c)
		}
	}
}

func TestIsNativeX12(t *testing.T) {
	tcs := []byte{'\r', '*', '>', ' ', '0', '9', 'A', 'Z'}
	fcs := []byte{'!', '@', '^', 'a', '~'}
	for _, c := range tcs {
		if isNativeX12(c) != true {
			t.Fatalf("isNativeX12(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeX12(c) != false {
			t.Fatalf("isNativeX12(%v) must false", c)
		}
	}
}

func TestIsNativeEDIFACT(t *testing.T) {
	tcs := []byte{' ', '!', '@', '0', 'A', '^'}
	fcs := []byte{'\r', '_', 'a', '~'}
	for _, c := range tcs {
		if isNativeEDIFACT(c) != true {
			t.Fatalf("isNativeEDIFACT(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeEDIFACT(c) != false {
			t.Fatalf("isNativeEDIFACT(%v) must false", c)
		}
	}
}

func TestIsSpecialB256(t *testing.T) {
	if isSpecialB256(0) {
		t.Fatalf("isSpecialB256 must false (not implemented yet)")
	}
}

func TestHighLevelEncoder_determineConsecutiveDigitCount(t *testing.T) {
	msg := []byte("   0123bc")
	if r := HighLevelEncoder_determineConsecutiveDigitCount(msg, 0); r != 0 {
		t.Fatalf("determineConsecutiveDigitCount = %v, expect 0", r)
	}
	if r := HighLevelEncoder_determineConsecutiveDigitCount(msg, 2); r != 0 {
		t.Fatalf("determineConsecutiveDigitCount = %v, expect 0", r)
	}
	if r := HighLevelEncoder_determineConsecutiveDigitCount(msg, 3); r != 4 {
		t.Fatalf("determineConsecutiveDigitCount = %v, expect 4", r)
	}
	if r := HighLevelEncoder_determineConsecutiveDigitCount(msg, 6); r != 1 {
		t.Fatalf("determineConsecutiveDigitCount = %v, expect 1", r)
	}
	if r := HighLevelEncoder_determineConsecutiveDigitCount(msg, 7); r != 0 {
		t.Fatalf("determineConsecutiveDigitCount = %v, expect 0", r)
	}
}

func TestHighLevelEncoder_lookAheadTest(t *testing.T) {
	msg := []byte{}
	pos := 0
	mode := HighLevelEncoder_ASCII_ENCODATION
	expect := HighLevelEncoder_ASCII_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte("000")
	pos = 0
	mode = HighLevelEncoder_BASE256_ENCODATION
	expect = HighLevelEncoder_ASCII_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte{0xff, 0xff, 0xff, 0xff}
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_BASE256_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte("^^^")
	pos = 1
	mode = HighLevelEncoder_EDIFACT_ENCODATION
	expect = HighLevelEncoder_EDIFACT_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte(" 0a")
	pos = 1
	mode = HighLevelEncoder_TEXT_ENCODATION
	expect = HighLevelEncoder_TEXT_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte("*>\r")
	pos = 1
	mode = HighLevelEncoder_X12_ENCODATION
	expect = HighLevelEncoder_X12_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte(" 0A")
	pos = 1
	mode = HighLevelEncoder_C40_ENCODATION
	expect = HighLevelEncoder_C40_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte("000000")
	pos = 1
	mode = HighLevelEncoder_BASE256_ENCODATION
	expect = HighLevelEncoder_ASCII_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_BASE256_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte("^^^^^^^^^^^")
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_EDIFACT_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte(" 0a 0a 0a")
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_TEXT_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte("*>\r*>\r*>\r")
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_X12_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte(" 0A 0A")
	pos = 1
	mode = HighLevelEncoder_C40_ENCODATION
	expect = HighLevelEncoder_C40_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte(" 0A 0A 0A 0A 0A 0")
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_C40_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte(" 0A 0A 0A 0A 0A 0*")
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_X12_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}

	msg = []byte(" 0A 0A 0A 0A 0A 0A^")
	pos = 1
	mode = HighLevelEncoder_ASCII_ENCODATION
	expect = HighLevelEncoder_C40_ENCODATION
	if r := HighLevelEncoder_lookAheadTest(msg, pos, mode); r != expect {
		t.Fatalf("lookAheadTest = %v, expect %v", r, expect)
	}
}

func TestEncodeHighLevel(t *testing.T) {
	shape := SymbolShapeHint_FORCE_NONE

	_, e := EncodeHighLevel("Mосква", shape, nil, nil)
	if e == nil {
		t.Fatalf("EncodeHighLevel(Mосква) must be error")
	}

	str := string(make([]byte, 1559))
	_, e = EncodeHighLevel(str, shape, nil, nil)
	if e == nil {
		t.Fatalf("EncodeHighLevel([1559]) must be error")
	}

	str = " 0A 0A 0A 0A000000"
	expect := []byte{
		230, 19, 111, 19, 111, 19, 111, 19, 111, 254, 130, 130, 130, 129, 87, 237, 133, 28,
	}
	min, _ := gozxing.NewDimension(16, 16)
	b, e := EncodeHighLevel(str, shape, min, nil)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("EncodeHighLevel = %v, expect %v", b, expect)
	}

	str = "[)>\u001E05\u001Daaaaaa\u001E\u0004"
	b, e = EncodeHighLevel(str, shape, nil, nil)
	expect = []byte{236, 239, 89, 191, 89, 191, 254, 129}
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("EncodeHighLevel = %v, expect %v", b, expect)
	}

	str = "[)>\u001E06\u001Daaaaaa\u001E\u0004"
	b, e = EncodeHighLevel(str, shape, nil, nil)
	expect = []byte{237, 239, 89, 191, 89, 191, 254, 129}
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("EncodeHighLevel = %v, expect %v", b, expect)
	}
}
