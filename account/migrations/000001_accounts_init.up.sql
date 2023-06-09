
CREATE TABLE IF NOT EXISTS account_types
(
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id           SERIAL PRIMARY KEY,
    name         VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts
(
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id           VARCHAR PRIMARY KEY,
    user_name    VARCHAR NOT NULL,
    phone_number    VARCHAR NOT NULL,
    account_type_id int NOT NULL ,
    profile_photo VARCHAR NOT NULL,
    FOREIGN KEY (account_type_id) REFERENCES account_types (id)

);

