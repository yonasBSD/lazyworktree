package utils

import "github.com/muesli/reflow/wrap"

// WrapANSIContent wraps terminal-styled content to the target width while
// preserving ANSI escape sequences.
func WrapANSIContent(content string, width int) string {
	if width <= 0 || content == "" {
		return content
	}
	return wrap.String(content, width)
}
