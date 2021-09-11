package gozxing

import (
	"fmt"
	"strings"
	"testing"

	errors "golang.org/x/xerrors"
)

func testChecksumExceptionType(t testing.TB, e error) {
	t.Helper()
	var ce ChecksumException
	if !errors.As(e, &ce) {
		t.Fatalf("Type must be FormatException")
	}
	var re ReaderException
	if !errors.As(e, &re) {
		t.Fatalf("Type must be ReaderException")
	}
	var fe FormatException
	if errors.As(e, &fe) {
		t.Fatalf("Type must not be FormatException")
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

func TestNewChecksumException(t *testing.T) {
	var e error = NewChecksumException()
	testChecksumExceptionType(t, e)

	s := fmt.Sprintf("%+v", e)
	cases := []string{
		"ChecksumException",
		"TestNewChecksumException",
		"checksum_exception_test.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"", c)
		}
	}
}

func TestWrapChecksumException(t *testing.T) {
	base := errors.New("newchecksumexceptionstring")
	var e error = WrapChecksumException(base)

	testChecksumExceptionType(t, e)

	if !errors.Is(e, base) {
		t.Fatalf("err(%v) is not base(%v)", e, base)
	}

	wants := "ChecksumException: newchecksumexceptionstring"
	if msg := fmt.Sprintf("%v", e); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
	e = WrapChecksumException(e)
	wants = "ChecksumException: ChecksumException: newchecksumexceptionstring"
	if msg := fmt.Sprintf("%v", e); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
}
