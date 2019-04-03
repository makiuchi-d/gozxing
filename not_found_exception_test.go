package gozxing

import (
	"fmt"
	"strings"
	"testing"

	errors "golang.org/x/xerrors"
)

func testNotFoundExceptionType(t *testing.T, e error) {
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

func TestNotFoundException(t *testing.T) {
	var e error = GetNotFoundExceptionInstance()
	testNotFoundExceptionType(t, e)

	s := fmt.Sprintf("%+v", e)
	cases := []string{
		"NotFoundException",
		"TestNotFoundException",
		"not_found_exception_test.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"", c)
		}
	}
}
