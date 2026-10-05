CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE readings (
  ts         TIMESTAMPTZ      NOT NULL,
  topic      TEXT             NOT NULL,
  value_num  DOUBLE PRECISION,
  value_bool BOOLEAN,
  value_text TEXT
);

SELECT create_hypertable('readings', by_range('ts'));

CREATE INDEX readings_topic_ts_idx ON readings (topic, ts DESC);

SELECT set_chunk_time_interval('readings', INTERVAL '1 day');

SELECT add_retention_policy('readings', INTERVAL '7 days');
