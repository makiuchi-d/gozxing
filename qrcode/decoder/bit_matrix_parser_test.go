package decoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

// formatinfo = 110100101110110 = 0x6976
//   ecl = 0b11 (L)
//   mask = 0b010
var qrstr = "" +
	"##############    ##  ####  ##############\n" +
	"##          ##  ####  ##    ##          ##\n" +
	"##  ######  ##  ####    ##  ##  ######  ##\n" +
	"##  ######  ##    ##  ##    ##  ######  ##\n" +
	"##  ######  ##  ##      ##  ##  ######  ##\n" +
	"##          ##  ##    ####  ##          ##\n" +
	"##############  ##  ##  ##  ##############\n" +
	"                ##########                \n" +
	"####  ##    ####  ####      ######  ####  \n" +
	"  ##########  ######        ##        ####\n" +
	"    ####  ########  ##  ####      ####  ##\n" +
	"      ##  ##    ##    ##          ##  ####\n" +
	"        ##  ####  ####  ##  ##  ##        \n" +
	"                ########      ####  ##  ##\n" +
	"##############  ######    ##  ##  ######  \n" +
	"##          ##    ##########  ####        \n" +
	"##  ######  ##    ##  ##    ######      ##\n" +
	"##  ######  ##  ##  ####      ##  ########\n" +
	"##  ######  ##    ####  ##      ##  ##  ##\n" +
	"##          ##  ######    ####            \n" +
	"##############  ##  ######    ##  ##  ##  "

func TestNewBitMatrixParser(t *testing.T) {
	img, _ := gozxing.NewSquareBitMatrix(20)
	_, e := NewBitMatrixParser(img)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("NewBitMatrixParser(20x20) must be FormatException, %T", e)
	}

	img, _ = gozxing.NewSquareBitMatrix(22)
	_, e = NewBitMatrixParser(img)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("NewBitMatrixParser(22x22) must be FormatException, %T", e)
	}

	img, _ = gozxing.NewSquareBitMatrix(21)
	p, e := NewBitMatrixParser(img)
	if e != nil {
		t.Fatalf("NewBitMatrixParser(21x21) returns error, %v", e)
	}
	if p.bitMatrix != img {
		t.Fatalf("p.bitMatrix = %p, expect %p", p.bitMatrix, img)
	}
	if p.parsedVersion != nil {
		t.Fatalf("p.parsedVersion is not nil, %p", p.parsedVersion)
	}
	if p.parsedFormatInfo != nil {
		t.Fatalf("p.parsedFormatInfo is not nil, %p", p.parsedFormatInfo)
	}
	if p.mirror != false {
		t.Fatalf("p.mirror is not false")
	}
}

func TestBitMatrixParser_setMirror(t *testing.T) {
	img, _ := gozxing.NewSquareBitMatrix(21)
	p, _ := NewBitMatrixParser(img)

	p.parsedVersion = &Version{}
	p.parsedFormatInfo = &FormatInformation{}

	p.SetMirror(true)
	if p.mirror != true {
		t.Fatalf("p.mirror must be true")
	}
	if p.parsedVersion != nil {
		t.Fatalf("p.parsedVersion must be nil")
	}
	if p.parsedFormatInfo != nil {
		t.Fatalf("p.parsedFormatInfo must be nil")
	}

	p.SetMirror(false)
	if p.mirror != false {
		t.Fatalf("p.mirror must be false")
	}
}

func TestBitMatrixParser_copyBit(t *testing.T) {
	img, _ := gozxing.NewSquareBitMatrix(21)
	for i := 0; i < 21; i++ {
		if i%2 == 0 {
			img.Set(i, 10)
		}
		if i%3 == 0 {
			img.Set(i, 11)
		}
	}
	p, _ := NewBitMatrixParser(img)

	bits := 0
	bits = p.copyBit(10, 10, bits)
	bits = p.copyBit(11, 10, bits)
	bits = p.copyBit(12, 10, bits)
	bits = p.copyBit(13, 10, bits)
	if bits != 10 {
		t.Fatalf("bits = %v, expect 10", bits)
	}

	p.SetMirror(true)
	bits = 0
	bits = p.copyBit(10, 6, bits)
	bits = p.copyBit(11, 6, bits)
	bits = p.copyBit(12, 6, bits)
	bits = p.copyBit(13, 6, bits)
	if bits != 12 {
		t.Fatalf("bits = %v, expect 12", bits)
	}
}

