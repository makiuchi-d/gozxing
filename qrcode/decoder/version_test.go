package decoder

import (
	"reflect"
	"testing"
)

func checkVersion(t testing.TB, v *Version, number, dimension int) {
	t.Helper()
	if v == nil {
		t.Fatalf("Version is nil, number=%d, dimension=%d", number, dimension)
	}
	num := v.GetVersionNumber()
	if num != number {
		t.Fatalf("version number is %v, expect %v", num, number)
	}
	if num > 1 {
		l := len(v.GetAlignmentPatternCenters())
		if l == 0 {
			t.Fatalf("AlignmentPatternCenters is empty")
		}
	}
	dim := v.GetDimensionForVersion()
	if dim != dimension {
		t.Fatalf("DimensionForVersion is %v, expect %v", dim, dimension)
	}

	if l := len(v.ecBlocks); l != 4 {
		t.Fatalf("ecBlocks must has 4 ECBlocks, %v", l)
	}

	if _, e := v.buildFunctionPattern(); e != nil {
		t.Fatalf("BuildFunctionPattern failed, (%d,%d), %v", number, dimension, e)
	}
}

func TestVersionForNumber(t *testing.T) {
	_, e := Version_GetVersionForNumber(0)
	if e == nil {
		t.Fatalf("GetVersionForNumber(0) must be error")
	}
	for i := 1; i <= 40; i++ {
		v, _ := Version_GetVersionForNumber(i)
		checkVersion(t, v, i, 4*i+17)
	}
}

func TestGetProvisionalVersionForDimension(t *testing.T) {
	for i := 1; i <= 40; i++ {
		v, e := Version_GetProvisionalVersionForDimension(4*i + 17)
		if e != nil {
			t.Fatalf("GetProvisiionalVersionForDimension(%v * i + 17) failed: %v", i, e)
		}
		if n := v.GetVersionNumber(); n != i {
			t.Fatalf("VersionNumber is %v, expect %v", n, i)
		}
	}
}

func doTestVersion(t testing.TB, exceptedVersion, mask int) {
	t.Helper()
	v, e := Version_decodeVersionInformation(mask)
	if e != nil {
		t.Fatalf("decodeVersionInformation(%v) failed: %v", mask, e)
	}
	if n := v.GetVersionNumber(); n != exceptedVersion {
		t.Fatalf("decodeVersionInformation(%v) version number is %v, expect %v", mask, n, exceptedVersion)
	}
}

func TestDecodeVersionInformation(t *testing.T) {
	doTestVersion(t, 7, 0x07C94)
	doTestVersion(t, 12, 0x0C762)
	doTestVersion(t, 17, 0x1145D)
	doTestVersion(t, 22, 0x168C9)
	doTestVersion(t, 27, 0x1B08E)
	doTestVersion(t, 32, 0x209D5)
}

func TestVersion1(t *testing.T) {
	v := NewVersion(1, []int{},
		ECBlocks{7, []ECB{{1, 19}}},
		ECBlocks{10, []ECB{{1, 16}}},
		ECBlocks{13, []ECB{{1, 13}}},
		ECBlocks{17, []ECB{{1, 9}}})

	if r := v.GetTotalCodewords(); r != 26 {
		t.Fatalf("totalCodewords = %v, expect 26", r)
	}

	ecbs := v.GetECBlocksForLevel(-1)
	if ecbs != nil {
		t.Fatalf("ECBlocksForLevel(-1) must return nil")
	}

	ecbs = v.GetECBlocksForLevel(ErrorCorrectionLevel_L)
	expect := &ECBlocks{7, []ECB{{1, 19}}}
	if !reflect.DeepEqual(ecbs, expect) {
		t.Fatalf("ECBlocksForLevel(L) is %v, expect %v", ecbs, expect)
	}
	ecbs = v.GetECBlocksForLevel(ErrorCorrectionLevel_M)
	expect = &ECBlocks{10, []ECB{{1, 16}}}
	if !reflect.DeepEqual(ecbs, expect) {
		t.Fatalf("ECBlocksForLevel(M) is %v, expect %v", ecbs, expect)
	}
	ecbs = v.GetECBlocksForLevel(ErrorCorrectionLevel_Q)
	expect = &ECBlocks{13, []ECB{{1, 13}}}
	if !reflect.DeepEqual(ecbs, expect) {
		t.Fatalf("ECBlocksForLevel(Q) is %v, expect %v", ecbs, expect)
	}
	ecbs = v.GetECBlocksForLevel(ErrorCorrectionLevel_H)
	expect = &ECBlocks{17, []ECB{{1, 9}}}
	if !reflect.DeepEqual(ecbs, expect) {
		t.Fatalf("ECBlocksForLevel(H) is %v, expect %v", ecbs, expect)
	}

	if r := ecbs.GetTotalECCodewords(); r != 17 {
		t.Fatalf("ECBlocks.GetNumBlocks is %v, expect 9", r)
	}
}

func TestVersion_GetProvisionalVersionForDimensionFail(t *testing.T) {
	_, e := Version_GetProvisionalVersionForDimension(3)
	if e == nil {
		t.Fatalf("GetProvisionalVersionForDimension(3) must be error")
	}

	_, e = Version_GetProvisionalVersionForDimension(181)
	if e == nil {
		t.Fatalf("GetProvisionalVersionForDimension(3) must be error")
	}
}

func TestVersion_BuildFunctionPatternFail(t *testing.T) {
	// dimension will be negative
	v := &Version{-5, []int{}, []ECBlocks{}, 0}
	_, e := v.buildFunctionPattern()
	if e == nil {
		t.Fatalf("BuildFunctionPatturn must be error")
	}
}

func TestDecodeVersionInformationUnmatch(t *testing.T) {
	v, e := Version_decodeVersionInformation(0x0C763) // best: 0x0C762, versionNumber=12
	if e != nil {
		t.Fatalf("decodeVersionInformation(0x0C763) failed :%v", e)
	}
	if n := v.GetVersionNumber(); n != 12 {
		t.Fatalf("best version number is %v, expect 12", n)
	}

	if _, e := Version_decodeVersionInformation(0); e == nil {
		t.Fatalf("decodeVersionInformation(0) must be error")
	}
}

func TestVersion_String(t *testing.T) {
	v, _ := Version_decodeVersionInformation(0x0C763) // best: 0x0C762, versionNumber=12
	if str := v.String(); str != "12" {
		t.Fatalf("String = \"%v\", expect \"12\"", str)
	}
	v = nil
	if str := v.String(); str != "" {
		t.Fatalf("String = \"%v\", expect \"\"", str)
	}
}
