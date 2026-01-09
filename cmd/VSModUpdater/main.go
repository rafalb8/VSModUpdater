package main

import (
	"flag"
	"fmt"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/modes"
)

func main() {
	flag.Parse()

	switch {
	case config.Version:
		fmt.Println(config.VersionNum)

	case config.Self:
		modes.Self()

	case config.List:
		modes.List()

	case config.Webpage:
		modes.Webpage()

	case config.WebConfigMode || config.ConfigModID != "" || config.ConfigAddCat != "" || config.ConfigEditCat != "" || config.ConfigDelCat != "":
		if err := modes.ConfigManagement(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}

	default:
		modes.Update()
	}
}
