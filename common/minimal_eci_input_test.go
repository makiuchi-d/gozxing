package common

import (
	"reflect"
	"testing"

	"golang.org/x/text/encoding/japanese"
)

func TestMinimalECIInput(t *testing.T) {
	inp, err := NewMinimalECIInput([]rune("\u001dÀabc"), nil, 0x1d)
	if err != nil {
		t.Fatalf("NewMinimalECIInput error: %v", err)
	}
	if fnc1 := inp.GetFNC1Character(); fnc1 != 0x1d {
		t.Errorf("GetFNC1Character = %#x wants %#x", fnc1, 0x1d)
	}
	if l, exp := inp.Length(), 5; l != exp {
		t.Errorf("Length = %v wants %v", l, exp)
	}
	if !inp.HaveNCharacters(0, 5) {
		t.Errorf("HaveNCharacters(0, 5) = false wants true")
	}
	if inp.HaveNCharacters(0, 6) {
		t.Errorf("HaveNCharacters(0, 6) = true wants false")
	}
	if c, err := inp.CharAt(0); err != nil {
		t.Errorf("CharAt(0) error: %v", err)
	} else if c != 0x1d {
		t.Errorf("CharAt(0) = %#x wants %#x", c, 0x1d)
	}
	if _, err := inp.CharAt(5); err == nil {
		t.Errorf("CharAt(5) must be error")
	}
	if ss, err := inp.SubSequence(1, 3); err != nil {
		t.Errorf("SubSequence(1, 3) error: %v", err)
	} else if exp := []int{0xc0, 0x61}; !reflect.DeepEqual(ss, exp) {
		t.Errorf("SubSequence(1, 3) = %q wants %q", ss, exp)
	}
	if _, err := inp.SubSequence(1, 6); err == nil {
		t.Errorf("SubSequence(1, 6) must be error")
	}
	if b, err := inp.IsFNC1(0); err != nil {
		t.Errorf("IsFNC1(0) error: %v", err)
	} else if !b {
		t.Errorf("IsFNC1(0) = false wants true")
	}
	if b, err := inp.IsFNC1(1); err != nil {
		t.Errorf("IsFNC1(1) error: %v", err)
	} else if b {
		t.Errorf("IsFNC1(1) = true wants false")
	}
	if _, err := inp.IsFNC1(-1); err == nil {
		t.Errorf("IsFNC1(-1) must be error")
	}
	if s, exp := inp.String(), "'\x1d', 192, 'a', 'b', 'c'"; s != exp {
		t.Errorf("String() = %q wants %q", s, exp)
	}
}

func TestMinimalECIInputWithECI(t *testing.T) {
	// c0, ECI(20), 82, a0, 61, 62, 63
	inp, err := NewMinimalECIInput([]rune("Àあabc"), nil, -1)
	if err != nil {
		t.Errorf("NewMinimalECIInput error: %v", err)
	}
	if l, exp := inp.Length(), 7; l != exp {
		t.Errorf("Length = %v wants %v", l, exp)
	}
	if !inp.HaveNCharacters(0, 1) {
		t.Errorf("HaveNCharacters(0, 1) = false wants true")
	}
	if inp.HaveNCharacters(0, 2) {
		t.Errorf("HaveNCharacters(0, 2) = true wants false")
	}
	if !inp.HaveNCharacters(2, 5) {
		t.Errorf("HaveNCharacters(2, 5) = false wants true")
	}
	if c, err := inp.CharAt(0); err != nil {
		t.Errorf("CharAt(0) error: %v", err)
	} else if c != 0xc0 {
		t.Errorf("CharAt(0) = %#x wants %#x", c, 0xc0)
	}
	if _, err := inp.CharAt(1); err == nil {
		t.Errorf("CharAt(1) must be error")
	}
	if ss, err := inp.SubSequence(0, 1); err != nil {
		t.Errorf("SubSequence(0, 1) error: %v", err)
	} else if exp := []int{0xc0}; !reflect.DeepEqual(ss, exp) {
		t.Errorf("SubSequence(0, 1) = %x wants %x", ss, exp)
	}
	if _, err := inp.SubSequence(0, 2); err == nil {
		t.Errorf("SubSequence(0, 2) must be error")
	}
	if ss, err := inp.SubSequence(2, 5); err != nil {
		t.Errorf("SubSequence(2, 5) error: %v", err)
	} else if exp := []int{0x82, 0xa0, 0x61}; !reflect.DeepEqual(ss, exp) {
		t.Errorf("SubSequence(2, 5) = %x wants %x", ss, exp)
	}
	if b, err := inp.IsECI(0); err != nil {
		t.Errorf("IsECI(0) error: %v", err)
	} else if b {
		t.Errorf("IsECI(0) = true wants false")
	}
	if b, err := inp.IsECI(1); err != nil {
		t.Errorf("IsECI(1) error: %v", err)
	} else if !b {
		t.Errorf("IsECI(1) = false wants true")
	}
	if _, err := inp.IsECI(-1); err == nil {
		t.Errorf("IsECI(-1) must be error")
	}
	if _, err := inp.GetECIValue(0); err == nil {
		t.Errorf("GetECIValue(0) must be error")
	}
	if _, err := inp.GetECIValue(-1); err == nil {
		t.Errorf("GetECIValue(-1) must be error")
	}
	if v, err := inp.GetECIValue(1); err != nil {
		t.Errorf("GetECIValue(1) error: %v", err)
	} else if exp := 20; v != exp {
		t.Errorf("GetECIValue(1) = %v wants %v", v, exp)
	}
	if s, exp := inp.String(), "192, ECI(20), 130, 160, 'a', 'b', 'c'"; s != exp {
		t.Errorf("String() = %q wants %q", s, exp)
	}
}

