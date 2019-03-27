package gozxing

import (
	"fmt"

	errors "golang.org/x/xerrors"
)

type ReaderException interface {
	error
	readerException()
}

type readerError struct {
	error
}

func (readerError) readerException() {}

func (e readerError) Unwrap() error {
	return e.error
}

func (e readerError) Format(s fmt.State, v rune) {
	errors.FormatError(e.error.(errors.Formatter), s, v)
}
