package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		buf, result string
		err         error
	)

	for _, c := range s {
		if unicode.IsDigit(c) {
			if len(buf) > 0 {
				value, _ := strconv.Atoi(string(c))
				result += strings.Repeat(buf, value)
				// result += strings.Repeat(buf, int(c-'0'))
				buf = ""
			} else {
				err = ErrInvalidString
				return "", err
			}
		} else {
			if len(buf) > 0 {
				result += buf
			}
			buf = string(c)
		}
	}
	result += buf
	return result, err
}
