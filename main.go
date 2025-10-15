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
)

var (
	token string
	api   *slack.Client
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	flag.StringVar(&token, "t", os.Getenv("SLACK_TOKEN"), "Slack API token (defaults to SLACK_TOKEN env var)")
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage:")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  slack send-message <channel|email> <message> - send a message to a user")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options:")
		flag.PrintDefaults()
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

	if token == "" {
		return fmt.Errorf("SLACK_TOKEN must be set")
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

	if _, _, err := api.PostMessageContext(ctx, channel, slack.MsgOptionText(body, false)); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("Message sent to %s (%s)\n", identifier, channel)
	return nil
}
