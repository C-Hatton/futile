package createTar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

// Create creates a tar archive from the input files and saves it to the destination.
// If a password is provided, it uses 7zip for password protection.
func Create(sources []string, dest, password string) error {
	if password != "" {
		// Use 7zip to create a password-protected tar archive
		return createPasswordProtectedTar(sources, dest, password)
	}

	// Standard tar archive creation
	return createStandardTar(sources, dest)
}

// createStandardTar creates a standard (non-password protected) tar archive.
func createStandardTar(sources []string, dest string) error {
	// Open the tar file for writing
	tarFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create tar file: %w", err)
	}
	defer func() {
		// Ensure the tarFile is closed and handle any closing error
		if closeErr := closeFile(tarFile); closeErr != nil {
			fmt.Printf("Warning: failed to close TAR file: %v\n", closeErr)
		}
	}()

	// Create a new tar writer
	tarWriter := tar.NewWriter(tarFile)
	defer func() {
		// Ensure the tarWriter is closed and handle any closing error
		if closeErr := closeFile(tarWriter); closeErr != nil {
			fmt.Printf("Warning: failed to close TAR writer: %v\n", closeErr)
		}
	}()

	// Loop over the input files and add them to the tar archive
	for _, filePath := range sources {
		// Open the file to be archived
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("could not open file %s: %w", filePath, err)
		}

		// Get the file info
		fileInfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("could not stat file %s: %w", filePath, err)
		}

		// Create a tar header for this file
		header := &tar.Header{
			Name:    filepath.Base(filePath),
			Size:    fileInfo.Size(),
			Mode:    int64(fileInfo.Mode()),
			ModTime: fileInfo.ModTime(),
		}

		// Write the header to the tar archive
		err = tarWriter.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("could not write header for file %s: %w", filePath, err)
		}

		// Write the file contents to the tar archive
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return fmt.Errorf("could not write contents of file %s: %w", filePath, err)
		}

		// Ensure the file is closed after processing
		if err := closeFile(file); err != nil {
			fmt.Printf("Warning: failed to close file %s: %v\n", filePath, err)
		}
	}

	return nil
}

// createPasswordProtectedTar uses 7zip to create a password-protected tar archive.
func createPasswordProtectedTar(sources []string, dest, password string) error {
	// Create a temporary tar file without password protection first
	tempTar := dest + ".temp"
	err := createStandardTar(sources, tempTar)
	if err != nil {
		return fmt.Errorf("failed to create temporary tar file: %w", err)
	}

	// Use 7zip to apply the password to the temporary tar file
	cmd := exec.Command("7z", "a", "-p"+password, dest, tempTar)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create password-protected tar archive with 7zip: %w", err)
	}

	// Clean up the temporary tar file
	err = os.Remove(tempTar)
	if err != nil {
		return fmt.Errorf("failed to remove temporary tar file: %w", err)
	}

	return nil
}
