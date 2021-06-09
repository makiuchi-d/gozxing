package decoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
)

func TestDecodedBitStreamParsser_parseECIValue(t *testing.T) {
	bits := common.NewBitSource([]byte{})
	_, e := DecodedBitStreamParser_parseECIValue(bits)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("parseECIValue must be FormatException, %T", e)
	}

	// just one byte
	bits = common.NewBitSource([]byte{0x01})
	r, e := DecodedBitStreamParser_parseECIValue(bits)
	if e != nil {
		t.Fatalf("parseECIValue returns error, %v", e)
	}
	if r != 0x01 {
		t.Fatalf("parseECIValue = 0x%02x, expect 0x01", r)
	}

	// two bytes
	bits = common.NewBitSource([]byte{0xa2})
	_, e = DecodedBitStreamParser_parseECIValue(bits)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("parseECIValue must be FormatException, %T", e)
	}
	bits = common.NewBitSource([]byte{0xa2, 0x22})
	r, e = DecodedBitStreamParser_parseECIValue(bits)
	if e != nil {
		t.Fatalf("parseECIValue returns error, %v", e)
	}
	if r != 0x2222 {
		t.Fatalf("parseECIValue = 0x%04x, expect 0x2222", r)
	}

	// three bytes
	bits = common.NewBitSource([]byte{0xd2, 0x34})
	_, e = DecodedBitStreamParser_parseECIValue(bits)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("parseECIValue must be FormatException, %T", e)
	}
	bits = common.NewBitSource([]byte{0xd2, 0x34, 0x56})
	r, e = DecodedBitStreamParser_parseECIValue(bits)
	if e != nil {
		t.Fatalf("parseECIValue returns error, %v", e)
	}
	if r != 0x123456 {
		t.Fatalf("parseECIValue = 0x%06x, expect 0x123456", r)
	}

	bits = common.NewBitSource([]byte{0xff})
	_, e = DecodedBitStreamParser_parseECIValue(bits)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("parseECIValue must be FormatException, %T", e)
	}
}

func TestDecodedBitStreamParser_decodeNumericSegment(t *testing.T) {
	var e error
	var bits *common.BitSource
	result := make([]byte, 0, 10)

	// count = 3n
	// less than 10bits
	bits = common.NewBitSource([]byte{0x1e})
	_, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 3)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeNumericSegment must be FormatException, %T", e)
	}
	// 1000 = 0x3e8 => 1111 1010 00
	bits = common.NewBitSource([]byte{0xfa, 0x00})
	_, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 3)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeNumericSegment must be FormatException, %T", e)
	}
	// 123456 = 0x07b*1000 + 0x1c8 => 00011110110111001000
	bits = common.NewBitSource([]byte{0x1e, 0xdc, 0x80})
	result, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 6)
	if e != nil {
		t.Fatalf("decodeNumericSegment returns error, %v", e)
	}
	if string(result) != "123456" {
		t.Fatalf("decodeNumericSegment = \"%s\", expect \"123456\"", string(result))
	}

	// count = 3n+2
	bits = common.NewBitSource([]byte{})
	_, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 2)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeNumericSegment must be FormatException, %T", e)
	}
	// 100 = 0x64 => 1100100
	bits = common.NewBitSource([]byte{0xc8})
	result, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 2)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeNumericSegment must be FormatException, %T", e)
	}
	// 99 = 0x63 => 110 0011 => 1100 0110
	bits = common.NewBitSource([]byte{0xc6})
	result, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 2)
	if e != nil {
		t.Fatalf("decodeNumericSegment returns error, %v", e)
	}
	if string(result) != "99" {
		t.Fatalf("decodeNumericSegment = \"%s\", expect \"99\"", string(result))
	}

	// count = 3n+1
	bits = common.NewBitSource([]byte{})
	_, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 1)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeNumericSegment must be FormatException, %T", e)
	}
	// 10 = 0xa
	bits = common.NewBitSource([]byte{0xa0})
	_, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 1)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeNumericSegment must be FormatException, %T", e)
	}
	// 9 = 0x9
	bits = common.NewBitSource([]byte{0x90})
	result, e = DecodedBitStreamParser_decodeNumericSegment(bits, result[:0], 1)
	if e != nil {
		t.Fatalf("decodeNumericSegment returns error, %v", e)
	}
	if string(result) != "9" {
		t.Fatalf("decodeNumericSegment = \"%s\", expect \"9\"", string(result))
	}
}

