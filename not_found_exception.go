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

func NewNotFoundException() NotFoundException {
	return notFoundException{
		newException("NotFoundException", nil),
	}
}

func WrapNotFoundException(e error) NotFoundException {
	return notFoundException{
		newException("NotFoundException: "+e.Error(), e),
	}
}
