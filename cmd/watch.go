package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
)

func Watch() {
	acct, err := config.Load()
	if err != nil {
		ui.Fatalf("%v", err)
	}

	client := api.NewWithToken(acct.Token)

	ui.Header("tmail")
	fmt.Printf("  %swatching: %s%s\n", ui.Dim, acct.Address, ui.Reset)
	fmt.Printf("  %spress Ctrl+C to stop%s\n\n", ui.Dim, ui.Reset)

	// Handle Ctrl+C gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		fmt.Println()
		ui.Status("stopped watching")
		os.Exit(0)
	}()

	backoff := time.Second
	maxBackoff := 30 * time.Second
	failures := 0

	for {
		err := client.Watch(acct.ID, func(msg api.SSEMessage) {
			// Reset backoff on successful message
			backoff = time.Second
			failures = 0

			from := msg.From.Address
			if msg.From.Name != "" {
				from = msg.From.Name + " <" + msg.From.Address + ">"
			}
			ui.Success(fmt.Sprintf("new message from %s%s%s", ui.Cyan, from, ui.Reset))
			ui.Status(msg.Subject)
			if msg.Intro != "" {
				ui.Dimf("%s", msg.Intro)
			}
			fmt.Println()
		})
		if err != nil {
			failures++
			if failures >= 5 {
				ui.Error(fmt.Sprintf("connection failed after %d attempts: %v", failures, err))
				os.Exit(1)
			}
			ui.Error(fmt.Sprintf("connection lost: %v, retrying in %s...", err, backoff))
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}
