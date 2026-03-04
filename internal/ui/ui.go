package ui

import (
	"fmt"
	"os"
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
