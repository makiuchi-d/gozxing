package decoder

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func testVersion(t testing.TB, row, col int,
	verNum, symbolRow, symbolCol, dataRegionRow, dataRegionCol, totalCW int,
	ecCW int, ecbs []ECB, str string) {
	t.Helper()

	v, e := getVersionForDimensions(row, col)
	if e != nil {
		t.Fatalf("getVersionForDimensions(%v, %v) returns error, %v", row, col, e)
	}

	if r := v.getVersionNumber(); r != verNum {
		t.Fatalf("versions(%v,%v) versionNumber = %v, expect %v", row, col, r, verNum)
	}
	if r := v.getSymbolSizeRows(); r != symbolRow {
		t.Fatalf("versions(%v,%v) symbolSizeRow = %v, expect %v", row, col, r, symbolRow)
	}
	if r := v.getSymbolSizeColumns(); r != symbolCol {
		t.Fatalf("versions(%v,%v) symbolSizeColumns = %v, expect %v", row, col, r, symbolCol)
	}
	if r := v.getDataRegionSizeRows(); r != dataRegionRow {
		t.Fatalf("versions(%v,%v) dataRegionSizeRow = %v, expect %v", row, col, r, dataRegionRow)
	}
	if r := v.getDataRegionSizeColumns(); r != dataRegionCol {
		t.Fatalf("versions(%v,%v) dataRegionSizeColumns = %v, expect %v", row, col, r, dataRegionCol)
	}
	if r := v.getTotalCodewords(); r != totalCW {
		t.Fatalf("versions(%v,%v) getTotalCodewords = %v, expect %v", row, col, r, totalCW)
	}
	ecb := v.getECBlocks()
	if r := ecb.getECCodewords(); r != ecCW {
		t.Fatalf("versions(%v,%v) ecCodewords = %v, expect %v", row, col, r, ecCW)
	}
	if r := ecb.getECBlocks(); !reflect.DeepEqual(r, ecbs) {
		t.Fatalf("versions(%v,%v) ecBlocks = %v, expect %v", row, col, r, ecbs)
	}

	if r := v.String(); r != str {
		t.Fatalf("versions(%v,%v) Strng = %v, expect %v", row, col, r, str)
	}
}

func TestVersion(t *testing.T) {
	_, e := getVersionForDimensions(11, 10)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("getVersionForDimensions(11, 10) must be FormatException, %T", e)
	}

	_, e = getVersionForDimensions(10, 11)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("getVersionForDimensions(10, 11) must be FormatException, %T", e)
	}

	_, e = getVersionForDimensions(146, 146)
	if _, ok := e.(gozxing.FormatException); !ok {
		t.Fatalf("getVersionForDimensions(146, 146) must be FormatException, %T", e)
	}

	// version 1
	testVersion(t, 10, 10, 1, 10, 10, 8, 8, 8, 5, []ECB{{1, 3}}, "1")

	// version 24
	testVersion(t, 144, 144, 24, 144, 144, 22, 22, 2178, 62, []ECB{{8, 156}, {2, 155}}, "24")

	// version 30
	testVersion(t, 16, 48, 30, 16, 48, 14, 22, 77, 28, []ECB{{1, 49}}, "30")
}
