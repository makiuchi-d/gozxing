package oned

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/testutil"
)

func TestUPCEReader_decodeMiddle(t *testing.T) {
	reader := &upcEReader{
		decodeMiddleCounters: make([]int, 4),
	}
	result := make([]byte, 4)
	startRange := []int{3, 6}

	row := gozxing.NewBitArray(60)

	_, _, e := reader.decodeMiddle(row, startRange, result[:0])
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeMiddle must be NotFoundException, %T", e)
	}

	// 0-123453-1
	// 0-1 : GGLGLL
	// [6-12] 1(G): 0110011
	row.Set(7)
	row.Set(8)
	row.Set(11)
	row.Set(12)
	// [13-19] 2(G): 0011011
	row.Set(15)
	row.Set(16)
	row.Set(18)
	row.Set(19)
	// [20-26] 3(L): 0111101
	row.Set(21)
	row.Set(22)
	row.Set(23)
	row.Set(24)
	row.Set(26)
	// [27-33] 4(G): 0011101
	row.Set(29)
	row.Set(30)
	row.Set(31)
	row.Set(33)
	// [34-40] 5(L): 0110001
	row.Set(35)
	row.Set(36)
	row.Set(40)
	// [41-47] 3(G, error): 0100001
	row.Set(42)
	row.Set(47)

	_, _, e = reader.decodeMiddle(row, startRange, result[:0])
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("decodeMiddle must be NotFoundException, %T", e)
	}

	// [41-47] 3(G): 0100001
	// [41-47] 3(L): 0111101
	row.Flip(43)
	row.Flip(44)
	row.Flip(45)

	offset, result, e := reader.decodeMiddle(row, startRange, result[:0])
	if e != nil {
		t.Fatalf("decodeMiddle returns error, %v", e)
	}
	if offset != 48 {
		t.Fatalf("decodeMiddle offset = %v, expect 48", offset)
	}
	if str := string(result); str != "01234531" {
		t.Fatalf("decodeMiddle result = \"%v\", expect \"01234531\"", str)
	}
}

func testConvertUPCEtoUPCA(t testing.TB, upce, expect string) {
	t.Helper()
	upca := convertUPCEtoUPCA(upce)
	if upca != expect {
		t.Fatalf("UPCEtoUPCE: %s => %s, expect %s", upce, upca, expect)
	}
}

func TestConvertUPCEtoUPCA(t *testing.T) {
	testConvertUPCEtoUPCA(t, "0123450", "01200000345")
	testConvertUPCEtoUPCA(t, "0123451", "01210000345")
	testConvertUPCEtoUPCA(t, "0123452", "01220000345")
	testConvertUPCEtoUPCA(t, "0123453", "01230000045")
	testConvertUPCEtoUPCA(t, "0123454", "01234000005")
	testConvertUPCEtoUPCA(t, "0123459", "01234500009")
	testConvertUPCEtoUPCA(t, "01234531", "012300000451")
}

func TestUPCEReader(t *testing.T) {
	// testdata from zxing core/src/test/resources/blackbox/upce-1/
	reader := NewUPCEReader()
	format := gozxing.BarcodeFormat_UPC_E

	tests := []struct {
		file  string
		wants string
	}{
		{"testdata/upce/1.png", "01234565"},
		{"testdata/upce/2.png", "00123457"},
		{"testdata/upce/4.png", "01234531"},
	}
	for _, test := range tests {
		testutil.TestFile(t, reader, test.file, test.wants, format, nil, nil)
	}
}
