package helpers

import (
	"strconv"
	"strings"
)

func ParseIDFromPath(path, prefix string) (int64, bool) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	s := strings.TrimPrefix(path, prefix)
	s = strings.TrimSuffix(s, "/")
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
