package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
	"github.com/rafalb8/VSModUpdater/internal/self"
)

func main() {
	flag.Parse()
	if config.Version {
		fmt.Println(config.VersionNum)
		return
	}

	if config.Self {
		err := self.Update()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	mods, err := mod.InfoFromPath(config.ModPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(mods) == 0 {
		fmt.Println("No Mods found")
		return
	}

	if config.List {
		List(mods)
	} else {
		Update(mods)
	}
}

func List(mods []*mod.Info) {
	sep := strings.Repeat("=", 80)
	for _, m := range mods {
		fmt.Println(sep)
		fmt.Println("Name:\t\t", m.Name)
		fmt.Println("Version:\t", m.Version)
		gameVer, ok := m.Dependencies["game"]
		if ok {
			if gameVer == "*" {
				gameVer = "any"
			}
			fmt.Println("Game Version:\t", gameVer)
		}
		fmt.Println("Authors:\t", strings.Join(m.Authors, ", "))
		fmt.Println("Description:\t", m.Description)
	}
	fmt.Println(sep)
}

func Update(mods []*mod.Info) {
	for _, m := range mods {
		update, err := m.CheckUpdates()
		switch {
		case err == mod.ErrNoUpdate:
			fmt.Println(m, "- SKIP")
			continue
		case err == mod.ErrNoModID:
			fmt.Println(m, "- Missing ModID")
			continue
		case err != nil:
			fmt.Println(err)
			return
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

	fmt.Println("DONE")
	time.Sleep(3 * time.Second)
}
