package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestUPCAReader_DecodeRow(t *testing.T) {
	reader := NewUPCAReader().(*upcAReader)
	row := gozxing.NewBitArray(105)

	_, e := reader.DecodeRow(10, row, nil)
	if e == nil {
		t.Fatalf("DecodeRow must be error")
	}

	// invalid UPCA (EAN13:1234567890128)
	for i, b := range "00010100100110111101001110101100010000101001000101010100100011101001110010110011011011001001000101000" {
		if b == '1' {
			row.Set(i)
		}
	}

	_, e = reader.DecodeRow(10, row, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("DecodeRow must be FormatException, %T", e)
	}

	// valid UPCA:123456789012
	for i, b := range "000101001100100100110111101010001101100010101111010101000100100100011101001110010110011011011001010000" {
		if b == '1' {
			row.Set(i)
		} else if row.Get(i) {
			row.Flip(i)
		}
	}

	result, e := reader.DecodeRow(10, row, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error, %v", e)
	}
	if format := result.GetBarcodeFormat(); format != gozxing.BarcodeFormat_UPC_A {
		t.Fatalf("DecodeRow format = %v, expect UPC_A", format)
	}
	if text := result.GetText(); text != "123456789012" {
		t.Fatalf("DecodeRow text = \"%v\",  expect \"123456789012\"", text)
	}
}

func TestUPCAReader_DecodeWithoutHint(t *testing.T) {
	bmp := testutil.NewBinaryBitmapFromFile("testdata/upca/2.png")

	result, e := NewUPCAReader().DecodeWithoutHints(bmp)
	if e != nil {
		t.Fatalf("DecodeWithoutHints returns error, %v", e)
	}
	if format := result.GetBarcodeFormat(); format != gozxing.BarcodeFormat_UPC_A {
		t.Fatalf("DecodeWithoutHints format = %v, expect UPC_A", format)
	}
	if text := result.GetText(); text != "036602301467" {
		t.Fatalf("DecodeWithoutHints text = \"%v\",  expect \"036602301467\"", text)
	}
}

func TestUPCAReader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/upca-1/
	reader := NewUPCAReader()
	format := gozxing.BarcodeFormat_UPC_A
	harder := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}

	// original zxing could't read too.
	bmp := testutil.NewBinaryBitmapFromFile("testdata/upca/1.png")
	_, e := reader.Decode(bmp, harder)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Decode \"testdata/upca/1.png\" must be NotFoundException, %T", e)
	}

	tests := []struct {
		file   string
		wants  string
		harder map[gozxing.DecodeHintType]interface{}
	}{
		{"testdata/upca/2.png", "036602301467", nil},
		{"testdata/upca/3.png", "070097025088", harder},
		// original zxing misread too.
		// {"testdata/upca/4.png", "070097025088", harder},
		{"testdata/upca/5.png", "070097025088", nil},
		{"testdata/upca/8.png", "071831007995", nil},
		{"testdata/upca/9.png", "071831007995", nil},
		{"testdata/upca/10.png", "027011006951", nil},
		{"testdata/upca/11.png", "027011006951", nil},
		{"testdata/upca/12.png", "781735802045", harder},
		{"testdata/upca/13.png", "781735802045", nil},
		{"testdata/upca/16.png", "456314319671", nil},
		{"testdata/upca/17.png", "434704791429", nil},
		{"testdata/upca/18.png", "024543136538", harder},
		// gozxing misread this image. (I don't know why....)
		// {"testdata/upca/19.png", "024543136538", harder},
		// original zxing could't read too.
		// {"testdata/upca/20.png", "752919460009", harder},
		{"testdata/upca/21.png", "752919460009", nil},
		{"testdata/upca/27.png", "606949762520", nil},
		{"testdata/upca/28.png", "061869053712", nil},
		{"testdata/upca/29.png", "619659023935", nil},
		{"testdata/upca/35.png", "045496442736", nil},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, test.harder, nil)
	}
}

func TestUpcAReader_getBarcodeFormat(t *testing.T) {
	reader := NewUPCAReader().(*upcAReader)
	wants := gozxing.BarcodeFormat_UPC_A
	if f := reader.getBarcodeFormat(); f != wants {
		t.Fatalf("getBarcodeFormat() = %v, wants =%v", f, wants)
	}
}
