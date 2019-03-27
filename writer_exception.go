package gozxing

import (
	"fmt"

	errors "golang.org/x/xerrors"
)

type WriterException interface {
	error
	writerException()
}

type writerError struct {
	error
}

func (writerError) writerException() {}

func (e writerError) Unwrap() error {
	return e.error
}

func (e writerError) Format(s fmt.State, v rune) {
	errors.FormatError(e.error.(errors.Formatter), s, v)
}

func NewWriterException(message string) WriterException {
	return writerError{
		errors.Errorf("WriterException: %s", message),
	}
}

func NewWriterExceptionWithError(err error) WriterException {
	return writerError{
		errors.Errorf("WriterException: %w", err),
	}
}
