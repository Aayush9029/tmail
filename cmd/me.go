package cmd

import (
	"fmt"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
)

func Me() {
	acct, err := config.Load()
	if err != nil {
		ui.Fatalf("%v", err)
	}

	client := api.NewWithToken(acct.Token)
	info, err := client.GetAccount(acct.ID)
	if err != nil {
		ui.Fatalf("failed to fetch account info: %v", err)
	}

	ui.Header("tmail")
	fmt.Println()
	fmt.Printf("  %sAddress:%s  %s%s%s\n", ui.Bold, ui.Reset, ui.Cyan, info.Address, ui.Reset)
	fmt.Printf("  %sID:%s       %s\n", ui.Bold, ui.Reset, info.ID)
	fmt.Printf("  %sCreated:%s  %s\n", ui.Bold, ui.Reset, formatDate(info.CreatedAt))
	fmt.Printf("  %sQuota:%s    %s / %s\n", ui.Bold, ui.Reset, formatBytes(info.Used), formatBytes(info.Quota))

	status := ui.Green + "active" + ui.Reset
	if info.IsDisabled {
		status = ui.Red + "disabled" + ui.Reset
	}
	fmt.Printf("  %sStatus:%s   %s\n", ui.Bold, ui.Reset, status)
	fmt.Println()
}

func formatBytes(b int) string {
	switch {
	case b >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	case b >= 1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	default:
		return fmt.Sprintf("%d B", b)
	}
}
