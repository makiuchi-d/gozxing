package encoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestSymbolInfo(t *testing.T) {
	s := NewSymbolInfo(false, 3, 5, 8, 8, 1)
	if r := s.GetSymbolDataWidth(); r != 8 {
		t.Fatalf("GetSymbolDataWidth = %v, expect 8", r)
	}
	if r := s.GetSymbolDataHeight(); r != 8 {
		t.Fatalf("GetSymbolDataHeight = %v, expect 8", r)
	}
	if r := s.GetSymbolWidth(); r != 10 {
		t.Fatalf("GetSymbolWidth = %v, expect 10", r)
	}
	if r := s.GetSymbolHeight(); r != 10 {
		t.Fatalf("GetSymbolHeight = %v, expect 10", r)
	}
	if r := s.GetCodewordCount(); r != 8 {
		t.Fatalf("GetCodewordCount = %v, expect 8", r)
	}
	if r := s.GetInterleavedBlockCount(); r != 1 {
		t.Fatalf("GetInterleavedBlockCount = %v, expect 1", r)
	}
	if r := s.GetDataCapacity(); r != 3 {
		t.Fatalf("GetDataCapacity = %v, expect 3", r)
	}
	if r := s.GetErrorCodewords(); r != 5 {
		t.Fatalf("GetErrorCodewords = %v, expect 5", r)
	}
	if r := s.GetMatrixWidth(); r != 8 {
		t.Fatalf("GetMatrixWidth = %v, expect 8", r)
	}
	if r := s.GetMatrixHeight(); r != 8 {
		t.Fatalf("GetMatrixHeight = %v, expect 8", r)
	}
	if r := s.GetDataLengthForInterleavedBlock(0); r != 3 {
		t.Fatalf("GetDataLengthForInterleavedBlock = %v, expect 3", r)
	}
	if r := s.GetErrorLengthForInterleavedBlock(0); r != 5 {
		t.Fatalf("GetErrorLengthForInterleavedBlock = %v, expect 5", r)
	}

	s = NewSymbolInfo(true, 5, 7, 16, 6, 1)
	if r := s.GetSymbolDataWidth(); r != 16 {
		t.Fatalf("GetSymbolDataWidth = %v, expect 16", r)
	}
	if r := s.GetSymbolDataHeight(); r != 6 {
		t.Fatalf("GetSymbolDataHeight = %v, expect 6", r)
	}
	if r := s.GetSymbolWidth(); r != 18 {
		t.Fatalf("GetSymbolWidth = %v, expect 18", r)
	}
	if r := s.GetSymbolHeight(); r != 8 {
		t.Fatalf("GetSymbolHeight = %v, expect 8", r)
	}
	if r := s.GetCodewordCount(); r != 12 {
		t.Fatalf("GetCodewordCount = %v, expect 12", r)
	}
	if r := s.GetInterleavedBlockCount(); r != 1 {
		t.Fatalf("GetInterleavedBlockCount = %v, expect 1", r)
	}
	if r := s.GetDataCapacity(); r != 5 {
		t.Fatalf("GetDataCapacity = %v, expect 5", r)
	}
	if r := s.GetErrorCodewords(); r != 7 {
		t.Fatalf("GetErrorCodewords = %v, expect 7", r)
	}
	if r := s.GetMatrixWidth(); r != 16 {
		t.Fatalf("GetMatrixWidth = %v, expect 16", r)
	}
	if r := s.GetMatrixHeight(); r != 6 {
		t.Fatalf("GetMatrixHeight = %v, expect 6", r)
	}
	if r := s.GetDataLengthForInterleavedBlock(0); r != 5 {
		t.Fatalf("GetDataLengthForInterleavedBlock = %v, expect 5", r)
	}
	if r := s.GetErrorLengthForInterleavedBlock(0); r != 7 {
		t.Fatalf("GetErrorLengthForInterleavedBlock = %v, expect 7", r)
	}

	s = NewSymbolInfoRS(false, 204, 84, 24, 24, 4, 102, 42)
	if r := s.GetSymbolDataWidth(); r != 48 {
		t.Fatalf("GetSymbolDataWidth = %v, expect 48", r)
	}
	if r := s.GetSymbolDataHeight(); r != 48 {
		t.Fatalf("GetSymbolDataHeight = %v, expect 48", r)
	}
	if r := s.GetSymbolWidth(); r != 52 {
		t.Fatalf("GetSymbolWidth = %v, expect 52", r)
	}
	if r := s.GetSymbolHeight(); r != 52 {
		t.Fatalf("GetSymbolHeight = %v, expect 52", r)
	}
	if r := s.GetCodewordCount(); r != 288 {
		t.Fatalf("GetCodewordCount = %v, expect 288", r)
	}
	if r := s.GetInterleavedBlockCount(); r != 2 {
		t.Fatalf("GetInterleavedBlockCount = %v, expect 2", r)
	}
	if r := s.GetDataCapacity(); r != 204 {
		t.Fatalf("GetDataCapacity = %v, expect 204", r)
	}
	if r := s.GetErrorCodewords(); r != 84 {
		t.Fatalf("GetErrorCodewords = %v, expect 84", r)
	}
	if r := s.GetMatrixWidth(); r != 24 {
		t.Fatalf("GetMatrixWidth = %v, expect 24", r)
	}
	if r := s.GetMatrixHeight(); r != 24 {
		t.Fatalf("GetMatrixHeight = %v, expect 24", r)
	}
	if r := s.GetDataLengthForInterleavedBlock(0); r != 102 {
		t.Fatalf("GetDataLengthForInterleavedBlock = %v, expect 102", r)
	}
	if r := s.GetErrorLengthForInterleavedBlock(0); r != 42 {
		t.Fatalf("GetErrorLengthForInterleavedBlock = %v, expect 42", r)
	}
}

