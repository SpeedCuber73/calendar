CREATE TABLE IF NOT EXISTS events(
    uuid         text,
    title        text NOT NULL,
    start_at      timestamp NOT NULL,
    duration     interval NOT NULL,
    descr  text,
    user_name         text NOT NULL,
    notify_before interval,
    CONSTRAINT events_pkey PRIMARY KEY (uuid)
);
