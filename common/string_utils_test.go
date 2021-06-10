package common

import (
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode"

	"github.com/makiuchi-d/gozxing"
)

func TestGessEncodingWithHint(t *testing.T) {
	hints := map[gozxing.DecodeHintType]interface{}{}

	hints[gozxing.DecodeHintType_CHARACTER_SET] = unicode.UTF8
	wants := unicode.UTF8
	guessedCharset, err := StringUtils_guessCharset([]byte{0x21}, hints)
	if err != nil {
		t.Fatalf("guessCharset error: %v", err)
	}
	if guessedCharset != wants {
		t.Fatalf("guessedCharset is %v, wants %v", guessedCharset, wants)
	}

	hints[gozxing.DecodeHintType_CHARACTER_SET] = "ASCII"
	wants, _ = ianaindex.IANA.Encoding("US-ASCII")
	guessedCharset, err = StringUtils_guessCharset([]byte{0x21}, hints)
	if err != nil {
		t.Fatalf("guessCharset error: %v", err)
	}
	if guessedCharset != wants {
		t.Fatalf("guessedCharset is %v, wants %v", guessedCharset, wants)
	}

	hints[gozxing.DecodeHintType_CHARACTER_SET] = "ISO-2022-JP"
	guessedName, err := StringUtils_guessEncoding([]byte{0x21}, hints)
	if err != nil {
		t.Fatalf("guessCharset error: %v", err)
	}
	if guessedName != "ISO-2022-JP" {
		t.Fatalf("guessedCharset is %v, wants %v", guessedCharset, "ISO-2022-JP")
	}

	hints[gozxing.DecodeHintType_CHARACTER_SET] = "Dummy"
	guessedName, err = StringUtils_guessEncoding([]byte{0x21}, hints)
	if err == nil {
		t.Fatalf("guessEncoding must be error: name=%v", guessedName)
	}
}

func doTest(t testing.TB, bytes []byte, charset encoding.Encoding, encoding string) {
	t.Helper()
	guessedCharset, err := StringUtils_guessCharset(bytes, nil)
	if err != nil {
		t.Fatalf("guessCharset error: %v", err)
	}
	guessedName, err := StringUtils_guessEncoding(bytes, nil)
	if err != nil {
		t.Fatalf("guessEncoding error: %v", err)
	}
	if guessedCharset != charset {
		t.Fatalf("guessedCharset is %v, expect %v", guessedCharset, charset)
	}
	if guessedName != encoding {
		t.Fatalf("guessedName is %v, expect %v", guessedName, encoding)
	}
}

func TestShortShiftJIS1(t *testing.T) {
	// 金魚
	doTest(t, []byte{0x8b, 0xe0, 0x8b, 0x9b}, StringUtils_SHIFT_JIS_CHARSET, "SJIS")
}

func TestShortISO985911(t *testing.T) {
	// båd
	doTest(t, []byte{0x62, 0xe5, 0x64}, charmap.ISO8859_1, "ISO8859_1")
}

func TestMixedShiftJIS1(t *testing.T) {
	// Hello 金!
	doTest(t, []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x8b, 0xe0, 0x21}, StringUtils_SHIFT_JIS_CHARSET, "SJIS")
}

func TestSJIS(t *testing.T) {
	// ｶﾀｶﾅ
	doTest(t, []byte{0xb6, 0xc0, 0xb6, 0xC5}, StringUtils_SHIFT_JIS_CHARSET, "SJIS")
	// 邂遘
	doTest(t, []byte{0xe7, 0xae, 0xe7, 0xa7}, StringUtils_SHIFT_JIS_CHARSET, "SJIS")
}

func TestLatin1(t *testing.T) {
	// Weiß
	doTest(t, []byte{0x57, 0x65, 0x69, 0xdf}, charmap.ISO8859_1, "ISO8859_1")
	// fræ
	doTest(t, []byte{0x66, 0x72, 0xe6}, charmap.ISO8859_1, "ISO8859_1")
}

func TestUTF8(t *testing.T) {
	// 🍣×🍻
	doTest(t, []byte{0xF0, 0x9F, 0x8D, 0xA3, 0xC3, 0x97, 0xF0, 0x9F, 0x8D, 0xBA}, unicode.UTF8, "UTF8")
}

func TestUTF16withBOM(t *testing.T) {
	// 调压柜
	doTest(t, []byte{0xfe, 0xff, 0x8c, 0x03, 0x53, 0x8b, 0x67, 0xdc},
		unicode.UTF16(unicode.BigEndian, unicode.UseBOM), "UTF-16")
	// 调压柜
	doTest(t, []byte{0xff, 0xfe, 0x03, 0x8c, 0x8b, 0x53, 0xdc, 0x67},
		unicode.UTF16(unicode.LittleEndian, unicode.UseBOM), "UTF-16")
}

func TestUnknown(t *testing.T) {
	guessed, _ := StringUtils_guessCharset([]byte{0xe1, 0xff, 0xff, 0xf8, 0x81}, nil)
	wants := StringUtils_PLATFORM_DEFAULT_ENCODING
	if guessed != wants {
		t.Fatalf("guessedEncoding is %v, expect %v", guessed, wants)
	}
}
