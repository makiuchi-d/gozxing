package gozxing

import (
	errors "golang.org/x/xerrors"
)

type ChecksumException interface {
	ReaderException
	checksumException()
}

type checksumError struct {
	readerError
}

func (checksumError) checksumException() {}

func (e checksumError) Unwrap() error {
	return e.readerError
}

func GetChecksumExceptionInstance() ChecksumException {
	return checksumError{
		readerError{
			errors.New("ChecksumException"),
		},
	}
}

func NewChecksumExceptionInstance(e error) ChecksumException {
	return checksumError{
		readerError{
			errors.Errorf("ChecksumException: %w", e),
		},
	}
}
