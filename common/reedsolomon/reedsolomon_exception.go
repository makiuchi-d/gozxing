package reedsolomon

import (
	"fmt"

	errors "golang.org/x/xerrors"
)

type ReedSolomonException interface {
	error
	reedSolomonException()
}

type reedSolomonException struct {
	msg   string
	next  error
	frame errors.Frame
}

func (reedSolomonException) reedSolomonException() {}

func (e reedSolomonException) Error() string {
	return e.msg
}

func (e reedSolomonException) Unwrap() error {
	return e.next
}

func (e reedSolomonException) Format(s fmt.State, v rune) {
	errors.FormatError(e, s, v)
}

func (e reedSolomonException) FormatError(p errors.Printer) error {
	p.Print(e.msg)
	e.frame.Format(p)
	return e.next
}

func NewReedSolomonException(msg string) ReedSolomonException {
	return reedSolomonException{
		"ReedSolomonException: " + msg,
		nil,
		errors.Caller(1),
	}
}

func WrapReedSolomonException(err error) ReedSolomonException {
	return reedSolomonException{
		"ReedSolomonException: " + err.Error(),
		err,
		errors.Caller(1),
	}
}
