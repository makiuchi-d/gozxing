package gozxing

import (
	"errors"
)

type NotFoundException interface {
	ReaderException
	NotFoundException()
}

type notFoundException struct {
	error
}

func (notFoundException) ReaderException()   {}
func (notFoundException) NotFoundException() {}

var notFoundInstance = notFoundException{errors.New("NotFoundException")}

func GetNotFoundExceptionInstance() NotFoundException {
	return notFoundInstance
}
