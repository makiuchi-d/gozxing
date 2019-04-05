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

func newTestException(args ...interface{}) ReaderException {
	return testException{
		newException("TestException", args...),
	}
}

func TestException_Format(t *testing.T) {
	re := newTestException("%d %x", 10, 10)

	s := fmt.Sprintf("%+v", re)
	cases := []string{
		"TestException: 10 a:",
		"gozxing.TestException_Format",
		"reader_exception_test.go:22",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"\n%s", c, s)
		}
	}
}
