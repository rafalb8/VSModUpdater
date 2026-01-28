package main

import (
	"fmt"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/modes"
)

func main() {
	switch {
	case config.Version:
		fmt.Println(config.VersionNum)

	case config.Self:
		modes.Self()

	case config.List:
		modes.List()
		
	case config.Import != "":
		modes.Import(config.Import)
		
	case config.Export != "":
		modes.Export(config.Export)

	default:
		modes.Update()
	}
}
