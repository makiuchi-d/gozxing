package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func testUPCEANExtension5Support_parseExtension5String(t testing.TB, raw, expect string) {
	t.Helper()
	e5s := NewUPCEANExtension5Support()
	str := e5s.parseExtension5String(raw)
	if str != expect {
		t.Fatalf("parseExtension5String(%s) = \"%v\", expect \"%v\"", raw, str, expect)
	}
}

func TestUPCEANExtension5Support_parseExtension5String(t *testing.T) {
	testUPCEANExtension5Support_parseExtension5String(t, "01203", "£12.03")
	testUPCEANExtension5Support_parseExtension5String(t, "50789", "$7.89")

	testUPCEANExtension5Support_parseExtension5String(t, "90000", "")
	testUPCEANExtension5Support_parseExtension5String(t, "99991", "0.00")
	testUPCEANExtension5Support_parseExtension5String(t, "99990", "Used")

	testUPCEANExtension5Support_parseExtension5String(t, "90123", "1.23")
	testUPCEANExtension5Support_parseExtension5String(t, "40123", "1.23")

	testUPCEANExtension5Support_parseExtension5String(t, "4----", "")
}

func TestUPCEANExtension5Support_parseExtensionString(t *testing.T) {
	e5s := NewUPCEANExtension5Support()

	meta := e5s.parseExtensionString("")
	if meta != nil {
		t.Fatalf("parseExtensionString(\"\") = %v, expect nil", meta)
	}

	meta = e5s.parseExtensionString("90000")
	if meta != nil {
		t.Fatalf("parseExtensionString(\"\") = %v, expect nil", meta)
	}

	meta = e5s.parseExtensionString("01234")
	price, ok := meta[gozxing.ResultMetadataType_SUGGESTED_PRICE]
	if !ok {
		t.Fatalf("metadata must contain key SUGGESTED_PRICE, %v", meta)
	}
	if price != "£12.34" {
		t.Fatalf("metadata[SUGGESTED_PRICE] = \"%v\", expect \"£12.34\"", price)
	}
}

func TestUPCEANExtension5Support_determineCheckDigit(t *testing.T) {
	e5s := NewUPCEANExtension5Support()

	_, e := e5s.determineCheckDigit(0)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("determineCheckDigit(0) must be NotFoundException, %T", e)
	}

	for expect, pattern := range []int{0x18, 0x14, 0x12, 0x11, 0x0C, 0x06, 0x03, 0x0A, 0x09, 0x05} {
		digit, e := e5s.determineCheckDigit(pattern)
		if e != nil {
			t.Fatalf("determineCheckDigit(%v) returns error, %v", pattern, e)
		}
		if digit != expect {
			t.Fatalf("determineCheckDigit(%v) = %v, expect %v", pattern, digit, expect)
		}
	}
}

func testUPCEANExtension5Support_extensionChecksum(t testing.TB, str string, expect int) {
	t.Helper()
	e5s := NewUPCEANExtension5Support()
	sum := e5s.extensionChecksum(str)
	if sum != expect {
		t.Fatalf("extensionChecksum(%v) = %v, expect %v", str, sum, expect)
	}
}

func TestUPCEANExtension5Support_extensionChecksum(t *testing.T) {
	testUPCEANExtension5Support_extensionChecksum(t, "00000", 0)
	testUPCEANExtension5Support_extensionChecksum(t, "12345", 1)
	testUPCEANExtension5Support_extensionChecksum(t, "67890", 6)
	testUPCEANExtension5Support_extensionChecksum(t, "99999", 3)
	testUPCEANExtension5Support_extensionChecksum(t, "56789", 9)
}

func TestUPCEANExtension5Support_decodeMiddle(t *testing.T) {
	row := gozxing.NewBitArray(50)

	//start pattern
	row.Set(1)
	row.Set(3)
	row.Set(4)
	startRange := []int{1, 5}

	e5s := NewUPCEANExtension5Support()
	_, e := e5s.decodeMiddle(row, startRange)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeMiddle must be NotFoundException, %T", e)
	}

	// "12345", parity=1 => 0x14(10100b)
	// [5-11] 1(G): 0110011 01
	row.Set(6)
	row.Set(7)
	row.Set(10)
	row.Set(11)

	e5s = NewUPCEANExtension5Support()
	_, e = e5s.decodeMiddle(row, startRange)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeMiddle must be NotFoundException, %T", e)
	}

	row.Set(13) // separator

	// [14-20] 2(L): 0010011 01
	row.Set(16)
	row.Set(19)
	row.Set(20)
	row.Set(22)
	// [23-29] 3(G): 0100001 01
	row.Set(24)
	row.Set(29)
	row.Set(31)
	// [32-38] 4(L): 0100011 01
	row.Set(33)
	row.Set(37)
	row.Set(38)
	row.Set(40)
	// [41-47] 5(L): 0110001
	row.Set(42)
	row.Set(43)
	row.Set(47)

	e5s = NewUPCEANExtension5Support()
	offset, e := e5s.decodeMiddle(row, startRange)
	if e != nil {
		t.Fatalf("decodeMiddle returns error, %v", e)
	}
	if offset != 48 {
		t.Fatalf("decodeMiddle offset = %v, expect 48", offset)
	}
	if str := string(e5s.decodeRowStringBuffer); str != "12345" {
		t.Fatalf("decodeMiddle string = \"%v\", expect \"12345\"", str)
	}

	// invalid checksum pattern: 0x14=>0x04
	// [5-11] 1(G):0110011 => 1(L):0011001
	row.Flip(6)
	row.Flip(8)
	row.Flip(10)
	e5s = NewUPCEANExtension5Support()
	_, e = e5s.decodeMiddle(row, startRange)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeMiddle must be NotFoundException, %T", e)
	}

	// invalid checksum value: 9=0x05
	// [41-47] 5(L):0110001 => 5(G):0111001
	row.Flip(44)
	e5s = NewUPCEANExtension5Support()
	_, e = e5s.decodeMiddle(row, startRange)
	if _, ok := e.(gozxing.ReaderException); !ok {
		t.Fatalf("decodeMiddle must be ReaderException, %T", e)
	}
}

func TestUPCEANExtension5Support_decodeRow(t *testing.T) {
	row := gozxing.NewBitArray(50)
	row.Set(1)
	row.Set(3)
	row.Set(4)
	startRange := []int{1, 5}

	e5s := NewUPCEANExtension5Support()
	_, e := e5s.decodeRow(10, row, startRange)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeRow must be NotFoundException, %T", e)
	}

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

	e5s = NewUPCEANExtension5Support()
	result, e := e5s.decodeRow(10, row, startRange)
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
	if x, y := points[0].GetX(), points[0].GetY(); x != 3 || y != 10 {
		t.Fatalf("result point start = (%v,%v), expect (3,10)", x, y)
	}
	if x, y := points[1].GetX(), points[1].GetY(); x != 48 || y != 10 {
		t.Fatalf("result point end = (%v,%v), expect (48,10)", x, y)
	}
}
