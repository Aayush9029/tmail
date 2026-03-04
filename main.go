package main

import (
	"fmt"
	"os"

	"github.com/Aayush9029/tmail/cmd"
	"github.com/Aayush9029/tmail/internal/ui"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	switch os.Args[1] {
	case "generate", "g":
		cmd.Generate()
	case "messages", "m":
		cmd.Messages()
	case "read", "r":
		cmd.Read(os.Args[2:])
	case "delete", "d":
		cmd.Delete()
	case "me":
		cmd.Me()
	case "watch", "w":
		cmd.Watch()
	case "domains":
		cmd.Domains()
	case "--version", "-v", "version":
		fmt.Printf("tmail %s\n", version)
	case "--help", "-h", "help":
		showHelp()
	default:
		ui.Error(fmt.Sprintf("unknown command: %s", os.Args[1]))
		fmt.Println()
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	if ui.IsTTY() {
		ui.Header("tmail")
		fmt.Printf("  %sdisposable email in your terminal%s\n", ui.Dim, ui.Reset)
	}
	fmt.Println()
	fmt.Printf("  %sUSAGE%s\n", ui.Blue, ui.Reset)
	fmt.Printf("    tmail <command> [options]\n")
	fmt.Println()
	fmt.Printf("  %sCOMMANDS%s\n", ui.Blue, ui.Reset)
	fmt.Printf("    generate, g    Create new disposable email\n")
	fmt.Printf("    messages, m    List inbox messages\n")
	fmt.Printf("    read, r <n>    Read message #n [--plain|-p] [--browser|-b]\n")
	fmt.Printf("    delete, d      Delete account\n")
	fmt.Printf("    me             Show account info\n")
	fmt.Printf("    watch, w       Watch for new messages (real-time)\n")
	fmt.Printf("    domains        List available domains\n")
	fmt.Println()
	fmt.Printf("  %sOPTIONS%s\n", ui.Blue, ui.Reset)
	fmt.Printf("    -h, --help     Show this help\n")
	fmt.Printf("    -v, --version  Show version\n")
	fmt.Println()
}
