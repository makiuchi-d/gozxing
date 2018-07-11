package decoder

import (
	"testing"
)

func testErrorCorrectionLevel_ForBits(t *testing.T, bits uint, expect ErrorCorrectionLevel) {
	r, e := ErrorCorrectionLevel_ForBits(bits)
	if e != nil {
		t.Fatalf("ForBits(0) returns error, %v", e)
	}
	if r != expect {
		t.Fatalf("ForBits(%v) != %v, expect %v", bits, r, expect)
	}
}

func TestErrorCorrectionLevel_ForBits(t *testing.T) {

	if _, e := ErrorCorrectionLevel_ForBits(4); e == nil {
		t.Fatalf("ForBits(4) must be error")
	}

	testErrorCorrectionLevel_ForBits(t, 0, ErrorCorrectionLevel_M)
	testErrorCorrectionLevel_ForBits(t, 1, ErrorCorrectionLevel_L)
	testErrorCorrectionLevel_ForBits(t, 2, ErrorCorrectionLevel_H)
	testErrorCorrectionLevel_ForBits(t, 3, ErrorCorrectionLevel_Q)
}

func TestErrorCorrectionLevel_GetBits(t *testing.T) {
	if r := ErrorCorrectionLevel_M.GetBits(); r != 0 {
		t.Fatalf("M.GetBits = %v, expect 0", r)
	}
	if r := ErrorCorrectionLevel_L.GetBits(); r != 1 {
		t.Fatalf("L.GetBits = %v, expect 1", r)
	}
	if r := ErrorCorrectionLevel_H.GetBits(); r != 2 {
		t.Fatalf("H.GetBits = %v, expect 2", r)
	}
	if r := ErrorCorrectionLevel_Q.GetBits(); r != 3 {
		t.Fatalf("Q.GetBits = %v, expect 3", r)
	}
}

func TestErrorCorrectionLevel_String(t *testing.T) {
	if s := ErrorCorrectionLevel_L.String(); s != "L" {
		t.Fatalf("ErrorCorrectionLevel_L string is %v, expect L", s)
	}
	if s := ErrorCorrectionLevel_M.String(); s != "M" {
		t.Fatalf("ErrorCorrectionLevel_M string is %v, expect M", s)
	}
	if s := ErrorCorrectionLevel_Q.String(); s != "Q" {
		t.Fatalf("ErrorCorrectionLevel_Q string is %v, expect Q", s)
	}
	if s := ErrorCorrectionLevel_H.String(); s != "H" {
		t.Fatalf("ErrorCorrectionLevel_H string is %v, expect H", s)
	}
	if s := ErrorCorrectionLevel(-1).String(); s != "" {
		t.Fatalf("invalid ErrorCorrectionLevel string must be \"\", %v", s)
	}
}
