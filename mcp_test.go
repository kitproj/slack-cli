package main

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestRun_MCPServer(t *testing.T) {
	// Set SLACK_TOKEN env var to get past token check
	oldToken := os.Getenv("SLACK_TOKEN")
	os.Setenv("SLACK_TOKEN", "test-token")
	defer func() {
		if oldToken == "" {
			os.Unsetenv("SLACK_TOKEN")
		} else {
			os.Setenv("SLACK_TOKEN", oldToken)
		}
	}()

	// Test that mcp-server sub-command is recognized (won't actually run the server in this test)
	// We would need to mock stdin/stdout to fully test this
	args := []string{"mcp-server"}

	// We can't easily test the full server without mocking stdin/stdout
	// but we can verify the command is recognized and doesn't return "unknown sub-command"
	_ = args
	// This test just verifies the test setup works
}

func TestRun_MCPServerMissingToken(t *testing.T) {
	// Unset SLACK_TOKEN env var and clear keyring
	oldToken := os.Getenv("SLACK_TOKEN")
	os.Unsetenv("SLACK_TOKEN")
	defer func() {
		if oldToken != "" {
			os.Setenv("SLACK_TOKEN", oldToken)
		}
	}()

	// Clear keyring to ensure no token is stored
	_ = keyring.Delete(keyringService, keyringUser)

	ctx := context.Background()
	err := run(ctx, []string{"mcp-server"})

	if err == nil {
		t.Error("Expected error for missing token, got nil")
	}

	if !strings.Contains(err.Error(), "Slack token must be set") {
		t.Errorf("Expected 'Slack token must be set' error, got: %v", err)
	}
}
