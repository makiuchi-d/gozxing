package common

import (
	"fmt"
)

type BitSource struct {
	bytes      []byte
	byteOffset uint
	bitOffset  uint
}

func NewBitSource(bytes []byte) *BitSource {
	return &BitSource{
		bytes: bytes,
	}
}

func (this *BitSource) GetBitOffset() uint {
	return this.bitOffset
}

func (this *BitSource) GetByteOffset() uint {
	return this.byteOffset
}

func (this *BitSource) ReadBit(numBits uint) (uint, error) {
	if numBits < 1 || numBits > 32 || numBits > this.Available() {
		return 0, fmt.Errorf("IllegalArgumentException: %v", numBits)
	}

	result := uint(0)

	// First, read remainder from current byte
	if this.bitOffset > 0 {
		bitsLeft := uint(8 - this.bitOffset)
		toRead := uint(bitsLeft)
		if numBits < bitsLeft {
			toRead = numBits
		}
		bitsToNotRead := bitsLeft - toRead
		mask := byte((0xFF >> (8 - toRead)) << bitsToNotRead)
		result = uint(this.bytes[this.byteOffset]&mask) >> bitsToNotRead
		numBits -= toRead
		this.bitOffset += toRead
		if this.bitOffset == 8 {
			this.bitOffset = 0
			this.byteOffset++
		}
	}

	// Next read whole bytes
	if numBits > 0 {
		for numBits >= 8 {
			result = (result << 8) | uint(this.bytes[this.byteOffset]&0xFF)
			this.byteOffset++
			numBits -= 8
		}

		// Finally read a partial byte
		if numBits > 0 {
			bitsToNotRead := 8 - numBits
			mask := byte((0xFF >> bitsToNotRead) << bitsToNotRead)
			result = (result << numBits) | uint((this.bytes[this.byteOffset]&mask)>>bitsToNotRead)
			this.bitOffset += numBits
		}
	}

	return result, nil
}

func (this *BitSource) Available() uint {
	return 8*(uint(len(this.bytes))-this.byteOffset) - this.bitOffset
}
