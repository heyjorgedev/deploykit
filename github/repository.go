package github

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type RepositoryService struct {
}

func NewRepositoryService() *RepositoryService {
	return &RepositoryService{}
}

func (s *RepositoryService) DownloadRepository(owner, repo string) error {
	downloadLink := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/master", owner, repo)
	fmt.Println("Downloading repository:", downloadLink)

	// Download the repository
	resp, err := http.Get(downloadLink)
	if err != nil {
		fmt.Println("Error downloading repository:", err)
		return err
	}
	defer resp.Body.Close()

	f, err := os.CreateTemp("", "repo-*.zip")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())

	b, err := io.Copy(f, resp.Body)
	if err != nil {
		fmt.Println("Error copying repository to file:", err)
		return err
	}

	// Unzip the repository
	r, err := zip.NewReader(f, b)
	if err != nil {
		fmt.Println("Error reading zip file:", err)
		return err
	}

	dest, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return err
	}
	defer os.RemoveAll(dest)
	fmt.Println("Unzipping repository to:", dest)

	for _, f := range r.File {
		filePath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
		}

		// 5. Create directory tree
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		// 6. Create a destination file for unzipped content
		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		// 7. Unzip the content of a file and copy it to the destination file
		zippedFile, err := f.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		if _, err := io.Copy(destinationFile, zippedFile); err != nil {
			return err
		}
	}

	fmt.Println("Repository downloaded successfully!")
	return nil
}
