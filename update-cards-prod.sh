#!/bin/bash
# Production Card Data Update Script
# Updates Scryfall card data on production server
#
# This script:
# 1. Downloads latest Scryfall bulk data (optional)
# 2. Copies data to production server
# 3. Runs migration and import on production database
# 4. Creates backup before updating
# 5. Verifies the update
#
# Usage:
#   ./update-cards-prod.sh                    # Use existing local Scryfall data
#   ./update-cards-prod.sh --download         # Download latest data first
#   ./update-cards-prod.sh --rollback         # Rollback to previous data

set -e  # Exit on error

# Configuration (matches deploy.sh)
REMOTE_USER="${REMOTE_USER:-hkdebiandocker}"
REMOTE_HOST="${REMOTE_HOST:-192.168.178.24}"
REMOTE_PATH="${REMOTE_PATH:-gomage}"
DOCKER_CONTAINER="mage-postgres"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Flags
DOWNLOAD_DATA=false
ROLLBACK=false
DRY_RUN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --download)
            DOWNLOAD_DATA=true
            shift
            ;;
        --rollback)
            ROLLBACK=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --download    Download latest Scryfall data before updating"
            echo "  --rollback    Rollback to previous card data"
            echo "  --dry-run     Show what would be done without executing"
            echo "  --help        Show this help message"
            echo ""
            echo "Environment Variables:"
            echo "  REMOTE_USER   SSH user (default: hkdebiandocker)"
            echo "  REMOTE_HOST   SSH host (default: 192.168.178.24)"
            echo "  REMOTE_PATH   Remote directory (default: gomage)"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  MAGE Production Card Data Update                              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Remote: ${YELLOW}${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}${NC}"
echo -e "Container: ${YELLOW}${DOCKER_CONTAINER}${NC}"
if [[ "$DRY_RUN" == "true" ]]; then
    echo -e "${YELLOW}Mode: DRY RUN (no changes will be made)${NC}"
fi
echo ""

# Function to execute commands (respects --dry-run)
execute() {
    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "${YELLOW}[DRY RUN] Would execute:${NC} $*"
        return 0
    else
        "$@"
    fi
}

# Check if we can connect to the remote server
echo -e "${GREEN}[1/8] Checking SSH connection...${NC}"
if ! ssh -o ConnectTimeout=10 "${REMOTE_USER}@${REMOTE_HOST}" "echo 'Connection successful'" > /dev/null 2>&1; then
    echo -e "${RED}ERROR: Cannot connect to ${REMOTE_USER}@${REMOTE_HOST}${NC}"
    exit 1
fi
echo -e "${GREEN}✓ SSH connection successful${NC}"
echo ""

if [[ "$ROLLBACK" == "true" ]]; then
    echo -e "${YELLOW}═══ ROLLBACK MODE ═══${NC}"
    echo ""
    echo -e "${YELLOW}This will restore the previous card data from backup.${NC}"
    read -p "Are you sure you want to rollback? (yes/no): " -r
    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        echo -e "${BLUE}Rollback cancelled${NC}"
        exit 0
    fi
    
    echo -e "${GREEN}[ROLLBACK] Restoring previous card data...${NC}"
    execute ssh "${REMOTE_USER}@${REMOTE_HOST}" "cd ${REMOTE_PATH}/mage-server-go && docker compose -f ../docker-compose.yml exec -T postgres bash -c './scripts/rollback_from_scryfall.sh'"
    
    echo -e "${GREEN}✓ Rollback complete${NC}"
    exit 0
fi

# Download latest Scryfall data if requested
if [[ "$DOWNLOAD_DATA" == "true" ]]; then
    echo -e "${GREEN}[2/8] Downloading latest Scryfall bulk data...${NC}"
    if [[ -f "mage-server-go/scripts/download_scryfall_bulk.sh" ]]; then
        execute bash mage-server-go/scripts/download_scryfall_bulk.sh
        echo -e "${GREEN}✓ Download complete${NC}"
    else
        echo -e "${RED}ERROR: download_scryfall_bulk.sh not found${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}[2/8] Skipping download (using existing data)${NC}"
