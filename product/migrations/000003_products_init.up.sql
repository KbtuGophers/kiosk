CREATE TABLE IF NOT EXISTS products
(
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id             VARCHAR PRIMARY KEY,
    category_id    VARCHAR NOT NULL,
    barcode         VARCHAR NOT NULL UNIQUE ,
    name        varchar NOT NULL,
    measure    VARCHAR,
    cost int,
    producer_country   VARCHAR,
    brand_name VARCHAR NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR NOT NULL,
    is_weighted BOOLEAN NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories (id)

    );
