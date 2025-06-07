# Desain Database

## 1. Deskripsi

Desain database bertujuan untuk memastikan penyimpanan data yang efisien, mudah diakses, serta mendukung seluruh fitur sistem, meliputi:

* **Pencatatan transaksi penjualan**
* **Pengelolaan stok (inbound/outbound)**
* **Pelacakan riwayat perubahan stok dan harga**

Desain relational database menyusun tabel-tabel dengan penentuan Primary Key (PK), Foreign Key (FK), dan constraints untuk memastikan:

* **Data integrity**: setiap record valid dan konsisten.
* **Efisiensi**: operasi CRUD optimal.
* **Scalability**: mudah dikembangkan fitur baru.

---

## 2. Entitas dan Relasi Utama

| Tabel                   | Deskripsi                                                               |
| ----------------------- | ----------------------------------------------------------------------- |
| **users**               | Data pengguna (admin, karyawan)                                         |
| **categories**          | Kategori produk (oli, ban, dll.)                                        |
| **products**            | Informasi produk: nama, stok snapshot, harga snapshot, lokasi, kategori |
| **activities**          | Transaksi penjualan (header)                                            |
| **activity\_items**     | Detail item pada transaksi penjualan                                    |
| **stock\_transactions** | Riwayat perubahan stok (inbound/outbound)                               |
| **price\_history**      | Riwayat perubahan harga produk                                          |

---

## 3. Struktur Tabel & Kunci

| Tabel                   | Primary Key (PK) | Foreign Key (FK)                                                                 | Deskripsi Singkat                                               |
| ----------------------- | ---------------- | -------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| **users**               | `id`             | —                                                                                | Data pengguna (admin, karyawan).                                |
| **categories**          | `id`             | —                                                                                | Kategori produk (oli, ban, dll.).                               |
| **products**            | `id`             | `category_id` → categories.id                                                    | Informasi produk (nama, stok snapshot, harga snapshot, lokasi). |
| **activities**          | `id`             | `user_id` → users.id                                                             | Header transaksi penjualan.                                     |
| **activity\_items**     | `id`             | `activity_id` → activities.id<br>`product_id` → products.id                      | Detail item transaksi (quantity, harga, diskon, final\_price).  |
| **stock\_transactions** | `id`             | `product_id` → products.id<br>`activity_item_id` → activity\_items.id (nullable) | Riwayat perubahan stok (inbound/outbound).                      |
| **price\_history**      | `id`             | `product_id` → products.id                                                       | Riwayat perubahan harga produk.                                 |

---

## 4. ERD (DBML)

```dbml
Table users {
  id          int      [pk, increment]
  name        varchar
  email       varchar  [unique]
  password    varchar
  role        varchar  // admin | karyawan
  created_at  datetime
  updated_at  datetime
}

Table categories {
  id          int      [pk, increment]
  name        varchar
  created_at  datetime
  updated_at  datetime
}

Table products {
  id          int      [pk, increment]
  name        varchar
  stock       int      // current snapshot
  price       decimal  // current snapshot
  location    varchar
  category_id int      [ref: > categories.id]
  created_at  datetime
  updated_at  datetime
}

Table activities {
  id          int      [pk, increment]
  user_id     int      [ref: > users.id]
  date        datetime
  status      varchar  // success | failed
  created_at  datetime
  updated_at  datetime
}

Table activity_items {
  id              int      [pk, increment]
  activity_id     int      [ref: > activities.id]
  product_id      int      [ref: > products.id]
  quantity        int
  price_at_time   decimal  // harga sebelum diskon
  discount_amount decimal  // potongan harga
  final_price     decimal  // harga setelah diskon
  created_at      datetime
  updated_at      datetime
}

Table stock_transactions {
  id                int      [pk, increment]
  product_id        int      [ref: > products.id]
  activity_item_id  int?     [ref: > activity_items.id]
  change_quantity   int      // positif=inbound, negatif=outbound
  date              datetime
  note              varchar  // misal: "Sale", "Restock"
  created_at        datetime
  updated_at        datetime
}

Table price_history {
  id            int      [pk, increment]
  product_id    int      [ref: > products.id]
  old_price     decimal
  new_price     decimal
  date_changed  datetime
  created_at    datetime
  updated_at    datetime
}

/*
Relationships:
users             1--* activities
categories        1--* products
products          1--* activity_items
activities        1--* activity_items
products          1--* stock_transactions
activity_items    1--* stock_transactions
products          1--* price_history
*/
```

---

## 5. Alur Interaksi Database

