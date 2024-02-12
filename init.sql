BEGIN;

DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS transactions;

CREATE TABLE IF NOT EXISTS clients (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "name" VARCHAR(80) NOT NULL,
    "limit" INTEGER NOT NULL,
    "balance" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "value" INTEGER NOT NULL,
    "type" CHAR(1) NOT NULL,
    "description" VARCHAR(10) NOT NULL,
    "client_id" INTEGER NOT NULL,
    CONSTRAINT fk_client FOREIGN KEY("client_id") REFERENCES clients("id")
);

INSERT INTO
    clients ("name", "limit", "balance")
VALUES
    ('o barato sai caro', 1000 * 100, 0),
    ('zan corp ltda', 800 * 100, 0),
    ('les cruders', 10000 * 100, 0),
    ('padaria joia de cocaia', 100000 * 100, 0),
    ('kid mais', 5000 * 100, 0);

COMMIT;