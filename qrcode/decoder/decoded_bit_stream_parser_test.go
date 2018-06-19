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
	// 123456 = 0x07b*1000 +  0x1c8 => 00011110110111001000
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
