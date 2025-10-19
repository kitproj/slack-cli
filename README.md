# Slack CLI

A Slack CLI that allows you to send Slack messages. Inspired by the GitHub CLI, it aims to provide a simple and efficient way to interact with Slack from the command line, without the need to install a runtime such as Node.js or Python.

It's aimed at coding agents with a very simple interface, and is not intended to be a full-featured Slack client.

## Installation

Download the binary for your platform from the release page:

```bash
sudo curl -fsL -o /usr/local/bin/slack https://github.com/kitproj/slack-cli/releases/download/v0.0.6/slack_v0.0.6_linux_arm64
sudo chmod +x /usr/local/bin/slack
```


## Features

### Markdown to Mrkdwn Conversion

The CLI automatically converts standard Markdown formatting to Slack's Mrkdwn format:

- **Bold**: `**text**` or `__text__` → `*text*`
- **Italic**: `_text_` → `_text_`
- **Strikethrough**: `~~text~~` → `~text~`
- **Links**: `[text](url)` → `<url|text>`
- **Code**: `` `code` `` → `` `code` ``
- **Code blocks**: ` ```lang\ncode\n``` ` → ` ```\ncode\n``` `
- **Unordered lists**: `* item` or `- item` → `• item`

### Suggested Future Features

The Slack API supports many additional message features that could be added:

1. **Attachments** - Add rich attachments with fields, colors, and images
2. **Blocks** - Use Slack's Block Kit for interactive elements (buttons, select menus, etc.)
3. **Thread replies** - Reply to existing messages in threads
4. **Reactions** - Add emoji reactions to messages
5. **File uploads** - Attach files to messages
6. **Message updates** - Edit or delete previously sent messages
7. **User mentions** - Tag users with `@username` syntax
8. **Channel mentions** - Tag channels with `#channel` syntax
9. **Message formatting options** - Custom colors, timestamps, footers
10. **Scheduled messages** - Send messages at a specific time

## Prompt

Add this to your prompt (e.g. `AGENTS.md`):

```markdown
- You can send messages to a Slack user by using the `slack send-message <channel|email> "<message>"` command.
- The message supports Markdown formatting which will be automatically converted to Slack's Mrkdwn format.
```

## Usage

```bash
Usage:
  slack send-message <channel|email> <message> - send a message to a user

Options:
  -t string
    	Slack API token (defaults to SLACK_TOKEN env var) (default "")
```
