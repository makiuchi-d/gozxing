package gozxing

import (
	"fmt"
	"strings"
	"testing"

	errors "golang.org/x/xerrors"
)

func TestReaderError_Format(t *testing.T) {
	re := readerError{
		errors.New("test error"),
	}

	s := fmt.Sprintf("%+v", re)
	cases := []string{
		"test error",
		"reader_exception_test.go:",
	}
	for _, c := range cases {
		if strings.Index(s, c) < 0 {
			t.Fatalf("error message must contains \"%s\"", c)
		}
	}
}
