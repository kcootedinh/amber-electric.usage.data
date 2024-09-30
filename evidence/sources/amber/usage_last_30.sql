SELECT *
FROM usage u
WHERE u.starttime > (NOW() - INTERVAL '30 DAY')