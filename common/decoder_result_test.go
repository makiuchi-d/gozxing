package common

import (
	"reflect"
	"testing"
)

func TestNewDecoderResult(t *testing.T) {
	rawBytes := make([]byte, 3)
	text := "teststring"
	byteSegments := make([][]byte, 5)
	ecLevel := "H"
	dr := NewDecoderResult(rawBytes, text, byteSegments, ecLevel)

	if !reflect.DeepEqual(dr.GetRawBytes(), rawBytes) {
		t.Fatalf("GetRawBytes() and rawBytes different")
	}
	if r := dr.GetNumBits(); r != 24 {
		t.Fatalf("numBits = %v, expect 24", r)
	}
	dr.SetNumBits(20)
	if r := dr.GetNumBits(); r != 20 {
		t.Fatalf("numBits = %v, expect 20", r)
	}
	if r := dr.GetText(); r != text {
		t.Fatalf("GetText() = %v, expect %v", r, text)
	}
	if !reflect.DeepEqual(dr.GetByteSegments(), byteSegments) {
		t.Fatalf("GetByteSegments() and byteSegments different")
	}
	if r := dr.GetECLevel(); r != ecLevel {
		t.Fatalf("GetECLevel() = %v, expect %v", r, ecLevel)
	}
	dr.SetErrorsCorrected(10)
	if r := dr.GetErrorsCorrected(); r != 10 {
		t.Fatalf("GetErrorsCorrected() = %v, expect %v", r, 10)
	}
	dr.SetErasures(15)
	if r := dr.GetErasures(); r != 15 {
		t.Fatalf("GetErasures() = %v, expect %v", r, 15)
	}
	other := struct{ num int }{25}
	dr.SetOther(other)
	if r := dr.GetOther(); r != other {
		t.Fatalf("GetErasures() = %v, expect %v", r, other)
	}
	if dr.HasStructuredAppend() {
		t.Fatalf("HasStructuredAppend() must be false")
	}
	dr.structuredAppendParity = 0
	dr.structuredAppendSequenceNumber = -1
	if dr.HasStructuredAppend() {
		t.Fatalf("HasStructuredAppend() must be false")
	}
	dr.structuredAppendParity = -1
	dr.structuredAppendSequenceNumber = 0
	if dr.HasStructuredAppend() {
		t.Fatalf("HasStructuredAppend() must be false")
	}
	dr.structuredAppendParity = 0
	dr.structuredAppendSequenceNumber = 0
	if !dr.HasStructuredAppend() {
		t.Fatalf("HasStructuredAppend() must be true")
	}
	dr.structuredAppendParity = 1
	dr.structuredAppendSequenceNumber = 2
	if !dr.HasStructuredAppend() {
		t.Fatalf("HasStructuredAppend() must be true")
	}
	if r := dr.GetStructuredAppendParity(); r != 1 {
		t.Fatalf("GetStructuredAppendParity() = %v, expect %v", r, 1)
	}
	if r := dr.GetStructuredAppendSequenceNumber(); r != 2 {
		t.Fatalf("GetStructuredAppendSequenceNumber() = %v, expect %v", r, 2)
	}

	dr.symbologyModifier = 1
	if r := dr.GetSymbologyModifier(); r != 1 {
		t.Fatalf("GetSymbologyModifier() = %v, expect %v", r, 1)
	}
}

func TestNewDecoderResultWithSymbologyModifier(t *testing.T) {
	symbologyModifier := 3
	dr := NewDecoderResultWithSymbologyModifier([]byte{}, "", [][]byte{}, "L", symbologyModifier)
	if r := dr.GetSymbologyModifier(); r != symbologyModifier {
		t.Fatalf("New WithSymbologyModifier GetSymbologyModifier() = %v, expect %v", r, symbologyModifier)
	}
}

func TestNewDecoderResultWithSA(t *testing.T) {
	saSeqNum := 4
	saParity := 2
	dr := NewDecoderResultWithSA([]byte{}, "", [][]byte{}, "L", saSeqNum, saParity)
	if r := dr.GetStructuredAppendSequenceNumber(); r != saSeqNum {
		t.Fatalf("New WithSA GetStructuredAppendSequenceNumber() = %v, expect %v", r, saSeqNum)
	}
	if r := dr.GetStructuredAppendParity(); r != saParity {
		t.Fatalf("New WithSA GetStructuredAppendParity() = %v, expect %v", r, saParity)
	}
}
