package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestUPCAReader_decodeRow(t *testing.T) {
	reader := NewUPCAReader().(*upcAReader)
	row := gozxing.NewBitArray(105)

	_, e := reader.decodeRow(10, row, nil)
	if e == nil {
		t.Fatalf("decodeRow must be error")
	}

	// invalid UPCA (EAN13:1234567890128)
	for i, b := range "00010100100110111101001110101100010000101001000101010100100011101001110010110011011011001001000101000" {
		if b == '1' {
			row.Set(i)
		}
	}

	_, e = reader.decodeRow(10, row, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeRow must be FormatException, %T", e)
	}

	// valid UPCA:123456789012
	for i, b := range "000101001100100100110111101010001101100010101111010101000100100100011101001110010110011011011001010000" {
		if b == '1' {
			row.Set(i)
		} else if row.Get(i) {
			row.Flip(i)
		}
	}

	result, e := reader.decodeRow(10, row, nil)
	if e != nil {
		t.Fatalf("decodeRow returns error, %v", e)
	}
	if format := result.GetBarcodeFormat(); format != gozxing.BarcodeFormat_UPC_A {
		t.Fatalf("decodeRow format = %v, expect UPC_A", format)
	}
	if text := result.GetText(); text != "123456789012" {
		t.Fatalf("decodeRow text = \"%v\",  expect \"123456789012\"", text)
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

	testFile(t, reader, "testdata/upca/2.png", "036602301467", format, nil)
	testFile(t, reader, "testdata/upca/3.png", "070097025088", format, harder)
	// original zxing misread too.
	// testFile(t, reader, "testdata/upca/4.png", "070097025088", format, harder)
	testFile(t, reader, "testdata/upca/5.png", "070097025088", format, nil)
	testFile(t, reader, "testdata/upca/8.png", "071831007995", format, nil)
	testFile(t, reader, "testdata/upca/9.png", "071831007995", format, nil)
	testFile(t, reader, "testdata/upca/10.png", "027011006951", format, nil)
	testFile(t, reader, "testdata/upca/11.png", "027011006951", format, nil)
	testFile(t, reader, "testdata/upca/12.png", "781735802045", format, harder)
	testFile(t, reader, "testdata/upca/13.png", "781735802045", format, nil)
	testFile(t, reader, "testdata/upca/16.png", "456314319671", format, nil)
	testFile(t, reader, "testdata/upca/17.png", "434704791429", format, nil)
	testFile(t, reader, "testdata/upca/18.png", "024543136538", format, harder)
	// gozxing misread this image. (I don't know why....)
	// testFile(t, reader, "testdata/upca/19.png", "024543136538", format, harder)
	// original zxing could't read too.
	// testFile(t, reader, "testdata/upca/20.png", "752919460009", format, harder)
	testFile(t, reader, "testdata/upca/21.png", "752919460009", format, nil)
	testFile(t, reader, "testdata/upca/27.png", "606949762520", format, nil)
	testFile(t, reader, "testdata/upca/28.png", "061869053712", format, nil)
	testFile(t, reader, "testdata/upca/29.png", "619659023935", format, nil)
	testFile(t, reader, "testdata/upca/35.png", "045496442736", format, nil)
}
