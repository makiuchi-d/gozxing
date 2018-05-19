package gozxing

type DecodeHintType int

const (
	DecodeHintType_OTHER DecodeHintType = iota
	DecodeHintType_PURE_BARCODE
	DecodeHintType_POSSIBLE_FORMATS
	DecodeHintType_TRY_HARDER
	DecodeHintType_CHARACTER_SET
	DecodeHintType_ALLOWED_LENGTHS
	DecodeHintType_ASSUME_CODE_39_CHECK_DIGIT
	DecodeHintType_ASSUME_GS1
	DecodeHintType_RETURN_CODABAR_START_END
	DecodeHintType_NEED_RESULT_POINT_CALLBACK
	DecodeHintType_ALLOWED_EAN_EXTENSIONS
)

func (t DecodeHintType) String() string {
	switch t {
	case DecodeHintType_OTHER:
		return "OTHER"
	case DecodeHintType_PURE_BARCODE:
		return "PURE_BARCODE"
	case DecodeHintType_POSSIBLE_FORMATS:
		return "POSSIBLE_FORMATS"
	case DecodeHintType_TRY_HARDER:
		return "TRY_HARDER"
	case DecodeHintType_CHARACTER_SET:
		return "CHARACTER_SET"
	case DecodeHintType_ALLOWED_LENGTHS:
		return "ALLOWED_LENGTHS"
	case DecodeHintType_ASSUME_CODE_39_CHECK_DIGIT:
		return "ASSUME_CODE_39_CHECK_DIGIT"
	case DecodeHintType_ASSUME_GS1:
		return "ASSUME_GS1"
	case DecodeHintType_RETURN_CODABAR_START_END:
		return "RETURN_CODABAR_START_END"
	case DecodeHintType_NEED_RESULT_POINT_CALLBACK:
		return "NEED_RESULT_POINT_CALLBACK"
	case DecodeHintType_ALLOWED_EAN_EXTENSIONS:
		return "ALLOWED_EAN_EXTENSIONS"
	}
	return "Unknown DecodeHintType"
}
