package archive

import (
	"fmt"
	createrar "futile/archive/create/rar"
	createsevenzip "futile/archive/create/sevenzip"
	createTar "futile/archive/create/tar"
	createzip "futile/archive/create/zip"
	extractrar "futile/archive/extract/rar"
	extractsevenzip "futile/archive/extract/sevenzip"
	extractTar "futile/archive/extract/tar"
	extractzip "futile/archive/extract/zip"
	"futile/utils"
)

// HandleExtract determines the archive type and calls the appropriate extraction function.
func HandleExtract(src, dest, password string) error {
	archiveType, err := utils.DetermineArchiveType(src)
	if err != nil {
		return fmt.Errorf("could not determine archive type: %w", err)
	}

	switch archiveType {
	case "zip":
		// If password is provided, call the password-protected ZIP extraction function
		if password != "" {
			return extractzip.ExtractPasswordProtected(src, dest, password)
		}
		return extractzip.Extract(src, dest) // Standard ZIP extraction
	case "rar":
		return extractrar.Extract(src, dest, password)
	case "7z":
		return extractsevenzip.Extract(src, dest, password)
	case "tar":
		return extractTar.Extract(src, dest, password)
	default:
		return fmt.Errorf("unsupported archive type for extraction: %s", archiveType)
	}
}

// HandleCreate determines the archive type and calls the appropriate creation function.
func HandleCreate(sources []string, dest, password string) error {
	archiveType, err := utils.DetermineArchiveType(dest)
	if err != nil {
		return fmt.Errorf("could not determine archive type: %w", err)
	}

	switch archiveType {
	case "zip":
		// If password is provided, call the password-protected ZIP creation function
		if password != "" {
			return createzip.CreatePasswordProtected(sources, dest, password)
		}
		return createzip.Create(sources, dest) // Standard ZIP creation
	case "rar":
		return createrar.Create(sources, dest, password)
	case "7z":
		return createsevenzip.Create(sources, dest, password)
	case "tar":
		return createTar.Create(sources, dest, password)
	default:
		return fmt.Errorf("unsupported archive type for creation: %s", archiveType)
	}
}
