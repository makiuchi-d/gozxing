package rss

import (
	"testing"
)

func TestDataCharacter(t *testing.T) {
	d := NewDataCharacter(60, 5000)

	if r, wants := d.GetValue(), 60; r != wants {
		t.Fatalf("GetValue() = %v, wants %v", r, wants)
	}
	if r, wants := d.GetChecksumPortion(), 5000; r != wants {
		t.Fatalf("GetChecksumPortion() = %v, wants %v", r, wants)
	}
	if r, wants := d.String(), "60(5000)"; r != wants {
		t.Fatalf("String() = %v, wants %v", r, wants)
	}
	if r, wants := d.HashCode(), 60^5000; r != wants {
		t.Fatalf("HashCode() = %v, wants %v", r, wants)
	}

	tests := []struct {
		data  interface{}
		wants bool
	}{
		{NewDataCharacter(60, 5000), true},
		{NewDataCharacter(60, 19572), false},
		{NewDataCharacter(97250, 5000), false},
		{NewDataCharacter(97250, 19572), false},
		{struct{}{}, false},
	}
	for _, test := range tests {
		if r := d.Equals(test.data); r != test.wants {
			t.Fatalf("Equals(%v) = %v, wants %v", test.data, r, test.wants)
		}
	}
}
