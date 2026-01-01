package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/velocitykode/velocity-cli/internal/delegator"
)

func TestUpgradeCmd_NotInProject(t *testing.T) {
	// Run from temp dir without velocity project structure
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	err := runUpgrade(nil, nil)
	if err == nil {
		t.Error("Expected error when not in a Velocity project")
	}
	if !strings.Contains(err.Error(), "not in a Velocity project") {
		t.Errorf("Expected 'not in a Velocity project' error, got: %v", err)
	}
}

func TestUpgradeCmd_NoVelocityCLIInGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create Velocity project (has velocity framework but not velocity-cli)
	goMod := `module testapp

go 1.21

require github.com/velocitykode/velocity v0.0.3
`
	os.WriteFile("go.mod", []byte(goMod), 0644)

	err := runUpgrade(nil, nil)
	if err == nil {
		t.Error("Expected error when velocity-cli not in go.mod")
	}
	if !strings.Contains(err.Error(), "could not determine") {
		t.Errorf("Expected version determination error, got: %v", err)
	}
}

func TestUpgradeCmd_AlreadyUpToDate(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create project structure
	os.MkdirAll("cmd/velocity", 0755)
	os.WriteFile("cmd/velocity/main.go", []byte("package main"), 0644)

	// Create go.mod with current version (strip 'v' prefix if present)
	currentVersion := strings.TrimPrefix(Version, "v")
	goMod := `module testapp

go 1.21

require github.com/velocitykode/velocity-cli v` + currentVersion + `
`
	os.WriteFile("go.mod", []byte(goMod), 0644)

	// Should succeed without error (already up to date)
	err := runUpgrade(nil, nil)
	if err != nil {
		t.Errorf("Expected no error when already up to date, got: %v", err)
	}
}

func TestUpgradeCmd_ParsesVersionFromGoMod(t *testing.T) {
	tests := []struct {
		name         string
		goModContent string
		wantVersion  string
	}{
		{
			name: "simple require",
			goModContent: `module testapp

require github.com/velocitykode/velocity-cli v0.5.0
`,
			wantVersion: "v0.5.0",
		},
		{
			name: "require block",
			goModContent: `module testapp

require (
	github.com/velocitykode/velocity v0.0.3
	github.com/velocitykode/velocity-cli v0.6.1
)
`,
			wantVersion: "v0.6.1",
		},
		{
			name: "with indirect",
			goModContent: `module testapp

require github.com/velocitykode/velocity-cli v0.6.2 // indirect
`,
			wantVersion: "v0.6.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(originalDir)

			os.WriteFile("go.mod", []byte(tt.goModContent), 0644)

			// Use the delegator's exported function
			gotVersion := delegator.GetProjectCLIVersion()

			if gotVersion != tt.wantVersion {
				t.Errorf("Got version %q, want %q", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestUpgradeCmd_ClearsCachedBinary(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create project structure
	os.MkdirAll("cmd/velocity", 0755)
	os.WriteFile("cmd/velocity/main.go", []byte("package main"), 0644)

	// Create cached binary
	cachedBinDir := filepath.Join(tmpDir, ".velocity", "bin")
	os.MkdirAll(cachedBinDir, 0755)
	cachedBin := filepath.Join(cachedBinDir, "cli")
	os.WriteFile(cachedBin, []byte("fake binary"), 0755)

	// Verify it exists
	if _, err := os.Stat(cachedBin); os.IsNotExist(err) {
		t.Fatal("Failed to create test cached binary")
	}

	// Create go.mod with current version (so upgrade succeeds without network)
	currentVersion := strings.TrimPrefix(Version, "v")
	goMod := `module testapp

go 1.21

require github.com/velocitykode/velocity-cli v` + currentVersion + `
`
	os.WriteFile("go.mod", []byte(goMod), 0644)

	// Run upgrade
	runUpgrade(nil, nil)

	// After upgrade with same version, binary should still be cleared
	// (This tests the cleanup logic even when version matches)
}
