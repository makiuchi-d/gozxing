package encoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

func TestEncoder_calculateMaskPenalty(t *testing.T) {
	matrix := makeCheckerMatrix(15)
	matrix.SetBool(5, 4, true)
	matrix.SetBool(6, 5, true)
	matrix.SetBool(1, 14, true)
	matrix.SetBool(3, 14, true)

	expect := 3 + 3 + 40 + 0
	if r := calculateMaskPenalty(matrix); r != expect {
		t.Fatalf("calculateMaskPenalty = %v, expect %v", r, expect)
	}

	for i := 8; i < 15; i++ {
		for j := 8; j < 15; j++ {
			matrix.SetBool(i, j, false)
		}
	}
	expect = 81 + 111 + 40 + 10
	if r := calculateMaskPenalty(matrix); r != expect {
		t.Fatalf("calculateMaskPenalty = %v, expect %v", r, expect)
	}
}

func TestEncoder_recommendVersion(t *testing.T) {
	headerBits := gozxing.NewBitArray(4) // mode indicator

	// 200 chars in alphanumeric mode
	// version=7, eclevel=L can contain 224 chars
	v, e := recommendVersion(
		decoder.ErrorCorrectionLevel_L,
		decoder.Mode_ALPHANUMERIC,
		headerBits,
		gozxing.NewBitArray((200/2)*11)) // 200 characters
	if e != nil {
		t.Fatalf("recommendVersion returns error: %v", e)
	}
	if v.GetVersionNumber() != 7 {
		t.Fatalf("recommendVersion = %v, expect 7", v)
	}

	// 1273 chars in byte mode.
	// the max chars of version=40, ecLevel=H
	v, e = recommendVersion(
		decoder.ErrorCorrectionLevel_H,
		decoder.Mode_BYTE,
		headerBits,
		gozxing.NewBitArray(1273*8))
	if e != nil {
		t.Fatalf("recommendVersion returns error: %v", e)
	}
	if v.GetVersionNumber() != 40 {
		t.Fatalf("recommendVersion = %v, expect 40", v)
	}

	// over capacity
	v, e = recommendVersion(
		decoder.ErrorCorrectionLevel_H,
		decoder.Mode_BYTE,
		headerBits,
		gozxing.NewBitArray(1274*8))
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("recommendVersion must be WriterException, %T", e)
	}
}

func TestEncoder_isOnlyDoubleByteKanji(t *testing.T) {
	if isOnlyDoubleByteKanji("fræ") {
		t.Fatalf("isOnlyDoubleByteKanji(fræ) must be false")
	}
	if isOnlyDoubleByteKanji("金魚!") {
		t.Fatalf("isOnlyDoubleByteKanji(金魚!) must be false")
	}
	if isOnlyDoubleByteKanji("ｶﾀｶﾅ") {
		t.Fatalf("isOnlyDoubleByteKanji(ｶﾀｶﾅ) must be false")
	}
	if !isOnlyDoubleByteKanji("邂遘") {
		t.Fatalf("isOnlyDoubleByteKanji(邂遘) must be true")
	}
}

func TestEncoder_chooseMode(t *testing.T) {
	content := "漢字モード"
	expect := decoder.Mode_KANJI
	if m := chooseMode(content, japanese.ShiftJIS); m != expect {
		t.Fatalf("chooseMode(%v, Shift_JIS) = %v mode, expect %v mode", content, m, expect)
	}

	content = "12345"
	expect = decoder.Mode_NUMERIC
	if m := chooseMode(content, nil); m != expect {
		t.Fatalf("chooseMode(%v) = %v mode, expect %v mode", content, m, expect)
	}

	content = "12345ABCDE"
	expect = decoder.Mode_ALPHANUMERIC
	if m := chooseMode(content, nil); m != expect {
		t.Fatalf("chooseMode(%v) = %v mode, expect %v mode", content, m, expect)
	}

	content = "12345ABCDEabcde"
	expect = decoder.Mode_BYTE
	if m := chooseMode(content, nil); m != expect {
		t.Fatalf("chooseMode(%v) = %v mode, expect %v mode", content, m, expect)
	}

	content = ""
	expect = decoder.Mode_BYTE
	if m := chooseMode(content, nil); m != expect {
		t.Fatalf("chooseMode(%v) = %v mode, expect %v mode", content, m, expect)
	}
}

