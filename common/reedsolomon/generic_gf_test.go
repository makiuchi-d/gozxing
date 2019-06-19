package reedsolomon

import (
	"testing"
)

func TestNewGenericGF(t *testing.T) {
	g := NewGenericGF(0x011D, 256, 0)

	if g.GetZero() == nil {
		t.Fatalf("GetZero() returns nil")
	}
	if g.GetOne() == nil {
		t.Fatalf("GetOne() returns nil")
	}
	if r := g.GetSize(); r != 256 {
		t.Fatalf("GetSize() = %v, expect 256", r)
	}
	if r := g.GetGeneratorBase(); r != 0 {
		t.Fatalf("GetGeneratorBase() = %v, expect 0", r)
	}
	if s := g.String(); s != "GF(0x11d,256)" {
		t.Fatalf("String() = %v, expect \"GF(0x11d,256)\"", s)
	}
}

func TestGenericGF_BuildMonomial(t *testing.T) {
	g := NewGenericGF(0x011D, 256, 0)

	if _, e := g.BuildMonomial(-1, 0); e == nil {
		t.Fatalf("BuildMonomial(-1, 0) must be error")
	}

	r, e := g.BuildMonomial(0, 0)
	if e != nil {
		t.Fatalf("BuildMonomial returns error, %v", e)
	}
	if !r.IsZero() {
		t.Fatalf("BuildMonomial returns non zero GenericGFPoly, %v", r)
	}

	degree := 2
	coefficient := 5
	r, e = g.BuildMonomial(degree, coefficient)
	if e != nil {
		t.Fatalf("BuildMonomial returns error, %v", e)
	}
	c := r.GetCoefficients()
	if len(c) <= degree {
		t.Fatalf("coefficients length not enough, %v", len(c))
	}
	if c[0] != coefficient {
		t.Fatalf("coefficients[0] = %v, expect %v", c[0], coefficient)
	}
}

func TestGenericGF_addOrSubtract(t *testing.T) {
	// 12 xor 10 = 6
	if r := GenericGF_addOrSubtract(12, 10); r != 6 {
		t.Fatalf("addOrSubtract(12, 10) = %v, expect 6", r)
	}
}

func testGenericGF_Exp(t testing.TB, g *GenericGF, a, expect int) {
	t.Helper()
	if r := g.Exp(a); r != expect {
		t.Fatalf("Exp(%v) = %v, expect %v", a, r, expect)
	}
}

func TestGenericGF_Exp(t *testing.T) {
	g := NewGenericGF(0x011D, 256, 0)
	testGenericGF_Exp(t, g, 0, 1)
	testGenericGF_Exp(t, g, 1, 2)
	testGenericGF_Exp(t, g, 7, 128)
	testGenericGF_Exp(t, g, 8, 29)
	testGenericGF_Exp(t, g, 128, 133)
	testGenericGF_Exp(t, g, 255, 1)
}

func testGenericGF_Log(t testing.TB, g *GenericGF, a, expect int) {
	t.Helper()
	r, e := g.Log(a)
	if e != nil {
		t.Fatalf("Log(%v) returns error, %v", a, e)
	}
	if r != expect {
		t.Fatalf("Log(%v) = %v, expect %v", a, r, expect)
	}
}

func TestGenericGF_Log(t *testing.T) {
	g := NewGenericGF(0x011D, 256, 0)
	if _, e := g.Log(0); e == nil {
		t.Fatalf("Log(0) must be error")
	}
	testGenericGF_Log(t, g, 1, 0)
	testGenericGF_Log(t, g, 10, 51)
	testGenericGF_Log(t, g, 90, 19)
	testGenericGF_Log(t, g, 128, 7)
	testGenericGF_Log(t, g, 170, 151)
	testGenericGF_Log(t, g, 255, 175)
}

func testGenericGF_Inverse(t testing.TB, g *GenericGF, a, expect int) {
	t.Helper()
	r, e := g.Inverse(a)
	if e != nil {
		t.Fatalf("Inverse(%v) returns error, %v", a, e)
	}
	if r != expect {
		t.Fatalf("Inverse(%v) = %v, expect %v", a, r, expect)
	}
}

func TestGenericGF_Inverse(t *testing.T) {
	g := NewGenericGF(0x011D, 256, 0)
	if _, e := g.Inverse(0); e == nil {
		t.Fatalf("Inverse(0) must be error")
	}
	testGenericGF_Inverse(t, g, 1, 1)
	testGenericGF_Inverse(t, g, 5, 167)
	testGenericGF_Inverse(t, g, 50, 111)
	testGenericGF_Inverse(t, g, 120, 219)
	testGenericGF_Inverse(t, g, 255, 253)
}

func testGenericGF_Multiply(t testing.TB, g *GenericGF, a, b, expect int) {
	t.Helper()
	r := g.Multiply(a, b)
	if r != expect {
		t.Fatalf("Multiply(%v, %v) = %v, expect %v", a, b, r, expect)
	}
}

func TestGenericGF_Multiply(t *testing.T) {
	g := NewGenericGF(0x011D, 256, 0)
	testGenericGF_Multiply(t, g, 0, 1, 0)
	testGenericGF_Multiply(t, g, 1, 0, 0)
	testGenericGF_Multiply(t, g, 1, 1, 1)
	testGenericGF_Multiply(t, g, 5, 3, 15)
	testGenericGF_Multiply(t, g, 20, 128, 210)
	testGenericGF_Multiply(t, g, 120, 50, 152)
	testGenericGF_Multiply(t, g, 255, 255, 226)
}
