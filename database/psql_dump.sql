CREATE TABLE receipts (
    id SERIAL PRIMARY KEY,
    retailer VARCHAR(255),
    purchase_date DATE,
    purchase_time TIME,
    total NUMERIC(10, 2)
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    receipt_id INTEGER REFERENCES receipts(id),
    short_description VARCHAR(255),
    price NUMERIC(10, 2)
);

INSERT INTO receipts (retailer, purchase_date, purchase_time, total) VALUES
('Walgreens', '2022-01-02', '08:13', 2.65),
('Target', '2022-01-02', '13:13', 1.25);

INSERT INTO items (receipt_id, short_description, price) VALUES
(1, 'Pepsi - 12-oz', 1.25),
(1, 'Dasani', 1.40),
(2, 'Pepsi - 12-oz', 1.25);