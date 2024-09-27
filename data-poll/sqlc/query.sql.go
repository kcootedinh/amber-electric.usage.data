// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertUsage = `-- name: InsertUsage :one
INSERT INTO usage(type,
                  duration,
                  spotperkwh,
                  perkwh,
                  kwh,
                  cost,
                  date,
                  nemtime,
                  starttime,
                  endtime,
                  renewables,
                  channeltype,
                  channelidentifier,
                  spikestatus,
                  descriptor,
                  quality,
                  tariffinformation,
                  demandwindow)
VALUES ($1, $2, $3, $4, $5, $6,
        $7, $8, $9, $10, $11,
        $12, $13, $14, $15,
        $16, $17, $18)
RETURNING usage_id, type, duration, spotperkwh, perkwh, kwh, cost, date, nemtime, starttime, endtime, renewables, channeltype, channelidentifier, spikestatus, descriptor, quality, tariffinformation, demandwindow
`

type InsertUsageParams struct {
	Type              pgtype.Text
	Duration          pgtype.Int4
	Spotperkwh        pgtype.Numeric
	Perkwh            pgtype.Numeric
	Kwh               pgtype.Numeric
	Cost              pgtype.Numeric
	Date              pgtype.Date
	Nemtime           pgtype.Timestamptz
	Starttime         pgtype.Timestamptz
	Endtime           pgtype.Timestamptz
	Renewables        pgtype.Numeric
	Channeltype       pgtype.Text
	Channelidentifier pgtype.Text
	Spikestatus       pgtype.Text
	Descriptor        pgtype.Text
	Quality           pgtype.Text
	Tariffinformation []byte
	Demandwindow      pgtype.Bool
}

func (q *Queries) InsertUsage(ctx context.Context, arg InsertUsageParams) (Usage, error) {
	row := q.db.QueryRow(ctx, insertUsage,
		arg.Type,
		arg.Duration,
		arg.Spotperkwh,
		arg.Perkwh,
		arg.Kwh,
		arg.Cost,
		arg.Date,
		arg.Nemtime,
		arg.Starttime,
		arg.Endtime,
		arg.Renewables,
		arg.Channeltype,
		arg.Channelidentifier,
		arg.Spikestatus,
		arg.Descriptor,
		arg.Quality,
		arg.Tariffinformation,
		arg.Demandwindow,
	)
	var i Usage
	err := row.Scan(
		&i.UsageID,
		&i.Type,
		&i.Duration,
		&i.Spotperkwh,
		&i.Perkwh,
		&i.Kwh,
		&i.Cost,
		&i.Date,
		&i.Nemtime,
		&i.Starttime,
		&i.Endtime,
		&i.Renewables,
		&i.Channeltype,
		&i.Channelidentifier,
		&i.Spikestatus,
		&i.Descriptor,
		&i.Quality,
		&i.Tariffinformation,
		&i.Demandwindow,
	)
	return i, err
}

const listUsages = `-- name: ListUsages :many
SELECT usage_id, type, duration, spotperkwh, perkwh, kwh, cost, date, nemtime, starttime, endtime, renewables, channeltype, channelidentifier, spikestatus, descriptor, quality, tariffinformation, demandwindow
FROM usage
`

func (q *Queries) ListUsages(ctx context.Context) ([]Usage, error) {
	rows, err := q.db.Query(ctx, listUsages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Usage
	for rows.Next() {
		var i Usage
		if err := rows.Scan(
			&i.UsageID,
			&i.Type,
			&i.Duration,
			&i.Spotperkwh,
			&i.Perkwh,
			&i.Kwh,
			&i.Cost,
			&i.Date,
			&i.Nemtime,
			&i.Starttime,
			&i.Endtime,
			&i.Renewables,
			&i.Channeltype,
			&i.Channelidentifier,
			&i.Spikestatus,
			&i.Descriptor,
			&i.Quality,
			&i.Tariffinformation,
			&i.Demandwindow,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}