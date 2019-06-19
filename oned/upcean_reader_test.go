package oned

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestUPCEANReader_findStartGuardPattern(t *testing.T) {
	row := gozxing.NewBitArray(10)

	_, e := upceanReader_findStartGuardPattern(row)
	if e == nil {
		t.Fatalf("findStartGuardPattern must be error")
	}

	row.Set(3)
	row.Set(5)
	startRange, e := upceanReader_findStartGuardPattern(row)
	if e != nil {
		t.Fatalf("findStartGuardPattern returns error, %v", e)
	}
	if !reflect.DeepEqual(startRange, []int{3, 6}) {
		t.Fatalf("findStartGuardPattern = %v, expect [3 6]", startRange)
	}
}

func TestUPCEANReader_checkChecksum(t *testing.T) {
	r, e := upceanReader_checkChecksum("")
	if e != nil {
		t.Fatalf("checkChecksum(\"\") returns error, %v", e)
	}
	if r != false {
		t.Fatalf("checkChecksum(\"\") =%v, expect false", r)
	}

	_, e = upceanReader_checkChecksum("01234567890a3")
	if e == nil {
		t.Fatalf("checkChecksum(\"01234567890a3\") must be error")
	}

	_, e = upceanReader_checkChecksum("0123456789a03")
	if e == nil {
		t.Fatalf("checkChecksum(\"01234567890a3\") must be error")
	}

	r, e = upceanReader_checkChecksum("0123456789012")
	if e != nil {
		t.Fatalf("checkChecksum(\"0123456789012\") returns error, %v", e)
	}
	if r != true {
		t.Fatalf("checkChecksum(\"0123456789012\") must be true")
	}

	r, e = upceanReader_checkChecksum("0123456789013")
	if e != nil {
		t.Fatalf("checkChecksum(\"0123456789013\") returns error, %v", e)
	}
	if r != false {
		t.Fatalf("checkChecksum(\"0123456789013\") must be false")
	}
}

func TestUPCEANReader_decodeEnd(t *testing.T) {
	row := gozxing.NewBitArray(10)
	row.Set(2)
	row.Set(5)
	row.Set(7)

	_, e := upceanReader_decodeEnd(row, 7)
	if e == nil {
		t.Fatalf("decodeEnd must be error")
	}

	endRange, e := upceanReader_decodeEnd(row, 4)
	if e != nil {
		t.Fatalf("decodeEnd returns error, %v", e)
	}
	if !reflect.DeepEqual(endRange, []int{5, 8}) {
		t.Fatalf("decodeEnd = %v, expect [5 8]", endRange)
	}
}

func TestUPCEANReader_findGuardPattern(t *testing.T) {
	pattern := []int{1, 3, 2, 1}
	row := gozxing.NewBitArray(20)
	// pattern {1,1,1,4,3,1,...} (start with white)
	// matches: [2-11] = {1,4,3,1} avgVariance:22.2%,
	row.Set(1)
	row.Set(3)
	row.Set(4)
	row.Set(5)
	row.Set(6)
	row.Set(10)
	prange, e := upceanReader_findGuardPattern(row, 0, true, pattern)
	if e != nil {
		t.Fatalf("findGuardPattern returns error, %v", e)
	}
	if !reflect.DeepEqual(prange, []int{2, 11}) {
		t.Fatalf("findGuardPattern = %v, expect [2 11]", prange)
	}
}

func TestUPCEANReader_decodeDigit(t *testing.T) {
	row := gozxing.NewBitArray(10)
	counters := make([]int, 4)

	_, e := upceanReader_decodeDigit(row, counters, 0, UPCEANReader_L_PATTERNS)
	if e == nil {
		t.Fatalf("decodeDigit must be error")
	}

	// {1,1,7,1}
	row.Set(1)
	row.Set(9)
	_, e = upceanReader_decodeDigit(row, counters, 0, UPCEANReader_L_PATTERNS)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeDigit must be NotFoundException, %T", e)
	}

	// {1,1,3,2} = 4
	row.Set(5)
	row.Set(6)
	i, e := upceanReader_decodeDigit(row, counters, 0, UPCEANReader_L_PATTERNS)
	if e != nil {
		t.Fatalf("decodeDigit returns error, %v", e)
	}
	if i != 4 {
		t.Fatalf("decodeDigit = %v, expect 4", i)
	}
}