func TestEncoder_getNumDataBytesAndNumECBytesForBlockID(t *testing.T) {
	ver, _ := decoder.Version_GetVersionForNumber(5)
	ecbs := ver.GetECBlocksForLevel(decoder.ErrorCorrectionLevel_H)
	numTotalBytes := ver.GetTotalCodewords() // 134
	numRSBlocks := ecbs.GetNumBlocks()       // 4

	_, _, e := getNumDataBytesAndNumECBytesForBlockID(numTotalBytes, 10, numRSBlocks, 5)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("getNumDataBytesAndNumECBytesForBlockID must be WriterException, %T", e)
	}

	// group1
	// dbytes = 10 % 4 = 2
	// ecbytes = (134/4) - dbytes = 31
	dbytes, ecbytes, e := getNumDataBytesAndNumECBytesForBlockID(numTotalBytes, 10, numRSBlocks, 0)
	if e != nil {
		t.Fatalf("getNumDataBytesAndNumECBytesForBlockID returns error, %v", e)
	}
	if dbytes != 2 || ecbytes != 31 {
		t.Fatalf("getNumDataBytesAndNumECBytesForBlockID = %v,%v, expect 2,31", dbytes, ecbytes)
	}

	// group2
	// dbytes = dbytes_group1 + 1 = 3
	// ecbytes = ecbytes_group1 = 31
	dbytes, ecbytes, e = getNumDataBytesAndNumECBytesForBlockID(numTotalBytes, 10, numRSBlocks, 3)
	if e != nil {
		t.Fatalf("getNumDataBytesAndNumECBytesForBlockID returns error, %v", e)
	}
	if dbytes != 3 || ecbytes != 31 {
		t.Fatalf("getNumDataBytesAndNumECBytesForBlockID = %v,%v, expect 3,31", dbytes, ecbytes)
	}
}

func TestEncoder_generateECBytes(t *testing.T) {
	_, e := generateECBytes([]byte{1, 2, 3}, 0)
	if e == nil {
		t.Fatalf("generateECBytes must be error")
	}

	ecbytes, e := generateECBytes([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, 5)
	expect := []int{52, 151, 95, 170, 87}
	if e != nil {
		t.Fatalf("generateECBytes returns error, %v", e)
	}
	if reflect.DeepEqual(ecbytes, expect) {
		t.Fatalf("generateECBytes = %v, expect %v", ecbytes, expect)
	}
}

func TestEncoder_interleaveWithECBytes(t *testing.T) {
	ver, _ := decoder.Version_GetVersionForNumber(3)
	totalcw := ver.GetTotalCodewords()
	numecb := ver.GetECBlocksForLevel(decoder.ErrorCorrectionLevel_H).GetNumBlocks()

	bits := gozxing.NewEmptyBitArray()

	_, e := interleaveWithECBytes(bits, totalcw, 1, numecb)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("interleaveWithECBytes must be WriterException, %T", e)
	}

	_, e = interleaveWithECBytes(bits, totalcw, 0, numecb)
	if e == nil {
		t.Fatalf("interleaveWithECBytes must be error")
	}

	// data bytes missmatch
	bits = gozxing.NewBitArray(40)
	_, e = interleaveWithECBytes(bits, totalcw, bits.GetSizeInBytes(), numecb)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("interleaveWithECBytes must be WriterException, %T", e)
	}

	bits = gozxing.NewEmptyBitArray()
	bits.AppendBits(0x12345678, 32)
	finalbits, e := interleaveWithECBytes(bits, totalcw, bits.GetSizeInBytes(), numecb)
	if e != nil {
		t.Fatalf("interleaveWithECBytes returns error, %v", e)
	}
	expects := "" +
		"...X..X..X.X.XX...XX.X...XXXX......XX...XXXX..XXXXXX..X...XXXX.X" +
		"XX.XX..X..X.X.X..X..XXXX..X.X.XXX..XX......X..X....XXXX....X..XX" +
		"XXXXX.XXX..XX...X....X..X.XXX...X.XXX.X.XXX.XXX...XXX.....XXXX.X" +
		".XXX.XXX...X.XX....XX..XXXXXX...X.XX.XX......X.XX..X.XXX..X.XX.X" +
		".X.XX.XX.X.XXXX.X.X....X.XXX..X..X..X.XXXX..X.XXX.X.X..X.X......" +
		"XXXX.X.X.....X.X..XXXXX....XXXX.X.XXX...XX.X....XXX.XXX...X.X.X." +
		"X..X.......XXXXX.X..X.....XXXXXX...X.XX....X.XX....XX.XXXX.X.X.X" +
		"..X....XXX...XX.XXX.X..XXX..XX......X..XX...XXX.XX....XX.XX.X.X." +
		"X..X...X..XXX.XXXX..XX.XXXX.XX...XX......X..X..X"
	if r := finalbits.GetSize(); r != len(expects) {
		t.Fatalf("interleaveWithECBytes result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := finalbits.Get(i), expects[i] == 'X'; r != expect {
			t.Fatalf("interleaveWithECBytes result[%v] = %v, expect %v", i, r, expect)
		}
	}
}

