package modes

import (
	"archive/zip"
	"io"
	"os"
)

func getFile(zipReader *zip.Reader) (io.ReadCloser, error) {
	for _, file := range zipReader.File {
		if file.Name != "VSModUpdater.exe" {
			continue
		}
		return file.Open()
	}
	return nil, os.ErrNotExist
}
