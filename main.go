package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

	cursor := ""
	pageCount := 0

	for {
		pageCount++
		fmt.Printf("Fetching page %d...\n", pageCount)

		url := "https://slack.com/api/users.list?limit=200"
		if cursor != "" {
			url += "&cursor=" + cursor
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to fetch users: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("failed to fetch users: HTTP %d", resp.StatusCode)
		}

		var result usersList
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("failed to parse Slack response: %w", err)
		}

		if !result.Ok {
			errorMsg := result.Error
			if errorMsg == "" {
				errorMsg = "unknown error"
			}
			return fmt.Errorf("slack API error (page %d): %s", pageCount, errorMsg)
		}

		for _, member := range result.Members {
			if !member.Deleted && !member.IsBot && member.Profile.Email != "" {
				fmt.Fprintf(file, "%s=%s\n", member.Name, member.Profile.Email)
			}
		}

		cursor = result.ResponseMetadata.NextCursor
		if cursor == "" {
			break
		}

		time.Sleep(4 * time.Second)
	}

	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
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

	url := fmt.Sprintf("https://slack.com/api/users.lookupByEmail?email=%s", email)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to lookup user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to lookup user: HTTP %d", resp.StatusCode)
	}

	var lookupResult lookupResult

	if err := json.NewDecoder(resp.Body).Decode(&lookupResult); err != nil {
		return fmt.Errorf("failed to parse Slack response: %w", err)
	}

	if lookupResult.User.ID == "" {
		return fmt.Errorf("user '%s' not found or missing email", username)
	}

	userID := lookupResult.User.ID

	payload := message{
		Channel: userID,
		Text:    body,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err = http.NewRequestWithContext(ctx, "POST", "https://slack.com/api/chat.postMessage", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to send message: HTTP %d", resp.StatusCode)
	}

	var sendResult sendResult

	if err := json.NewDecoder(resp.Body).Decode(&sendResult); err != nil {
		return fmt.Errorf("failed to parse Slack response: %w", err)
	}

	if !sendResult.Ok {
		errorMsg := sendResult.Error
		if errorMsg == "" {
			errorMsg = "unknown error"
		}
		return fmt.Errorf("failed to send message: %s", errorMsg)
	}

	fmt.Printf("Message sent to %s\n", username)
	return nil
}
