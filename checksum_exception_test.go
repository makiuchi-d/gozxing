package gozxing

import (
	"errors"
	"testing"
)

func TestChecksumException(t *testing.T) {
	var e error = GetChecksumExceptionInstance()

	if _, ok := e.(ChecksumException); !ok {
		t.Fatalf("Type is not ChecksumException")
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Type is not ReaderException")
	}
	if _, ok := e.(NotFoundException); ok {
		t.Fatalf("Type must not be NotFoundException")
	}
}

func TestNewChecksumException(t *testing.T) {
	base := errors.New("newchecksumexceptionstring")
	var e error = NewChecksumExceptionInstance(base)

	if _, ok := e.(ChecksumException); !ok {
		t.Fatalf("Type is not ChecksumException")
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Type is not ReaderException")
	}
	if _, ok := e.(NotFoundException); ok {
		t.Fatalf("Type must not be NotFoundException")
	}
	if e.Error() != base.Error() {
		t.Fatalf("e.Error() = %v, expect %v", e.Error(), base.Error())
	}
}
