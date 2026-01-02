package encoder

import (
	"github.com/makiuchi-d/gozxing"
)

// AztecCode Aztec 2D code representation
type AztecCode struct {
	compact  bool
	size     int
	layers   int
	codeWords int
	matrix   *gozxing.BitMatrix
}

func NewAztecCode() *AztecCode {
	return &AztecCode{}
}

// IsCompact returns true if compact instead of full mode
func (this *AztecCode) IsCompact() bool {
	return this.compact
}

func (this *AztecCode) SetCompact(compact bool) {
	this.compact = compact
}

// GetSize returns size in pixels (width and height)
func (this *AztecCode) GetSize() int {
	return this.size
}

func (this *AztecCode) SetSize(size int) {
	this.size = size
}

// GetLayers returns number of levels
func (this *AztecCode) GetLayers() int {
	return this.layers
}

func (this *AztecCode) SetLayers(layers int) {
	this.layers = layers
}

// GetCodeWords returns number of data codewords
func (this *AztecCode) GetCodeWords() int {
	return this.codeWords
}

func (this *AztecCode) SetCodeWords(codeWords int) {
	this.codeWords = codeWords
}

// GetMatrix returns the symbol image
func (this *AztecCode) GetMatrix() *gozxing.BitMatrix {
	return this.matrix
}

func (this *AztecCode) SetMatrix(matrix *gozxing.BitMatrix) {
	this.matrix = matrix
}