func TestDecodedBitStreamParser_decodeAlphanumericSegment(t *testing.T) {
	var e error
	var bits *common.BitSource
	result := make([]byte, 0, 10)

	bits = common.NewBitSource([]byte{0x00})
	_, e = DecodedBitStreamParser_decodeAlphanumericSegment(bits, result[:0], 3, false)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeAlphanumericSegment must be FormatException, %T", e)
	}
	// 45*45 = 0x7e9 => 1111 1101 001
	bits = common.NewBitSource([]byte{0xfd, 0x20})
	_, e = DecodedBitStreamParser_decodeAlphanumericSegment(bits, result[:0], 2, false)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeAlphanumericSegment must be FormatException, %T", e)
	}

	// count==1
	bits = common.NewBitSource([]byte{})
	_, e = DecodedBitStreamParser_decodeAlphanumericSegment(bits, result[:0], 1, false)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeAlphanumericSegment must be FormatException, %T", e)
	}
	// 45 = 0x2d => 10 1101
	bits = common.NewBitSource([]byte{0xb4})
	_, e = DecodedBitStreamParser_decodeAlphanumericSegment(bits, result[:0], 1, false)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeAlphanumericSegment must be FormatException, %T", e)
	}

	// ABC => AB(45*10 + 11 = 0x1cd), C(0xc) => 00111001101, 001100
	bits = common.NewBitSource([]byte{0x39, 0xa6, 0x00})
	result, e = DecodedBitStreamParser_decodeAlphanumericSegment(bits, result[:0], 3, false)
	if e != nil {
		t.Fatalf("decodeAlphanumericSegment returns error, %v", e)
	}
	if string(result) != "ABC" {
		t.Fatalf("decodeAlphanumericSegment = \"%s\", expect \"ABC\"", string(result))
	}

	// "A%%B%C" => 10, 38, 38, 11, 38, 12 => 45*10+38, 45*38+11, 45*38+12
	// => 1e8, 6b9, 6ba => 0011 1101 000,1 1010 1110 01,11 0101 1101 0
	bits = common.NewBitSource([]byte{0x3d, 0x1a, 0xe7, 0x5d, 0x00})
	expect := []byte{'A', '%', 'B', 0x1d, 'C'}
	result, e = DecodedBitStreamParser_decodeAlphanumericSegment(bits, result[:0], 6, true)
	if e != nil {
		t.Fatalf("decodeAlphanumericSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("decodeAlphanumericSegment = %v, expect %v", result, expect)
	}
}

func TestDecodedBitStreamParser_decodeByteSegment(t *testing.T) {
	var e error
	var bits *common.BitSource
	result := make([]byte, 0, 10)
	byteSegments := make([][]byte, 0, 1)

	bits = common.NewBitSource([]byte{})
	_, _, e = DecodedBitStreamParser_decodeByteSegment(bits, result, 1, nil, byteSegments, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeByteSegment must be FormatException, %T", e)
	}

	// invalid charset
	bits = common.NewBitSource([]byte{0x8b, 0xe0, 0x8b, 0x9b})
	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_CHARACTER_SET] = "dummy"
	_, _, e = DecodedBitStreamParser_decodeByteSegment(bits, result, 4, nil, byteSegments, hints)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeByteSegment must be FormatException, %T", e)
	}

	// Shift_JIS 金魚
	bytes := []byte{0x8b, 0xe0, 0x8b, 0x9b}
	bits = common.NewBitSource(bytes)
	expect := []byte("金魚")
	r, bs, e := DecodedBitStreamParser_decodeByteSegment(bits, result, 4, common.CharacterSetECI_SJIS, byteSegments, nil)
	if e != nil {
		t.Fatalf("decodeByteSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeByteSegment result = %v, expect %v", r, expect)
	}
	if !reflect.DeepEqual(bs, [][]byte{bytes}) {
		t.Fatalf("decodeByteSegment byteSegments = %v, expect %v", bs, [][]byte{bytes})
	}

	// UTF-8 金魚
	bits = common.NewBitSource(expect)
	r, bs, e = DecodedBitStreamParser_decodeByteSegment(bits, result, len(expect), nil, byteSegments, nil)
	if e != nil {
		t.Fatalf("decodeByteSegment returns error, %v", e)
	}
	if !reflect.DeepEqual(r, expect) {
		t.Fatalf("decodeByteSegment result = %v, expect %v", r, expect)
	}
	if !reflect.DeepEqual(bs, [][]byte{expect}) {
		t.Fatalf("decodeByteSegment byteSegments = %v, expect %v", bs, [][]byte{expect})
	}
}

