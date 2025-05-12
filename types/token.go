package types

type Token struct {
	Data TokenData `json:"data"`
}

type TokenData struct {
	Token string `json:"token"`
}