func TestEncoder_terminateBits(t *testing.T) {
	var bits *gozxing.BitArray
	var e error

	bits = gozxing.NewBitArray(128)
	e = terminateBits(15, bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("terminateBits must be WriterException, %T", e)
	}

	bits = gozxing.NewEmptyBitArray()
	bits.AppendBits(0x1234, 16)
	e = terminateBits(8, bits)

	buf := make([]byte, 8)
	bits.ToBytes(0, buf, 0, 8)
	if buf[0] != 0x12 {
		t.Fatalf("terminateBits array bytes[0] = 0x%02x, expect 0x12", buf[0])
	}
	if buf[1] != 0x34 {
		t.Fatalf("terminateBits array bytes[1] = 0x%02x, expect 0x34", buf[0])
	}
	if buf[2] != 0x00 {
		t.Fatalf("terminateBits array bytes[2] = 0x%02x, expect 0x00", buf[0])
	}
	if buf[3] != 0xec {
		t.Fatalf("terminateBits array bytes[3] = 0x%02x, expect 0xec", buf[0])
	}
	if buf[4] != 0x11 {
		t.Fatalf("terminateBits array bytes[4] = 0x%02x, expect 0x11", buf[0])
	}
	if buf[5] != 0xec {
		t.Fatalf("terminateBits array bytes[5] = 0x%02x, expect 0xec", buf[0])
	}
	if buf[6] != 0x11 {
		t.Fatalf("terminateBits array bytes[6] = 0x%02x, expect 0x11", buf[0])
	}
	if buf[7] != 0xec {
		t.Fatalf("terminateBits array bytes[7] = 0x%02x, expect 0xec", buf[0])
	}
}

func TestEncoder_appendNumericBytes(t *testing.T) {
	bits := gozxing.NewEmptyBitArray()
	appendNumericBytes("12345", bits)
	// 123 = 0001111011
	// 45  = 0101101
	expects := "00011110110101101"
	if r := bits.GetSize(); r != len(expects) {
		t.Fatalf("appendNumericBytes result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := bits.Get(i), expects[i] == '1'; r != expect {
			t.Fatalf("appendNumericBytes result[%v] = %v, expect %v", i, r, expect)
		}
	}

	bits = gozxing.NewEmptyBitArray()
	appendNumericBytes("1234", bits)
	// 123 = 0001111011
	// 4  = 0100
	expects = "00011110110100"
	if r := bits.GetSize(); r != len(expects) {
		t.Fatalf("appendNumericBytes result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := bits.Get(i), expects[i] == '1'; r != expect {
			t.Fatalf("appendNumericBytes result[%v] = %v, expect %v", i, r, expect)
		}
	}
}

func TestEncoder_appendAlphanumericBytes(t *testing.T) {
	var e error
	bits := gozxing.NewEmptyBitArray()

	e = appendAlphanumericBytes("a", bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("appendAlphanumericBytes must be WriterException, %T", e)
	}
	e = appendAlphanumericBytes("1a", bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("appendAlphanumericBytes must be WriterException, %T", e)
	}

	e = appendAlphanumericBytes("%1A", bits)
	if e != nil {
		t.Fatalf("appendAlphanumericBytes returns error, %v", e)
	}
	// %1 = 38*45 + 1 = 11010101111
	// A  = 10        = 001010
	expects := "11010101111001010"
	if r := bits.GetSize(); r != len(expects) {
		t.Fatalf("appendAlphanumericBytes result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := bits.Get(i), expects[i] == '1'; r != expect {
			t.Fatalf("appendAlphanumericBytes result[%v] = %v, expect %v", i, r, expect)
		}
	}
}

