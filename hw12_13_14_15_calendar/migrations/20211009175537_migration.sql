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
 title       VARCHAR(100) NOT NULL,
 description text,
 start_date  TIMESTAMPTZ default now(),
 end_date    TIMESTAMPTZ default now(),
 notified_at TIMESTAMPTZ,
 created_at  TIMESTAMPTZ NOT NULL default now(),
 updated_at  TIMESTAMPTZ NOT NULL default now()
);

create index owner_idx on events (user_id);
create index start_idx on events using btree (start_date);

-- +goose Down
DROP TABLE events;
DROP TABLE users;
