#!/bin/bash

# Farben für die Ausgabe
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Speichere PIDs für Cleanup
BACKEND_PID=""
FRONTEND_PID=""

# Cleanup Funktion
cleanup() {
    echo -e "\n${YELLOW}Stopping services...${NC}"
    
    # Stoppe Frontend
    if [ -n "$FRONTEND_PID" ] && kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "Stopping frontend..."
        kill $FRONTEND_PID
    fi
    
    # Stoppe Backend
    if [ -n "$BACKEND_PID" ] && kill -0 $BACKEND_PID 2>/dev/null; then
        echo "Stopping backend..."
        kill $BACKEND_PID
    fi
    
    # Stoppe und entferne Docker Container
    if docker ps -q -f name=statuspage-db >/dev/null; then
        echo "Stopping database..."
        docker stop statuspage-db >/dev/null 2>&1
        docker rm statuspage-db >/dev/null 2>&1
    fi
    
    echo -e "${GREEN}All services stopped${NC}"
    exit 0
}

# Registriere Cleanup für Ctrl+C
trap cleanup INT TERM

echo -e "${GREEN}Starting Statuspage Setup...${NC}"

# Prüfe ob Docker installiert ist
if ! command -v docker &> /dev/null; then
    echo "Docker ist nicht installiert. Bitte installiere Docker und versuche es erneut."
    exit 1
fi

# Prüfe ob die erforderlichen Ports verfügbar sind
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null ; then
        echo "Port $1 wird bereits verwendet. Bitte stelle sicher, dass der Port frei ist."
        exit 1
    fi
}

check_port 5432  # PostgreSQL
check_port 8080  # Backend
check_port 5173  # Frontend

# PostgreSQL mit Docker starten
echo -e "${YELLOW}Starting PostgreSQL...${NC}"

# Erstelle ein Netzwerk für die Anwendung, wenn es noch nicht existiert
if ! docker network ls | grep -q statuspage-network; then
    echo -e "${YELLOW}Creating docker network...${NC}"
    docker network create statuspage-network
fi

# Prüfe, ob der Datenbank-Container existiert und starte ihn ggf. oder erstelle ihn neu
if ! docker ps -a --format '{{.Names}}' | grep -qw "^statuspage-db$"; then
    echo -e "${YELLOW}Database container 'statuspage-db' not found. Creating new one...${NC}"
    docker run --name statuspage-db \
        --network statuspage-network \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_DB=statuspage \
        -p 5432:5432 \
        -d postgres:latest
else
    echo -e "${YELLOW}Database container 'statuspage-db' found.${NC}"
    # Prüfe, ob der Container läuft
    if ! docker ps --format '{{.Names}}' | grep -qw "^statuspage-db$"; then
        echo -e "${YELLOW}Container 'statuspage-db' is not running. Starting it...${NC}"
        docker start statuspage-db
    else
        echo -e "${GREEN}Container 'statuspage-db' is already running.${NC}"
    fi
fi

# Warte bis die Datenbank wirklich bereit ist
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
until docker exec statuspage-db pg_isready -U postgres >/dev/null 2>&1; do
    echo -n "."
    sleep 1
done
echo -e "\n${GREEN}Database is ready!${NC}"

# Erstelle Logs-Verzeichnisse
mkdir -p /workspaces/statuspage-test/logs
mkdir -p /workspaces/statuspage-test/frontend/logs

# Backend Dependencies installieren und starten
echo -e "${YELLOW}Installing backend dependencies...${NC}"
cd /workspaces/statuspage-test/backend
go mod download
go mod tidy

# Frontend Dependencies installieren
echo -e "${YELLOW}Installing frontend dependencies...${NC}"
cd ../frontend
npm install

# Backend und Frontend starten
echo -e "${GREEN}Starting backend and frontend...${NC}"

# Erstelle ein Verzeichnis für die Logs
mkdir -p logs

# Starte Backend
echo -e "${YELLOW}Starting backend server...${NC}"
cd /workspaces/statuspage-test/backend
go run . > /workspaces/statuspage-test/logs/backend.log 2>&1 &
BACKEND_PID=$!

# Warte länger auf das Backend wegen Datenbankinitialisierung
sleep 5
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${RED}Backend failed to start. Check /workspaces/statuspage-test/logs/backend.log for details${NC}"
    cleanup
    exit 1
fi

# Starte Frontend
echo -e "${YELLOW}Starting frontend server...${NC}"
cd /workspaces/statuspage-test/frontend
npm run dev > /workspaces/statuspage-test/logs/frontend.log 2>&1 &
FRONTEND_PID=$!

# Prüfe ob das Frontend erfolgreich gestartet ist
sleep 2
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    echo -e "${RED}Frontend failed to start. Check logs/frontend.log for details${NC}"
    cleanup
    exit 1
fi

echo -e "${GREEN}Setup complete!${NC}"
echo -e "Frontend runs on: ${YELLOW}http://localhost:5173${NC}"
echo -e "Backend runs on:  ${YELLOW}http://localhost:8080${NC}"
echo -e "\nPress Ctrl+C to stop all services"

# Cleanup Funktion
cleanup() {
    echo -e "\n${YELLOW}Stopping services...${NC}"
    
    # Stoppe Frontend
    if [ -n "$FRONTEND_PID" ] && kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "Stopping frontend..."
        kill $FRONTEND_PID
    fi
    
    # Stoppe Backend
    if [ -n "$BACKEND_PID" ] && kill -0 $BACKEND_PID 2>/dev/null; then
        echo "Stopping backend..."
        kill $BACKEND_PID
    fi
    
    # Stoppe und entferne Docker Container
    if docker ps -q -f name=statuspage-db >/dev/null; then
        echo "Stopping database..."
        docker stop statuspage-db >/dev/null 2>&1
        docker rm statuspage-db >/dev/null 2>&1
    fi
    
    echo -e "${GREEN}All services stopped${NC}"
    exit 0
}

# Registriere Cleanup für Ctrl+C
trap cleanup INT TERM

# Halte das Script am Laufen
while true; do sleep 1; done
