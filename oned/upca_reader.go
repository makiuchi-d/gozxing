package oned

import (
	"github.com/makiuchi-d/gozxing"
)

type upcAReader struct {
	*OneDReader
}

type upcAInternalDecoder struct {
	*ean13Reader
}

func NewUPCAReader() gozxing.Reader {
	decoder := &upcAInternalDecoder{
		&ean13Reader{
			decodeMiddleCounters: make([]int, 4),
		},
	}
	return &upcAReader{
		NewOneDReader(NewUPCEANReader(decoder)),
	}
}

func (this *upcAReader) decodeRow(rowNumber int, row *gozxing.BitArray, hints map[gozxing.DecodeHintType]interface{}) (*gozxing.Result, error) {
	return maybeReturnResult(this.OneDReader.decodeRow(rowNumber, row, hints))
}

func (this *upcAReader) decodeRowWithStartRange(
	rowNumber int, row *gozxing.BitArray, startGuardRange []int,
	hints map[gozxing.DecodeHintType]interface{}) (*gozxing.Result, error) {
	decoder := this.OneDReader.rowDecoder.(*upceanReader)
	return maybeReturnResult(decoder.decodeRowWithStartRange(rowNumber, row, startGuardRange, hints))
}

func (this *upcAReader) DecodeWithoutHints(image *gozxing.BinaryBitmap) (*gozxing.Result, error) {
	return maybeReturnResult(this.OneDReader.DecodeWithoutHints(image))
}

func (this *upcAReader) Decode(image *gozxing.BinaryBitmap, hints map[gozxing.DecodeHintType]interface{}) (*gozxing.Result, error) {
	return maybeReturnResult(this.OneDReader.Decode(image, hints))
}

func (this *upcAInternalDecoder) getBarcodeFormat() gozxing.BarcodeFormat {
	return gozxing.BarcodeFormat_UPC_A
}

func maybeReturnResult(result *gozxing.Result, e error) (*gozxing.Result, error) {
	if e != nil {
		return nil, e
	}
	text := result.GetText()
	if text[0] == '0' {
		upcaResult := gozxing.NewResult(
			text[1:],
			nil,
			result.GetResultPoints(),
			gozxing.BarcodeFormat_UPC_A)
		if result.GetResultMetadata() != nil {
			upcaResult.PutAllMetadata(result.GetResultMetadata())
		}
		return upcaResult, nil
	} else {
		return nil, gozxing.NewFormatException()
	}
}
