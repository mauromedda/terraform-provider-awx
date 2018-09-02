package awx

import (
	"strconv"
)

// AtoipOr takes a string and a defaultValue. If the string cannot be converted, defaultValue is returned
func AtoipOr(s string, defaultValue *int) *int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return &n
}
