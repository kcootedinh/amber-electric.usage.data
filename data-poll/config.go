package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ServerUrl string
	ApiKey    string
	Site      string
	LogLevel  slog.Level
}

func loadConfig(getenv func(string) string) (Config, error) {
	url := getenv("SERVER_URL")
	if url == "" {
		return Config{}, fmt.Errorf("SERVER_URL not set")
	}

	apiKey := getenv("API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("API_KEY not set")
	}

	site := getenv("SITE")
	if site == "" {
		return Config{}, fmt.Errorf("SITE not set")
	}

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "0"
	}

	logLevelInt, err := strconv.ParseInt(logLevel, 10, 32)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse LOG_LEVEL, %s", err.Error()))
		logLevelInt = 0
	}

	return Config{
		ServerUrl: url,
		ApiKey:    apiKey,
		Site:      site,
		LogLevel:  slog.Level(logLevelInt),
	}, nil
}
