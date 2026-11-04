# Portfolio Backend (Go + Fiber + GORM + PostgreSQL)

API untuk portfolio. Frontend (Next.js) mengambil **Featured Projects**, **Pengalaman**, dan **Tools** dari sini,
dan **Admin Panel** di `/admin` dipakai untuk mengelolanya tanpa menyentuh kode.

```text
Pengunjung ─► Frontend (Next.js :3000) ──GET /api/projects|experiences|tools──► Backend (Go :8080) ──► PostgreSQL
                       ▲ cadangan: data/*.ts di frontend bila backend mati
Admin ───────► Admin Panel (:8080/admin) ──POST/PUT/DELETE (token JWT)──► Backend
```

## Menjalankan (pengembangan lokal)

Butuh: Go 1.26+, PostgreSQL. Dari folder `backend`:

```powershell
# 1. Database khusus dev (port 5433, terpisah dari PostgreSQL lain, tanpa password)
.\scripts\dev-db.ps1 start

# 2. Salin pengaturan lalu isi
copy .env.example .env        # contoh dev: DB_PORT=5433, DB_USER=portfolio_dev, DB_PASSWORD=apa-saja

# 3. Jalankan
go run .
```

Backend aktif di `http://localhost:8080`, Admin Panel di `http://localhost:8080/admin`.
Saat pertama jalan: tabel dibuat otomatis, data awal di-*seed* (5 project, 4 pengalaman, 18 tools), dan admin
dibuat dari `ADMIN_USERNAME` / `ADMIN_PASSWORD` bila tabel user masih kosong.

`.\scripts\dev-db.ps1 stop | status | reset` untuk mematikan, mengecek, atau mengosongkan database dev.

> Sudah punya PostgreSQL sendiri? Cukup isi `DB_*` di `.env` (atau `DATABASE_URL`) dan lewati langkah 1.

### Aman untuk database yang sudah berisi data

Migrasi hanya **menambah** tabel/kolom baru (`experiences`, `tools`, dan kolom `category`, `accent`, `image_fit`,
`sort_order` di `projects`). Data lama tidak diubah. Data awal hanya diisi bila tabelnya belum pernah berisi apa pun
(termasuk yang sudah dihapus), jadi item yang sengaja kamu hapus lewat Admin tidak muncul lagi saat backend dinyalakan ulang.
Tetap buat cadangan (`pg_dump`) sebelum deploy.

## Pengaturan (`.env`)

| Variabel | Fungsi |
| --- | --- |
| `DB_HOST` `DB_PORT` `DB_USER` `DB_PASSWORD` `DB_NAME` `DB_SSLMODE` | Koneksi PostgreSQL. Atau pakai satu `DATABASE_URL`. |
| `JWT_SECRET` | Kunci penandatangan token login. **Wajib kuat di produksi.** |
| `ADMIN_USERNAME` `ADMIN_PASSWORD` | Akun admin pertama (dibuat sekali bila belum ada user). |
| `ENABLE_REGISTRATION` | `true` membuka `/api/register`. Biarkan `false`. |
| `ALLOWED_ORIGINS` | Domain frontend yang boleh memanggil API (CORS), pisahkan dengan koma. |
| `SITE_URL` | Alamat website portfolio. Admin Panel memakainya untuk thumbnail gambar dan tautan "Lihat website". |
| `PORT` | Port backend (default 8080). |

`.env` berisi rahasia dan **tidak boleh di-commit** (sudah ada di `.gitignore`).

## Endpoint

| Method | Path | Auth | Fungsi |
| --- | --- | --- | --- |
| GET | `/api/projects` `/api/experiences` `/api/tools` | publik | Daftar data, urut `sort_order` lalu `id`. Kosong = `[]` |
| GET | `/api/projects/:id` | publik | Satu project |
| GET | `/api/site` | publik | `{ "site_url": "..." }` untuk Admin Panel |
| POST | `/api/login` | publik | `{username,password}` → `{token}` (berlaku 24 jam) |
| POST | `/api/register` | publik* | *Nonaktif (404) kecuali `ENABLE_REGISTRATION=true` |
| POST | `/api/projects` `/api/experiences` `/api/tools` | token | Tambah |
| PUT | `/api/projects/:id` `/api/experiences/:id` `/api/tools/:id` | token | Ubah |
| DELETE | `/api/projects/:id` `/api/experiences/:id` `/api/tools/:id` | token | Hapus (soft delete) |

