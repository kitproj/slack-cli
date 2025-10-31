# Slack CLI & MPC Server

A Slack CLI and MCP server that allows you to send Slack messages. Inspired by the GitHub CLI, it aims to provide a simple and efficient way for humans and AI to interact with Slack from the command line or via MPC.

Like `jq`, it is packaged in a tiny (10Mb) binary, without the need to install a runtime such as Node.js or Python, and securely stores your secret in the key-ring, rather than in plain-text.

It's aimed at coding agents with a very simple interface, and is not intended to be a full-featured Slack client.

## Installation

Download the binary for your platform from the release page:

```bash
sudo curl -fsL -o /usr/local/bin/slack https://github.com/kitproj/slack-cli/releases/download/v0.0.10/slack_v0.0.10_linux_arm64
sudo chmod +x /usr/local/bin/slack
```

## Configuration

### Getting Your Slack API Token

1. Visit https://api.slack.com/apps
2. Create a new app or select an existing one
3. Navigate to "OAuth & Permissions"
4. Add the following Bot Token Scopes:
   - `chat:write` - Send messages
   - `users:read.email` - Look up users by email
5. Install the app to your workspace
6. Copy the "Bot User OAuth Token" (starts with `xoxb-`)

### Configuring the Token

For security, the Slack token is stored in your system keyring (login keyring). Configure it once:

```bash
echo "xoxb-your-slack-token" | slack configure
```

Or configure it interactively (token input will be hidden):

```bash
slack configure
```

Alternatively, you can use the `SLACK_TOKEN` environment variable:

```bash
export SLACK_TOKEN="xoxb-your-slack-token"
```

**Note:** Using the keyring is more secure in multi-user systems as environment variables are visible in the process list.

## Usage

### Direct CLI Usage

```bash
Usage:
  slack configure                                   - configure Slack token (reads from stdin)
  slack send-message <channel|email> <message>      - send a message to a user
  slack mcp-server                                  - start MCP server (Model Context Protocol)
```

**Sending to a User by Email:**
```bash
slack send-message alex_collins@intuit.com "I love this tool! It makes Slack integration so easy."
```

**Sending to a Channel by ID:**
```bash
slack send-message C1234567890 "Hello team! 👋"
```

**Using Markdown Formatting:**
```bash
slack send-message alex_collins@intuit.com "**Bold**, *italic*, ~~strikethrough~~, [link](https://example.com)"
```

### Markdown Support

Messages automatically convert Markdown to Slack's Mrkdwn format. Supported features:

- **Bold**: `**text**` or `__text__` → `*text*`
- **Italic**: `*text*` → `_text_`
- **Strikethrough**: `~~text~~` → `~text~`
- **Inline code**: `` `code` `` (unchanged)
- **Links**: `[text](url)` → `<url|text>`
- **Code blocks**: ` ```language\ncode\n``` ` (language identifier removed)
- **Unordered lists**: `* item` or `- item` → `• item`
- **Ordered lists**: `1. item` (unchanged)

### Finding Channel IDs

To get a channel ID in Slack:
1. Right-click on the channel name
2. Select "Copy" → "Copy link"
3. The channel ID is the part after the last slash (e.g., `C1234567890`)

### MCP Server Mode

The MCP (Model Context Protocol) server allows AI assistants and other tools to interact with Slack. This enables seamless integration with AI coding assistants and other automation tools.

**Setup:**

1. First, configure your Slack token (stored securely in the system keyring):
   ```bash
   echo "xoxb-your-slack-token" | slack configure
   ```

2. Add the MCP server configuration to your MCP client:
   ```json
   {
     "mcpServers": {
       "slack": {
         "command": "slack",
         "args": ["mcp-server"]
       }
     }
   }
   ```

The server exposes a `send_message` tool that accepts:
- `identifier` - Slack channel ID (e.g., 'C1234567890') or user email address (e.g., 'user@example.com')
- `message` - The message to send (supports Markdown formatting)

**Example usage from an AI assistant:**
> "Slack alex_collins@intuit.com to say how much you like this tool."
