package gozxing

import (
	"errors"
	"testing"
)

func TestFormatException(t *testing.T) {
	var e error = GetFormatExceptionInstance()

	if _, ok := e.(FormatException); !ok {
		t.Fatalf("Type is not FormatException")
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Type is not ReaderException")
	}
	if _, ok := e.(NotFoundException); ok {
		t.Fatalf("Type must not be NotFoundException")
	}

	e.(FormatException).FormatException()
	e.(FormatException).ReaderException()
}

func TestNewFormatException(t *testing.T) {
	base := errors.New("newformatexception")
	var e error = NewFormatExceptionInstance(base)

	if _, ok := e.(FormatException); !ok {
		t.Fatalf("Type is not FormatException")
	}
	if _, ok := e.(ReaderException); !ok {
		t.Fatalf("Type is not ReaderException")
	}
	if _, ok := e.(NotFoundException); ok {
		t.Fatalf("Type must not be NotFoundException")
	}
	if e.Error() != base.Error() {
		t.Fatalf("e.Error() = %v, expect %v", e.Error(), base.Error())
	}
}
