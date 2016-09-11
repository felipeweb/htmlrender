package htmlrender

import (
	"testing"
	"reflect"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected ||%#v|| (type %v) - Got ||%#v|| (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func expectNil(t *testing.T, a interface{}) {
	if a != nil {
		t.Errorf("Expected ||nil|| - Got ||%#v|| (type %v)", a, reflect.TypeOf(a))
	}
}

func expectNotNil(t *testing.T, a interface{}) {
	if a == nil {
		t.Errorf("Expected ||not nil|| - Got ||nil|| (type %v)", reflect.TypeOf(a))
	}
}