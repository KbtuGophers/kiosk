CREATE TABLE IF NOT EXISTS otps
(
    id VARCHAR NOT NULL PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    key             VARCHAR NOT NULL UNIQUE,
    secret          VARCHAR NOT NULL,
    phone_number    VARCHAR NOT NULL,
    send_at         INTEGER NOT NULL,
    confirmed_at    INTEGER NOT NULL,
    attempts        INT NOT NULL DEFAULT 0,
    status          INT NOT NULL DEFAULT 1
)
