package encoder

import (
	"github.com/makiuchi-d/gozxing"
)

type SimpleToken struct {
	tokenBase
	value    int
	bitCount int
}

func NewSimpleToken(previous Token, value, bitCount int) *SimpleToken {
	return &SimpleToken{
		tokenBase: tokenBase{previous: previous},
		value:    value,
		bitCount: bitCount,
	}
}

// Add overrides tokenBase.Add() to use 'this' (the SimpleToken) as the previous,
// not 'this.previous'. This matches the Java behavior where Token.add() uses 'this'.
func (this *SimpleToken) Add(value, bitCount int) Token {
	return NewSimpleToken(this, value, bitCount)
}

func (this *SimpleToken) AppendTo(bitArray *gozxing.BitArray, text []byte) {
	_ = bitArray.AppendBits(this.value, this.bitCount)
}

