package main

import (
	"regexp"
	"strings"
)

// convertMarkdownToMrkdwn converts Markdown format to Slack's Mrkdwn format
// See: https://api.slack.com/reference/surfaces/formatting
func convertMarkdownToMrkdwn(markdown string) string {
	text := markdown

	// Placeholder to protect already-converted bold text
	const boldPlaceholder = "\x00BOLD\x00"

	// Store bold conversions with placeholders
	boldMatches := []string{}

	// Convert bold: **text** or __text__ -> placeholder
	boldPattern := regexp.MustCompile(`\*\*(.+?)\*\*`)
	text = boldPattern.ReplaceAllStringFunc(text, func(match string) string {
		content := boldPattern.FindStringSubmatch(match)
		if len(content) > 1 {
			replacement := "*" + content[1] + "*"
			boldMatches = append(boldMatches, replacement)
			return boldPlaceholder
		}
		return match
	})

	boldUnderscorePattern := regexp.MustCompile(`__(.+?)__`)
	text = boldUnderscorePattern.ReplaceAllStringFunc(text, func(match string) string {
		content := boldUnderscorePattern.FindStringSubmatch(match)
		if len(content) > 1 {
			replacement := "*" + content[1] + "*"
			boldMatches = append(boldMatches, replacement)
			return boldPlaceholder
		}
		return match
	})

	// Convert italic: single *text* -> _text_ (only single asterisks)
	italicAsteriskPattern := regexp.MustCompile(`\*([^*\n]+?)\*`)
	text = italicAsteriskPattern.ReplaceAllString(text, `_${1}_`)

	// Convert italic: _text_ -> _text_ (already in Mrkdwn format, no change needed)

	// Restore bold text from placeholders
	for _, boldText := range boldMatches {
		text = strings.Replace(text, boldPlaceholder, boldText, 1)
	}

	// Convert strikethrough: ~~text~~ -> ~text~
	strikethroughPattern := regexp.MustCompile(`~~(.+?)~~`)
	text = strikethroughPattern.ReplaceAllString(text, `~${1}~`)

	// Convert links: [text](url) -> <url|text>
	linkPattern := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	text = linkPattern.ReplaceAllString(text, `<${2}|${1}>`)

	// Convert code blocks: ```lang\ncode\n``` -> ```code```
	// Remove language specifier from code blocks (supports alphanumeric, hyphens, plus, etc.)
	codeBlockPattern := regexp.MustCompile("```[a-zA-Z0-9+#\\-]*\n")
	text = codeBlockPattern.ReplaceAllString(text, "```\n")

	// Convert unordered lists: * item or - item -> • item
	listPattern := regexp.MustCompile(`(?m)^[\*\-]\s+`)
	text = listPattern.ReplaceAllString(text, "• ")

	// Convert ordered lists: 1. item -> 1. item (no change needed)

	return text
}
