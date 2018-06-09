package common

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestGessEncodingWithHint(t *testing.T) {
	hints := map[gozxing.DecodeHintType]interface{}{}
	charsetName := "Dummy"
	hints[gozxing.DecodeHintType_CHARACTER_SET] = charsetName
	guessedName := StringUtils_guessEncoding([]byte{0x21}, hints)
	if guessedName != charsetName {
		t.Fatalf("guessedName is %v, expect %v", guessedName, charsetName)
	}
}

func doTest(t *testing.T, bytes []byte, charsetName string) {
	guessedName := StringUtils_guessEncoding(bytes, nil)
	if guessedName != charsetName {
		t.Fatalf("guessedName is %v, expect %v", guessedName, charsetName)
	}
}

func TestShortShiftJIS1(t *testing.T) {
	// 金魚
	doTest(t, []byte{0x8b, 0xe0, 0x8b, 0x9b}, StringUtils_SHIFT_JIS)
}

func TestShortISO985911(t *testing.T) {
	// båd
	doTest(t, []byte{0x62, 0xe5, 0x64}, StringUtils_ISO88591)
}

func TestMixedShiftJIS1(t *testing.T) {
	// Hello 金!
	doTest(t, []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x8b, 0xe0, 0x21}, StringUtils_SHIFT_JIS)
}

func TestSJIS(t *testing.T) {
	// ｶﾀｶﾅ
	doTest(t, []byte{0xb6, 0xc0, 0xb6, 0xC5}, StringUtils_SHIFT_JIS)
	// 邂遘
	doTest(t, []byte{0xe7, 0xae, 0xe7, 0xa7}, StringUtils_SHIFT_JIS)
}

func TestLatin1(t *testing.T) {
	// Weiß
	doTest(t, []byte{0x57, 0x65, 0x69, 0xdf}, StringUtils_ISO88591)
	// fræ
	doTest(t, []byte{0x66, 0x72, 0xe6}, StringUtils_ISO88591)
}

func TestUTF8(t *testing.T) {
	// 🍣×🍻
	doTest(t, []byte{0xF0, 0x9F, 0x8D, 0xA3, 0xC3, 0x97, 0xF0, 0x9F, 0x8D, 0xBA}, StringUtils_UTF8)
}

func TestUnknown(t *testing.T) {
	doTest(t, []byte{0xe1, 0xff, 0xff, 0xf8, 0x81}, StringUtils_PLATFORM_DEFAULT_ENCODING)
}