func TestEncoder_append8BitBytes(t *testing.T) {
	var e error
	bits := gozxing.NewEmptyBitArray()

	e = append8BitBytes("金魚", bits, charmap.ISO8859_1)
	if e == nil {
		t.Fatalf("append8BitBytes must be error")
	}

	e = append8BitBytes("金魚", bits, common.StringUtils_SHIFT_JIS_CHARSET)
	if e != nil {
		t.Fatalf("append8BitBytes returns error, %v", e)
	}
	// 0x8b, 0xe0, 0x8b, 0x9b = 10001011 11100000 10001011 10011011
	expects := "10001011111000001000101110011011"
	if r := bits.GetSize(); r != len(expects) {
		t.Fatalf("append8BitBytes result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := bits.Get(i), expects[i] == '1'; r != expect {
			t.Fatalf("append8BitBytes result[%v] = %v, expect %v", i, r, expect)
		}
	}
}

func TestEncoder_appendKanjiBytes(t *testing.T) {
	var e error
	bits := gozxing.NewEmptyBitArray()

	e = appendKanjiBytes("båd", bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("appendKanjiBytes must be WriterException, %T", e)
	}

	e = appendKanjiBytes("金魚!", bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("appendKanjiBytes must be WriterException, %T", e)
	}

	e = appendKanjiBytes("ｶﾀｶﾅ", bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("appendKanjiBytes must be WriterException, %T", e)
	}

	e = appendKanjiBytes("金魚邂遘", bits)
	if e != nil {
		t.Fatalf("appendKanjiBytes returns error, %v", e)
	}
	// 0x8be0 - 0x8140 = 0xaa0, 0xa*0xc0 + 0xa0 = 0x820 = 0100000100000
	// 0x8b9b - 0x8140 = 0xa5b, 0xa*0xc0 + 0x5b = 0x7db = 0011111011011
	// 0xe7ae - 0xc140 = 0x266e, 0x26*0xc0 + 0x6e = 0x1cee = 1110011101110
	// 0xe7a7 - 0xc140 = 0x2667, 0x26*0xc0 + 0x67 = 0x1ce7 = 1110011100111
	expects := "0100000100000001111101101111100111011101110011100111"
	if r := bits.GetSize(); r != len(expects) {
		t.Fatalf("appendKanjiBytes result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := bits.Get(i), expects[i] == '1'; r != expect {
			t.Fatalf("appendKanjiBytes result[%v] = %v, expect %v", i, r, expect)
		}
	}
}

func TestEncoder_appendLengthInfo(t *testing.T) {
	// kanji-mode version1: 8bits for length
	ver, _ := decoder.Version_GetVersionForNumber(1)
	bits := gozxing.NewEmptyBitArray()
	e := appendLengthInfo(256, ver, decoder.Mode_KANJI, bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("appendLengthInfo must be WriterException, %T", e)
	}
	e = appendLengthInfo(10, ver, decoder.Mode_KANJI, bits)
	if e != nil {
		t.Fatalf("appendLengthInfo returns error, %v", e)
	}
	expects := "00001010"
	if r := bits.GetSize(); r != len(expects) {
		t.Fatalf("appendLengthInfo result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := bits.Get(i), expects[i] == '1'; r != expect {
			t.Fatalf("appendLengthInfo result[%v] = %v, expect %v", i, r, expect)
		}
	}

	// numeric-mode version28: 14bits for length
	ver, _ = decoder.Version_GetVersionForNumber(28)
	bits = gozxing.NewEmptyBitArray()
	e = appendLengthInfo(16384, ver, decoder.Mode_NUMERIC, bits)
	if _, ok := e.(gozxing.WriterException); !ok {
		t.Fatalf("appendLengthInfo must be WriterException, %T", e)
	}
	e = appendLengthInfo(3669, ver, decoder.Mode_NUMERIC, bits)
	if e != nil {
		t.Fatalf("appendLengthInfo returns error, %v", e)
	}
	expects = "00111001010101"
	if r := bits.GetSize(); r != len(expects) {
		t.Fatalf("appendLengthInfo result size = %v, expect %v", r, len(expects))
	}
	for i := 0; i < len(expects); i++ {
		if r, expect := bits.Get(i), expects[i] == '1'; r != expect {
			t.Fatalf("appendLengthInfo result[%v] = %v, expect %v", i, r, expect)
		}
	}
}

