SET search_path TO fatch;

CREATE TABLE logs (
  id BIGSERIAL,
  timestamp TIMESTAMPTZ NOT NULL,

  -- hot, frequently filtered fields
  level TEXT NOT NULL,
  service TEXT NOT NULL,

  -- flexible payload
  log_data JSONB NOT NULL,
  PRIMARY KEY (timestamp, id)
) PARTITION BY RANGE (timestamp);


-- time index
CREATE INDEX idx_logs_timestamp
ON logs (timestamp);

-- hot filters
CREATE INDEX idx_logs_level
ON logs (level);

CREATE INDEX idx_logs_service
ON logs (service);

CREATE OR REPLACE FUNCTION insert_log(
  p_timestamp TIMESTAMPTZ,
  p_level TEXT,
  p_service TEXT,
  p_log_data JSONB
) RETURNS VOID
LANGUAGE plpgsql AS
$$
DECLARE
  partition_date DATE := date_trunc('day', p_timestamp);
  partition_name TEXT := 'logs_' || to_char(partition_date, 'YYYY_MM_DD');
BEGIN
  -- create partition if it does not exist
  EXECUTE format(
    'CREATE TABLE IF NOT EXISTS %I PARTITION OF logs
     FOR VALUES FROM (%L) TO (%L)',
    partition_name,
    partition_date,
    partition_date + INTERVAL '1 day'
  );

  -- safe to insert
  INSERT INTO logs (timestamp, level, service, log_data)
  VALUES (p_timestamp, p_level, p_service, p_log_data);
END;
$$;

-- select insert_log(
--   '2025-12-13T16:52:55.927765423Z',
--   'INFO',
--   'health',
--   '{"msg":"HTTP request","method":"GET","url":"/health","status":200,"remote_addr":"[::1]:39322","duration_ms":1}'
-- );

-- select * from logs
-- order by id DESC;
--
-- TRUNCATE TABLE logs RESTART IDENTITY CASCADE;
