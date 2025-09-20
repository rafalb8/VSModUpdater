package modes

import (
	"fmt"
	"os"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func Update(interactive bool) {
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

		if interactive {
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
	if interactive {
		fmt.Print("Press any key ")
		fmt.Scanln()
	}
}
