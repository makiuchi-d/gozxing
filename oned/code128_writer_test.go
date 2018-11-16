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
	ctype := code128FindCType("abcd", 4)
	if ctype != code128CType_UNCODABLE {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_UNCODABLE)
	}

	ctype = code128FindCType("\xf1", 0)
	if ctype != code128CType_FNC_1 {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_FNC_1)
	}

	ctype = code128FindCType("abc", 0)
	if ctype != code128CType_UNCODABLE {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_UNCODABLE)
	}

	ctype = code128FindCType("abc0", 3)
	if ctype != code128CType_ONE_DIGIT {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_ONE_DIGIT)
	}

	ctype = code128FindCType("abc0a", 3)
	if ctype != code128CType_ONE_DIGIT {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_ONE_DIGIT)
	}

	ctype = code128FindCType("abc00", 3)
	if ctype != code128CType_TWO_DIGITS {
		t.Fatalf("code128FindCType = %v, expect %v", ctype, code128CType_TWO_DIGITS)
	}
}

func TestCode128ChooseCode(t *testing.T) {
	r := code128ChooseCode("0a", 0, code128CODE_CODE_A)
	expect := code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("0a", 0, code128CODE_CODE_C)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("\x00", 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("A", 0, code128CODE_CODE_A)
	expect = code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("A", 0, code128CODE_CODE_C)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("00", 0, code128CODE_CODE_C)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("\xf1", 0, code128CODE_CODE_A)
	expect = code128CODE_CODE_A
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("\xf1", 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("012a", 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("01\xf123", 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("01\xf10a", 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("0123456a", 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_B
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("01234567a", 0, code128CODE_CODE_B)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("00", 0, 0)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("\xf100", 0, 0)
	expect = code128CODE_CODE_C
	if r != expect {
		t.Fatalf("chooseCode = %v, expect %v", r, expect)
	}

	r = code128ChooseCode("\xf10a", 0, 0)
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
	r, e := enc.encode("\n\xf1\xf2\xf3\xf41\n")
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
	r, e = enc.encode("a\xf1\xf2\xf3\xf4a")
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
