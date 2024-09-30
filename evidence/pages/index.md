---
title: Amber Electric - Usage Last 30 days
---

```sql costPerDay
SELECT date, sum(cost) / 100 as cost, sum(kwh) as kwh
FROM usage_last_30
GROUP BY date
ORDER BY date desc
```

<BarChart
    data={costPerDay}
    title="Cost per day"
    x=date
    y=cost
    yFmt=aud2
    y2=kwh
    y2SeriesType=line
/>


```sql usagePerHour
SELECT date_part('hour', u.starttime) as hourOfDay, avg(u.kwh) * 2 as kwhAvg
FROM usage_last_30 u
GROUP BY date_part('hour', u.starttime)
ORDER BY date_part('hour', u.starttime);
```

<LineChart
    data={usagePerHour}
    x=hourOfDay
    xAxisTitle="Hour of Day"
    y=kwhAvg
    yAxisTitle="Average kwh used per hour"
/>

```sql pricingPerHour
SELECT date_part('hour', u.starttime) as "Hour of Day",
       avg(u.perkwh)    / 100               as "Average per kwh",
       avg(u.spotperkwh)  / 100             as "Average spot per kwh"
FROM usage_last_30 u
GROUP BY date_part('hour', u.starttime)
ORDER BY date_part('hour', u.starttime);
```

<LineChart
    data={pricingPerHour}
    x="Hour of Day"
    xAxisTitle="Hour of Day"
    y="Average per kwh"
    yFmt=aud2
    yAxisTitle="Average per kwh"
    y2="Average spot per kwh"
    y2Fmt=aud2
    y2AxisTitle="Average spot per kwh"
/>

```sql maxPricePerHour
SELECT date_part('hour', u.starttime) as "Hour of Day",
       max(u.perkwh) / 100            as "Max per kwh",
       max(u.spotperkwh) / 100        as "Max spot per kwh"
FROM usage_last_30 u
GROUP BY date_part('hour', u.starttime)
ORDER BY date_part('hour', u.starttime);
```

<LineChart
    data={maxPricePerHour}
    x="Hour of Day"
    xAxisTitle="Hour of Day"
    y="Max per kwh"
    yFmt=aud2
    yAxisTitle="Max per kwh"
    y2="Max spot per kwh"
    y2Fmt=aud2
    y2AxisTitle="Max spot per kwh"
/>