func TestDecodedBitStreamParser_decodeKanjiSegment(t *testing.T) {
	var e error
	var bits *common.BitSource
	result := make([]byte, 0, 10)

	bits = common.NewBitSource([]byte{0xf0})
	_, e = DecodedBitStreamParser_decodeKanjiSegment(bits, result, 1)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeKanjiSegment must be FormatException, %T", e)
	}

	// Shift_JIS "餃子" = 0xe9 0x4c, 0x8e 0x71
	// e94c - c140 = 280c => 28*c0 + 0x0c = 1e0c = 1 1110 0000 1100
	// 8e71 - 8140 =  d31 => 0d*c0 + 0x31 =  9f1 = 0 1001 1111 0001
	// 1111 0000 0110 0010 0111 1100 01
	bits = common.NewBitSource([]byte{0xf0, 0x62, 0x7c, 0x60})
	result, e = DecodedBitStreamParser_decodeKanjiSegment(bits, result[:0], 2)
	if e != nil {
		t.Fatalf("decodeKanjiSegment returns error, %v", e)
	}
	if string(result) != "餃子" {
		t.Fatalf("decodeKanjiSegment = \"%s\", expect \"餃子\"", string(result))
	}
}

func TestDecodedBitStreamParser_decodeHanziSegment(t *testing.T) {
	var e error
	var bits *common.BitSource
	result := make([]byte, 0, 10)

	bits = common.NewBitSource([]byte{0xf0})
	_, e = DecodedBitStreamParser_decodeHanziSegment(bits, result, 1)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("decodeHanjiSegment must be FormatException, %T", e)
	}

	// ♂ a1e1 => a1e1-a1a1=0040, 00*60 + 40 = 0040, 0 0000 0100 0000
	// 安 b0b2 => b0b2-a6a1=0a11, 0a*60 + 11 = 03d1, 0 0011 1101 0001
	// 0000 0010 0000 0000 1111 0100 01
	bits = common.NewBitSource([]byte{0x02, 0x00, 0xf4, 0x40})
	result, e = DecodedBitStreamParser_decodeHanziSegment(bits, result[:0], 2)
	if e != nil {
		t.Fatalf("decodeHanziSegment returns error, %v", e)
	}
	if string(result) != "♂安" {
		t.Fatalf("decodeHanziSegment = \"%s\", expect \"♂安\"", string(result))
	}

	// ア a5a2 => a5a2-a1a1=0401, 04*60 + 01 = 0181, 0 0001 1000 0001
	bits = common.NewBitSource([]byte{0x0c, 0x08})
	result, e = DecodedBitStreamParser_decodeHanziSegment(bits, result[:0], 1)
	if e != nil {
		t.Fatalf("decodeHanziSegment returns error, %v", e)
	}
	if string(result) != "ア" {
		t.Fatalf("decodeHanziSegment = \"%s\", expect \"ア\"", string(result))
	}
}

