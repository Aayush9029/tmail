package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const MercureURL = "https://mercure.mail.tm/.well-known/mercure"

type SSEMessage struct {
	Type    string         `json:"@type"`
	ID      string         `json:"id"`
	From    MessageAddr    `json:"from"`
	Subject string         `json:"subject"`
	Intro   string         `json:"intro"`
}

// Watch subscribes to the Mercure SSE endpoint and calls onMessage for each new message.
// It blocks until the context is cancelled or an error occurs.
func (c *Client) Watch(accountID string, onMessage func(SSEMessage)) error {
	url := fmt.Sprintf("%s?topic=/accounts/%s", MercureURL, accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "text/event-stream")

	httpClient := &http.Client{Timeout: 0} // no timeout for SSE
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SSE connection failed: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	// Increase buffer for large messages
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var dataLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
			continue
		}

		// Empty line = end of event
		if line == "" && len(dataLines) > 0 {
			raw := strings.Join(dataLines, "\n")
			dataLines = nil

			var msg SSEMessage
			if err := json.Unmarshal([]byte(raw), &msg); err != nil {
				continue // skip malformed events
			}
			if msg.Type == "Message" || msg.ID != "" {
				onMessage(msg)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// If we reach here, the connection was closed — retry after a short delay
	time.Sleep(2 * time.Second)
	return fmt.Errorf("SSE connection closed")
}
