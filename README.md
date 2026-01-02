# Course Planner API

Backend REST API sederhana untuk pengelolaan perkuliahan (course planner) berbasis Go, Fiber, dan PostgreSQL.  
Fitur utama:

- Autentikasi user (mahasiswa, dosen, admin) dengan JWT.
- Manajemen data user, course, room, class, dan KRS (struktur model).
- Fitur admin untuk mengelola kelas:
  - Membuat kelas baru.
  - Mengedit (patch) kelas.
  - Melihat detail kelas.
  - Menghapus kelas.
  - Melihat list semua kelas.
- Seed data awal (user, course, room, class).

---

## Teknologi

- **Bahasa**: Go
- **Framework HTTP**: [Fiber](https://github.com/gofiber/fiber)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Autentikasi**: JWT (JSON Web Token)
- **Konfigurasi**: `.env` (opsional, dengan fallback default di code)

---

## Struktur Project

```text
.
├── main.go
├── go.mod / go.sum
├── config/
│   └── config.go          # Koneksi database (PostgreSQL)
├── cmd/
│   └── seed/
│       └── main.go        # Seeder data awal (user, course, room, class)
└── internal/
    ├── models/            # Definisi model GORM
    │   ├── user.go
    │   ├── course.go
    │   ├── room.go
    │   ├── class.go
    │   ├── krs.go
    │   └── krs_item.go
    ├── repository/        # Akses database (repository pattern)
    │   ├── user_repository.go
    │   └── class_repository.go
    ├── service/           # Business logic
    │   ├── auth_service.go
    │   └── class_service.go
    ├── handler/           # HTTP handler / controller
    │   ├── auth_handler.go
    │   └── class_handler.go
    └── router/
        └── router.go      # Definisi route dan middleware
```

---

## Persiapan Lingkungan

### 1. Prasyarat

- Go terinstal (minimal versi yang sesuai dengan `go.mod`).
- PostgreSQL berjalan (lokal atau remote).
- Git (untuk clone repo).

### 2. Clone Repository

```bash
git clone <URL_REPO_GITHUB_ANDA>
cd course-planner-api
```

### 3. Konfigurasi Database

Secara default, koneksi DB di `config/config.go` menggunakan environment variables:

- `DB_HOST` (default: `localhost`)
- `DB_USER` (default: `postgres`)
- `DB_PASSWORD` (default: kosong)
- `DB_NAME` (default: `course_planner`)
- `DB_PORT` (default: `5432`)

Buat file `.env` di root project jika ingin override default:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password_kamu
DB_NAME=course_planner
DB_PORT=5432

JWT_SECRET=supersecretjwt
```

> `JWT_SECRET` digunakan untuk menandatangani token JWT. Jika tidak di-set, code fallback ke `"secret"`.

### 4. Setup Database

- Buat database di PostgreSQL:

```sql
CREATE DATABASE course_planner;
```

---

## Menjalankan Aplikasi

### 1. Instal Dependencies

```bash
go mod tidy
```

### 2. Jalankan Server

```bash
go run main.go
```

Secara default server akan berjalan di:

```text
http://localhost:8080
```

Saat startup, `main.go` akan menjalankan `AutoMigrate` untuk model:

- `User`
- `Course`
- `Room`
- `Class`
- `KRS`
- `KRSItem`

---

## Seeder (Data Awal)

### Menjalankan Seeder

```bash
go run ./cmd/seed
```

Seeder akan:

- Membuat beberapa user:
  - Admin:
    - `admin@example.com` / `password123`
    - `admin2@example.com` / `password123`
  - Dosen:
    - `dosen@example.com` / `password123` (`nidn`: `1234567890`)
    - `dosen2@example.com` / `password123` (`nidn`: `0987654321`)
  - Mahasiswa:
    - `mahasiswa@example.com` / `password123` (`nim`: `2024000001`)
    - `mahasiswa2@example.com` / `password123` (`nim`: `2024000002`)
    - Keduanya di-assign ke `dosen@example.com` sebagai dosen PA.
- Membuat course:
  - `JSI60214` – `Aplikasi Berbasis Layanan` – 3 SKS
  - `JSI60204` – `Tata Kelola` – 3 SKS
- Membuat room:
  - `H1.1`
  - `H1.2`
  - `H1.3`
- Membuat satu data kelas (class) dengan ID statis (harus cocok dengan data di DB):
  - `course_id`: `1c1cf54c-e380-4369-a038-5d2bcf0926c0`
  - `dosen_id`: `cba38bab-3e52-4f06-9bfe-112ae81e32cf`
  - `nama_kelas`: `B`
  - `hari`: `Senin`
  - `jam_mulai`: `08:00`
  - `jam_selesai`: `10:00`
  - `room_id`: `dc72adcf-c666-4753-994a-a26f6e0718d3`
  - `kuota`: `30`
  - `semester_penawaran`: `ganjil`

> Jika ID di atas belum ada di DB, seeder kelas bisa gagal karena foreign key. Silakan sesuaikan ID dengan data yang ada atau ubah seeder agar mengambil ID berdasarkan email/kode/nama.

---

## Autentikasi & JWT

Sistem menggunakan JWT untuk autentikasi:

- Login menghasilkan token JWT.
- Token harus dikirim di header:

```http
Authorization: Bearer <TOKEN>
```

- Middleware JWT di-setup di `router/router.go`.
- `adminOnlyMiddleware` memastikan hanya user dengan `role = "admin"` yang bisa akses route `/api/admin/**`.

---

## API Endpoint

### Prefix

Semua endpoint diawali dengan:

```text
/api
```

### 1. Auth

#### POST `/api/auth/register`

Register user baru (saat ini field dasar saja, bisa dikembangkan).

**Body contoh:**

```json
{
  "name": "User Baru",
  "email": "userbaru@example.com",
  "password": "password123"
}
```

**Response (201 Created):**

```json
{
  "message": "User registered successfully",
  "user": {
    "id": "...",
    "name": "User Baru",
    "email": "userbaru@example.com",
    "role": "",
    "nim": "",
    "nidn": "",
    "dosen_pa_id": null,
    "created_at": "...",
    "updated_at": "..."
  }
}
```

#### POST `/api/auth/login`

Login dan mendapatkan JWT token.

**Body:**

```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "token": "<JWT_TOKEN>",
  "user": {
    "id": "...",
    "name": "Admin Satu",
    "email": "admin@example.com",
    "role": "admin",
    "nim": "",
    "nidn": "",
    "dosen_pa_id": null,
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### 2. Endpoint Protected Sample

#### GET `/api/me`

Contoh endpoint protected sederhana, hanya mengembalikan pesan jika token valid.  
Header:

```http
Authorization: Bearer <TOKEN>
```

---

## Manajemen Kelas (Admin Only)

Semua endpoint ini hanya bisa diakses oleh user dengan `role = "admin"` dan token JWT yang valid.

Base path:

```text
/api/admin/classes
```

### 1. GET `/api/admin/classes`

List semua kelas.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
```

**Response (200 OK):**

```json
[
  {
    "id": "...",
    "course_id": "...",
    "course": {
      "id": "...",
      "kode": "JSI60214",
      "nama": "Aplikasi Berbasis Layanan",
      "sks": 3
    },
    "dosen_id": "...",
    "dosen": {
      "id": "...",
      "name": "Dosen Satu",
      "email": "dosen@example.com",
      "role": "dosen",
      "nim": "",
      "nidn": "1234567890"
    },
    "nama_kelas": "B",
    "hari": "Senin",
    "jam_mulai": "...",
    "jam_selesai": "...",
    "room_id": "...",
    "room": {
      "id": "...",
      "nama": "H1.1"
    },
    "kuota": 30,
    "semester_penawaran": "ganjil"
  }
]
```

### 2. POST `/api/admin/classes`

Membuat kelas baru.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
Content-Type: application/json
```

**Body contoh:**

```json
{
  "course_id": "UUID_COURSE",
  "dosen_id": "UUID_DOSEN",
  "nama_kelas": "B",
  "hari": "Senin",
  "jam_mulai": "08:00",
  "jam_selesai": "10:00",
  "room_id": "UUID_ROOM",
  "kuota": 30,
  "semester_penawaran": "ganjil"
}
```

**Catatan rules:**

- Format `jam_mulai` dan `jam_selesai`: `"HH:MM"`.
- Sistem akan menolak jika:
  - Di hari yang sama,
  - Di ruangan yang sama,
  - Rentang waktu overlap dengan kelas lain yang sudah ada.

Jika bentrok:

```json
{
  "error": "room already used at this time"
}
```

(status `400 Bad Request`).

### 3. GET `/api/admin/classes/:id`

Detail kelas berdasarkan ID.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
```

**Response (200 OK):** sama seperti satu item di list.

Jika tidak ditemukan: `404 Not Found`.

### 4. PATCH `/api/admin/classes/:id`

Update / edit kelas (partial update). Hanya field yang dikirim yang akan diubah.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
Content-Type: application/json
```

**Body contoh (ubah jam & kuota saja):**

```json
{
  "jam_mulai": "09:00",
  "jam_selesai": "11:00",
  "kuota": 40
}
```

**Catatan rules:**

- Semua field optional: `course_id`, `dosen_id`, `nama_kelas`, `hari`, `jam_mulai`, `jam_selesai`, `room_id`, `kuota`, `semester_penawaran`.
- Setelah semua perubahan diterapkan ke nilai final, sistem cek bentrok dengan kelas lain di ruangan & hari yang sama (mengabaikan kelas ini sendiri).
- Jika bentrok → `400 Bad Request` dengan pesan yang sama seperti create.

### 5. DELETE `/api/admin/classes/:id`

Hapus kelas.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
```

**Response:**

- `204 No Content` jika berhasil.

---

## Manajemen Mata Kuliah / Courses (Admin Only)

Base path:

```text
/api/admin/courses
```

### 1. GET `/api/admin/courses`

List semua mata kuliah.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
```

**Response (200 OK):**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "kode": "IF101",
    "nama": "Algoritma dan Pemrograman",
    "sks": 3
  }
]
```

### 2. POST `/api/admin/courses`

Membuat mata kuliah baru.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
Content-Type: application/json
```

**Body contoh:**

```json
{
  "kode": "IF101",
  "nama": "Algoritma dan Pemrograman",
  "sks": 3
}
```

**Response (201 Created):**

```json
{
  "message": "Course created successfully",
  "course": {
    "id": "...",
    "kode": "IF101",
    "nama": "Algoritma dan Pemrograman",
    "sks": 3
  }
}
```

### 3. GET `/api/admin/courses/:id`

Detail mata kuliah berdasarkan ID.

**Response (200 OK):** sama seperti satu item di list.

Jika tidak ditemukan: `404 Not Found`.

### 4. PATCH `/api/admin/courses/:id`

Update mata kuliah (partial update).

**Body contoh:**

```json
{
  "kode": "IF102",
  "nama": "Struktur Data",
  "sks": 4
}
```

Semua field optional.

### 5. DELETE `/api/admin/courses/:id`

Hapus mata kuliah.

**Response:**

- `204 No Content` jika berhasil.

---

## Manajemen Ruangan / Rooms (Admin Only)

Base path:

```text
/api/admin/rooms
```

### 1. GET `/api/admin/rooms`

List semua ruangan.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
```

**Response (200 OK):**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "nama": "Lab Komputer A"
  }
]
```

### 2. POST `/api/admin/rooms`

Membuat ruangan baru.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
Content-Type: application/json
```

**Body contoh:**

```json
{
  "nama": "Lab Komputer A"
}
```

**Response (201 Created):**

```json
{
  "message": "Room created successfully",
  "room": {
    "id": "...",
    "nama": "Lab Komputer A"
  }
}
```

### 3. GET `/api/admin/rooms/:id`

Detail ruangan berdasarkan ID.

**Response (200 OK):** sama seperti satu item di list.

Jika tidak ditemukan: `404 Not Found`.

### 4. PATCH `/api/admin/rooms/:id`

Update ruangan (partial update).

**Body contoh:**

```json
{
  "nama": "Lab Komputer B"
}
```

### 5. DELETE `/api/admin/rooms/:id`

Hapus ruangan.

**Response:**

- `204 No Content` jika berhasil.

---

## Manajemen Dosen (Admin Only)

Base path:

```text
/api/admin/dosen
```

### 1. GET `/api/admin/dosen`

List semua dosen.

**Header:**

```http
Authorization: Bearer <ADMIN_TOKEN>
```

**Response (200 OK):**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Dr. John Doe",
    "email": "johndoe@university.ac.id",
    "role": "dosen",
    "nidn": "0012345678",
    "created_at": "...",
    "updated_at": "..."
  }
]
```

### 2. GET `/api/admin/dosen/:id`

Detail dosen berdasarkan ID.

**Response (200 OK):** sama seperti satu item di list.

Jika tidak ditemukan: `404 Not Found`.

### 3. PATCH `/api/admin/dosen/:id`

Update data dosen (partial update).

**Body contoh:**

```json
{
  "name": "Dr. Jane Doe",
  "nidn": "0087654321"
}
```

**Response (200 OK):**

```json
{
  "message": "Dosen updated successfully",
  "dosen": {
    "id": "...",
    "name": "Dr. Jane Doe",
    "email": "...",
    "role": "dosen",
    "nidn": "0087654321",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

---

## Model Data (Ringkas)

### Users (`internal/models/user.go`)

Kolom utama:

- `id` (UUID, PK)
- `name`
- `email` (unik)
- `password`
- `role` (`mahasiswa`, `dosen`, `admin`)
- `nim` (untuk mahasiswa)
- `nidn` (untuk dosen)
- `dosen_pa_id` (UUID, refer ke `users.id` dosen PA)
- `created_at`, `updated_at` (`timestamp without time zone`)

### Courses (`internal/models/course.go`)

- `id` (UUID, PK)
- `kode`
- `nama`
- `sks`

### Rooms (`internal/models/room.go`)

- `id` (UUID, PK)
- `nama`

### Classes (`internal/models/class.go`)

- `id` (UUID, PK)
- `course_id` (FK → `courses.id`)
- `dosen_id` (FK → `users.id`)
- `nama_kelas`
- `hari`
- `jam_mulai`, `jam_selesai` (`timestamp without time zone`)
- `room_id` (FK → `rooms.id`)
- `kuota`
- `semester_penawaran` (`ganjil` / `genap`)

### KRS & KRS Items

Struktur sudah disiapkan di model (`krs.go`, `krs_item.go`) untuk fitur KRS:

- `krs`:
  - `mahasiswa_id`
  - `semester`
  - `status` (`draft`, `menunggu_verifikasi`, `disetujui`, `ditolak`)
  - `catatan_dosen`
  - `created_at`, `verified_at`
- `krs_items`:
  - `krs_id`
  - `class_id`
  - `status` (`aktif`, `diajukan_batal`, `batal`)
  - `created_at`, `diajukan_batal_at`, `dibatalkan_at`

---

## Referensi Buku / Books (External API - UAS Feature)

> **[UAS]** Fitur ini mengkonsumsi **Google Books API** (Public API) dengan autentikasi **API Key**.

### Konfigurasi

Tambahkan API Key di file `.env`:

```env
GOOGLE_BOOKS_API_KEY=your_google_books_api_key_here
```

> Dapatkan API Key gratis di [Google Cloud Console](https://console.cloud.google.com/) → APIs & Services → Credentials → Create Credentials → API Key → Enable "Books API"

### API yang Dikonsumsi

| Provider | Base URL | Auth Type |
|----------|----------|-----------|
| Google Books API | `https://www.googleapis.com/books/v1/volumes` | API Key |

Base path:

```text
/api/books
```

### 1. GET `/api/books?query=...`

Cari buku berdasarkan kata kunci.

**Header:**

```http
Authorization: Bearer <TOKEN>
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | ✅ | Kata kunci pencarian (contoh: "algoritma", "pemrograman") |
| `maxResults` | int | ❌ | Jumlah maksimal hasil (default: 10, max: 40) |

**Response (200 OK):**

```json
{
  "total_items": 150,
  "books": [
    {
      "id": "zyTCAlFPjgYC",
      "title": "Introduction to Algorithms",
      "authors": ["Thomas H. Cormen", "Charles E. Leiserson"],
      "publisher": "MIT Press",
      "published_date": "2009",
      "description": "Buku panduan lengkap tentang algoritma...",
      "thumbnail": "https://books.google.com/books/content?id=...",
      "info_link": "https://books.google.com/books?id=...",
      "isbn_10": "0262033844",
      "isbn_13": "9780262033848"
    }
  ]
}
```

### 2. GET `/api/books/:id`

Detail buku berdasarkan Google Books Volume ID.

**Header:**

```http
Authorization: Bearer <TOKEN>
```

**Response (200 OK):**

```json
{
  "id": "zyTCAlFPjgYC",
  "title": "Introduction to Algorithms",
  "subtitle": "Third Edition",
  "authors": ["Thomas H. Cormen", "Charles E. Leiserson"],
  "publisher": "MIT Press",
  "published_date": "2009-07-31",
  "description": "This book covers a broad range of algorithms in depth...",
  "page_count": 1312,
  "categories": ["Computers", "Programming"],
  "average_rating": 4.5,
  "ratings_count": 127,
  "language": "en",
  "thumbnail": "https://books.google.com/books/content?id=...",
  "preview_link": "https://books.google.com/books?id=...",
  "info_link": "https://books.google.com/books?id=...",
  "isbn_10": "0262033844",
  "isbn_13": "9780262033848"
}
```

Jika tidak ditemukan: `404 Not Found`.

---

## Pengembangan Lanjutan (Ide)

Beberapa pengembangan yang bisa dilakukan selanjutnya:

- ✅ ~~Endpoint CRUD untuk `courses`, `rooms` (admin).~~ (Sudah diimplementasikan)
- ✅ ~~Endpoint manajemen dosen (admin).~~ (Sudah diimplementasikan)
- ✅ ~~Dokumentasi API dengan Swagger / OpenAPI.~~ (Tersedia di `/docs`)
- ✅ ~~Integrasi External API (UAS).~~ (Google Books API - Sudah diimplementasikan)
- Endpoint KRS:
  - Mahasiswa buat KRS.
  - Tambah/hapus class ke KRS.
  - Dosen PA verifikasi / tolak KRS.
- Pagination dan filtering untuk list classes (by hari, dosen, course, semester).
- Validasi tambahan (misal: dosen hanya boleh ajar di jam tertentu, max SKS per mahasiswa, dll).

---

## Lisensi

Silakan sesuaikan bagian ini dengan lisensi yang kamu mau (MIT, Apache 2.0, dsb).  
Contoh:

```text
MIT License – silakan modifikasi dan gunakan sesuai kebutuhan.
```

