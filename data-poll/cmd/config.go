package main

import (
	"fmt"
	"log/slog"
	"strconv"
)

type Config struct {
	AmberUrl    string
	AmberApiKey string
	Site        string
	Frequency   int
	DbHost      string
	DbPort      int
	DbName      string
	DbUser      string
	DbPassord   string
	LogLevel    slog.Level
}

func loadConfig(getenv func(string) string) (Config, error) {
	url := getenv("AMBER_URL")
	if url == "" {
		return Config{}, fmt.Errorf("AMBER_URL not set")
	}

	apiKey := getenv("AMBER_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("AMBER_API_KEY not set")
	}

	site := getenv("SITE")
	if site == "" {
		return Config{}, fmt.Errorf("SITE not set")
	}

	frequency := getenv("JOB_FREQUENCY")
	if frequency == "" {
		frequency = "0"
	}

	dbHost := getenv("DB_HOST")
	if dbHost == "" {
		return Config{}, fmt.Errorf("DB_HOST not set")
	}

	dbPort, err := strconv.Atoi(getenv("DB_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("DB_PORT not an integer")
	}

	dbName := getenv("DB_NAME")
	if dbName == "" {
		return Config{}, fmt.Errorf("DB_NAME not set")
	}

	dbUser := getenv("DB_USER")
	if dbUser == "" {
		return Config{}, fmt.Errorf("DB_USER not set")
	}

	dbPassord := getenv("DB_PASSWORD")
	if dbPassord == "" {
		return Config{}, fmt.Errorf("DB_PASSWORD not set")
	}

	frequencyInt, err := strconv.ParseInt(frequency, 10, 32)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse JOB_FREQUENCY, %s", err.Error()))
		frequencyInt = 0
	}

	logLevel := getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "0"
	}

	logLevelInt, err := strconv.ParseInt(logLevel, 10, 32)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse LOG_LEVEL, %s", err.Error()))
		logLevelInt = 0
	}

	return Config{
		AmberUrl:    url,
		AmberApiKey: apiKey,
		Site:        site,
		Frequency:   int(frequencyInt),
		DbHost:      dbHost,
		DbPort:      dbPort,
		DbName:      dbName,
		DbUser:      dbUser,
		DbPassord:   dbPassord,
		LogLevel:    slog.Level(logLevelInt),
	}, nil
}
