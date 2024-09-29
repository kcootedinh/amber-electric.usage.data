---
title: Amber Electric - Usage 
---

```sql costPerDay
SELECT date, sum(cost) / 100 as cost
FROM amber.usage
GROUP BY date
ORDER BY date;
```

<BarChart
    data={costPerDay}
    title="Cost per day"
    x=date
    y=cost
/>