package mod

import "strings"

// https://apidocs.vintagestory.at/api/Vintagestory.API.Common.EnumAppSide.html
type AppSide uint8

const (
	Server    AppSide = 1
	Client    AppSide = 2
	Universal AppSide = Server | Client
)

func (a *AppSide) MarshalJSON() ([]byte, error) {
	switch *a {
	case Server:
		return []byte("Server"), nil
	case Client:
		return []byte("Client"), nil
	default:
		return []byte("Universal"), nil
	}
}

func (a *AppSide) UnmarshalJSON(side []byte) error {
	switch strings.ToLower(strings.Trim(string(side), `"`)) {
	case "server", "1":
		*a = Server
	case "client", "2":
		*a = Client
	default:
		*a = Universal
	}
	return nil
}
