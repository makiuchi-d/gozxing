package oned

import (
	"fmt"
	"strconv"

	"github.com/makiuchi-d/gozxing"
)

const (
	ean13Writer_CODE_WIDTH = 3 + // start guard
		(7 * 6) + // left bars
		5 + // middle guard
		(7 * 6) + // right bars
		3 // end guard
)

type ean13Encoder struct{}

type EAN13Writer struct {
	*OneDimensionalCodeWriter
}

func NewEAN13Writer() gozxing.Writer {
	return NewUPCEANWriter(ean13Encoder{})
}

func (ean13Encoder) getFormat() gozxing.BarcodeFormat {
	return gozxing.BarcodeFormat_EAN_13
}

func (ean13Encoder) encodeContents(contents string) ([]bool, error) {
	length := len(contents)
	switch length {
	case 12:
		// No check digit present, calculate it and add it
		check, e := upceanReader_getStandardUPCEANChecksum(contents)
		if e != nil {
			return nil, fmt.Errorf("IllegalArgumentException: %s", e.Error())
		}
		contents += strconv.Itoa(check)
		break
	case 13:
		ok, e := upceanReader_checkStandardUPCEANChecksum(contents)
		if e != nil {
			return nil, fmt.Errorf("IllegalArgumentException: Illegal contents, %v", e)
		}
		if !ok {
			return nil, fmt.Errorf("IllegalArgumentException: Contents do not pass checksum")
		}
		break
	default:
		return nil, fmt.Errorf("IllegalArgumentException: "+
			"Requested contents should be 12 or 13 digits long, but got %v", length)
	}

	firstDigit := contents[0] - '0'
	parities := ean13Reader_FIRST_DIGIT_ENCODINGS[firstDigit]
	result := make([]bool, ean13Writer_CODE_WIDTH)
	pos := 0

	pos += onedWriter_appendPattern(result, pos, UPCEANReader_START_END_PATTERN, true)

	// See EAN13Reader for a description of how the first digit & left bars are encoded
	for i := 1; i <= 6; i++ {
		digit := contents[i] - '0'
		if ((parities >> uint(6-i)) & 1) == 1 {
			digit += 10
		}
		pos += onedWriter_appendPattern(result, pos, UPCEANReader_L_AND_G_PATTERNS[digit], false)
	}

	pos += onedWriter_appendPattern(result, pos, UPCEANReader_MIDDLE_PATTERN, false)

	for i := 7; i <= 12; i++ {
		digit := contents[i] - '0'
		pos += onedWriter_appendPattern(result, pos, UPCEANReader_L_PATTERNS[digit], true)
	}
	onedWriter_appendPattern(result, pos, UPCEANReader_START_END_PATTERN, true)

	return result, nil
}
