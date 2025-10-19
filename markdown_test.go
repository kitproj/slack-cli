package main

import (
	"testing"
)

func TestConvertMarkdownToMrkdwn(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "bold with double asterisks",
			markdown: "This is **bold** text",
			expected: "This is *bold* text",
		},
		{
			name:     "bold with double underscores",
			markdown: "This is __bold__ text",
			expected: "This is *bold* text",
		},
		{
			name:     "strikethrough",
			markdown: "This is ~~strikethrough~~ text",
			expected: "This is ~strikethrough~ text",
		},
		{
			name:     "inline code",
			markdown: "This is `code` text",
			expected: "This is `code` text",
		},
		{
			name:     "link",
			markdown: "Check out [Google](https://google.com)",
			expected: "Check out <https://google.com|Google>",
		},
		{
			name:     "code block with language",
			markdown: "```python\nprint('hello')\n```",
			expected: "```\nprint('hello')\n```",
		},
		{
			name:     "code block without language",
			markdown: "```\ncode here\n```",
			expected: "```\ncode here\n```",
		},
		{
			name:     "unordered list with asterisk",
			markdown: "* Item 1\n* Item 2",
			expected: "• Item 1\n• Item 2",
		},
		{
			name:     "unordered list with dash",
			markdown: "- Item 1\n- Item 2",
			expected: "• Item 1\n• Item 2",
		},
		{
			name:     "ordered list",
			markdown: "1. First\n2. Second",
			expected: "1. First\n2. Second",
		},
		{
			name:     "mixed formatting",
			markdown: "This is **bold** and ~~strike~~ with a [link](https://example.com)",
			expected: "This is *bold* and ~strike~ with a <https://example.com|link>",
		},
		{
			name:     "italic with single asterisk",
			markdown: "This is *italic* text",
			expected: "This is _italic_ text",
		},
		{
			name:     "code block with complex language specifier",
			markdown: "```c++\nint main() {}\n```",
			expected: "```\nint main() {}\n```",
		},
		{
			name:     "code block with c# language",
			markdown: "```c#\npublic void Main() {}\n```",
			expected: "```\npublic void Main() {}\n```",
		},
		{
			name:     "code block with hyphenated language",
			markdown: "```objective-c\n@interface MyClass\n```",
			expected: "```\n@interface MyClass\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertMarkdownToMrkdwn(tt.markdown)
			if result != tt.expected {
				t.Errorf("convertMarkdownToMrkdwn() = %q, want %q", result, tt.expected)
			}
		})
	}
}