### 5.1 Proses Penjualan (Outbound Stock)

```
    actor Karyawan
    participant API
    participant activities as "activities"
    participant items as "activity_items"
    participant stocks as "stock_transactions"
    participant products as "products"

    Karyawan->>API: POST /activities {user_id, date, status="success"}
    API->>activities: INSERT activity

    loop setiap item
      API->>items: INSERT activity_item {activity_id, product_id, quantity, price_at_time, discount_amount, final_price}
      API->>stocks: INSERT stock_transaction {product_id, activity_item_id, change_quantity: -quantity, date, note: "Sale"}
      API->>products: UPDATE products SET stock = stock - quantity
    end
```

### 5.2 Proses Restock (Inbound Stock)

```
    actor Admin
    participant API
    participant stocks as "stock_transactions"
    participant products as "products"

    Admin->>API: POST /stock_transactions {product_id, change_quantity: +n, note: "Restock", date}
    API->>stocks: INSERT stock_transaction
    API->>products: UPDATE products SET stock = stock + n
```

### 5.3 Proses Ubah Harga

```
    actor Admin
    participant API
    participant prices as "price_history"
    participant products as "products"

    Admin->>API: PUT /products/:id/price {new_price}
    API->>prices: INSERT price_history {product_id, old_price, new_price, date_changed}
    API->>products: UPDATE products SET price = new_price
```

---

## 6. Constraints & Rules

| Tabel                   | Kolom                                               | Constraint                             |
| ----------------------- | --------------------------------------------------- | -------------------------------------- |
| **users**               | `id`                                                | PK, auto increment, NOT NULL           |
|                         | `name`, `password`                                  | NOT NULL                               |
|                         | `email`                                             | UNIQUE, NOT NULL                       |
|                         | `role`                                              | ENUM('admin','karyawan'), NOT NULL     |
| **categories**          | `id`                                                | PK, auto increment, NOT NULL           |
|                         | `name`                                              | UNIQUE, NOT NULL                       |
| **products**            | `id`                                                | PK, auto increment, NOT NULL           |
|                         | `name`, `location`                                  | NOT NULL                               |
|                         | `stock`                                             | NOT NULL, CHECK(`stock` ≥ 0)           |
|                         | `price`                                             | NOT NULL, CHECK(`price` ≥ 0)           |
|                         | `category_id`                                       | NOT NULL                               |
| **activities**          | `id`                                                | PK, auto increment, NOT NULL           |
|                         | `user_id`, `date`, `status`                         | NOT NULL                               |
|                         | `status`                                            | ENUM('success','failed'), NOT NULL     |
| **activity\_items**     | `id`                                                | PK, auto increment, NOT NULL           |
|                         | `activity_id`, `product_id`                         | NOT NULL, FK                           |
|                         | `quantity`                                          | NOT NULL, CHECK(`quantity` > 0)        |
|                         | `price_at_time`                                     | NOT NULL, CHECK(`price_at_time` ≥ 0)   |
|                         | `discount_amount`                                   | NOT NULL, CHECK(`discount_amount` ≥ 0) |
|                         | `final_price`                                       | NOT NULL, CHECK(`final_price` ≥ 0)     |
| **stock\_transactions** | `id`                                                | PK, auto increment, NOT NULL           |
|                         | `product_id`, `change_quantity`, `date`             | NOT NULL                               |
|                         | `change_quantity`                                   | CHECK(`change_quantity` ≠ 0)           |
| **price\_history**      | `id`                                                | PK, auto increment, NOT NULL           |
|                         | `product_id`,`old_price`,`new_price`,`date_changed` | NOT NULL                               |
|                         | `old_price`,`new_price`                             | CHECK(≥ 0)                             |

---

## 7. Relasi Antar Tabel

* **One-to-Many**

  * users → activities
  * categories → products
  * products → activity\_items
  * activities → activity\_items
  * products → stock\_transactions
  * activity\_items → stock\_transactions
  * products → price\_history

* **Many-to-Many**

  * products ↔ activities via activity\_items (junction table)

* **One-to-One**

  * (tidak ada dalam model ini)

---

## 8. Normalisasi

1. **1NF**: Setiap kolom berisi satu nilai atomic, tidak ada array/list.
2. **2NF**: Semua tabel menggunakan PK single-column (`id`), sehingga tidak ada partial dependency.
3. **3NF**: Tidak ada transitive dependency—atribut non-key bergantung langsung pada PK.

> Contoh: Histori harga dipindahkan ke tabel `price_history` agar tabel `products` hanya menyimpan snapshot harga terkini.

---