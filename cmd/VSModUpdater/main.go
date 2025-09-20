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

	case config.Interactive && len(config.Ignored) == 0:
		modes.Update(true)

	default:
		modes.Update(false)
	}
}
