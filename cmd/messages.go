package cmd

import (
	"fmt"
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

	// Table header
	fmt.Printf("  %s%-4s %-25s %-35s %s%s\n", ui.Dim, "#", "FROM", "SUBJECT", "DATE", ui.Reset)
	fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", 75), ui.Reset)

	for i, m := range msgs {
		from := m.From.Address
		if len(from) > 24 {
			from = from[:21] + "..."
		}
		subject := m.Subject
		if len(subject) > 34 {
			subject = subject[:31] + "..."
		}
		date := formatDate(m.CreatedAt)

		marker := " "
		if !m.Seen {
			marker = ui.Cyan + "●" + ui.Reset
		}

		fmt.Printf(" %s%-4d %-25s %-35s %s\n", marker, i+1, from, subject, date)
	}
	fmt.Println()
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
