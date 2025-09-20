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
		err := modes.Self()
		if err != nil {
			fmt.Println(err)
		}

	case config.List:
		modes.List()

	default:
		modes.Update()
	}
}
