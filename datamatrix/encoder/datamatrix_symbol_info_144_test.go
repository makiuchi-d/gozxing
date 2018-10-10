package encoder

import (
	"testing"
)

func TestDatamatrixSymbolInfo144(t *testing.T) {
	s := NewDataMatrixSymbolInfo144()

	if r := s.GetInterleavedBlockCount(); r != 10 {
		t.Fatalf("GetInterleavedBlockCount = %v, expect 10", r)
	}

	if r := s.GetDataLengthForInterleavedBlock(8); r != 156 {
		t.Fatalf("GetDataLengthForInterleavedBlock(8) = %v, expect 156", r)
	}

	if r := s.GetDataLengthForInterleavedBlock(9); r != 155 {
		t.Fatalf("GetDataLengthForInterleavedBlock(9) = %v, expect 155", r)
	}
}
