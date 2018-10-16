package encoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestNewEncoderContxt(t *testing.T) {
	_, e := NewEncoderContext("あ")
	if e == nil {
		t.Fatalf("NewEncoderContext must be error")
	}

	ctx, e := NewEncoderContext("båd")
	if e != nil {
		t.Fatalf("NewEncoderContext returns error: %v", e)
	}
	expectmsg := []byte{0x62, 0xe5, 0x64}
	if msg := ctx.GetMessage(); !reflect.DeepEqual(msg, expectmsg) {
		t.Fatalf("NewEncoderContext msg = %v, expect %v", msg, expectmsg)
	}
	if c := ctx.GetCurrentChar(); c != 0x62 {
		t.Fatalf("GetCurrentChar = 0x%2x, expect 0x62", c)
	}
	ctx.pos++
	if c := ctx.GetCurrent(); c != 0xe5 {
		t.Fatalf("GetCurrentChar = 0x%2x, expect 0xe5", c)
	}
}

func TestEncoderContext_Codewords(t *testing.T) {
	ctx, _ := NewEncoderContext("abcdefg")

	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, []byte{}) {
		t.Fatalf("GetCodewords = %v, expect []", r)
	}

	ctx.WriteCodewords([]byte("abc"))
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, []byte("abc")) {
		t.Fatalf("GetCodewords = %v, expect [97 98 99]", r)
	}

	ctx.WriteCodeword('d')
	if r := ctx.GetCodewords(); !reflect.DeepEqual(r, []byte("abcd")) {
		t.Fatalf("GetCodewords = %v, expect [97 98 99 100]", r)
	}

	if r := ctx.GetCodewordCount(); r != 4 {
		t.Fatalf("GetCodewordCount = %v, expect 4", r)
	}
}

func TestEncoderContext_NewEncoding(t *testing.T) {
	ctx, _ := NewEncoderContext("abcdefg")

	if r := ctx.GetNewEncoding(); r != -1 {
		t.Fatalf("GetNewEncoding = %v, expect -1", r)
	}

	ctx.SignalEncoderChange(1)
	if r := ctx.GetNewEncoding(); r != 1 {
		t.Fatalf("GetNewEncoding = %v, expect 1", r)
	}

	ctx.ResetEncoderSignal()
	if r := ctx.GetNewEncoding(); r != -1 {
		t.Fatalf("GetNewEncoding = %v, expect -1", r)
	}
}

func TestEncoderContext_CharCount(t *testing.T) {
	ctx, _ := NewEncoderContext("abcdefg")

	if r := ctx.GetRemainingCharacters(); r != 7 {
		t.Fatalf("GetRemainingCharacters = %v, expect 7", r)
	}

	ctx.SetSkipAtEnd(2)
	if r := ctx.GetRemainingCharacters(); r != 5 {
		t.Fatalf("GetRemainingCharacters = %v, expect 5", r)
	}

	ctx.pos += 2
	if r := ctx.GetRemainingCharacters(); r != 3 {
		t.Fatalf("GetRemainingCharacters = %v, expect 5", r)
	}
	if !ctx.HasMoreCharacters() {
		t.Fatalf("HasMoreCharacters = false, expect true")
	}

	ctx.pos += 3
	if r := ctx.GetRemainingCharacters(); r != 0 {
		t.Fatalf("GetRemainingCharacters = %v, expect 0", r)
	}
	if ctx.HasMoreCharacters() {
		t.Fatalf("HasMoreCharacters = true, expect false")
	}
}

func TestEncoderContext_SymbolInfo(t *testing.T) {
	ctx, _ := NewEncoderContext("abcdefg")
	ctx.WriteCodewords([]byte("abcdefg"))

	if si := ctx.GetSymbolInfo(); si != nil {
		t.Fatalf("GetSymbolInfo must be nil, %v", si)
	}

	if e := ctx.UpdateSymbolInfo(); e != nil {
		t.Fatalf("UpdateSymbolInfo returns error: %v", e)
	}

	si := ctx.GetSymbolInfo()
	if si == nil {
		t.Fatalf("GetSymbolInfo = nil")
	}
	if w, h := si.GetSymbolWidth(), si.GetSymbolHeight(); w != 14 || h != 14 {
		t.Fatalf("GetSymbolInfo wxh = %vx%v, expect 14x14", w, h)
	}

	ctx.ResetSymbolInfo()
	if si := ctx.GetSymbolInfo(); si != nil {
		t.Fatalf("GetSymbolInfo = %v, expect nil", si)
	}

	ctx.SetSymbolShape(SymbolShapeHint_FORCE_RECTANGLE)
	if e := ctx.UpdateSymbolInfo(); e != nil {
		t.Fatalf("UpdateSymbolInfo returns error: %v", e)
	}
	si = ctx.GetSymbolInfo()
	if si == nil {
		t.Fatalf("GetSymbolInfo = nil")
	}
	if w, h := si.GetSymbolWidth(), si.GetSymbolHeight(); w != 32 || h != 8 {
		t.Fatalf("GetSymbolInfo wxh = %vx%v, expect 32x8", w, h)
	}

	ctx.ResetSymbolInfo()
	ctx.SetSymbolShape(SymbolShapeHint_FORCE_SQUARE)
	minSize, _ := gozxing.NewDimension(16, 16)
	ctx.SetSizeConstraints(minSize, nil)
	if e := ctx.UpdateSymbolInfo(); e != nil {
		t.Fatalf("UpdateSymbolInfo returns error: %v", e)
	}
	si = ctx.GetSymbolInfo()
	if si == nil {
		t.Fatalf("GetSymbolInfo = nil")
	}
	if w, h := si.GetSymbolWidth(), si.GetSymbolHeight(); w != 16 || h != 16 {
		t.Fatalf("GetSymbolInfo wxh = %vx%v, expect 16x16", w, h)
	}

	if e := ctx.UpdateSymbolInfoByLength(1559); e == nil {
		t.Fatalf("UpdateSymbolInfoByLength(1559) must be error")
	}
}
