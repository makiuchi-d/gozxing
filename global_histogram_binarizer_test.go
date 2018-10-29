package gozxing

import (
	"testing"
)

func newTestLuminanceSource2(size int) LuminanceSource {
	return &testLuminanceSource{
		LuminanceSourceBase{size, size},
		func(x, y int) byte {
			if (y+x)%2 == 0 {
				return 10 + byte(50*x/size)
			}
			return 250 - byte(50*x/size)
		},
	}
}

func newTestBlackSource(size int) LuminanceSource {
	return &testLuminanceSource{
		LuminanceSourceBase{size, size},
		func(x, y int) byte { return 0 },
	}
}

func TestGlobalHistgramBinarizer(t *testing.T) {
	size := 32
	src := newTestLuminanceSource2(size)
	ghb := NewGlobalHistgramBinarizer(src)

	if s := ghb.GetLuminanceSource(); s != src {
		t.Fatalf("GetLuminanceSource = %p, expect %p", s, src)
	}
	if w, h := ghb.GetWidth(), ghb.GetHeight(); w != size || h != size {
		t.Fatalf("GetWidth,GetHeight = %v,%v, expect %v,%v", w, h, size, size)
	}
}

func TestGlobalHistgramBinarizer_estimateBlackPoint(t *testing.T) {
	g := GlobalHistogramBinarizer{}

	// single peak
	buckets := []int{0, 0, 0, 15, 12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	_, e := g.estimateBlackPoint(buckets)
	if _, ok := e.(NotFoundException); !ok {
		t.Fatalf("estimateBlackPoint must be NotFoundException, %T", e)
	}

	buckets = []int{0, 0, 0, 15, 12, 12, 5, 14, 16, 19, 20, 18, 0, 0, 0, 0}
	valley := 6 << LUMINANCE_SHIFT
	r, e := g.estimateBlackPoint(buckets)
	if e != nil {
		t.Fatalf("estimateBlackPoint returns error, %v", e)
	}
	if r != valley {
		t.Fatalf("estimateBlackPoint = %v, expect %v", r, valley)
	}
}

func TestGlobalHistgramBinarizer_GetBlackRow(t *testing.T) {
	src := newTestLuminanceSource2(16)
	ghb := NewGlobalHistgramBinarizer(src)

	if _, e := ghb.GetBlackRow(16, nil); e == nil {
		t.Fatalf("GetBlackRow must be error")
	}

	r, e := ghb.GetBlackRow(1, nil)
	if e != nil {
		t.Fatalf("GetBlackRow returns error, %v", e)
	}
	expect := " .X.X.X.X .X.X.X.."
	if r.String() != expect {
		t.Fatalf("GetBlackRow = \"%v\", expect \"%v\"", r, expect)
	}

	// white image
	ghb = ghb.CreateBinarizer(newTestBlackSource(16))
	_, e = ghb.GetBlackRow(0, r)
	if _, ok := e.(NotFoundException); !ok {
		t.Fatalf("GetBlackRow must be NotFoundException, %T", e)
	}

	// small image
	ghb = ghb.CreateBinarizer(newTestLuminanceSource2(2))
	r, e = ghb.GetBlackRow(0, nil)
	if e != nil {
		t.Fatalf("GetBlackRow returns error, %v", e)
	}
	expect = " X."
	if r.String() != expect {
		t.Fatalf("GetBlackRow = \"%v\", expect \"%v\"", r, expect)
	}
}

func TestGlobalHistgramBinarizer_GetBlackMatrix(t *testing.T) {
	ghb := NewGlobalHistgramBinarizer(newTestLuminanceSource2(0))
	_, e := ghb.GetBlackMatrix()
	if e == nil {
		t.Fatalf("GetBlackMatrix must be error")
	}

	ghb = NewGlobalHistgramBinarizer(newTestBlackSource(16))
	_, e = ghb.GetBlackMatrix()
	if _, ok := e.(NotFoundException); !ok {
		t.Fatalf("GetBlackMatrix must be NotFoundException, %T", e)
	}

	src := newTestLuminanceSource2(16)
	rawmatrix := src.GetMatrix()
	ghb = NewGlobalHistgramBinarizer(src)
	m, e := ghb.GetBlackMatrix()
	if e != nil {
		t.Fatalf("GetBlackMatrix returns error, %v", e)
	}
	for w := 0; w < m.GetWidth(); w++ {
		for h := 0; h < m.GetHeight(); h++ {
			expect := rawmatrix[w+m.GetHeight()*h] < 128
			if r := m.Get(w, h); r != expect {
				t.Fatalf("GetBlackMatrix [%v,%v] is %v, expect %v", w, h, r, expect)
			}
		}
	}
}
