package tests

import (
	"go-service/src"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test content
	testContent := "Hello, World!"
	if _, err := tmpfile.Write([]byte(testContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test reading the file
	content, err := src.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("ReadFile returned error: %v", err)
	}
	if content != testContent {
		t.Errorf("ReadFile returned %q, want %q", content, testContent)
	}
}

func TestReadFileNonExistent(t *testing.T) {
	_, err := src.ReadFile("nonexistent.txt")
	if err == nil {
		t.Error("Expected error when reading non-existent file")
	}
}

func TestWriteFile(t *testing.T) {
	filename := "test_write.txt"
	testContent := "Test content"

	// Clean up before and after
	os.Remove(filename)
	defer os.Remove(filename)

	// Write file
	err := src.WriteFile(filename, testContent)
	if err != nil {
		t.Errorf("WriteFile returned error: %v", err)
	}

	// Verify file was written
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != testContent {
		t.Errorf("Written content was %q, want %q", string(content), testContent)
	}
}

func TestFileExists(t *testing.T) {
	// Create a test file
	tmpfile, err := os.CreateTemp("", "exists*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if !src.FileExists(tmpfile.Name()) {
		t.Error("FileExists returned false for existing file")
	}

	if src.FileExists("nonexistent.txt") {
		t.Error("FileExists returned true for non-existent file")
	}
}

func TestDeleteFile(t *testing.T) {
	// Create a test file
	tmpfile, err := os.CreateTemp("", "delete*.txt")
	if err != nil {
		t.Fatal(err)
	}
	filename := tmpfile.Name()
	tmpfile.Close()

	// Verify file exists
	if !src.FileExists(filename) {
		t.Fatal("Test file was not created")
	}

	// Delete the file
	err = src.DeleteFile(filename)
	if err != nil {
		t.Errorf("DeleteFile returned error: %v", err)
	}

	// Verify file no longer exists
	if src.FileExists(filename) {
		t.Error("File still exists after deletion")
	}
}

func TestDeleteFileNonExistent(t *testing.T) {
	err := src.DeleteFile("nonexistent.txt")
	if err == nil {
		t.Error("Expected error when deleting non-existent file")
	}
}
