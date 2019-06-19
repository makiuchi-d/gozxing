package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestMultiFormatUPCEANReader_DecodeRow(t *testing.T) {
	reader := NewMultiFormatUPCEANReader(nil).(*multiFormatUPCEANReader)

	row := gozxing.NewBitArray(30)

	_, e := reader.DecodeRow(0, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("error must be NotFoundException, %T, %+v", e, e)
	}

	row = testutil.NewBitArrayFromString("0000010100000000000000000")
	_, e = reader.DecodeRow(0, row, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("error must be NotFoundException, %T, %+v", e, e)
	}

	// valid UPCA:123456789012
	row = testutil.NewBitArrayFromString(
		"000101001100100100110111101010001101100010101111010101000100100100011101001110010110011011011001010000")
	r, e := reader.DecodeRow(0, row, nil)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %+v", e)
	}
	if txt, wants := r.GetText(), "0123456789012"; txt != wants {
		t.Fatalf("result text = \"%v\", wants \"%v\"", txt, wants)
	}
	if format, wants := r.GetBarcodeFormat(), gozxing.BarcodeFormat_EAN_13; format != wants {
		t.Fatalf("result format = %v, wants %v", format, wants)
	}

	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_POSSIBLE_FORMATS: []gozxing.BarcodeFormat{
			gozxing.BarcodeFormat_EAN_13, gozxing.BarcodeFormat_UPC_A,
		},
	}
	r, e = reader.DecodeRow(0, row, hints)
	if e != nil {
		t.Fatalf("DecodeRow returns error: %+v", e)
	}
	if txt, wants := r.GetText(), "123456789012"; txt != wants {
		t.Fatalf("result text = \"%v\", wants \"%v\"", txt, wants)
	}
	if format, wants := r.GetBarcodeFormat(), gozxing.BarcodeFormat_UPC_A; format != wants {
		t.Fatalf("result format = %v, wants %v", format, wants)
	}
}

func TestMultiFormatUPCEANReader(t *testing.T) {
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_POSSIBLE_FORMATS: []gozxing.BarcodeFormat{
			gozxing.BarcodeFormat_EAN_13,
			gozxing.BarcodeFormat_EAN_8,
			gozxing.BarcodeFormat_UPC_A,
			gozxing.BarcodeFormat_UPC_E,
		},
	}
	reader := NewMultiFormatUPCEANReader(hints)

	tests := []struct {
		file   string
		text   string
		format gozxing.BarcodeFormat
	}{
		{"testdata/ean13/1.png", "8413000065504", gozxing.BarcodeFormat_EAN_13},
		{"testdata/ean8/1.png", "48512343", gozxing.BarcodeFormat_EAN_8},
		{"testdata/upca/2.png", "036602301467", gozxing.BarcodeFormat_UPC_A},
		{"testdata/upce/1.png", "01234565", gozxing.BarcodeFormat_UPC_E},
	}
	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			img := testutil.NewBinaryBitmapFromFile(test.file)
			r, e := reader.Decode(img, hints)
			if e != nil {
				t.Fatalf("reader.Decode returns error: %+v", e)
			}
			if text := r.GetText(); text != test.text {
				t.Fatalf("result text = \"%v\", wants \"%v\"", text, test.text)
			}
			if format := r.GetBarcodeFormat(); format != test.format {
				t.Fatalf("result format = %v, wants %v", format, test.format)
			}
		})
	}
}
