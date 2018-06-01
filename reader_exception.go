package gozxing

type ReaderException interface {
	error
	ReaderException()
}
