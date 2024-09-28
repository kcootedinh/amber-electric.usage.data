package main

import (
	"amber-electric.usage.data/sqlc"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"io"
	"log/slog"
	"strconv"
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

		for _, u := range data {
			tariffJson, err := json.Marshal(u.TariffInformation)
			if err != nil {
				slog.Error("failed to marshal TariffInformation for db", slog.Any("TariffInformation", u.TariffInformation))
			}

			spotPerKwh := pgtype.Numeric{}
			if err := spotPerKwh.Scan(strconv.FormatFloat(u.SpotPerKwh, 'f', -1, 64)); err != nil {
				return fmt.Errorf("failed converting spotPerKwh for db: %w", err)
			}

			perKwh := pgtype.Numeric{}
			if err := perKwh.Scan(strconv.FormatFloat(u.PerKwh, 'f', -1, 64)); err != nil {
				return fmt.Errorf("failed converting perKwh for db: %w", err)
			}

			kwh := pgtype.Numeric{}
			if err := kwh.Scan(strconv.FormatFloat(u.Kwh, 'f', -1, 64)); err != nil {
				return fmt.Errorf("failed converting kwh for db: %w", err)
			}

			cost := pgtype.Numeric{}
			if err := cost.Scan(strconv.FormatFloat(u.Cost, 'f', -1, 64)); err != nil {
				return fmt.Errorf("failed converting cost for db: %w", err)
			}

			renewables := pgtype.Numeric{}
			if err := renewables.Scan(strconv.FormatFloat(u.Renewables, 'f', -1, 64)); err != nil {
				return fmt.Errorf("failed converting renewables for db: %w", err)
			}

			date, err := time.Parse(time.DateOnly, u.Date)
			if err != nil {
				return fmt.Errorf("failed converting date for db: %w", err)
			}

			_, err = db.InsertUsage(ctx, sqlc.InsertUsageParams{
				Type:              u.Type,
				Duration:          u.Duration,
				Spotperkwh:        spotPerKwh,
				Perkwh:            perKwh,
				Kwh:               kwh,
				Cost:              cost,
				Date:              pgtype.Date{Time: date, Valid: true},
				Nemtime:           pgtype.Timestamptz{Time: u.NemTime, Valid: true},
				Starttime:         pgtype.Timestamptz{Time: u.StartTime, Valid: true},
				Endtime:           pgtype.Timestamptz{Time: u.EndTime, Valid: true},
				Renewables:        renewables,
				Channeltype:       u.ChannelType,
				Channelidentifier: u.ChannelIdentifier,
				Spikestatus:       u.SpikeStatus,
				Descriptor:        u.Descriptor,
				Quality:           u.Quality,
				Tariffinformation: tariffJson,
				Demandwindow:      u.TariffInformation.DemandWindow,
			})
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintln(w, data)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to print usage data: %s", err.Error()))
			return err
		}

		return nil
	}
}
