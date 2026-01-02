package encoder

import (
	"github.com/makiuchi-d/gozxing"
)

type Token interface {
	GetPrevious() Token
	Add(value, bitCount int) Token
	AddBinaryShift(start, byteCount int) Token
	AppendTo(bitArray *gozxing.BitArray, text []byte)
}

type tokenBase struct {
	previous Token
}

func (this *tokenBase) GetPrevious() Token {
	return this.previous
}

// Add creates a new SimpleToken with the current token as the previous
// Note: In Go, when this method is called on an embedded tokenBase within
// a SimpleToken or BinaryShiftToken, we need to get the containing struct
// as the Token interface. However, since Go doesn't allow this directly,
// we need to handle it differently. The Java version uses 'this' directly.
// For now, we'll use a workaround: we need to pass the Token interface
// that contains this tokenBase. But we can't get that from tokenBase alone.
// 
// Actually, the real issue is that tokenBase.Add() should receive the
// containing Token, not use this.previous. Let's fix this properly.
func (this *tokenBase) Add(value, bitCount int) Token {
	// The Java version: return new SimpleToken(this, value, bitCount);
	// So 'this' (the current Token) becomes the previous of the new token
	// But in Go, 'this' is *tokenBase, not Token. We need the containing Token.
	// Since we can't get it from tokenBase, we need to change the approach.
	// 
	// Actually, wait - when called on a SimpleToken, the receiver is still
	// *tokenBase (the embedded field). So we can't get the SimpleToken from here.
	// 
	// The solution: SimpleToken and BinaryShiftToken should override Add()
	// to use themselves as the previous token.
	return NewSimpleToken(this, value, bitCount)
}

func (this *tokenBase) AddBinaryShift(start, byteCount int) Token {
	return NewBinaryShiftToken(this.previous, start, byteCount)
}

func (this *tokenBase) AppendTo(bitArray *gozxing.BitArray, text []byte) {
	// This should never be called on tokenBase directly
	panic("AppendTo called on tokenBase")
}

var emptyToken = &tokenBase{previous: nil}

func GetEmptyToken() Token {
	return emptyToken
}

