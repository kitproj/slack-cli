package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/slack-go/slack"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const (
	keyringService = "slack-cli"
	keyringUser    = "SLACK_TOKEN"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage:")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  slack configure                                   - configure Slack token (reads from stdin)")
		fmt.Fprintln(w, "  slack send-message <channel|email> <message>      - send a message to a user")
		fmt.Fprintln(w, "  slack mcp-server                                  - start MCP server (Model Context Protocol)")
		fmt.Fprintln(w)
	}
	flag.Parse()

	if err := run(ctx, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing sub-command")
	}

	switch args[0] {
	case "configure":
		return configureToken(ctx)
	case "mcp-server":
		return runMCPServer(ctx)
	case "send-message":
		if len(args) < 3 {
			return fmt.Errorf("usage: slack send-message <channel|email> <message>")
		}
		
		token := getToken()
		if token == "" {
			return fmt.Errorf("Slack token must be set (use 'slack configure' or set SLACK_TOKEN env var)")
		}

		// disable HTTP/2 support as it causes issues with some proxies
		http.DefaultTransport.(*http.Transport).ForceAttemptHTTP2 = false
		api := slack.New(token)
		
		return sendMessage(ctx, api, args[1], args[2])
	default:
		return fmt.Errorf("unknown sub-command: %s", args[0])
	}
}

func getToken() string {
	// Get token from env var first, then fall back to keyring
	if token := os.Getenv("SLACK_TOKEN"); token != "" {
		return token
	}
	
	keyringToken, err := keyring.Get(keyringService, keyringUser)
	if err == nil && keyringToken != "" {
		return keyringToken
	}
	
	return ""
}

func sendMessage(ctx context.Context, api *slack.Client, identifier, body string) error {
	var channel string
	if strings.Contains(identifier, "@") {
		user, err := api.GetUserByEmailContext(ctx, identifier)
		if err != nil {
			return fmt.Errorf("failed to lookup user: %w", err)
		}
		channel = user.ID
	} else {
		channel = identifier
	}

	// Convert Markdown to Mrkdwn format
	mrkdwnBody := convertMarkdownToMrkdwn(body)

	if _, _, err := api.PostMessageContext(ctx, channel, slack.MsgOptionText(mrkdwnBody, false)); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("Message sent to %s (%s)\n", identifier, channel)
	return nil
}

func configureToken(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "To get your Slack API token, visit: https://api.slack.com/apps")
	fmt.Fprintln(os.Stderr, "Create an app, install it to your workspace, and copy the Bot User OAuth Token")
	fmt.Fprint(os.Stderr, "Enter your Slack API token: ")
	
	var token string
	
	// Check if stdin is a terminal
	if term.IsTerminal(int(os.Stdin.Fd())) {
		// Read password without echoing to terminal
		tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr) // Print newline after password input
		
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		token = strings.TrimSpace(string(tokenBytes))
	} else {
		// If not a terminal (e.g., piped input), read normally
		var line string
		if _, err := fmt.Fscanln(os.Stdin, &line); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		token = strings.TrimSpace(line)
		fmt.Fprintln(os.Stderr) // Print newline for consistency
	}

	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Store the token in the keyring
	if err := keyring.Set(keyringService, keyringUser, token); err != nil {
		return fmt.Errorf("failed to store token in keyring: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Token successfully stored in keyring")
	return nil
}
