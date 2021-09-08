package common

import (
	"testing"

	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode"

	"github.com/makiuchi-d/gozxing"
)

func TestCharacterSetECI(t *testing.T) {
	_, e := GetCharacterSetECIByValue(-1)
	if e == nil {
		t.Fatalf("GetCharacterSetECIByValue(-1) must be error")
	}
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("GetCharacterSetECIByValue(-1) must be FormatException, %T", e)
	}

	_, e = GetCharacterSetECIByValue(900)
	if e == nil {
		t.Fatalf("GetCharacterSetECIByValue(900) must be error")
	}
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("GetCharacterSetECIByValue(900) must be FormatException, %T", e)
	}

	c, e := GetCharacterSetECIByValue(0)
	if e != nil {
		t.Fatalf("GetCharacterSetECIByValue(0) returns error, %v", e)
	}
	if c == nil {
		t.Fatalf("GetCharacterSetECIByValue(0) returns nil without error")
	}
	if r := c.GetValue(); r != 0 {
		t.Fatalf("GetCharacterSetECIByValue(0) value = %v, expect 0", r)
	}

	c, e = GetCharacterSetECIByValue(3)
	if e != nil {
		t.Fatalf("GetCharacterSetECIByValue(3) returns error, %v", e)
	}
	if c == nil {
		t.Fatalf("GetCharacterSetECIByValue(3) returns nil without error")
	}
	if r := c.GetValue(); r != 1 {
		t.Fatalf("GetCharacterSetECIByValue(3) value = %v, expect 1", r)
	}

	c, e = GetCharacterSetECIByValue(899)
	if e != nil {
		t.Fatalf("GetCharacterSetECIByValue(899) returns error, %v", e)
	}
	if c != nil {
		t.Fatalf("GetCharacterSetECIByValue(899) must return nil without error")
	}

	c, ok := GetCharacterSetECIByName("ISO-8859-1")
	if !ok || c == nil {
		t.Fatalf("GetCharacterSetECIByName(ISO-8859-1) returns nil")
	}
	if r := c.GetValue(); r != 1 {
		t.Fatalf("GetCharacterSetECIByValue(ISO-8859-1) value = %v, expect 1", r)
	}

	c, ok = GetCharacterSetECIByName("UTF-16BE")
	if !ok || c == nil {
		t.Fatalf("GetCharacterSetECIByName(UTF-16BE) returns nil")
	}
	if r := c.GetValue(); r != 25 {
		t.Fatalf("GetCharacterSetECIByValue(UTF-16E) value = %v, expect 25", r)
	}

	c, ok = GetCharacterSetECIByName("UnicodeBig")
	if !ok || c == nil {
		t.Fatalf("GetCharacterSetECIByName(UnicodeBig) returns nil")
	}
	if r := c.GetValue(); r != 25 {
		t.Fatalf("GetCharacterSetECIByValue(UnicodeBig) value = %v, expect 25", r)
	}

	c, ok = GetCharacterSetECIByName("UTF-8")
	if !ok || c == nil {
		t.Fatalf("GetCharacterSetECIByName(UTF-8) returns nil")
	}
	if r := c.GetValue(); r != 26 {
		t.Fatalf("GetCharacterSetECIByValue(UTF-8) value = %v, expect 26", r)
	}

	c, ok = GetCharacterSetECIByName("")
	if ok || c != nil {
		t.Fatalf("GetCharacterSetECIByName(\"\") must be nil")
	}

	_, ok = GetCharacterSetECI(nil)
	if ok {
		t.Fatalf("GetCharacterSetECI(nil) must be not ok")
	}

	c, ok = GetCharacterSetECI(unicode.UTF8)
	if !ok || c == nil {
		t.Fatalf("GetCharacterSetECI(UTF8) failed")
	}
	if r := c.GetValue(); r != 26 {
		t.Fatalf("GetCharacterSetECI(UTF-8) value = %v, expect 26", r)
	}
}

