package api

import "encoding/json"

type MessageSummary struct {
	ID        string      `json:"id"`
	From      MessageAddr `json:"from"`
	Subject   string      `json:"subject"`
	Intro     string      `json:"intro"`
	CreatedAt string      `json:"createdAt"`
	Seen      bool        `json:"seen"`
}

type MessageAddr struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type MessageFull struct {
	ID        string        `json:"id"`
	From      MessageAddr   `json:"from"`
	To        []MessageAddr `json:"to"`
	Subject   string        `json:"subject"`
	Intro     string        `json:"intro"`
	Text      string        `json:"text"`
	HTML      []string      `json:"html"`
	CreatedAt string        `json:"createdAt"`
	Seen      bool          `json:"seen"`
}

type messagesResponse struct {
	Members []MessageSummary `json:"hydra:member"`
}

func (c *Client) ListMessages() ([]MessageSummary, error) {
	data, err := c.get("/messages")
	if err != nil {
		return nil, err
	}
	var resp messagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Members, nil
}

func (c *Client) GetMessage(id string) (*MessageFull, error) {
	data, err := c.get("/messages/" + id)
	if err != nil {
		return nil, err
	}
	var msg MessageFull
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
