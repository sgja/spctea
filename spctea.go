package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/jsonq"
)

func main() {
	// Create a Resty Client
	client := resty.New()

	//register(client, "Nexor")
	//register(client, "Nexor")
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiTkVYT1IiLCJ2ZXJzaW9uIjoidjIuMS4wIiwicmVzZXRfZGF0ZSI6IjIwMjMtMTAtMjgiLCJpYXQiOjE2OTg5NjA1NTksInN1YiI6ImFnZW50LXRva2VuIn0.FBoqYd2o-HhMHyinsrDzFL6Xj5ybsZVUhS2GGWeYTckIHFm5ltjqksr_diSY5FJUfw34sz-BCC2cCn2ZmjOjEiAtsjpaILZuzN9KUr46G1RraQ0kylLCVOweQASzwvcGQnHfZVjCB4cDb2qFfDGI_wPPSQzbbNNKIFCf112wOj_wVjT2z21QyZYA3oHX-rBf-xsz1prP57Q_1hN1jOLHbPkrPQqUBP8Ira18jGm6UCW2r-r7L2XCEbNeF4iurqPKwG3ll98eR8j_wq5Lkeh5L7z89auW9NII-xAA1GW7N88y5MZgFMMEZPG6nRnuAdwCCSsiBs0RfATC_0D9vIhrdw"
	factions, err := list_factions(client, token)
	if err != nil {
		fmt.Println(err)
	} else {
		show_factions(client, factions)
	}
}

type Agent struct {
	Symbol  string `json:"symbol"`
	Faction string `json:"faction"`
}

func (a *Agent) to_body() ([]byte, error) {
	return json.Marshal(a)
}

type Token struct {
	Data TokenData `json:"data"`
}

type TokenData struct {
	Token string `json:"token"`
}

func register(client *resty.Client, callsign string, faction string) (string, string, error) {
	agent := Agent{Symbol: callsign, Faction: faction}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(agent).
		Post("https://api.spacetraders.io/v2/register")
	if err != nil {
		return "", "", err
	} else {
		data := make(map[string]interface{})
		err = json.Unmarshal(resp.Body(), &data)
		if err != nil {
			return "", "", err
		}
		jq := jsonq.NewQuery(data)
		token, e := jq.String("data", "token")
		return token, string(resp.Body()), e
	}

}

func list_factions(client *resty.Client, token string) ([]Faction, error) {
	url := "https://api.spacetraders.io/v2/factions"
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}
	data := Factions{}
	err = json.Unmarshal(resp.Body(), &data)
	return data.Data, err
}

type Factions struct {
	Data []Faction `json:"data"`
}

type Faction struct {
	Symbol       string         `json:"symbol"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Headquarters string         `json:"headquarters"`
	Traits       []FactionTrait `json:"traits"`
}

type FactionTrait struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
