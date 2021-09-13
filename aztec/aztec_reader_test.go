package aztec

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestAztecReader_Decode(t *testing.T) {
	d := NewAztecReader()
	d.Reset()

	img, _ := gozxing.NewBitMatrix(1, 1)
	bmp := testutil.NewBinaryBitmapFromBitMatrix(img)
	r, e := d.DecodeWithoutHints(bmp)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	// not found
	img, _ = gozxing.ParseStringToBitMatrix(""+
		"                                  \n"+
		"                                  \n"+
		"      ######          ##  ##      \n"+
		"      ######################      \n"+
		"        ##              ##        \n"+
		"        ##  ##########  ##        \n"+
		"        ##  ##      ##  ##        \n"+
		"        ##  ##  ##  ##  ##        \n"+
		"        ##  ##      ##  ##        \n"+
		"        ##  ##########  ##        \n"+
		"        ##              ##        \n"+
		"        ####################      \n"+
		"                    ##            \n"+
		"                                  \n"+
		"                                  \n",
		"##", "  ")
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 3))
	_, e = d.Decode(bmp, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Decode must return NotFoundException: %+v", e)
	}

	// invalid data
	img, _ = gozxing.ParseStringToBitMatrix(""+
		"    ####      ##  ##  ##  ####\n"+
		"            ##            ##  \n"+
		"    ####            ##  ######\n"+
		"  ########################    \n"+
		"      ##              ##    ##\n"+
		"    ####  ##########  ##  ####\n"+
		"  ##  ##  ##      ##  ####    \n"+
		"####  ##  ##  ##  ##  ##      \n"+
		"####  ##  ##      ##  ########\n"+
		"##    ##  ##########  ######  \n"+
		"  ######              ##    ##\n"+
		"      ######################  \n"+
		"  ##          ####        ##  \n"+
		"      ####        ##  ##  ####\n"+
		"##  ##      ######  ####    ##\n",
		"##", "  ")
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 3))
	_, e = d.Decode(bmp, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must return FormatException: %+v", e)
	}

	// mirror image
	img = testutil.MirrorBitMatrix(img)
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 3))
	_, e = d.Decode(bmp, nil)
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("Decode must return FormatException: %+v", e)
	}

	// correct image
	img, _ = gozxing.ParseStringToBitMatrix(""+
		"                                  \n"+
		"      ##    ##  ####        ##    \n"+
		"    ######    ##  ######      ##  \n"+
		"      ####        ##  ##  ##      \n"+
		"  ##########################      \n"+
		"  ####  ##              ##        \n"+
		"      ####  ##########  ##  ##    \n"+
		"    ##  ##  ##      ##  ##        \n"+
		"    ######  ##  ##  ##  ########  \n"+
		"    ######  ##      ##  ##        \n"+
		"    ######  ##########  ####      \n"+
		"      ####              ######    \n"+
		"  ##    ####################  ##  \n"+
		"  ##        ##    ##  ##          \n"+
		"  ####      ######  ##  ##    ##  \n"+
		"  ########    ####  ####  ##  ##  \n"+
		"                                  \n",
		"##", "  ")
	bmp = testutil.NewBinaryBitmapFromBitMatrix(testutil.ExpandBitMatrix(img, 3))

	points := make([]gozxing.ResultPoint, 0)
	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_NEED_RESULT_POINT_CALLBACK] = gozxing.ResultPointCallback(
		func(point gozxing.ResultPoint) {
			points = append(points, point)
		})

	r, e = d.Decode(bmp, hints)
	if e != nil {
		t.Fatalf("Decode error: %+v", e)
	}
	if txt, wants := r.GetText(), "Histórico"; txt != wants {
		t.Fatalf("GetText: %v, wants %v", txt, wants)
	}
	if bs, wants := r.GetRawBytes(), []byte{79, 21, 74, 252, 62, 115, 81, 33, 128}; !reflect.DeepEqual(bs, wants) {
		t.Fatalf("GetRawBytes: %v, wants %v", bs, wants)
	}
	if num, wants := r.GetNumBits(), 65; num != wants {
		t.Fatalf("GetNumBits: %v, wants %v", num, wants)
	}
	if format, wants := r.GetBarcodeFormat(), gozxing.BarcodeFormat_AZTEC; format != wants {
		t.Fatalf("GetBarcodeFormat: %v, wants %v", format, wants)
	}
	pswants := []gozxing.ResultPoint{
		gozxing.NewResultPoint(47.5, 2.5), gozxing.NewResultPoint(47.5, 47.5),
		gozxing.NewResultPoint(2.5, 47.5), gozxing.NewResultPoint(2.5, 2.5),
	}
	if ps := r.GetResultPoints(); !reflect.DeepEqual(ps, pswants) {
		t.Fatalf("GetResultPoint(): %v, wants %v", ps, pswants)
	}
	if !reflect.DeepEqual(points, pswants) {
		t.Fatalf("ResultPointCallback: %v, wants %v", points, pswants)
	}
	metadata := r.GetResultMetadata()
	if ecl, wants := metadata[gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL], "35%"; ecl != wants {
		t.Fatalf("Metadata[ERROR_CORRECTION_LEVEL]: %v, wants %v", ecl, wants)
	}
	if si, wants := metadata[gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER], "]z0"; si != wants {
		t.Fatalf("Metadata[SYMBOLOGY_IDENTIFIER]: %v, wants %v", si, wants)
	}
}

