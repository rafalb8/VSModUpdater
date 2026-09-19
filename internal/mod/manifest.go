package mod

import (
	"archive/zip"
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rafalb8/VSModUpdater/v2/internal/config"
	"github.com/tailscale/hujson"
)

// Manifest contains mod metadata read from modinfo.json
//   - [Wiki](https://wiki.vintagestory.at/Modding:Modinfo)
//   - [Docs](https://apidocs.vintagestory.at/api/Vintagestory.API.Common.Manifest.html)
type Manifest struct {
	Path  string   `json:"-"`
	Error error    `json:"-"`
	Page  *ModPage `json:"-"`

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

// ManifestsFromPath returns manifests from zip files or mod folders
func ManifestsFromPath(root string) ([]*Manifest, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	results := make(chan *Manifest, len(entries))
	sem := make(chan struct{}, 128)

	for _, e := range entries {
		wg.Go(func() {
			path := filepath.Join(root, e.Name())
			var modFS fs.FS

			sem <- struct{}{}
			defer func() { <-sem }()

			switch {
			case e.IsDir():
				modFS = os.DirFS(path)

			case filepath.Ext(path) == ".zip":
				r, err := zip.OpenReader(path)
				if err != nil {
					results <- &Manifest{Path: path, Error: err}
					return
				}
				defer r.Close()
				modFS = r

			default:
				return
			}

			results <- parseModFS(modFS, path)
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	mods := make([]*Manifest, 0, len(entries))
	for info := range results {
		mods = append(mods, info)
	}
	slices.SortFunc(mods, func(a, b *Manifest) int { return cmp.Compare(a.Path, b.Path) })
	return mods, nil
}

func parseModFS(modFS fs.FS, path string) *Manifest {
	m := &Manifest{Path: path}

	data, err := fs.ReadFile(modFS, "modinfo.json")
	if err != nil {
		m.Error = err
		return m
	}

	// Sometimes some editors add BOM (Byte Order Mark) to signal endianess.
	// hujson doesn't like that.
	data = bytes.TrimPrefix(data, []byte("\ufeff"))

	// Workaround for non-compliant JSON:
	// Stripping trailing commas here, as a few mods continue
	// to adhere to a looser standard than the parser.
	data, err = hujson.Standardize(data)
	if err != nil {
		m.Error = err
		return m
	}

	err = json.Unmarshal(data, m)
	m.Error = err
	return m
}

// PageURL returns mod page url
func (m *Manifest) PageURL() string {
	if m.Page.UrlAlias != nil {
		uri, _ := url.JoinPath("https://mods.vintagestory.at/", *m.Page.UrlAlias)
		return uri
	}

	uri, _ := url.JoinPath("https://mods.vintagestory.at/show/mod/", strconv.Itoa(m.Page.AssetID))
	return uri
}

func (m *Manifest) String() string {
	if m.Name == "" {
		// Fallback to extracting name from file path
		name := filepath.Base(m.Path)
		return name[:len(name)-len(filepath.Ext(name))]
	}
	return m.Name + "@" + m.Version.String()
}

// Details returns detailed mod info string
func (m *Manifest) Details() string {
	var sb strings.Builder

	// Pre-allocating a rough estimate of the buffer size to avoid dynamic reallocations
	sb.Grow(256)

	if m.Error != nil {
		sb.WriteString("File:\t\t")
		sb.WriteString(filepath.Base(m.Path))
		sb.WriteString("\nError:\t\t")
		sb.WriteString(m.Error.Error())
		return sb.String()
	}

	sb.WriteString("Name:\t\t")
	sb.WriteString(m.Name)

	sb.WriteString("\nModID:\t\t")
	sb.WriteString(m.ModID)

	sb.WriteString("\nVersion:\t")
	sb.WriteString(m.Version.String())

	if gameVer, ok := m.Dependencies["game"]; ok {
		if gameVer == "*" || gameVer == "" {
			gameVer = "any"
		}
		sb.WriteString("\nGame Version:\t")
		sb.WriteString(gameVer)
	}

	sb.WriteString("\nAuthors:\t")
	sb.WriteString(strings.Join(m.Authors, ", "))

	sb.WriteString("\nDescription:\t")
	sb.WriteString(m.Description)

	sb.WriteString("\nURL:\t\t")
	sb.WriteString(m.PageURL())

	return sb.String()
}

// CheckUpdates returns the url to the latest compatible mod version.
func (m *Manifest) CheckUpdates() (Update, error) {
	if m.ModID == "" {
		return Update{}, ErrNoModID
	}

	mod, err := m.FetchModPage()
	if err != nil {
		return Update{}, fmt.Errorf("Info.CheckUpdates: %w", err)
	}

	allowDev := cmp.Or(m.Version.PreRelease(), config.PreRelease)
	return m.findLatestUpdate(mod, allowDev)
}

func (m *Manifest) FetchModPage() (*ModPage, error) {
	if m.Page != nil {
		return m.Page, nil
	}

	cache := filepath.Join(os.TempDir(), "VSModUpdater")
	file := filepath.Join(cache, m.ModID+".json")
	os.MkdirAll(cache, 0o755)

	api := &APIResponse{}

	stat, err := os.Stat(file)
	if err == nil {
		if time.Since(stat.ModTime()) < 15*time.Minute {
			f, err := os.Open(file)
			if err == nil {
				err := json.NewDecoder(f).Decode(api)
				if err == nil {
					m.Page = &api.Mod
					return &api.Mod, nil
				}
			}
		}
	}

	uri, err := url.JoinPath("https://mods.vintagestory.at/api/mod/", m.ModID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, api)
	if err != nil {
		return nil, err
	}

	go os.WriteFile(file, body, 0o644)

	m.Page = &api.Mod
	return &api.Mod, nil
}

func (m *Manifest) findLatestUpdate(mod *ModPage, allowDev bool) (Update, error) {
	err := ErrNoUpdate
	upd := Update{Name: mod.Name}

	for _, rel := range mod.Releases {
		if !allowDev {
			if rel.ModVersion.PreRelease() {
				if err == ErrNoUpdate {
					err = ErrPreReleaseSkip
					upd.Version = rel.ModVersion
				}
				continue
			}

			if IsAllPreRelease(rel.Tags) {
				if err == ErrNoUpdate {
					err = ErrUnstableSkip
					upd.Version = rel.ModVersion
				}
				continue
			}
		}

		// if ModVersion > local, we found update
		if rel.ModVersion.Compare(m.Version) > 0 {
			upd.URL = rel.Mainfile
			upd.Version = rel.ModVersion
			upd.Filename = rel.Filename
			return upd, nil
		}

		break
	}

	return upd, err
}

func (m *Manifest) Backup() error {
	err := os.MkdirAll(config.BackupPath, 0o755)
	if err != nil {
		return err
	}

	oldPath := m.Path
	m.Path = filepath.Join(config.BackupPath, filepath.Base(m.Path))
	return os.Rename(oldPath, m.Path)
}

func (m *Manifest) Restore() error {
	err := os.MkdirAll(config.ModPath, 0o755)
	if err != nil {
		return err
	}

	oldPath := m.Path
	m.Path = filepath.Join(config.ModPath, filepath.Base(m.Path))
	return os.Rename(oldPath, m.Path)
}