func TestMinimalECIInput_encodeMinimally(t *testing.T) {
	sjis := NewECIEncoderSet([]rune("あ"), japanese.ShiftJIS, 0x1d)
	latin1 := NewECIEncoderSet([]rune("abc"), nil, 0x1d)

	str := []rune("\u001dÀabc")
	exp := []int{fnc1Value, 0xc0, 0x61, 0x62, 0x63}
	bytes, err := minimalECIInput_encodeMinimally(str, latin1, 0x1d)
	if err != nil {
		t.Errorf("encodeMinimally(%q) error: %v", string(str), err)
	} else if !reflect.DeepEqual(bytes, exp) {
		t.Fatalf("encodeMinimally(%q): %v wants %v", string(str), bytes, exp)
	}

	str = []rune("Àあ")
	_, err = minimalECIInput_encodeMinimally(str, latin1, -1)
	if err == nil {
		t.Errorf("encodeMinimally(%q) must be error", string(str))
	}

	exp = []int{0xc0, 256 + 20, 0x82, 0xa0}
	bytes, err = minimalECIInput_encodeMinimally(str, sjis, -1)
	if err != nil {
		t.Errorf("encodeMinimally(%q) error: %v", string(str), err)
	} else if !reflect.DeepEqual(bytes, exp) {
		t.Errorf("encodeMinimally(%q) = %v wants %v", string(str), bytes, exp)
	}
}

func TestInputEdge(t *testing.T) {
	latin1 := NewECIEncoderSet([]rune("abc"), nil, 0x1d)
	utf8 := NewECIEncoderSet([]rune("あ"), nil, 0x1d)

	tests := map[string]struct {
		c      rune
		encset *ECIEncoderSet
		idx    int
		prev   *inputEdge
		fnc1   int

		cVal      int
		encIdx    int
		totalSize int
		isFnc1    bool
	}{
		"latin1 a": {
			c:         'a',
			encset:    latin1,
			idx:       0,
			prev:      nil,
			fnc1:      -1,
			cVal:      0x61,
			encIdx:    0,
			totalSize: 1,
			isFnc1:    false,
		},
		"latin1 fnc1": {
			c:         '\u001d',
			encset:    latin1,
			idx:       0,
			prev:      nil,
			fnc1:      0x1d,
			cVal:      fnc1Value,
			encIdx:    0,
			totalSize: 1,
			isFnc1:    true,
		},
		"utf8-latin1": {
			c:      'あ',
			encset: utf8,
			idx:    1,
			prev: &inputEdge{
				c:               'a',
				encoderIndex:    0,
				previous:        nil,
				cachedTotalSize: 1,
			},
			fnc1:      0x1d,
			cVal:      'あ',
			encIdx:    1,
			totalSize: 6, // length + costPerECI + prevSize
			isFnc1:    false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			e := newInputEdge(test.c, test.encset, test.idx, test.prev, test.fnc1)

			if e.c != test.cVal {
				t.Errorf("edge.c = %v wants %v", e.c, test.cVal)
			}
			if e.encoderIndex != test.encIdx {
				t.Errorf("edge.encoderIndex = %v wants %v", e.encoderIndex, test.encIdx)
			}
			if e.previous != test.prev {
				t.Errorf("edge.previous = %v wants %v", e.previous, test.prev)
			}
			if e.cachedTotalSize != test.totalSize {
				t.Errorf("edge.cachedTotalSize = %v wants %v", e.cachedTotalSize, test.totalSize)
			}
			if f := e.isFNC1(); f != test.isFnc1 {
				t.Errorf("edge.isFnc1 = %v wants %v", f, test.isFnc1)
			}
		})
	}
}