func TestAztecReader_Decode_Blackbox(t *testing.T) {
	reader := NewAztecReader()
	format := gozxing.BarcodeFormat_AZTEC

	tests := []struct {
		file  string
		wants string
	}{
		// testdata from zxing core/src/test/resources/blackbox/aztec-[12]/
		{"testdata/aztec-1/7.png", "Code 2D!"},
		{"testdata/aztec-1/Historico.png", "Histórico"},
		{"testdata/aztec-1/HistoricoLong.png", "Históóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóóórico"},
		{"testdata/aztec-1/abc-19x19C.png", "abcdefghijklmnopqrstuvwxyz"},
		{"testdata/aztec-1/abc-37x37.png", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"testdata/aztec-1/dlusbs.png", "3333h3i3jITIT"},
		{"testdata/aztec-1/hello.png", "hello"},
		{"testdata/aztec-1/lorem-075x075.png", "In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p"},
		{"testdata/aztec-1/lorem-105x105.png", "In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p.In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p"},
		{"testdata/aztec-1/lorem-131x131.png", "In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p.In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p. In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p.In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tris. In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo e."},
		{"testdata/aztec-1/lorem-151x151.png", "In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p.In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p. In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tristique ac, p.In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem. Nulla ultrices, elit pellentesque aliquet laoreet, justo erat pulvinar nisi, id elementum sapien dolor et diam. Donec ac nunc sodales elit placerat eleifend. Sed ornare luctus ornare. Vestibulum vehicula, massa at pharetra fringilla, risus justo faucibus erat, nec porttitor nibh tellus sed est. Ut justo diam, lobortis eu tris. In ut magna vel mauris malesuada dictum. Nulla ullamcorper metus quis diam cursus facilisis. Sed mollis quam id justo rutrum sagittis. Donec laoreet rutrum est, nec convallis mauris condimentum sit amet. Phasellus gravida, justo et congue auctor, nisi ipsum viverra erat, eget hendrerit felis turpis nec lorem."},
		{"testdata/aztec-1/tableShifts.png", "AhUUDgdy672;..:8KjHH776JHHn3g. 8lm/%22Nn873R2897ks4JKDJ9JJaza2323!::;09UJRrhDQSKJDKdSJSdskjdslkEdjseze:ze"},
		{"testdata/aztec-1/tag.png", "Ceci est un tag."},
		{"testdata/aztec-1/texte.png", "Ceci est un texte!"},

		{"testdata/aztec-2/01.png", "This is a real world Aztec barcode test."},
		{"testdata/aztec-2/02.png", "This is a real world Aztec barcode test."},
		//{"testdata/aztec-2/03.png", "This is a real world Aztec barcode test."},
		//{"testdata/aztec-2/04.png", "This is a real world Aztec barcode test."},
		{"testdata/aztec-2/05.png", "This is a real world Aztec barcode test."},
		//{"testdata/aztec-2/06.png", "This is a real world Aztec barcode test."},
		//{"testdata/aztec-2/07.png", "This is a real world Aztec barcode test."},
		//{"testdata/aztec-2/08.png", "This is a real world Aztec barcode test."},
		{"testdata/aztec-2/09.png", "mailto:zxing@googlegroups.com"},
		//{"testdata/aztec-2/10.png", "mailto:zxing@googlegroups.com"},
		//{"testdata/aztec-2/11.png", "mailto:zxing@googlegroups.com"},
		//{"testdata/aztec-2/12.png", "mailto:zxing@googlegroups.com"},
		//{"testdata/aztec-2/13.png", "mailto:zxing@googlegroups.com"},
		//{"testdata/aztec-2/14.png", "mailto:zxing@googlegroups.com"},
		//{"testdata/aztec-2/15.png", "mailto:zxing@googlegroups.com"},
		//{"testdata/aztec-2/16.png", "http://code.google.com/p/zxing/source/browse/trunk/android/src/com/google/zxing/client/android/result/URIResultHandler.java"},
		//{"testdata/aztec-2/17.png", "http://code.google.com/p/zxing/source/browse/trunk/android/src/com/google/zxing/client/android/result/URIResultHandler.java"},
		{"testdata/aztec-2/18.png", "http://code.google.com/p/zxing/source/browse/trunk/android/src/com/google/zxing/client/android/result/URIResultHandler.java"},
		//{"testdata/aztec-2/19.png", "http://code.google.com/p/zxing/source/browse/trunk/android/src/com/google/zxing/client/android/result/URIResultHandler.java"},
		//{"testdata/aztec-2/20.png", "http://code.google.com/p/zxing/source/browse/trunk/android/src/com/google/zxing/client/android/result/URIResultHandler.java"},
		//{"testdata/aztec-2/21.png", "http://code.google.com/p/zxing/source/browse/trunk/android/src/com/google/zxing/client/android/result/URIResultHandler.java"},
		//{"testdata/aztec-2/22.png", "http://code.google.com/p/zxing/source/browse/trunk/android/src/com/google/zxing/client/android/result/URIResultHandler.java"},
	}

	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, nil, nil)
	}
}
