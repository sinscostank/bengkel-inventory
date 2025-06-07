-- 0001_create_tables.up.sql

CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role ENUM('admin','karyawan') NOT NULL,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE products (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  stock INT NOT NULL CHECK (stock >= 0),
  price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
  location VARCHAR(255) NOT NULL,
  category_id INT NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE activities (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  date DATETIME NOT NULL,
  status ENUM('success','failed') NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE activity_items (
  id INT AUTO_INCREMENT PRIMARY KEY,
  activity_id INT NOT NULL,
  product_id INT NOT NULL,
  quantity INT NOT NULL CHECK (quantity > 0),
  price_at_time DECIMAL(10,2) NOT NULL CHECK (price_at_time >= 0),
  discount_amount DECIMAL(10,2) NOT NULL CHECK (discount_amount >= 0),
  final_price DECIMAL(10,2) NOT NULL CHECK (final_price >= 0),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (activity_id) REFERENCES activities(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE stock_transactions (
  id INT AUTO_INCREMENT PRIMARY KEY,
  product_id INT NOT NULL,
  activity_item_id INT,
  change_quantity INT NOT NULL CHECK (change_quantity <> 0),
  date DATETIME NOT NULL,
  note VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (product_id) REFERENCES products(id),
  FOREIGN KEY (activity_item_id) REFERENCES activity_items(id)
);

CREATE TABLE price_history (
  id INT AUTO_INCREMENT PRIMARY KEY,
  product_id INT NOT NULL,
  old_price DECIMAL(10,2) NOT NULL CHECK (old_price >= 0),
  new_price DECIMAL(10,2) NOT NULL CHECK (new_price >= 0),
  date_changed DATETIME NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (product_id) REFERENCES products(id)
);
