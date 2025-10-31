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
- You can reply to a message in a thread by adding the thread timestamp as a third parameter: `slack send-message <channel|email> "<message>" <thread-ts>`.
- The message supports Markdown formatting which will be automatically converted to Slack's Mrkdwn format.
- For AI assistants supporting MCP (Model Context Protocol), you can use `slack mcp-server` to enable tool-based Slack integration.
```

## Usage

### Direct CLI Usage

```bash
Usage:
  slack configure                                            - configure Slack token (reads from stdin)
  slack send-message <channel|email> <message> [thread-ts]   - send a message (optionally reply to a thread)
  slack mcp-server                                           - start MCP server (Model Context Protocol)
```

**Examples:**
```bash
# Send a message
slack send-message alex_collins@intuit.com "I love this tool! It makes Slack integration so easy."
# Output includes:
# Message sent to alex_collins@intuit.com (U12345678)
# thread-ts: 1234567890.123456

# Reply to a message in a thread (use the thread-ts from the previous message)
slack send-message alex_collins@intuit.com "Thanks for the feedback!" "1234567890.123456"
# Output includes:
# Reply sent to alex_collins@intuit.com (U12345678) in thread 1234567890.123456
# thread-ts: 1234567890.654321
```

The `thread-ts` is printed after each message is sent, allowing you to use it to continue the conversation in a thread.

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

The server exposes the `send_message` tool with the following parameters:
- `identifier` - Slack channel ID (e.g., 'C1234567890') or user email address (e.g., 'user@example.com')
- `message` - The message to send (supports Markdown formatting)
- `thread_ts` - Optional: The thread timestamp of the parent message to reply to (e.g., '1234567890.123456'). When provided, the message will be sent as a threaded reply.

**Example usage from an AI assistant:**
> "Slack alex_collins@intuit.com to say how much you like this tool."
> "Reply to that Slack message with a thumbs up emoji."
