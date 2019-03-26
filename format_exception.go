package gozxing

import (
	errors "golang.org/x/xerrors"
)

type FormatException interface {
	ReaderException
	formatException()
}

type FormatError struct {
	ReaderError
}

func (FormatError) formatException() {}

func (e FormatError) Unwrap() error {
	return e.ReaderError
}

func GetFormatExceptionInstance() FormatException {
	return FormatError{
		ReaderError{
			errors.New("FormatException"),
		},
	}
}

func WrapFormatExceptionInstance(e error) FormatException {
	return FormatError{
		ReaderError{
			errors.Errorf("FormatException: %w", e),
		},
	}
}
