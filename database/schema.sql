CREATE TABLE city (
  id serial PRIMARY KEY,
  name varchar(255) UNIQUE NOT NULL,
  latitude decimal NOT NULL,
  longitude decimal NOT NULL
);

CREATE TABLE temperature (
  id serial PRIMARY KEY,
  city_id int NOT NULL,
  max decimal NOT NULL,
  min decimal NOT NULL,
  timestamp timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE webhook (
  id serial PRIMARY KEY,
  city_id int NOT NULL,
  callback_url varchar(255) NOT NULL
);

ALTER TABLE temperature ADD FOREIGN KEY (city_id) REFERENCES city (id) ON DELETE CASCADE;

ALTER TABLE webhook ADD FOREIGN KEY (city_id) REFERENCES city (id) ON DELETE CASCADE;

CREATE INDEX city_index_0 ON city (name);

CREATE INDEX temperature_index_1 ON temperature (city_id);

CREATE INDEX webhook_index_2 ON webhook (city_id);
