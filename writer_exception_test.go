package gozxing

import (
	"errors"
	"testing"
)

func TestWriteerException(t *testing.T) {
	var e error

	e = NewWriterException("test message")
	if _, ok := e.(WriterException); !ok {
		t.Fatalf("Not WriterException, %T", e)
	}
	if _, ok := e.(ReaderException); ok {
		t.Fatalf("Type must not be ReaderException")
	}

	e = NewWriterExceptionWithError(errors.New("test2 message"))
	if _, ok := e.(WriterException); !ok {
		t.Fatalf("Not WriterException, %T", e)
	}
	if _, ok := e.(ReaderException); ok {
		t.Fatalf("Type must not be ReaderException")
	}

	e.(WriterException).WriterException()
}
