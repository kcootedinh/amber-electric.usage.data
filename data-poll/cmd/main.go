package main

import (
	"amber-electric.usage.data/cmd/backfill"
	"amber-electric.usage.data/cmd/cron"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"amber-electric.usage.data/internal/amber"
	"amber-electric.usage.data/sqlc"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func run(ctx context.Context, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := loadConfig(getenv)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	slog.SetLogLoggerLevel(cfg.LogLevel)

	s, err := gocron.NewScheduler()

	usage := amber.NewUsageService(cfg.AmberUrl, cfg.AmberApiKey, cfg.Site)
	dbConn, err := connectDb(ctx, cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassord, cfg.DbName)

	defer func() {
		if err := dbConn.Close(ctx); err != nil {
			slog.Error("error closing db connection", "error", err)
		}
	}()

	queries := sqlc.New(dbConn)

	if cfg.BackfillStart != "" {
		return backfill.Handler(ctx, usage, queries, cfg.BackfillStart)
	}

	if cfg.Frequency <= 0 {
		slog.Debug("no frequency set, exiting")
		return nil
	}

	interval := time.Duration(cfg.Frequency) * time.Minute
	slog.Debug(fmt.Sprintf("scheduling job every %s minutes", interval))

	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(10, 30, 0))),
		gocron.NewTask(cron.Handler(ctx, usage, queries)),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					slog.Error(fmt.Sprintf("job %s with ID %s ended with error: %s", jobName, jobID, err))
				},
			),
		))
	if err != nil {
		return fmt.Errorf("error creating cron job %w", err)
	}

	s.Start()

	// Wait for CTRL-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// We block here until a CTRL-C / SigInt is received
	// Once received, we exit and the server is cleaned up
	<-sigChan

	slog.Info("shutting down")
	err = s.Shutdown()
	if err != nil {
		return err
	}

	return nil
}

func connectDb(ctx context.Context, host string, port int, user string, password string, databaseName string) (*pgx.Conn, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, databaseName)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to connect to database: %v", err))
		return nil, err
	}

	return conn, nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
