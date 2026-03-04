package cmd

import (
	"fmt"
	"os"
	"os/signal"

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

	for {
		err := client.Watch(acct.ID, func(msg api.SSEMessage) {
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
			ui.Error(fmt.Sprintf("connection lost: %v, reconnecting...", err))
		}
	}
}
