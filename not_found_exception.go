package gozxing

import (
	"errors"
)

type NotFoundException struct {
	error
}

var notFoundInstance NotFoundException = NotFoundException{errors.New("NotFoundException")}

func NotFoundException_GetNotFoundInstance() NotFoundException {
	return notFoundInstance
}
