package main

import (
	"context"
	"fmt"
	"io"
	"log"
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
		log.Fatal(fmt.Errorf("error loading config: %w", err))
	}

	usage := amber.NewUsageService(cfg.ServerUrl, cfg.ApiKey, cfg.Site)
	data, err := usage.GetUsage(time.Now().Add(-2*time.Hour*24), time.Now().Add(-2*time.Hour*24))

	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, data)
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
