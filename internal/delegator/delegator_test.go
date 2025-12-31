package delegator

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestShouldDelegate_GlobalCommands(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"new command", []string{"new", "myapp"}, false},
		{"init command", []string{"init"}, false},
		{"help flag", []string{"--help"}, false},
		{"version flag", []string{"--version"}, false},
		{"config command", []string{"config"}, false},
		{"empty args", []string{}, false},
		{"upgrade command", []string{"upgrade"}, false},
		{"self-update command", []string{"self-update"}, false},
		{"-h flag", []string{"-h"}, false},
		{"-v flag", []string{"-v"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldDelegate(tt.args)
			if got != tt.want {
				t.Errorf("ShouldDelegate(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestGlobalCommands(t *testing.T) {
	// Ensure all expected global commands are registered
	expectedCmds := []string{
		"new", "init", "upgrade", "self-update",
		"help", "--help", "-h",
		"version", "--version", "-v",
		"config",
	}

	for _, cmd := range expectedCmds {
		if !GlobalCommands[cmd] {
			t.Errorf("Expected %q to be a global command", cmd)
		}
	}
}

func TestGetProjectCLIVersion_NoGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	version := getProjectCLIVersion()
	if version != "" {
		t.Errorf("getProjectCLIVersion() = %q, want empty string when no go.mod", version)
	}
}

func TestGetProjectCLIVersion_WithGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create go.mod with velocity-cli dependency
	goMod := `module testproject

go 1.21

require (
	github.com/velocitykode/velocity-cli v0.5.0
	github.com/other/pkg v1.0.0
)
`
	os.WriteFile("go.mod", []byte(goMod), 0644)

	version := getProjectCLIVersion()
	if version != "v0.5.0" {
		t.Errorf("getProjectCLIVersion() = %q, want %q", version, "v0.5.0")
	}
}

func TestGetProjectCLIVersion_NoDependency(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create go.mod without velocity-cli
	goMod := `module testproject

go 1.21

require github.com/other/pkg v1.0.0
`
	os.WriteFile("go.mod", []byte(goMod), 0644)

	version := getProjectCLIVersion()
	if version != "" {
		t.Errorf("getProjectCLIVersion() = %q, want empty string", version)
	}
}

func TestNeedsRebuild_NoBinary(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "nonexistent")

	if !needsRebuild(binPath) {
		t.Error("needsRebuild() should return true when binary doesn't exist")
	}
}

func TestNeedsRebuild_NewerGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create old binary
	binPath := filepath.Join(tmpDir, "cli")
	os.WriteFile(binPath, []byte("binary"), 0755)

	// Set binary time to past
	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(binPath, oldTime, oldTime)

	// Create newer go.mod
	os.WriteFile("go.mod", []byte("module test"), 0644)

	if !needsRebuild(binPath) {
		t.Error("needsRebuild() should return true when go.mod is newer")
	}
}

func TestNeedsRebuild_OlderGoMod(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create go.mod first
	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.WriteFile("go.sum", []byte(""), 0644)

	// Set go.mod/go.sum time to past
	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes("go.mod", oldTime, oldTime)
	os.Chtimes("go.sum", oldTime, oldTime)

	// Create cmd/velocity directory (needed for the walk)
	os.MkdirAll("cmd/velocity", 0755)
	os.WriteFile("cmd/velocity/main.go", []byte("package main"), 0644)
	os.Chtimes("cmd/velocity/main.go", oldTime, oldTime)

	// Create newer binary
	binPath := filepath.Join(tmpDir, "cli")
	os.WriteFile(binPath, []byte("binary"), 0755)

	if needsRebuild(binPath) {
		t.Error("needsRebuild() should return false when binary is newer than all sources")
	}
}

func TestIsNewer(t *testing.T) {
	tmpDir := t.TempDir()

	// Create file
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	// Get file time
	info, _ := os.Stat(filePath)
	fileTime := info.ModTime()

	// Test with older time
	olderTime := fileTime.Add(-1 * time.Hour)
	if !isNewer(filePath, olderTime) {
		t.Error("isNewer() should return true when file is newer than time")
	}

	// Test with newer time
	newerTime := fileTime.Add(1 * time.Hour)
	if isNewer(filePath, newerTime) {
		t.Error("isNewer() should return false when file is older than time")
	}

	// Test with nonexistent file
	if isNewer(filepath.Join(tmpDir, "nonexistent"), olderTime) {
		t.Error("isNewer() should return false for nonexistent file")
	}
}

func TestCheckVersionMismatch_NoWarningWhenMatch(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create go.mod with matching version
	goMod := `module testproject

require github.com/velocitykode/velocity-cli v0.5.0
`
	os.WriteFile("go.mod", []byte(goMod), 0644)

	// Should not panic
	CheckVersionMismatch("v0.5.0")
}

func TestCheckVersionMismatch_WarningWhenDifferent(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create go.mod with different version
	goMod := `module testproject

require github.com/velocitykode/velocity-cli v0.4.0
`
	os.WriteFile("go.mod", []byte(goMod), 0644)

	// Should not panic (warning is printed but we can't easily capture it)
	CheckVersionMismatch("v0.5.0")
}

