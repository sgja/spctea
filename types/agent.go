package types

import "encoding/json"

type Agent struct {
	AccountId       string `json:"accountId,omitempty"`
	Symbol          string `json:"symbol,omitempty"`
	Headquarters    string `json:"headquarters"`
	Credits         int64  `json:"credits"`
	StartingFaction string `json:"startingFaction,omitempty"`
	ShipCount       int32  `json:"shipCount"`
}

func (a *Agent) ToBody() ([]byte, error) {
	return json.Marshal(a)
}
