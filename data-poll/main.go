package main

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"amber-electric.usage.data/internal/amber"
)

func run(ctx context.Context, w io.Writer, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := loadConfig(getenv)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	slog.SetLogLoggerLevel(cfg.LogLevel)

	s, err := gocron.NewScheduler()

	usage := amber.NewUsageService(cfg.ServerUrl, cfg.ApiKey, cfg.Site)

	// first run
	err = handler(ctx, w, usage)()
	if err != nil {
		return err
	}

	if cfg.Frequency <= 0 {
		slog.Debug("no frequency set, exiting")
		return nil
	}

	interval := time.Duration(cfg.Frequency) * time.Minute
	slog.Debug(fmt.Sprintf("scheduling job every %s minutes", interval))

	_, err = s.NewJob(gocron.DurationJob(interval), gocron.NewTask(handler(ctx, w, usage)),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					slog.Error(fmt.Sprintf("job %s ended with error: %s", jobName, err))
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

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Getenv); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func handler(ctx context.Context, w io.Writer, usage amber.Service) func() error {
	return func() error {
		data, err := usage.GetUsage(time.Now().Add(-2*time.Hour*24), time.Now().Add(-2*time.Hour*24))

		if err != nil {
			slog.Error(fmt.Sprintf("failed to retrieve usage data: %s", err.Error()))
			return err
		}

		_, err = fmt.Fprintln(w, data)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to print usage data: %s", err.Error()))
			return err
		}

		return nil
	}
}
