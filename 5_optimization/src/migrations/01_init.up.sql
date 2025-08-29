CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       customer_id VARCHAR(255) NOT NULL UNIQUE,
                       name VARCHAR(255) NOT NULL,
                       phone VARCHAR(16) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE
);
CREATE TABLE addresses (
                           id SERIAL PRIMARY KEY,
                           customer_id VARCHAR(255) NOT NULL REFERENCES users(customer_id),
                           zip          VARCHAR(20) NOT NULL,
                           city         VARCHAR(100) NOT NULL,
                           address      TEXT NOT NULL,
                           region       VARCHAR(100) NOT NULL,
    UNIQUE(customer_id, zip, city, address, region)
);
CREATE TABLE users_addresses (
                                 user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 address_id INTEGER NOT NULL REFERENCES addresses(id),
    PRIMARY KEY(user_id, address_id)
);

CREATE TABLE payments (
                          transaction VARCHAR(255) PRIMARY KEY,
                          request_id VARCHAR(255) NOT NULL,
                          currency VARCHAR(3) NOT NULL,
                          provider VARCHAR(100) NOT NULL,
                          amount INTEGER NOT NULL,
                          payment_dt BIGINT NOT NULL,
                          bank VARCHAR(100) NOT NULL,
                          delivery_cost INTEGER,
                          goods_total INTEGER,
                          custom_fee INTEGER
);

CREATE TABLE orders (
                        order_uid         VARCHAR(255) PRIMARY KEY,
                        track_number      VARCHAR(255) NOT NULL UNIQUE,
                        entry             VARCHAR(50) NOT NULL,
                        delivery INTEGER NOT NULL REFERENCES addresses(id) ON DELETE CASCADE,
                        payment VARCHAR(255) NOT NULL REFERENCES payments(transaction) ON DELETE CASCADE,
                        locale VARCHAR(8),
                        internal_signature TEXT,
                        customer_id VARCHAR(255) NOT NULL REFERENCES users(customer_id) ON DELETE CASCADE,
                        delivery_service VARCHAR(50) NOT NULL,
                        shardkey TEXT NOT NULL,
                        sm_id INTEGER NOT NULL,
                        date_created TIMESTAMP NOT NULL,
                        oof_shard TEXT NOT NULL
);

CREATE TABLE items (
                       nm_id INTEGER PRIMARY KEY,
                       chrt_id INTEGER NOT NULL,
                       price INTEGER,
                       name VARCHAR(255) NOT NULL,
                       size VARCHAR(50) NOT NULL,
                       brand VARCHAR(255) NOT NULL
);


CREATE TABLE order_items (
                             order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
                             item_id INTEGER NOT NULL REFERENCES items(nm_id) ON DELETE CASCADE,
                             rid VARCHAR(255) NOT NULL,
                             track_number VARCHAR(255) NOT NULL REFERENCES orders(track_number) ON DELETE CASCADE,
                             sale INTEGER NOT NULL,
                             total_price INTEGER,
                             status INTEGER NOT NULL NOT NULL,
                             PRIMARY KEY(order_uid, item_id)
);
