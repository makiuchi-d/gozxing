package gozxing

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/stretchr/testify/assert"
)

// scanningQrCode helper
func scanningQrCodeFile(path string) (output string, err error) {

	// open and decode image file
	file, err := os.Open(path)

	if err != nil {
		return
	}

	img, _, err := image.Decode(file)

	if err != nil {
		return
	}

	// prepare BinaryBitmap
	bmp, err := NewBinaryBitmapFromImage(img)

	if err != nil {
		return
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)

	output = result.String()
	fmt.Println("The Result: ", result)

	return
}

// ScanningQrCodeFile helper
func TestScanningQrCodeFile(t *testing.T) {

	testFiles := []string{
		"testdata/qrcode.png",
		"testdata/qrcode-A4.png",
		"testdata/qrcode-border-bottom-center-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-border-bottom-left-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-border-bottom-right-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-border-top-center-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-border-top-left-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-border-top-right-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-bottom-center-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-bottom-left-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-bottom-right-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-center-left-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-center-right-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-top-center-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-top-left-a4-size-paper-white-background-images.jpeg",
		"testdata/qrcode-top-right-a4-size-paper-white-background-images.jpeg",
		// "testdata/qrcode-center-center-a4-size-paper-white-background-images.jpeg",
		// "testdata/a4-size-paper-white-background-images.jpeg",

		"testdata/qrcode-border-bottom-center-a4-size-paper-white-background-images.png",
		"testdata/qrcode-border-bottom-left-a4-size-paper-white-background-images.png",
		"testdata/qrcode-border-bottom-right-a4-size-paper-white-background-images.png",
		"testdata/qrcode-border-top-center-a4-size-paper-white-background-images.png",
		"testdata/qrcode-border-top-left-a4-size-paper-white-background-images.png",
		"testdata/qrcode-border-top-right-a4-size-paper-white-background-images.png",
		"testdata/qrcode-bottom-center-a4-size-paper-white-background-images.png",
		"testdata/qrcode-bottom-left-a4-size-paper-white-background-images.png",
		"testdata/qrcode-bottom-right-a4-size-paper-white-background-images.png",
		"testdata/qrcode-center-left-a4-size-paper-white-background-images.png",
		"testdata/qrcode-center-right-a4-size-paper-white-background-images.png",
		"testdata/qrcode-top-center-a4-size-paper-white-background-images.png",
		"testdata/qrcode-top-left-a4-size-paper-white-background-images.png",
		"testdata/qrcode-top-right-a4-size-paper-white-background-images.png",
		// "testdata/qrcode-center-center-a4-size-paper-white-background-images.png",
		// "testdata/a4-size-paper-white-background-images.png",
	}

	t.Parallel()

	for _, file := range testFiles {

		result, err := scanningQrCodeFile(file)

		assert.Nil(t, err, "Should be nil")
		assert.NotEmpty(t, result, "Should not be Empty")
	}
}
