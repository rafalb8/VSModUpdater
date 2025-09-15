package mod

import (
	"strconv"
	"strings"
)

type Bool bool

func (b *Bool) MarshalJSON() ([]byte, error) {
	boolean := false
	if b != nil {
		boolean = (bool)(*b)
	}
	return []byte(strconv.FormatBool(boolean)), nil
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	boolean := strings.ToLower(strings.Trim(string(data), `"`))

	switch boolean {
	case "true", "t", "1":
		*b = true
	default:
		*b = false
	}

	return nil
}
