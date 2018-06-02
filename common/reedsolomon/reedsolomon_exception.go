package reedsolomon

import (
	"errors"
)

type ReedSolomonException interface {
	error
	ReedSolomonException()
}

type reedSolomonException struct {
	error
}

func (reedSolomonException) ReedSolomonException() {}

func NewReedSolomonException(msg string) ReedSolomonException {
	return reedSolomonException{errors.New(msg)}
}
