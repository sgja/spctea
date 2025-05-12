package types

import "encoding/json"

type Agent struct {
	Symbol  string `json:"symbol"`
	Faction string `json:"faction"`
}

func (a *Agent) ToBody() ([]byte, error) {
	return json.Marshal(a)
}
