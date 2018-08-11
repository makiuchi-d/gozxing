package oned

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

func readFile(reader *OneDReader, filename string, hints map[gozxing.DecodeHintType]interface{}) (*gozxing.Result, error) {
	file, e := os.Open(filename)
	if e != nil {
		return nil, e
	}
	img, _, e := image.Decode(file)
	if e != nil {
		return nil, e
	}
	src := gozxing.NewLuminanceSourceFromImage(img)
	bmp, _ := gozxing.NewBinaryBitmap(common.NewHybridBinarizer(src))
	if e != nil {
		return nil, e
	}

	return reader.Decode(bmp, hints)
}

func testFile(t *testing.T, reader *OneDReader, file, expect string, hints map[gozxing.DecodeHintType]interface{}) {
	result, e := readFile(reader, file, hints)
	if e != nil {
		t.Fatalf("testFail(%v) readFile failed: %v", file, e)
	}
	if txt := result.GetText(); txt != expect {
		t.Fatalf("testFile(%v) = %v, expect %v", file, txt, expect)
	}
}

func TestRecordPattern(t *testing.T) {
	row := gozxing.NewBitArray(20)
	counters := make([]int, 4)

	e := recordPattern(row, 20, counters)
	if e == nil {
		t.Fatalf("recordPattern must be error")
	}

	// [1-7] 0011001 (UPC/EAN 1L)
	row.Set(3)
	row.Set(4)
	row.Set(7)

	// [11-17] 1001000 (UPC/EAN 8R)
	row.Set(11)
	row.Set(14)
	row.Set(18)

	e = recordPattern(row, 1, counters)
	if e != nil {
		t.Fatalf("recordPattern returns error, %v", e)
	}
	if !reflect.DeepEqual(counters, []int{2, 2, 2, 1}) {
		t.Fatalf("recordPattern = %v, expect [2 2 2 1]", counters)
	}

	e = recordPattern(row, 11, counters)
	if e != nil {
		t.Fatalf("recordPattern returns error, %v", e)
	}
	if !reflect.DeepEqual(counters, []int{1, 2, 1, 3}) {
		t.Fatalf("recordPattern = %v, expect [1 2 1 3]", counters)
	}
}

func TestRecordPatternInReverse(t *testing.T) {
	row := gozxing.NewBitArray(30)
	counters := make([]int, 4)

	row.Set(0)
	// [1-7] 0011001 (1L)
	row.Set(3)
	row.Set(4)
	row.Set(7)
	// [11-17] 1001000 (8R)
	row.Set(11)
	row.Set(14)
	row.Set(18)

	e := recordPatternInReverse(row, 3, counters)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("recordPatternInReverse must be NotFoundException, %T", e)
	}

	e = recordPatternInReverse(row, 8, counters)
	if e != nil {
		t.Fatalf("recordPattern returns error, %v", e)
	}
	if !reflect.DeepEqual(counters, []int{2, 2, 2, 1}) {
		t.Fatalf("recordPattern = %v, expect [2 2 2 1]", counters)
	}

	e = recordPatternInReverse(row, 18, counters)
	if e != nil {
		t.Fatalf("recordPattern returns error, %v", e)
	}
	if !reflect.DeepEqual(counters, []int{1, 2, 1, 3}) {
		t.Fatalf("recordPattern = %v, expect [1 2 1 3]", counters)
	}
}

type testBitSource struct {
	gozxing.LuminanceSourceBase
	bits string
}

func newTestBitSource(height int, bits string) gozxing.LuminanceSource {
	return &testBitSource{
		gozxing.LuminanceSourceBase{len(bits), height},
		bits,
	}
}
func (this *testBitSource) GetRow(y int, row []byte) ([]byte, error) {
	w := len(this.bits)
	if w <= 0 {
		return nil, fmt.Errorf("GetRow error: width=%v", w)
	}

	if len(row) < w {
		row = make([]byte, w)
	}
	for i, b := range this.bits {
		if b == '1' {
			row[i] = 0 //black
		} else {
			row[i] = 255 // white
		}
	}
	return row, nil
}
func (this *testBitSource) GetMatrix() []byte {
	w := this.GetWidth()
	h := this.GetHeight()
	matrix := make([]byte, w*h)
	row, _ := this.GetRow(0, nil)
	for y := 0; y < h; y++ {
		copy(matrix[y*w:], row)
	}
	return matrix
}
func (this *testBitSource) Invert() gozxing.LuminanceSource {
	return gozxing.LuminanceSourceInvert(this)
}
func (this *testBitSource) String() string {
	return gozxing.LuminanceSourceString(this)
}

func TestOneDReader_doDecode(t *testing.T) {
	src := newTestBitSource(1, "")
	bmp, _ := gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))

	reader := NewEAN8Reader().(*OneDReader)

	_, e := reader.doDecode(bmp, nil)
	if e == nil {
		t.Fatalf("doDecode must be error")
	}

	src = newTestBitSource(1, "0000")
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
	_, e = reader.doDecode(bmp, nil)
	if e == nil {
		t.Fatalf("doDecode must be error")
	}

	// reverse with resultpointcallback
	points := make([]gozxing.ResultPoint, 0)
	callback := func(p gozxing.ResultPoint) { points = append(points, p) }
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER:                 true,
		gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK: gozxing.ResultPointCallback(callback),
	}

	// upside down of "12345670"
	src = newTestBitSource(10,
		"000010101001110010001000010101110010101011000101011110110010010011001010000")
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
	r, e := reader.doDecode(bmp, hints)
	if e != nil {
		t.Fatalf("doDecode returns error, %v", e)
	}
	if txt := r.GetText(); txt != "12345670" {
		t.Fatalf("doDecode text = \"%v\", expect \"12345670\"", txt)
	}
	points = r.GetResultPoints()
	if x, y := points[0].GetX(), points[0].GetY(); x != 68.5 || y != 5 {
		t.Fatalf("doDecode resultpoint[0] = (%v,%v), expect (68.5,5)", x, y)
	}
	if x, y := points[1].GetX(), points[1].GetY(); x != 4.5 || y != 5 {
		t.Fatalf("doDecode resultpoint[0] = (%v,%v), expect (4.5,5)", x, y)
	}
}

func TestOneDReader_Reset(t *testing.T) {
	NewEAN8Reader().Reset() // do nothing
}

type testUnrotatableBitSource struct {
	*testBitSource
}

func (*testUnrotatableBitSource) IsRotateSupported() bool {
	return true // but, not implement RotateCounterClockwise()
}

func TestOneDReader_DecodeFail(t *testing.T) {
	reader := NewEAN8Reader()

	src := newTestBitSource(1, "")
	bmp, _ := gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
	_, e := reader.DecodeWithoutHints(bmp)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	src = newTestBitSource(1, "01010")
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}
	_, e = reader.Decode(bmp, hints)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	src = &testUnrotatableBitSource{src.(*testBitSource)}
	bmp, _ = gozxing.NewBinaryBitmap(common.NewGlobalHistgramBinarizer(src))
	_, e = reader.Decode(bmp, hints)
	if e == nil {
		t.Fatalf("Decode must be error")
	}
}
