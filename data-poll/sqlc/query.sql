-- name: ListUsages :many
SELECT *
FROM usage;

-- name: InsertUsage :one
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
RETURNING *;