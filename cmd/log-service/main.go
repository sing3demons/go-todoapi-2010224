package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dirPath := "../../logs/details"

	// Read the directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, entry := range entries {
		// Process only .log.gz files
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".log.gz") {
			// Construct the full file path
			filePath := filepath.Join(dirPath, entry.Name())

			// Read and print the content of each .log.gz file
			fmt.Printf("Reading file: %s\n", filePath)
			if err := readGzFile(filePath); err != nil {
				log.Printf("Error reading file %s: %v", filePath, err)
			}
		}
	}
}

func readGzFile(filePath string) error {
	// Open the .gz file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Read and print each line
	scanner := bufio.NewScanner(gzReader)
	for scanner.Scan() {
		fmt.Println("--------------------")
		fmt.Println(scanner.Text())
	}

	// remove the file
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}
