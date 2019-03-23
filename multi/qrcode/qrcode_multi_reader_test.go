package qrcode

import (
	"reflect"
	"sort"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

var qrstr = "" +
	"##############      ##  ##  ##############        ##############      ##  ##  ##############\n" +
	"##          ##          ##  ##          ##        ##          ##          ##  ##          ##\n" +
	"##  ######  ##  ##  ##      ##  ######  ##        ##  ######  ##  ##  ##      ##  ######  ##\n" +
	"##  ######  ##          ##  ##  ######  ##        ##  ######  ##          ##  ##  ######  ##\n" +
	"##  ######  ##    ##  ####  ##  ######  ##        ##  ######  ##    ##  ####  ##  ######  ##\n" +
	"##          ##    ######    ##          ##        ##          ##    ######    ##          ##\n" +
	"##############  ##  ##  ##  ##############        ##############  ##  ##  ##  ##############\n" +
	"                ##  ##                                            ##  ##                    \n" +
	"######  ##########  ##  ######      ##            ######  ##########  ##  ######      ##    \n" +
	"  ##  ##        ########  ##  ##      ####          ##  ##        ########  ##  ##      ####\n" +
	"##    ####  ##  ########  ######  ########        ##    ####  ##  ########  ######  ########\n" +
	"    ####  ##  ####    ######  ####    ##              ####  ##  ####    ######  ####    ##  \n" +
	"        ##  ##    ##  ##  ######                          ##  ##    ##  ##  ######          \n" +
	"                ##  ##      ####    ######                        ##  ##      ####    ######\n" +
	"##############  ##  ##  ##      ##  ######        ##############  ##  ##  ##      ##  ######\n" +
	"##          ##  ######      ######    ####        ##          ##  ######      ######    ####\n" +
	"##  ######  ##  ####    ##  ##        ####        ##  ######  ##  ####    ##  ##        ####\n" +
	"##  ######  ##    ######  ##  ##    ####          ##  ######  ##    ######  ##  ##    ####  \n" +
	"##  ######  ##  ########  ####  ##  ##  ##        ##  ######  ##  ########  ####  ##  ##  ##\n" +
	"##          ##  ##  ########    ##    ##          ##          ##  ##  ########    ##    ##  \n" +
	"##############  ########  ######      ####        ##############  ########  ######      ####\n"

func TestNewSAComparator(t *testing.T) {
	r1 := gozxing.NewResult("r1", []byte{}, []gozxing.ResultPoint{}, gozxing.BarcodeFormat_QR_CODE)
	r2 := gozxing.NewResult("r2", []byte{}, []gozxing.ResultPoint{}, gozxing.BarcodeFormat_QR_CODE)
	r3 := gozxing.NewResult("r3", []byte{}, []gozxing.ResultPoint{}, gozxing.BarcodeFormat_QR_CODE)
	r1.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE, (0<<4)+2)
	r2.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE, (1<<4)+2)
	r3.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE, (2<<4)+2)

	results := []*gozxing.Result{r3, r1, r2}
	sort.Slice(results, newSAComparator(results))

	wants := []*gozxing.Result{r1, r2, r3}
	if !reflect.DeepEqual(results, wants) {
		t.Fatalf("sorted results %v, wants %v", results, wants)
	}
}

func resultsContains(results []*gozxing.Result, str string) bool {
	for _, r := range results {
		if r.GetText() == str {
			return true
		}
	}
	return false
}

func TestProcessStructuredAppend(t *testing.T) {
	sa1 := gozxing.NewResult("SA1", []byte{}, []gozxing.ResultPoint{}, gozxing.BarcodeFormat_QR_CODE)
	sa2 := gozxing.NewResult("SA2", []byte{}, []gozxing.ResultPoint{}, gozxing.BarcodeFormat_QR_CODE)
	sa3 := gozxing.NewResult("SA3", []byte{}, []gozxing.ResultPoint{}, gozxing.BarcodeFormat_QR_CODE)
	nsa := gozxing.NewResult("NotSA", []byte{}, []gozxing.ResultPoint{}, gozxing.BarcodeFormat_QR_CODE)
	sa1.PutMetadata(gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL, "L")
	sa2.PutMetadata(gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL, "L")
	sa3.PutMetadata(gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL, "L")
	nsa.PutMetadata(gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL, "L")

	expectedByteSegments := [][]byte{[]byte("ByteSegment")}
	sa2.PutMetadata(gozxing.ResultMetadataType_BYTE_SEGMENTS, expectedByteSegments)
	nsa.PutMetadata(gozxing.ResultMetadataType_BYTE_SEGMENTS, expectedByteSegments)

	// no structured append
	results := processStructuredAppend([]*gozxing.Result{sa3, sa1, nsa, sa2})
	if len(results) != 4 {
		t.Fatalf("processed results count=%v, wants %v", len(results), 4)
	}
	for _, str := range []string{"SA1", "SA2", "SA3", "NotSA"} {
		if !resultsContains(results, str) {
			t.Fatalf("results dose not contain \"%s\"", str)
		}
	}

	// with structured append
	sa1.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE, (0<<4)+2)
	sa2.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE, (1<<4)+2)
	sa3.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE, (2<<4)+2)
	results = processStructuredAppend([]*gozxing.Result{sa3, sa1, nsa, sa2})
	if len(results) != 2 {
		t.Fatalf("processed results count=%v, wants %v", len(results), 2)
	}
	for _, str := range []string{"SA1SA2SA3", "NotSA"} {
		if !resultsContains(results, str) {
			t.Fatalf("results dose not contain \"%s\"", str)
		}
	}

	for _, r := range results {
		byteSetments, ok := r.GetResultMetadata()[gozxing.ResultMetadataType_BYTE_SEGMENTS].([][]byte)
		if !ok {
			t.Fatalf("result(%v) must have BYTE_SEGMENTS metadata", r)
		}
		if !reflect.DeepEqual(byteSetments, expectedByteSegments) {
			t.Fatalf("result(%v) byteSegments = %v, want %v", r, byteSetments, expectedByteSegments)
		}
	}
}

