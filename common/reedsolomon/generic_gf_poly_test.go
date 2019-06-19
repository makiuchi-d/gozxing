package reedsolomon

import (
	"reflect"
	"testing"
)

func TestNewGenericGFPoly(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256

	if _, e := NewGenericGFPoly(gf, []int{}); e == nil {
		t.Fatalf("NewGenericGFPoly must be error")
	}

	coefficients := []int{0, 0, 0}
	p, e := NewGenericGFPoly(gf, coefficients)
	if e != nil {
		t.Fatalf("NewGenericGFPoly returns error, %v", e)
	}
	if len(p.coefficients) != 1 || p.coefficients[0] != 0 {
		t.Fatalf("coefficients is %v, expect %v", p.coefficients, coefficients)
	}
	if !p.IsZero() {
		t.Fatalf("IsZero() must be true")
	}

	coefficients = []int{0, 1, 2}
	p, e = NewGenericGFPoly(gf, coefficients)
	if e != nil {
		t.Fatalf("NewGenericGFPoly returns error, %v", e)
	}
	if !reflect.DeepEqual(p.GetCoefficients(), coefficients[1:]) {
		t.Fatalf("coefficients is %v, expect %v", p.GetCoefficients(), coefficients[1:])
	}

	coefficients = []int{1, 2, 3}
	p, e = NewGenericGFPoly(gf, coefficients)
	if e != nil {
		t.Fatalf("NewGenericGFPoly returns error, %v", e)
	}
	if !reflect.DeepEqual(p.GetCoefficients(), coefficients) {
		t.Fatalf("coefficients is %v, expect %v", p.GetCoefficients(), coefficients)
	}

	if r := p.GetDegree(); r != 2 {
		t.Fatalf("GetDegree = %v, expect 2", r)
	}
	if r := p.GetCoefficient(0); r != 3 {
		t.Fatalf("GetCoefficient(0) = %v, expect 3", r)
	}
	if r := p.GetCoefficient(1); r != 2 {
		t.Fatalf("GetCoefficient(1) = %v, expect 2", r)
	}
	if r := p.GetCoefficient(2); r != 1 {
		t.Fatalf("GetCoefficient(2) = %v, expect 1", r)
	}
}

func testGenericGFPoly_EvaluateAt(t testing.TB, p *GenericGFPoly, a, expect int) {
	t.Helper()
	if r := p.EvaluateAt(a); r != expect {
		t.Fatalf("EvaluateAt(%v) = %v, expect %v", a, r, expect)
	}
}

func TestGenericGFPoly_EvaluateAt(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256
	coefficients := []int{2, 3, 5}
	p, _ := NewGenericGFPoly(gf, coefficients)

	testGenericGFPoly_EvaluateAt(t, p, 0, 5)
	testGenericGFPoly_EvaluateAt(t, p, 1, 4)
	testGenericGFPoly_EvaluateAt(t, p, 2, 11)
	testGenericGFPoly_EvaluateAt(t, p, 3, 10)
	testGenericGFPoly_EvaluateAt(t, p, 120, 88)
	testGenericGFPoly_EvaluateAt(t, p, 255, 192)
}

func testGenericGFPoly_AddOrSubtract(t testing.TB, p, other *GenericGFPoly, expect []int) {
	t.Helper()
	r, e := p.AddOrSubtract(other)
	if e != nil {
		t.Fatalf("AddOrSubtract returns error, %v", e)
	}
	if r.field != p.field {
		t.Fatalf("AddOrSubtract returns different GF, %v, expect %v", r.field, p.field)
	}
	if !reflect.DeepEqual(r.coefficients, expect) {
		t.Fatalf("AddOrSubtract coefficients = '%v', expect '%v'", r.coefficients, expect)
	}
}

