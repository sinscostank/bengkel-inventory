# 🛠️ Sistem Inventori & Penjualan Bengkel

Capstone Bootcamp – Backend Developer  
RESTful API untuk mengelola inventori produk, transaksi penjualan, dan laporan performa bisnis pada bengkel otomotif keluarga.

---

## 🚀 Tech Stack

- **Backend Framework**: Go (Gin)
- **ORM**: GORM
- **Database**: MySQL
- **Authentication**: JWT
- **Documentation**: Postman & Swagger (Swaggo)
- **Architecture**: Modular (controller, service, repository, etc.)

---

## 📁 Project Structure
/controller ← Handler untuk endpoint API
/service ← Logika bisnis
/repository ← Akses ke database (query layer)
/entity ← Struktur data / model
/config ← Konfigurasi database dan environment
/middleware ← JWT, validasi, otorisasi
/reports ← Laporan & ringkasan penjualan
/docs ← ERD, flow bisnis, dokumentasi tambahan
main.go
route.go


---

## 📌 Features

- ✅ **User Management**
  - Register & Login (JWT)
  - Role-based Authorization (Admin & Karyawan)

- 📦 **Manajemen Produk**
  - Tambah/Ubah/Hapus/Lihat produk
  - Informasi stok & lokasi penyimpanan
  - Validasi stok tidak negatif & nama unik

- 📑 **Transaksi Penjualan**
  - Pencatatan aktivitas penjualan
  - Pengurangan stok otomatis saat transaksi sukses
  - Cek ketersediaan stok sebelum transaksi

- 📊 **Laporan**
  - Produk terlaris
  - Produk dengan stok rendah
  - Transaksi berdasarkan waktu

- 🔍 **Filter, Search & Pagination**
  - Filtering by kategori, status, tanggal
  - Pagination untuk listing produk & transaksi

---

## 📮 API Documentation

- 📘 **Postman Collection**: [Klik di sini](#) *(upload ke GitHub atau share link Postman)*  
- 📚 **Swagger UI**: `http://localhost:8080/swagger/index.html` *(jika tersedia)*  
- ERD & Flow Diagram: Lihat di folder `/docs/`

---

## 🧠 Analisis Kebutuhan Sistem

### Deskripsi Umum
Sistem ini dirancang untuk membantu manajemen inventori dan transaksi penjualan pada sebuah bengkel otomotif keluarga. Selama ini pengecekan stok dilakukan secara manual dan tidak efisien, serta tidak ada informasi digital terkait lokasi penyimpanan barang. Sistem ini bertujuan untuk mengatasi permasalahan tersebut dengan digitalisasi proses stok, transaksi, dan pelacakan lokasi.

### Tujuan Sistem
- Mengelola stok produk secara real-time.
- Meningkatkan efisiensi transaksi penjualan.
- Menyediakan informasi lokasi penyimpanan barang.
- Menyediakan laporan penjualan untuk evaluasi performa produk.

### Aktor Sistem
| Aktor         | Peran                                                                 |
|--------------|-----------------------------------------------------------------------|
| Admin         | Mengelola data produk, pengguna, dan laporan                         |
| Karyawan      | Melayani pesanan pelanggan dan mencatat transaksi                    |

### Kebutuhan Fungsional
- Autentikasi & otorisasi (login, register, JWT).
- CRUD produk (nama, harga, stok, lokasi, kategori).
- CRUD aktivitas transaksi pembelian.
- Validasi stok sebelum melakukan transaksi.
- Update stok otomatis jika transaksi berhasil.
- Laporan penjualan: produk terlaris, stok rendah.
- Filter dan pagination data produk & transaksi.
- Cegah penghapusan produk yang sudah ditransaksikan.

### Kebutuhan Non-Fungsional
- Struktur proyek modular (controller, service, dsb).
- Dokumentasi Swagger/Postman.
- Validasi input user (stok tidak negatif, nama unik, dsb).
- Response API konsisten (dengan metadata & status).
- API RESTful dan siap untuk di-deploy.

### Lingkup & Batasan
**Lingkup**:
- Backend-only REST API.
- Fokus pada inventori bengkel dalam satu cabang.
- Login, transaksi, laporan, dan filter data.

**Batasan**:
- Tidak mendukung pembayaran elektronik atau multi-cabang.
- Tidak ada sistem pengadaan dari supplier.

### Stack Teknologi
- Golang (Gin), GORM, MySQL, JWT
- Swagger + Postman
- Modular folder structure

