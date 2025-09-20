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
	"runtime"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
)

func Self() error {
	selfLocation, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get self location: %w", err)
	}

	selfInfo := mod.Info{
		ModID:   "5060",
		Version: strings.TrimPrefix(config.VersionNum, "v"),
	}

	fmt.Println("checking for update...")
	update, err := selfInfo.CheckUpdates()
	if err == mod.ErrNoUpdate {
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to check for update: %w", err)
	}

	fmt.Printf("update found, old: %s, new: %s\n", config.VersionNum, update.Version)

	backupLocation := fmt.Sprintf("%s.bak", selfLocation)
	fmt.Println("backing up old version to:", backupLocation)
	err = os.Rename(selfLocation, backupLocation)
	if err != nil {
		return fmt.Errorf("failed to rename old version: %w", err)
	}
	defer func() {
		if err != nil {
			err = os.Rename(backupLocation, selfLocation)
			if err != nil {
				fmt.Printf("failed to undo renaming of '%s', please remove '.bak' manualy", selfLocation)
			}
			return
		}
		// remove old on success
		fmt.Println("removing old", backupLocation)
		os.Remove(backupLocation)
	}()

	fmt.Println("downloading new version...")
	resp, err := http.Get(update.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DownloadFile: status: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(b), resp.ContentLength)
	if err != nil {
		return err
	}

	var f *zip.File
	for _, file := range zipReader.File {
		switch runtime.GOOS {
		case "linux":
			if filepath.Ext(file.Name) == "" {
				f = file
			}
		case "windows":
			if filepath.Ext(file.Name) == "exe" {
				f = file
			}
		}
	}

	newSelf, err := os.Create(selfLocation)
	if err != nil {
		return fmt.Errorf("failed to update self: %w", err)
	}

	zipFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to update self: %w", err)
	}

	_, err = newSelf.ReadFrom(zipFile)
	if err != nil {
		return fmt.Errorf("failed to update self: %w", err)
	}

	err = newSelf.Chmod(0o755)
	if err != nil {
		return fmt.Errorf("failed to update self: %w", err)
	}
	newSelf.Close()

	fmt.Println("testing new version...")
	fmt.Printf("running `%s --help`\n\n", selfLocation)
	cmd := exec.Command(selfLocation, "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to test self: %w", err)
	}

	fmt.Println(string(out))

	return nil
}
