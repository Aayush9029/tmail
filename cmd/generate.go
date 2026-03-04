package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
)

const (
	addrLen = 8
	passLen = 12
	charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func randomString(length int, chars string) string {
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

func copyToClipboard(text string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return false
	}
	cmd.Stdin = nil
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return false
	}
	if err := cmd.Start(); err != nil {
		return false
	}
	pipe.Write([]byte(text))
	pipe.Close()
	return cmd.Wait() == nil
}

func Generate() {
	ui.Header("tmail")

	// Check if account already exists
	if existing, _ := config.Load(); existing != nil {
		ui.Error(fmt.Sprintf("account already exists: %s", existing.Address))
		ui.Status("run 'tmail delete' first to remove it")
		return
	}

	client := api.New()

	// Get domains
	ui.Waiting("fetching domains...")
	domains, err := client.ListDomains()
	if err != nil {
		ui.Fatalf("failed to fetch domains: %v", err)
	}

	var domain string
	for _, d := range domains {
		if d.IsActive {
			domain = d.Domain
			break
		}
	}
	if domain == "" {
		ui.Fatalf("no active domains available")
	}

	// Generate credentials
	address := randomString(addrLen, charset) + "@" + domain
	password := randomString(passLen, charset+"ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%")

	// Create account
	ui.Waiting("creating account...")
	acct, err := client.CreateAccount(address, password)
	if err != nil {
		ui.Fatalf("failed to create account: %v", err)
	}

	// Get token
	token, _, err := client.GetToken(address, password)
	if err != nil {
		ui.Fatalf("failed to get token: %v", err)
	}

	// Save config
	if err := config.Save(&config.Account{
		ID:       acct.ID,
		Address:  address,
		Password: password,
		Token:    token,
	}); err != nil {
		ui.Fatalf("failed to save config: %v", err)
	}

	fmt.Println()
	ui.Success(fmt.Sprintf("created %s%s%s", ui.Cyan, address, ui.Reset))

	if copyToClipboard(address) {
		ui.Status("copied to clipboard")
	}
	fmt.Println()
}
