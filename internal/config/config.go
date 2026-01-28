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
	PreRelease  bool
	Ignored     map[string]struct{}
)

// Modes
var (
	Version bool
	Self    bool
	List    bool
	Import  string
	Export  string
)

func init() {
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfgPath = filepath.Join(cfgPath, "VintagestoryData")

	// Flags
	flag.StringVar(&ModPath, "mod-path", filepath.Join(cfgPath, "Mods"), "path to VS mod directory")
	flag.StringVar(&BackupPath, "backup-path", "", "path to VS mod backup directory")
	flag.BoolVar(&DryRun, "dry-run", false, "run the updater without actually doing anything")
	flag.BoolVar(&Backup, "backup", false, "backup mods instead of removing them")
	flag.BoolVar(&Interactive, "interactive", runtime.GOOS != "linux", "interactive update mode")
	flag.BoolVar(&PreRelease, "pre-release", false, "allow updating to pre-release mod versions (enabled if mod is already pre-release)")
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
	flag.StringVar(&Import, "import", "", "import mod list")
	flag.StringVar(&Export, "export", "", "export mod list")

	flag.Parse()

	// Make sure modpath is absolute path
	ModPath, err = filepath.Abs(ModPath)
	if err != nil {
		panic(err)
	}

	if BackupPath == "" {
		// Set backup path as a sibling of mod path
		BackupPath = filepath.Join(filepath.Dir(ModPath), "ModBackups")
	}
}
