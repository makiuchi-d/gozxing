package oned

import (
	"testing"
)

func TestNewUPCEANWriter(t *testing.T) {
	writer := NewUPCEANWriter(dummyEncoder{})
	if writer.defaultMargin != 9 {
		t.Fatalf("defaultMargine = %v, expect 9", writer.defaultMargin)
	}
}