func TestGenericGFPoly_AddOrSubtract(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256

	p, _ := NewGenericGFPoly(gf, []int{0})
	other, _ := NewGenericGFPoly(GenericGF_AZTEC_PARAM, []int{1, 2, 3})
	if _, e := p.AddOrSubtract(other); e == nil {
		t.Fatalf("AddOrSubtract must be error")
	}

	other, _ = NewGenericGFPoly(gf, []int{1, 2, 3})
	testGenericGFPoly_AddOrSubtract(t, p, other, other.coefficients)

	p, _ = NewGenericGFPoly(gf, []int{1, 2, 3})
	other, _ = NewGenericGFPoly(gf, []int{0})
	testGenericGFPoly_AddOrSubtract(t, p, other, p.coefficients)

	other, _ = NewGenericGFPoly(gf, []int{3, 5, 7})
	testGenericGFPoly_AddOrSubtract(t, p, other, []int{2, 7, 4})

	other, _ = NewGenericGFPoly(gf, []int{3, 5})
	testGenericGFPoly_AddOrSubtract(t, p, other, []int{1, 1, 6})
}

func testGenericGFPoly_Multiply(t testing.TB, p, o *GenericGFPoly, expect []int) {
	t.Helper()
	r, e := p.Multiply(o)
	if e != nil {
		t.Fatalf("Multiply(%v) returns error, %v", o, e)
	}
	if r.field != p.field {
		t.Fatalf("Multiply(%v) returns different GF, %v, expect %v", o, r.field, p.field)
	}
	if !reflect.DeepEqual(r.coefficients, expect) {
		t.Fatalf("Multiply(%v) coefficients = '%v', expect '%v'", o, r.coefficients, expect)
	}
}

func TestGenericGFPoly_Multiply(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256

	p, _ := NewGenericGFPoly(gf, []int{0})
	o, _ := NewGenericGFPoly(GenericGF_AZTEC_PARAM, []int{1, 2, 3})
	if _, e := p.Multiply(o); e == nil {
		t.Fatalf("Multiply(%v) must be error", o)
	}

	o, _ = NewGenericGFPoly(gf, []int{1, 2, 3})
	testGenericGFPoly_Multiply(t, p, o, []int{0})

	p, _ = NewGenericGFPoly(gf, []int{1, 2, 3})
	o, _ = NewGenericGFPoly(gf, []int{0})
	testGenericGFPoly_Multiply(t, p, o, []int{0})

	p, _ = NewGenericGFPoly(gf, []int{1, 2, 3})
	o, _ = NewGenericGFPoly(gf, []int{7, 11})
	testGenericGFPoly_Multiply(t, p, o, []int{7, 5, 31, 29})
}

func testGenericGFPoly_MultiplyBy(t testing.TB, p *GenericGFPoly, s int, expect []int) {
	t.Helper()
	r := p.MultiplyBy(s)
	if r.field != p.field {
		t.Fatalf("MultiplyBy returns different GF, %v, expect %v", r.field, p.field)
	}
	if !reflect.DeepEqual(r.coefficients, expect) {
		t.Fatalf("MultiplyBy(%v) coefficients = '%v', expect '%v'", s, r.coefficients, expect)
	}
}

func TestGenericGFPoly_MultiplyBy(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256
	p, _ := NewGenericGFPoly(gf, []int{2, 3, 5, 7})

	testGenericGFPoly_MultiplyBy(t, p, 0, []int{0})
	testGenericGFPoly_MultiplyBy(t, p, 1, []int{2, 3, 5, 7})
	testGenericGFPoly_MultiplyBy(t, p, 3, []int{6, 5, 15, 9})
}

func testGenericGFPoly_MultiplyMonomial(t testing.TB, p *GenericGFPoly, d, c int, expect []int) {
	t.Helper()
	r, e := p.MultiplyByMonomial(d, c)
	if e != nil {
		t.Fatalf("MultiplyByMonomial(%v,%v) returns error, %v", d, c, e)
	}
	if r.field != p.field {
		t.Fatalf("MultiplyByMonomial(%v,%v) returns different GF, %v, expect %v", d, c, r.field, p.field)
	}
	if !reflect.DeepEqual(r.coefficients, expect) {
		t.Fatalf("MultiplyByMonomial(%v,%v) coefficients = '%v', expect '%v'", d, c, r.coefficients, expect)
	}
}

