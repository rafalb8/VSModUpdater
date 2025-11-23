package mod

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/rafalb8/VSModUpdater/internal/config"
)

type Update struct {
	URL      string
	Version  SemVer
	Filename string
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
