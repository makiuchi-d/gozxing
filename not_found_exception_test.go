package gozxing

import (
	"reflect"
	"testing"
)

func TestNotFoundException_GetNotFoundInstance(t *testing.T) {
	var e error = NotFoundException_GetNotFoundInstance()
	switch e.(type) {
	case NotFoundException:
		break
	default:
		t.Fatalf("Type is not NotFoundException, %v", reflect.TypeOf(e))
	}
}
