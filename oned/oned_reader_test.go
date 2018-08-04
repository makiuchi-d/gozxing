package oned

import (
	"image"
	_ "image/png"
	"os"
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