func TestEncoder_chooseMaskPattern(t *testing.T) {
	eclevel := decoder.ErrorCorrectionLevel_M
	ver, _ := decoder.Version_GetVersionForNumber(1)

	matrix := NewByteMatrix(15, 15)
	bits := gozxing.NewEmptyBitArray()
	_, e := chooseMaskPattern(bits, eclevel, ver, matrix)
	if e == nil {
		t.Fatalf("chooseMaskPattern must be error")
	}

	bits = gozxing.NewEmptyBitArray()
	for i := 0; i < 25; i++ {
		bits.AppendBits(0xaa, 8)
	}
	matrix = NewByteMatrix(21, 21)
	mask, e := chooseMaskPattern(bits, eclevel, ver, matrix)
	if e != nil {
		t.Fatalf("chooseMaskPattern returns error, %v", e)
	}
	if mask != 1 {
		t.Fatalf("chooseMaskPattern = %v, expect 1", mask)
	}

	bits = gozxing.NewEmptyBitArray()
	for i := 0; i < 25; i++ {
		bits.AppendBits(0xa5, 8)
	}
	matrix = NewByteMatrix(21, 21)
	mask, e = chooseMaskPattern(bits, eclevel, ver, matrix)
	if e != nil {
		t.Fatalf("chooseMaskPattern returns error, %v", e)
	}
	if mask != 7 {
		t.Fatalf("chooseMaskPattern = %v, expect 7", mask)
	}
}

func testDecode(t testing.TB, qr *QRCode, content string) {
	t.Helper()
	matrix, _ := gozxing.ParseStringToBitMatrix(qr.matrix.String(), " 1", " 0")
	result, e := decoder.NewDecoder().Decode(matrix, nil)
	if e != nil {
		t.Fatalf("testDecode(%v) cannot decode, %v", content, e)
	}
	if result.GetText() != content {
		t.Fatalf("testDecode(%v) result missmatch, %v", content, result.GetText())
	}
}

