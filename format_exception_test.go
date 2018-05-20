package gozxing

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewFormatException(t *testing.T) {
	var e error = FormatException_GetFormatInstance()
	switch e.(type) {
	case FormatException:
		break
	default:
		t.Fatalf("Type is not FormatException, %v", reflect.TypeOf(e))
	}

	base := errors.New("test")
	e = FormatException_GetFormatInstanceWithError(base)
	if _, ok := e.(FormatException); !ok {
		t.Fatalf("Type is not FormatException, %v", reflect.TypeOf(e))
	}
	if e.Error() != base.Error() {
		t.Fatalf("error is not inhelited, %v, expect %v", e, base)
	}
}
