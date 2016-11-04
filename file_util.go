package logger

import (
	"os"
	"path/filepath"
	"regexp"
)

// File contains helper functions for files.
var File = &fileUtil{}

type fileUtil struct{}

// CreateOrOpen creates or opens a file.
func (fu fileUtil) CreateOrOpen(filePath string) (*os.File, error) {
	f, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return os.Create(filePath)
	}
	return f, err
}

// CreateAndClose creates and closes a file.
func (fu fileUtil) CreateAndClose(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// RemoveMany removes an array of files.
func (fu fileUtil) RemoveMany(filePaths ...string) error {
	var err error
	for _, path := range filePaths {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	return err
}

func (fu fileUtil) List(path string, expr *regexp.Regexp) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(fullFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if expr == nil {
			files = append(files, fullFilePath)
		} else if expr.MatchString(info.Name()) {
			files = append(files, fullFilePath)
		}
		return nil
	})
	return files, err
}
