package gozxing

import (
	"errors"
)

type FormatException interface {
	ReaderException
	FormatException()
}

type formatException struct {
	error
}

func (formatException) ReaderException() {}
func (formatException) FormatException() {}

var formatInstance = formatException{errors.New("FormatException")}

func GetFormatExceptionInstance() FormatException {
	return formatInstance
}

func NewFormatExceptionInstance(e error) FormatException {
	return formatException{e}
}
