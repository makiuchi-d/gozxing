package gozxing

import (
	"fmt"
	"strings"
	"testing"

	errors "golang.org/x/xerrors"

)

func testWriterErrorType(t *testing.T, e error) {
	var we WriterException
	if !errors.As(e, &we) {
		t.Fatalf("Type must be WriterException")
	}
	var re ReaderException
	if errors.As(e, &re) {
		t.Fatalf("Type must not be ReaderException")
	}

	if _, ok := e.(WriterException); !ok {
		t.Fatalf("Type must be WriterException")
	}
	if _, ok := e.(ReaderException); ok {
		t.Fatalf("Type must not be ReaderException")
	}

	we.writerException()
}


func TestNewWriterException(t *testing.T) {
	var e error = NewWriterException("test message")
	testWriterErrorType(t, e)
}

func TestNewWriterExceptionWithError(t *testing.T) {
	base := errors.New("test error")
	var e error = NewWriterExceptionWithError(base)

	testWriterErrorType(t, e)

	if !errors.Is(e, base) {
		t.Fatalf("err(%v) is not base(%v)", e, base)
	}
}

func TestWriterError_Format(t *testing.T) {
	e := NewWriterException("test error")

	s := fmt.Sprintf("%+v", e)
	cases := []string{
		"test error",
		"writer_exception.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"", c)
		}
	}
}
