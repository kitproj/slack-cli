package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/slack-go/slack"
)

var (
	token          string
	userEmailsFile string
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	flag.StringVar(&token, "t", os.Getenv("SLACK_TOKEN"), "Slack API token (defaults to SLACK_TOKEN env var)")
	flag.StringVar(&userEmailsFile, "c", "/var/local/slack/user_emails", "File to cache Slack user emails")
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage:")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  slack send-message <username> <message> - send a message to a user")
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

	switch args[0] {
	case "cache-users":
		return cacheUsers(ctx)
	case "send-message":
		if len(args) < 3 {
			return fmt.Errorf("usage: slack send-message <username> <message>")
		}
		return sendMessage(ctx, args[1], args[2])
	default:
		return fmt.Errorf("unknown sub-command: %s", args[0])
	}
}

func cacheUsers(ctx context.Context) error {
	fmt.Println("Caching Slack users...")
	fmt.Println("Starting user list fetch from Slack API...")

	file, err := os.Create(userEmailsFile)
	if err != nil {
		return fmt.Errorf("failed to create user emails cache: %w", err)
	}
	defer file.Close()

	api := slack.New(token)

	users, err := api.GetUsersContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	count := 0
	for _, user := range users {
		if !user.Deleted && !user.IsBot && user.Profile.Email != "" {
			fmt.Fprintf(file, "%s=%s\n", user.Name, user.Profile.Email)
			count++
		}
	}

	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	fmt.Fprintf(os.Stderr, "\n✅ Success! Exported %d email addresses to '%s'\n", count, userEmailsFile)

	return nil
}

func sendMessage(ctx context.Context, username, body string) error {
	file, err := os.Open(userEmailsFile)
	if err != nil {
		return fmt.Errorf("failed to open user emails cache: %w", err)
	}
	defer file.Close()

	var email string
	scanner := bufio.NewScanner(file)
	prefix := username + "="
	for scanner.Scan() {
		line := scanner.Text()
		if suffix, found := strings.CutPrefix(line, prefix); found {
			email = suffix
			break
		}
	}

	if email == "" {
		return fmt.Errorf("user '%s' not found in cache, please run 'slack cache-users' first", username)
	}

	api := slack.New(token)

	user, err := api.GetUserByEmailContext(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to lookup user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user '%s' not found or missing email", username)
	}

	_, _, err = api.PostMessageContext(ctx, user.ID, slack.MsgOptionText(body, false))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("Message sent to %s\n", username)
	return nil
}
