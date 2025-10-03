package config

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var VersionNum string = "v0.0.0"

// Flags
var (
	ModPath     string
	BackupPath  string
	DryRun      bool
	Backup      bool
	Interactive bool
	Ignored     map[string]struct{}
)

// Modes
var (
	Version bool
	Self    bool
	List    bool
)

func init() {
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfgPath = filepath.Join(cfgPath, "VintagestoryData")

	// Flags
	flag.StringVar(&ModPath, "modPath", filepath.Join(cfgPath, "Mods"), "path to VS mod directory")
	flag.StringVar(&BackupPath, "backupPath", filepath.Join(cfgPath, "ModBackups"), "path to VS mod backup directory")
	flag.BoolVar(&DryRun, "dry-run", false, "run the updater without actually doing anything")
	flag.BoolVar(&Backup, "backup", false, "backup mods instead of removing them")
	flag.BoolVar(&Interactive, "interactive", runtime.GOOS == "windows", "interactive update mode")
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
}
