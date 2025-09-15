package config

import (
	"flag"
	"os"
	"path/filepath"
)

var ConfigPath string

var VersionNum string = "dev"

// Flags
var (
	BackupPath string
	ModPath    string
	Backup     bool
	Self       bool
	Version    bool
)

func init() {
	var err error
	ConfigPath, err = os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	ConfigPath = filepath.Join(ConfigPath, "/VintagestoryData")

	flag.StringVar(&ModPath, "modPath", filepath.Join(ConfigPath, "Mods"), "path to VS mod directory")
	flag.StringVar(&BackupPath, "backupPath", filepath.Join(ConfigPath, "ModBackups"), "path to VS mod backup directory")
	flag.BoolVar(&Backup, "backup", false, "backup mods instead of removing them")
	flag.BoolVar(&Self, "self", false, "update VSModUpdater")
	flag.BoolVar(&Version, "version", false, "print version")
}
