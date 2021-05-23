package encoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestC40Encoder_getEncodingMode(t *testing.T) {
	enc := NewC40Encoder()
	if r := enc.getEncodingMode(); r != HighLevelEncoder_C40_ENCODATION {
		t.Fatalf("getEncodingMode %v, expect %v", r, HighLevelEncoder_C40_ENCODATION)
	}
}

func TestC40EncodeToCodewords(t *testing.T) {
	sb := []byte{3, 4, 5}
	r := c40EncodeToCodewords(sb)
	expect := []byte{19, 102}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("c40EncodeToCodewords = %v, expect %v", r, expect)
	}
}

func testC40EncodeChar(t testing.TB, c byte, expn int, expb []byte) {
	t.Helper()
	n, b := c40EncodeChar(c, make([]byte, 0))
	if n != expn {
		t.Fatalf("c40EncodeChar(%v) num=%v, expect %v", c, n, expn)
	}
	if !reflect.DeepEqual(b, expb) {
		t.Fatalf("c40EncodeChar(%v) buffer=%v, expect %v", c, b, expb)
	}
}

func TestC40EncodeChar(t *testing.T) {
	testC40EncodeChar(t, ' ', 1, []byte{3})
	testC40EncodeChar(t, '5', 1, []byte{9})
	testC40EncodeChar(t, 'K', 1, []byte{24})
	testC40EncodeChar(t, '\n', 2, []byte{0, 10})
	testC40EncodeChar(t, '%', 2, []byte{1, 4})
	testC40EncodeChar(t, '=', 2, []byte{1, 18})
	testC40EncodeChar(t, '^', 2, []byte{1, 25})
	testC40EncodeChar(t, 'i', 2, []byte{2, 9})
	testC40EncodeChar(t, 0xda, 3, []byte{1, 0x1e, 39})
	testC40EncodeChar(t, 0x8f, 4, []byte{1, 0x1e, 0, 0x0f})
}

func TestC40HandleEOD(t *testing.T) {

	// error on UpdateSymbolInfo
	ctx, _ := NewEncoderContext("")
	ctx.SignalEncoderChange(HighLevelEncoder_C40_ENCODATION)
	e := c40HandleEOD(ctx, make([]byte, 3000))
	if e == nil {
		t.Fatalf("c40HandleEOD must be error")
	}

	ctx, _ = NewEncoderContext("")
	ctx.SignalEncoderChange(HighLevelEncoder_C40_ENCODATION)
	e = c40HandleEOD(ctx, make([]byte, 10))
	if e == nil {
		t.Fatalf("c40HandleEOD must be error")
	}

	// lest = 2
	ctx, _ = NewEncoderContext("     000000")
	ctx.pos = 5
	ctx.SignalEncoderChange(HighLevelEncoder_C40_ENCODATION)
	e = c40HandleEOD(ctx, []byte{3, 3, 3, 3, 3})
	if e != nil {
		t.Fatalf("c40HandleEOD returns error: %v", e)
	}
	expect := []byte{19, 60, 19, 57, 254}
	if cws := ctx.GetCodewords(); !reflect.DeepEqual(cws, expect) {
		t.Fatalf("c40HandleEOD codewords=%v, expect=%v", cws, expect)
	}
	if ctx.pos != 5 {
		t.Fatalf("c40HandleEOD ctx.pos = %v, expect 6", ctx.pos)
	}

	// lest = 1, available = 1
	ctx, _ = NewEncoderContext("       000000")
	ctx.pos = 7
	ctx.SignalEncoderChange(HighLevelEncoder_C40_ENCODATION)
	e = c40HandleEOD(ctx, []byte{3, 3, 3, 3, 3, 3, 3})
	if e != nil {
		t.Fatalf("c40HandleEOD returns error: %v", e)
	}
	expect = []byte{19, 60, 19, 60, 254}
	if cws := ctx.GetCodewords(); !reflect.DeepEqual(cws, expect) {
		t.Fatalf("c40HandleEOD codewords=%v, expect=%v", cws, expect)
	}
	if ctx.pos != 6 {
		t.Fatalf("c40HandleEOD ctx.pos = %v, expect 6", ctx.pos)
	}

	// rest = 0
	ctx, _ = NewEncoderContext("      000000")
	ctx.pos = 6
	ctx.SignalEncoderChange(HighLevelEncoder_C40_ENCODATION)
	e = c40HandleEOD(ctx, []byte{3, 3, 3, 3, 3, 3})
	if e != nil {
		t.Fatalf("c40HandleEOD returns error: %v", e)
	}
	expect = []byte{19, 60, 19, 60, 254}
	if cws := ctx.GetCodewords(); !reflect.DeepEqual(cws, expect) {
		t.Fatalf("c40HandleEOD codewords=%v, expect=%v", cws, expect)
	}
	if ctx.pos != 6 {
		t.Fatalf("c40HandleEOD ctx.pos = %v, expect 6", ctx.pos)
	}
}

func TestC40Encoder_encode(t *testing.T) {
	enc := NewC40Encoder()

	ctx, _ := NewEncoderContext("0A 0A 0A 0A 0A")
	dim, _ := gozxing.NewDimension(10, 10)
	ctx.SetSizeConstraints(dim, dim)
	e := enc.encode(ctx)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	ctx, _ = NewEncoderContext("0A 0A 0A 0A 0A")
	ctx.SignalEncoderChange(HighLevelEncoder_C40_ENCODATION)
	expect := []byte{27, 52, 27, 52, 27, 52, 27, 52, 254}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if cws := ctx.GetCodewords(); !reflect.DeepEqual(cws, expect) {
		t.Fatalf("encode codewords=%v, expect=%v", cws, expect)
	}

	ctx, _ = NewEncoderContext(" A A A A A A000000")
	ctx.SignalEncoderChange(HighLevelEncoder_C40_ENCODATION)
	expect = []byte{20, 244, 88, 7, 20, 244, 88, 7, 254}
	e = enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if cws := ctx.GetCodewords(); !reflect.DeepEqual(cws, expect) {
		t.Fatalf("encode codewords=%v, expect=%v", cws, expect)
	}
	if ctx.pos != 12 {
		t.Fatalf("encode context.pos = %v, expect 12", ctx.pos)
	}
}