type testMiddleDecoder struct {
	*upceanReader
	data  []byte
	ean13 bool
}

func newTestMiddleDecoder(data []byte, ean13 bool) *testMiddleDecoder {
	this := &testMiddleDecoder{
		data:  data,
		ean13: ean13,
	}
	this.upceanReader = newUPCEANReader(this)
	return this
}

func (t *testMiddleDecoder) decodeMiddle(row *gozxing.BitArray, startRange []int, result []byte) (int, []byte, error) {
	if len(t.data) == 0 {
		return 0, result, gozxing.NewFormatException()
	}
	result = append(result, t.data...)
	return startRange[1] + len(t.data), result, nil
}
func (t *testMiddleDecoder) getBarcodeFormat() gozxing.BarcodeFormat {
	if t.ean13 {
		return gozxing.BarcodeFormat_EAN_13
	}
	return gozxing.BarcodeFormat_EAN_8
}
func (this *testMiddleDecoder) decodeEnd(row *gozxing.BitArray, endStart int) ([]int, error) {
	return upceanReader_decodeEnd(row, endStart)
}
func (this *testMiddleDecoder) checkChecksum(s string) (bool, error) {
	return upceanReader_checkChecksum(s)
}

func TestUPCEANReader_DecodeRow(t *testing.T) {
	// no start guard
	row := gozxing.NewBitArray(13)
	reader := newTestMiddleDecoder([]byte{}, false)
	_, e := reader.DecodeRow(5, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// decodeMiddle error
	row.Set(3)
	row.Set(5)
	_, e = reader.DecodeRow(5, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// no end guard
	data := []byte{'0', '1', '2'}
	reader = newTestMiddleDecoder(data, false)
	_, e = reader.DecodeRow(5, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// less end quiet
	row.Set(9)
	row.Set(11)
	_, e = reader.DecodeRow(5, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// black on end quiet
	row = gozxing.NewBitArray(25)
	row.Set(3)
	row.Set(5)
	row.Set(9)
	row.Set(11)
	row.Set(13)
	_, e = reader.DecodeRow(5, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	// less result length
	row.Flip(13)
	_, e = reader.DecodeRow(5, row, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("DecodeRow must be FormatException, %T", e)
	}

	// checksum calc error
	data = []byte("a1234567")
	reader = newTestMiddleDecoder(data, false)
	row.Set(15)
	row.Set(17)
	_, e = reader.DecodeRow(5, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// checksum error
	data = []byte("01234567")
	reader = newTestMiddleDecoder(data, false)
	_, e = reader.DecodeRow(5, row, nil)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("DecodeRow must be ChecksumException, %T", e)
	}

	// success
	data = []byte("01234565")
	reader = newTestMiddleDecoder(data, false)
	result, e := reader.DecodeRow(5, row, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error, %v", e)
	}
	if str := result.GetText(); str != "01234565" {
		t.Fatalf("result text = \"%v\", expect \"01234565\"", str)
	}
	rps := result.GetResultPoints()
	if x, y := rps[0].GetX(), rps[0].GetY(); x != 4.5 || y != 5 {
		t.Fatalf("result point[0] = (%v, %v), expect (4.5, 5)", x, y)
	}
	if x, y := rps[1].GetX(), rps[1].GetY(); x != 16.5 || y != 5 {
		t.Fatalf("result point[1] = (%v, %v), expect (16.5, 5)", x, y)
	}
	if format := result.GetBarcodeFormat(); format != gozxing.BarcodeFormat_EAN_8 {
		t.Fatalf("result format = %v, expect EAN_8", format)
	}
}

func TestUPCEANReader_DecodeRowWithResultPointCallback(t *testing.T) {
	row := gozxing.NewBitArray(25)
	row.Set(3)
	row.Set(5)
	row.Set(15)
	row.Set(17)
	data := []byte("01234565")
	reader := newTestMiddleDecoder(data, false)

	points := make([]gozxing.ResultPoint, 0)
	callback := gozxing.ResultPointCallback(func(p gozxing.ResultPoint) {
		points = append(points, p)
	})

	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK] = callback
	_, e := reader.DecodeRow(7, row, hints)
	if e != nil {
		t.Fatalf("DecodeRow returns error, %v", e)
	}

	if len(points) != 3 {
		t.Fatalf("callbacked points count = %v, expect 3", len(points))
	}
	if x, y := points[0].GetX(), points[0].GetY(); x != 4.5 || y != 7 {
		t.Fatalf("result point[0] = (%v, %v), expect (4.5, 7)", x, y)
	}
	if x, y := points[1].GetX(), points[1].GetY(); x != 14 || y != 7 {
		t.Fatalf("result point[1] = (%v, %v), expect (14, 7)", x, y)
	}
	if x, y := points[2].GetX(), points[2].GetY(); x != 16.5 || y != 7 {
		t.Fatalf("result point[2] = (%v, %v), expect (16.5, 7)", x, y)
	}
}

func TestUPCEANReader_DecodeRowWithExtension(t *testing.T) {
	row := gozxing.NewBitArray(55)
	row.Set(3)
	row.Set(5)
	row.Set(15)
	row.Set(17)
	// extension2
	row.Set(31)
	row.Set(33)
	row.Set(34)
	row.Set(37)
	row.Set(38)
	row.Set(41)
	row.Set(43)
	row.Set(46)
	row.Set(49)
	row.Set(50)

	data := []byte("01234565")
	reader := newTestMiddleDecoder(data, false)

	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_ALLOWED_EAN_EXTENSIONS] = []int{5}
	_, e := reader.DecodeRow(7, row, hints)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeRow must be NotFoundException, %T", e)
	}

	hints[gozxing.DecodeHintType_ALLOWED_EAN_EXTENSIONS] = []int{2, 5}
	result, e := reader.DecodeRow(7, row, hints)
	if e != nil {
		t.Fatalf("DecodeRow returns error, %v", e)
	}
	meta := result.GetResultMetadata()
	ext, ok := meta[gozxing.ResultMetadataType_UPC_EAN_EXTENSION]
	if !ok {
		t.Fatalf("metadata must have key UPC_EAN_EXTENSION")
	}
	if s := ext.(string); s != "12" {
		t.Fatalf("metadata[UPC_EAN_EXTENSION] = \"%v\", expect \"12\"", s)
	}
	issue, ok := meta[gozxing.ResultMetadataType_ISSUE_NUMBER]
	if !ok {
		t.Fatalf("metadata must have key ISSUE_NUMBER")
	}
	if i := issue.(int); i != 12 {
		t.Fatalf("metadata[UPC_EAN_EXTENSION] = %v, expect 12", i)
	}
}

func TestUPCEANReader_DecodeRowEAN13(t *testing.T) {
	row := gozxing.NewBitArray(36)
	row.Set(3)
	row.Set(5)
	row.Set(29)
	row.Set(31)

	data := []byte("0123456789012")
	reader := newTestMiddleDecoder(data, true)
	result, e := reader.DecodeRow(7, row, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error, %v", e)
	}
	meta := result.GetResultMetadata()
	country, ok := meta[gozxing.ResultMetadataType_POSSIBLE_COUNTRY]
	if !ok {
		t.Fatalf("metadata must have key POSSIBLE_COUNTRY")
	}
	if country != "US/CA" {
		t.Fatalf("metadata[POSSIBLE_COUNTRY] = %v, expect US/CA", country)
	}
}
