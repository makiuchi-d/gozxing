package encoder

import (
	"reflect"
	"testing"
)

func TestX12Encoder_getEncodingMode(t *testing.T) {
	enc := NewX12Encoder()
	if r := enc.getEncodingMode(); r != HighLevelEncoder_X12_ENCODATION {
		t.Fatalf("getEncodingMode != %v, expect %v", r, HighLevelEncoder_X12_ENCODATION)
	}
}

func testX12EncodeChar(t testing.TB, c, expect byte) {
	t.Helper()
	sb, e := x12EncodeChar(c, []byte{})
	if e != nil {
		t.Fatalf("x12EncodeChar(%v) returns error: %v", c, e)
	}
	if !reflect.DeepEqual(sb, []byte{expect}) {
		t.Fatalf("x12EncodeChar(%v) sb = %v, expect [%v]", c, sb, expect)
	}
}

func TestX12EncodeChar(t *testing.T) {
	_, e := x12EncodeChar('^', []byte{})
	if e == nil {
		t.Fatalf("x12EncodeChar must be error")
	}

	testX12EncodeChar(t, '\r', 0)
	testX12EncodeChar(t, '*', 1)
	testX12EncodeChar(t, '>', 2)
	testX12EncodeChar(t, ' ', 3)
	testX12EncodeChar(t, '0', 4)
	testX12EncodeChar(t, '9', 13)
	testX12EncodeChar(t, 'A', 14)
	testX12EncodeChar(t, 'Z', 39)
}

func TestX12HandleEOD(t *testing.T) {
	ctx, _ := NewEncoderContext("")
	ctx.codewords = make([]byte, 1559)
	e := x12HandleEOD(ctx, []byte{})
	if e == nil {
		t.Fatalf("x12HandleEOD must be error")
	}

	ctx, _ = NewEncoderContext("AAA0")
	ctx.symbolInfo = symbols[10]
	ctx.pos = 4
	ctx.codewords = []byte{89, 191}
	e = x12HandleEOD(ctx, []byte{1})
	if e != nil {
		t.Fatalf("x12HandleEOD returns error: %v", e)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, []byte{89, 191, 254}) {
		t.Fatalf("context codewords = %v, expect [89 191 254]", r)
	}
	if ctx.pos != 3 {
		t.Fatalf("context pos = %v, expect 3", ctx.pos)
	}
}

func TestX12Encoder_encode(t *testing.T) {
	enc := NewX12Encoder()

	ctx, _ := NewEncoderContext("^")
	ctx.SignalEncoderChange(HighLevelEncoder_X12_ENCODATION)
	e := enc.encode(ctx)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	ctx, _ = NewEncoderContext("AAA000000")
	ctx.SignalEncoderChange(HighLevelEncoder_X12_ENCODATION)
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	expect := []byte{89, 191, 254}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("context codewords = %v, expect %v", r, expect)
	}
	if ctx.pos != 3 {
		t.Fatalf("context pos = %v, expect 3", ctx.pos)
	}
	if r := ctx.GetNewEncoding(); r != HighLevelEncoder_ASCII_ENCODATION {
		t.Fatalf("context nextEncoding = %v, expect %v", r, HighLevelEncoder_ASCII_ENCODATION)
	}
}