func TestGenericGFPoly_MultiplyByMonomial(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256
	p, _ := NewGenericGFPoly(gf, []int{3, 5, 7})

	if _, e := p.MultiplyByMonomial(-1, 3); e == nil {
		t.Fatalf("MultiplyByMonomial(-1, 3) must be error")
	}
	testGenericGFPoly_MultiplyMonomial(t, p, 1, 0, []int{0})
	testGenericGFPoly_MultiplyMonomial(t, p, 0, 1, []int{3, 5, 7})
	testGenericGFPoly_MultiplyMonomial(t, p, 0, 3, []int{5, 15, 9})
	testGenericGFPoly_MultiplyMonomial(t, p, 1, 1, []int{3, 5, 7, 0})
	testGenericGFPoly_MultiplyMonomial(t, p, 1, 3, []int{5, 15, 9, 0})
	testGenericGFPoly_MultiplyMonomial(t, p, 3, 5, []int{15, 17, 27, 0, 0, 0})
}

func testGenericGFPoly_Divide(t testing.TB, p, o *GenericGFPoly, expquot, exprem []int) {
	t.Helper()
	quot, rem, e := p.Divide(o)
	if e != nil {
		t.Fatalf("Divide(%v) returns error, %v", o, e)
	}
	if quot == nil || rem == nil {
		t.Fatalf("Divide(%v) returns nil, %v, %v", o, quot, rem)
	}
	if quot.field != p.field {
		t.Fatalf("Divide(%v) quotient has different GF, %v", o, quot.field)
	}
	if rem.field != p.field {
		t.Fatalf("Divide(%v) remainder has different GF, %v", o, rem.field)
	}
	if !reflect.DeepEqual(quot.coefficients, expquot) {
		t.Fatalf("Divide(%v) quotient coefficients = '%v', expect '%v'", o, quot.coefficients, expquot)
	}
	if !reflect.DeepEqual(rem.coefficients, exprem) {
		t.Fatalf("Divide(%v) reminder coefficients = '%v', expect '%v'", o, rem.coefficients, exprem)
	}
}

func TestGenericGFPoly_Divide(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256
	p, _ := NewGenericGFPoly(gf, []int{11, 7, 5, 0, 3})
	o, _ := NewGenericGFPoly(GenericGF_AZTEC_PARAM, []int{1, 2, 3})

	if _, _, e := p.Divide(o); e == nil {
		t.Fatalf("Divide by different field must be error")
	}

	o, _ = NewGenericGFPoly(gf, []int{0})
	if _, _, e := p.Divide(o); e == nil {
		t.Fatalf("Divide by 0 must be error")
	}

	testGenericGFPoly_Divide(t, p, p, []int{1}, []int{0})

	o, _ = NewGenericGFPoly(gf, []int{1, 0, 0, 0, 0, 0})
	testGenericGFPoly_Divide(t, p, o, []int{0}, p.coefficients)

	o, _ = NewGenericGFPoly(gf, []int{3, 2, 1})
	testGenericGFPoly_Divide(t, p, o, []int{242, 161, 147}, []int{154, 144})
}

func TestGenericGFPoly_String(t *testing.T) {
	gf := GenericGF_QR_CODE_FIELD_256

	if s := gf.GetZero().String(); s != "0" {
		t.Fatalf("string of zero must be \"0\", %s", s)
	}

	p, _ := NewGenericGFPoly(gf, []int{3, 0, -2, 1, 1})
	expect := "a^25x^4 - ax^2 + x + 1"
	if s := p.String(); s != expect {
		t.Fatalf("string is %s, expect %s", s, expect)
	}

	p, _ = NewGenericGFPoly(gf, []int{-1})
	expect = "-1"
	if s := p.String(); s != expect {
		t.Fatalf("string is %s, expect %s", s, expect)
	}

	p, _ = NewGenericGFPoly(gf, []int{3})
	expect = "a^25"
	if s := p.String(); s != expect {
		t.Fatalf("string is %s, expect %s", s, expect)
	}
}
