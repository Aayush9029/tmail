package api

import "encoding/json"

type Domain struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
}

type domainsResponse struct {
	Members []Domain `json:"hydra:member"`
}

func (c *Client) ListDomains() ([]Domain, error) {
	data, err := c.get("/domains")
	if err != nil {
		return nil, err
	}
	var resp domainsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Members, nil
}
