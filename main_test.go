package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with version command (doesn't require project setup)
	os.Args = []string{"velocity", "version"}

	// Run main - it should not panic
	main()
}
