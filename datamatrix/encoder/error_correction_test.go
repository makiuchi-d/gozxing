package encoder

import (
	"reflect"
	"testing"
)

func TestCreateECCBlock(t *testing.T) {
	_, e := createECCBlock([]byte{}, 0)
	if e == nil {
		t.Fatalf("createECCBlock must be error")
	}

	b, e := createECCBlock([]byte{1, 2, 3, 4, 113}, 7)
	expect := []byte{207, 172, 22, 211, 108, 152, 0}
	if e != nil {
		t.Fatalf("createECCBlock returns error: %v", e)
	}
	if !reflect.DeepEqual(b, expect) {
		t.Fatalf("createECCBlock = %v, expect %v", b, expect)
	}
}

func TestEncodeECC200(t *testing.T) {
	_, e := ErrorCorrection_EncodeECC200([]byte{}, symbols[0])
	if e == nil {
		t.Fatalf("encodeECC200 must be error")
	}

	_, e = ErrorCorrection_EncodeECC200([]byte{1}, NewSymbolInfo(false, 1, 1, 10, 10, 1))
	if e == nil {
		t.Fatalf("encodeECC200 must be error")
	}

	cw := []byte{1, 2, 3}
	expect := []byte{1, 2, 3, 53, 46, 95, 0, 75}
	sb, e := ErrorCorrection_EncodeECC200(cw, symbols[0])
	if e != nil {
		t.Fatalf("encodeECC200 returns error: %v", e)
	}
	if !reflect.DeepEqual(sb, expect) {
		t.Fatalf("encodeECC200 = %v, expect %v", sb, expect)
	}

	cw = make([]byte, 204)
	for i := range cw {
		cw[i] = byte(i % 256)
	}
	expect = []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
		60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
		100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115,
		116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131,
		132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147,
		148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163,
		164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179,
		180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195,
		196, 197, 198, 199, 200, 201, 202, 203, 66, 250, 127, 172, 140, 182, 251, 108,
		195, 81, 219, 5, 68, 215, 97, 22, 237, 200, 215, 8, 137, 118, 91, 106, 105, 1,
		52, 172, 95, 125, 35, 242, 142, 53, 6, 198, 132, 225, 183, 248, 101, 102, 124,
		165, 207, 191, 240, 190, 189, 155, 124, 104, 108, 137, 132, 219, 59, 90, 18, 251,
		104, 115, 121, 48, 44, 225, 25, 62, 207, 94, 129, 153, 228, 143, 112, 127, 234,
		76, 117, 244, 144, 241, 209, 202,
	}
	sb, e = ErrorCorrection_EncodeECC200(cw, symbols[20])
	if e != nil {
		t.Fatalf("encodeECC200 returns error: %v", e)
	}
	if !reflect.DeepEqual(sb, expect) {
		t.Fatalf("encodeECC200 = %v, expect %v", sb, expect)
	}
}
