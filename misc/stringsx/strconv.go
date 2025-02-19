package stringsx

import (
	"strconv"
	"strings"
)

func Atoi(s string) (int, error) {
	if strings.HasPrefix(s, "0b") {
		i, err := strconv.ParseInt(s[2:], 2, 64)
		return int(i), err
	}
	if strings.HasPrefix(s, "0x") {
		i, err := strconv.ParseInt(s[2:], 16, 64)
		return int(i), err
	}
	if strings.HasPrefix(s, "0") {
		i, err := strconv.ParseInt(s[1:], 8, 64)
		return int(i), err
	}
	return strconv.Atoi(s)
}
