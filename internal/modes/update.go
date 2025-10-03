package modes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func Update() {
	if config.Interactive {
		defer func() {
			fmt.Print("Press any key ")
			fmt.Scanln()
		}()
	}

	fmt.Println("Updating mods:", config.ModPath)
	mods, err := mod.InfoFromPath(config.ModPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(mods) == 0 {
		fmt.Println("No Mods found")
		return
	}

	for _, m := range mods {
		if _, ignored := config.Ignored[m.ModID]; ignored {
			fmt.Println(m, "- Ignore")
			continue
		}

		if m.Error != nil {
			fmt.Print("\033[0;31m") // Red
			fmt.Println("!!!", filepath.Base(m.Path), "- Failed:", m.Error)
			fmt.Print("\033[0m") // Reset
			continue
		}

		update, err := m.CheckUpdates()
		switch {
		case err == mod.ErrNoUpdate:
			fmt.Println(m, "- Skip")
			continue
		case err == mod.ErrNoModID:
			fmt.Println(m, "- Missing ModID")
			continue
		case err != nil:
			fmt.Println(err)
			return
		}

		if config.DryRun {
			fmt.Printf("%s - Update v%s found!\n", m, update.Version)
			continue
		}

		if config.Interactive {
			shouldUpdate := ""
			fmt.Printf("Update %s: %s => %s? [Y/n] ", m.Name, m.Version, update.Version)
			fmt.Scanf("%s", &shouldUpdate)
			if len(shouldUpdate) > 0 && shouldUpdate[0]|' ' == 'n' {
				fmt.Println(m, "- SKIP")
				continue
			}
		}

		fmt.Printf("Downloading %s: %s => %s - ", m.Name, m.Version, update.Version)
		err = update.Download()
		if err != nil {
			fmt.Println("FAIL")
			fmt.Println(err)
			continue
		}
		fmt.Println("SUCCESS")

		if config.Backup {
			fmt.Printf("Backing up %s - ", m)
			err = m.Backup()
		} else {
			fmt.Printf("Removing %s - ", m)
			err = os.Remove(m.Path)
		}

		if err != nil {
			fmt.Println("FAIL")
			fmt.Println(err)
			continue
		}
		fmt.Println("SUCCESS")
	}

	fmt.Println("Finished Updating.")
}
