package utils

import "testing"

func TestToCamelCase(t *testing.T) {
	t.Run("should convert to CamelCase", func(t *testing.T) {
		cases := map[string]string{
			"HelloWorld":  "helloWorld",
			"hello World": "helloWorld",
			"hello world": "helloWorld",
		}

		for request, expect := range cases {
			got := ToCamelCase(request)
			if got != expect {
				t.Errorf("ToCamelCase(%s) got %s, want %s", request, got, expect)
			}
		}
	})
}

func TestFirstLetterToLower(t *testing.T) {
	t.Run("should make first letter to lower case", func(t *testing.T) {
		cases := map[string]string{
			"HelloWorld": "helloWorld",
			"Hi Mom":     "hi Mom",
		}

		for request, expect := range cases {
			got := FirstLetterToLower(request)
			if got != expect {
				t.Errorf("FirstLetterToLower(%s) got %s, want %s", request, got, expect)
			}
		}
	})
}
