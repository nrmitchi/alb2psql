-- Select request time buckets (30 second intervals) grouped by minute/response code
SELECT
    date_trunc('minute', requesttime),
    backendresponsecode,
	width_bucket(backendprocessingtime, 0, 90, 180),
	count(*)
FROM
	alb_logs
WHERE
	requesttime > '2017-06-06 03:00:00' and
	requesttime < '2017-06-06 03:50:00'
GROUP BY
	date_trunc('minute', requesttime),
	backendresponsecode,
	width_bucket
ORDER BY
    date_trunc('minute', requesttime),
    backendresponsecode,
    width_bucket;
