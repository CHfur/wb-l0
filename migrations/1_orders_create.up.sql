CREATE TABLE IF NOT EXISTS orders
(
    uid  varchar(255)
        CONSTRAINT firstkey PRIMARY KEY,
    data jsonb not null
);
