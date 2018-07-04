package gozxing

type Writer interface {
	EncodeWithoutHint(contents string, format BarcodeFormat, width, height int) (*BitMatrix, error)
	Encode(contents string, format BarcodeFormat, width, height int, hints map[EncodeHintType]interface{}) (*BitMatrix, error)
}
