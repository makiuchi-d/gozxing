package gozxing

type NotFoundException interface {
	ReaderException
	notFoundException()
}

type notFoundException struct {
	exception
}

func (notFoundException) readerException()   {}
func (notFoundException) notFoundException() {}

func GetNotFoundExceptionInstance() NotFoundException {
	return notFoundException{
		newException("NotFoundException", nil),
	}
}
