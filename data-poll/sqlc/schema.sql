CREATE TABLE IF NOT EXISTS usage
(
    usage_id          serial PRIMARY KEY,
    type              varchar(16),
    duration          integer,

    spotPerKwh        decimal(15, 5),
    perKwh            decimal(15, 5),
    kwh               decimal(15, 5),
    cost              decimal(15, 5),

    date              date,
    nemTime           timestamp with time zone,
    startTime         timestamp with time zone,
    endTime           timestamp with time zone,

    renewables        decimal(15, 20),

    channelType       varchar(16),
    channelIdentifier varchar(16),

    spikeStatus       varchar(16),
    descriptor        varchar(16),
    quality           varchar(16),

    tariffInformation jsonb,
    demandWindow      bool
)