package gozxing

import (
	"errors"
)

type FormatException struct {
	error
}

func NewFormatException(msg string) FormatException {
	return FormatException{errors.New(msg)}
}
