#!/bin/bash
# Copy-only script for MAGE Go server and web client
# Copies files to remote server via SSH (no docker build/restart)

set -e  # Exit on error

# Configuration (kept in sync with deploy.sh)
REMOTE_USER="hkdebiandocker"
REMOTE_HOST="192.168.178.24"
REMOTE_PATH="gomage"  # Will be resolved to $HOME/gomage on remote
DOCKER_COMPOSE_FILE="docker-compose.prod.yml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== MAGE Copy Script (no deploy) ===${NC}"
echo -e "Remote: ${YELLOW}${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}${NC}"
echo ""

# Check if we can connect to the remote server
echo -e "${GREEN}[1/3] Checking SSH connection...${NC}"
if ! ssh -o ConnectTimeout=10 "${REMOTE_USER}@${REMOTE_HOST}" "echo 'Connection successful'"; then
    echo -e "${RED}ERROR: Cannot connect to ${REMOTE_USER}@${REMOTE_HOST}${NC}"
    exit 1
fi
echo -e "${GREEN}✓ SSH connection successful${NC}"
echo ""

# Create remote directory structure
echo -e "${GREEN}[2/3] Creating remote directory structure...${NC}"
ssh "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p ${REMOTE_PATH}"
echo -e "${GREEN}✓ Remote directory created${NC}"
echo ""

# Copy necessary files to remote server
echo -e "${GREEN}[3/3] Copying files to remote server...${NC}"
echo "  - Copying docker-compose.prod.yml..."
scp "${DOCKER_COMPOSE_FILE}" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/docker-compose.yml"

echo "  - Copying mage-server-go directory..."
rsync -avz --progress \
    --chmod=Du=rwx,Dgo=rx,Fu=rw,Fgo=r \
    --exclude='bin' \
    --exclude='*.log' \
    --exclude='.git' \
    --exclude='tmp' \
    ./mage-server-go/ "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/mage-server-go/"

echo "  - Copying mage-client-web directory..."
rsync -avz --progress \
    --chmod=Du=rwx,Dgo=rx,Fu=rw,Fgo=r \
    --exclude='node_modules' \
    --exclude='.svelte-kit' \
    --exclude='build' \
    --exclude='.git' \
    --exclude='*.log' \
    ./mage-client-web/ "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/mage-client-web/"

echo -e "${GREEN}✓ Files copied successfully${NC}"
echo ""

echo -e "${GREEN}=== Copy Complete ===${NC}"
echo ""
echo -e "${YELLOW}Next steps (optional):${NC}"
echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_PATH} && docker compose build'"
echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_PATH} && docker compose up -d'"
echo ""
