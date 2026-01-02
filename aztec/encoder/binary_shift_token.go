package encoder

import (
	"github.com/makiuchi-d/gozxing"
)

type BinaryShiftToken struct {
	tokenBase
	binaryShiftStart    int
	binaryShiftByteCount int
}

func NewBinaryShiftToken(previous Token, binaryShiftStart, binaryShiftByteCount int) *BinaryShiftToken {
	return &BinaryShiftToken{
		tokenBase:            tokenBase{previous: previous},
		binaryShiftStart:     binaryShiftStart,
		binaryShiftByteCount: binaryShiftByteCount,
	}
}

// Add overrides tokenBase.Add() to use 'this' (the BinaryShiftToken) as the previous,
// not 'this.previous'. This matches the Java behavior where Token.add() uses 'this'.
func (this *BinaryShiftToken) Add(value, bitCount int) Token {
	return NewSimpleToken(this, value, bitCount)
}

func (this *BinaryShiftToken) AppendTo(bitArray *gozxing.BitArray, text []byte) {
	bsbc := this.binaryShiftByteCount
	for i := 0; i < bsbc; i++ {
		if i == 0 || (i == 31 && bsbc <= 62) {
			// We need a header before the first character, and before
			// character 31 when the total byte code is <= 62
			_ = bitArray.AppendBits(31, 5) // BINARY_SHIFT
			if bsbc > 62 {
				_ = bitArray.AppendBits(bsbc-31, 16)
			} else if i == 0 {
				// 1 <= binaryShiftByteCode <= 62
				min := bsbc
				if min > 31 {
					min = 31
				}
				_ = bitArray.AppendBits(min, 5)
			} else {
				// 32 <= binaryShiftCount <= 62 and i == 31
				_ = bitArray.AppendBits(bsbc-31, 5)
			}
		}
		_ = bitArray.AppendBits(int(text[this.binaryShiftStart+i])&0xFF, 8)
	}
}

