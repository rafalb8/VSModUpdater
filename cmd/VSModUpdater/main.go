package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func main() {
	flag.Parse()

	mods := []*mod.Info{}
	filepath.WalkDir(config.ModPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path == config.ModPath {
				return nil
			}
			return fs.SkipDir
		}

		name := d.Name()
		if filepath.Ext(name) != ".zip" {
			return nil
		}

		info, err := mod.InfoFromZip(path)
		if err != nil {
			return err
		}
		mods = append(mods, info)
		return nil
	})

	if len(mods) == 0 {
		fmt.Println("No Mods found")
		return
	}

	for _, m := range mods {
		update, err := m.CheckUpdates()
		if err == mod.ErrNoUpdate {
			fmt.Println(m, "- SKIP")
			continue
		}

		if err != nil {
			panic(err)
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
}
