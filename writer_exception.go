package gozxing

import (
	"errors"
)

type WriterException interface {
	error
	WriterException()
}

type writerException struct {
	error
}

func (writerException) WriterException() {}

func NewWriterException(message string) WriterException {
	return writerException{errors.New(message)}
}

func NewWriterExceptionWithError(err error) WriterException {
	return writerException{err}
}
