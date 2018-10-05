package decoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestDecodeBoolMap(t *testing.T) {
	d := NewDecoder()

	_, e := d.DecodeBoolMap([][]bool{{}})
	if e == nil {
		t.Fatalf("DecodeBoolMap must be error")
	}

	bits, _ := gozxing.ParseStringToBitMatrix(dm7str, "##", "  ")
	bmap := make([][]bool, bits.GetHeight())
	for j := 0; j < len(bmap); j++ {
		bmap[j] = make([]bool, bits.GetWidth())
		for i := 0; i < len(bmap[j]); i++ {
			bmap[j][i] = bits.Get(i, j)
		}
	}
	r, e := d.DecodeBoolMap(bmap)
	if e != nil {
		t.Fatalf("DecodeBoolMap returns error, %v", e)
	}
	expect := "0123456789"
	if s := r.GetText(); s != expect {
		t.Fatalf("DecodeBoolMap text = \"%v\", expect \"%v\"", s, expect)
	}
}

func TestDecode(t *testing.T) {
	d := NewDecoder()

	bits, _ := gozxing.NewBitMatrix(1, 1)
	_, e := d.Decode(bits)
	if e == nil {
		t.Fatalf("Decode must be error")
	}

	bits, _ = gozxing.ParseStringToBitMatrix(dm4str, "##", "  ")
	r, e := d.Decode(bits)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	expect := "Hello World"
	if s := r.GetText(); s != expect {
		t.Fatalf("Decode text = \"%v\", expect \"%v\"", s, expect)
	}

	// error collection
	bits.SetRegion(3, 3, 5, 5)
	r, e = d.Decode(bits)
	if e != nil {
		t.Fatalf("Decode returns error, %v", e)
	}
	if s := r.GetText(); s != expect {
		t.Fatalf("Decode text = \"%v\", expect \"%v\"", s, expect)
	}

	// checksum exception
	bits.SetRegion(3, 3, 10, 10)
	r, e = d.Decode(bits)
	if _, ok := e.(gozxing.ChecksumException); !ok {
		t.Fatalf("Decode must be ChecksumException, %T", e)
	}
}