func TestBitMatrixParser_ReadFormatInformation(t *testing.T) {
	img, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	p, _ := NewBitMatrixParser(img)

	info, e := p.ReadFormatInformation()
	if e != nil {
		t.Fatalf("ReadFormatInformation returns error, %v", e)
	}
	if r := info.GetErrorCorrectionLevel(); r != ErrorCorrectionLevel_L {
		t.Fatalf("ErrorCollectionLevel = %v, expect L", r)
	}
	if r := info.GetDataMask(); r != 7 {
		t.Fatalf("DataMask = %v, expect 7", r)
	}

	info2, e := p.ReadFormatInformation()
	if e != nil {
		t.Fatalf("ReadFormatInformation returns error, %v", e)
	}
	if info != info2 {
		t.Fatalf("info and info2 must be same, %p, %p", info, info2)
	}

	for i := 0; i < 21; i++ {
		img.Unset(i, 8)
	}
	p, _ = NewBitMatrixParser(img)

	info, e = p.ReadFormatInformation()
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("ReadFormatInformation must return FormatException, %T", e)
	}

}

func testBitMatrixParser_ReadVersion(t testing.TB, img *gozxing.BitMatrix, expect_version int) {
	t.Helper()
	p, _ := NewBitMatrixParser(img)
	ver, e := p.ReadVersion()
	if e != nil {
		t.Fatalf("ReadVersion returns error, %v", e)
	}
	vernum := ver.GetVersionNumber()
	if vernum != expect_version {
		t.Fatalf("VersionNumber = %v, expect %v", vernum, expect_version)
	}
}

func TestBitMatrixParser_ReadVersion(t *testing.T) {

	img, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	p, _ := NewBitMatrixParser(img)
	ver, e := p.ReadVersion()
	if e != nil {
		t.Fatalf("ReadVersion returns error, %v", e)
	}
	if r := ver.GetVersionNumber(); r != 1 {
		t.Fatalf("VersionNumber = %v, expect 1", r)
	}

	img, _ = gozxing.NewSquareBitMatrix(41)
	testBitMatrixParser_ReadVersion(t, img, 6)

	// write version bits on right-top
	// ver17: 85x85 0x1145D = 010 001 010 001 011 101
	dim := 85
	img, _ = gozxing.NewSquareBitMatrix(dim)
	img.Set(dim-10, 5)
	img.Set(dim-11, 4)
	img.Set(dim-10, 3)
	img.Set(dim-11, 2)
	img.Set(dim-10, 1)
	img.Set(dim-11, 1)
	img.Set(dim-9, 0)
	img.Set(dim-11, 0)
	testBitMatrixParser_ReadVersion(t, img, 17)

	// write version bits on left-bottom
	// ver31: 0x1F250 = 011 111 001 001 010 000
	dim = 141
	img, _ = gozxing.NewSquareBitMatrix(dim)
	img.Set(5, dim-10)
	img.Set(5, dim-11)
	img.Set(4, dim-9)
	img.Set(4, dim-10)
	img.Set(4, dim-11)
	img.Set(3, dim-11)
	img.Set(2, dim-11)
	img.Set(1, dim-10)
	testBitMatrixParser_ReadVersion(t, img, 31)

	p, _ = NewBitMatrixParser(img)
	ver1, _ := p.ReadVersion()
	ver2, _ := p.ReadVersion()
	if ver1 != ver2 {
		t.Fatalf("ver1(%p) != ver2(%p)", ver1, ver2)
	}

	img.Clear()
	p, _ = NewBitMatrixParser(img)
	_, e = p.ReadVersion()
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("ReadVersion must be FormatException, %T", e)
	}
}

func TestBitMatrixParser_ReadCodewords(t *testing.T) {
	img, _ := gozxing.NewSquareBitMatrix(85)
	img.SetRegion(8, 0, 1, 85)
	p, _ := NewBitMatrixParser(img)
	_, e := p.ReadCodewords()
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("ReadCodewords() must be FormatException, %T", e)
	}

	img.Clear()
	// set format information on left-top
	img.Set(0, 7)
	img.Set(1, 7)
	img.Set(4, 7)
	img.Set(5, 7)
	img.Set(7, 5)
	img.Set(7, 3)
	img.Set(7, 2)
	img.Set(7, 1)
	img.Set(7, 0)
	p, _ = NewBitMatrixParser(img)
	_, e = p.ReadCodewords()
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("ReadCodewords() must be FormatException, %T", e)
	}

	img, _ = gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	p, _ = NewBitMatrixParser(img)
	words, e := p.ReadCodewords()
	if e != nil {
		t.Fatalf("ReadCodewords returns error, %v", e)
	}
	expect := []byte{
		0x40, 0x56, 0x86, 0x56, 0xc6, 0xc6, 0xf0, 0xec, 0x11, 0xec,
		0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0xec, 0x11, 0x25,
		0x19, 0xd0, 0xd2, 0x68, 0x59, 0x39,
	}
	if !reflect.DeepEqual(words, expect) {
		t.Fatalf("codewords different:\n %v\nexpect:\n %v", words, expect)
	}
}

