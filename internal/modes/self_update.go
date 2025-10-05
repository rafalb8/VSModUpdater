package modes

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func Self() {
	selfPath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get self location:", err)
		return
	}

	m := mod.Info{
		ModID:   "5060",
		Version: strings.TrimPrefix(config.VersionNum, "v"),
	}

	fmt.Print("Checking for update - ")
	update, err := m.CheckUpdates()
	if err == mod.ErrNoUpdate {
		fmt.Println("SUCCESS")
		fmt.Println("No updates")
		return
	}
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}
	fmt.Println("SUCCESS")

	fmt.Printf("Downloading: %s => %s - ", m.Version, update.Version)
	resp, err := http.Get(update.URL)
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("FAIL")
		fmt.Println("HTTP status:", resp.Status)
		return
	}
	fmt.Println("SUCCESS")

	fmt.Print("Unzipping - ")
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf), resp.ContentLength)
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}

	f, err := getFile(zipReader)
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}

	basename := filepath.Base(selfPath)
	ext := filepath.Ext(basename)
	newName := fmt.Sprintf("%s_v%s%s", basename[:len(basename)-len(ext)], update.Version, ext)
	newPath := filepath.Join(filepath.Dir(selfPath), newName)

	newSelf, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}

	_, err = newSelf.ReadFrom(f)
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}

	newSelf.Close()
	fmt.Println("SUCCESS")

	fmt.Print("Testing new version - ")
	cmd := exec.Command(newPath, "--help")
	err = cmd.Run()
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}
	fmt.Println("SUCCESS")

	fmt.Printf("Replacing %s => %s - ", filepath.Base(newPath), filepath.Base(selfPath))
	err = os.Rename(newPath, selfPath)
	if err != nil {
		fmt.Println("FAIL")
		fmt.Println(err)
		return
	}
	fmt.Println("SUCCESS")
}
