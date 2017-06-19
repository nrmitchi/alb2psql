SELECT 
	width_bucket(backendprocessingtime, 0, 90, 90), 
	count(*)
FROM 
	alb_logs 
WHERE 
	requesttime > '2017-6-6' and 
	requesttime < '2017-6-7' 
GROUP BY 
	date_trunc('minute', requesttime),
	width_bucket;


SELECT 
	date_trunc('minute', requesttime),
	backendresponsecode,
	count(*)
FROM 
	alb_logs 
WHERE 
	requesttime > '2017-6-6' and 
	requesttime < '2017-6-7' 
GROUP BY 
	date_trunc('minute', requesttime),
	backendresponsecode;


