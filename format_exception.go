package gozxing

import (
	errors "golang.org/x/xerrors"
)

type FormatException interface {
	ReaderException
	formatException()
}

type formatError struct {
	readerError
}

func (formatError) formatException() {}

func (e formatError) Unwrap() error {
	return e.readerError
}

func GetFormatExceptionInstance() FormatException {
	return formatError{
		readerError{
			errors.New("FormatException"),
		},
	}
}

func WrapFormatExceptionInstance(e error) FormatException {
	return formatError{
		readerError{
			errors.Errorf("FormatException: %w", e),
		},
	}
}
