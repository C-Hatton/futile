package extractrar

import (
	"fmt"
	"os"
	"os/exec"
)

// Extract extracts the contents of a RAR archive.
func Extract(src, dest, password string) error {
	// Ensure the destination directory exists
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Prepare the 7zip command arguments for extracting a RAR archive
	cmdArgs := []string{"x", src, "-o" + dest}

	// If password is provided, add the -p flag with the password
	if password != "" {
		cmdArgs = append(cmdArgs, "-p"+password)
	}

	// Create the 7zip command
	cmd := exec.Command("7z", cmdArgs...)
	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to extract RAR archive with 7zip: %w", err)
	}

	return nil
}
