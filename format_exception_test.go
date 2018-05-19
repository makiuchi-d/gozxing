package gozxing

import (
	"reflect"
	"testing"
)

func TestNewFormatException(t *testing.T) {
	var e error = NewFormatException("test")
	if _, ok := e.(FormatException); !ok {
		t.Fatalf("Type is not FormatException, %v", reflect.TypeOf(e))
	}
}
