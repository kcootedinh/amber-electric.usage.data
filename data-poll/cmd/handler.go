package main

import (
	"amber-electric.usage.data/sqlc"
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"amber-electric.usage.data/internal/amber"
)

type dbQueries interface {
	ListUsages(ctx context.Context) ([]sqlc.Usage, error)
	InsertUsage(ctx context.Context, arg sqlc.InsertUsageParams) (sqlc.Usage, error)
}

func handler(ctx context.Context, w io.Writer, usage amber.Service, db dbQueries) func() error {
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