func TestEncoder_encode(t *testing.T) {
	hints := map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_CHARACTER_SET: "SJIS"}
	qr, e := Encoder_encode("漢字モード", decoder.ErrorCorrectionLevel_H, hints)
	if e != nil {
		t.Fatalf("encode returns error, %v", e)
	}
	if r := qr.GetMode(); r != decoder.Mode_KANJI {
		t.Fatalf("encoded mode = %v, expect %v", r, decoder.Mode_KANJI)
	}
	if r := qr.GetECLevel(); r != decoder.ErrorCorrectionLevel_H {
		t.Fatalf("encoded ecLevel = %v, expect %v", r, decoder.ErrorCorrectionLevel_H)
	}
	testDecode(t, qr, "漢字モード")

	hints = map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_CHARACTER_SET: "UTF-8"}
	qr, e = Encoder_encode("8Byteモード", decoder.ErrorCorrectionLevel_M, hints)
	if e != nil {
		t.Fatalf("encode returns error, %v", e)
	}
	if r := qr.GetMode(); r != decoder.Mode_BYTE {
		t.Fatalf("encoded mode = %v, expect %v", r, decoder.Mode_BYTE)
	}
	if r := qr.GetECLevel(); r != decoder.ErrorCorrectionLevel_M {
		t.Fatalf("encoded ecLevel = %v, expect %v", r, decoder.ErrorCorrectionLevel_M)
	}
	testDecode(t, qr, "8Byteモード")

	hints = map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_GS1_FORMAT: "True"}
	qr, e = Encoder_encode("01049123451234591597033130128%10ABC123", decoder.ErrorCorrectionLevel_Q, hints)
	if e != nil {
		t.Fatalf("encode returns error, %v", e)
	}
	if r := qr.GetMode(); r != decoder.Mode_ALPHANUMERIC {
		t.Fatalf("encoded mode = %v, expect %v", r, decoder.Mode_ALPHANUMERIC)
	}
	if r := qr.GetECLevel(); r != decoder.ErrorCorrectionLevel_Q {
		t.Fatalf("encoded ecLevel = %v, expect %v", r, decoder.ErrorCorrectionLevel_Q)
	}
	testDecode(t, qr, "01049123451234591597033130128\x1d10ABC123")

	hints = map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_QR_VERSION: "7"}
	qr, e = Encoder_encode("01234567890", decoder.ErrorCorrectionLevel_L, hints)
	if e != nil {
		t.Fatalf("encode returns error, %v", e)
	}
	if r := qr.GetMode(); r != decoder.Mode_NUMERIC {
		t.Fatalf("encoded mode = %v, expect %v", r, decoder.Mode_NUMERIC)
	}
	if r := qr.GetECLevel(); r != decoder.ErrorCorrectionLevel_L {
		t.Fatalf("encoded ecLevel = %v, expect %v", r, decoder.ErrorCorrectionLevel_L)
	}
	if r := qr.GetVersion(); r.GetVersionNumber() != 7 {
		t.Fatalf("encoded version = %v, expect %v", r, 7)
	}
	testDecode(t, qr, "01234567890")

	qr, e = Encoder_encodeWithoutHint("http://example.com", decoder.ErrorCorrectionLevel_H)
	if r := qr.GetECLevel(); r != decoder.ErrorCorrectionLevel_H {
		t.Fatalf("encoded ecLevel = %v, expect %v", r, decoder.ErrorCorrectionLevel_H)
	}
	testDecode(t, qr, "http://example.com")

	hints = make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_QR_MASK_PATTERN] = 1
	qr, e = Encoder_encode("http://example.com", decoder.ErrorCorrectionLevel_H, hints)
	if e != nil {
		t.Fatalf("Encoder_encoder returns error: %v", e)
	}
	if mask, wants := qr.GetMaskPattern(), 1; mask != wants {
		t.Fatalf("Encoder_encode maskPattern = %v, wants %v", mask, wants)
	}

	hints[gozxing.EncodeHintType_QR_MASK_PATTERN] = "2"
	qr, e = Encoder_encode("http://example.com", decoder.ErrorCorrectionLevel_H, hints)
	if e != nil {
		t.Fatalf("Encoder_encoder returns error: %v", e)
	}
	if mask, wants := qr.GetMaskPattern(), 2; mask != wants {
		t.Fatalf("Encoder_encode maskPattern = %v, wants %v", mask, wants)
	}

	hints[gozxing.EncodeHintType_QR_MASK_PATTERN] = 10
	qr, e = Encoder_encode("http://example.com", decoder.ErrorCorrectionLevel_H, hints)
	if e != nil {
		t.Fatalf("Encoder_encoder returns error: %v", e)
	}
	if mask, wants := qr.GetMaskPattern(), 7; mask != wants {
		t.Fatalf("Encoder_encode maskPattern = %v, wants %v", mask, wants)
	}
}

func TestEncoder_encodeFail(t *testing.T) {
	// fail in append8BitBytes
	hints := map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_CHARACTER_SET: "ISO-8859-1"}
	_, e := Encoder_encode("エラー", decoder.ErrorCorrectionLevel_H, hints)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	// fail in GetVersionForNumber
	hints = map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_QR_VERSION: "41"}
	_, e = Encoder_encode("エラー", decoder.ErrorCorrectionLevel_H, hints)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	// willFit is false
	hints = map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_QR_VERSION: "1"}
	_, e = Encoder_encode("123456789012345678", decoder.ErrorCorrectionLevel_H, hints)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	// failed in recommendVersion
	hints = map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_CHARACTER_SET: "UTF-8"}
	content := "🍣🍻🍣🍻🍣🍻🍣🍻🍣🍻🍣🍻"
	for len(content) < 1273 {
		content = content + content
	}
	_, e = Encoder_encode(content, decoder.ErrorCorrectionLevel_H, hints)
	if e == nil {
		t.Fatalf("encode must be error")
	}

	hints = map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_CHARACTER_SET: "Dummy"}
	_, e = Encoder_encode("Dummy", decoder.ErrorCorrectionLevel_H, hints)
	if e == nil {
		t.Fatalf("encode must be error")
	}
}
