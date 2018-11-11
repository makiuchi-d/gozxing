package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestCode39ToIntArray(t *testing.T) {
	arr := make([]int, 9)
	code39ToIntArray(0x034, arr)
	expect := []int{1, 1, 1, 2, 2, 1, 2, 1, 1}
	if !reflect.DeepEqual(arr, expect) {
		t.Fatalf("code39ToIntArray = %v, expect %v", arr, expect)
	}
}

func TestCode39TryToConvertToExtendedMode(t *testing.T) {
	_, e := code39TryToConvertToExtendedMode("\x80")
	if e == nil {
		t.Fatalf("code39TryToConvertToExtendedMode must be error")
	}

	r, e := code39TryToConvertToExtendedMode("\x00 @`\n\x1b!:1;A[a{\x7f")
	expect := "%U %V%W$J%A/A/Z1%FA%K+A%P%T"
	if e != nil {
		t.Fatalf("code39TryToConvertToExtendedMode returns error: %v", e)
	}
	if r != expect {
		t.Fatalf("code39TryToConvertToExtendedMode = \"%v\", expect \"%v\"", r, expect)
	}
}

func TestCode39Encoder_encode(t *testing.T) {
	enc := code39Encoder{}

	_, e := enc.encode(
		"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-. $/+%0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-.")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("\x80")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	_, e = enc.encode("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@.")
	if e == nil {
		t.Fatalf("encode must be error")
	}

	r, e := enc.encode("1!")
	expect := []bool{
		true, false, false, true, false, true, true, false, true, true, false, true, false, // *
		true, true, false, true, false, false, true, false, true, false, true, true, false, // 1
		true, false, false, true, false, false, true, false, true, false, false, true, false, // /
		true, true, false, true, false, true, false, false, true, false, true, true, false, // A
		true, false, false, true, false, true, true, false, true, true, false, true, // *
	}
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("encoded: %v\nexpect: %v", r, expect)
	}
}

func TestCode39Writer(t *testing.T) {
	writer := NewCode39Writer()
	format := gozxing.BarcodeFormat_CODE_39

	testEncode(t, writer, format,
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		"000001001011011010110101001011010110100101101101101001010101011001011011010110010101"+
			"011011001010101010011011011010100110101011010011010101011001101011010101001101011010"+
			"100110110110101001010101101001101101011010010101101101001010101011001101101010110010"+
			"101101011001010101101100101100101010110100110101011011001101010101001011010110110010"+
			"110101010011011010101010011011010110100101011010110010101101101100101010101001101011"+
			"01101001101010101100110101010100101101101101001011010101100101101010010110110100000")
	testEncode(t, writer, format,
		"\u0000\u0001\u0002\u0003\u0004\u0005\u0006\u0007\b\t\n\u000b\f\r\u000e\u000f\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001a\u001b\u001c\u001d\u001e\u001f",
		"000001001011011010101001001001011001010101101001001001010110101001011010010010010101"+
			"011010010110100100100101011011010010101001001001010101011001011010010010010101101011"+
			"001010100100100101010110110010101001001001010101010011011010010010010101101010011010"+
			"100100100101010110100110101001001001010101011001101010010010010101101010100110100100"+
			"100101010110101001101001001001010110110101001010010010010101010110100110100100100101"+
			"011010110100101001001001010101101101001010010010010101010101100110100100100101011010"+
			"101100101001001001010101101011001010010010010101010110110010100100100101011001010101"+
			"101001001001010100110101011010010010010101100110101010100100100101010010110101101001"+
			"001001010110010110101010010010010101001101101010101001001001011010100101101010010010"+
			"010101101001011010100100100101101101001010101001001001010101100101101010010010010110"+
			"101100101010010110110100000")
}
