package mod

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	Version          string            `json:"version"`
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

func InfoFromZip(path string) (*Info, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != "modinfo.json" {
			continue
		}

		fr, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer fr.Close()

		data, err := io.ReadAll(fr)
		if err != nil {
			return nil, err
		}

		// Sometimes some editors add BOM (Byte Order Mark) to signal endianess.
		// hujson doesn't like that.
		data = bytes.TrimPrefix(data, []byte("\ufeff"))

		// Workaround for non-compliant JSON:
		// Stripping trailing commas here, as a few mods continue
		// to adhere to a looser standard than the parser.
		data, err = hujson.Standardize(data)
		if err != nil {
			return nil, err
		}

		info := &Info{Path: path}
		return info, json.Unmarshal(data, info)
	}
	return nil, fmt.Errorf("mod.InfoFromZip: no files found in %s", path)
}

func InfoFromPath(path string) ([]*Info, error) {
	mods := []*Info{}
	err := filepath.WalkDir(config.ModPath, func(path string, d fs.DirEntry, err error) error {
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

		info, err := InfoFromZip(path)
		if err != nil {
			info = &Info{Path: path, Error: err}
		}
		mods = append(mods, info)
		return nil
	})
	return mods, err
}

func (i *Info) String() string {
	if i.Name == "" {
		// Fallback to extracting name from file path
		name := filepath.Base(i.Path)
		return name[:len(name)-len(filepath.Ext(name))]
	}
	return i.Name + "@v" + i.Version
}

// Details returns detailed mod info string
func (i *Info) Details() string {
	sb := strings.Builder{}
	if i.Error != nil {
		sb.WriteString("File:\t\t" + filepath.Base(i.Path) + "\n")
		sb.WriteString("Error:\t\t" + i.Error.Error())
		return sb.String()
	}

	sb.WriteString("Name:\t\t" + i.Name + "\n")
	sb.WriteString("ModID:\t\t" + i.ModID + "\n")
	sb.WriteString("Version:\t" + i.Version + "\n")
	gameVer, ok := i.Dependencies["game"]
	if ok {
		if gameVer == "*" {
			gameVer = "any"
		}
		sb.WriteString("Game Version:\t" + gameVer + "\n")
	}
	sb.WriteString("Authors:\t" + strings.Join(i.Authors, ", ") + "\n")
	sb.WriteString("Description:\t" + i.Description)
	return sb.String()
}

// CheckUpdates returns url to latest mod version
func (i *Info) CheckUpdates() (Update, error) {
	if i.ModID == "" {
		return Update{}, ErrNoModID
	}

	uri, err := url.JoinPath("https://mods.vintagestory.at/api/mod/", i.ModID)
	if err != nil {
		return Update{}, fmt.Errorf("Info.CheckUpdates: %w", err)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return Update{}, fmt.Errorf("Info.CheckUpdates: %w", err)
	}
	defer resp.Body.Close()

	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return Update{}, fmt.Errorf("Info.CheckUpdates: %w", err)
	}

	for _, release := range r.Mod.Releases {
		if release.ModVersion > i.Version {
			return Update{
				URL:      release.Mainfile,
				Version:  release.ModVersion,
				Filename: release.Filename,
			}, nil
		} else {
			return Update{}, ErrNoUpdate
		}
	}
	return Update{}, fmt.Errorf("Info.CheckUpdates: no release found for %s", i.ModID)
}

func (i *Info) Backup() error {
	err := os.MkdirAll(config.BackupPath, 0o755)
	if err != nil {
		return err
	}

	return os.Rename(i.Path, filepath.Join(config.BackupPath, filepath.Base(i.Path)))
}
