package sevenzip

import (
	"fmt"
	"os"
	"os/exec"
)

// Extract extracts the contents of a 7z archive to the specified destination directory.
// Supports password-protected and split archives.
func Extract(src, dest, password string) error {
	// Ensure the destination directory exists
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Build the command for 7z extraction
	args := []string{"x", src, "-o" + dest, "-y"} // "-y" auto answers "yes" to all prompts
	if password != "" {
		args = append(args, "-p"+password) // Add password flag if provided
	}

	// Run the 7z command to extract the archive
	cmd := exec.Command("7z", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract 7z archive: %w", err)
	}

	return nil
}
