CREATE TABLE products (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    stockQunantity INT
);

CREATE TABLE stockHistory (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    productID INT,
    FOREIGN KEY (productID) REFERENCES products(id),
    dateChanged TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    quantityChanged INT
);

CREATE TABLE salesHistory (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    productID INT,
    FOREIGN KEY (productID) REFERENCES products(id),
    dateSold TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    quantitySold INT,
    totalAmount INT
);
