package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestCodabarEncoder_encode(t *testing.T) {
	enc := codabarEncoder{}

	failtests := []string{
		"A123T",
		"E123B",
		"123*",
		"???",
	}
	for _, test := range failtests {
		_, e := enc.encode(test)
		if e == nil {
			t.Fatalf("encode(%v) must be error", test)
		}
	}
}

func TestCodabarWriter(t *testing.T) {
	writer := NewCodaBarWriter()

	tests := []struct {
		content string
		wants   string
	}{
		{
			"0",
			"00000" + "10110010010" + "1010100110" + "10110010010" + "0000",
		},
		{
			"T+N",
			"00000" + "10110010010" + "10110110110" + "10010010110" + "0000",
		},
		{
			"*-5$e",
			"00000" + "10100100110" + "1010011010" + "1101010010" + "1011001010" + "10100110010" + "0000",
		},
	}
	for _, test := range tests {
		testEncode(t, writer, gozxing.BarcodeFormat_CODABAR, test.content, test.wants)
	}
}
