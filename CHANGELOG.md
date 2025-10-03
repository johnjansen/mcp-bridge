# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- Install script now properly detects and uses directories from user's PATH
- Install script is OS-aware (macOS vs Linux) with platform-specific directory preferences
- Install script correctly handles missing installation directories by creating them
- Fixed extracted binary filename detection in install script

## [0.1.0] - 2025-10-03

### Added
- Initial release of MCP Bridge
- Core MCP bridging functionality between stdio and bidirectional HTTP streaming
- Streaming transport support with automatic fallback to HTTP POST
- Debug logging for both client and server-side messages
- Full MCP protocol support via official Go SDK
- Command line flags for server URL, API key, and debug options
- Support for both streaming and HTTP POST transport modes
- Automatic transport negotiation with fallback
- Official MCP Go SDK integration
- BDD testing framework with godog
- Comprehensive documentation in README.md
- Pre-commit hooks for code quality
- GitHub Actions CI workflow

### Changed
- Updated transport documentation to accurately reflect MCP protocol implementation
- Clarified transport terminology to consistently use bidirectional HTTP streaming when referring to the transport layer

[0.1.0]: https://github.com/johnjansen/mcp-bridge/releases/tag/v0.1.0
