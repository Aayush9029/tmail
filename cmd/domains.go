package cmd

import (
	"fmt"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/ui"
)

func Domains() {
	client := api.New()
	domains, err := client.ListDomains()
	if err != nil {
		ui.Fatalf("failed to fetch domains: %v", err)
	}

	if len(domains) == 0 {
		ui.Error("no domains available")
		return
	}

	ui.Header("tmail")
	fmt.Println()
	for _, d := range domains {
		status := ui.Green + "active" + ui.Reset
		if !d.IsActive {
			status = ui.Dim + "inactive" + ui.Reset
		}
		fmt.Printf("  %s%-30s%s  %s\n", ui.Bold, d.Domain, ui.Reset, status)
	}
	fmt.Println()
}
