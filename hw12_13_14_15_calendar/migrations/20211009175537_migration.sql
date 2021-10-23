-- +goose Up
CREATE TABLE users
(id         SERIAL PRIMARY KEY NOT NULL,
 email      VARCHAR(1000) NOT NULL,
 first_name VARCHAR(1000),
 last_name  VARCHAR(1000)
);

CREATE TABLE events
(id          SERIAL PRIMARY KEY NOT NULL,
 user_id     INT NOT NULL REFERENCES users(id),
 title       VARCHAR(1000) NOT NULL,
 description VARCHAR(1000),
 start_date  DATE NOT NULL,
 start_time  TIME default '00:00:00',
 end_date  DATE NOT NULL,
 end_time  TIME default '00:00:00',
 notified_at TIMESTAMPTZ,
 created_at  TIMESTAMPTZ NOT NULL default now(),
 updated_at  TIMESTAMPTZ NOT NULL default now()
);

create index owner_idx on events (user_id);
create index start_idx on events using btree (start_date, start_time);

-- +goose Down
DROP TABLE events;
DROP TABLE users;
