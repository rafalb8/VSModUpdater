package config

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/pflag"
)

// Flags
var (
	ModPath     string
	BackupPath  string
	DryRun      bool
	Backup      bool
	Interactive bool
	PreRelease  bool
	Ignored     = map[string]struct{}{}
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
	pflag.StringVarP(&ModPath, "mod-path", "m", filepath.Join(cfgPath, "Mods"), "path to VS mod directory")
	pflag.StringVar(&BackupPath, "backup-path", "", "path to VS mod backup directory")
	pflag.BoolVarP(&DryRun, "dry-run", "p", false, "run the updater without actually doing anything")
	pflag.BoolVarP(&Backup, "backup", "b", false, "backup mods instead of removing them")
	pflag.BoolVarP(&Interactive, "interactive", "t", runtime.GOOS != "linux", "interactive update mode")
	pflag.BoolVar(&PreRelease, "pre-release", false, "allow updating to pre-release mod versions (enabled if mod is already pre-release)")
	pflag.FuncP("ignore", "x", "disable updates: modID1,modID2,...", func(s string) error {
		for modID := range strings.SplitSeq(s, ",") {
			modID = strings.TrimSpace(modID)
			if modID != "" {
				Ignored[modID] = struct{}{}
			}
		}
		return nil
	})

	// Modes
	pflag.BoolVar(&Self, "self", false, "update VSModUpdater")
	pflag.BoolVarP(&Version, "version", "v", false, "print version")
	pflag.BoolVarP(&List, "list", "l", false, "list mods")
	pflag.StringVarP(&Import, "import", "i", "", "import mod list")
	pflag.StringVarP(&Export, "export", "e", "", "export mod list")

	// Parse flags
	pflag.Parse()

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

var version = "v0.0.0"

func BuildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}
