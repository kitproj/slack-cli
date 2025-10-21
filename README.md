# Slack CLI

A Slack CLI that allows you to send Slack messages. Inspired by the GitHub CLI, it aims to provide a simple and efficient way to interact with Slack from the command line, without the need to install a runtime such as Node.js or Python.

It's aimed at coding agents with a very simple interface, and is not intended to be a full-featured Slack client.

## Installation

Download the binary for your platform from the release page:

```bash
sudo curl -fsL -o /usr/local/bin/slack https://github.com/kitproj/slack-cli/releases/download/v0.0.6/slack_v0.0.6_linux_arm64
sudo chmod +x /usr/local/bin/slack
```

## Configuration

For security, the Slack token is stored in your system keyring (login keyring). Configure it once:

```bash
echo "xoxb-your-slack-token" | slack configure
```

Alternatively, you can use the `SLACK_TOKEN` environment variable:

```bash
export SLACK_TOKEN="xoxb-your-slack-token"
```

**Note:** Using the keyring is more secure in multi-user systems as environment variables are visible in the process list.


## Prompt

Add this to your prompt (e.g. `AGENTS.md`):

```markdown
- You can send messages to a Slack user by using the `slack send-message <channel|email> "<message>"` command.
- The message supports Markdown formatting which will be automatically converted to Slack's Mrkdwn format.
```

## Usage

```bash
Usage:
  slack configure                                   - configure Slack token (reads from stdin)
  slack send-message <channel|email> <message>      - send a message to a user
```
