package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
)

func Messages() {
	acct, err := config.Load()
	if err != nil {
		ui.Fatalf("%v", err)
	}

	client := api.NewWithToken(acct.Token)
	msgs, err := client.ListMessages()
	if err != nil {
		ui.Fatalf("failed to fetch messages: %v", err)
	}

	ui.Header("tmail")
	fmt.Printf("  %sinbox: %s%s\n\n", ui.Dim, acct.Address, ui.Reset)

	if len(msgs) == 0 {
		ui.Dimf("no messages yet")
		fmt.Println()
		return
	}

	w := ui.TermWidth()
	printMessageTable(msgs, w)

	// Interactive mode: prompt to select a message
	if !ui.IsTTY() {
		return
	}
	fmt.Printf("  %senter message # to open in browser (or press enter to skip):%s ", ui.Dim, ui.Reset)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return
	}
	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(msgs) {
		ui.Fatalf("invalid message number: %s", input)
	}

	msg, err := client.GetMessage(msgs[num-1].ID)
	if err != nil {
		ui.Fatalf("failed to read message: %v", err)
	}
	openMessageInBrowser(msg)
}

func printMessageTable(msgs []api.MessageSummary, termWidth int) {
	// Layout: " ● NN  FROM  SUBJECT  DATE"
	// Fixed: marker=2, num=4, date=6, spacing=6 → overhead ~18
	// Remaining split: from gets 30%, subject gets 70%
	overhead := 18
	avail := termWidth - overhead
	if avail < 30 {
		avail = 30
	}
	fromW := avail * 3 / 10
	if fromW < 12 {
		fromW = 12
	}
	subjW := avail - fromW
	if subjW < 12 {
		subjW = 12
	}

	hdrFmt := fmt.Sprintf("  %%s%%-%ds %%-%ds %%-%ds %%s%%s\n", 4, fromW, subjW)
	fmt.Printf(hdrFmt, ui.Dim, "#", "FROM", "SUBJECT", "DATE", ui.Reset)
	fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", termWidth-4), ui.Reset)

	rowFmt := fmt.Sprintf(" %%s%%-%dd %%-%ds %%-%ds %%s\n", 4, fromW, subjW)
	for i, m := range msgs {
		from := m.From.Address
		if m.From.Name != "" {
			from = m.From.Name
		}
		from = truncate(from, fromW)
		subject := truncate(m.Subject, subjW)
		date := formatDate(m.CreatedAt)

		marker := " "
		if !m.Seen {
			marker = ui.Cyan + "●" + ui.Reset
		}

		fmt.Printf(rowFmt, marker, i+1, from, subject, date)
	}
	fmt.Println()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func formatDate(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso[:10]
	}
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04")
	}
	return t.Format("Jan 02")
}

func openMessageInBrowser(msg *api.MessageFull) {
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
		cmd = exec.Command("open", "-a", "Safari", tmpFile)
	default:
		cmd = exec.Command("xdg-open", tmpFile)
	}
	if err := cmd.Start(); err != nil {
		ui.Fatalf("failed to open browser: %v", err)
	}
	ui.Success("opened in Safari")
}
