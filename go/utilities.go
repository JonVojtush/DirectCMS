package main

import (
	"os"
	"strings"
)

// pathExists checks if a file or directory exists at the given path
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// isMediaFile checks if a file has one of the specified media extensions
func isMediaFile(fileName string) bool {
	for _, ext := range mediaExtensions {
		if strings.HasSuffix(fileName, "."+ext) {
			return true
		}
	}
	return false
}
