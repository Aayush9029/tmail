package ui

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	"golang.org/x/term"
)

// ANSI color constants matching the bash tool palette
const (
	Green  = "\033[1;32m"
	Red    = "\033[1;31m"
	Yellow = "\033[1;33m"
	Cyan   = "\033[1;36m"
	Blue   = "\033[1;34m"
	Dim    = "\033[2m"
	Bold   = "\033[1m"
	Reset  = "\033[0m"
)

// IsTTY returns true if stdout is a terminal
func IsTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// TermWidth returns the terminal width, or 80 as default.
func TermWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

func Header(name string) {
	if IsTTY() {
		fmt.Printf("%s%s⚡ %s%s\n", Cyan, Bold, name, Reset)
	}
}

func Success(msg string) {
	fmt.Printf("%s✓%s %s\n", Green, Reset, msg)
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s✗%s %s\n", Red, Reset, msg)
}

func Status(msg string) {
	fmt.Printf("  %s→%s %s\n", Green, Reset, msg)
}

func Dimf(format string, a ...any) {
	fmt.Printf("  %s"+format+"%s\n", append([]any{Dim}, append(a, Reset)...)...)
}

func Waiting(msg string) {
	fmt.Printf("%s⏳%s %s\n", Yellow, Reset, msg)
}

// Fatalf prints an error and exits
func Fatalf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
	os.Exit(1)
}

// HTML-to-text conversion

var (
	reBr      = regexp.MustCompile(`(?i)<br\s*/?>`)
	reP       = regexp.MustCompile(`(?i)</p>`)
	reA       = regexp.MustCompile(`(?i)<a\s[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`)
	rePreCode = regexp.MustCompile(`(?is)<pre><code>(.*?)</code></pre>`)
	reTag     = regexp.MustCompile(`<[^>]+>`)
	reSpaces  = regexp.MustCompile(`[ \t]+`)
	reBlank   = regexp.MustCompile(`\n{3,}`)
)

// StripHTML converts email HTML to readable plain text.
func StripHTML(s string) string {
	if s == "" {
		return ""
	}

	// <pre><code>...</code></pre> → indented code block
	s = rePreCode.ReplaceAllStringFunc(s, func(match string) string {
		inner := rePreCode.FindStringSubmatch(match)
		if len(inner) < 2 {
			return match
		}
		code := html.UnescapeString(inner[1])
		lines := strings.Split(code, "\n")
		var b strings.Builder
		b.WriteString("\n")
		for _, line := range lines {
			b.WriteString("    ")
			b.WriteString(line)
			b.WriteString("\n")
		}
		return b.String()
	})

	// <a href="url">text</a> → text (url) — but skip if text == url
	s = reA.ReplaceAllStringFunc(s, func(match string) string {
		parts := reA.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		href, text := parts[1], parts[2]
		text = strings.TrimSpace(text)
		if text == "" || text == href {
			return href
		}
		return text + " (" + href + ")"
	})

	// <br> → newline
	s = reBr.ReplaceAllString(s, "\n")

	// </p> → double newline
	s = reP.ReplaceAllString(s, "\n\n")

	// Strip remaining tags
	s = reTag.ReplaceAllString(s, "")

	// Unescape HTML entities
	s = html.UnescapeString(s)

	// Collapse horizontal whitespace (not newlines)
	s = reSpaces.ReplaceAllString(s, " ")

	// Collapse excessive blank lines
	s = reBlank.ReplaceAllString(s, "\n\n")

	return strings.TrimSpace(s)
}

// WordWrap wraps text to the given width, preserving indented lines.
func WordWrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var b strings.Builder
	for _, paragraph := range strings.Split(s, "\n") {
		if strings.HasPrefix(paragraph, "    ") {
			b.WriteString(paragraph)
			b.WriteString("\n")
			continue
		}
		col := 0
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			b.WriteString("\n")
			continue
		}
		for _, w := range words {
			wLen := len(w)
			if col > 0 && col+1+wLen > width {
				b.WriteString("\n")
				col = 0
			}
			if col > 0 {
				b.WriteString(" ")
				col++
			}
			b.WriteString(w)
			col += wLen
		}
		b.WriteString("\n")
	}
	return b.String()
}
