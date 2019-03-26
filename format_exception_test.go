package gozxing

import (
	"testing"

	errors "golang.org/x/xerrors"
)

func testFormatErrorType(t *testing.T, e error) {
	var fe FormatError
	if !errors.As(e, &fe) {
		t.Fatalf("Type must be FormatError")
	}
	var re ReaderError
	if !errors.As(e, &re) {
		t.Fatalf("Type must be ReaderError")
	}
	var ne NotFoundError
	if errors.As(e, &ne) {
		t.Fatalf("Type must not be NotFoundError")
	}

	if _, ok := e.(FormatException); !ok {
		t.Fatalf("Type must be FormatException")
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Type must be ReaderException")
	}
	if _, ok := e.(ChecksumException); ok {
		t.Fatalf("Type must not be ChecksumException")
	}

	fe.formatException()
	fe.readerException()
}

func TestFormatException(t *testing.T) {
	var e error = GetFormatExceptionInstance()
	testFormatErrorType(t, e)
}

func TestNewFormatException(t *testing.T) {
	base := errors.New("newformatexception")
	var e error = WrapFormatExceptionInstance(base)

	testFormatErrorType(t, e)

	if !errors.Is(e, base) {
		t.Fatalf("err(%v) is not base(%v)", e, base)
	}
}
