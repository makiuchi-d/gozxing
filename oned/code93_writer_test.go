package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestCode93AppendPattern(t *testing.T) {
	target := make([]bool, 18)
	pos := 0
	pos += code93AppendPattern(target, pos, code93AsteriskEncoding)
	pos += code93AppendPattern(target, pos, 0x1A8) // A
	expect := []bool{
		true, false, true, false, true, true, true, true, false, // *
		true, true, false, true, false, true, false, false, false, // A
	}
	if pos != 18 {
		t.Fatalf("pos = %v, expect 18", pos)
	}
	if !reflect.DeepEqual(target, expect) {
		t.Fatalf("target = %v, expect %v", target, expect)
	}
}

func TestCode93ComputeChecksumIndex(t *testing.T) {
	contents := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	checksum := code93ComputeChecksumIndex(contents, 15)
	if checksum != 12 {
		t.Fatalf("checksum = %v, expect %v", checksum, 12)
	}
	checksum = code93ComputeChecksumIndex(contents, 20)
	if checksum != 0 {
		t.Fatalf("checksum = %v, expect %v", checksum, 0)
	}
}

func TestCode93ConvertToExtended(t *testing.T) {
	_, e := code93ConvertToExtended("\x80")
	if e == nil {
		t.Fatalf("code93ConvertToExtended must be error")
	}

	src := "\x00\x01\x1a\x1b\x1f $%+!,09:;@AZ[_`az{\x7f"
	expect := "bUaAaZbAbE $%+cAcL09cZbFbVAZbKbObWdAdZbPbT"
	r, e := code93ConvertToExtended(src)
	if e != nil {
		t.Fatalf("code93ConvertToExtended returns error: %v", e)
	}
	if r != expect {
		t.Fatalf("return = \"%v\", expect \"%v\"", r, expect)
	}
}

func TestCode93Encoder_encode(t *testing.T) {
	enc := code93Encoder{}

	_, e := enc.encode("\x80")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@1")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	r, e := enc.encode("!A")
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	// *(/)AA*
	expect := []bool{
		true, false, true, false, true, true, true, true, false, // *
		true, true, true, false, true, false, true, true, false, // (/)
		true, true, false, true, false, true, false, false, false, // A
		true, true, false, true, false, true, false, false, false, // A
		true, false, false, true, false, true, true, false, false, // checksum1 = 24
		true, false, false, true, true, false, true, false, false, // checksum2 = 19
		true, false, true, false, true, true, true, true, false, // *
		true,
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode = \n%v, expect \n%v", r, expect)
	}
}

func TestCode93Writer(t *testing.T) {
	writer := NewCode93Writer()
	format := gozxing.BarcodeFormat_CODE_93

	testEncode(t, writer, format,
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		"000001010111101101010001101001001101000101100101001100100101100010101011010001011001"+
			"001011000101001101001000110101010110001010011001010001101001011001000101101101101001"+
			"101100101101011001101001101100101101100110101011011001011001101001101101001110101000"+
			"101001010010001010001001010000101001010001001001001001000101010100001000100101000010"+
			"10100111010101000010101011110100000")
	testEncode(t, writer, format,
		"\x00\x01\x1a\x1b\x1f $%+!,09:;@AZ[_`az{\x7f",
		"00000"+"101011110"+ // *
			"111011010"+"110010110"+"100100110"+"110101000"+ // bU aA
			"100100110"+"100111010"+"111011010"+"110101000"+ // aZ bA
			"111011010"+"110010010"+"111010010"+"111001010"+ // bE space $
			"110101110"+"101110110"+"111010110"+"110101000"+ // % + cA
			"111010110"+"101011000"+"100010100"+"100001010"+ // cL 0 9
			"111010110"+"100111010"+"111011010"+"110001010"+ // cZ bF
			"111011010"+"110011010"+"110101000"+"100111010"+ // bV A Z
			"111011010"+"100011010"+"111011010"+"100101100"+ // bK bO
			"111011010"+"101101100"+"100110010"+"110101000"+ // bW dA
			"100110010"+"100111010"+"111011010"+"100010110"+ // dZ bP
			"111011010"+"110100110"+ // bT
			"110100010"+"110101100"+ // checksum: 12 28
			"101011110"+"100000") // *
}
