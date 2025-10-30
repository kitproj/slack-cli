package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/slack-go/slack"
)

// runMCPServer starts the MCP server that communicates over stdio using the mcp-go library
func runMCPServer(ctx context.Context) error {
	token := getToken()
	if token == "" {
		return fmt.Errorf("Slack token must be set (use 'slack configure' or set SLACK_TOKEN env var)")
	}

	// disable HTTP/2 support as it causes issues with some proxies
	http.DefaultTransport.(*http.Transport).ForceAttemptHTTP2 = false
	api := slack.New(token)

	// Create a new MCP server
	s := server.NewMCPServer(
		"slack-cli-mcp-server",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Define the send_message tool
	sendMessageTool := mcp.NewTool("send_message",
		mcp.WithDescription("Send a message to a Slack channel or user. You can specify either a channel ID or a user's email address. The message supports Markdown formatting which will be automatically converted to Slack's Mrkdwn format."),
		mcp.WithString("identifier",
			mcp.Required(),
			mcp.Description("The Slack channel ID (e.g., 'C1234567890') or user email address (e.g., 'user@example.com')"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The message to send. Supports Markdown formatting."),
		),
	)

	// Add the tool handler
	s.AddTool(sendMessageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		identifier, err := request.RequireString("identifier")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid 'identifier' argument: %v", err)), nil
		}

		message, err := request.RequireString("message")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Missing or invalid 'message' argument: %v", err)), nil
		}

		// Send the message using the existing sendMessage function
		err = sendMessage(ctx, api, identifier, message)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Message sent successfully to %s", identifier)), nil
	})

	// Start the stdio server
	return server.ServeStdio(s)
}
