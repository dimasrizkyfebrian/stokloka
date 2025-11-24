# 📝 Blueprint Proyek: StokLoka

Dokumen ini adalah "sumber kebenaran" (single source of truth) untuk perencanaan dan pengembangan aplikasi StokLoka.

---

## 1. 🎯 Konsep Inti & Peran Pengguna (Aktor)

### Konsep Inti

"StokLoka" adalah sistem backend untuk melacak pergerakan inventaris secara presisi di berbagai lokasi gudang. Tujuannya adalah menyediakan satu "sumber kebenaran" untuk data stok, mulai dari barang diterima dari pemasok hingga dikirim ke pelanggan atau digunakan.

### 👥 Peran Pengguna (Aktor)

1.  **Admin (`admin`):**

    - Akses penuh ke semua modul.
    - Mengelola pengguna (staf, manajer).
    - Mengelola data master (kategori produk, unit, gudang, pemasok).
    - Melihat semua laporan.

2.  **Manajer Gudang (`manager`):**

    - Melakukan persetujuan (approval) untuk Purchase Order (PO) atau permintaan barang.
    - Melihat laporan stok dan pergerakan barang.
    - Tidak bisa mengelola pengguna atau data master kritis.

3.  **Staf Gudang (`staff`):**
    - Pengguna harian.
    - Mencatat barang masuk (dari PO).
    - Mencatat barang keluar (untuk _delivery order_).
    - Melakukan mutasi stok antar gudang.
    - Melakukan _stock opname_ (penyesuaian).

---

## 2. 📦 Fitur Aplikasi (Daftar Fungsionalitas)

### Modul 1: Autentikasi & Pengguna (`auth`)

- **Registrasi Pengguna:** (Hanya oleh Admin) Mendaftarkan Staf atau Manajer baru.
- **Login Pengguna:** Mendapatkan token JWT berdasarkan email & password.
- **Manajemen Peran:** (Hanya oleh Admin) Menetapkan peran (admin, manager, staff) ke pengguna.

### Modul 2: Data Master (`product` & `supplier`)

- **Manajemen Kategori:** CRUD untuk kategori produk.
- **Manajemen Unit:** CRUD untuk unit/satuan (misal: "pcs", "kg", "box").
- **Manajemen Produk:** CRUD untuk data inti produk (SKU, Nama, Deskripsi, Kategori, Unit).
- **Manajemen Pemasok (Supplier):** CRUD untuk data pemasok (Kode, Nama, Kontak, Alamat).

### Modul 3: Transaksi Inti (`inventory` & `po`)

- **Manajemen Gudang:** (Oleh Admin) CRUD untuk lokasi gudang.
- **Manajemen Purchase Order (PO):**
  - Staf membuat draf PO baru ke Supplier tertentu.
  - Manajer me-review dan _approve_ PO.
- **Barang Masuk (Inbound):**
  - Staf mencatat penerimaan barang berdasarkan PO.
- **Barang Keluar (Outbound):**
  - Staf membuat permintaan barang keluar.
- **Mutasi Stok (Transfer):**
  - Pemindahan stok antar gudang.
- **Penyesuaian Stok (Stock Opname):**
  - Penyesuaian jumlah fisik vs sistem.

---

## 3. 🗺️ Pemetaan Awal Skema Tabel (Per Service)

Setiap _service_ memiliki database PostgreSQL-nya sendiri.

### 🏛️ `auth` (Database: `auth_db`)

- `users`: (id, nama, email, password_hash, role_id, created_at)
- `roles`: (id, nama_peran) (misal: "admin", "manager", "staff")

### 📦 `product` (Database: `product_db`)

- `categories`: (id, nama_kategori, deskripsi)
- `units`: (id, nama_unit, singkatan)
- `products`: (id, sku, nama_produk, deskripsi, category_id, unit_id)

### 🚚 `supplier` (Database: `supplier_db`)

- `suppliers`: (id, kode_pemasok, nama_pemasok, kontak_person, email, telepon, alamat)

### 🏦 `inventory` (Database: `inventory_db`)

- `warehouses`: (id, nama_gudang, kode_gudang, alamat)
- `stocks`: (id, product_id, warehouse_id, quantity)
- `stock_movements`: (id, product_id, warehouse_id, quantity, type, reference_id)

### 🧾 `po` (Database: `po_db`)

- `purchase_orders`: (id, po_number, supplier_id, warehouse_id, status, created_by)
- `po_items`: (id, po_id, product_id, quantity)

---

## 4. 🛠️ Tech Stack & Arsitektur (Status Implementasi)

Berikut adalah tumpukan teknologi yang digunakan dan status implementasi saat ini (per 25 November 2025).

### 4.1. Tumpukan Teknologi (Tech Stack)

- **Backend:** Go (Golang) v1.25
- **Framework:** Go Fiber
- **ORM:** GORM
- **Arsitektur:** Microservices (Pola: Database-per-service)
- **Pola Desain:** Clean Architecture (Model, DTO, Repository, Service, Handler, Middleware)
- **Database:** PostgreSQL 16
- **Autentikasi:** JWT (JSON Web Tokens)
- **Keamanan:** Bcrypt (hashing), Rate Limiter, Role-Based Access Control (RBAC)
- **Caching:** Redis 7 (Pola Cache-Aside & Invalidation)
- **Eventing:** RabbitMQ 3.13 (Pola Topic Exchange Publisher)
- **Konkurensi:** Goroutine (non-blocking) untuk _event publishing_.
- **Containerization:** Docker & Docker Compose
- **Hot Reload:** `air-verse` (via `.air.toml`)
- **Testing:** Unit Testing (Go standard lib) & Mocking (`testify/mock`)

### 4.2. Status Progres

- **[✅] Infrastruktur & Lingkungan (Root)**

  - Setup Docker Compose lengkap (Multi-DB, Networking, Volumes, Healthcheck).
  - Skrip inisialisasi Multi-DB PostgreSQL (`init-multiple-dbs.sh`).
  - Konfigurasi Hot Reload (Air) per _service_.

- **[✅] `auth-service` (Database: `auth_db`)**

  - **Status:** Selesai & Teruji.
  - **Fitur:** Login, Register (Admin-only), Health Check.
  - **Keamanan:** Rate Limiter, JWT Middleware, RBAC.
  - **Testing:** Unit Test Service Layer.

- **[✅] `product-service` (Database: `product_db`)**

  - **Status:** Selesai & Teruji.
  - **Fitur:** CRUD lengkap (Product, Category, Unit).
  - **Teknologi:** Redis Caching, RabbitMQ Publisher (Goroutine), RBAC Admin-Only.
  - **Testing:** Unit Test Service Layer.

- **[✅] `supplier-service` (Database: `supplier_db`)**

  - **Status:** Selesai & Teruji.
  - **Fitur:** CRUD lengkap Supplier.
  - **Teknologi:** Redis Caching, RabbitMQ Publisher (Goroutine), RBAC Admin-Only.
  - **Testing:** Unit Test Service Layer (dengan penanganan _race condition_ goroutine).

- **[⬜️] `po-service` (Database: `po_db`)**

  - Belum dimulai. (Target Selanjutnya)

- **[⬜️] `inventory-service` (Database: `inventory_db`)**
  - Belum dimulai.
