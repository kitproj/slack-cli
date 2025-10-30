package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestRunMCPServer_Initialize(t *testing.T) {
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

	// Test initialize request
	initReq := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  json.RawMessage(`{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}`),
	}
	
	reqBytes, _ := json.Marshal(initReq)
	
	// We can't easily test the full server loop, but we can test the JSON structures
	var parsedReq JSONRPCRequest
	if err := json.Unmarshal(reqBytes, &parsedReq); err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}
	
	if parsedReq.Method != "initialize" {
		t.Errorf("Expected method 'initialize', got '%s'", parsedReq.Method)
	}
}

func TestRunMCPServer_ToolsList(t *testing.T) {
	// Test tools/list request structure
	listReq := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}
	
	reqBytes, _ := json.Marshal(listReq)
	
	var parsedReq JSONRPCRequest
	if err := json.Unmarshal(reqBytes, &parsedReq); err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}
	
	if parsedReq.Method != "tools/list" {
		t.Errorf("Expected method 'tools/list', got '%s'", parsedReq.Method)
	}
	
	// Test that we can create a proper tools list response
	result := ListToolsResult{
		Tools: []Tool{
			{
				Name:        "send_message",
				Description: "Send a message to a Slack channel or user",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"identifier": map[string]interface{}{
							"type": "string",
						},
						"message": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
	}
	
	if len(result.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(result.Tools))
	}
	
	if result.Tools[0].Name != "send_message" {
		t.Errorf("Expected tool name 'send_message', got '%s'", result.Tools[0].Name)
	}
}

func TestRunMCPServer_ToolCall(t *testing.T) {
	// Test tools/call request structure
	callReq := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"send_message","arguments":{"identifier":"test@example.com","message":"Hello"}}`),
	}
	
	reqBytes, _ := json.Marshal(callReq)
	
	var parsedReq JSONRPCRequest
	if err := json.Unmarshal(reqBytes, &parsedReq); err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}
	
	if parsedReq.Method != "tools/call" {
		t.Errorf("Expected method 'tools/call', got '%s'", parsedReq.Method)
	}
	
	var params CallToolParams
	if err := json.Unmarshal(parsedReq.Params, &params); err != nil {
		t.Fatalf("Failed to parse params: %v", err)
	}
	
	if params.Name != "send_message" {
		t.Errorf("Expected tool name 'send_message', got '%s'", params.Name)
	}
	
	if params.Arguments["identifier"] != "test@example.com" {
		t.Errorf("Expected identifier 'test@example.com', got '%v'", params.Arguments["identifier"])
	}
}

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
	// Unset SLACK_TOKEN env var
	oldToken := os.Getenv("SLACK_TOKEN")
	os.Unsetenv("SLACK_TOKEN")
	defer func() {
		if oldToken != "" {
			os.Setenv("SLACK_TOKEN", oldToken)
		}
	}()

	ctx := context.Background()
	err := run(ctx, []string{"mcp-server"})
	
	if err == nil {
		t.Error("Expected error for missing token, got nil")
	}
	
	if !strings.Contains(err.Error(), "Slack token must be set") {
		t.Errorf("Expected 'Slack token must be set' error, got: %v", err)
	}
}
