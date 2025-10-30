# Slack CLI

A Slack CLI that allows you to send Slack messages. Inspired by the GitHub CLI, it aims to provide a simple and efficient way to interact with Slack from the command line, without the need to install a runtime such as Node.js or Python.

It's aimed at coding agents with a very simple interface, and is not intended to be a full-featured Slack client.

## Installation

Download the binary for your platform from the release page:

```bash
sudo curl -fsL -o /usr/local/bin/slack https://github.com/kitproj/slack-cli/releases/download/v0.0.10/slack_v0.0.10_linux_arm64
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
- For AI assistants supporting MCP (Model Context Protocol), you can use `slack mcp-server` to enable tool-based Slack integration.
```

## Usage

### Direct CLI Usage

```bash
Usage:
  slack configure                                   - configure Slack token (reads from stdin)
  slack send-message <channel|email> <message>      - send a message to a user
  slack mcp-server                                  - start MCP server (Model Context Protocol)
```

**Example:**
```bash
slack send-message alex_collins@intuit.com "I love this tool! It makes Slack integration so easy."
```

### MCP Server Mode

The MCP (Model Context Protocol) server allows AI assistants and other tools to interact with Slack through a standardized JSON-RPC protocol over stdio. This enables seamless integration with AI coding assistants and other automation tools.

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
