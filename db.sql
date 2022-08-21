CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    locale VARCHAR(255),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(255)
    );

CREATE TABLE IF NOT EXISTS item (
    id SERIAL PRIMARY KEY,
    chrt_id INTEGER,
    track_number VARCHAR(255),
    price INTEGER,
    rid VARCHAR(255),
    name VARCHAR(255),
    size VARCHAR(255),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(255),
    status INTEGER
);

CREATE TABLE IF NOT EXISTS order_item (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    item_id INTEGER REFERENCES item(id) ON DELETE CASCADE,
    sale INTEGER
    );

CREATE TABLE IF NOT EXISTS delivery (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255),
    phone VARCHAR(255),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
    );

CREATE TABLE IF NOT EXISTS payment (
    transaction VARCHAR(255) PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt INTEGER,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
    );

CREATE OR REPLACE PROCEDURE add_order(
    order_uid VARCHAR(255),
    track_number VARCHAR(255),
    entry VARCHAR(255),
    locale VARCHAR(255),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(255)
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO orders
    VALUES (order_uid, track_number, entry, locale, internal_signature, customer_id,
            delivery_service, shardkey, sm_id, date_created, oof_shard);
END$$;

CREATE OR REPLACE PROCEDURE add_order_item(
    order_uid VARCHAR(255),
    sale INTEGER,
    chrt_id INTEGER,
    track_number VARCHAR(255),
    price INTEGER,
    name VARCHAR(255),
    rid VARCHAR(255),
    size VARCHAR(255),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(255),
    status INTEGER
)
LANGUAGE plpgsql
AS $$
DECLARE
    item_id INTEGER;
BEGIN
    INSERT INTO item
    VALUES (chrt_id, track_number, price, rid, name, size, total_price, nm_id, brand, status)
    ON CONFLICT DO NOTHING
    RETURNING id INTO item_id;

    INSERT INTO order_item
    VALUES (item_id, order_uid, sale)
    ON CONFLICT DO NOTHING;
END$$;

CREATE OR REPLACE PROCEDURE add_delivery(
    order_uid VARCHAR(255),
    name VARCHAR(255),
    phone VARCHAR(255),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO delivery
    VALUES (order_uid, name, phone, zip, city, address, region, email);
END$$;

CREATE OR REPLACE PROCEDURE add_payment(
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt INTEGER,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO payment
    VALUES (transaction, request_id, currency, provider, amount, payment_dt, bank,
            delivery_cost, goods_total, custom_fee);
END$$;