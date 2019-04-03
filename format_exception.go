package gozxing

import (
	"fmt"
)

type FormatException interface {
	ReaderException
	formatException()
}

type formatException struct {
	exception
}

func (formatException) readerException() {}
func (formatException) formatException() {}

func NewFormatException(args ...interface{}) FormatException {
	msg := "FormatException"
	if len(args) > 0 {
		msg += ": " + fmt.Sprintf(args[0].(string), args[1:]...)
	}
	return formatException{
		newException(msg, nil),
	}
}

func NewFormatExceptionWithError(e error) FormatException {
	return formatException{
		newException("FormatException: "+e.Error(), e),
	}
}
