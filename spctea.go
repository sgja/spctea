package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"sgja/spctea/backend"

	tea "github.com/charmbracelet/bubbletea"
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

	//developmentCfg := zap.NewDevelopmentEncoderConfig()
	//developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

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
	return strings.Trim(string(content), "\n"), nil
}

func main() {
	logger := createLogger()
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	token, err := load_token()
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Info("token loaded successfully")
	app := backend.NewApp(logger, token)

	if _, err := tea.NewProgram(RootScreen(&app), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
