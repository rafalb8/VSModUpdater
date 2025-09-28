package modes

import (
	"fmt"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func List() {
	mods, err := mod.InfoFromPath(config.ModPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(mods) == 0 {
		fmt.Println("No Mods found")
		return
	}

	sep := strings.Repeat("=", 80)
	for _, m := range mods {
		fmt.Println(sep)

		if m.Error != nil {
			fmt.Print("\033[0;31m") // Red
			fmt.Println(m.Details())
			fmt.Print("\033[0m") // Reset
			continue
		}

		fmt.Println(m.Details())
	}
	fmt.Println(sep)
}
