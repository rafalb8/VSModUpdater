package modes

import (
	"bufio"
	"fmt"
	"iter"
	"os"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/filter"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

type update struct {
	*mod.Info
	Update mod.Update
}

func Update() {
	mods, err := mod.InfoFromPath(config.ModPath)
	if err != nil {
		fmt.Println("Error loading mods:", err)
		return
	}

	if len(mods) == 0 {
		fmt.Println("No Mods found")
		return
	}

	fmt.Println(":: Searching for updates...")

	var (
		updates     = make([]update, 0, len(mods))
		preReleases = []update{} // pre-release mod version
		unstable    = []update{} // pre-release game version
		errors      = map[string]error{}
		upToDate    = 0
	)

	for _, m := range mods {
		if _, ignored := config.Ignored[m.ModID]; ignored {
			fmt.Printf(" %s - Ignored\n", m.String())
			continue
		}

		if m.Error != nil {
			errors[m.Name] = m.Error
			continue
		}

		u, err := m.CheckUpdates()
		upd := update{m, u}

		switch err {
		case nil:
			updates = append(updates, upd)

		case mod.ErrNoUpdate:
			upToDate += 1

		case mod.ErrPreReleaseSkip:
			preReleases = append(preReleases, upd)

		case mod.ErrUnstableSkip:
			unstable = append(unstable, upd)

		default:
			errors[m.Name] = err
		}
	}

	if len(errors) > 0 {
		fmt.Println(":: Errors encountered during check:")
		for name, err := range errors {
			fmt.Printf(" %s: %v\n", name, err)
		}
	}

	if len(preReleases) > 0 {
		fmt.Println(":: Pre-release updates skipped:")
		for _, m := range preReleases {
			fmt.Printf(" %s (%s -> %s)\n", m.Name, m.Version, m.Update.Version)
		}
	}

	if len(unstable) > 0 {
		fmt.Println(":: Unstable updates skipped:")
		for _, m := range unstable {
			fmt.Printf(" %s (%s -> %s)\n", m.Name, m.Version, m.Update.Version)
		}
	}

	fmt.Printf(":: %d updates available (%d are up to date).\n\n", len(updates), upToDate)
	if len(updates) == 0 {
		return
	}

	for i, m := range updates {
		fmt.Printf("[%d] %s (%s -> %s) - %s\n", i+1, m.Name, m.Version, m.Update.Version, m.Page())
	}

	s := bufio.NewScanner(os.Stdin)
	fmt.Println("\n=> Mods to EXCLUDE from update: (e.g. 1 2 3, 1-3, ^4)")
	fmt.Print("=> ")
	s.Scan()
	fmt.Println()

	filter, err := filter.NewExclusion[update](s.Text())
	if err != nil {
		fmt.Println("Invalid exclude expression:", err)
		return
	}

	fmt.Println(":: Updating mods...")
	for m := range filter.Filter(OneBased(updates)) {
		fmt.Printf(" %s@%s", m.Name, m.Update.Version)

		if config.DryRun {
			fmt.Println(" - OK")
			continue
		}

		// Backup before download. New file might have the same filename
		err = m.Backup()
		if err != nil {
			fmt.Println("-", err)
			continue
		}

		err = m.Update.Download()
		if err != nil {
			fmt.Println("-", err)

			// Try to restore the backup
			err = m.Restore()
			if err != nil {
				fmt.Println("Restore:", err)
			}
			continue
		}

		fmt.Println("- OK")

		if config.Backup {
			continue
		}

		// Remove the backup
		err = os.RemoveAll(m.Path)
		if err != nil {
			fmt.Println("Cleanup:", err)
			continue
		}
	}
}

func OneBased[T any](s []T) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range s {
			if !yield(i+1, v) {
				return
			}
		}
	}
}
