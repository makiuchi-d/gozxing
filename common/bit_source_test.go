package common

import (
	"testing"
)

func TestBitSource(t *testing.T) {
	b := NewBitSource([]byte{0xde, 0xad, 0xbe, 0xaf})

	if _, e := b.ReadBit(0); e == nil {
		t.Fatalf("ReadBit(0) must be error")
	}
	if _, e := b.ReadBit(33); e == nil {
		t.Fatalf("ReadBit(33) must be error")
	}

	r, e := b.ReadBit(4)
	if e != nil {
		t.Fatalf("ReadBit(4) returns error, %v", e)
	}
	if r != 0xd {
		t.Fatalf("ReadBit(4) returns 0x%02x, expect 0x0d", r)
	}
	if r := b.GetBitOffset(); r != 4 {
		t.Fatalf("GetBitOffset = %v, expect 4", r)
	}

	r, e = b.ReadBit(16)
	if e != nil {
		t.Fatalf("ReadBit(16) returns error, %v", e)
	}
	if r != 0xeadb {
		t.Fatalf("ReadBit(16) returns 0x%04x, expect 0xeadb", r)
	}
	if r := b.GetBitOffset(); r != 4 {
		t.Fatalf("GetBitOffset = %v, expect 4", r)
	}
	if r := b.GetByteOffset(); r != 2 {
		t.Fatalf("GetByteOffset = %v, expect 2", r)
	}

	r, e = b.ReadBit(3)
	if e != nil {
		t.Fatalf("ReadBit(3) returns error, %v", e)
	}
	if r != 7 {
		t.Fatalf("ReadBit(3) returns 0x%02x, expect 0x07", r)
	}
	if r := b.GetBitOffset(); r != 7 {
		t.Fatalf("GetBitOffset = %v, expect 7", r)
	}
	if r := b.GetByteOffset(); r != 2 {
		t.Fatalf("GetByteOffset = %v, expect 3", r)
	}
	if r := b.Available(); r != 9 {
		t.Fatalf("Available = %v, expect 9", r)
	}

	if _, e := b.ReadBit(10); e == nil {
		t.Fatalf("ReadBit(10) must be error")
	}
}
