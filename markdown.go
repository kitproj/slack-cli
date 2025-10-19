package main

import (
	"regexp"
)

// convertMarkdownToMrkdwn converts Markdown format to Slack's Mrkdwn format
// See: https://api.slack.com/reference/surfaces/formatting
func convertMarkdownToMrkdwn(markdown string) string {
	text := markdown

	// Convert bold: **text** or __text__ -> *text*
	boldPattern := regexp.MustCompile(`\*\*(.+?)\*\*`)
	text = boldPattern.ReplaceAllString(text, `*$1*`)
	
	boldUnderscorePattern := regexp.MustCompile(`__(.+?)__`)
	text = boldUnderscorePattern.ReplaceAllString(text, `*$1*`)

	// Convert italic: single * or _ -> _ (but not if part of bold)
	// We only convert underscores for italic since asterisks are used for bold in Mrkdwn
	italicUnderscorePattern := regexp.MustCompile(`\b_([^_]+?)_\b`)
	text = italicUnderscorePattern.ReplaceAllString(text, `_$1_`)

	// Convert strikethrough: ~~text~~ -> ~text~
	strikethroughPattern := regexp.MustCompile(`~~(.+?)~~`)
	text = strikethroughPattern.ReplaceAllString(text, `~$1~`)

	// Convert links: [text](url) -> <url|text>
	linkPattern := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	text = linkPattern.ReplaceAllString(text, `<$2|$1>`)

	// Convert code blocks: ```lang\ncode\n``` -> ```code```
	// Remove language specifier from code blocks
	codeBlockPattern := regexp.MustCompile("```[a-zA-Z]*\n")
	text = codeBlockPattern.ReplaceAllString(text, "```\n")

	// Convert unordered lists: * item or - item -> • item
	listPattern := regexp.MustCompile(`(?m)^[\*\-]\s+`)
	text = listPattern.ReplaceAllString(text, "• ")

	// Convert ordered lists: 1. item -> 1. item (no change needed)

	return text
}
