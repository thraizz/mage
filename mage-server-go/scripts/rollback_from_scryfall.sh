#!/usr/bin/env bash
#
# Rollback from Scryfall migration back to XMage data
#
# This script:
# 1. Drops the Scryfall compatibility view
# 2. Restores the original XMage cards table
# 3. Optionally drops Scryfall tables
#
# Usage:
#   ./mage-server-go/scripts/rollback_from_scryfall.sh [--drop-scryfall]
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

DROP_SCRYFALL=false
if [[ "${1:-}" == "--drop-scryfall" ]]; then
    DROP_SCRYFALL=true
    shift
fi

# Database connection settings
DB_NAME="${DB_NAME:-mage}"
DB_USER="${DB_USER:-mage}"
DB_PASSWORD="${DB_PASSWORD:-mage}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DOCKER_CONTAINER="${DOCKER_CONTAINER:-mage-postgres}"

# Determine if using Docker or direct connection
MODE="direct"
if docker ps --format '{{.Names}}' | grep -q "^${DOCKER_CONTAINER}$"; then
    MODE="docker"
    echo "Detected running Docker container: $DOCKER_CONTAINER"
fi

psql_exec() {
    local query="$1"
    if [[ "$MODE" == "docker" ]]; then
        docker exec -i "$DOCKER_CONTAINER" env PAGER=cat \
            psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
    else
        PAGER=cat PGPASSWORD="$DB_PASSWORD" \
            psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
    fi
}

echo "==================================="
echo "  Rollback from Scryfall Migration"
echo "==================================="
echo ""
echo "Database: $DB_NAME@$DB_HOST:$DB_PORT"
echo "Mode: $MODE"
echo ""

# Check if backup exists
BACKUP_EXISTS=$(psql_exec "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'cards_xmage_backup');" 2>/dev/null || echo "f")

if [ "$BACKUP_EXISTS" != "t" ]; then
    echo "ERROR: No backup table found (cards_xmage_backup)"
    echo "Cannot rollback without a backup!"
    exit 1
fi

# Confirm with user
echo "This will:"
echo "  1. Drop the 'cards' view (Scryfall compatibility layer)"
echo "  2. Restore 'cards_xmage_backup' as 'cards' table"
if [ "$DROP_SCRYFALL" == "true" ]; then
    echo "  3. Drop 'scryfall_cards' and 'scryfall_sets' tables"
fi
echo ""
read -p "Are you sure you want to rollback? (yes/no): " -r
echo
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "Rollback cancelled"
    exit 0
fi

# Step 1: Drop cards view
echo "Step 1: Dropping cards view..."
psql_exec "DROP VIEW IF EXISTS cards CASCADE;" 2>&1 | grep -v "^$" || true
echo "✓ View dropped"
echo ""

# Step 2: Restore backup table
echo "Step 2: Restoring XMage cards table..."
psql_exec "ALTER TABLE cards_xmage_backup RENAME TO cards;" 2>&1 | grep -v "^$" || true
echo "✓ Table restored"
echo ""

# Step 3: Optionally drop Scryfall tables
if [ "$DROP_SCRYFALL" == "true" ]; then
    echo "Step 3: Dropping Scryfall tables..."
    psql_exec "DROP TABLE IF EXISTS scryfall_cards CASCADE;" 2>&1 | grep -v "^$" || true
    psql_exec "DROP TABLE IF EXISTS scryfall_sets CASCADE;" 2>&1 | grep -v "^$" || true
    echo "✓ Scryfall tables dropped"
    echo ""
fi

# Verify
CARD_COUNT=$(psql_exec "SELECT COUNT(*) FROM cards;" 2>/dev/null || echo "0")

echo "==================================="
echo "  Rollback Complete!"
echo "==================================="
echo ""
echo "Statistics:"
echo "  Cards in restored table: $CARD_COUNT"
echo ""
echo "Next steps:"
echo "  - Test your application with the restored XMage data"
if [ "$DROP_SCRYFALL" == "false" ]; then
    echo "  - Scryfall tables are still present (use --drop-scryfall to remove)"
fi
echo "  - To migrate to Scryfall again: ./scripts/migrate_to_scryfall.sh"
