package reedsolomon

import (
	"fmt"

	errors "golang.org/x/xerrors"
)

type ReedSolomonException interface {
	error
	reedSolomonException()
}

type reedSolomonError struct {
	error
}

func (reedSolomonError) reedSolomonException() {}

func (e reedSolomonError) Unwrap() error {
	return e.error
}

func (e reedSolomonError) Format(s fmt.State, v rune) {
	errors.FormatError(e.error.(errors.Formatter), s, v)
}

func NewReedSolomonException(msg string) ReedSolomonException {
	return reedSolomonError{
		errors.Errorf("ReedSolomonException: %s", msg),
	}
}
