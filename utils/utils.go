package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DetermineArchiveType determines the type of archive based on the file extension
func DetermineArchiveType(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".zip":
		return "zip", nil
	case ".tar":
		return "tar", nil
	case ".rar":
		return "rar", nil
	case ".7z":
		return "7z", nil
	default:
		return "", fmt.Errorf("unsupported or unknown archive type")
	}
}
