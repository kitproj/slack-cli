# Slack CLI

A Slack CLI that allows you to send Slack messages. Inspired by the GitHub CLI, it aims to provide a simple and efficient way to interact with Slack from the command line, without the need to install a runtime such as Node.js or Python.

It's aimed at coding agents with a very simple interface, and is not intended to be a full-featured Slack client.

## Installation

Download the binary for your platform from the release page:

```bash
sudo curl -fsL -o /usr/local/bin/slack https://github.com/kitproj/slack-cli/releases/download/v0.0.3/slack_v0.0.3_linux_arm64
sudo chmod +x /usr/local/bin/slack
```


## Prompt

Add this to your prompt (e.g. `AGENTS.md`):

```markdown
- You can send messages to a Slack user by using the `slack send-message <channel|email> "<message>"` command.
```

## Usage

```bash
Usage:
  slack send-message <channel|email> <message> - send a message to a user

Options:
  -t string
    	Slack API token (defaults to SLACK_TOKEN env var) (default "")
```
