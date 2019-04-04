package gozxing

import (
	"fmt"
)

type NotFoundException interface {
	ReaderException
	notFoundException()
}

type notFoundException struct {
	exception
}

func (notFoundException) readerException()   {}
func (notFoundException) notFoundException() {}

func NewNotFoundException(args ...interface{}) NotFoundException {
	msg := "NotFoundException"
	if len(args) > 0 {
		msg += ": " + fmt.Sprintf(args[0].(string), args[1:]...)
	}
	return notFoundException{
		newException(msg, nil),
	}
}

func WrapNotFoundException(e error) NotFoundException {
	return notFoundException{
		newException("NotFoundException: "+e.Error(), e),
	}
}
