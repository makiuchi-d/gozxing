package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestUPCEANExtension2Support_ParseExtensionString(t *testing.T) {
	e2s := NewUPCEANExtension2Support()

	if r := e2s.parseExtensionString(""); r != nil {
		t.Fatalf("parseExtensionString(\"\") = %v, expect nil", r)
	}

	if r := e2s.parseExtensionString("ab"); r != nil {
		t.Fatalf("parseExtensionString(\"ab\") = %v, expect nil", r)
	}

	r := e2s.parseExtensionString("01")
	n, ok := r[gozxing.ResultMetadataType_ISSUE_NUMBER]
	if !ok {
		t.Fatalf("parseExtensionString(\"01\") must contain key ISSUE_NUMBER, %v", r)
	}
	if n != 1 {
		t.Fatalf("parseExtensionString(\"01\") issue number = %v, expect 1", n)
	}
}

func TestUPCEANExtension2Support_decodeMiddle(t *testing.T) {
	row := gozxing.NewBitArray(30)

	// start pattern
	row.Set(1)
	row.Set(3)
	row.Set(4)
	startRange := []int{1, 5}

	e2s := NewUPCEANExtension2Support()
	_, e := e2s.decodeMiddle(row, startRange)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeMiddle must be NotFoundException, %T", e)
	}

	// [5-11] 0(L): 0001101
	row.Set(8)
	row.Set(9)
	row.Set(11)
	e2s = NewUPCEANExtension2Support()
	_, e = e2s.decodeMiddle(row, startRange)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeMiddle must be NotFoundException, %T", e)
	}

	// [12-13] separator: 01
	row.Set(13)
	// [14-20] 1(G): 0110011
	row.Set(15)
	row.Set(16)
	row.Set(19)
	row.Set(20)
	e2s = NewUPCEANExtension2Support()
	offset, e := e2s.decodeMiddle(row, startRange)
	if e != nil {
		t.Fatalf("decodeMiddle returns error, %v", e)
	}
	if offset != 21 {
		t.Fatalf("decodeMiddle offset = %v, expect 21", offset)
	}
	if s := string(e2s.decodeRowStringBuffer); s != "01" {
		t.Fatalf("decodeMiddle string = \"%v\", expect \"01\"", s)
	}

	// parity error
	// [14-20] 1(L): 0011001
	row.Flip(15)
	row.Flip(17)
	row.Flip(19)
	e2s = NewUPCEANExtension2Support()
	_, e = e2s.decodeMiddle(row, startRange)
	if _, ok := e.(gozxing.ReaderException); !ok {
		t.Fatalf("decodeMiddle must be ReaderException, %T", e)
	}
}

func TestUPCEANExtension2Support_decodeRow(t *testing.T) {
	row := gozxing.NewBitArray(30)
	startRange := []int{1, 5}

	e2s := NewUPCEANExtension2Support()
	_, e := e2s.decodeRow(10, row, startRange)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeRow must be NotFoundException, %T", e)
	}

	// "12": 01011 0011001 01 0010011
	row.Set(1)
	row.Set(3)
	row.Set(4)
	row.Set(7)
	row.Set(8)
	row.Set(11)
	row.Set(13)
	row.Set(16)
	row.Set(19)
	row.Set(20)

	e2s = NewUPCEANExtension2Support()
	result, e := e2s.decodeRow(10, row, startRange)
	if e != nil {
		t.Fatalf("decodeRow returns error, %v", e)
	}
	if txt := result.GetText(); txt != "12" {
		t.Fatalf("result text = \"%v\", expect \"12\"", txt)
	}
	if format := result.GetBarcodeFormat(); format != gozxing.BarcodeFormat_UPC_EAN_EXTENSION {
		t.Fatalf("result format = %v, expect UPC_EAN_EXTENSION", format)
	}
	meta := result.GetResultMetadata()
	num, ok := meta[gozxing.ResultMetadataType_ISSUE_NUMBER]
	if !ok {
		t.Fatalf("result metadata must contain key ISSUE_NUMBER, %v", meta)
	}
	if num != 12 {
		t.Fatalf("result metadata[ISSUE_NUMBER] = %v, expect 12", num)
	}
	points := result.GetResultPoints()
	if len(points) != 2 {
		t.Fatalf("result points length = %v, expect 2", len(points))
	}
	if x, y := points[0].GetX(), points[0].GetY(); x != 3 || y != 10 {
		t.Fatalf("result point start = (%v,%v), expect (3,10)", x, y)
	}
	if x, y := points[1].GetX(), points[1].GetY(); x != 21 || y != 10 {
		t.Fatalf("result point end = (%v,%v), expect (21,10)", x, y)
	}
}
