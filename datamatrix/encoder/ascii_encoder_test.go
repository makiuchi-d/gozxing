package encoder

import (
	"reflect"
	"testing"
)

func TestASCIIEncoder_getEncodingMode(t *testing.T) {
	enc := NewASCIIEncoder()
	if r := enc.getEncodingMode(); r != HighLevelEncoder_ASCII_ENCODATION {
		t.Fatalf("getEncodingMode = %v, expect %v", r, HighLevelEncoder_ASCII_ENCODATION)
	}
}

func TestEncodeASCIIDigits(t *testing.T) {
	_, e := encodeASCIIDigits('0', 'a')
	if e == nil {
		t.Fatalf("encodeASCIIDigits must be error")
	}

	r, e := encodeASCIIDigits('9', '1')
	if e != nil {
		t.Fatalf("encodeASCIIDigits returns error: %v", e)
	}
	if r != 221 {
		t.Fatalf("encodeASCIIDigits = %v, expect 139", r)
	}
}

func TestASCIIEncoder_encode(t *testing.T) {
	enc := NewASCIIEncoder()

	ctx, _ := NewEncoderContext("12")
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode := HighLevelEncoder_ASCII_ENCODATION
	expect := []byte{142}
	e := enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}

	ctx, _ = NewEncoderContext("")
	ctx.msg = []byte{'A', 0xff, 0xff, 0xff, 0xff}
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode = HighLevelEncoder_BASE256_ENCODATION
	expect = []byte{HighLevelEncoder_LATCH_TO_BASE256}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}

	ctx, _ = NewEncoderContext(" 0A 0A")
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode = HighLevelEncoder_C40_ENCODATION
	expect = []byte{HighLevelEncoder_LATCH_TO_C40}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}

	ctx, _ = NewEncoderContext("*>\r*>\r*>\r")
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode = HighLevelEncoder_X12_ENCODATION
	expect = []byte{HighLevelEncoder_LATCH_TO_ANSIX12}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}

	ctx, _ = NewEncoderContext(" 0a 0a 0a")
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode = HighLevelEncoder_TEXT_ENCODATION
	expect = []byte{HighLevelEncoder_LATCH_TO_TEXT}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}

	ctx, _ = NewEncoderContext("^^^^^^^^^^^^")
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode = HighLevelEncoder_EDIFACT_ENCODATION
	expect = []byte{HighLevelEncoder_LATCH_TO_EDIFACT}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}

	ctx, _ = NewEncoderContext("")
	ctx.msg = []byte{128, '0', '0', '0', '0', '0', '0'}
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode = HighLevelEncoder_ASCII_ENCODATION
	expect = []byte{HighLevelEncoder_UPPER_SHIFT, 1}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}

	ctx, _ = NewEncoderContext("aa")
	ctx.SignalEncoderChange(HighLevelEncoder_ASCII_ENCODATION)
	newMode = HighLevelEncoder_ASCII_ENCODATION
	expect = []byte{0x62}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if r := ctx.GetNewEncoding(); r != newMode {
		t.Fatalf("encode newEncoding = %v, expect %v", r, newMode)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}
}
