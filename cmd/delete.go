package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
)

func Delete() {
	acct, err := config.Load()
	if err != nil {
		ui.Fatalf("%v", err)
	}

	ui.Header("tmail")
	fmt.Printf("\n  delete %s%s%s? [y/N] ", ui.Cyan, acct.Address, ui.Reset)

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		ui.Status("cancelled")
		return
	}

	client := api.NewWithToken(acct.Token)
	if err := client.DeleteAccount(acct.ID); err != nil {
		ui.Fatalf("failed to delete account: %v", err)
	}

	if err := config.Delete(); err != nil {
		ui.Fatalf("failed to remove config: %v", err)
	}

	fmt.Println()
	ui.Success("account deleted")
	fmt.Println()
}
