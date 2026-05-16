package mod

import (
	"strconv"
	"strings"
	"unsafe"
)

type Bool bool

func (b Bool) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatBool(bool(b))), nil
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	s := unsafe.String(unsafe.SliceData(data), len(data))
	boolean, err := strconv.ParseBool(strings.Trim(s, `"`))
	if err != nil {
		return err
	}
	*b = Bool(boolean)
	return nil
}
