package gozxing

import (
	"testing"
)

func TestNotFoundException(t *testing.T) {
	var e error = GetNotFoundExceptionInstance()

	if _, ok := e.(NotFoundException); !ok {
		t.Fatalf("Not NotFoundException, %T", e)
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Not ReaderException, %T", e)
	}
	if _, ok := e.(FormatException); ok {
		t.Fatalf("Type must not be FormatException")
	}

	e.(NotFoundException).NotFoundException()
	e.(NotFoundException).ReaderException()
}
