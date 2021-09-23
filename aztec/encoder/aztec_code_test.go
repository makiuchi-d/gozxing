package encoder

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestAztecCode(t *testing.T) {
	c := newAztecCode()

	c.setCompact(true)
	if r, wants := c.isCompact(), true; r != wants {
		t.Fatalf("isCompact = %v, wants %v", r, wants)
	}

	c.setSize(15)
	if r, wants := c.getSize(), 15; r != wants {
		t.Fatalf("getSize= %v, wants %v", r, wants)
	}

	c.setLayers(2)
	if r, wants := c.getLayers(), 2; r != wants {
		t.Fatalf("getLayers= %v, wants %v", r, wants)
	}

	c.setCodeWords(10)
	if r, wants := c.getCodeWords(), 10; r != wants {
		t.Fatalf("getCodeWords= %v, wants %v", r, wants)
	}

	bm, _ := gozxing.NewBitMatrix(1, 1)
	c.setMatrix(bm)
	if r := c.getMatrix(); r != bm {
		t.Fatalf("getMatrix = %p, wants %p", r, bm)
	}
}
