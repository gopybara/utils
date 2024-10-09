package utils

import (
	"reflect"
	"testing"
)

func TestArrayValueIndex(t *testing.T) {
	req := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	t.Run("should return offset of value", func(t *testing.T) {
		offset, exists := ArrayValueIndex(req, 11)
		if !exists {
			t.Error("should exists")
		}

		if offset != 10 {
			t.Errorf("expect 12, got %d", offset)
		}
	})

	t.Run("should return false if not found", func(t *testing.T) {
		offset, exists := ArrayValueIndex(req, 22)
		if exists {
			t.Error("should not exists")
		}

		if offset != 0 {
			t.Errorf("expect 0, got %d", offset)
		}
	})
}

func TestArrayKeysDelete(t *testing.T) {
	req := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	t.Run("should delete value", func(t *testing.T) {
		expect := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15}
		result := ArrayKeysDelete(req, 11, 16)
		if !reflect.DeepEqual(expect, result) {
			t.Errorf("expect %v, got %v", expect, result)
		}
	})
}

func TestArrayContains(t *testing.T) {
	req := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	t.Run("should return true if array contains", func(t *testing.T) {
		result := ArrayContains(req, 11)
		if ArrayContains(req, 11) != true {
			t.Errorf("expect %v, got %v", true, result)
		}
	})
}

func TestArrayEquals(t *testing.T) {
	req := []float64{1, 2, 3, 4, 5, 6, 7}

	t.Run("should return true if array compare", func(t *testing.T) {
		if !ArrayEquals(req, []float64{1, 2, 3, 4, 5, 6, 7}) {
			t.Errorf("expect arrays to be equal")
		}
	})
}

func BenchmarkArrayEquals(t *testing.B) {
	req := []float64{1, 2, 3, 4, 5, 6, 7}
	t.Run("should return true if array compare", func(t *testing.B) {
		// should be faster than reflect.DeepEqual
		result := ArrayEquals(req, []float64{1, 2, 3, 4, 5, 6, 7})
		if !result {
			t.Errorf("expect arrays to be equal")
		}
	})
}
