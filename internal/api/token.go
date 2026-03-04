package api

import "encoding/json"

type tokenRequest struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
	ID    string `json:"id"`
}

func (c *Client) GetToken(address, password string) (string, string, error) {
	data, err := c.post("/token", tokenRequest{
		Address:  address,
		Password: password,
	})
	if err != nil {
		return "", "", err
	}
	var resp tokenResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", "", err
	}
	return resp.Token, resp.ID, nil
}
