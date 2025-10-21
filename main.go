package main

import (
	"bufio"
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
)

const (
	keyringService = "slack-cli"
	keyringUser    = "SLACK_TOKEN"
)

var (
	api *slack.Client
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

	// Handle configure command first (doesn't need token)
	if args[0] == "configure" {
		return configure(ctx)
	}

	// Get token from keyring first, then fall back to env var
	token := ""
	keyringToken, err := keyring.Get(keyringService, keyringUser)
	if err == nil && keyringToken != "" {
		token = keyringToken
	} else {
		token = os.Getenv("SLACK_TOKEN")
	}

	if token == "" {
		return fmt.Errorf("SLACK_TOKEN must be set (use 'slack configure' or set SLACK_TOKEN env var)")
	}

	// disable HTTP/2 support as it causes issues with some proxies
	http.DefaultTransport.(*http.Transport).ForceAttemptHTTP2 = false
	api = slack.New(token)

	switch args[0] {
	case "send-message":
		if len(args) < 3 {
			return fmt.Errorf("usage: slack send-message <channel|email> <message>")
		}
		return sendMessage(ctx, args[1], args[2])
	default:
		return fmt.Errorf("unknown sub-command: %s", args[0])
	}
}

func sendMessage(ctx context.Context, identifier, body string) error {

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

func configure(ctx context.Context) error {
	fmt.Fprintln(os.Stderr, "Enter your Slack API token:")
	
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		return fmt.Errorf("no token provided")
	}

	token := strings.TrimSpace(scanner.Text())
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
