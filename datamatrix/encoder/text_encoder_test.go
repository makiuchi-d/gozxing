package encoder

import (
	"reflect"
	"testing"
)

func TestTextEncoder_getEncodingMode(t *testing.T) {
	enc := NewTextEncoder()
	if r := enc.getEncodingMode(); r != HighLevelEncoder_TEXT_ENCODATION {
		t.Fatalf("getEncodingMode = %v, expect %v", r, HighLevelEncoder_TEXT_ENCODATION)
	}
}

func testTextEncodeChar(t testing.TB, c byte, expn int, expb []byte) {
	t.Helper()
	n, b := textEncodeChar(c, []byte{})
	if n != expn {
		t.Fatalf("textEncodeChar(%v) n = %v, expect %v", c, n, expn)
	}
	if !reflect.DeepEqual(b, expb) {
		t.Fatalf("textEncodeChar(%v) bs = %v, expect %v", c, b, expb)
	}
}

func TestTextEncodeChar(t *testing.T) {
	testTextEncodeChar(t, ' ', 1, []byte{3})
	testTextEncodeChar(t, '0', 1, []byte{4})
	testTextEncodeChar(t, '9', 1, []byte{13})
	testTextEncodeChar(t, 'a', 1, []byte{14})
	testTextEncodeChar(t, 'z', 1, []byte{39})
	testTextEncodeChar(t, '\n', 2, []byte{0, 10})
	testTextEncodeChar(t, '*', 2, []byte{1, 9})
	testTextEncodeChar(t, '=', 2, []byte{1, 18})
	testTextEncodeChar(t, '^', 2, []byte{1, 25})
	testTextEncodeChar(t, '`', 2, []byte{2, 0})
	testTextEncodeChar(t, 'A', 2, []byte{2, 1})
	testTextEncodeChar(t, 'Z', 2, []byte{2, 26})
	testTextEncodeChar(t, '~', 2, []byte{2, 30})
	testTextEncodeChar(t, 0xb0, 3, []byte{1, 0x1e, 4})
	testTextEncodeChar(t, 0xfe, 4, []byte{1, 0x1e, 2, 30})
}

func TestTextEncoder_encode(t *testing.T) {
	enc := NewTextEncoder()
	ctx, _ := NewEncoderContext("Zxing")
	ctx.SignalEncoderChange(HighLevelEncoder_TEXT_ENCODATION)

	e := enc.encode(ctx)
	if e != nil {
		t.Fatalf("encode returns error: %v", e)
	}
	if ctx.HasMoreCharacters() {
		t.Fatalf("context HasMoreCharacters = true, expect false")
	}
	expect := []byte{16, 182, 141, 205, 254}
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, expect) {
		t.Fatalf("context codewords = %v, expect %v", r, expect)
	}
}
