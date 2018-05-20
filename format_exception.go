package gozxing

import (
	"errors"
)

type FormatException struct {
	error
}

var formatInstance FormatException = FormatException{errors.New("")}

func FormatException_GetFormatInstance() FormatException {
	return formatInstance
}

func FormatException_GetFormatInstanceWithError(err error) FormatException {
	return FormatException{err}
}
