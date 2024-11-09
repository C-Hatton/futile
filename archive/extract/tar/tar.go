package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// closeFile is a helper function to close files and handle errors.
func closeFile(f io.Closer) error {
	if f != nil {
		err := f.Close()
		if err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
	}
	return nil
}

// Extract extracts the contents of a TAR archive.
// If a password is provided, it uses 7zip for extraction.
func Extract(src, dest, password string) error {
	if password != "" {
		// Use 7zip to extract password-protected tar archive
		return extractPasswordProtectedTar(src, dest, password)
	}

	// Standard tar extraction
	return extractStandardTar(src, dest)
}

// extractStandardTar extracts a non-password-protected tar archive.
func extractStandardTar(src, dest string) error {
	// Open the TAR archive
	tarFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open TAR file: %w", err)
	}
	defer func() {
		// Ensure the tarFile is closed and handle any closing error
		if closeErr := closeFile(tarFile); closeErr != nil {
			fmt.Printf("Warning: failed to close TAR file: %v\n", closeErr)
		}
	}()

	// Create a new tar.Reader to read the TAR archive
	tarReader := tar.NewReader(tarFile)

	// Extract the contents
	return extractTarContents(tarReader, dest)
}

// extractPasswordProtectedTar uses 7zip to extract a password-protected tar archive.
func extractPasswordProtectedTar(src, dest, password string) error {
	// Use 7zip to extract the password-protected tar archive
	cmd := exec.Command("7z", "x", "-p"+password, src, "-o"+dest)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to extract password-protected tar archive with 7zip: %w", err)
	}

	return nil
}

// extractTarContents extracts the content of the TAR archive using the tar.Reader.
func extractTarContents(tarReader *tar.Reader, dest string) error {
	// Loop through the TAR file
	for {
		// Get the next file in the TAR archive
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of the TAR archive
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read TAR header: %w", err)
		}

		// Build the destination file path
		destPath := fmt.Sprintf("%s/%s", dest, header.Name)

		// Handle directories
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			continue
		}

		// Handle regular files
		file, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}

		// Ensure the file is closed after processing
		if err := closeFile(file); err != nil {
			return fmt.Errorf("failed to close file %s: %w", destPath, err)
		}

		// Copy the file contents
		_, err = io.Copy(file, tarReader)
		if err != nil {
			_ = closeFile(file)
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}
	}

	return nil
}
