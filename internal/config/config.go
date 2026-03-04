package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Account struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "tmail")
}

func configPath() string {
	return filepath.Join(configDir(), "account.json")
}

func Load() (*Account, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return nil, fmt.Errorf("no account found — run 'tmail generate' first")
	}
	var acct Account
	if err := json.Unmarshal(data, &acct); err != nil {
		return nil, fmt.Errorf("corrupt config: %w", err)
	}
	return &acct, nil
}

func Save(acct *Account) error {
	if err := os.MkdirAll(configDir(), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(acct, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0o600)
}

func Delete() error {
	return os.RemoveAll(configDir())
}
