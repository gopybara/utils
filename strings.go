package utils

import (
	"bytes"
	"strings"
)

func FirstLetterToLower(s string) string {
	bts := []byte(s)
	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}

func ToCamelCase(str string) string {
	str = strings.ReplaceAll(str, "_", " ")
	str = strings.ReplaceAll(str, "-", " ")
	str = strings.Title(str)
	str = strings.ReplaceAll(str, " ", "")

	return FirstLetterToLower(str)
}
