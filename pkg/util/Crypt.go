package util

import (
	"strconv"

	"github.com/segmentio/fasthash/fnv1a"
)

func HashHex(source string) string {
	return strconv.FormatUint(fnv1a.HashString64(source), 16)
}
