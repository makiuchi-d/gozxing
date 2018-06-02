package reedsolomon

import (
	"errors"
)

type GenericGFPoly struct {
	field        *GenericGF
	coefficients []int
}

func NewGenericGFPoly(field *GenericGF, coefficients []int) (*GenericGFPoly, error) {
	if len(coefficients) == 0 {
		return nil, errors.New("IllegalArgumentException")
	}
	this := &GenericGFPoly{field: field}

	coefficientsLength := len(coefficients)
	if coefficientsLength > 1 && coefficients[0] == 0 {
		// Leading term must be non-zero for anything except the constant polynomial "0"
		firstNonZero := 1
		for firstNonZero < coefficientsLength && coefficients[firstNonZero] == 0 {
			firstNonZero++
		}
		if firstNonZero == coefficientsLength {
			this.coefficients = []int{0}
		} else {
			this.coefficients = coefficients[firstNonZero:]
		}
	} else {
		this.coefficients = coefficients
	}

	return this, nil
}

func (this *GenericGFPoly) GetCoefficients() []int {
	return this.coefficients
}

func (this *GenericGFPoly) Getdegree() int {
	return len(this.coefficients) - 1
}

func (this *GenericGFPoly) IsZero() bool {
	return this.coefficients[0] == 0
}

func (this *GenericGFPoly) GetCoefficient(degree int) int {
	return this.coefficients[len(this.coefficients)-1-degree]
}
