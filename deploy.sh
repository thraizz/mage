#!/bin/bash
# Deploy script for MAGE Go server and web client
# Deploys to remote server via SSH using docker-compose
#
# Usage:
#   ./deploy.sh              # Deploy both frontend and backend
#   ./deploy.sh --frontend-only   # Deploy only the frontend

set -e  # Exit on error

# Parse command-line arguments
FRONTEND_ONLY=false
if [[ "$1" == "--frontend-only" ]]; then
    FRONTEND_ONLY=true
fi

# Configuration
REMOTE_USER="hkdebiandocker"
REMOTE_HOST="192.168.178.24"
REMOTE_PATH="gomage"  # Will be resolved to $HOME/gomage on remote
DOCKER_COMPOSE_FILE="docker-compose.prod.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== MAGE Deployment Script ===${NC}"
echo -e "Remote: ${YELLOW}${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}${NC}"
if [[ "$FRONTEND_ONLY" == true ]]; then
    echo -e "Mode: ${YELLOW}Frontend Only${NC}"
else
    echo -e "Mode: ${YELLOW}Full Deployment (Frontend + Backend)${NC}"
fi
echo ""

# Check if we can connect to the remote server
echo -e "${GREEN}[1/6] Checking SSH connection...${NC}"
if ! ssh -o ConnectTimeout=10 "${REMOTE_USER}@${REMOTE_HOST}" "echo 'Connection successful'"; then
    echo -e "${RED}ERROR: Cannot connect to ${REMOTE_USER}@${REMOTE_HOST}${NC}"
    exit 1
fi
echo -e "${GREEN}✓ SSH connection successful${NC}"
echo ""

# Create remote directory structure
echo -e "${GREEN}[2/6] Creating remote directory structure...${NC}"
ssh "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p ${REMOTE_PATH}"
echo -e "${GREEN}✓ Remote directory created${NC}"
echo ""

# Copy necessary files to remote server
echo -e "${GREEN}[3/6] Copying files to remote server...${NC}"

if [[ "$FRONTEND_ONLY" == false ]]; then
    echo "  - Copying docker-compose.prod.yml..."
    scp "${DOCKER_COMPOSE_FILE}" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/docker-compose.yml"

    echo "  - Copying mage-server-go directory..."
    rsync -avz --progress \
        --delete \
        --chmod=Du=rwx,Dgo=rx,Fu=rw,Fgo=r \
        --exclude='bin' \
        --exclude='*.log' \
        --exclude='.git' \
        --exclude='tmp' \
        ./mage-server-go/ "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/mage-server-go/"
fi

echo "  - Copying mage-client-web directory..."
rsync -avz --progress \
    --delete \
    --chmod=Du=rwx,Dgo=rx,Fu=rw,Fgo=r \
    --exclude='node_modules' \
    --exclude='.svelte-kit' \
    --exclude='build' \
    --exclude='.git' \
    --exclude='*.log' \
    ./mage-client-web/ "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/mage-client-web/"

echo -e "${GREEN}✓ Files copied successfully${NC}"
echo ""

# Build new images (while old containers are still running)
echo -e "${GREEN}[4/6] Building new Docker images...${NC}"
if [[ "$FRONTEND_ONLY" == true ]]; then
    ssh "${REMOTE_USER}@${REMOTE_HOST}" "cd ${REMOTE_PATH} && docker compose build mage-client"
else
    ssh "${REMOTE_USER}@${REMOTE_HOST}" "cd ${REMOTE_PATH} && docker compose build"
fi
echo -e "${GREEN}✓ Images built successfully${NC}"
echo ""

# Recreate and restart containers with new images (rolling restart)
echo -e "${GREEN}[5/6] Restarting containers with new images...${NC}"
if [[ "$FRONTEND_ONLY" == true ]]; then
    ssh "${REMOTE_USER}@${REMOTE_HOST}" "cd ${REMOTE_PATH} && docker compose up -d mage-client"
else
    ssh "${REMOTE_USER}@${REMOTE_HOST}" "cd ${REMOTE_PATH} && docker compose up -d"
fi
echo -e "${GREEN}✓ Containers restarted${NC}"
echo ""

# Show status
echo -e "${GREEN}[6/6] Checking container status...${NC}"
ssh "${REMOTE_USER}@${REMOTE_HOST}" "cd ${REMOTE_PATH} && docker compose ps"
echo ""

echo -e "${GREEN}=== Deployment Complete ===${NC}"
echo ""
echo -e "${YELLOW}Services are running at:${NC}"
echo -e "  Frontend: http://${REMOTE_HOST}:38216"
if [[ "$FRONTEND_ONLY" == false ]]; then
    echo -e "  Backend gRPC: http://${REMOTE_HOST}:17171"
    echo -e "  Backend WebSocket: ws://${REMOTE_HOST}:17179"
    echo -e "  Backend Health: http://${REMOTE_HOST}:8080/health"
    echo -e "  Backend Status: http://${REMOTE_HOST}:17171/status"
    echo -e "  Metrics: http://${REMOTE_HOST}:9090/metrics"
fi
echo ""
if [[ "$FRONTEND_ONLY" == false ]]; then
    echo -e "${GREEN}Testing status endpoint...${NC}"
    if curl -sf "http://${REMOTE_HOST}:17171/status" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Status endpoint responding${NC}"
        echo -e "${YELLOW}Response:${NC}"
        curl -s "http://${REMOTE_HOST}:17171/status" | python3 -m json.tool 2>/dev/null || curl -s "http://${REMOTE_HOST}:17171/status"
    else
        echo -e "${RED}✗ Status endpoint not responding${NC}"
    fi
    echo ""
fi
echo -e "${YELLOW}To view logs:${NC}"
if [[ "$FRONTEND_ONLY" == true ]]; then
    echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_PATH} && docker compose logs -f mage-client'"
else
    echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_PATH} && docker compose logs -f'"
fi
echo ""
echo -e "${YELLOW}To stop services:${NC}"
if [[ "$FRONTEND_ONLY" == true ]]; then
    echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_PATH} && docker compose stop mage-client'"
else
    echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_PATH} && docker compose down'"
fi
echo ""
