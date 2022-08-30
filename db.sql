DROP TABLE IF EXISTS orders CASCADE;

CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY UNIQUE,
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

DROP TABLE IF EXISTS delivery CASCADE;

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

DROP TABLE IF EXISTS payment CASCADE;

CREATE TABLE IF NOT EXISTS payment (
    transaction VARCHAR(255) PRIMARY KEY UNIQUE,
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

DROP TABLE IF EXISTS item CASCADE;

-- could be composite primary key sql
CREATE TABLE IF NOT EXISTS item (
    chrt_id INTEGER PRIMARY KEY UNIQUE,
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

DROP TABLE IF EXISTS order_item CASCADE;

CREATE TABLE IF NOT EXISTS order_item (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER REFERENCES item(chrt_id) ON DELETE CASCADE,
    sale INTEGER
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
    order_uid VARCHAR(255),
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
    VALUES (transaction, order_uid, request_id, currency, provider, amount, payment_dt, bank,
            delivery_cost, goods_total, custom_fee);
END$$;

DROP PROCEDURE IF EXISTS add_order_item(character varying,integer,integer,character varying,integer,character varying,character varying,character varying,integer,integer,character varying,integer);

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
BEGIN
    INSERT INTO item
    VALUES (chrt_id, track_number, price, rid, name, size, total_price, nm_id, brand, status)
    ON CONFLICT DO NOTHING;

    INSERT INTO order_item
    VALUES (chrt_id, order_uid, sale)
    ON CONFLICT DO NOTHING;
END$$;

CREATE OR REPLACE PROCEDURE get_order(uid VARCHAR(255))
    LANGUAGE plpgsql
AS $$
BEGIN
    SELECT
        json_build_object(
            orders.*,
            'delivery', to_json(delivery.*),
            'payment', to_json(payment.*),
            'items', (SELECT json_agg(order_items.*)
                      FROM (SELECT item.*, order_item.sale
                            FROM order_item,
                                 item
                            WHERE order_item.order_uid = uid
                              AND order_item.chrt_id = item.chrt_id) as order_items)
        )
    FROM orders,
         payment,
         delivery
    WHERE orders.order_uid = uid AND
          delivery.order_uid = uid AND
          payment.order_uid = uid
    LIMIT 1;
END$$;

DROP PROCEDURE IF EXISTS get_all_orders();

CREATE OR REPLACE PROCEDURE get_all_orders()
LANGUAGE plpgsql
AS $$
BEGIN
    SELECT
        json_agg(
            json_build_object(
                orders.*,
                'delivery', to_json(delivery.*),
                'payment', to_json(payment.*),
                'items', (SELECT json_agg(order_items.*)
                          FROM (SELECT item.*, order_item.sale
                                FROM order_item,
                                     item
                                WHERE order_item.order_uid = orders.order_uid
                                  AND order_item.chrt_id = item.chrt_id) as order_items)
            )
        )
    FROM orders,
         payment,
         delivery;
END$$;
