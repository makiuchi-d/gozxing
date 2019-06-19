package decoder

import (
	"testing"
)

func testModeForBits(t testing.TB, bits int, expect *Mode) {
	t.Helper()
	mode, e := ModeForBits(bits)
	if e != nil {
		t.Fatalf("ModeForBits(%v) returns error, %v", bits, e)
	}
	if mode != expect {
		t.Fatalf("ModeForBits(%v) = %v, expect %v", bits, mode, expect)
	}
}

func TestModeForBits(t *testing.T) {
	testModeForBits(t, 0, Mode_TERMINATOR)
	testModeForBits(t, 1, Mode_NUMERIC)
	testModeForBits(t, 2, Mode_ALPHANUMERIC)
	testModeForBits(t, 3, Mode_STRUCTURED_APPEND)
	testModeForBits(t, 4, Mode_BYTE)
	testModeForBits(t, 5, Mode_FNC1_FIRST_POSITION)
	testModeForBits(t, 7, Mode_ECI)
	testModeForBits(t, 8, Mode_KANJI)
	testModeForBits(t, 9, Mode_FNC1_SECOND_POSITION)
	testModeForBits(t, 0xD, Mode_HANZI)

	if _, e := ModeForBits(6); e == nil {
		t.Fatalf("ModeForBits(6) must be error")
	}
}

func TestModeNumeric(t *testing.T) {
	mode := Mode_NUMERIC
	if r := mode.GetBits(); r != 1 {
		t.Fatalf("Bits = %v, expect 1", r)
	}

	ver, _ := Version_GetVersionForNumber(9)
	if r := mode.GetCharacterCountBits(ver); r != 10 {
		t.Fatalf("CharacterCountBits(ver=9) = %v, expect 10", r)
	}
	ver, _ = Version_GetVersionForNumber(10)
	if r := mode.GetCharacterCountBits(ver); r != 12 {
		t.Fatalf("CharacterCountBits(ver=10) = %v, expect 12", r)
	}
	ver, _ = Version_GetVersionForNumber(26)
	if r := mode.GetCharacterCountBits(ver); r != 12 {
		t.Fatalf("CharacterCountBits(ver=26) = %v, expect 12", r)
	}
	ver, _ = Version_GetVersionForNumber(27)
	if r := mode.GetCharacterCountBits(ver); r != 14 {
		t.Fatalf("CharacterCountBits(ver=23) = %v, expect 14", r)
	}
}

func TestModeAlphaNumeric(t *testing.T) {
	mode := Mode_ALPHANUMERIC
	if r := mode.GetBits(); r != 2 {
		t.Fatalf("Bits = %v, expect 2", r)
	}

	ver, _ := Version_GetVersionForNumber(9)
	if r := mode.GetCharacterCountBits(ver); r != 9 {
		t.Fatalf("CharacterCountBits(ver=9) = %v, expect 9", r)
	}
	ver, _ = Version_GetVersionForNumber(10)
	if r := mode.GetCharacterCountBits(ver); r != 11 {
		t.Fatalf("CharacterCountBits(ver=10) = %v, expect 11", r)
	}
	ver, _ = Version_GetVersionForNumber(26)
	if r := mode.GetCharacterCountBits(ver); r != 11 {
		t.Fatalf("CharacterCountBits(ver=26) = %v, expect 11", r)
	}
	ver, _ = Version_GetVersionForNumber(27)
	if r := mode.GetCharacterCountBits(ver); r != 13 {
		t.Fatalf("CharacterCountBits(ver=23) = %v, expect 13", r)
	}
}

func testMode_String(t testing.TB, mode *Mode, expect string) {
	t.Helper()
	str := mode.String()
	if str != expect {
		t.Fatalf("String = \"%v\", expect \"%v\"", str, expect)
	}
}

func TestMode_String(t *testing.T) {
	testMode_String(t, Mode_TERMINATOR, "TERMINATOR")
	testMode_String(t, Mode_NUMERIC, "NUMERIC")
	testMode_String(t, Mode_ALPHANUMERIC, "ALPHANUMERIC")
	testMode_String(t, Mode_STRUCTURED_APPEND, "STRUCTURED_APPEND")
	testMode_String(t, Mode_BYTE, "BYTE")
	testMode_String(t, Mode_ECI, "ECI")
	testMode_String(t, Mode_KANJI, "KANJI")
	testMode_String(t, Mode_FNC1_FIRST_POSITION, "FNC1_FIRST_POSITION")
	testMode_String(t, Mode_FNC1_SECOND_POSITION, "FNC1_SECOND_POSITION")
	testMode_String(t, Mode_HANZI, "HANZI")

	testMode_String(t, nil, "")
}
