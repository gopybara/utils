package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestMapClone(t *testing.T) {
	t.Run("should clone map", func(t *testing.T) {
		m := map[string]interface{}{
			"foo": "bar",
			"bar": "baz",
		}
		cloned := MapClone(m)
		if !reflect.DeepEqual(m, cloned) {
			t.Errorf("cloned map did not clone properly")
		}
	})
}

func TestMapKeys(t *testing.T) {
	t.Run("should get keys", func(t *testing.T) {
		expect := []string{"bar", "foo"}
		m := map[string]interface{}{
			"foo": "bar",
			"bar": "baz",
		}

		keys := MapKeys(m)

		sort.Strings(keys)

		if !reflect.DeepEqual(keys, expect) {
			t.Errorf("keys did not get properly")
		}
	})
}
