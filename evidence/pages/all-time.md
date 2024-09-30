---
title: All Time Usage 
---

```sql costPerDay
SELECT date, sum(cost) / 100 as cost, sum(kwh) as kwh
FROM amber.usage
GROUP BY date
ORDER BY date;
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


<CalendarHeatmap
data={costPerDay}
date=date
value=cost
/>