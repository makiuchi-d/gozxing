package common

import (
	"reflect"
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestDefaultGridSampler_SampleGrid(t *testing.T) {
	var e error
	s := NewDefaultGridSampler()
	image, _ := gozxing.NewBitMatrix(300, 300)

	_, e = s.SampleGrid(image, 0, 10,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0)
	if e == nil {
		t.Fatalf("SampleGrid must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("error must be NotFoundException, %v", reflect.TypeOf(e))
	}

	_, e = s.SampleGrid(image, 10, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0)
	if e == nil {
		t.Fatalf("SampleGrid must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("error must be NotFoundException, %v", reflect.TypeOf(e))
	}

	_, e = s.SampleGrid(image, 3, 3,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0)
	if e == nil {
		t.Fatalf("SampleGrid must be error")
	}
	if _, ok := e.(gozxing.NotFoundException); !ok {
		t.Fatalf("error must be NotFoundException, %v", reflect.TypeOf(e))
	}

	image.SetRegion(10, 10, 20, 30)
	bits, e := s.SampleGrid(image, 30, 30,
		1, 0, 0, 1, 0, 0, 1, 1,
		1, 0, 0, 1, 0, 0, 1, 1)
	if e != nil {
		t.Fatalf("SampleGrid returns error, %v", e)
	}
	if !bits.Get(25, 25) {
		t.Fatalf("sampling failed: bits[25,25] = false")
	}
}
