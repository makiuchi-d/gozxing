package reedsolomon

import (
	"fmt"
	"strings"
	"testing"

	errors "golang.org/x/xerrors"

	"github.com/makiuchi-d/gozxing"
)

func testReedSolomonErrorType(t *testing.T, e error) {
	var rse ReedSolomonException
	if !errors.As(e, &rse) {
		t.Fatalf("Type must be ReedSolomonException")
	}
	var re gozxing.ReaderException
	if errors.As(e, &re) {
		t.Fatalf("Type must not be RederException")
	}

	if _, ok := e.(ReedSolomonException); !ok {
		t.Fatalf("Type must be ReedSolomonException")
	}
	if _, ok := e.(gozxing.ReaderException); ok {
		t.Fatalf("Type must no be ReaderException")
	}

	rse.reedSolomonException()
}

func TestReedsolomonException(t *testing.T) {
	var e error = NewReedSolomonException("newreedsolomonexception")

	testReedSolomonErrorType(t, e)

	s := fmt.Sprintf("%+v", e)
	cases := []string{
		"newreedsolomonexception",
		"common/reedsolomon/reedsolomon_exception.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"\n", c)
		}
	}
}

func TestReedSolomonError_Unwrap(t *testing.T) {
	e := NewReedSolomonException("test error").(reedSolomonError)

	if e.Unwrap() != e.error {
		t.Fatalf("e.Unwrap() must return e.error")
	}
}