func TestSymbolInfo_Lookup(t *testing.T) {
	_, e := SymbolInfo_Lookup(1559, SymbolShapeHint_FORCE_NONE, nil, nil, true)
	if e == nil {
		t.Fatalf("SymbolInfo_Lookup must be error")
	}
	r, e := SymbolInfo_Lookup(1559, SymbolShapeHint_FORCE_NONE, nil, nil, false)
	if r != nil || e != nil {
		t.Fatalf("SymbolInfo_Lookup = (%v, %v), expect (nil, nil)", r, e)
	}

	r, e = SymbolInfo_Lookup(10, SymbolShapeHint_FORCE_SQUARE, nil, nil, true)
	if e != nil {
		t.Fatalf("SymbolInfo_Lookup returns error: %v", e)
	}
	if r == nil {
		t.Fatalf("SymbolInfo_Lookup returns nil")
	}
	if r.rectangular {
		t.Fatalf("SymbolInfo_Lookup shape must be square, %v", r)
	}
	if r.dataCapacity != 12 {
		t.Fatalf("SymbolInfo_Lookup data capacity != 12, %v", r)
	}

	r, e = SymbolInfo_Lookup(5, SymbolShapeHint_FORCE_RECTANGLE, nil, nil, true)
	if e != nil {
		t.Fatalf("SymbolInfo_Lookup returns error: %v", e)
	}
	if r == nil {
		t.Fatalf("SymbolInfo_Lookup returns nil")
	}
	if !r.rectangular {
		t.Fatalf("SymbolInfo_Lookup shape must be rectangular, %v", r)
	}
	if r.dataCapacity != 5 {
		t.Fatalf("SymbolInfo_Lookup data capacity != 5, %v", r)
	}
	if w, h := r.GetSymbolDataWidth(), r.GetSymbolDataHeight(); w != 16 || h != 6 {
		t.Fatalf("SymbolInfo_Lookup symbol datasize = %vx%v, expect 16x6", w, h)
	}

	minSize, _ := gozxing.NewDimension(32, 10)
	r, e = SymbolInfo_Lookup(10, SymbolShapeHint_FORCE_NONE, minSize, nil, true)
	if e != nil {
		t.Fatalf("SymbolInfo_Lookup returns error: %v", e)
	}
	if r == nil {
		t.Fatalf("SymbolInfo_Lookup returns nil")
	}
	if w, h := r.GetSymbolWidth(), r.GetSymbolHeight(); w != 36 || h != 12 {
		t.Fatalf("SymbolInfo_Lookup symbol size = %vx%v, expect 36x12", w, h)
	}

	maxSize, _ := gozxing.NewDimension(30, 20)
	r, e = SymbolInfo_Lookup(10, SymbolShapeHint_FORCE_NONE, nil, maxSize, true)
	if e != nil {
		t.Fatalf("SymbolInfo_Lookup returns error: %v", e)
	}
	if r == nil {
		t.Fatalf("SymbolInfo_Lookup returns nil")
	}
	if w, h := r.GetSymbolWidth(), r.GetSymbolHeight(); w != 16 || h != 16 {
		t.Fatalf("SymbolInfo_Lookup symbol size = %vx%v, expect 16x16", w, h)
	}
}