func TestDecodedBitStreamParsser_DecodeNormalModes(t *testing.T) {
	var bytes []byte
	var e error
	var result *common.DecoderResult
	var ver, _ = Version_GetVersionForNumber(1)

	// invalid mode bits
	bytes = []byte{0x60}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	// numeric      1234  0001, 0000000100, 0001111011 0100
	// alphanumeric ABC   0010, 000000011, 00111001101 001100
	// kanji        餃子  1000, 00000010, 1111000001100 0 100111110001
	// byte(sjis)   邂遘  0100, 00000100, 0xe7 0xae 0xe7 0xa7
	// terminator         0000
	// 0001 0000 0001 0000 0111 1011 0100 0010  0000 0001 1001 1100 1101 0011 0010 0000
	// 0000 1011 1100 0001 1000 1001 1111 0001  0100 0000 0100 1110 0111 1010 1110 1110
	// 0111 1010 0111 0000
	bytes = []byte{
		0x10, 0x10, 0x7b, 0x42, 0x01, 0x9c, 0xd3, 0x20,
		0x0b, 0xc1, 0x89, 0xf1, 0x40, 0x4e, 0x7a, 0xee,
		0x7a, 0x70,
	}
	result, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if s := result.GetText(); s != "1234ABC餃子邂遘" {
		t.Fatalf("Decode result = \"%v\", expect \"1234ABC餃子邂遘\"", s)
	}

	// numeric 1234   0001, 0000000100, 0001111011 0100
	// character count bits error
	bytes = []byte{0x10}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}
	// numeriv value error
	bytes = []byte{0x10, 0x13, 0xff, 0xff}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	// alphanumeric error
	// 0010, 0000 0001 1 00111001101 001100
	bytes = []byte{0x20, 0x1f}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	// kanji error
	bytes = []byte{0x80, 0x6f, 0xff, 0xff}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	// byte segment error
	bytes = []byte{0x40, 0x20}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}
}

func TestDecodedBitStreamParser_DecodeWithCharcterSet(t *testing.T) {
	var bytes []byte
	var e error
	var result *common.DecoderResult
	var ver, _ = Version_GetVersionForNumber(1)

	// ISO-8859-1  0111, 00000001
	// byte Weiß   0100, 00000100, 0x57 0x65 0x69 0xdf
	// Shift_JIS   0111, 00010100,
	// byte 金魚   0100, 00000100, 0x8b 0xe0 0x8b 0x9b
	// without terminator
	bytes = []byte{
		0x70, 0x14, 0x04, 0x57, 0x65, 0x69, 0xdf,
		0x71, 0x44, 0x04, 0x8b, 0xe0, 0x8b, 0x9b,
	}
	result, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if s := result.GetText(); s != "Weiß金魚" {
		t.Fatalf("Decode result = \"%v\", expect \"Weiß金魚\"", s)
	}

	// invalid eci segment  0111, 1111 1111
	bytes = []byte{0x7f, 0xff}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	// invalid eci number (1000) 0111, 1000 0100 0000 0000
	bytes = []byte{0x78, 0x40, 0x00}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}
}

func TestDecodedBitStreamPars_decodeHanziSegment(t *testing.T) {
	var bytes []byte
	var e error
	var result *common.DecoderResult
	var ver, _ = Version_GetVersionForNumber(1)

	// hanzi(gb2312)  1101 0001
	// count          00000011
	// ♂安ア         0000001000000, 0001111010001, 0000110000001
	// 0000 0010 0000 0000 1111 0100 0100 0011 0000 0010
	bytes = []byte{0xd1, 0x03, 0x02, 0x00, 0xf4, 0x43, 0x02}
	result, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if s := result.GetText(); s != "♂安ア" {
		t.Fatalf("Decode result = \"%v\", expect \"♂安ア\"", s)
	}

	// read subset error
	// numeric"1"  0001, 0000 0000 01, 00 01
	// hanzi       11 01, 00(error here)
	bytes = []byte{0x10, 0x04, 0x74}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	// read count error
	// numeric"1"  0001, 0000 0000 01, 00 01
	// hanzi       11 01 00 01, 00 0000 (error here)
	bytes = []byte{0x10, 0x04, 0x74, 0x40}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}

	// read hanzi error
	// numeric"1"  0001, 0000 0000 01, 00 01
	// hanzi       11 01 00 01, 00 0000 11 00 0000 (error here)
	bytes = []byte{0x10, 0x04, 0x74, 0x40, 0xc0}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}
}

