
CREATE TABLE IF NOT EXISTS account.account_types
(
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id           SERIAL PRIMARY KEY,
    name         VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS account.accounts
(
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id           VARCHAR PRIMARY KEY,
    user_name    VARCHAR NOT NULL,
    phone_number    VARCHAR NOT NULL UNIQUE ,
    account_type_id int NOT NULL,
    latitude    NUMERIC,
    longitude   NUMERIC,
    profile_photo VARCHAR,
    FOREIGN KEY (account_type_id) REFERENCES account_types (id)

);

