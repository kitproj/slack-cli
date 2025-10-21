package main

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestConfigure_EmptyToken(t *testing.T) {
	// Note: This test will fail in non-terminal environments (like CI)
	// because term.ReadPassword requires a real terminal
	t.Skip("Skipping test that requires terminal for password input")
}

func TestConfigure_ValidToken(t *testing.T) {
	// Note: This test will fail in non-terminal environments (like CI)
	// because term.ReadPassword requires a real terminal
	t.Skip("Skipping test that requires terminal for password input")
}

func TestConfigure_WhitespaceToken(t *testing.T) {
	// Note: This test will fail in non-terminal environments (like CI)
	// because term.ReadPassword requires a real terminal
	t.Skip("Skipping test that requires terminal for password input")
}

func TestRun_MissingSubCommand(t *testing.T) {
	ctx := context.Background()
	err := run(ctx, []string{})
	
	if err == nil {
		t.Error("Expected error for missing sub-command, got nil")
	}
	
	if !strings.Contains(err.Error(), "missing sub-command") {
		t.Errorf("Expected 'missing sub-command' error, got: %v", err)
	}
}

func TestRun_UnknownSubCommand(t *testing.T) {
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

	ctx := context.Background()
	err := run(ctx, []string{"unknown-command"})
	
	if err == nil {
		t.Error("Expected error for unknown sub-command, got nil")
	}
	
	if !strings.Contains(err.Error(), "unknown sub-command") {
		t.Errorf("Expected 'unknown sub-command' error, got: %v", err)
	}
}

func TestConvertMarkdownToMrkdwn_Integration(t *testing.T) {
	// This is a simple integration test to ensure the function is accessible
	result := convertMarkdownToMrkdwn("**bold**")
	if result != "*bold*" {
		t.Errorf("Expected '*bold*', got '%s'", result)
	}
}
