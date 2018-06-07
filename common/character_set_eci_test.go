package common

import (
	"testing"

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

	c = GetCharacterSetECIByName("ISO-8859-1")
	if c == nil {
		t.Fatalf("GetCharacterSetECIByName(ISO-8859-1) returns nil")
	}
	if r := c.GetValue(); r != 1 {
		t.Fatalf("GetCharacterSetECIByValue(ISO-8859-1) value = %v, expect 1", r)
	}

	c = GetCharacterSetECIByName("UTF-16BE")
	if c == nil {
		t.Fatalf("GetCharacterSetECIByName(UTF-16BE) returns nil")
	}
	if r := c.GetValue(); r != 25 {
		t.Fatalf("GetCharacterSetECIByValue(UTF-16E) value = %v, expect 25", r)
	}

	c = GetCharacterSetECIByName("UnicodeBig")
	if c == nil {
		t.Fatalf("GetCharacterSetECIByName(UnicodeBig) returns nil")
	}
	if r := c.GetValue(); r != 25 {
		t.Fatalf("GetCharacterSetECIByValue(UnicodeBig) value = %v, expect 25", r)
	}

	c = GetCharacterSetECIByName("UTF-8")
	if c == nil {
		t.Fatalf("GetCharacterSetECIByName(UTF-8) returns nil")
	}
	if r := c.GetValue(); r != 26 {
		t.Fatalf("GetCharacterSetECIByValue(UTF-8) value = %v, expect 26", r)
	}

	c = GetCharacterSetECIByName("")
	if c != nil {
		t.Fatalf("GetCharacterSetECIByName(\"\") must be nil")
	}
}
