package gozxing

import (
	"fmt"
	"testing"

	errors "golang.org/x/xerrors"
)

func testWriterErrorType(t testing.TB, e error) {
	t.Helper()
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

func TestWrapWriterException(t *testing.T) {
	base := errors.New("test error")
	var e error = WrapWriterException(base)

	testWriterErrorType(t, e)

	if !errors.Is(e, base) {
		t.Fatalf("err(%v) is not base(%v)", e, base)
	}

	wants := "WriterException: test error"
	if msg := fmt.Sprintf("%v", e); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
	e = WrapWriterException(e)
	wants = "WriterException: WriterException: test error"
	if msg := fmt.Sprintf("%v", e); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
}
