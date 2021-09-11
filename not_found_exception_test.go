package gozxing

import (
	"fmt"
	"strings"
	"testing"

	errors "golang.org/x/xerrors"
)

func testNotFoundExceptionType(t testing.TB, e error) {
	t.Helper()
	var ne NotFoundException
	if !errors.As(e, &ne) {
		t.Fatalf("Type must be NotFoundException")
	}
	var re ReaderException
	if !errors.As(e, &re) {
		t.Fatalf("Type must be ReaderException")
	}
	var ce ChecksumException
	if errors.As(e, &ce) {
		t.Fatalf("Type must not be ChecksumException")
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

func TestNewNotFoundException(t *testing.T) {
	var e error = NewNotFoundException()
	testNotFoundExceptionType(t, e)

	s := fmt.Sprintf("%+v", e)
	cases := []string{
		"NotFoundException",
		"TestNewNotFoundException",
		"not_found_exception_test.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"", c)
		}
	}

	e = NewNotFoundException("not %s", "found")
	msg := e.Error()
	wants := "NotFoundException: not found"
	if msg != wants {
		t.Fatalf("Error() = \"%s\", wants \"%s\"", msg, wants)
	}
}

func TestWrapNotFoundException(t *testing.T) {
	base := errors.New("newnotfoundexception")
	var e error = WrapNotFoundException(base)

	testNotFoundExceptionType(t, e)

	if !errors.Is(e, base) {
		t.Fatalf("err(%v) is not base(%v)", e, base)
	}

	wants := "NotFoundException: newnotfoundexception"
	if msg := fmt.Sprintf("%v", e); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
	e = WrapNotFoundException(e)
	wants = "NotFoundException: NotFoundException: newnotfoundexception"
	if msg := fmt.Sprintf("%v", e); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
}
