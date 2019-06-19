package reedsolomon

import (
	"fmt"
	"strings"
	"testing"

	errors "golang.org/x/xerrors"

	"github.com/makiuchi-d/gozxing"
)

func testReedSolomonExceptionType(t testing.TB, e error) {
	t.Helper()
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

func TestNewReedSolomonException(t *testing.T) {
	var e error = NewReedSolomonException("test error")

	testReedSolomonExceptionType(t, e)

	wants := "ReedSolomonException: test error"
	if s := e.Error(); s != wants {
		t.Fatalf("e.Error = \"%s\", wants \"%s\"", s, wants)
	}

	s := fmt.Sprintf("%+v", e)
	cases := []string{
		"test error",
		"reedsolomon.TestNewReedSolomonException",
		"common/reedsolomon/reedsolomon_exception_test.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"\n", c)
		}
	}
}

func testError() error {
	return errors.New("test error")
}

func TestWrapReedSolomonException(t *testing.T) {
	testerr := testError()
	e := WrapReedSolomonException(testerr)

	if !errors.Is(e, testerr) {
		t.Fatalf("err is not testerr")
	}

	s := fmt.Sprintf("%+v", e)
	cases := []string{
		"test error",
		"reedsolomon.TestWrapReedSolomonException",
		"common/reedsolomon/reedsolomon_exception_test.go:63",
		"reedsolomon.testError",
		"common/reedsolomon/reedsolomon_exception_test.go:58",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"\n", c)
		}
	}

	wants := "ReedSolomonException: test error"
	if msg := e.Error(); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
	e = WrapReedSolomonException(e)
	if msg := e.Error(); msg != wants {
		t.Fatalf("e.Error() = \"%s\", wants \"%s\"", msg, wants)
	}
}
