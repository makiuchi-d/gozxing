package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestUPCEANExtensionSupport_decodeRowFail(t *testing.T) {
	row := gozxing.NewBitArray(20)

	// no guard pattern
	s := NewUPCEANExtensionSupport()
	_, e := s.decodeRow(10, row, 1)
	if e == nil {
		t.Fatalf("decodeRow must return error")
	}

	// guard pattern only
	row.Set(10)
	row.Set(12)
	row.Set(13)
	s = NewUPCEANExtensionSupport()
	_, e = s.decodeRow(10, row, 1)
	if e == nil {
		t.Fatalf("decodeRow must return error")
	}
}

func TestUPCEANExtensionSupport_decodeRowTwoSupport(t *testing.T) {
	row := gozxing.NewBitArray(30)
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

	s := NewUPCEANExtensionSupport()
	result, e := s.decodeRow(5, row, 0)
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
	if x, y := points[0].GetX(), points[0].GetY(); x != 3 || y != 5 {
		t.Fatalf("result point start = (%v,%v), expect (3,5)", x, y)
	}
	if x, y := points[1].GetX(), points[1].GetY(); x != 21 || y != 5 {
		t.Fatalf("result point end = (%v,%v), expect (21,5)", x, y)
	}
}

func TestUPCEANExtensionSupport_decodeRowFiveSupport(t *testing.T) {
	row := gozxing.NewBitArray(50)
	row.Set(1)
	row.Set(3)
	row.Set(4)
	// "56789", checksum=9(pattern:0x05=00101b)
	// [5-11] 5(L): 0110001
	row.Set(6)
	row.Set(7)
	row.Set(11)
	row.Set(13)
	// [14-20] 6(L): 0101111
	row.Set(15)
	row.Set(17)
	row.Set(18)
	row.Set(19)
	row.Set(20)
	row.Set(22)
	// [23-29] 7(G): 0010001
	row.Set(25)
	row.Set(29)
	row.Set(31)
	// [32-38] 8(L): 0110111
	row.Set(33)
	row.Set(34)
	row.Set(36)
	row.Set(37)
	row.Set(38)
	row.Set(40)
	// [41-47] 9(G): 0010111
	row.Set(43)
	row.Set(45)
	row.Set(46)
	row.Set(47)

	s := NewUPCEANExtensionSupport()
	result, e := s.decodeRow(7, row, 0)
	if e != nil {
		t.Fatalf("decodeRow return error, %v", e)
	}
	if txt := result.GetText(); txt != "56789" {
		t.Fatalf("decodeRow text = \"%v\", expect \"56789\"", txt)
	}
	if format := result.GetBarcodeFormat(); format != gozxing.BarcodeFormat_UPC_EAN_EXTENSION {
		t.Fatalf("result format = %v, expect UPC_EAN_EXTENSION", format)
	}
	meta := result.GetResultMetadata()
	price, ok := meta[gozxing.ResultMetadataType_SUGGESTED_PRICE]
	if !ok {
		t.Fatalf("metadata must contain key SUGGESTED_PRICE")
	}
	if price != "$67.89" {
		t.Fatalf("metadata[SUGGESTED_PRICE] = %v, expect $67.89", price)
	}
	points := result.GetResultPoints()
	if len(points) != 2 {
		t.Fatalf("result points length = %v, expect 2", len(points))
	}
	if x, y := points[0].GetX(), points[0].GetY(); x != 3 || y != 7 {
		t.Fatalf("result point start = (%v,%v), expect (3,7)", x, y)
	}
	if x, y := points[1].GetX(), points[1].GetY(); x != 48 || y != 7 {
		t.Fatalf("result point end = (%v,%v), expect (48,7)", x, y)
	}
}
