package encoder

import (
	"reflect"
	"testing"
)

func TestEdifactEncoder_getEncodingMode(t *testing.T) {
	enc := NewEdifactEncoder()
	if r := enc.getEncodingMode(); r != HighLevelEncoder_EDIFACT_ENCODATION {
		t.Fatalf("getEncodingMode = %v, expect %v", r, HighLevelEncoder_EDIFACT_ENCODATION)
	}
}

func TestEdifactEncodeToCodewords(t *testing.T) {
	_, e := edifactEncodeToCodewords([]byte{})
	if e == nil {
		t.Fatalf("edifactEncodeToCodewords must be error")
	}

	sb, e := edifactEncodeToCodewords([]byte{10})
	expect := []byte{40}
	if e != nil {
		t.Fatalf("edifactEncodeToCodewords returns error: %v", e)
	}
	if !reflect.DeepEqual(sb, expect) {
		t.Fatalf("edifactEncodeToCodewords = %v, expect %v", sb, expect)
	}

	sb, e = edifactEncodeToCodewords([]byte{10, 20, 30, 40})
	expect = []byte{41, 71, 168}
	if e != nil {
		t.Fatalf("edifactEncodeToCodewords returns error: %v", e)
	}
	if !reflect.DeepEqual(sb, expect) {
		t.Fatalf("edifactEncodeToCodewords = %v, expect %v", sb, expect)
	}
}

func TestEdifactEncodeChar(t *testing.T) {
	sb := make([]byte, 0)

	_, e := edifactEncodeChar(31, sb)
	if e == nil {
		t.Fatalf("edifactEncodeChar must be error")
	}

	_, e = edifactEncodeChar(95, sb)
	if e == nil {
		t.Fatalf("edifactEncodeChar must be error")
	}

	r, e := edifactEncodeChar(32, sb)
	if e != nil {
		t.Fatalf("edifactEncodeChar returns error: %v", e)
	}
	if !reflect.DeepEqual(r, []byte{32}) {
		t.Fatalf("edifactEncodeChar = %v, expect [32]", r)
	}

	r, e = edifactEncodeChar(63, sb)
	if e != nil {
		t.Fatalf("edifactEncodeChar returns error: %v", e)
	}
	if !reflect.DeepEqual(r, []byte{63}) {
		t.Fatalf("edifactEncodeChar = %v, expect [63]", r)
	}

	r, e = edifactEncodeChar(64, sb)
	if e != nil {
		t.Fatalf("edifactEncodeChar returns error: %v", e)
	}
	if !reflect.DeepEqual(r, []byte{0}) {
		t.Fatalf("edifactEncodeChar = %v, expect [0]", r)
	}

	r, e = edifactEncodeChar(94, sb)
	if e != nil {
		t.Fatalf("edifactEncodeChar returns error: %v", e)
	}
	if !reflect.DeepEqual(r, []byte{30}) {
		t.Fatalf("edifactEncodeChar = %v, expect [30]", r)
	}
}

func TestEdifactHandleEOD(t *testing.T) {
	ctx, _ := NewEncoderContext("")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	e := edifactHandleEOD(ctx, []byte{})
	if e != nil {
		t.Fatalf("edifactHandleEOD returns error: %v", e)
	}
	if r := ctx.GetCodewordCount(); r != 0 {
		t.Fatalf("context codeword count = %v, expect 0", r)
	}

	ctx, _ = NewEncoderContext("")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	ctx.codewords = make([]byte, 1559)
	e = edifactHandleEOD(ctx, []byte{31})
	if e == nil {
		t.Fatalf("edifactHandleEOD must be error")
	}

	ctx, _ = NewEncoderContext("@")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	ctx.codewords = make([]byte, 1558)
	e = edifactHandleEOD(ctx, []byte{31})
	if e == nil {
		t.Fatalf("edifactHandleEOD must be error")
	}

	ctx, _ = NewEncoderContext("@@@")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	ctx.pos = 1
	ctx.codewords = make([]byte, 3)
	e = edifactHandleEOD(ctx, []byte{31})
	if e != nil {
		t.Fatalf("edifactHandleEOD returns error: %v", e)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, []byte{0, 0, 0}) {
		t.Fatalf("context codewords = %v, expect [0 0 0]", r)
	}

	ctx, _ = NewEncoderContext("AAAA")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	ctx.pos = 4
	e = edifactHandleEOD(ctx, []byte{1, 1, 1, 1, 31})
	if e == nil {
		t.Fatalf("edifactHandleEOD must be error")
	}

	ctx, _ = NewEncoderContext("AA")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	ctx.pos = 2
	ctx.codewords = make([]byte, 1557)
	e = edifactHandleEOD(ctx, []byte{1, 1, 31})
	if e == nil {
		t.Fatalf("edifactHandleEOD must be error")
	}

	ctx, _ = NewEncoderContext("??")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	ctx.pos = 2
	e = edifactHandleEOD(ctx, []byte{63, 63, 31})
	if e != nil {
		t.Fatalf("edifactHandleEOD returns error: %v", e)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, []byte{255, 247, 192}) {
		t.Fatalf("context codewords = %v, expect [255 247 192]", r)
	}

	ctx, _ = NewEncoderContext("??")
	ctx.SignalEncoderChange(HighLevelEncoder_EDIFACT_ENCODATION)
	ctx.pos = 2
	ctx.codewords = make([]byte, 1)
	e = edifactHandleEOD(ctx, []byte{63, 63, 31})
	if e != nil {
		t.Fatalf("edifactHandleEOD returns error: %v", e)
	}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, []byte{0}) {
		t.Fatalf("context codewords = %v, expect [0]", r)
	}
	if ctx.pos != 0 {
		t.Fatalf("context pos = %v, expect 0", ctx.pos)
	}
}

func TestEdifactEncoder_encode(t *testing.T) {
	enc := NewEdifactEncoder()

	ctx, _ := NewEncoderContext("~")
	e := enc.encode(ctx)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	ctx, _ = NewEncoderContext("AAAA00000000")
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	expect := []byte{4, 16, 65, 124}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("encode codewords = %v, expect %v", r, expect)
	}
	if ctx.pos != 4 {
		t.Fatalf("encode context.pos = %v, expect 4", ctx.pos)
	}
	if r := ctx.GetNewEncoding(); r != HighLevelEncoder_ASCII_ENCODATION {
		t.Fatalf("encode context.newEncoding = %v, expect %v", r, HighLevelEncoder_ASCII_ENCODATION)
	}
}
