package src

import (
	"os"
)

// ReadFile reads the content of a file and returns it as a string.
// Returns an error if the file cannot be read.
func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes the given content to a file.
// Creates the file if it doesn't exist, overwrites if it does.
// Returns an error if the file cannot be written.
func WriteFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// FileExists checks if a file exists at the given path.
// Returns true if the file exists, false otherwise.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// DeleteFile deletes the file at the given path.
// Returns an error if the file cannot be deleted.
func DeleteFile(filename string) error {
	return os.Remove(filename)
}
