package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Create creates a standard ZIP archive from the provided source files and directories.
func Create(sources []string, dest string) error {
	// Create the ZIP file
	zipFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file %s: %w", dest, err)
	}
	// Ensure the zipFile is closed
	defer func() {
		if closeErr := zipFile.Close(); closeErr != nil {
			fmt.Printf("Error closing ZIP file %s: %v\n", dest, closeErr)
		}
	}()

	// Create a new zip.Writer
	zipWriter := zip.NewWriter(zipFile)
	// Ensure the zipWriter is closed
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			fmt.Printf("Error closing ZIP writer for %s: %v\n", dest, closeErr)
		}
	}()

	// Iterate over each source file or directory
	for _, source := range sources {
		err := addToZip(source, zipWriter)
		if err != nil {
			return fmt.Errorf("failed to add %s to ZIP: %w", source, err)
		}
	}

	return nil
}

// addToZip adds a file or directory to the ZIP archive.
func addToZip(source string, zipWriter *zip.Writer) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("failed to stat source %s: %w", source, err)
	}

	if info.IsDir() {
		// Use manual file opening and closing inside loop to avoid `defer` in loops
		dirWalkErr := filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking through directory %s: %w", source, err)
			}
			if file == source {
				return nil
			}
			return addFileToZip(file, source, zipWriter)
		})
		if dirWalkErr != nil {
			return dirWalkErr
		}
	} else {
		// For files, add them directly
		return addFileToZip(source, "", zipWriter)
	}

	return nil
}

// addFileToZip adds a single file to the ZIP archive.
func addFileToZip(file, baseDir string, zipWriter *zip.Writer) error {
	fileToZip, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file, err)
	}
	// Close the file explicitly
	closeErr := fileToZip.Close()
	if closeErr != nil {
		fmt.Printf("Error closing file %s: %v\n", file, closeErr)
	}

	relativePath := file
	if baseDir != "" {
		relativePath, _ = filepath.Rel(baseDir, file)
	}
	fileHeader := &zip.FileHeader{
		Name:   relativePath,
		Method: zip.Deflate,
	}

	writer, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return fmt.Errorf("failed to create header for file %s: %w", file, err)
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return fmt.Errorf("failed to write file %s to ZIP: %w", file, err)
	}

	return nil
}

// CreatePasswordProtected creates a password-protected ZIP archive using 7z.
func CreatePasswordProtected(sources []string, dest, password string) error {
	cmd := exec.Command("7z", "a", "-p"+password, dest)
	cmd.Args = append(cmd.Args, sources...)

	// Execute the 7z command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create password-protected ZIP file %s: %w", dest, err)
	}

	return nil
}
