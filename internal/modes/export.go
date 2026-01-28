package modes

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func Export(output string) {
	if config.Interactive {
		defer func() {
			fmt.Print("Press any key ")
			fmt.Scanln()
		}()
	}

	err := os.MkdirAll(filepath.Dir(output), 0o755)
	if err != nil {
		fmt.Println(err)
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

	modlist := make([]string, 0, len(mods))

	for _, m := range mods {
		if m.ModID == "" {
			fmt.Println(m, "-", mod.ErrNoModID)
			continue
		}
		modlist = append(modlist, m.ModID+"@"+m.Version.String())
	}

	err = os.WriteFile(output, []byte(strings.Join(modlist, "\n")), 0o644)
	if err != nil {
		fmt.Println(err)
		return
	}

	abspath, _ := filepath.Abs(output)
	fmt.Println("Finished export", cmp.Or(abspath, output))
}
