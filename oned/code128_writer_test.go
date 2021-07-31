package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestCode128CType_String(t *testing.T) {
	if s := code128CType_UNCODABLE.String(); s != "UNCODABLE" {
		t.Fatalf("code128CType_UNCODABLE string=\"%s\", expect \"UNCODABLE\"", s)
	}
	if s := code128CType_ONE_DIGIT.String(); s != "ONE_DIGIT" {
		t.Fatalf("code128CType_ONE_DIGIT string=\"%s\", expect \"ONE_DIGIT\"", s)
	}
	if s := code128CType_TWO_DIGITS.String(); s != "TWO_DIGITS" {
		t.Fatalf("code128CType_TWO_DIGITS string=\"%s\", expect \"TWO_DIGITS\"", s)
	}
	if s := code128CType_FNC_1.String(); s != "FNC_1" {
		t.Fatalf("code128CType_FNC_1 string=\"%s\", expect \"FNC_1\"", s)
	}
	if s := code128CType(-1).String(); s != "" {
		t.Fatalf("code128CType(-1) string=\"%s\", expect \"\"", s)
	}
}

func TestCode128FindCType(t *testing.T) {
	ctype := code128FindCType([]rune("abcd"), 4)
	if ctype != code128CType_UNCODABLE {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_UNCODABLE)
	}

	ctype = code128FindCType([]rune("\u00f1"), 0)
	if ctype != code128CType_FNC_1 {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_FNC_1)
	}

	ctype = code128FindCType([]rune("abc"), 0)
	if ctype != code128CType_UNCODABLE {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_UNCODABLE)
	}

	ctype = code128FindCType([]rune("abc0"), 3)
	if ctype != code128CType_ONE_DIGIT {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_ONE_DIGIT)
	}

	ctype = code128FindCType([]rune("abc0a"), 3)
	if ctype != code128CType_ONE_DIGIT {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_ONE_DIGIT)
	}

	ctype = code128FindCType([]rune("abc00"), 3)
	if ctype != code128CType_TWO_DIGITS {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_TWO_DIGITS)
	}
}

func TestCode128ChooseCode(t *testing.T) {
	r := code128ChooseCode([]rune("0a"), 0, code128CODE_CODE_A)
	expect := code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("0a)"), 0, code128CODE_CODE_C)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("\x00"), 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("A"), 0, code128CODE_CODE_A)
	expect = code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("A"), 0, code128CODE_CODE_C)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("00"), 0, code128CODE_CODE_C)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("\u00f1"), 0, code128CODE_CODE_A)
	expect = code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("\u00f1"), 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("012a"), 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("01\u00f123"), 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("01\u00f10a"), 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("0123456a"), 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("01234567a"), 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("00"), 0, 0)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("\u00f100"), 0, 0)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode([]rune("\u00f10a"), 0, 0)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}
}

func TestCode128Encoder_encode(t *testing.T) {
	enc := code128Encoder{}

	_, e := enc.encode("")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("012345678901234567890123456789012345678901234567890123456789012345678901234567890")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("\x80")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	// FNCs and single number in codesetA with FNCs
	r, e := enc.encode("\n\u00f1\u00f2\u00f3\u00f41\n")
	expect := []bool{
		true, true, false, true, false, false, false, false, true, false, false, // StartA
		true, false, false, false, false, true, true, false, false, true, false, // LF
		true, true, true, true, false, true, false, true, true, true, false, // FNC1
		true, true, true, true, false, true, false, true, false, false, false, // FNC2
		true, false, true, true, true, true, false, false, false, true, false, // FNC3
		true, true, true, false, true, false, true, true, true, true, false, // FNC4A
		true, false, false, true, true, true, false, false, true, true, false, // 1
		true, false, false, false, false, true, true, false, false, true, false, // LF
		true, true, false, false, true, true, true, false, false, true, false, // Checksum 18
		true, true, false, false, false, true, true, true, false, true, false, true, true, // Stop
	}
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode = %v, expect %v", r, expect)
	}

	// CodesetB with FNCs
	r, e = enc.encode("a\u00f1\u00f2\u00f3\u00f4a")
	expect = []bool{
		true, true, false, true, false, false, true, false, false, false, false, // StartB
		true, false, false, true, false, true, true, false, false, false, false, // a
		true, true, true, true, false, true, false, true, true, true, false, // FNC1
		true, true, true, true, false, true, false, true, false, false, false, // FNC2
		true, false, true, true, true, true, false, false, false, true, false, // FNC3
		true, false, true, true, true, true, false, true, true, true, false, // FNC4B
		true, false, false, true, false, true, true, false, false, false, false, // a
		true, false, false, true, true, true, true, false, true, false, false, //Checksum 84
		true, true, false, false, false, true, true, true, false, true, false, true, true, // Stop
	}
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode = %v, expect %v", r, expect)
	}

	// Codeset C -> B -> A
	r, e = enc.encode("00a\n")
	expect = []bool{
		true, true, false, true, false, false, true, true, true, false, false, // StartC
		true, true, false, true, true, false, false, true, true, false, false, // 00
		true, false, true, true, true, true, false, true, true, true, false, // CodeB
		true, false, false, true, false, true, true, false, false, false, false, // a
		true, true, true, false, true, false, true, true, true, true, false, // CodeA
		true, false, false, false, false, true, true, false, false, true, false, // LF
		true, false, false, false, true, true, false, false, false, true, false, // Checksum 38
		true, true, false, false, false, true, true, true, false, true, false, true, true, // Stop
	}
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode = %v, expect %v", r, expect)
	}
}