func TestDecodedBitStreamPars_decodeStructuredAppend(t *testing.T) {
	var bytes []byte
	var e error
	var result *common.DecoderResult
	var ver, _ = Version_GetVersionForNumber(1)

	// mode  0011
	// symbol sequence 1 of 3 => 0001 0011
	// parity 0xAB
	bytes = []byte{0x31, 0x3a, 0xb0}
	result, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if parity := result.GetStructuredAppendParity(); parity != 0xab {
		t.Fatalf("structured append parity = 0x%02x, expect 0xab", parity)
	}
	if seqNum := result.GetStructuredAppendSequenceNumber(); seqNum != 0x13 {
		t.Fatalf("structured append parity = 0x%02x, expect 0x13", seqNum)
	}

	// less parity
	bytes = []byte{0x31, 0x3a}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}
	// less sequence number
	bytes = []byte{0x31}
	_, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("Decode must be FormatException, %T", e)
	}
}

func TestDecodedBitStreamPars_decodeGS1(t *testing.T) {
	var bytes []byte
	var e error
	var result *common.DecoderResult
	var ver, _ = Version_GetVersionForNumber(1)

	// FNC1(1st position) 01049123451234591597033130128%10ABC123
	// fnc1    1001
	// numeric 0001
	// count 29 = 0000011101
	// 010     0000001010
	// 491     0111101011
	// 234     0011101010
	// 512     1000000000
	// 345     0101011001
	// 915     1110010011
	// 970     1111001010
	// 331     0101001011
	// 301     0100101101
	// 28      0011100
	// alphanumeric 0010
	// count 9 = 000001001
	// %1      11010101111
	// 0A      00000001010
	// BC      00111111011
	// 12      00000101111
	// 3       000011
	// 1001 0001 0000 0111 0100 0000 1010 0111  1010 1100 1110 1010 1000 0000 0001 0101
	// 1001 1110 0100 1111 1100 1010 0101 0010  1101 0010 1101 0011 1000 0100 0000 1001
	// 1101 0101 1110 0000 0010 1000 1111 1101  1000 0010 1111 0000 11
	bytes = []byte{
		0x91, 0x07, 0x40, 0xa7, 0xac, 0xea, 0x80, 0x15,
		0x9e, 0x4f, 0xca, 0x52, 0xd2, 0xd3, 0x84, 0x09,
		0xd5, 0xe0, 0x28, 0xfd, 0x82, 0xf0, 0xc0,
	}
	result, e = DecodedBitStreamParser_Decode(bytes, ver, ErrorCorrectionLevel_M, nil)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	expect := "01049123451234591597033130128\x1d10ABC123" // % replaced to 0x1D
	if text := result.GetText(); text != expect {
		t.Fatalf("Decode result: \"%v\", expect \"%v\" ", text, expect)
	}
}

func TestDecodedBitStreamParser_Decode_SymbologyModifier(t *testing.T) {
	tests := []struct {
		symbologyModifier int
		bytes             []byte
	}{
		{4, []byte{0x57, 0x1a, 0x00}}, // FNC1-1stPos, ECI(26), TERM
		{6, []byte{0x71, 0xa9, 0x00}}, // ECI(26), FNC1-2ndPos, TERM
		{2, []byte{0x71, 0xa0}},       // ECI(26), TERM
		{3, []byte{0x50}},             // FNC1-1stPos, TERM
		{5, []byte{0x90}},             // FNC1-2ndPos, TERM
		{1, []byte{0x00}},
	}
	for _, test := range tests {
		r, e := DecodedBitStreamParser_Decode(test.bytes, VERSIONS[0], ErrorCorrectionLevel_L, nil)
		if e != nil {
			t.Fatalf("decode error (%v): %v", test.symbologyModifier, e)
		}
		if sm := r.GetSymbologyModifier(); sm != test.symbologyModifier {
			t.Fatalf("symbologyModifier = %v, wants %v", sm, test.symbologyModifier)
		}
	}
}
