package reedsolomon

import (
	"testing"
)

func TestReedsolomonException(t *testing.T) {
	var e error = NewReedSolomonException("newreedsolomonexception")
	if _, ok := e.(ReedSolomonException); !ok {
		t.Fatalf("error must be ReedSolomonException, %T", e)
	}

	e.(ReedSolomonException).ReedSolomonException()
}
