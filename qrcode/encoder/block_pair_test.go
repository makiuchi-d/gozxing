package encoder

import (
	"reflect"
	"testing"
)

func TestBlockPair(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	ecb := []byte{3, 5, 7}

	b := NewBlockPair(data, ecb)

	if r := b.GetDataBytes(); !reflect.DeepEqual(r, data) {
		t.Fatalf("GetDataBytes = %v, expect %v", r, data)
	}
	if r := b.GetErrorCorrectionBytes(); !reflect.DeepEqual(r, ecb) {
		t.Fatalf("GetErrorCorrectionBytes = %v, expect %v", r, ecb)
	}
}
