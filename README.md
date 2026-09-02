# Article Management Service - Sharing Vision Backend Test

Microservice REST API untuk pengelolaan data artikel (*Post Article*) yang dibangun menggunakan **Go (Golang)**, framework **Gin**, ORM **GORM**, dan basis data **MySQL**.

---

## Daftar Isi
- [Teknologi](#-teknologi)
- [Struktur Proyek](#-struktur-proyek)
- [Spesifikasi Database](#-spesifikasi-database)
- [Persyaratan Sistem](#-persyaratan-sistem-prerequisites)
- [Panduan Setup & Instalasi](#-panduan-setup--instalasi)
- [Panduan Migrasi Database](#-panduan-migrasi-database)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [Menjalankan Pengujian (Automated Tests)](#-menjalankan-pengujian-automated-tests)
- [Dokumentasi API Endpoint](#-dokumentasi-api-endpoint)
- [Aturan Validasi Input](#-aturan-validasi-input)
- [Postman Collection](#-postman-collection)

---

## Teknologi

- **Bahasa Pemrograman**: Go (Golang) v1.22+
- **HTTP Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database & ORM**: MySQL 8.0 & [GORM](https://gorm.io/)
- **Database Migration**: [golang-migrate](https://github.com/golang-migrate/migrate) & GORM AutoMigrate
- **Data Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Containerization**: Docker & Docker Compose

---

## Struktur Proyek

```text
backend-sharing-vision-test/
├── cmd/
│   ├── api/
│   │   └── main.go                  # Entry point REST API server
│   └── migrate/
│       └── main.go                  # CLI runner untuk migrasi database
├── config/
│   └── config.go                    # Pengelolaan konfigurasi environment
├── database/
│   ├── manual_create_posts.sql      # Skrip DDL SQL pembuatan tabel manual
│   └── migrations/
│       ├── 000001_create_posts_table.up.sql    # Migrasi UP
│       └── 000001_create_posts_table.down.sql  # Migrasi DOWN (Rollback)
├── internal/
│   ├── dto/
│   │   ├── request/                 # Struct payload request & validasi
│   │   └── response/                # Struct response API
│   ├── handler/
│   │   └── http/                    # Controller / HTTP Request Handler & Test
│   ├── middleware/                  # Middleware (CORS, Logger, Recovery, Rate Limiting)
│   ├── models/                      # GORM Model database
│   ├── repository/                  # Data access layer (CRUD MySQL)
│   ├── routes/                      # Definisi routing endpoint Gin
│   └── service/                     # Business logic layer & Test
├── pkg/
│   ├── database/
│   │   └── mysql/                   # Koneksi DB & helper migrasi MySQL
│   └── utils/                       # Validator & Response helpers & Test
├── postman/
│   ├── Article_API.postman_collection.json    # Postman Collection lengkap
│   └── Article_API.postman_environment.json   # Postman Environment
├── .env.example                     # Template konfigurasi environment
├── docker-compose.yml               # Konfigurasi container Docker MySQL & API
├── Dockerfile                       # Dockerfile untuk aplikasi Go
├── go.mod                           # Go modules manifest
└── README.md                        # Dokumentasi proyek
```

---

## Spesifikasi Database

Nama Database: **`article`**  
Nama Tabel: **`posts`**

| No | Nama Kolom | Tipe Data | Keterangan |
|---|---|---|---|
| 1 | `id` | INT | Auto increment, Primary Key |
| 2 | `title` | VARCHAR(200) | Not Null |
| 3 | `content` | TEXT | Not Null |
| 4 | `category` | VARCHAR(100) | Not Null |
| 5 | `created_date` | TIMESTAMP | Default `CURRENT_TIMESTAMP` |
| 6 | `updated_date` | TIMESTAMP | Default `CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |
| 7 | `status` | VARCHAR(100) | `Publish` \| `Draft` \| `Thrash` |

---

## Persyaratan Sistem (Prerequisites)

Sebelum memulai, pastikan perangkat Anda telah terpasang:
- **Go**: Versi 1.22 atau lebih baru ([Download Go](https://go.dev/dl/))
- **MySQL Server**: Versi 8.0+ ([Download MySQL](https://dev.mysql.com/downloads/)) atau **Docker** ([Download Docker Desktop](https://www.docker.com/products/docker-desktop/))
- **Git**

---

## Panduan Setup & Instalasi

### 1. Clone Repository
```bash
git clone https://github.com/teukufaandii/backend-sharing-vision-test.git
cd backend-sharing-vision-test
```

### 2. Konfigurasi File `.env`
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
*(Pada Windows PowerShell gunakan: `Copy-Item .env.example .env`)*

Sesuaikan nilai environment di dalam `.env` sesuai pengaturan database lokal Anda:
```env
PORT=8080
ENVIRONMENT=development

# Konfigurasi MySQL
DATABASE_URL=root:password@tcp(localhost:3306)/article?charset=utf8mb4&parseTime=True&loc=Local
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=article

# Path file migrasi SQL
MIGRATIONS_PATH=database/migrations
```

### 3. Download Dependencies
```bash
go mod download
```

---

## Panduan Migrasi Database

Tersedia beberapa cara untuk menjalankan migrasi skema tabel `posts`:

### Opsi A: Menjalankan Migrasi via Go CLI Tool (Direkomendasikan)
Jalankan perintah berikut untuk mengeksekusi file migrasi SQL (`.up.sql`):
```bash
# Menjalankan migrasi UP
go run cmd/migrate/main.go up

# Menjalankan rollback migrasi DOWN (1 langkah)
go run cmd/migrate/main.go down

# Menjalankan GORM AutoMigrate
go run cmd/migrate/main.go auto
```

### Opsi B: Eksekusi SQL Manual
Jika ingin membuat skema tabel secara manual langsung ke MySQL server:
```bash
mysql -u root -p < database/manual_create_posts.sql
```

---

## Menjalankan Aplikasi

### Opsi 1: Menjalankan Langsung dengan Go
Pastikan service MySQL sudah aktif, lalu jalankan:
```bash
go run cmd/api/main.go
```
Server akan aktif di: `http://localhost:8080`

### Opsi 2: Menjalankan Menggunakan Docker Compose
Docker Compose akan otomatis menyiapkan container MySQL 8.0 (dengan inisialisasi database `article`) dan container API Go:
```bash
docker-compose up --build -d
```
Untuk menghentikan container:
```bash
docker-compose down
```

---

## Menjalankan Pengujian (Automated Tests)

Jalankan seluruh pengujian unit & integrasi untuk memverifikasi validasi dan endpoint REST:
```bash
go test -v ./...
```

Contoh output pengujian yang berhasil:
```text
=== RUN   TestArticleEndpoints
--- PASS: TestArticleEndpoints (0.00s)
=== RUN   TestArticleService
--- PASS: TestArticleService (0.00s)
=== RUN   TestValidationRules
--- PASS: TestValidationRules (0.00s)
PASS
```

---

## Dokumentasi API Endpoint

| No | Endpoint / URL | Method | Request Body | Response Body | Status HTTP | Deskripsi |
|---|---|---|---|---|---|---|
| 1 | `/article/` | `POST` | `JSON Article` | `{}` | `201 Created` | Membuat article baru |
| 2 | `/article/<limit>/<offset>` | `GET` | - | `Array Article` | `200 OK` | Menampilkan seluruh article dengan paging |
| 3 | `/article/<id>` | `GET` | - | `JSON Article` | `200 OK` / `404 Not Found` | Menampilkan detail article berdasarkan ID |
| 4 | `/article/<id>` | `PUT` / `PATCH` / `POST` | `JSON Article` | `{}` | `200 OK` / `404 Not Found` | Mengubah data article berdasarkan ID |
| 5 | `/article/<id>` | `DELETE` / `POST` | - | `{}` | `200 OK` / `404 Not Found` | Menghapus data article berdasarkan ID |
| 6 | `/health` | `GET` | - | `JSON Health Status` | `200 OK` | Cek status server & runtime |

---

## Aturan Validasi Input

Sebelum artikel disimpan atau diubah, payload JSON harus memenuhi kriteria berikut:
1. **`title`**: Wajib diisi (*required*), minimal **20 karakter**, maksimal 200 karakter.
2. **`content`**: Wajib diisi (*required*), minimal **200 karakter**.
3. **`category`**: Wajib diisi (*required*), minimal **3 karakter**, maksimal 100 karakter.
4. **`status`**: Wajib diisi (*required*), harus memilih salah satu dari: **`Publish`**, **`Draft`**, atau **`Thrash`** (case-insensitive).

---

## Postman Collection

File koleksi Postman telah disediakan di folder [`postman/`](./postman):
1. **Koleksi Request**: [`postman/Article_API.postman_collection.json`](./postman/Article_API.postman_collection.json)
2. **Environment Variable**: [`postman/Article_API.postman_environment.json`](./postman/Article_API.postman_environment.json)

### Cara Menggunakan Postman:
1. Buka aplikasi **Postman**.
2. Klik tombol **Import** (di pojok kiri atas).
3. Pilih dan impor kedua file di atas (`postman/Article_API.postman_collection.json` dan `postman/Article_API.postman_environment.json`).
4. Pilih environment **Article API - Local Environment**.
5. Anda dapat langsung menjalankan seluruh skenario (Positive Case & Negative Validation Test Case).