func TestNeedsRebuild_NewerGoFileInCmdVelocity(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create old binary
	binPath := filepath.Join(tmpDir, "cli")
	os.WriteFile(binPath, []byte("binary"), 0755)
	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(binPath, oldTime, oldTime)

	// Create old go.mod/go.sum
	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.WriteFile("go.sum", []byte(""), 0644)
	os.Chtimes("go.mod", oldTime, oldTime)
	os.Chtimes("go.sum", oldTime, oldTime)

	// Create cmd/velocity with NEWER .go file
	os.MkdirAll("cmd/velocity", 0755)
	os.WriteFile("cmd/velocity/main.go", []byte("package main"), 0644)
	// main.go is newer than binary (created just now)

	if !needsRebuild(binPath) {
		t.Error("needsRebuild() should return true when .go file in cmd/velocity is newer")
	}
}

func TestNeedsRebuild_NoCmdVelocityDir(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create binary and go.mod but NO cmd/velocity
	binPath := filepath.Join(tmpDir, "cli")
	os.WriteFile(binPath, []byte("binary"), 0755)
	os.WriteFile("go.mod", []byte("module test"), 0644)

	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(binPath, oldTime, oldTime)
	os.Chtimes("go.mod", oldTime, oldTime)

	if !needsRebuild(binPath) {
		t.Error("needsRebuild() should return true when cmd/velocity doesn't exist")
	}
}

func TestNeedsRebuild_NewerGoSum(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create old binary
	binPath := filepath.Join(tmpDir, "cli")
	os.WriteFile(binPath, []byte("binary"), 0755)
	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(binPath, oldTime, oldTime)

	// Create old go.mod but NEW go.sum
	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.Chtimes("go.mod", oldTime, oldTime)
	os.WriteFile("go.sum", []byte("checksums"), 0644)
	// go.sum is newer (just created)

	if !needsRebuild(binPath) {
		t.Error("needsRebuild() should return true when go.sum is newer")
	}
}

func TestShouldDelegate_InVelocityProject(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create a Velocity project marker
	os.WriteFile("go.mod", []byte("module myapp\n\nrequire github.com/velocitykode/velocity v1.0.0\n"), 0644)

	// Non-global command should delegate when in Velocity project
	if !ShouldDelegate([]string{"serve"}) {
		t.Error("ShouldDelegate() should return true for 'serve' in Velocity project")
	}

	if !ShouldDelegate([]string{"migrate"}) {
		t.Error("ShouldDelegate() should return true for 'migrate' in Velocity project")
	}

	if !ShouldDelegate([]string{"build"}) {
		t.Error("ShouldDelegate() should return true for 'build' in Velocity project")
	}
}

func TestShouldDelegate_NotInVelocityProject(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// No Velocity markers - just a plain Go project
	os.WriteFile("go.mod", []byte("module myapp\n\ngo 1.21\n"), 0644)

	// Should NOT delegate when not in Velocity project
	if ShouldDelegate([]string{"serve"}) {
		t.Error("ShouldDelegate() should return false when not in Velocity project")
	}
}

func TestNeedsRebuild_WalkError(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create old binary
	binPath := filepath.Join(tmpDir, "cli")
	os.WriteFile(binPath, []byte("binary"), 0755)
	oldTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(binPath, oldTime, oldTime)

	// Create old go.mod/go.sum
	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.WriteFile("go.sum", []byte(""), 0644)
	os.Chtimes("go.mod", oldTime, oldTime)
	os.Chtimes("go.sum", oldTime, oldTime)

	// Create cmd/velocity with a subdirectory
	os.MkdirAll("cmd/velocity/subdir", 0755)
	os.WriteFile("cmd/velocity/main.go", []byte("package main"), 0644)
	os.WriteFile("cmd/velocity/subdir/helper.go", []byte("package main"), 0644)
	os.Chtimes("cmd/velocity/main.go", oldTime, oldTime)
	os.Chtimes("cmd/velocity/subdir/helper.go", oldTime, oldTime)

	// Make subdirectory unreadable to cause walk error
	os.Chmod("cmd/velocity/subdir", 0000)
	defer os.Chmod("cmd/velocity/subdir", 0755)

	// Walk should error when trying to enter unreadable subdir
	if !needsRebuild(binPath) {
		t.Error("needsRebuild() should return true when walk encounters error")
	}
}

func TestDelegate_BuildFails(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create project without valid cmd/velocity
	os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644)
	os.MkdirAll("cmd/velocity", 0755)
	// Invalid Go code - will fail to build
	os.WriteFile("cmd/velocity/main.go", []byte("invalid go code"), 0644)

	err := Delegate([]string{"test"})
	// Should fail because build fails and go run also fails
	if err == nil {
		t.Error("Delegate() should error when build fails")
	}
}

func TestDelegate_BuildSucceeds(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create a valid project with cmd/velocity that exits immediately
	os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644)
	os.MkdirAll("cmd/velocity", 0755)
	os.MkdirAll(".velocity/bin", 0755)
	os.WriteFile("cmd/velocity/main.go", []byte(`package main

func main() {
	// Exit successfully
}
`), 0644)

	err := Delegate([]string{})
	// Should succeed - builds and runs
	if err != nil {
		t.Errorf("Delegate() error = %v", err)
	}
}

func TestRunWithGoRun_InvalidProject(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// No cmd/velocity - go run will fail
	os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644)

	err := runWithGoRun([]string{})
	if err == nil {
		t.Error("runWithGoRun() should error when cmd/velocity doesn't exist")
	}
}

func TestRunWithGoRun_ValidProject(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create valid cmd/velocity
	os.WriteFile("go.mod", []byte("module testproject\n\ngo 1.21\n"), 0644)
	os.MkdirAll("cmd/velocity", 0755)
	os.WriteFile("cmd/velocity/main.go", []byte(`package main

func main() {}
`), 0644)

	err := runWithGoRun([]string{})
	if err != nil {
		t.Errorf("runWithGoRun() error = %v", err)
	}
}