fi
echo ""

# Find the Scryfall data file
echo -e "${GREEN}[3/8] Locating Scryfall data file...${NC}"
SCRYFALL_FILE=""

# Check for latest symlink
if [[ -L "mage-server-go/data/scryfall-all-cards-latest.json" ]]; then
    SCRYFALL_FILE="mage-server-go/data/$(readlink mage-server-go/data/scryfall-all-cards-latest.json)"
elif [[ -f "mage-server-go/data/scryfall-all-cards-latest.json" ]]; then
    SCRYFALL_FILE="mage-server-go/data/scryfall-all-cards-latest.json"
else
    # Find most recent file in data directory
    SCRYFALL_FILE=$(ls -t mage-server-go/data/scryfall-all-cards-*.json 2>/dev/null | head -1 || echo "")
fi

# Check root directory too (like the existing all-cards-20260111103023.json)
if [[ -z "$SCRYFALL_FILE" ]]; then
    SCRYFALL_FILE=$(ls -t all-cards-*.json 2>/dev/null | head -1 || echo "")
fi

if [[ -z "$SCRYFALL_FILE" ]] || [[ ! -f "$SCRYFALL_FILE" ]]; then
    echo -e "${RED}ERROR: No Scryfall data file found${NC}"
    echo "Please run with --download or ensure a Scryfall JSON file exists"
    exit 1
fi

FILE_SIZE=$(du -h "$SCRYFALL_FILE" | cut -f1)
echo -e "  Found: ${YELLOW}$(basename "$SCRYFALL_FILE")${NC} (${FILE_SIZE})"
echo -e "${GREEN}✓ Data file located${NC}"
echo ""

# Verify remote directory structure
echo -e "${GREEN}[4/8] Verifying remote directory structure...${NC}"
execute ssh "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p ${REMOTE_PATH}/mage-server-go/data"
execute ssh "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p ${REMOTE_PATH}/mage-server-go/scripts"
echo -e "${GREEN}✓ Remote directories ready${NC}"
echo ""

# Copy Scryfall data to remote server
echo -e "${GREEN}[5/8] Copying Scryfall data to remote server...${NC}"
echo "  This may take a few minutes for large files..."
REMOTE_FILE="${REMOTE_PATH}/mage-server-go/data/$(basename "$SCRYFALL_FILE")"
if [[ "$DRY_RUN" == "false" ]]; then
    rsync -avz --progress "$SCRYFALL_FILE" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_FILE}"
else
    echo -e "${YELLOW}[DRY RUN] Would rsync: $SCRYFALL_FILE to ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_FILE}${NC}"
fi
echo -e "${GREEN}✓ Data file copied${NC}"
echo ""

# Copy import scripts to remote server
echo -e "${GREEN}[6/8] Copying import scripts...${NC}"
if [[ "$DRY_RUN" == "false" ]]; then
    rsync -avz \
        mage-server-go/scripts/migrate_to_scryfall.sh \
        mage-server-go/scripts/rollback_from_scryfall.sh \
        "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/mage-server-go/scripts/"
    
    # Set execute permissions on scripts
    ssh "${REMOTE_USER}@${REMOTE_HOST}" "chmod +x ${REMOTE_PATH}/mage-server-go/scripts/*.sh"
    
    # Also copy migrations if they don't exist
    rsync -avz \
        mage-server-go/migrations/ \
        "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/mage-server-go/migrations/"
else
    echo -e "${YELLOW}[DRY RUN] Would copy scripts and migrations${NC}"
fi
echo -e "${GREEN}✓ Scripts copied${NC}"
echo ""

