package sevenzip

import (
	"fmt"
	"os"
	"os/exec"
)

// Create creates a 7z archive from the provided source files and directories.
// Supports password protection and split archives.
func Create(sources []string, dest, password string) error {
	// Build the command arguments for 7z
	args := append([]string{"a", dest}, sources...)
	if password != "" {
		args = append(args, "-p"+password) // Add password flag if provided
	}

	// Run the 7z command to create the archive
	cmd := exec.Command("7z", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create 7z archive: %w", err)
	}

	return nil
}
