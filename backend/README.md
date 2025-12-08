# Shosha Finance Backend

Backend API untuk aplikasi Shosha Finance dengan arsitektur offline-first dan sinkronisasi cloud.

## Arsitektur

```
┌─────────────────────┐         ┌─────────────────────┐
│   Desktop App 1     │         │   Desktop App 2     │
│   (User A)          │         │   (User B)          │
│   + Local API       │         │   + Local API       │
│   + SQLite          │         │   + SQLite          │
└──────────┬──────────┘         └──────────┬──────────┘
           │                               │
           │  push/pull    ┌───────────┐   │  push/pull
           └──────────────►│ Cloud API │◄──┘
                           │ PostgreSQL│
                           │ (Deploy)  │
                           └───────────┘
```

## Dua API yang Berbeda

### 1. Local API (`cmd/local/main.go`)
- **Lokasi**: Jalan di PC user bersamaan dengan Electron
- **Database**: SQLite (file lokal)
- **Port**: localhost:8080
- **Fungsi**: 
  - CRUD data secara langsung
  - Menyimpan data offline
  - Sync worker untuk push/pull ke cloud

### 2. Cloud API (`cmd/cloud/main.go`)
- **Lokasi**: Deploy ke server (VPS/Railway/Render/dll)
- **Database**: PostgreSQL
- **Port**: 3000 (configurable)
- **Fungsi**:
  - Menyimpan data terpusat
  - Endpoint sync untuk semua desktop app
  - Dashboard admin (opsional)

## Menjalankan Local API

```bash
cd backend

# Set environment (opsional)
export CLOUD_API_URL=https://your-cloud-api.com
export SYNC_INTERVAL=30

# Jalankan
go run cmd/local/main.go
```

### Environment Variables Local API

| Variable | Default | Keterangan |
|----------|---------|------------|
| PORT | 8080 | Port local API |
| SQLITE_PATH | ./shosha_finance.db | Path file SQLite |
| CLOUD_API_URL | - | URL Cloud API untuk sync |
| SYNC_INTERVAL | 30 | Interval sync dalam detik |
| JWT_SECRET | shosha-finance-secret-key-2024 | Secret untuk JWT |

## Deploy Cloud API

### 1. Setup PostgreSQL

Buat database PostgreSQL:

```sql
CREATE DATABASE shosha_finance;
```

### 2. Environment Variables Cloud API

| Variable | Default | Keterangan |
|----------|---------|------------|
| PORT | 3000 | Port Cloud API |
| DB_HOST | localhost | Host PostgreSQL |
| DB_PORT | 5432 | Port PostgreSQL |
| DB_USER | postgres | User PostgreSQL |
| DB_PASS | - | Password PostgreSQL |
| DB_NAME | shosha_finance | Nama database |
| JWT_SECRET | shosha-finance-cloud-secret-2024 | Secret untuk JWT |

### 3. Jalankan Cloud API

```bash
cd backend

# Set environment
export DB_HOST=your-db-host
export DB_USER=your-db-user
export DB_PASS=your-db-password
export DB_NAME=shosha_finance

# Jalankan
go run cmd/cloud/main.go
```

### 4. Deploy dengan Docker

```bash
cd backend
docker build -t shosha-cloud -f Dockerfile .
docker run -p 3000:3000 \
  -e DB_HOST=your-db-host \
  -e DB_USER=your-db-user \
  -e DB_PASS=your-db-password \
  -e DB_NAME=shosha_finance \
  shosha-cloud
```

## Flow Sinkronisasi

1. **User input data** → Simpan ke SQLite lokal
2. **Sync Worker** (setiap 30 detik):
   - **Pull**: Ambil data terbaru dari Cloud API
   - **Push**: Kirim data yang belum sync ke Cloud API
3. **Data tersinkronisasi** → Semua user bisa melihat data yang sama

## API Endpoints

### Local API (localhost:8080)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | /api/v1/auth/login | Login |
| GET | /api/v1/auth/me | Get current user |
| GET | /api/v1/branches | List semua unit |
| POST | /api/v1/branches | Buat unit baru |
| GET | /api/v1/transactions | List transaksi |
| POST | /api/v1/transactions | Buat transaksi |
| GET | /api/v1/dashboard/summary | Ringkasan dashboard |
| GET | /api/v1/system/status | Status online/offline |

### Cloud API (your-domain:3000)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/health | Health check |
| POST | /api/v1/sync/push | Terima data dari local |
| GET | /api/v1/sync/pull | Kirim data ke local |
| POST | /api/v1/auth/login | Login (admin) |
| GET | /api/v1/branches | List unit |
| GET | /api/v1/transactions | List transaksi |
| GET | /api/v1/dashboard/summary | Dashboard |

## Default Users

Aplikasi otomatis membuat user default:

| Username | Password | Role |
|----------|----------|------|
| admin | admin123 | Admin |
| manager | manager123 | Manager |
| staff | staff123 | Staff |

## Struktur Folder

```
backend/
├── cmd/
│   ├── local/          # Local API (jalan di desktop)
│   │   └── main.go
│   └── cloud/          # Cloud API (deploy ke server)
│       └── main.go
├── internal/
│   ├── config/         # Konfigurasi
│   ├── database/       # Koneksi database
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # JWT, CORS
│   ├── models/         # Data models
│   ├── repository/     # Database queries
│   ├── service/        # Business logic
│   └── worker/         # Sync worker
├── Dockerfile          # Untuk deploy cloud
├── go.mod
└── README.md
```
