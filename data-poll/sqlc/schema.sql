CREATE TABLE IF NOT EXISTS usage
(
    usage_id          serial PRIMARY KEY,
    type              varchar(16)              NOT NULL,
    duration          integer                  NOT NULL,

    spotPerKwh        decimal(15, 5)           NOT NULL,
    perKwh            decimal(15, 5)           NOT NULL,
    kwh               decimal(15, 5)           NOT NULL,
    cost              decimal(15, 5)           NOT NULL,

    date              date                     NOT NULL,
    nemTime           timestamp with time zone NOT NULL,
    startTime         timestamp with time zone NOT NULL,
    endTime           timestamp with time zone NOT NULL,

    renewables        decimal(15, 5)           NOT NULL,

    channelType       varchar(16)              NOT NULL,
    channelIdentifier varchar(16)              NOT NULL,

    spikeStatus       varchar(16)              NOT NULL,
    descriptor        varchar(16)              NOT NULL,
    quality           varchar(16)              NOT NULL,

    tariffInformation jsonb                    NOT NULL,
    demandWindow      bool                     NOT NULL
)