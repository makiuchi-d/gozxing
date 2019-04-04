package gozxing

import (
	"fmt"
)

type WriterException interface {
	error
	writerException()
}

type writerException struct {
	exception
}

func (writerException) writerException() {}

func NewWriterException(message string, args ...interface{}) WriterException {
	return writerException{
		newException(
			fmt.Sprintf("WriterException: "+message, args...),
			nil),
	}
}

func WrapWriterException(err error) WriterException {
	return writerException{
		newException("WriterException: "+err.Error(), err),
	}
}
