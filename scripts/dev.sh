#!/bin/bash

# Shosha Finance Development Script
# Menjalankan Cloud API, Local API, dan Frontend sekaligus

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  Shosha Finance Development Server  ${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping all services...${NC}"
    kill $CLOUD_PID $LOCAL_PID $FRONTEND_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start Cloud API (SQLite for dev)
echo -e "${GREEN}[1/3] Starting Cloud API on port 3000...${NC}"
cd "$BACKEND_DIR"
DB_DRIVER=sqlite SQLITE_PATH=./shosha_cloud.db PORT=3000 go run cmd/cloud/main.go &
CLOUD_PID=$!
sleep 2

# Start Local API
echo -e "${GREEN}[2/3] Starting Local API on port 8080...${NC}"
cd "$BACKEND_DIR"
CLOUD_API_URL=http://localhost:3000 SYNC_INTERVAL=10 go run cmd/local/main.go &
LOCAL_PID=$!
sleep 2

# Start Frontend
echo -e "${GREEN}[3/3] Starting Frontend...${NC}"
cd "$FRONTEND_DIR"
npm run dev &
FRONTEND_PID=$!

echo ""
echo -e "${BLUE}======================================${NC}"
echo -e "${GREEN}All services started!${NC}"
echo ""
echo -e "  Cloud API:  ${YELLOW}http://localhost:3000${NC}"
echo -e "  Local API:  ${YELLOW}http://localhost:8080${NC}"
echo -e "  Frontend:   ${YELLOW}http://localhost:5173${NC}"
echo ""
echo -e "  Login: ${YELLOW}admin / admin123${NC}"
echo ""
echo -e "${BLUE}======================================${NC}"
echo -e "Press ${RED}Ctrl+C${NC} to stop all services"
echo ""

# Wait for all processes
wait
