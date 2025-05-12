package backend

import "go.uber.org/zap"

type App struct {
	SpcClient *Client
	Logger    *zap.Logger
}

func NewApp(logger *zap.Logger, token string) App {
	client := NewClient(token)
	return App{&client, logger}
}
