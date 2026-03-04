package api

import "encoding/json"

type createAccountRequest struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type AccountInfo struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	Quota     int    `json:"quota"`
	Used      int    `json:"used"`
	IsDisabled bool  `json:"isDisabled"`
	CreatedAt string `json:"createdAt"`
}

func (c *Client) CreateAccount(address, password string) (*AccountInfo, error) {
	data, err := c.post("/accounts", createAccountRequest{
		Address:  address,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	var acct AccountInfo
	if err := json.Unmarshal(data, &acct); err != nil {
		return nil, err
	}
	return &acct, nil
}

func (c *Client) GetAccount(id string) (*AccountInfo, error) {
	data, err := c.get("/accounts/" + id)
	if err != nil {
		return nil, err
	}
	var acct AccountInfo
	if err := json.Unmarshal(data, &acct); err != nil {
		return nil, err
	}
	return &acct, nil
}

func (c *Client) DeleteAccount(id string) error {
	_, err := c.delete("/accounts/" + id)
	return err
}
