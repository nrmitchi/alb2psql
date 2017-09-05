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

WITH total_req_counts AS (
SELECT 
	date_trunc('hour', requesttime) bucket,
	targetgroup tg,
	count(*) c
FROM 
	alb_logs
GROUP BY
	date_trunc('hour', requesttime),
	targetgroup
)
SELECT 
	date_trunc('hour', requesttime) requesttime,
	elbresponsecode,
	backendresponsecode,
	count(*),
	count(*)::float / total_req_counts.c
FROM 
	alb_logs JOIN
	total_req_counts ON (total_req_counts.bucket = date_trunc('hour', requesttime) AND total_req_counts.tg = targetgroup)
WHERE 
	targetgroup = 'arn:aws:elasticloadbalancing:us-east-1:584106525078:targetgroup/vendor-https/c82f9357c85097f8' AND
	(backendresponsecode in ('400', '500', '501', '502', '503', '504' ) OR elbresponsecode in ('400', '500', '501', '502', '503', '504' )) AND
	requesttime > '2017-06-29 00:00:00' and
	requesttime < '2017-07-01 00:00:00'
GROUP BY
	date_trunc('hour', requesttime),
	elbresponsecode,
	backendresponsecode,
	total_req_counts.c
ORDER BY 
	date_trunc('hour', requesttime)


WITH total_req_counts AS (
SELECT 
	date_trunc('hour', requesttime) bucket,
	targetgroup tg,
	count(*) c
FROM 
	alb_logs
GROUP BY
	date_trunc('hour', requesttime),
	targetgroup
)
SELECT 
	date_trunc('hour', requesttime) requesttime,
	sum(case when elbresponsecode = '400' then 1 else 0 end) as x400s,
	sum(case when elbresponsecode in ('500','501','502','503','504') AND backendresponsecode != '-' then 1 else 0 end) as x5xxs,
	sum(case when elbresponsecode = '502' AND backendresponsecode = '-' then 1 else 0 end) as timeouts
FROM 
	alb_logs
WHERE 
	targetgroup = 'arn:aws:elasticloadbalancing:us-east-1:584106525078:targetgroup/vendor-https/c82f9357c85097f8' AND
	(backendresponsecode in ('400', '500', '501', '502', '503', '504' ) OR elbresponsecode in ('400', '500', '501', '502', '503', '504' )) AND
	requesttime > '2017-06-29 00:00:00' and
	requesttime < '2017-07-01 00:00:00'
GROUP BY
	date_trunc('hour', requesttime)
ORDER BY 
	date_trunc('hour', requesttime)




sum(case when place  = 'home' then 1 else 0 end) as Home,