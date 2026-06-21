package handlers

import (
	"archive/zip"
	"backend/internal/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (h *Handler) BackupHandler(w http.ResponseWriter, r *http.Request) {

	timestamp := time.Now().Format("2006-01-02_150405")
	filename := fmt.Sprintf("backup_%s.zip", timestamp)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, filename),
	)

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// backup data directory
	if err := addDirectoryToZip(zipWriter, "./data", "data"); err != nil {
		utils.ErrorJson(
			w,
			http.StatusInternalServerError,
			"failed to backup data directory",
		)
		return
	}

	// backup public directory
	if err := addDirectoryToZip(zipWriter, "./public", "public"); err != nil {
		utils.ErrorJson(
			w,
			http.StatusInternalServerError,
			"failed to backup public directory",
		)
		return
	}
}

func addDirectoryToZip(
	zipWriter *zip.Writer,
	sourceDir string,
	zipRoot string,
) error {

	return filepath.Walk(
		sourceDir,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			relativePath, err := filepath.Rel(sourceDir, path)

			if err != nil {
				return err
			}

			zipPath := filepath.Join(
				zipRoot,
				relativePath,
			)

			return addFileToZip(
				zipWriter,
				path,
				zipPath,
			)
		},
	)
}

func addFileToZip(
	zipWriter *zip.Writer,
	filePath string,
	zipName string,
) error {

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	writer, err := zipWriter.Create(zipName)

	if err != nil {
		return err
	}

	_, err = io.Copy(
		writer,
		file,
	)

	return err
}
