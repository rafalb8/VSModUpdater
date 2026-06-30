package mod

import (
	"archive/zip"
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/tailscale/hujson"
)

// Info contains mod metadata
//   - [Wiki](https://wiki.vintagestory.at/Modding:Modinfo)
//   - [Docs](https://apidocs.vintagestory.at/api/Vintagestory.API.Common.Info.html)
type Info struct {
	Path  string `json:"-"`
	Error error  `json:"-"`

	Type             Type              `json:"type"`
	Name             string            `json:"name"`
	ModID            string            `json:"modid,omitempty"`
	Version          SemVer            `json:"version"`
	NetworkVersion   string            `json:"networkVersion,omitempty"`
	TextureSize      int               `json:"textureSize,omitempty"`
	Description      string            `json:"description,omitempty"`
	Website          string            `json:"website,omitempty"`
	IconPath         string            `json:"iconPath,omitempty"`
	Authors          []string          `json:"authors,omitempty"`
	Contributors     []string          `json:"contributors,omitempty"`
	Side             AppSide           `json:"side,omitempty"`
	RequiredOnClient Bool              `json:"requiredOnClient,omitempty"`
	RequiredOnServer Bool              `json:"requiredOnServer,omitempty"`
	Dependencies     map[string]string `json:"dependencies,omitempty"`
}

// Returns Info slice from zip files
func InfoFromPath(path string) ([]*Info, error) {
	mods := []*Info{}
	err := filepath.WalkDir(config.ModPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		var modFS fs.FS

		switch {
		case d.IsDir():
			if path == config.ModPath {
				return nil
			}
			modFS = os.DirFS(path)
			err = fs.SkipDir

		case filepath.Ext(path) == ".zip":
			r, err := zip.OpenReader(path)
			if err != nil {
				mods = append(mods, &Info{Path: path, Error: err})
				return nil
			}
			defer r.Close()
			modFS = r

		default:
			return nil
		}

		mods = append(mods, parseModFS(modFS, path))
		return err
	})
	return mods, err
}

func parseModFS(modFS fs.FS, path string) *Info {
	info := &Info{Path: path}

	data, err := fs.ReadFile(modFS, "modinfo.json")
	if err != nil {
		info.Error = err
		return info
	}

	// Sometimes some editors add BOM (Byte Order Mark) to signal endianess.
	// hujson doesn't like that.
	data = bytes.TrimPrefix(data, []byte("\ufeff"))

	// Workaround for non-compliant JSON:
	// Stripping trailing commas here, as a few mods continue
	// to adhere to a looser standard than the parser.
	data, err = hujson.Standardize(data)
	if err != nil {
		info.Error = err
		return info
	}

	err = json.Unmarshal(data, info)
	info.Error = err
	return info
}

func (i *Info) String() string {
	if i.Name == "" {
		// Fallback to extracting name from file path
		name := filepath.Base(i.Path)
		return name[:len(name)-len(filepath.Ext(name))]
	}
	return i.Name + "@" + i.Version.String()
}

// Details returns detailed mod info string
func (i *Info) Details() string {
	var sb strings.Builder

	// Pre-allocating a rough estimate of the buffer size to avoid dynamic reallocations
	sb.Grow(256)

	if i.Error != nil {
		sb.WriteString("File:\t\t")
		sb.WriteString(filepath.Base(i.Path))
		sb.WriteString("\nError:\t\t")
		sb.WriteString(i.Error.Error())
		return sb.String()
	}

	sb.WriteString("Name:\t\t")
	sb.WriteString(i.Name)

	sb.WriteString("\nModID:\t\t")
	sb.WriteString(i.ModID)

	sb.WriteString("\nVersion:\t")
	sb.WriteString(i.Version.String())

	if gameVer, ok := i.Dependencies["game"]; ok {
		if gameVer == "*" || gameVer == "" {
			gameVer = "any"
		}
		sb.WriteString("\nGame Version:\t")
		sb.WriteString(gameVer)
	}

	sb.WriteString("\nAuthors:\t")
	sb.WriteString(strings.Join(i.Authors, ", "))

	sb.WriteString("\nDescription:\t")
	sb.WriteString(i.Description)

	return sb.String()
}

// CheckUpdates returns the url to the latest compatible mod version.
func (i *Info) CheckUpdates() (Update, error) {
	if i.ModID == "" {
		return Update{}, ErrNoModID
	}

	mod, err := i.fetchReleases()
	if err != nil {
		return Update{}, fmt.Errorf("Info.CheckUpdates: %w", err)
	}

	allowDev := cmp.Or(i.Version.PreRelease(), config.PreRelease)

	return i.findLatestUpdate(mod, allowDev)
}

func (i *Info) fetchReleases() (*Mod, error) {
	uri, err := url.JoinPath("https://mods.vintagestory.at/api/mod/", i.ModID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	return &r.Mod, nil
}

func (i *Info) findLatestUpdate(mod *Mod, allowDev bool) (Update, error) {
	hasSkippedDev := false
	upd := Update{Name: mod.Name}

	for _, rel := range mod.Releases {
		// skip pre-releases if they aren't allowed
		if !allowDev && (rel.ModVersion.PreRelease() || IsAllPreRelease(rel.Tags)) {
			hasSkippedDev = true
			continue
		}

		// if ModVersion > local, we found update
		if rel.ModVersion.Compare(i.Version) > 0 {
			upd.URL = rel.Mainfile
			upd.Version = rel.ModVersion
			upd.Filename = rel.Filename
			return upd, nil
		}

		break
	}

	if hasSkippedDev {
		return upd, ErrPreReleaseSkip
	}
	return upd, ErrNoUpdate
}

func (i *Info) Backup() error {
	err := os.MkdirAll(config.BackupPath, 0o755)
	if err != nil {
		return err
	}

	oldPath := i.Path
	i.Path = filepath.Join(config.BackupPath, filepath.Base(i.Path))
	return os.Rename(oldPath, i.Path)
}

func (i *Info) Restore() error {
	err := os.MkdirAll(config.ModPath, 0o755)
	if err != nil {
		return err
	}

	oldPath := i.Path
	i.Path = filepath.Join(config.ModPath, filepath.Base(i.Path))
	return os.Rename(oldPath, i.Path)
}
