package decoder

import (
	"testing"
)

func TestGetDataBlocks(t *testing.T) {
	version, _ := getVersionForDimensions(10, 10)
	rawCW := make([]byte, version.getTotalCodewords()+10)
	_, e := DataBlocks_getDataBlocks(rawCW, version)
	if e == nil {
		t.Fatalf("getDataBlocks must be error")
	}

	// version 1
	version, _ = getVersionForDimensions(10, 10)
	rawCW = make([]byte, version.getTotalCodewords())
	dbs, e := DataBlocks_getDataBlocks(rawCW, version)
	if e != nil {
		t.Fatalf("getDataBlocks(ver1) returns error, %v", e)
	}
	if r := len(dbs); r != 1 {
		t.Fatalf("getDataBlocks(ver1) len = %v, expect 1", r)
	}
	for _, d := range dbs {
		if r := d.getNumDataCodewords(); r != 3 {
			t.Fatalf("getDataBlocks(ver1) numDataCodewords = %v, expect 3", r)
		}
		if r := d.getCodewords(); len(r) != 8 {
			t.Fatalf("getDataBlocks(ver1) len(codewords) = %v, expect 8", len(r))
		}
	}

	// version 15
	version, _ = getVersionForDimensions(52, 52)
	rawCW = make([]byte, version.getTotalCodewords())
	dbs, e = DataBlocks_getDataBlocks(rawCW, version)
	if e != nil {
		t.Fatalf("getDataBlocks(ver1) returns error, %v", e)
	}
	if r := len(dbs); r != 2 {
		t.Fatalf("getDataBlocks(ver1) len = %v, expect 2", r)
	}
	for _, d := range dbs {
		if r := d.getNumDataCodewords(); r != 102 {
			t.Fatalf("getDataBlocks(ver1) numDataCodewords = %v, expect 102", r)
		}
		if r := d.getCodewords(); len(r) != 144 {
			t.Fatalf("getDataBlocks(ver1) len(codewords) = %v, expect 144", len(r))
		}
	}

	// version 24 (special version)
	version, _ = getVersionForDimensions(144, 144)
	rawCW = make([]byte, version.getTotalCodewords())
	dbs, e = DataBlocks_getDataBlocks(rawCW, version)
	if e != nil {
		t.Fatalf("getDataBlocks(ver1) returns error, %v", e)
	}
	if r := len(dbs); r != 8+2 {
		t.Fatalf("getDataBlocks(ver1) len = %v, expect 10", r)
	}
	for i := 0; i < 8; i++ {
		d := dbs[i]
		if r := d.getNumDataCodewords(); r != 156 {
			t.Fatalf("getDataBlocks(ver1) numDataCodewords = %v, expect 3", r)
		}
		if r := d.getCodewords(); len(r) != 62+156 {
			t.Fatalf("getDataBlocks(ver1) len(codewords) = %v, expect 62+156", len(r))
		}
	}
	for i := 8; i < 10; i++ {
		d := dbs[i]
		if r := d.getNumDataCodewords(); r != 155 {
			t.Fatalf("getDataBlocks(ver1) numDataCodewords = %v, expect 155", r)
		}
		if r := d.getCodewords(); len(r) != 62+155 {
			t.Fatalf("getDataBlocks(ver1) len(codewords) = %v, expect 62+155", len(r))
		}
	}
}
