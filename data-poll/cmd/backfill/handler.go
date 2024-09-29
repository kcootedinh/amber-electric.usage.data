package backfill

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"amber-electric.usage.data/internal/amber"
	"amber-electric.usage.data/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/ratelimit"
)

type dbQueries interface {
	ListUsages(ctx context.Context) ([]sqlc.Usage, error)
	GetUsagesForDate(ctx context.Context, date pgtype.Date) ([]sqlc.Usage, error)
	InsertUsage(ctx context.Context, arg sqlc.InsertUsageParams) (sqlc.Usage, error)
}

func Handler(ctx context.Context, usage amber.Service, db dbQueries, backfillStart string) error {
	date, err := time.Parse(time.DateOnly, backfillStart)
	if err != nil {
		return fmt.Errorf("parsing backfill start date: %w", err)
	}

	rl := ratelimit.New(1) // per second

	for ; date.Before(time.Now()); date = date.Add(time.Hour * 24) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_ = rl.Take()

			dbUsages, err := db.GetUsagesForDate(ctx, pgtype.Date{Time: date, Valid: true})
			if err != nil {
				return err
			}

			if len(dbUsages) != 0 {
				slog.Warn(fmt.Sprintf("Found %d usages on db for date %v, skipping", len(dbUsages), date), "date", date)
				continue
			}

			data, err := usage.GetUsage(date, date)

			if err != nil {
				slog.Error(fmt.Sprintf("failed to retrieve usage data: %s", err.Error()))
				return err
			}

			if len(data) == 0 {
				slog.Info(fmt.Sprintf("no usage data found for date %s", date))
				continue
			}

			slog.Debug(fmt.Sprintf("inserting %d rows usage data", len(data)), "date", date)

			for _, u := range data {
				insertUsage, err := toSqlcUsage(u)
				if err != nil {
					return err
				}

				_, err = db.InsertUsage(ctx, insertUsage)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func toSqlcUsage(u amber.Usage) (sqlc.InsertUsageParams, error) {
	tariffJson, err := json.Marshal(u.TariffInformation)
	if err != nil {
		slog.Error("failed to marshal TariffInformation for db", slog.Any("TariffInformation", u.TariffInformation))
	}

	spotPerKwh := pgtype.Numeric{}
	if err := spotPerKwh.Scan(strconv.FormatFloat(u.SpotPerKwh, 'f', -1, 64)); err != nil {
		return sqlc.InsertUsageParams{}, fmt.Errorf("failed converting spotPerKwh for db: %w", err)
	}

	perKwh := pgtype.Numeric{}
	if err := perKwh.Scan(strconv.FormatFloat(u.PerKwh, 'f', -1, 64)); err != nil {
		return sqlc.InsertUsageParams{}, fmt.Errorf("failed converting perKwh for db: %w", err)
	}

	kwh := pgtype.Numeric{}
	if err := kwh.Scan(strconv.FormatFloat(u.Kwh, 'f', -1, 64)); err != nil {
		return sqlc.InsertUsageParams{}, fmt.Errorf("failed converting kwh for db: %w", err)
	}

	cost := pgtype.Numeric{}
	if err := cost.Scan(strconv.FormatFloat(u.Cost, 'f', -1, 64)); err != nil {
		return sqlc.InsertUsageParams{}, fmt.Errorf("failed converting cost for db: %w", err)
	}

	renewables := pgtype.Numeric{}
	if err := renewables.Scan(strconv.FormatFloat(u.Renewables, 'f', -1, 64)); err != nil {
		return sqlc.InsertUsageParams{}, fmt.Errorf("failed converting renewables for db: %w", err)
	}

	date, err := time.Parse(time.DateOnly, u.Date)
	if err != nil {
		return sqlc.InsertUsageParams{}, fmt.Errorf("failed converting date for db: %w", err)
	}

	insertUsage := sqlc.InsertUsageParams{
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
	}
	return insertUsage, nil
}
