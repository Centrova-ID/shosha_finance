.PHONY: dev dev-cloud dev-local dev-frontend dev-offline clean help

# ====== DEVELOPMENT ======

# Jalankan semua service sekaligus (online mode)
dev:
	@./scripts/dev.sh

# Jalankan Cloud API saja (development dengan SQLite)
dev-cloud:
	@echo "======================================"
	@echo "  Starting Cloud API on port 3000"
	@echo "======================================"
	@cd backend && DB_DRIVER=sqlite SQLITE_PATH=./shosha_cloud.db PORT=3000 go run cmd/cloud/main.go

# Jalankan Local API saja (dengan sync ke cloud)
dev-local:
	@echo "======================================"
	@echo "  Starting Local API on port 8080"
	@echo "  Sync to: http://localhost:3000"
	@echo "======================================"
	@cd backend && CLOUD_API_URL=http://localhost:3000 SYNC_INTERVAL=10 go run cmd/local/main.go

# Jalankan Frontend saja (akan auto start Local API)
dev-frontend:
	@echo "======================================"
	@echo "  Starting Frontend + Local API"
	@echo "======================================"
	@cd frontend && npm run dev

# Jalankan dalam mode offline (tanpa cloud sync)
dev-offline:
	@echo "======================================"
	@echo "  Starting in OFFLINE mode"
	@echo "  No sync to cloud"
	@echo "======================================"
	@cd frontend && CLOUD_API_URL= npm run dev

# ====== TESTING ======

# Test online: Cloud API + Frontend
test-online:
	@echo "======================================"
	@echo "  ONLINE MODE TEST"
	@echo "======================================"
	@echo "Terminal 1: make dev-cloud"
	@echo "Terminal 2: make dev-frontend"
	@echo ""
	@echo "Pastikan Cloud API jalan dulu di terminal 1"
	@echo "======================================"

# Test offline: Frontend saja
test-offline:
	@echo "======================================"
	@echo "  OFFLINE MODE TEST"
	@echo "======================================"
	@$(MAKE) dev-offline

# ====== BUILD ======

# Build backend
build-backend:
	@echo "Building backend..."
	@cd backend && go build -o bin/local cmd/local/main.go
	@cd backend && go build -o bin/cloud cmd/cloud/main.go

# Build frontend
build-frontend:
	@echo "Building frontend..."
	@cd frontend && npm run build

# Build semua
build: build-backend build-frontend

# ====== UTILS ======

# Clean database files
clean:
	@echo "Cleaning database files..."
	@rm -f backend/shosha_finance.db backend/shosha_cloud.db
	@rm -f ~/.config/shosha-finance/shosha_finance.db
	@echo "Done!"

# Help
help:
	@echo "======================================"
	@echo "  Shosha Finance - Development Commands"
	@echo "======================================"
	@echo ""
	@echo "DEVELOPMENT:"
	@echo "  make dev           - Jalankan semua (Cloud + Local + Frontend)"
	@echo "  make dev-cloud     - Cloud API saja (port 3000)"
	@echo "  make dev-local     - Local API saja (port 8080)"
	@echo "  make dev-frontend  - Frontend + Local API"
	@echo "  make dev-offline   - Frontend tanpa sync (offline mode)"
	@echo ""
	@echo "TESTING:"
	@echo "  make test-online   - Instruksi test online mode"
	@echo "  make test-offline  - Jalankan offline mode"
	@echo ""
	@echo "BUILD:"
	@echo "  make build         - Build backend dan frontend"
	@echo ""
	@echo "UTILS:"
	@echo "  make clean         - Hapus file database"
	@echo ""