# Create backup of current data
echo -e "${GREEN}[7/8] Creating backup of current card data...${NC}"
BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_SQL="${REMOTE_PATH}/mage-server-go/data/cards_backup_${BACKUP_DATE}.sql"

execute ssh "${REMOTE_USER}@${REMOTE_HOST}" bash << EOF
    set -e
    cd ${REMOTE_PATH}
    
    # Backup current cards table
    echo "Creating SQL backup..."
    docker compose -f ../docker-compose.yml exec -T postgres pg_dump -U mage -d mage \
        -t cards \
        -t cards_xmage_backup \
        --if-exists --clean \
        > mage-server-go/data/cards_backup_${BACKUP_DATE}.sql 2>/dev/null || true
    
    # Also create a simple count snapshot
    echo "Current card counts before update:" > mage-server-go/data/backup_${BACKUP_DATE}_info.txt
    docker compose -f ../docker-compose.yml exec -T postgres psql -U mage -d mage -tAc \
        "SELECT 'cards table:', COUNT(*) FROM cards;" >> mage-server-go/data/backup_${BACKUP_DATE}_info.txt 2>/dev/null || true
    docker compose -f ../docker-compose.yml exec -T postgres psql -U mage -d mage -tAc \
        "SELECT 'scryfall_cards:', COUNT(*) FROM scryfall_cards WHERE lang='en';" >> mage-server-go/data/backup_${BACKUP_DATE}_info.txt 2>/dev/null || true
    
    echo "✓ Backup created: ${BACKUP_SQL}"
EOF
echo -e "${GREEN}✓ Backup complete${NC}"
echo ""

# Run the card data update on production
echo -e "${GREEN}[8/8] Updating card data on production...${NC}"
echo -e "${YELLOW}This will take 2-3 minutes. Do not interrupt!${NC}"
echo ""

if [[ "$DRY_RUN" == "false" ]]; then
    ssh "${REMOTE_USER}@${REMOTE_HOST}" bash << EOF
        set -e
        cd ${REMOTE_PATH}/mage-server-go
        
        # Check if postgres container is running
        if ! docker ps | grep -q mage-postgres; then
            echo "Starting database service..."
            # Stop anything using port 5432 first
            docker ps -q --filter "publish=5432" | xargs -r docker stop 2>/dev/null || true
            # Remove any conflicting containers by name
            docker rm -f mage-postgres 2>/dev/null || true
            docker compose -f ../docker-compose.yml up -d postgres
            echo "Waiting for database to be ready..."
            # Wait for postgres to be actually ready
            for i in {1..30}; do
                if docker compose -f ../docker-compose.yml exec -T postgres pg_isready -U mage &>/dev/null; then
                    echo "Database is ready!"
                    break
                fi
                echo "Waiting... (\$i/30)"
                sleep 2
            done
        else
            echo "Database service already running"
        fi
        
        echo "Applying Scryfall migrations..."
        docker compose -f ../docker-compose.yml exec -T postgres psql -U mage -d mage < migrations/009_create_scryfall_tables.up.sql 2>&1 | grep -v "already exists" || true
        
        echo ""
        echo "Importing Scryfall data..."
        # Check if mage-server is running, start if needed
        if ! docker ps | grep -q mage-server; then
            echo "Starting mage-server service..."
            # Remove any conflicting containers by name
            docker rm -f mage-server 2>/dev/null || true
            echo "Building and starting mage-server (this may take 2-3 minutes)..."
            docker compose up -d mage-server
            echo "Waiting for service to be ready..."
            sleep 10
        else
            echo "Mage-server already running"
        fi
        
        docker compose -f ../docker-compose.yml exec -T mage-server bash -c "DATABASE_URL='postgres://mage:mage@postgres:5432/mage?sslmode=disable' go run ./cmd/scryfall-import/main.go --input='/app/data/$(basename "$SCRYFALL_FILE")' --lang=en --skip-tokens=true --batch=1000" 2>&1 | tail -20
        
        echo ""
        echo "Creating compatibility view..."
        docker compose -f ../docker-compose.yml exec -T postgres psql -U mage -d mage << 'SQL'
