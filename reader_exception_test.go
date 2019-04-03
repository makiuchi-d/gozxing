package gozxing

import (
	"fmt"
	"strings"
	"testing"
)

type testException struct {
	exception
}

func (testException) readerException() {}

func newTestException(msg string) ReaderException {
	return testException{
		newException(msg, nil),
	}
}

func TestException_Format(t *testing.T) {
	re := newTestException("test error")

	s := fmt.Sprintf("%+v", re)
	cases := []string{
		"test error",
		"TestException_Format",
		"reader_exception_test.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"", c)
		}
	}
}