func TestCode128Writer(t *testing.T) {
	writer := NewCode128Writer()
	format := gozxing.BarcodeFormat_CODE_128

	testEncode(t, writer, format,
		"Code128-012345",
		"00000"+"11010010000"+ // StartB
			"10001000110"+"10001111010"+"10000100110"+"10110010000"+ // Code
			"10011100110"+"11001110010"+"11101001100"+"10011011100"+ // 1 2 8 -
			"10111011110"+"11001101100"+"11101101110"+"10111011000"+ // CodeC 01 23 45
			"11011100010"+ // Checksum 52
			"1100011101011"+"00000") // Stop
}

func TestCode128WriterFnc3(t *testing.T) {
	enc := code128Encoder{}
	data := "\u00f3FOO"
	expect := []bool{
		true, true, false, true, false, false, true, false, false, false, false, // StartB
		true, false, true, true, true, true, false, false, false, true, false, // FNC3
		true, false, false, false, true, true, false, false, false, true, false, // F
		true, false, false, false, true, true, true, false, true, true, false, // O
		true, false, false, false, true, true, true, false, true, true, false, // O
		true, true, false, true, true, true, true, false, true, true, false, // Checksum
		true, true, false, false, false, true, true, true, false, true, false, true, true, // Stop
	}

	r, e := enc.encode(data)
	if e != nil {
		t.Fatalf("encode = %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode = %v, expect %v", r, expect)
	}
}

func TestCode128WriterForceCodeSetFailure(t *testing.T) {
	enc := code128Encoder{}

	tests := []struct {
		toEncode string
		codeset  string
	}{
		{"ASDFx0123", "A"},
		{"ASdf\00123", "B"},
		{"123a5678", "C"},
		{"123\u00f2a678", "C"},
		{"123456789", "C"},
		{"123456789", "D"},
	}
	for _, test := range tests {
		hints := map[gozxing.EncodeHintType]interface{}{
			gozxing.EncodeHintType_FORCE_CODE_SET: test.codeset,
		}
		_, e := enc.encodeWithHints(test.toEncode, hints)
		if e == nil {
			t.Fatalf("encode %q with force code %q must be error", test.toEncode, test.codeset)
		}
	}
}

func TestCode128WriterWithForceCodeSet(t *testing.T) {
	enc := code128Encoder{}

	tests := []struct {
		toEncode string
		codeset  string
		expected string
	}{
		{
			"AB123", "A",
			"11010000100" + "10100011000" + "10001011000" + "10011100110" + "11001110010" + "11001011100" + "11001000100" + "1100011101011",
		},
		{
			"1234", "B",
			"11010010000" + "10011100110" + "11001110010" + "11001011100" + "11001001110" + "11110010010" + "1100011101011",
		},
	}
	for _, test := range tests {
		hints := map[gozxing.EncodeHintType]interface{}{
			gozxing.EncodeHintType_FORCE_CODE_SET: test.codeset,
		}
		r, e := enc.encodeWithHints(test.toEncode, hints)
		if e != nil {
			t.Fatalf("encode(%q, %q): %v", test.toEncode, test.codeset, e)
		}
		if lr, le := len(r), len(test.expected); lr != le {
			t.Fatalf("encode(%q, %q): len(ret) = %v, wants %v", test.toEncode, test.codeset, lr, le)
		}
		for i := range r {
			ex := test.expected[i] == '1'
			if r[i] != ex {
				t.Fatalf("encode(%q, %q) [%v] = %v, wants %v", test.toEncode, test.codeset, i, r[i], ex)
			}
		}
	}
}
