package config

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var ConfigPath string

var VersionNum string = "dev"

// Flags
var (
	BackupPath string
	ModPath    string
	Ignored    map[string]struct{}
	Backup     bool
)

// Modes
var (
	Version     bool
	Self        bool
	List        bool
	Interactive bool
)

func init() {
	var err error
	ConfigPath, err = os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	ConfigPath = filepath.Join(ConfigPath, "VintagestoryData")

	// Flags
	flag.StringVar(&ModPath, "modPath", filepath.Join(ConfigPath, "Mods"), "path to VS mod directory")
	flag.StringVar(&BackupPath, "backupPath", filepath.Join(ConfigPath, "ModBackups"), "path to VS mod backup directory")
	flag.BoolVar(&Backup, "backup", false, "backup mods instead of removing them")
	flag.Func("ignore", "disable updates: modID1,modID2,...", func(s string) error {
		mods := strings.Split(s, ",")
		Ignored = make(map[string]struct{}, len(mods))
		for _, modID := range mods {
			Ignored[modID] = struct{}{}
		}
		return nil
	})

	// Modes
	flag.BoolVar(&Self, "self", false, "update VSModUpdater")
	flag.BoolVar(&Version, "version", false, "print version")
	flag.BoolVar(&List, "list", false, "list mods")
	flag.BoolVar(&Interactive, "interactive", runtime.GOOS == "windows", "interactive update mode")
}
