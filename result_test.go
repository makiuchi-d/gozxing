package gozxing

import (
	"reflect"
	"testing"
	"time"
)

func TestResult(t *testing.T) {
	text := "teststr"
	rawBytes := []byte{1, 2, 3, 4}
	resultPoints := []ResultPoint{
		NewResultPoint(10, 20),
		NewResultPoint(20, 30),
	}
	format := BarcodeFormat_AZTEC

	nowmilli := time.Now().UnixNano() / int64(time.Millisecond)
	result := NewResult(text, rawBytes, resultPoints, format)

	if r := result.GetText(); r != text {
		t.Fatalf("GetText = %v, expect %v", r, text)
	}
	if r := result.GetRawBytes(); !reflect.DeepEqual(r, rawBytes) {
		t.Fatalf("GetRawBytes = %v, expect %v", r, rawBytes)
	}
	if r := result.GetNumBits(); r != len(rawBytes)*8 {
		t.Fatalf("GetNumBits = %v, expect %v", r, len(rawBytes)*8)
	}
	if r := result.GetResultPoints(); !reflect.DeepEqual(r, resultPoints) {
		t.Fatalf("GetResultPoints = %v, expect %v", r, resultPoints)
	}
	if r := result.GetBarcodeFormat(); r != format {
		t.Fatalf("GetBarcodeFormat = %v, expect %v", r, format)
	}
	if r := result.GetResultMetadata(); r != nil {
		t.Fatalf("GetResulMetadata = %v, expect %v", r, nil)
	}
	if r := result.GetTimestamp(); r < nowmilli-100 || nowmilli+100 < r {
		t.Fatalf("GetTimestamp = %v, expect %v", r, nowmilli)
	}
	if r := result.String(); r != text {
		t.Fatalf("String = %v, expect %v", r, text)
	}

	metadata := map[ResultMetadataType]interface{}{
		ResultMetadataType_ORIENTATION:                90,
		ResultMetadataType_ERROR_CORRECTION_LEVEL:     "M",
		ResultMetadataType_STRUCTURED_APPEND_SEQUENCE: 0x14,
		ResultMetadataType_STRUCTURED_APPEND_PARITY:   0xab,
	}
	key := ResultMetadataType_STRUCTURED_APPEND_SEQUENCE
	metadata1 := make(map[ResultMetadataType]interface{})
	metadata1[key] = metadata[key]

	result.PutMetadata(key, metadata[key])
	if r := result.GetResultMetadata(); !reflect.DeepEqual(r, metadata1) {
		t.Fatalf("metadata = %v, expect %v", r, metadata1)
	}
	result.PutAllMetadata(nil)
	if r := result.GetResultMetadata(); !reflect.DeepEqual(r, metadata1) {
		t.Fatalf("metadata = %v, expect %v", r, metadata1)
	}
	result.PutAllMetadata(metadata)
	if r := result.GetResultMetadata(); !reflect.DeepEqual(r, metadata) {
		t.Fatalf("metadata = %v, expect %v", r, metadata)
	}

	newPoints := []ResultPoint{NewResultPoint(30, 40)}
	allPoints := append(resultPoints, newPoints...)
	result.AddResultPoints(newPoints)
	if r := result.GetResultPoints(); !reflect.DeepEqual(r, allPoints) {
		t.Fatalf("resultPoints = %v, expect %v", r, allPoints)
	}

	// test GetNumBits when RawBytes is nil
	result = NewResult("", nil, nil, 1)
	if r := result.GetNumBits(); r != 0 {
		t.Fatalf("GetNumBits = %v, expect %v", r, 0)
	}
	// put all metadata when metadata=nil
	result.PutAllMetadata(metadata)
	if r := result.GetResultMetadata(); !reflect.DeepEqual(r, metadata) {
		t.Fatalf("metadata = %v, expect %v", r, metadata)
	}
	// add resultpoint to empty
	result.AddResultPoints(newPoints)
	if r := result.GetResultPoints(); !reflect.DeepEqual(r, newPoints) {
		t.Fatalf("resultPoints = %v, expect %v", r, newPoints)
	}
}