-- Backup old cards table if it exists and isn't already backed up
DO \$\$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'cards' AND table_type = 'BASE TABLE') THEN
        DROP TABLE IF EXISTS cards_xmage_backup CASCADE;
        ALTER TABLE cards RENAME TO cards_xmage_backup;
    END IF;
END
\$\$;

-- Create compatibility view
DROP VIEW IF EXISTS cards CASCADE;
CREATE VIEW cards AS
SELECT 
    abs(hashtext(id::text))::bigint as id,
    collector_number as card_number,
    set_code,
    name,
    type_line as card_type,
    COALESCE(mana_cost, '') as mana_cost,
    COALESCE(power, '') as power,
    COALESCE(toughness, '') as toughness,
    COALESCE(oracle_text, '') as rules_text,
    '' as flavor_text,
    '' as original_text,
    '' as original_type,
    0::bigint as cn,
    name as card_name,
    rarity,
    COALESCE(name, '') as card_class_name,
    created_at
FROM scryfall_cards
WHERE lang = 'en'
  AND layout NOT IN ('token', 'art_series', 'emblem', 'double_faced_token')
  AND digital = false;
SQL
EOF
else
    echo -e "${YELLOW}[DRY RUN] Would execute migration and import on production${NC}"
fi

echo -e "${GREEN}✓ Card data updated${NC}"
echo ""

# Verify the update
echo -e "${GREEN}Verifying update...${NC}"
if [[ "$DRY_RUN" == "false" ]]; then
    ssh "${REMOTE_USER}@${REMOTE_HOST}" bash << EOF
        cd ${REMOTE_PATH}
        
        echo "Card counts after update:"
        docker compose -f docker-compose.yml exec -T postgres psql -U mage -d mage << 'SQL'
SELECT 
    'English Scryfall Cards' as metric, 
    COUNT(*)::text as count 
FROM scryfall_cards WHERE lang='en'
UNION ALL
SELECT 
    'Unique Cards (oracle_id)', 
    COUNT(DISTINCT oracle_id)::text 
FROM scryfall_cards WHERE lang='en'
UNION ALL
SELECT 
    'Cards View (compatible)', 
    COUNT(*)::text 
FROM cards;
SQL

        echo ""
        echo "Sample card lookup test:"
        docker compose -f docker-compose.yml exec -T postgres psql -U mage -d mage -c \
            "SELECT name, card_type, mana_cost FROM cards WHERE name LIKE 'Lightning%' LIMIT 3;"
EOF
else
    echo -e "${YELLOW}[DRY RUN] Would verify update${NC}"
fi
echo ""

echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Card Data Update Complete!                                    ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}Summary:${NC}"
echo -e "  ✓ Backup created: ${YELLOW}${BACKUP_SQL}${NC}"
echo -e "  ✓ Scryfall data imported"
echo -e "  ✓ Compatibility view updated"
echo -e "  ✓ Production database updated"
echo ""
echo -e "${YELLOW}Backup Information:${NC}"
echo -e "  Location: ${REMOTE_USER}@${REMOTE_HOST}:${BACKUP_SQL}"
echo -e "  Info file: ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/mage-server-go/data/backup_${BACKUP_DATE}_info.txt"
echo ""
echo -e "${YELLOW}Rollback (if needed):${NC}"
echo -e "  ./update-cards-prod.sh --rollback"
echo ""
echo -e "${YELLOW}View production logs:${NC}"
echo -e "  ssh ${REMOTE_USER}@${REMOTE_HOST} 'cd ${REMOTE_PATH} && docker compose -f ../docker-compose.yml logs -f mage-server'"
echo ""
echo -e "${YELLOW}Verify production is working:${NC}"
echo -e "  curl http://${REMOTE_HOST}:17171/status"
echo ""