func compareBitMatrix(t testing.TB, img, expect *gozxing.BitMatrix) {
	t.Helper()
	if img.GetWidth() != expect.GetWidth() || img.GetHeight() != expect.GetHeight() {
		t.Fatalf("BitMatrix size different, (%v, %v), expect (%v, %v)",
			img.GetWidth(), img.GetHeight(), expect.GetWidth(), expect.GetHeight())
	}
	for y := 0; y < img.GetHeight(); y++ {
		for x := 0; x < img.GetWidth(); x++ {
			if img.Get(x, y) != expect.Get(x, y) {
				t.Fatalf("BitMatrix different on %v, %v", x, y)
			}
		}
	}
}

func TestBitMatrixParser_Remask(t *testing.T) {
	masked, _ := gozxing.ParseStringToBitMatrix(""+
		"  ##  ##  ##    ########      ##  ##  ##  \n"+
		"##    ########  ##  ##      ########    ##\n"+
		"    ####  ##    ########    ##          ##\n"+
		"######  ########          ########  ######\n"+
		"  ##  ####    ##          ##    ####  ##  \n"+
		"########    ####  ##  ######  ##      ##  \n"+
		"  ##  ##  ##                  ##  ##  ##  \n"+
		"      ######    ##      ##    ######      \n"+
		"  ##  ######  ##  ##  ####  ####  ##  ##  \n"+
		"    ##  ##      ##    ##  ######  ##    ##\n"+
		"####  ##  ##        ##      ##    ##  ##  \n"+
		"  ####    ##  ##  ##  ##  ######  ##      \n"+
		"##  ##        ######                ##  ##\n"+
		"      ######    ##                ####  ##\n"+
		"  ######        ####  ######  ####    ##  \n"+
		"####  ##  ######    ##  ##      ####  ##  \n"+
		"  ##  ####    ######  ######  ####  ####  \n"+
		"####    ##  ####  ######  ####    ####    \n"+
		"      ##        ####        ##            \n"+
		"##    ########  ##    ##  ##########      \n"+
		"  ######        ##            ####  ####  ", "##", "  ")
	unmasked, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")

	img, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	p, _ := NewBitMatrixParser(img)

	p.Remask()
	compareBitMatrix(t, img, unmasked)

	p.ReadFormatInformation()
	p.Remask()
	compareBitMatrix(t, img, masked)

	p.Remask()
	compareBitMatrix(t, img, unmasked)
}

func TestBitMatrixParser_Mirror(t *testing.T) {
	mirrored, _ := gozxing.ParseStringToBitMatrix(""+
		"##############  ##          ##############\n"+
		"##          ##  ####        ##          ##\n"+
		"##  ######  ##    ####      ##  ######  ##\n"+
		"##  ######  ##  ########    ##  ######  ##\n"+
		"##  ######  ##    ##    ##  ##  ######  ##\n"+
		"##          ##    ######    ##          ##\n"+
		"##############  ##  ##  ##  ##############\n"+
		"                ######  ##                \n"+
		"  ####  ########  ######  ####    ##  ####\n"+
		"########      ######    ##########  ####  \n"+
		"            ######  ##  ########  ########\n"+
		"####  ##  ##  ##      ##  ##  ######    ##\n"+
		"##  ##  ########    ##  ##    ##    ##  ##\n"+
		"                    ##      ####      ##  \n"+
		"##############  ####    ##      ##    ##  \n"+
		"##          ##  ##        ##########    ##\n"+
		"##  ######  ##  ##      ####  ####  ##    \n"+
		"##  ######  ##      ####    ##    ##    ##\n"+
		"##  ######  ##  ##  ##    ####    ####    \n"+
		"##          ##  ####  ##    ##    ##    ##\n"+
		"##############    ######  ##    ######    ", "##", "  ")
	unmirrored, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")
	img, _ := gozxing.ParseStringToBitMatrix(qrstr, "##", "  ")

	p, _ := NewBitMatrixParser(img)

	p.Mirror()
	compareBitMatrix(t, img, mirrored)

	p.Mirror()
	compareBitMatrix(t, img, unmirrored)
}
