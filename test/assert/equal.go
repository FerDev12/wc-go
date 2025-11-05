package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal(t *testing.T, wants, got any, msg ...string) {
	t.Helper()

	if reflect.DeepEqual(wants, got) {
		return
	}

	t.Errorf("%s\nexpected: %v\n     got: %v\n", strings.Join(msg, ": "), wants, got)
}