func TestQRCodeMultiReader_DecodeMultiple(t *testing.T) {
	reader := NewQRCodeMultiReader()

	bmp, _ := gozxing.NewBitMatrix(1, 1)
	img := testutil.NewBinaryBitmapFromBitMatrix(bmp)
	_, e := reader.DecodeMultiple(img, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeMultiple must be NotFoundException, %T", e)
	}

	bmp, _ = gozxing.NewBitMatrix(10, 10)
	img = testutil.NewBinaryBitmapFromBitMatrix(bmp)
	_, e = reader.DecodeMultiple(img, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("DecodeMultiple must be NotFoundException, %T", e)
	}

	bmp, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	img = testutil.NewBinaryBitmapFromBitMatrix(bmp)
	results, e := reader.DecodeMultiple(img, nil)
	if e != nil {
		t.Fatalf("DecodeMultiple returns error: %v", e)
	}
	if n := len(results); n != 2 {
		t.Fatalf("DecodeMultiple len(results) = %v, wants 2", n)
	}
	for i, r := range results {
		if r.GetText() != "hello\n" {
			t.Fatalf("DecodeMultiple results[%v] = \"%v\", wants \"hello\\n\"", i, r.GetText())
		}
	}

	bmp, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	bmp = testutil.MirrorBitMatrix(bmp)
	img = testutil.NewBinaryBitmapFromBitMatrix(bmp)
	results, e = reader.DecodeMultiple(img, nil)
	if e != nil {
		t.Fatalf("DecodeMultiple returns error: %v", e)
	}
	if n := len(results); n != 2 {
		t.Fatalf("DecodeMultiple len(results) = %v, wants 2", n)
	}
	for i, r := range results {
		if r.GetText() != "hello\n" {
			t.Fatalf("DecodeMultiple results[%v] = \"%v\", wants \"hello\\n\"", i, r.GetText())
		}
	}

	bmp, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	bmp.SetRegion(8, 8, 10, 10)
	img = testutil.NewBinaryBitmapFromBitMatrix(bmp)
	results, e = reader.DecodeMultiple(img, nil)
	if e != nil {
		t.Fatalf("DecodeMultiple returns error: %v", e)
	}
	if n := len(results); n != 1 {
		t.Fatalf("DecodeMultiple len(results) = %v, wants 1", n)
	}
	if txt := results[0].GetText(); txt != "hello\n" {
		t.Fatalf("DecodeMultiple results[0] = \"%v\", wants \"hello\\n\"", txt)
	}
}

func TestQRCodeMultiReader_DecodeMultipleWithoutHint(t *testing.T) {
	reader := NewQRCodeMultiReader()

	testResults := []struct {
		file     string
		contents []string
	}{
		// https://github.com/zxing/zxing/tree/master/core/src/test/resources/blackbox/multi-qrcode-1
		{"testdata/1.png", []string{
			"You get to CREATE OUR JOURNAL PROMPT FOR THE DAY!  Yay!  Way to go!  ",
			"You earned the class 5 EXTRA MINUTES OF RECESS!!  Fabulous!!  Way to go!!",
			"You earned the class a 5 MINUTE DANCE PARTY!!  Awesome!  Way to go!  Let's boogie!",
			"You get to SIT AT MRS. SIGMON'S DESK FOR A DAY!!  Awesome!!  Way to go!! Guess I better clean up! :)",
		}},
		// ISO/IEC 18004:2000 Figure 22
		{"testdata/sa.png", []string{
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		}},
	}
	for _, test := range testResults {
		bmp := testutil.NewBinaryBitmapFromFile(test.file)
		results, e := reader.DecodeMultipleWithoutHint(bmp)
		if e != nil {
			t.Fatalf("DecodeMultiple returns error: %v", e)
		}
		if nr, ne := len(results), len(test.contents); nr != ne {
			t.Fatalf("len(results) = %v, wants %v", nr, ne)
		}
		for i, r := range results {
			if txt, expect := r.GetText(), test.contents[i]; txt != expect {
				t.Fatalf("results[%v] = \n\"%v\", wants \n\"%v\"", i, txt, expect)
			}
		}
	}
}
