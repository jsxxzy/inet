package xfs

import (
	"fmt"
	"os"
)

// Exists check file/dir exists
//
// https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateIfNotExists auto create files
//
// https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}
