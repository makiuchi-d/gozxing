package gozxing

import (
	errors "golang.org/x/xerrors"
)

type NotFoundException interface {
	ReaderException
	notFoundException()
}

type NotFoundError struct {
	ReaderError
}

func (NotFoundError) notFoundException() {}

func (e NotFoundError) Unwrap() error {
	return e.ReaderError
}

func GetNotFoundExceptionInstance() NotFoundError {
	return NotFoundError{
		ReaderError{
			errors.New("NotFoundException"),
		},
	}
}
