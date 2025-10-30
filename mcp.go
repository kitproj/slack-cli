package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

// MCP JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP Protocol structures
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    ServerCapabilities     `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

type ServerCapabilities struct {
	Tools map[string]interface{} `json:"tools,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type CallToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// runMCPServer starts the MCP server that communicates over stdio using JSON-RPC 2.0
func runMCPServer(ctx context.Context) error {
	token := getToken()
	if token == "" {
		return fmt.Errorf("Slack token must be set (use 'slack configure' or set SLACK_TOKEN env var)")
	}

	// disable HTTP/2 support as it causes issues with some proxies
	http.DefaultTransport.(*http.Transport).ForceAttemptHTTP2 = false
	api := slack.New(token)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read from stdin: %w", err)
		}

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			sendError(writer, nil, -32700, "Parse error", nil)
			continue
		}

		var response JSONRPCResponse
		response.JSONRPC = "2.0"
		response.ID = req.ID

		switch req.Method {
		case "initialize":
			var params InitializeParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				sendError(writer, req.ID, -32602, "Invalid params", nil)
				continue
			}

			result := InitializeResult{
				ProtocolVersion: "2024-11-05",
				Capabilities: ServerCapabilities{
					Tools: map[string]interface{}{},
				},
				ServerInfo: ServerInfo{
					Name:    "slack-cli-mcp-server",
					Version: "0.1.0",
				},
			}
			response.Result = result

		case "tools/list":
			result := ListToolsResult{
				Tools: []Tool{
					{
						Name:        "send_message",
						Description: "Send a message to a Slack channel or user. You can specify either a channel ID or a user's email address. The message supports Markdown formatting which will be automatically converted to Slack's Mrkdwn format.",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"identifier": map[string]interface{}{
									"type":        "string",
									"description": "The Slack channel ID (e.g., 'C1234567890') or user email address (e.g., 'user@example.com')",
								},
								"message": map[string]interface{}{
									"type":        "string",
									"description": "The message to send. Supports Markdown formatting.",
								},
							},
							"required": []string{"identifier", "message"},
						},
					},
				},
			}
			response.Result = result

		case "tools/call":
			var params CallToolParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				sendError(writer, req.ID, -32602, "Invalid params", nil)
				continue
			}

			switch params.Name {
			case "send_message":
				identifier, ok := params.Arguments["identifier"].(string)
				if !ok {
					sendError(writer, req.ID, -32602, "Missing or invalid 'identifier' argument", nil)
					continue
				}

				message, ok := params.Arguments["message"].(string)
				if !ok {
					sendError(writer, req.ID, -32602, "Missing or invalid 'message' argument", nil)
					continue
				}

				err := sendMessage(ctx, api, identifier, message)
				if err != nil {
					response.Result = CallToolResult{
						Content: []ContentItem{
							{
								Type: "text",
								Text: fmt.Sprintf("Error: %v", err),
							},
						},
						IsError: true,
					}
				} else {
					response.Result = CallToolResult{
						Content: []ContentItem{
							{
								Type: "text",
								Text: fmt.Sprintf("Message sent successfully to %s", identifier),
							},
						},
						IsError: false,
					}
				}

			default:
				sendError(writer, req.ID, -32601, fmt.Sprintf("Unknown tool: %s", params.Name), nil)
				continue
			}

		case "notifications/initialized":
			// Client sends this after receiving initialize response, we just ignore it
			continue

		default:
			sendError(writer, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method), nil)
			continue
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			sendError(writer, req.ID, -32603, "Internal error", nil)
			continue
		}

		if _, err := writer.Write(responseBytes); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
		if _, err := writer.Write([]byte("\n")); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush response: %w", err)
		}
	}
}

func sendError(writer *bufio.Writer, id interface{}, code int, message string, data interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	responseBytes, _ := json.Marshal(response)
	writer.Write(responseBytes)
	writer.Write([]byte("\n"))
	writer.Flush()
}
