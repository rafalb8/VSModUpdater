package mod

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
)

type Update struct {
	Name     string
	URL      string
	Version  SemVer
	Filename string
}

func UpdateFromString(line string) (upd Update, err error) {
	modid, version, found := strings.Cut(line, "@")
	if !found {
		return upd, fmt.Errorf("failed to parse info")
	}

	semver, err := SemVerFromString(version)
	if err != nil {
		return upd, err
	}

	uri, err := url.JoinPath("https://mods.vintagestory.at/api/mod/", modid)
	if err != nil {
		return upd, fmt.Errorf("UpdateFromString: %w", err)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return upd, fmt.Errorf("UpdateFromString: %w", err)
	}
	defer resp.Body.Close()

	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return upd, fmt.Errorf("UpdateFromString: %w", err)
	}

	upd.Name = r.Mod.Name
	for _, release := range r.Mod.Releases {
		if release.ModVersion.Compare(semver) == 0 {
			upd.URL = release.Mainfile
			upd.Version = release.ModVersion
			upd.Filename = release.Filename
			return
		}
	}
	return upd, fmt.Errorf("UpdateFromString: no release found for %s", modid)
}

func (upd Update) Download() error {
	req, err := http.NewRequest(http.MethodGet, upd.URL, nil)
	if err != nil {
		return err
	}

	// Make sure queries are escaped
	req.URL.RawQuery = url.QueryEscape(req.URL.RawQuery)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Download: status: %s", resp.Status)
	}

	out, err := os.Create(filepath.Join(config.ModPath, upd.Filename))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
