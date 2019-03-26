package gozxing

import (
	"testing"

	errors "golang.org/x/xerrors"
)

func testChecksumErrorType(t *testing.T, e error) {
	var ce ChecksumError
	if !errors.As(e, &ce) {
		t.Fatalf("Type must be FormatError")
	}
	var re ReaderError
	if !errors.As(e, &re) {
		t.Fatalf("Type must be ReaderError")
	}
	var fe FormatError
	if errors.As(e, &fe) {
		t.Fatalf("Type must not be FormatError")
	}

	if _, ok := e.(ChecksumException); !ok {
		t.Fatalf("Type must be ChecksumException")
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Type must be ReaderException")
	}
	if _, ok := e.(NotFoundException); ok {
		t.Fatalf("Type must not be NotFoundException")
	}

	ce.checksumException()
	ce.readerException()
}

func TestChecksumException(t *testing.T) {
	var e error = GetChecksumExceptionInstance()
	testChecksumErrorType(t, e)
}

func TestNewChecksumException(t *testing.T) {
	base := errors.New("newchecksumexceptionstring")
	var e error = NewChecksumExceptionInstance(base)

	testChecksumErrorType(t, e)

	if !errors.Is(e, base) {
		t.Fatalf("err(%v) is not base(%v)", e, base)
	}
}