func TestSymbolInfo_getHorizontalDataRegions(t *testing.T) {
	s := NewSymbolInfo(false, 3, 5, 8, 8, 1)
	if r := s.getHorizontalDataRegions(); r != 1 {
		t.Fatalf("getHorizontalDataRegions = %v, expect 1", r)
	}

	s.dataRegions = 2
	if r := s.getHorizontalDataRegions(); r != 2 {
		t.Fatalf("getHorizontalDataRegions = %v, expect 2", r)
	}

	s.dataRegions = 4
	if r := s.getHorizontalDataRegions(); r != 2 {
		t.Fatalf("getHorizontalDataRegions = %v, expect 2", r)
	}

	s.dataRegions = 16
	if r := s.getHorizontalDataRegions(); r != 4 {
		t.Fatalf("getHorizontalDataRegions = %v, expect 4", r)
	}

	s.dataRegions = 36
	if r := s.getHorizontalDataRegions(); r != 6 {
		t.Fatalf("getHorizontalDataRegions = %v, expect 6", r)
	}

	s.dataRegions = 0
	if r := s.getHorizontalDataRegions(); r != 0 {
		t.Fatalf("getHorizontalDataRegions = %v, expect 0", r)
	}

	for _, s := range symbols {
		r := s.getHorizontalDataRegions()
		if r == 0 {
			t.Fatalf("getHorizontalDataRegions == 0, SymbolInfo=%v", s)
		}
	}
}

func TestSymbolInfo_getVerticalDataRegions(t *testing.T) {
	s := NewSymbolInfo(false, 3, 5, 8, 8, 1)
	if r := s.getVerticalDataRegions(); r != 1 {
		t.Fatalf("getVerticalDataRegions = %v, expect 1", r)
	}

	s.dataRegions = 2
	if r := s.getVerticalDataRegions(); r != 1 {
		t.Fatalf("getVerticalDataRegions = %v, expect 1", r)
	}

	s.dataRegions = 4
	if r := s.getVerticalDataRegions(); r != 2 {
		t.Fatalf("getVerticalDataRegions = %v, expect 2", r)
	}

	s.dataRegions = 36
	if r := s.getVerticalDataRegions(); r != 6 {
		t.Fatalf("getVerticalDataRegions = %v, expect 6", r)
	}

	s.dataRegions = 5
	if r := s.getVerticalDataRegions(); r != 0 {
		t.Fatalf("getVerticalDataRegions = %v, expect 0", r)
	}

	for _, s := range symbols {
		r := s.getVerticalDataRegions()
		if r == 0 {
			t.Fatalf("getVerticalDataRegions == 0, SymbolInfo=%v", s)
		}
	}
}

func TestSymbolInfo_String(t *testing.T) {
	s := NewSymbolInfo(false, 3, 5, 8, 8, 1)
	expect := "Square Symbpl: data region 8x8, symbol size 10x10, symbol data size 8x8, codewords 3+5"
	if str := s.String(); str != expect {
		t.Fatalf("String = \"%v\", expect \"%v\"", str, expect)
	}

	s = NewSymbolInfo(true, 5, 7, 16, 6, 1)
	expect = "Rectangular Symbpl: data region 16x6, symbol size 18x8, symbol data size 16x6, codewords 5+7"
	if str := s.String(); str != expect {
		t.Fatalf("String = \"%v\", expect \"%v\"", str, expect)
	}
}
