# GitHub Copilot Instructions for slack-cli

## Project Overview

This is a Slack CLI tool written in Go that enables users to send Slack messages from the command line. It's designed to be simple and efficient, particularly for coding agents, without requiring runtimes like Node.js or Python.

## Project Structure

- `main.go` - Main entry point and core functionality
- `go.mod` - Go module dependencies
- `.github/workflows/` - CI/CD workflows for releases

## Development Setup

### Prerequisites
- Go 1.24.4 or later
- A Slack API token (set via `SLACK_TOKEN` environment variable)

### Building
```bash
go build -o slack .
```

### Testing
```bash
# Build and run a test
go build -o /tmp/slack .
/tmp/slack send-message <channel|email> "test message"
```

## Code Style and Conventions

### Go Best Practices
- Follow standard Go formatting (use `gofmt`)
- Use meaningful variable and function names
- Keep functions small and focused
- Handle errors explicitly
- Add context to errors using `fmt.Errorf` with `%w`

### Project-Specific Guidelines
- The CLI is designed to be minimal and focused
- Use the `slack-go/slack` library for Slack API interactions
- Context should be passed through for cancellation support
- HTTP/2 is disabled intentionally due to proxy compatibility issues

## Key Features

### Current Functionality
- `send-message <channel|email> <message>` - Send a message to a Slack channel or user by email

### Design Principles
1. **Simplicity** - Keep the interface minimal and easy to use
2. **Agent-friendly** - Optimized for use by coding agents
3. **Zero runtime dependencies** - Compiled binary with no external runtime requirements
4. **Graceful shutdown** - Support for signal handling (SIGTERM, SIGINT)

## Testing Guidelines

- Test the CLI with both channel IDs and email addresses
- Verify error handling for missing tokens
- Test signal handling for graceful shutdown
- Ensure proper error messages are displayed to users

## Dependencies

- `github.com/slack-go/slack` - Official Slack Go client library
- `github.com/gorilla/websocket` - Required by slack-go/slack

## Common Patterns

### Adding New Commands
When adding new sub-commands:
1. Add a new case in the `run()` function switch statement
2. Create a dedicated function for the command logic
3. Accept `context.Context` as the first parameter
4. Return errors with descriptive messages
5. Update the `flag.Usage` function with new command documentation

### Error Handling
```go
if err != nil {
    return fmt.Errorf("descriptive message: %w", err)
}
```

## Release Process

Releases are automated via GitHub Actions (`.github/workflows/release.yml`). The workflow builds binaries for multiple platforms when a new tag is pushed.

## Notes for Copilot

- This is a production tool used by coding agents
- Maintain backward compatibility in the CLI interface
- Keep the codebase minimal and focused
- Error messages should be clear and actionable
- The tool must work reliably in automated environments
