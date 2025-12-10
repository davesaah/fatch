-- query statistics
SELECT userid, dbid, query, calls, mean_exec_time
FROM fatch.pg_stat_statements;
