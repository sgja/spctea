package backend

import (
	"encoding/json"
	"fmt"

	"sgja/spctea/types"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/jsonq"
)

const BaseURL = "https://api.spacetraders.io/v2/"

func GetURL(endpoint string) string {
	return BaseURL + endpoint
}

func (c *Client) BuildRequestWithAuth() *resty.Request {
	return c.rc.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.token))
}

func (c *Client) BuildRequest() *resty.Request {
	return c.rc.R().
		SetHeader("Content-Type", "application/json")
}

type Client struct {
	rc    *resty.Client
	token string
}

func NewClient(token string) Client {
	return Client{resty.New(), token}
}

func (c *Client) Register(callsign string, faction string) (string, string, error) {
	agent := types.Agent{Symbol: callsign, Faction: faction}
	body, err := agent.ToBody()
	if err != nil {
		return "", "", err
	}
	resp, err := c.BuildRequestWithAuth().
		SetBody(string(body)).
		Post(GetURL("register"))
	if err != nil {
		return "", "", err
	} else {
		data := make(map[string]any)
		err = json.Unmarshal(resp.Body(), &data)
		if err != nil {
			return "", "", err
		}
		jq := jsonq.NewQuery(data)
		token, e := jq.String("data", "token")
		return token, string(resp.Body()), e
	}

}

func (c *Client) ListFactions() ([]types.Faction, error) {
	resp, err := c.BuildRequest().
		Get(GetURL("factions"))
	if err != nil {
		return nil, err
	}
	data := types.Factions{}
	err = json.Unmarshal(resp.Body(), &data)
	return data.Data, err
}
