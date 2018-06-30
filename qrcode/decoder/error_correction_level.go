package decoder

import (
	"errors"
)

type ErrorCorrectionLevel int

const (
	ErrorCorrectionLevel_L ErrorCorrectionLevel = 0x01 // ~7% correction
	ErrorCorrectionLevel_M ErrorCorrectionLevel = 0x00 // ~15% correction
	ErrorCorrectionLevel_Q ErrorCorrectionLevel = 0x03 // ~25% correction
	ErrorCorrectionLevel_H ErrorCorrectionLevel = 0x02 // ~30% correction
)

func ErrorCorrectionLevel_ForBits(bits uint) (ErrorCorrectionLevel, error) {
	switch bits {
	case 0:
		return ErrorCorrectionLevel_M, nil
	case 1:
		return ErrorCorrectionLevel_L, nil
	case 2:
		return ErrorCorrectionLevel_H, nil
	case 3:
		return ErrorCorrectionLevel_Q, nil
	}
	return -1, errors.New("IllegalArgumentException")
}

func (e ErrorCorrectionLevel) String() string {
	switch e {
	case ErrorCorrectionLevel_M:
		return "M"
	case ErrorCorrectionLevel_L:
		return "L"
	case ErrorCorrectionLevel_H:
		return "H"
	case ErrorCorrectionLevel_Q:
		return "Q"
	}
	return ""
}
