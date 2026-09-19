package modes

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rafalb8/VSModUpdater/v2/internal/config"
	"github.com/rafalb8/VSModUpdater/v2/internal/mod"
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

	// Cache AssetID
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 10)
	for _, m := range mods {
		wg.Go(func() {
			sem <- struct{}{}
			m.FetchMod()
			<-sem
		})
	}
	wg.Wait()

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
