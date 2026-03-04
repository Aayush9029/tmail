package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
	"github.com/k3a/html2text"
)

func Read(args []string) {
	fs := flag.NewFlagSet("read", flag.ExitOnError)
	plain := fs.Bool("plain", false, "strip HTML and show plain text")
	fs.BoolVar(plain, "p", false, "strip HTML and show plain text")
	browser := fs.Bool("browser", false, "open message in browser")
	fs.BoolVar(browser, "b", false, "open message in browser")
	fs.Parse(args)

	if fs.NArg() < 1 {
		ui.Fatalf("usage: tmail read <number> [--plain|-p] [--browser|-b]")
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

	// Get message list to resolve index → ID
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
		openInBrowser(msg)
		return
	}

	ui.Header("tmail")
	fmt.Println()
	fmt.Printf("  %sFrom:%s    %s\n", ui.Bold, ui.Reset, msg.From.Address)
	fmt.Printf("  %sSubject:%s %s\n", ui.Bold, ui.Reset, msg.Subject)
	fmt.Printf("  %sDate:%s    %s\n", ui.Bold, ui.Reset, formatDate(msg.CreatedAt))
	fmt.Printf("  %s%s%s\n\n", ui.Dim, strings.Repeat("─", 60), ui.Reset)

	if *plain || len(msg.HTML) == 0 {
		text := msg.Text
		if text == "" && len(msg.HTML) > 0 {
			text = html2text.HTML2Text(strings.Join(msg.HTML, "\n"))
		}
		if text == "" {
			ui.Dimf("(empty message)")
		} else {
			fmt.Println(text)
		}
	} else {
		// Render HTML as plain text (default behavior for terminal)
		htmlContent := strings.Join(msg.HTML, "\n")
		fmt.Println(html2text.HTML2Text(htmlContent))
	}
	fmt.Println()
}

func openInBrowser(msg *api.MessageFull) {
	htmlContent := strings.Join(msg.HTML, "\n")
	if htmlContent == "" {
		htmlContent = "<pre>" + msg.Text + "</pre>"
	}

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "tmail-message.html")
	if err := os.WriteFile(tmpFile, []byte(htmlContent), 0o644); err != nil {
		ui.Fatalf("failed to write temp file: %v", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", tmpFile)
	case "linux":
		cmd = exec.Command("xdg-open", tmpFile)
	default:
		ui.Fatalf("unsupported platform for browser open")
	}
	if err := cmd.Start(); err != nil {
		ui.Fatalf("failed to open browser: %v", err)
	}
	ui.Success(fmt.Sprintf("opened message #%s in browser", msg.ID))
}
