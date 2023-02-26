-------------------------------------------------------------------------------------------------------
-- \gexec: Sends the current query buffer to the server,
-- then treats each column of each row of the query's output (if any) as a SQL statement to be executed.
-- \gexec is a psql meta command: https://www.postgresql.org/docs/current/app-psql.html
--------------------------------------------------------------------------------------------------------
SELECT 'CREATE DATABASE user_service' WHERE NOT EXISTS(SELECT FROM pg_database WHERE datname = 'user_service');
\gexec
