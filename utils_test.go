package utils

import "testing"

func TestPtr(t *testing.T) {
	t.Run("should return pointer", func(t *testing.T) {
		testStr := "test"

		v := Ptr(testStr)
		if *v != testStr {
			t.Errorf("got %v, want %v", *v, testStr)
		}
	})
}
