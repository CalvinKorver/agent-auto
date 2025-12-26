#!/bin/bash

# Development startup script for Lolo AI
# Starts both frontend and backend servers

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting Lolo AI development environment...${NC}\n"

# Function to cleanup background processes on exit
cleanup() {
    echo -e "\n${YELLOW}Shutting down servers...${NC}"
    kill $(jobs -p) 2>/dev/null
    exit
}

trap cleanup SIGINT SIGTERM

# Start backend with hot reload
echo -e "${GREEN}Starting backend server with hot reload...${NC}"
cd backend

# Add Go bin directories to PATH
GOPATH=$(go env GOPATH 2>/dev/null || echo "$HOME/go")
export PATH="$PATH:$GOPATH/bin:$HOME/go/bin"

# Function to find air binary
find_air() {
    if command -v air &> /dev/null; then
        echo "air"
        return 0
    fi
    # Check common Go bin locations
    for path in "$GOPATH/bin/air" "$HOME/go/bin/air"; do
        if [ -f "$path" ]; then
            echo "$path"
            return 0
        fi
    done
    return 1
}

# Check if air is installed
AIR_CMD=$(find_air)
if [ -z "$AIR_CMD" ]; then
    echo -e "${YELLOW}⚠ Air (hot reload) is not installed.${NC}"
    echo -e "${YELLOW}Installing air...${NC}"
    go install github.com/air-verse/air@latest
    
    # Try to find air again after installation
    AIR_CMD=$(find_air)
    if [ -z "$AIR_CMD" ]; then
        echo -e "${RED}✗ Failed to install or find air.${NC}"
        echo -e "${YELLOW}Please install manually: go install github.com/air-verse/air@latest${NC}"
        echo -e "${YELLOW}Make sure $GOPATH/bin or ~/go/bin is in your PATH${NC}"
        echo -e "${YELLOW}Or run without hot reload: go run cmd/server/main.go${NC}"
        cd ..
        exit 1
    fi
    echo -e "${GREEN}✓ Air installed successfully!${NC}"
fi

# Run backend with air (hot reload)
$AIR_CMD &
BACKEND_PID=$!
cd ..

# Give backend a moment to start
sleep 1

# Check if backend process is still running
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${RED}✗ Backend server failed to start!${NC}"
    exit 1
fi

# Start frontend
echo -e "${GREEN}Starting frontend server...${NC}"
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo -e "\n${GREEN}✓ Servers started!${NC}"
echo -e "${BLUE}Frontend:${NC} http://localhost:3000"
echo -e "${BLUE}Backend:${NC}  http://localhost:8080"
echo -e "\n${YELLOW}Press Ctrl+C to stop all servers${NC}\n"

# Wait for all background processes
wait
