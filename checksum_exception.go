package gozxing

import (
	errors "golang.org/x/xerrors"
)

type ChecksumException interface {
	ReaderException
	checksumException()
}

type ChecksumError struct {
	ReaderError
}

func (ChecksumError) checksumException() {}

func (e ChecksumError) Unwrap() error {
	return e.ReaderError
}

func GetChecksumExceptionInstance() ChecksumException {
	return ChecksumError{
		ReaderError{
			errors.New("ChecksumException"),
		},
	}
}

func NewChecksumExceptionInstance(e error) ChecksumException {
	return ChecksumError{
		ReaderError{
			errors.Errorf("ChecksumException: %w", e),
		},
	}
}
