## Why

No backend exists yet. Need to establish the base server framework, configuration management, and test foundation before any API endpoints can be built.

## What Changes

- Initialize Go module (`github.com/notepass/vmprov-web`)
- Set up Echo web framework as the base server
- Config file (`config.yaml`) + env var support for: server port, DB connect string, DB username, DB password
- Makefile for build/test/run workflow
- Testify as the test framework for all backend components

## Capabilities

### New Capabilities
- `server-bootstrap`: Echo server setup, config loading, project structure
- `config-management`: YAML + env var configuration with viper
- `test-foundation`: testing + testify setup for all backend packages

### Modified Capabilities
None.

## Impact

- New `go.mod`, project layout, dependencies (Echo, viper, slog, testify)
- Establishes the foundation for all future backend work
