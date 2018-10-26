-- script from https://github.com/tyleragreen/gtfs-schema

DROP TABLE agency;
DROP TABLE stops;
DROP TABLE routes;
DROP TABLE calendar_dates;
DROP TABLE trips;
DROP TABLE stop_times;

CREATE TABLE agency
(
  agency_id         text UNIQUE NULL,
  agency_name       text NOT NULL,
  agency_url        text NOT NULL,
  agency_timezone   text NOT NULL,
  agency_lang       text NULL
);

CREATE TABLE stops
(
  stop_id           text PRIMARY KEY,
  stop_name         text NOT NULL,
  stop_lat          double precision NOT NULL,
  stop_lon          double precision NOT NULL,
  location_type     text NULL
);

CREATE INDEX idx_stop_id ON stop_times (stop_id);

CREATE TABLE routes
(
  route_id          text PRIMARY KEY,
  agency_id         text NULL,
  route_short_name  text NULL,
  route_long_name   text NULL,
  route_type        integer NULL,
  route_url         text NULL
);

CREATE TABLE calendar_dates
(
  service_id text NOT NULL,
  date numeric(8) NOT NULL,
  exception_type integer NOT NULL
);

CREATE TABLE trips
(
  route_id          text NOT NULL,
  service_id        text NOT NULL,
  trip_id           text NOT NULL PRIMARY KEY,
  trip_headsign     text NULL,
  trip_short_name   text NULL
);

CREATE TABLE stop_times
(
  trip_id           text NOT NULL,
  arrival_time      interval NOT NULL,
  departure_time    interval NOT NULL,
  stop_id           text NOT NULL,
  stop_sequence     integer NOT NULL,
  pickup_type       integer NULL CHECK(pickup_type >= 0 and pickup_type <=3),
  drop_off_type     integer NULL CHECK(drop_off_type >= 0 and drop_off_type <=3)
);

CREATE INDEX idx_trip_id ON stop_times (trip_id);

\copy agency from './agency.txt' with csv header
\copy stops from './stops.txt' with csv header
\copy routes from './routes.txt' with csv header
\copy calendar_dates from './calendar_dates.txt' with csv header
\copy trips from './trips.txt' with csv header
\copy stop_times from './stop_times.txt' with csv header