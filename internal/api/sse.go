package api

import "time"

// Watch polls for new messages and calls onMessage for each unseen one.
// Polls every interval. Blocks until stopped externally.
func (c *Client) Watch(onMessage func(MessageSummary)) error {
	seen := make(map[string]bool)

	// Seed with existing messages so we don't re-notify
	msgs, err := c.ListMessages()
	if err != nil {
		return err
	}
	for _, m := range msgs {
		seen[m.ID] = true
	}

	for {
		time.Sleep(5 * time.Second)

		msgs, err := c.ListMessages()
		if err != nil {
			return err
		}

		for i := len(msgs) - 1; i >= 0; i-- {
			m := msgs[i]
			if !seen[m.ID] {
				seen[m.ID] = true
				onMessage(m)
			}
		}
	}
}
