package mod

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rafalb8/VSModUpdater/internal/config"
)

// Info contains mod metadata
//   - [Wiki](https://wiki.vintagestory.at/Modding:Modinfo)
//   - [Docs](https://apidocs.vintagestory.at/api/Vintagestory.API.Common.Info.html)
type Info struct {
	Path string `json:"-"`

	Type    Type   `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
	ModID   string `json:"modid,omitempty"`

	// TODO: fix unmarshal string bools
	// NetworkVersion   string            `json:"networkVersion,omitempty"`
	// TextureSize      int               `json:"textureSize,omitempty"`
	// Description      string            `json:"description,omitempty"`
	// Website          string            `json:"website,omitempty"`
	// IconPath         string            `json:"iconPath,omitempty"`
	// Authors          []string          `json:"authors,omitempty"`
	// Contributors     []string          `json:"contributors,omitempty"`
	// Side             AppSide           `json:"side,omitempty"`
	// RequiredOnClient bool              `json:"requiredOnClient,omitempty"`
	// RequiredOnServer bool              `json:"requiredOnServer,omitempty"`
	// Dependencies     map[string]string `json:"dependencies,omitempty"`
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

		info := &Info{Path: path}
		return info, json.NewDecoder(fr).Decode(info)
	}
	return nil, fmt.Errorf("mod.InfoFromZip: no files found in %s", path)
}

func (i *Info) String() string {
	return i.Name + "@" + i.Version
}

// CheckUpdates returns url to latest mod version
func (i *Info) CheckUpdates() (Update, error) {
	// TODO: Fix search when ModID is missing
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
