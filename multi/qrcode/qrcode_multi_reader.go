package qrcode

import (
	"sort"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/multi"
	"github.com/makiuchi-d/gozxing/multi/qrcode/detector"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
)

// This implementation can detect and decode multiple QR Codes in an image.

var (
	noPoints = []gozxing.ResultPoint{}
)

type QRCodeMultiReader struct {
	*qrcode.QRCodeReader
}

func NewQRCodeMultiReader() multi.MultipleBarcodeReader {
	return &QRCodeMultiReader{
		qrcode.NewQRCodeReader().(*qrcode.QRCodeReader),
	}
}

func (this *QRCodeMultiReader) DecodeMultipleWithoutHint(image *gozxing.BinaryBitmap) ([]*gozxing.Result, error) {
	return this.DecodeMultiple(image, nil)
}

func (this *QRCodeMultiReader) DecodeMultiple(image *gozxing.BinaryBitmap, hints map[gozxing.DecodeHintType]interface{}) ([]*gozxing.Result, error) {
	results := make([]*gozxing.Result, 0)
	matrix, e := image.GetBlackMatrix()
	if e != nil {
		return results, e
	}
	detectorResults, e := detector.NewMultiDetector(matrix).DetectMulti(hints)
	if e != nil {
		return results, e
	}
	for _, detectorResult := range detectorResults {
		decoderResult, e := this.GetDecoder().Decode(detectorResult.GetBits(), hints)
		if e != nil {
			if _, ok := e.(gozxing.ReaderException); ok {
				// ignore and continue
				continue
			} else {
				return results, e
			}
		}
		points := detectorResult.GetPoints()
		// If the code was mirrored: swap the bottom-left and the top-right points.
		if metadata, ok := decoderResult.GetOther().(*decoder.QRCodeDecoderMetaData); ok {
			metadata.ApplyMirroredCorrection(points)
		}
		result := gozxing.NewResult(decoderResult.GetText(), decoderResult.GetRawBytes(), points,
			gozxing.BarcodeFormat_QR_CODE)
		byteSegments := decoderResult.GetByteSegments()
		if byteSegments != nil {
			result.PutMetadata(gozxing.ResultMetadataType_BYTE_SEGMENTS, byteSegments)
		}
		ecLevel := decoderResult.GetECLevel()
		if ecLevel != "" {
			result.PutMetadata(gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL, ecLevel)
		}
		if decoderResult.HasStructuredAppend() {
			result.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE,
				decoderResult.GetStructuredAppendSequenceNumber())
			result.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_PARITY,
				decoderResult.GetStructuredAppendParity())
		}
		results = append(results, result)
	}
	if len(results) != 0 {
		results = processStructuredAppend(results)
	}
	return results, nil
}

func processStructuredAppend(results []*gozxing.Result) []*gozxing.Result {
	hasSA := false

	// first, check, if there is at least on SA result in the list
	for _, result := range results {
		metadata := result.GetResultMetadata()
		if _, ok := metadata[gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE]; ok {
			hasSA = true
			break
		}
	}
	if !hasSA {
		return results
	}

	// it is, second, split the lists and built a new result list
	newResults := make([]*gozxing.Result, 0)
	saResults := make([]*gozxing.Result, 0)
	for _, result := range results {
		metadata := result.GetResultMetadata()
		if _, ok := metadata[gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE]; ok {
			saResults = append(saResults, result)
		} else {
			newResults = append(newResults, result)
		}
	}
	// sort and concatenate the SA list items
	sort.Slice(saResults, newSAComparator(saResults))
	concatedText := make([]byte, 0)
	rawBytesLen := 0
	byteSegmentLength := 0
	for _, saResult := range saResults {
		concatedText = append(concatedText, []byte(saResult.GetText())...)
		rawBytesLen += len(saResult.GetRawBytes())
		metadata := saResult.GetResultMetadata()
		if byteSegments, ok := metadata[gozxing.ResultMetadataType_BYTE_SEGMENTS].([][]byte); ok {
			for _, segment := range byteSegments {
				byteSegmentLength += len(segment)
			}
		}
	}
	newRawBytes := make([]byte, rawBytesLen)
	newByteSegment := make([]byte, byteSegmentLength)
	newRawBytesIndex := 0
	byteSegmentIndex := 0
	for _, saResult := range saResults {
		copy(newRawBytes[newRawBytesIndex:], saResult.GetRawBytes())
		newRawBytesIndex += len(saResult.GetRawBytes())

		metadata := saResult.GetResultMetadata()
		if byteSegments, ok := metadata[gozxing.ResultMetadataType_BYTE_SEGMENTS].([][]byte); ok {
			for _, segment := range byteSegments {
				copy(newByteSegment[byteSegmentIndex:], segment)
				byteSegmentIndex += len(segment)
			}
		}
	}
	newResult := gozxing.NewResult(string(concatedText), newRawBytes, noPoints, gozxing.BarcodeFormat_QR_CODE)
	if byteSegmentLength > 0 {
		byteSegmentList := [][]byte{newByteSegment}
		newResult.PutMetadata(gozxing.ResultMetadataType_BYTE_SEGMENTS, byteSegmentList)
	}
	newResults = append(newResults, newResult)
	return newResults
}

func newSAComparator(results []*gozxing.Result) func(int, int) bool {
	return func(a, b int) bool {
		aNumber, _ := results[a].GetResultMetadata()[gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE].(int)
		bNumber, _ := results[b].GetResultMetadata()[gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE].(int)
		return aNumber < bNumber
	}
}
