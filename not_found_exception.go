package gozxing

import (
	"errors"
)

type NotFoundException struct {
	error
}

var notFoundInstance NotFoundException = NotFoundException{errors.New("")}

func NotFoundException_GetNotFoundInstance() NotFoundException {
	return notFoundInstance
}
