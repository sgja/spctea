package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/jsonq"
	zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createLogger() *zap.Logger {
	//stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/spctea.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		//zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return zap.New(core)
}

func load_token() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	token_path := path.Join(homedir, ".spctea", "token")
	content, err := os.ReadFile(token_path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {
	logger := createLogger()
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	// Create a Resty Client
	client := resty.New()

	token, err := load_token()
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Info("token loaded successfully")
	factions, err := list_factions(client)
	if err != nil {
		fmt.Println(err)
	} else {
		show_factions(client, logger, token, factions)
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

func register(client *resty.Client, account_token string, callsign string, faction string) (string, string, error) {
	agent := Agent{Symbol: callsign, Faction: faction}
	body, err := agent.to_body()
	if err != nil {
		return "", "", err
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", account_token)).
		SetBody(string(body)).
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

func list_factions(client *resty.Client) ([]Faction, error) {
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
