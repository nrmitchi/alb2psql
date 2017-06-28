-- Select request time buckets (30 second intervals) grouped by minute/response code
-- Pull it from a subquery so that we can multiple the bucket but the millisecond response time the bucket represents
SELECT 
	requesttime, 
	backendresponsecode, 
	width_bucket * 50 as bucket_ms, 
	count 
FROM (
	SELECT
		date_trunc('minute', requesttime) requesttime,
		backendresponsecode,
		width_bucket(backendprocessingtime, 0, 90, 1800),
		count(*) 
	FROM
		alb_logs
	WHERE
		requesttime > '2017-06-06 18:00:00' and
		requesttime < '2017-06-06 20:00:00'
	GROUP BY
		date_trunc('minute', requesttime),
		backendresponsecode,
		width_bucket
	ORDER BY
		date_trunc('minute', requesttime),
		backendresponsecode,
		width_bucket
) a;