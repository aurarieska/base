# latihan-fiber — Tugas Mandiri Modul 3

REST API `students` dengan Go + Fiber + PostgreSQL mengikuti pola repository.

## Prasyarat

- Go 1.21 atau lebih baru — https://go.dev/dl/
- PostgreSQL 16 atau lebih baru — https://www.postgresql.org/download/
- Git

## Variabel Environment

Seluruh kredensial disimpan di berkas `.env` di root project. Berkas itu **tidak** di-commit; lihat `.env.example` sebagai contoh.

| Variabel       | Contoh nilai        | Keterangan                                           |
|----------------|---------------------|------------------------------------------------------|
| `APP_PORT`     | `3000`              | Port HTTP tempat aplikasi mendengarkan              |
| `DB_HOST`      | `localhost`         | Host PostgreSQL                                      |
| `DB_PORT`      | `5432`              | Port PostgreSQL                                      |
| `DB_USER`      | `postgres`          | Nama pengguna PostgreSQL                             |
| `DB_PASSWORD`  | `rahasia`           | Kata sandi — **jangan** di-commit                    |
| `DB_NAME`      | `praktikum_backend` | Nama basis data                                      |
| `DB_SSLMODE`   | `disable`           | `disable` untuk lokal; di server sungguhan `require` |
| `DB_MAX_CONNS` | `10`                | Jumlah maksimum koneksi di pool                      |

## Skema Tabel `students`

| Kolom        | Tipe            | Batasan                                        |
|--------------|-----------------|------------------------------------------------|
| `id`         | `SERIAL`        | PRIMARY KEY, dibuat basis data                 |
| `nim`        | `VARCHAR(10)`   | NOT NULL, **unik** via `UNIQUE INDEX(LOWER(nim))` |
| `name`       | `VARCHAR(100)`  | NOT NULL                                       |
| `grade`      | `INTEGER`       | NOT NULL DEFAULT 0, `CHECK (grade BETWEEN 0 AND 100)` |
| `is_active`  | `BOOLEAN`       | NOT NULL DEFAULT TRUE                          |
| `created_at` | `TIMESTAMPTZ`   | NOT NULL DEFAULT NOW()                         |

**Indeks:**
- `students_nim_lower_key` (UNIQUE) — menjamin keunikan NIM di basis data (bukan di Go).
- `students_name_lower_idx` — mempercepat pencarian `name ILIKE '%...%'`.

## Cara Menyiapkan dari Nol

```powershell
# 1. Clone repositori ini
git clone <url-repo>
cd latihan-fiber

# 2. Buat basis data kosong
psql -U postgres -c "CREATE DATABASE praktikum_backend;"

# 3. Jalankan migrasi
psql -U postgres -d praktikum_backend -f migrations/001_create_students.sql

# 4. Siapkan .env
Copy-Item .env.example .env
# lalu edit .env, isi DB_PASSWORD dengan kata sandi PostgreSQL Anda

# 5. Unduh dependency dan jalankan
go mod tidy
go run .
```

Server akan mendengarkan di `http://localhost:3000`.

## Daftar Endpoint

| Method | Path                     | Keterangan                                      |
|--------|--------------------------|-------------------------------------------------|
| GET    | `/api/v1/health`         | Cek status server **dan** basis data (503 bila DB mati) |
| GET    | `/api/v1/students/`      | Daftar student dengan `?search`, `?is_active`, `?sort`, `?order`, `?page`, `?limit` |
| GET    | `/api/v1/students/:id`   | Ambil satu student                              |
| POST   | `/api/v1/students/`      | Buat student baru (201 + header `Location`)     |
| PUT    | `/api/v1/students/:id`   | Ganti seluruh field                             |
| PATCH  | `/api/v1/students/:id`   | Ubah sebagian field                             |
| DELETE | `/api/v1/students/:id`   | Hapus (204 No Content)                          |

## Pemetaan Status HTTP

| Situasi                                            | Status |
|----------------------------------------------------|--------|
| id tidak ada pada GET/PUT/PATCH/DELETE             | 404    |
| NIM ganda pada POST/PUT/PATCH                      | 409    |
| PostgreSQL mati (dicek oleh `/health`)             | 503    |
| Kegagalan basis data lainnya                       | 500    |

## Struktur Project

```
latihan-fiber/
├── app/
│   ├── model/student.go                struct entitas, request, respons
│   └── repository/student_repository.go  kontrak + implementasi PostgreSQL
├── config/env.go                       pembaca variabel environment
├── database/postgres.go                connection pool + Ping
├── migrations/001_create_students.sql  skema tabel
├── main.go                             perakitan dan daftar route
├── handler.go                          penghubung HTTP ↔ repository
└── helper.go                           reqCtx, parse query, pembalas respons
```