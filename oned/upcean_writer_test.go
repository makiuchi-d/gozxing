package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestNewUPCEANWriter(t *testing.T) {
	writer := NewUPCEANWriter(dummyEncoder{}, gozxing.BarcodeFormat_UPC_A)
	if writer.defaultMargin != 9 {
		t.Fatalf("defaultMargine = %v, expect 9", writer.defaultMargin)
	}
}
