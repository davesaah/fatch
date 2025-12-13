SET search_path TO fatch;

-- LOGS 

CREATE VIEW log_summary AS
SELECT level, COUNT(*) AS value
FROM fatch.logs
GROUP BY level
ORDER BY value DESC;

CREATE VIEW error_logs AS
SELECT service, log_data->>'msg' as msg, log_data->>'trace' as trace FROM fatch.logs
where level = 'ERROR'
order by timestamp DESC;

CREATE VIEW warning_logs AS
SELECT service, log_data->>'msg' as msg, log_data->>'trace' as trace FROM fatch.logs
where level = 'WARN'
order by timestamp DESC;

CREATE VIEW debug_logs AS
SELECT service, log_data->>'msg' as msg, log_data->>'trace' as trace FROM fatch.logs
where level = 'DEBUG'
order by timestamp DESC;

CREATE VIEW info_logs AS
SELECT service, log_data FROM fatch.logs
where level = 'INFO'
order by timestamp DESC;
