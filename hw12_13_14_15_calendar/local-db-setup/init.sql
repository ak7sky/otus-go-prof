CREATE DATABASE calendar;
\connect calendar;
CREATE SCHEMA IF NOT EXISTS calendar;
CREATE USER calendar_user WITH ENCRYPTED PASSWORD 'calendar_pswd';
GRANT CONNECT ON DATABASE calendar TO calendar_user;
GRANT USAGE, CREATE ON SCHEMA calendar TO calendar_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA calendar TO calendar_user;