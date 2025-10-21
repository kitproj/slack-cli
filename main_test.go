package main

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestConfigure_EmptyToken(t *testing.T) {
	// Mock stdin with empty input
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write([]byte("\n"))
	w.Close()

	ctx := context.Background()
	err := configure(ctx)
	
	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}
	
	if !strings.Contains(err.Error(), "token cannot be empty") {
		t.Errorf("Expected 'token cannot be empty' error, got: %v", err)
	}
}

func TestConfigure_ValidToken(t *testing.T) {
	// Mock stdin with valid token
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r
	testToken := "xoxb-test-token-123"
	w.Write([]byte(testToken + "\n"))
	w.Close()

	// Capture stderr to suppress the output during tests
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()
	
	_, tmpFile, _ := os.Pipe()
	os.Stderr = tmpFile

	ctx := context.Background()
	err := configure(ctx)
	
	// We expect this to fail in the test environment due to keyring access,
	// but the token reading logic should work
	if err != nil && !strings.Contains(err.Error(), "keyring") {
		t.Errorf("Expected keyring error or nil, got: %v", err)
	}
}

func TestConfigure_WhitespaceToken(t *testing.T) {
	// Mock stdin with whitespace-only input
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write([]byte("   \n"))
	w.Close()

	ctx := context.Background()
	err := configure(ctx)
	
	if err == nil {
		t.Error("Expected error for whitespace-only token, got nil")
	}
	
	if !strings.Contains(err.Error(), "token cannot be empty") {
		t.Errorf("Expected 'token cannot be empty' error, got: %v", err)
	}
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
	// Set a dummy token to get past token check
	oldToken := token
	token = "test-token"
	defer func() { token = oldToken }()

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
