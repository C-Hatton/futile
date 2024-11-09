package createrar

import (
	"fmt"
	"os/exec"
)

// Create creates a RAR archive from the input files and saves it to the destination.
func Create(sources []string, dest, password string) error {
	// Prepare the 7zip command arguments for creating a RAR archive
	cmdArgs := []string{"a", dest}

	// If password is provided, add the -p flag with the password
	if password != "" {
		cmdArgs = append(cmdArgs, "-p"+password)
	}

	// Add the source files to the command arguments
	cmdArgs = append(cmdArgs, sources...)

	// Create the 7zip command
	cmd := exec.Command("7z", cmdArgs...)
	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create RAR archive with 7zip: %w", err)
	}

	return nil
}
