ALTER TABLE products
MODIFY name VARCHAR(8) NOT NULL UNIQUE CHECK (name REGEXP '^[A-Za-z]{1,8}$');