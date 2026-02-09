package common

import (
	"reflect"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
)

func TestCanEncodeChar(t *testing.T) {
	latin1, _ := ianaindex.IANA.Encoding("ISO-8859-1")
	sjis, _ := ianaindex.IANA.Encoding("Shift_JIS")
	utf8, _ := ianaindex.IANA.Encoding("UTF-8")

	tests := map[string]struct {
		enc  encoding.Encoding
		char rune
		exp  bool
	}{
		"ISO-8859-1 ok": {latin1, 'À', true},
		"ISO-8859-1 ng": {latin1, 'あ', false},
		"Shift_JIS ok":  {sjis, 'あ', true},
		"Shift_JIS ng":  {sjis, 'À', false},
		"UTF-8 ok":      {utf8, '😀', true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := canEncodeChar(test.enc, test.char)

			if r != test.exp {
				t.Errorf("canEncodeChar = %v wants %v", r, test.exp)
			}
		})
	}
}

func TestNewECIEncoderSet(t *testing.T) {
	latin1 := charmap.ISO8859_1
	sjis := japanese.ShiftJIS
	utf8 := unicode.UTF8
	utf16be, _ := ianaindex.IANA.Encoding("UTF-16BE")

	tests := map[string]struct {
		str     []rune
		charset encoding.Encoding
		fnc1    int

		encs    []encoding.Encoding
		names   []string
		ecis    []int
		pencidx int
	}{
		"Latin1-only": {
			[]rune("abc"), nil, -1,
			[]encoding.Encoding{latin1},
			[]string{"ISO_8859-1:1987"},
			[]int{1},
			-1,
		},
		"ShiftJIS": {
			[]rune("あ"), sjis, -1,
			[]encoding.Encoding{latin1, sjis, utf8, utf16be},
			[]string{"ISO_8859-1:1987", "Shift_JIS", "UTF-8", "UTF-16BE"},
			[]int{1, 20, 26, 25},
			1,
		},
		"NeedUnicode": {
			[]rune("😀"), nil, 1,
			[]encoding.Encoding{latin1, utf8, utf16be},
			[]string{"ISO_8859-1:1987", "UTF-8", "UTF-16BE"},
			[]int{1, 26, 25},
			-1,
		},
		"FNC1": {
			[]rune("\u001dabc"), nil, 0x1d,
			[]encoding.Encoding{latin1},
			[]string{"ISO_8859-1:1987"},
			[]int{1},
			-1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			encset := NewECIEncoderSet(test.str, test.charset, test.fnc1)

			l := encset.Length()
			if l != len(test.encs) {
				t.Fatalf("Length = %v wants %v", l, len(test.encs))
			}
			for i := range l {
				if e := encset.GetCharset(i); e != test.encs[i] {
					t.Errorf("GetCharset(%d) = %v wants %v", i, e, test.encs[i])
				}
				if n := encset.GetCharsetName(i); n != test.names[i] {
					t.Errorf("GetCharsetName(%d) = %q wants %q", i, n, test.names[i])
				}
				if v := encset.GetECIValue(i); v != test.ecis[i] {
					t.Errorf("GetECIValue(%d) = %v wants %v", i, v, test.ecis[i])
				}
			}
			if p := encset.GetPriorityEncoderIndex(); p != test.pencidx {
				t.Errorf("GetPriorityEncoderIndex = %v wants %v", p, test.pencidx)
			}
		})
	}
}

func TestECIEncoderSet_MethodsFail(t *testing.T) {
	es := NewECIEncoderSet([]rune("abc"), nil, -1)

	if n := es.GetCharsetName(2); n != "" {
		t.Errorf("GetCharsetName must be empty: %q", n)
	}
	if e := es.GetCharset(2); e != nil {
		t.Errorf("Getcharset must be nil: %v", e)
	}
	if v := es.GetECIValue(2); v != -1 {
		t.Errorf("GetECIValue must be -1: %v", v)
	}
}

func TestECIEncoderSet_Encode(t *testing.T) {
	es := NewECIEncoderSet([]rune("abc"), nil, -1)
	str := "Àabc"
	expchar := []byte{0xc0}
	expstr := []byte{0xc0, 0x61, 0x62, 0x63}

	char := []rune(str)[0]
	if !es.CanEncode(char, 0) {
		t.Errorf("CanEncode must be true")
	}
	if es.CanEncode(char, 1) {
		t.Errorf("CanEncode must be false")
	}

	encchar, err := es.EncodeChar(char, 0)
	if err != nil {
		t.Errorf("EncodeChar error: %v", err)
	} else if !reflect.DeepEqual(encchar, expchar) {
		t.Errorf("EncodeChar(%c, 0) = %v wants %v", char, encchar, expchar)
	}
	_, err = es.EncodeChar(char, 1)
	if err == nil {
		t.Errorf("EncodeChar(%c, 1) must be error", char)
	}
	_, err = es.EncodeChar('あ', 0)
	if err == nil {
		t.Errorf("EncodeChar(あ, 0) must be error")
	}

	encstr, err := es.EncodeString(str, 0)
	if err != nil {
		t.Errorf("EncodeString(%q, 0) error: %v", str, err)
	} else if !reflect.DeepEqual(encstr, expstr) {
		t.Errorf("EncodeString(%q, 0) = %v wants %v", str, encstr, expstr)
	}
	_, err = es.EncodeString(str, 1)
	if err == nil {
		t.Errorf("EncodeString(%q, 1) must be error", str)
	}
	_, err = es.EncodeString("あ", 0)
	if err == nil {
		t.Errorf("EncodeString(\"あ\", 0) must be error")
	}
}
