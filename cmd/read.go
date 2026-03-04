package cmd

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
)

func Read(args []string) {
	fs := flag.NewFlagSet("read", flag.ExitOnError)
	browser := fs.Bool("browser", false, "open message in browser")
	fs.BoolVar(browser, "b", false, "open message in browser")
	fs.Parse(args)

	if fs.NArg() < 1 {
		ui.Fatalf("usage: tmail read <number> [--browser|-b]")
	}

	num, err := strconv.Atoi(fs.Arg(0))
	if err != nil || num < 1 {
		ui.Fatalf("invalid message number: %s", fs.Arg(0))
	}

	acct, err := config.Load()
	if err != nil {
		ui.Fatalf("%v", err)
	}

	client := api.NewWithToken(acct.Token)

	msgs, err := client.ListMessages()
	if err != nil {
		ui.Fatalf("failed to fetch messages: %v", err)
	}
	if num > len(msgs) {
		ui.Fatalf("message #%d not found (inbox has %d messages)", num, len(msgs))
	}

	msg, err := client.GetMessage(msgs[num-1].ID)
	if err != nil {
		ui.Fatalf("failed to read message: %v", err)
	}

	if *browser {
		openMessageInBrowser(msg)
		return
	}

	w := ui.TermWidth()

	ui.Header("tmail")
	fmt.Println()

	from := msg.From.Address
	if msg.From.Name != "" {
		from = msg.From.Name + " <" + msg.From.Address + ">"
	}
	fmt.Printf("  %sFrom:%s    %s\n", ui.Bold, ui.Reset, from)
	fmt.Printf("  %sSubject:%s %s\n", ui.Bold, ui.Reset, msg.Subject)
	fmt.Printf("  %sDate:%s    %s\n", ui.Bold, ui.Reset, formatDate(msg.CreatedAt))
	divider := w - 4
	if divider < 20 {
		divider = 20
	}
	fmt.Printf("  %s%s%s\n\n", ui.Dim, strings.Repeat("─", divider), ui.Reset)

	// Render body
	htmlContent := strings.Join(msg.HTML, "\n")
	var body string
	if htmlContent != "" {
		body = ui.StripHTML(htmlContent)
	} else {
		body = msg.Text
	}

	if body == "" {
		ui.Dimf("(empty message)")
	} else {
		wrapWidth := w - 4
		if wrapWidth < 40 {
			wrapWidth = 40
		}
		wrapped := ui.WordWrap(body, wrapWidth)
		// Indent each line by 2 spaces
		for _, line := range strings.Split(wrapped, "\n") {
			if line == "" {
				fmt.Println()
			} else {
				fmt.Printf("  %s\n", line)
			}
		}
	}
	fmt.Println()
}
