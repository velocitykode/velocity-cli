# Changelog

All notable changes to Velocity CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.5] - 2025-12-30

### Changed
- refactor: remove dots pattern, simplify project creation output

## [0.3.4] - 2025-12-30

### Changed
- refactor: standardize CLI output with lipgloss ui package, fixes #6

## [0.3.3] - 2025-12-30

### Fixed
- detach dev server processes to survive CLI exit

## [0.3.2] - 2025-12-30

### Fixed
- clean up code comments

### Changed
- ci: use VELOCITY_TOKEN org secret
- ci: use PAT for tag push to trigger release workflow

## [0.3.1] - 2025-12-30

### Fixed
- trigger release
- release trigger
- trigger release workflow
- standardize default port to 4000, fixes #4

### Changed
- ci: use GITHUB_TOKEN instead of PAT for releases
- chore: trigger release

## [0.3.0] - 2025-12-29

### Added
- Add auto-release workflow with conventional commit changelog, fixes #2
- Add Docker-style spinner for project creation steps

## [0.2.1] - 2025-12-29

### Fixed
- Test output capture using `cmd.OutOrStdout()` instead of direct stdout
- Controller test expected filename (snake_case)
- Removed lint job from CI

## [0.2.0] - 2025-12-29

### Added
- Auto-install `air` for hot reloading during project setup
- Automated releases with GoReleaser and GitHub Actions

### Fixed
- Use `go mod edit` instead of `go get` to set framework version (fixes project creation on fresh machines)
- Fallback to `go run` if `air` not installed

## [0.1.2] - 2025-12-29

### Fixed
- Pin velocity framework to v0.0.3

## [0.1.1] - 2025-12-29

### Added
- Go version check on CLI startup - validates system's installed Go 1.25+ before running any command
- Clear error message with upgrade instructions when Go version is too old
- Error handling for when Go is not installed

### Fixed
- Cryptic "failed to build migration runner" error when Go version is incompatible - now shows friendly message immediately

## [0.1.0] - 2025-12-27

### Added
- Initial release
- `velocity new` - Create new Velocity projects
- `velocity init` - Initialize Velocity in existing projects
- `velocity serve` - Development server with hot reload
- `velocity build` - Production builds with cross-compilation
- `velocity migrate` - Run database migrations
- `velocity migrate:fresh` - Reset and re-run migrations
- `velocity make:controller` - Generate controllers
- `velocity key:generate` - Generate encryption keys
- `velocity config` - CLI configuration management
