package gozxing

import (
	"testing"

	errors "golang.org/x/xerrors"
)

func testNotFoundErrorType(t *testing.T, e error) {
	var ne NotFoundError
	if !errors.As(e, &ne) {
		t.Fatalf("Type must be NotFoundError")
	}
	var re ReaderError
	if !errors.As(e, &re) {
		t.Fatalf("Type must be ReaderError")
	}
	var ce ChecksumError
	if errors.As(e, &ce) {
		t.Fatalf("Type must not be ChecksumError")
	}

	if _, ok := e.(NotFoundException); !ok {
		t.Fatalf("Type must be NotFoundException")
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Type must be ReaderException")
	}
	if _, ok := e.(FormatException); ok {
		t.Fatalf("Type must not be FormatException")
	}

	ne.notFoundException()
	ne.readerException()
}

func TestNotFoundException(t *testing.T) {
	var e error = GetNotFoundExceptionInstance()
	testNotFoundErrorType(t, e)
}
