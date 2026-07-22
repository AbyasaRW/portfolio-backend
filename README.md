# Portfolio Backend API (Go + Fiber + GORM + PostgreSQL)

Backend RESTful API untuk mengelola data portofolio dan menyediakan autentikasi admin. Dibuat menggunakan **Go (Golang)**, **Go Fiber v2**, **GORM ORM**, dan database **PostgreSQL**.

---

## 🚀 Fitur Utama

- **Authentication & Authorization**: Registrasi & Login Admin menggunakan bcrypt password hashing dan JWT (JSON Web Token).
- **CRUD Project Management**: Rest API lengkap untuk Menampilkan, Menambah, Mengubah, dan Menghapus data project.
- **Embedded Admin Panel (SPA)**: Dashboard admin bawaan yang di-embed langsung ke dalam Go binary (`/admin`) menggunakan `embed.FS`.
- **Production-Ready CORS**: Dynamic & secure CORS middleware yang di-whitelist khusus untuk domain Vercel Frontend dan fleksibel lewat environment variable.
- **Database Support**: Mendukung koneksi PostgreSQL lokal maupun cloud deployment (seperti Railway, Render, Fly.io) melalui `DATABASE_URL` dan konfigurasi SSL.
- **Automated Testing**: Dilengkapi dengan White-box Unit Test (HTTP Test `app.Test()`, Mock Database `go-sqlmock`, dan Middleware CORS Test).

---

## 🛠️ Teknologi yang Digunakan

- **Language**: Go (Golang) 1.26+
- **Web Framework**: [Go Fiber v2](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Authentication**: JWT (`golang-jwt/jwt/v5`) & Bcrypt (`golang.org/x/crypto`)
- **Testing**: `stretchr/testify` & `DATA-DOG/go-sqlmock`

---

## 📁 Struktur Folder Project

```text
portfolio-backend/
├── admin/                      # Embedded HTML/JS Single Page Admin Panel
│   └── index.html
├── config/                     # Konfigurasi database & Mock DB Test
│   ├── database.go
│   └── database_test.go
├── controllers/                # Logic Handler / Controller API & Testing
│   ├── admin_handler.go
│   ├── auth_controller.go
│   ├── project_controller.go
│   └── project_controller_test.go
├── middleware/                 # Middleware Auth & CORS
│   ├── auth_middleware.go
│   ├── cors_middleware.go
│   └── cors_test.go
├── models/                     # GORM Database Models (User, Project)
│   ├── project.go
│   └── user.go
├── routes/                     # Definisi Endpoint API
│   └── routes.go
├── utils/                      # Helper JWT Token
│   └── jwt.go
├── .env.example                # Template Environment Variables
├── Dockerfile                  # Multi-stage Docker build config
├── main.go                     # Entrypoint aplikasi
└── seed.go                     # Data seeder awal (Projects & Admin User)
```

---

## 📌 API Endpoints

### 🔓 Public Routes (Tanpa Auth)

| Method | Endpoint | Deskripsi |
| --- | --- | --- |
| `POST` | `/api/register` | Registrasi user admin baru |
| `POST` | `/api/login` | Login admin & mendapatkan JWT token |
| `GET` | `/api/projects` | Mengambil semua daftar project |
| `GET` | `/api/projects/:id` | Mengambil detail project berdasarkan ID |
| `GET` | `/admin` | Mengakses Embedded Dashboard Admin (UI) |

### 🔒 Protected Routes (Butuh Header `Authorization: Bearer <token>`)

| Method | Endpoint | Deskripsi |
| --- | --- | --- |
| `POST` | `/api/projects` | Menambahkan project baru |
| `PUT` | `/api/projects/:id` | Mengubah data project berdasarkan ID |
| `DELETE` | `/api/projects/:id` | Menghapus project berdasarkan ID |

---

## 💻 Panduan Jalankan Secara Lokal

### 1. Prasyarat
- Installed Go 1.20+
- Installed PostgreSQL (atau via Docker)

### 2. Clone & Install Dependencies
```bash
git clone https://github.com/AbyasaRW/portfolio-backend.git
cd portfolio-backend
go mod download
```

### 3. Setup Environment Variables
Buat file `.env` di root project berdasarkan `.env.example`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=portfolio_db
DB_SSLMODE=disable

JWT_SECRET=supersecretjwtkey123
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000
```

### 4. Jalankan Aplikasi
```bash
go run main.go
```
*Aplikasi akan otomatis melakukan AutoMigrate tabel & Seeding data default jika tabel masih kosong (Default Login: `username: admin`, `password: admin123`).*

---

## 🧪 Jalankan Unit Test (White-Box Testing)

Untuk menjalankan seluruh rangkaian pengujian otomatis:

```bash
# Run all tests
go test ./... -v

# Run tests with coverage summary
go test ./... -cover
```

---

## 🌐 Deployment (Railway)

Project ini telah dilengkapi dengan `Dockerfile` multi-stage build yang siap di-deploy ke Railway/Render/Fly.io.

1. Hubungkan repository GitHub ke Railway.
2. Tambahkan Plugin Database **PostgreSQL** di Railway.
3. Tambahkan Environment Variable di Railway:
   - `JWT_SECRET` : Secret key produksi kamu.
   - `ALLOWED_ORIGINS` : `https://abyasa-portfolio.vercel.app` (domain Vercel frontend kamu).
4. Railway secara otomatis membaca `Dockerfile` dan menyambungkan `DATABASE_URL`.
