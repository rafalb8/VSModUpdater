package mod

import (
	"fmt"
	"io"
	"net/http"
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
	resp, err := http.Get(upd.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DownloadFile: status: %s", resp.Status)
	}

	out, err := os.Create(filepath.Join(config.ModPath, upd.Filename))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
