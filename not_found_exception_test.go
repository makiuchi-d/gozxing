package gozxing

import (
	"reflect"
	"testing"
)

func TestNewNotFoundException(t *testing.T) {
	var e error = NewNotFoundException("test")
	switch e.(type) {
	case NotFoundException:
		break
	default:
		t.Fatalf("Type is not NotFoundException, %v", reflect.TypeOf(e))
	}
}
