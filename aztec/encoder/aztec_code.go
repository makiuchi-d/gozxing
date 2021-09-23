package encoder

import (
	"github.com/makiuchi-d/gozxing"
)

// AztecCode : Aztec 2D code representation
//
// @author Rustam Abdullaev
//
type AztecCode struct {
	compact   bool
	size      int
	layers    int
	codeWords int
	matrix    *gozxing.BitMatrix
}

func newAztecCode() *AztecCode {
	return &AztecCode{}
}

// isCompact @return {@code true} if compact instead of full mode
//
func (this *AztecCode) isCompact() bool {
	return this.compact
}

func (this *AztecCode) setCompact(compact bool) {
	this.compact = compact
}

// getSize @return size in pixels (width and height)
//
func (this *AztecCode) getSize() int {
	return this.size
}

func (this *AztecCode) setSize(size int) {
	this.size = size
}

// getLayers @return number of levels
//
func (this *AztecCode) getLayers() int {
	return this.layers
}

func (this *AztecCode) setLayers(layers int) {
	this.layers = layers
}

// getCodeWords @return number of data codewords
//
func (this *AztecCode) getCodeWords() int {
	return this.codeWords
}

func (this *AztecCode) setCodeWords(codeWords int) {
	this.codeWords = codeWords
}

// getMatrix @return the symbol image
//
func (this *AztecCode) getMatrix() *gozxing.BitMatrix {
	return this.matrix
}

func (this *AztecCode) setMatrix(matrix *gozxing.BitMatrix) {
	this.matrix = matrix
}
