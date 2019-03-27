package gozxing

import (
	errors "golang.org/x/xerrors"
)

type NotFoundException interface {
	ReaderException
	notFoundException()
}

type notFoundError struct {
	readerError
}

func (notFoundError) notFoundException() {}

func (e notFoundError) Unwrap() error {
	return e.readerError
}

func GetNotFoundExceptionInstance() NotFoundException {
	return notFoundError{
		readerError{
			errors.New("NotFoundException"),
		},
	}
}
