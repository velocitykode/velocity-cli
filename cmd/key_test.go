package cmd

import (
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Mock reader that fails
type failingReader struct{}

func (f failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock read error")
}

func TestGenerateKey(t *testing.T) {
	key, err := generateKey()
	if err != nil {
		t.Fatalf("generateKey() error = %v", err)
	}

	if len(key) != 44 {
		t.Errorf("generateKey() key length = %d, want 44", len(key))
	}
}

func TestGenerateKeyUniqueness(t *testing.T) {
	key1, _ := generateKey()
	key2, _ := generateKey()

	if key1 == key2 {
		t.Error("generateKey() should generate unique keys")
	}
}

func TestGenerateKeyError(t *testing.T) {
	// Save original
	origReader := randReader
	defer func() { randReader = origReader }()

	randReader = failingReader{}

	_, err := generateKey()
	if err == nil {
		t.Error("generateKey() should error with failing reader")
	}
}

func TestUpdateEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Test: no .env file
	err := updateEnvFile("base64:testkey")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Error("updateEnvFile() should error when .env doesn't exist")
	}

	// Test: .env without CRYPTO_KEY
	os.WriteFile(envPath, []byte("APP_NAME=Test\n"), 0644)
	err = updateEnvFile("base64:testkey")
	if err == nil || !strings.Contains(err.Error(), "CRYPTO_KEY not found") {
		t.Error("updateEnvFile() should error when CRYPTO_KEY not found")
	}

	// Test: successful update
	os.WriteFile(envPath, []byte("APP_NAME=Test\nCRYPTO_KEY=old\nPORT=4000\n"), 0644)
	err = updateEnvFile("base64:newkey")
	if err != nil {
		t.Fatalf("updateEnvFile() error = %v", err)
	}

	content, _ := os.ReadFile(envPath)
	if !strings.Contains(string(content), "CRYPTO_KEY=base64:newkey") {
		t.Error("updateEnvFile() did not update CRYPTO_KEY")
	}
	if !strings.Contains(string(content), "APP_NAME=Test") {
		t.Error("updateEnvFile() should preserve other values")
	}
}

func TestUpdateEnvFileReadError(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Create unreadable file
	os.WriteFile(envPath, []byte("test"), 0644)
	os.Chmod(envPath, 0000)
	defer os.Chmod(envPath, 0644)

	err := updateEnvFile("base64:testkey")
	if err == nil {
		t.Error("updateEnvFile() should error on unreadable file")
	}
}

func TestRunKeyGenerateShowOnly(t *testing.T) {
	origShowOnly := showOnly
	defer func() { showOnly = origShowOnly }()

	showOnly = true

	cmd := &cobra.Command{}
	runKeyGenerate(cmd, []string{})
}

func TestRunKeyGenerateUpdateEnv(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	os.WriteFile(envPath, []byte("CRYPTO_KEY=old\n"), 0644)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	origShowOnly := showOnly
	defer func() { showOnly = origShowOnly }()
	showOnly = false

	cmd := &cobra.Command{}
	runKeyGenerate(cmd, []string{})

	content, _ := os.ReadFile(envPath)
	if !strings.Contains(string(content), "CRYPTO_KEY=base64:") {
		t.Error("runKeyGenerate() should have set a base64 key")
	}
}

func TestRunKeyGenerateKeyError(t *testing.T) {
	// Mock exit function
	var exitCode int
	origExit := exitFunc
	exitFunc = func(code int) { exitCode = code }
	defer func() { exitFunc = origExit }()

	// Mock failing reader
	origReader := randReader
	randReader = failingReader{}
	defer func() { randReader = origReader }()

	cmd := &cobra.Command{}
	runKeyGenerate(cmd, []string{})

	if exitCode != 1 {
		t.Errorf("runKeyGenerate() exitCode = %d, want 1", exitCode)
	}
}

func TestRunKeyGenerateEnvError(t *testing.T) {
	tmpDir := t.TempDir()

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Mock exit function
	var exitCode int
	origExit := exitFunc
	exitFunc = func(code int) { exitCode = code }
	defer func() { exitFunc = origExit }()

	// Restore rand reader
	origReader := randReader
	randReader = rand.Reader
	defer func() { randReader = origReader }()

	origShowOnly := showOnly
	showOnly = false
	defer func() { showOnly = origShowOnly }()

	cmd := &cobra.Command{}
	runKeyGenerate(cmd, []string{}) // No .env file

	if exitCode != 1 {
		t.Errorf("runKeyGenerate() exitCode = %d, want 1", exitCode)
	}
}

func TestKeyCmd(t *testing.T) {
	if KeyCmd == nil {
		t.Fatal("KeyCmd should not be nil")
	}
	if KeyCmd.Use != "key:generate" {
		t.Errorf("KeyCmd.Use = %s, want key:generate", KeyCmd.Use)
	}
}
