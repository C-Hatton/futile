package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Extract extracts the contents of a standard ZIP archive to the destination.
func Extract(src, dest string) error {
	archive, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file %s: %w", src, err)
	}
	// Close the archive explicitly
	closeErr := archive.Close()
	if closeErr != nil {
		fmt.Printf("Error closing ZIP archive %s: %v\n", src, closeErr)
	}

	for _, file := range archive.File {
		if strings.HasPrefix(file.Name, "__MACOSX") || strings.HasPrefix(file.Name, "._") {
			continue
		}

		destFilePath := filepath.Join(dest, file.Name)

		destDir := filepath.Dir(destFilePath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		inFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s in archive: %w", file.Name, err)
		}

		// Explicitly close input file
		closeErr = inFile.Close()
		if closeErr != nil {
			fmt.Printf("Error closing input file %s: %v\n", file.Name, closeErr)
		}

		outFile, err := os.Create(destFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destFilePath, err)
		}

		// Explicitly close output file
		closeErr = outFile.Close()
		if closeErr != nil {
			fmt.Printf("Error closing output file %s: %v\n", destFilePath, closeErr)
		}

		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return fmt.Errorf("failed to extract file %s to %s: %w", file.Name, destFilePath, err)
		}
	}

	return nil
}

// ExtractPasswordProtected extracts a password-protected ZIP archive using 7z.
func ExtractPasswordProtected(src, dest, password string) error {
	cmd := exec.Command("7z", "x", "-p"+password, src, "-o"+dest)

	// Execute the 7z command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to extract password-protected ZIP file %s: %w", src, err)
	}

	return nil
}
