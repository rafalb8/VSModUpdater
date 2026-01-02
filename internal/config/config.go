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
	Version       bool
	Self          bool
	List          bool
	Webpage       bool
	WebConfigMode bool
	ConfigModID   string
	ConfigAddCat  string
	ConfigEditCat string
	ConfigDelCat  string
)

// Webpage options
var (
	WebpageOutput     string
	WebpageTitle      string
	WebpageDeploy     string
	WebpageDeployFlag bool
	WebConfigFile     string
)

func init() {
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfgPath = filepath.Join(cfgPath, "VintagestoryData")

	// Flags
	flag.StringVar(&ModPath, "mod-path", filepath.Join(cfgPath, "Mods"), "path to VS mod directory")
	flag.StringVar(&BackupPath, "backup-path", filepath.Join(cfgPath, "ModBackups"), "path to VS mod backup directory")
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
	flag.BoolVar(&Webpage, "webpage", false, "generate static HTML modlist webpage")
	flag.BoolVar(&WebConfigMode, "web-config", false, "interactive configuration mode for mod descriptions and categories")
	flag.StringVar(&ConfigModID, "edit-mod", "", "edit specific mod configuration (e.g., -edit-mod=alchemy)")
	flag.StringVar(&ConfigAddCat, "add-category", "", "add a new category (format: 'name:color')")
	flag.StringVar(&ConfigEditCat, "edit-category", "", "edit existing category (format: 'oldname:newname:color')")
	flag.StringVar(&ConfigDelCat, "delete-category", "", "delete a category by name")

	// Webpage options
	flag.StringVar(&WebpageOutput, "output", "modlist.html", "output filename for webpage")
	flag.StringVar(&WebpageTitle, "title", "", "title for the webpage (overrides config file)")
	flag.BoolVar(&WebpageDeployFlag, "deploy", false, "deploy to Cloudflare Pages (uses projectName from config)")
	flag.StringVar(&WebpageDeploy, "deploy-project", "", "deploy to Cloudflare Pages with specified project name")
	flag.StringVar(&WebConfigFile, "web-config-file", "webpage_config.json", "path to webpage configuration file")
}
