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
--   {
--     "type": "Usage",
--     "duration": 30,
--     "date": "2024-09-05",
--     "endTime": "2024-09-04T14:30:00Z",
--     "quality": "billable",
--     "kwh": 0.18,
--     "nemTime": "2024-09-05T00:30:00+10:00",
--     "perKwh": 13.04418,
--     "channelType": "general",
--     "channelIdentifier": "E1",
--     "cost": 2.348,
--     "renewables": 21.243000000000002,
--     "spotPerKwh": 5.53742,
--     "startTime": "2024-09-04T14:00:01Z",
--     "spikeStatus": "none",
--     "tariffInformation": {
--       "demandWindow": false
--     },
--     "descriptor": "veryLow"
--   },
--     "type": "Usage",
--     "duration": 5,
--     "spotPerKwh": 6.12,
--     "perKwh": 24.33,
--     "date": "2021-05-05",
--     "nemTime": "2021-05-06T12:30:00+10:00",
--     "startTime": "2021-05-05T02:00:01Z",
--     "endTime": "2021-05-05T02:30:00Z",
--     "renewables": 45,
--     "channelType": "general",
--     "tariffInformation": "string",
--     "spikeStatus": "none",
--     "descriptor": "negative",
--     "channelIdentifier": "E1",
--     "kwh": 0,
--     "quality": "estimated",
--     "cost": 0