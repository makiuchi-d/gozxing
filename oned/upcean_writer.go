package oned

import (
	"github.com/makiuchi-d/gozxing"
)

func NewUPCEANWriter(enc encoder, format gozxing.BarcodeFormat) *OneDimensionalCodeWriter {
	writer := NewOneDimensionalCodeWriter(enc, format)
	writer.defaultMargin = 9
	return writer
}
