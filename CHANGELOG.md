# Changelog

All notable changes to Velocity CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
