package gozxing

import (
	"errors"
)

type ChecksumException interface {
	ReaderException
	ChecksumException()
}

type checksumException struct {
	error
}

func (checksumException) ReaderException()   {}
func (checksumException) ChecksumException() {}

var checksumInstance = checksumException{errors.New("ChecksumException")}

func GetChecksumExceptionInstance() ChecksumException {
	return checksumInstance
}

func NewChecksumExceptionInstance(e error) ChecksumException {
	return checksumException{e}
}
