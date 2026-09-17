# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Daemon Configuration**: Added `aax service install` to instantly configure and boot Aaxion as a systemd background daemon on Linux.
- **Dynamic Installer**: `install.sh` now uses progress bars, fetches dynamic version tags, and conditionally installs the systemd daemon if Linux is detected.
- **Admin Management**: Added `--force` flag to `aax create-admin` to allow easy overriding and resetting of the admin user and sessions.
- **Cross-Platform Builds**: Created `build.sh` script to cross-compile the binary natively for Windows, Linux, and macOS across multiple architectures.
- **System Information**: Added the `aax version` command to output cleanly formatted build, OS, and runtime statistics.

### Changed

- **Database Migration**: Relocated the `.aaxion.db` SQLite database from the local execution directory to a secure `~/.aaxion/` directory to prevent permission conflicts with background daemons.
- **Auth Middleware**: Refactored the token query parameter parsing (`?tkn=`) to cleanly iterate over an array of permitted endpoints (`thumbnail`, `stream`, etc.), resolving strict-matching 401 Unauthorized errors.
- **Website UI**: Added dynamic OS-detection `<InstallCommand>` toggles to the documentation frontend.
- **Storage Compatibility**: Split `syscall.Statfs` implementations using Go build tags (`_unix.go` and `_windows.go`) to prevent cross-compilation panics on Windows.

### Fixed

- Fixed Next.js hydration mismatches on the dashboard caused by browser extensions injecting into SVG and HTML elements.
- Fixed an authentication bypass loop where the middleware would throw 401 Unauthorized on thumbnail streams by strictly expecting short-lived access tokens instead of allowing primary session tokens.
