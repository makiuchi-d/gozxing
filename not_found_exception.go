package gozxing

import (
	"errors"
)

type NotFoundException struct {
	error
}

func NewNotFoundException(msg string) NotFoundException {
	return NotFoundException{errors.New(msg)}
}
