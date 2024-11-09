package main

import (
	"flag"
	"fmt"
	"futile/archive"
	"log"
	"path/filepath"
)

// Print help message
func printUsage() {
	fmt.Println(`Usage: futile [options]

Options:
  -o, --operation      Operation to perform ('create' or 'extract')
  -i, --input          Archive input (for extraction or creation)
  -d, --destination    Destination directory or path for extraction or archive creation
  -p, --password       Password for password-protected archives
  -h, --help           Show help message
  -v, --version        Show version information`)
}

// Print version information
func printVersion() {
	fmt.Println("Futile version 1.0.0")
}

func main() {
	// Declare flags
	operation := flag.String("o", "", "Operation: 'create' or 'extract'")
	destination := flag.String("d", "", "Destination directory or path for extraction or archive creation")
	password := flag.String("p", "", "Password for password-protected archives")
	help := flag.Bool("h", false, "Show help message")
	version := flag.Bool("v", false, "Show version information")

	// Custom flag to capture multiple input files for creation or extraction
	var inputFiles []string
	flag.Func("i", "Archive input (for extraction or creation)", func(s string) error {
		inputFiles = append(inputFiles, s)
		return nil
	})

	// Parse the command line arguments
	flag.Parse()

	// Handle help or version flags
	if *help || *version {
		if *help {
			printUsage()
		} else {
			printVersion()
		}
		return
	}

	// Ensure operation flag is provided
	if *operation == "" {
		log.Fatal("Operation (-o) is required")
	}

	// Ensure either create or extract operation is chosen
	if *operation != "extract" && *operation != "create" {
		log.Fatal("Invalid operation. Use 'extract' or 'create'.")
	}

	// Validate flags based on the selected operation
	if *operation == "create" {
		// For 'create' operation, ensure destination (-d) and input files (-i) are provided
		if *destination == "" || len(inputFiles) == 0 {
			log.Fatal("Both destination (-d) and at least one input file (-i) are required for creating an archive")
		}
	} else if *operation == "extract" {
		// For 'extract' operation, ensure archive input (-i) is provided
		if len(inputFiles) == 0 {
			log.Fatal("Archive input (-i) is required for extraction")
		}

		// Set the destination to the directory of the archive if not provided
		if *destination == "" {
			// Get the directory of the input archive
			dir := filepath.Dir(inputFiles[0])
			// Set destination to the same directory as the archive
			*destination = dir
		}
	}

	// Handle the operation based on user input
	var err error
	switch *operation {
	case "extract":
		// Handle extraction with password
		err = archive.HandleExtract(inputFiles[0], *destination, *password)
	case "create":
		// Handle creation with password
		err = archive.HandleCreate(inputFiles, *destination, *password)
	}

	// If there was an error, log and exit
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Operation completed successfully")
}