Token dikirim sebagai header `Authorization: Bearer <token>`.

### Bentuk data

Semua isian dirapikan (spasi dibuang, nilai tidak valid diganti default). `ID` pada body diabaikan: tidak bisa menimpa ID asli.

**Project** `{ title*, stack, description, link, image, category, accent, image_fit, sort_order }`
- `accent`: `blue` `cyan` `violet` `orange` `slate` (lainnya → `blue`)
- `image_fit`: `cover` atau `contain` (lainnya → `cover`)
- `image`: path di folder `public/` **frontend** (mis. `/Projects/x.png`) atau URL penuh

**Pengalaman** `{ company*, role, full_role, place, period, categories, tone, initials, logo, summary, skills, points, sort_order }`
- `categories`: `broadcast` `tech` `operations`, pisahkan koma (yang tidak dikenal dibuang)
- `tone`: `blue` `cyan` `orange` `violet` · `skills`: pisahkan koma · `points`: satu tugas per baris
- `summary`: `**kata**` ditampilkan tebal di website

**Tool** `{ name*, icon, group, sort_order }`
- `group`: `frontend` `backend` `design` `dev` (lainnya → `frontend`) · `icon`: mis. `/tools/nextjs.svg`

(`*` = wajib diisi.)

## Struktur

```text
backend/
├── main.go             Titik masuk: koneksi DB, migrasi, seed, route, Admin Panel
├── seed.go             Data awal project + admin, dan backfill kolom baru
├── seed_content.go     Data awal pengalaman & tools (dibaca dari seed_content.json)
├── seed_content.json   Dibuat otomatis: node --no-warnings scripts/gen-seed.mjs (dari data frontend)
├── config/             Koneksi database
├── models/             Struktur tabel + Normalize(): Project, Experience, Tool, User
├── controllers/        Logika endpoint (project, content = CRUD generik pengalaman & tools, auth, site, admin)
├── middleware/         Cek token JWT, CORS
├── routes/             Daftar endpoint
├── utils/              Pembuat/pemeriksa JWT
├── admin/              Admin Panel (index.html, admin.css, js/*.js) — ikut ter-embed di binary
└── scripts/            dev-db.ps1 (PostgreSQL dev), gen-seed.mjs (membuat seed_content.json)
```

Bila data di `frontend/data/` (profile, experience, tools) berubah, jalankan ulang `gen-seed.mjs`
agar data awal ikut berubah, lalu `go test ./...`.

## Admin Panel

`http://localhost:8080/admin` — masuk dengan akun admin. Tiga tab: **Projects**, **Pengalaman**, **Tools**.

- **Baru / Ubah** membuka editor dengan **pratinjau langsung** yang meniru tampilan di website.
- **Hapus** meminta konfirmasi. Esc atau klik di luar menutup editor tanpa menyimpan.
- Perubahan tampil di website dalam ±30 detik (frontend menyimpan cache 30 detik).
- Isi data ditampilkan sebagai teks biasa, jadi tidak bisa menyisipkan HTML/skrip (XSS).
- Kode dipecah per modul di `admin/js/` (`api`, `ui`, `resources`, `editor`, `app`). Menambah jenis konten baru
  cukup menambah satu objek di `resources.js`.

## Test

```powershell
go test ./...          # unit test (database ditiru dengan sqlmock, tidak butuh PostgreSQL)
go vet ./...
```

## Deploy

`Dockerfile` multi-stage siap dipakai. Di produksi isi `DATABASE_URL` (atau `DB_*`), `JWT_SECRET`,
`ADMIN_PASSWORD`, `ALLOWED_ORIGINS` (domain frontend), dan `SITE_URL` (domain website). Jangan memakai nilai contoh
dari `.env` lokal.