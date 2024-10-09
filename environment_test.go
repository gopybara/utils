package utils

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Run("should return default value", func(t *testing.T) {
		data := GetEnv("FOO", "BAR")
		if data != "BAR" {
			t.Error("Expected BAR, got ", data)
		}
	})

	t.Run("should return int value", func(t *testing.T) {
		os.Setenv("INTVALUE", "1")

		data := GetEnv("INTVALUE", 2)
		if data != 1 {
			t.Error("Expected 1, got ", data)
		}
	})

	t.Run("should return bool value", func(t *testing.T) {
		os.Setenv("BOOLVALUE", "true")

		data := GetEnv("BOOLVALUE", true)
		if data != true {
			t.Error("Expected true, got ", data)
		}
	})

	t.Run("should return string value", func(t *testing.T) {
		os.Setenv("STRINGVALUE", "test")

		data := GetEnv("STRINGVALUE", "test")
		if data != "test" {
			t.Error("Expected test, got ", data)
		}
	})

	t.Run("should return slice value", func(t *testing.T) {
		os.Setenv("SLICEVALUE", "test1,test2")
		data := GetEnv("SLICEVALUE", []string{"test1", "test3"})
		if !ArrayEquals(data, []string{"test1", "test2"}) {
			t.Error("Expected test1, test2, got ", data)
		}
	})
}
