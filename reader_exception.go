package gozxing

import (
	"fmt"

	errors "golang.org/x/xerrors"
)

type ReaderException interface {
	error
	readerException()
}

type ReaderError struct {
	error
}

func (ReaderError) readerException() {}

func (e ReaderError) Unwrap() error {
	return e.error
}

func (e ReaderError) Format(s fmt.State, v rune) {
	errors.FormatError(e.error.(errors.Formatter), s, v)
}