func TestCharasetECIName(t *testing.T) {
	if n := CharacterSetECI_Cp437.Name(); n != "Cp437" {
		t.Fatalf("invalid Cp437 name, %v", n)
	}
	if n := CharacterSetECI_ISO8859_1.Name(); n != "ISO-8859-1" {
		t.Fatalf("invalid ISO-8859-1 name, %v", n)
	}
	if n := CharacterSetECI_ISO8859_2.Name(); n != "ISO-8859-2" {
		t.Fatalf("invalid ISO-8859-2 name, %v", n)
	}
	if n := CharacterSetECI_ISO8859_3.Name(); n != "ISO-8859-3" {
		t.Fatalf("invalid ISO-8859-3 name, %v", n)
	}
	if n := CharacterSetECI_ISO8859_4.Name(); n != "ISO-8859-4" {
		t.Fatalf("invalid ISO-8859-4 name, %v", n)
	}
	if n := CharacterSetECI_ISO8859_5.Name(); n != "ISO-8859-5" {
		t.Fatalf("invalid ISO-8859-5 name, %v", n)
	}
	//if n := CharacterSetECI_ISO8859_6.Name(); n != "ISO-8859-6" {
	//	t.Fatalf("invalid ISO-8859-6 name, %v", n)
	//}
	if n := CharacterSetECI_ISO8859_7.Name(); n != "ISO-8859-7" {
		t.Fatalf("invalid ISO-8859-7 name, %v", n)
	}
	//if n := CharacterSetECI_ISO8859_8.Name(); n != "ISO-8859-8" {
	//	t.Fatalf("invalid ISO-8859-8 name, %v", n)
	//}
	if n := CharacterSetECI_ISO8859_9.Name(); n != "ISO-8859-9" {
		t.Fatalf("invalid ISO-8859-9 name, %v", n)
	}
	//if n := CharacterSetECI_ISO8859_10.Name(); n != "ISO-8859-10" {
	//	t.Fatalf("invalid ISO-8859-10 name, %v", n)
	//}
	if n := CharacterSetECI_ISO8859_13.Name(); n != "ISO-8859-13" {
		t.Fatalf("invalid ISO-8859-13 name, %v", n)
	}
	//if n := CharacterSetECI_ISO8859_14.Name(); n != "ISO-8859-14" {
	//	t.Fatalf("invalid ISO-8859-14 name, %v", n)
	//}
	if n := CharacterSetECI_ISO8859_15.Name(); n != "ISO-8859-15" {
		t.Fatalf("invalid ISO-8859-15 name, %v", n)
	}
	if n := CharacterSetECI_ISO8859_16.Name(); n != "ISO-8859-16" {
		t.Fatalf("invalid ISO-8859-16 name, %v", n)
	}
	if n := CharacterSetECI_SJIS.Name(); n != "Shift_JIS" {
		t.Fatalf("invalid Shift_JIS name, %v", n)
	}
	if n := CharacterSetECI_Cp1250.Name(); n != "windows-1250" {
		t.Fatalf("invalid windows-1250 name, %v", n)
	}
	if n := CharacterSetECI_Cp1251.Name(); n != "windows-1251" {
		t.Fatalf("invalid windows-1251 name, %v", n)
	}
	if n := CharacterSetECI_Cp1252.Name(); n != "windows-1252" {
		t.Fatalf("invalid windows-1252 name, %v", n)
	}
	if n := CharacterSetECI_Cp1256.Name(); n != "windows-1256" {
		t.Fatalf("invalid windows-1256 name, %v", n)
	}
	if n := CharacterSetECI_UnicodeBigUnmarked.Name(); n != "UTF-16BE" {
		t.Fatalf("invalid UTF-16BE name, %v", n)
	}
	if n := CharacterSetECI_UTF8.Name(); n != "UTF-8" {
		t.Fatalf("invalid UTF-8 name, %v", n)
	}
	if n := CharacterSetECI_ASCII.Name(); n != "ASCII" {
		t.Fatalf("invalid ASCII name, %v", n)
	}
	if n := CharacterSetECI_Big5.Name(); n != "Big5" {
		t.Fatalf("invalid Big5 name, %v", n)
	}
	if n := CharacterSetECI_GB18030.Name(); n != "GB18030" {
		t.Fatalf("invalid GB18030 name, %v", n)
	}
	if n := CharacterSetECI_EUC_KR.Name(); n != "EUC-KR" {
		t.Fatalf("invalid EUC-KR name, %v", n)
	}
}

func TestCharacterSetECI_GetEncoderDecoder(t *testing.T) {
	for _, charsetECI := range valueToECI {
		name := charsetECI.Name()
		enc, e := ianaindex.IANA.Encoding(name)
		if name == "ASCII" {
			continue
		}
		if e != nil {
			t.Fatalf("IANA.Encoding(%s) returns error, %v", name, e)
		}
		if enc.NewEncoder() == nil {
			t.Fatalf("%s encoder is nil", name)
		}
		if enc.NewDecoder() == nil {
			t.Fatalf("%s encoder is nil", name)
		}
	}
}
