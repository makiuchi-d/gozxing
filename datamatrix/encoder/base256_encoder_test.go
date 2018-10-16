package encoder

import (
	"testing"
)

func TestBase256Encoder_getEncodingMode(t *testing.T) {
	enc := NewBase256Encoder()
	mode := enc.getEncodingMode()
	if mode != HighLevelEncoder_BASE256_ENCODATION {
		t.Fatalf("getEncodingMode = %v, expect %v", mode, HighLevelEncoder_BASE256_ENCODATION)
	}
}

func TestRandomize255State(t *testing.T) {
	ch := byte(10)
	pos := 0
	expect := byte(11)
	if r := base256Randomize255State(ch, pos); r != expect {
		t.Fatalf("randomize255State(%v,%v) = %v, expect %v", ch, pos, r, expect)
	}

	ch = byte(105)
	pos = 1
	expect = byte(255)
	if r := base256Randomize255State(ch, pos); r != expect {
		t.Fatalf("randomize255State(%v,%v) = %v, expect %v", ch, pos, r, expect)
	}

	ch = byte(25)
	pos = 5
	expect = byte(5)
	if r := base256Randomize255State(ch, pos); r != expect {
		t.Fatalf("randomize255State(%v,%v) = %v, expect %v", ch, pos, r, expect)
	}
}

func unrandom255State(ch byte, pos int) byte {
	pr := ((149 * pos) % 255) + 1
	v := int(ch) - pr
	if v < 0 {
		v += 256
	}
	return byte(v)
}

func TestBase256Encoder_encode(t *testing.T) {
	enc := NewBase256Encoder()

	ctx, _ := NewEncoderContext("")
	ctx.msg = make([]byte, 1560)
	e := enc.encode(ctx)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	ctx, _ = NewEncoderContext("")
	ctx.SignalEncoderChange(HighLevelEncoder_BASE256_ENCODATION)
	ctx.msg = make([]byte, 1556)
	e = enc.encode(ctx)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	ctx, _ = NewEncoderContext("")
	ctx.SignalEncoderChange(HighLevelEncoder_BASE256_ENCODATION)
	ctx.msg = make([]byte, 250)
	for i := range ctx.msg {
		ctx.msg[i] = 'A'
	}
	e = enc.encode(ctx)
	if r := ctx.GetCodewordCount(); r != 252 {
		t.Fatalf("encode codeword count = %v, expect 252", r)
	}
	if r := unrandom255State(ctx.codewords[0], 1); r != 250 {
		t.Fatalf("encode unrandom([0]) = %v, expect 250", r)
	}
	if r := unrandom255State(ctx.codewords[1], 2); r != 0 {
		t.Fatalf("encode unrandom([1]) = %v, expect 0", r)
	}
	for i := 2; i < len(ctx.codewords); i++ {
		if r := unrandom255State(ctx.codewords[i], i+1); r != 'A' {
			t.Fatalf("encode urandom([%d]) = %v, expect %v", i, r, 'A')
		}
	}

	ctx, _ = NewEncoderContext("")
	ctx.SignalEncoderChange(HighLevelEncoder_BASE256_ENCODATION)
	ctx.msg = make([]byte, 250)
	for i := 0; i < 100; i++ {
		ctx.msg[i] = 'A'
	}
	for i := 0; i < 10; i++ { // reset to ascii encodation
		ctx.msg[100+i] = '0'
	}
	e = enc.encode(ctx)
	if r := ctx.GetNewEncoding(); r != HighLevelEncoder_ASCII_ENCODATION {
		t.Fatalf("encode newEncoding = %v, expect %v", r, HighLevelEncoder_ASCII_ENCODATION)
	}
	if r := ctx.GetCodewordCount(); r != 101 {
		t.Fatalf("encode codeword count = %v, expect 101", r)
	}
	if r := unrandom255State(ctx.codewords[0], 1); r != 100 {
		t.Fatalf("encode unrandom([0]) = %v, expect 100", r)
	}
	for i := 1; i < len(ctx.codewords); i++ {
		if r := unrandom255State(ctx.codewords[i], i+1); r != 'A' {
			t.Fatalf("encode urandom([%d]) = %v, expect %v", i, r, 'A')
		}
	}
}
