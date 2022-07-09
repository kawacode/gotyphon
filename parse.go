package gotyphon

import (
	"strings"
)

// It takes a string, and two other strings, and returns the string between the two other strings.
func Parse(str string, left string, right string) string {
	return strings.Split(strings.Split(str, left)[1], right)[0]
}
