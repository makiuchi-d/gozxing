package gozxing

type ChecksumException interface {
	ReaderException
	checksumException()
}

type checksumException struct {
	exception
}

func (checksumException) readerException()   {}
func (checksumException) checksumException() {}

func NewChecksumException() ChecksumException {
	return checksumException{
		newException("ChecksumException", nil),
	}
}

func WrapChecksumException(e error) ChecksumException {
	return checksumException{
		newException("ChecksumException"+e.Error(), e),
	}
}
