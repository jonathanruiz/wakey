# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Run the application
go run main.go

# Build the application
go build -o wakey main.go

# Run all tests
go test ./tests/...

# Run a specific test
go test ./tests/ -run TestReadConfig

# Install dependencies
go mod download
```

## Architecture

This is a Go TUI application for Wake-on-LAN device management, built with the Bubble Tea framework.

### Core Structure

- `main.go` - Entry point; initializes config and starts the Bubble Tea program
- `internal/view.go` - Root model that manages view switching between Devices and Groups views
- `internal/config/` - JSON config management for devices/groups stored at `~/.wakey_config.json`
- `internal/devices/` - Devices list view and device form for create/edit
- `internal/groups/` - Groups list view and group form for create/edit
- `internal/common/` - Shared components: keybindings, styles, popup dialogs, WoL implementation

### Bubble Tea Pattern

Each view follows the Bubble Tea Model-View-Update pattern:
- Models implement `tea.Model` interface (`Init`, `Update`, `View`)
- Views return `tea.Model` to switch between screens (e.g., list to form and back)
- Keybindings are defined in `keybindings.go` files within each package

### Key Packages

- `internal/common/wol/` - Wake-on-LAN magic packet implementation and device ping checks
- `internal/common/style/` - Lipgloss styles and terminal width calculations
- `internal/common/popup/` - Confirmation dialog component
- `internal/common/status/` - Global status message state

### Config Structure

Devices and groups are stored in `~/.wakey_config.json`:
- Devices have: ID, DeviceName, Description, MacAddress, IPAddress, State
- Groups reference device IDs in their Devices array
