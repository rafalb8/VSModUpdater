package mod

import (
	"fmt"
	"strings"
)

// https://apidocs.vintagestory.at/api/Vintagestory.API.Common.EnumModType.html
type Type uint8

const (
	Theme Type = iota
	Content
	Code
)

func (t *Type) MarshalJSON() ([]byte, error) {
	switch *t {
	case Theme:
		return []byte("theme"), nil
	case Content:
		return []byte("content"), nil
	case Code:
		return []byte("code"), nil
	default:
		return nil, fmt.Errorf("mod.Type: unknown type: %v", t)
	}
}

func (t *Type) UnmarshalJSON(typ []byte) error {
	typStr := strings.Trim(string(typ), `"`)
	switch strings.ToLower(typStr) {
	case "theme", "0":
		*t = Theme
	case "content", "1":
		*t = Content
	case "code", "2":
		*t = Code
	default:
		return fmt.Errorf("mod.Type: unknown value: %s", typStr)
	}
	return nil
}
