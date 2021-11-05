-- +goose Up
CREATE TABLE events_status
(
  id           SERIAL PRIMARY KEY NOT NULL,
  event_id     INT NOT NULL REFERENCES events(id),
  status       VARCHAR(100) NOT NULL default 'New',
  processed_at TIMESTAMPTZ,
  sent_at      TIMESTAMPTZ
);

create index event_idx on events_status using btree (event_id, status);

-- +goose Down
DROP TABLE events_status;