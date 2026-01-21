# Automatic Scryfall Data Import

The production Docker setup now automatically imports Scryfall card data on first startup, eliminating the need for manual import steps.

## How It Works

### 1. Database Initialization (Postgres Container)
When the postgres container starts for the first time with an empty data directory:
- `docker-initdb/00_migrate_and_import.sh` runs all migrations
- Creates `scryfall_cards` table via migration 009
- Creates `schema_migrations` tracking table

### 2. Automatic Import (Mage-Server Container)
The mage-server container includes a `docker-entrypoint.sh` script that runs before the server starts:

```bash
[1/4] Wait for database to be ready
[2/4] Check if scryfall_cards table is empty
      - If empty: Run scryfall-import tool automatically
      - If populated: Skip import
[3/4] Create/update compatibility view (cards -> scryfall_cards)
[4/4] Start MAGE server
```

### 3. Built-in Tools
The mage-server Docker image now includes:
- `/app/mage-server` - Main server binary
- `/app/scryfall-import` - Scryfall import tool
- `/app/data/` - Mounted from host (contains Scryfall JSON files)
- `/app/migrations/` - Database migrations
- `postgresql-client` - For running SQL commands

## Data File Setup

Ensure you have Scryfall data in `mage-server-go/data/`:

```bash
# Download latest Scryfall data
cd mage-server-go
./scripts/download_scryfall_bulk.sh

# Or use existing file
ls -lh data/scryfall-all-cards-*.json
```

The entrypoint script searches for:
1. `data/scryfall-all-cards-latest.json` (symlink or file)
2. Most recent `data/scryfall-all-cards-*.json` file

## First Startup

```bash
# Fresh deployment (no data directory)
docker compose -f docker-compose.prod.yml up -d

# Container startup sequence:
# 1. postgres: Runs migrations (30 seconds)
# 2. mage-server: Waits for postgres
# 3. mage-server: Detects empty scryfall_cards table
# 4. mage-server: Imports Scryfall data (2-3 minutes)
# 5. mage-server: Creates compatibility view
# 6. mage-server: Starts server

# Monitor progress
docker logs -f mage-server
```

Expected output:
```
======================================
MAGE Server - Docker Entrypoint
======================================

[1/4] Waiting for database to be ready...
✓ Database is ready

[2/4] Checking Scryfall card data...
⚠ No Scryfall card data found. Starting import...
  Found: scryfall-all-cards-20260119230948.json (2.3G)
  This will take 2-3 minutes. Please wait...

[Import logs...]

✓ Scryfall import complete

[3/4] Creating compatibility view...
✓ Compatibility view created

[4/4] Verifying setup...
✓ Cards available through compatibility view: 95000

======================================
Starting MAGE Server
======================================
```

## Subsequent Startups

On subsequent restarts, the entrypoint detects existing data and skips import:

```
[2/4] Checking Scryfall card data...
✓ Found 95000 English Scryfall cards

[3/4] Creating compatibility view...
✓ Compatibility view created
```

Startup time: ~10 seconds (no import)

## Healthcheck Configuration

The healthcheck is configured with extended startup time for first boot:
```yaml
healthcheck:
  start_period: 180s  # 3 minutes for import on first boot
  interval: 30s
  timeout: 10s
  retries: 3
```

## Manual Import (Alternative)

If you prefer manual import or need to update data:

```bash
# Run import from mage-server container
docker compose exec mage-server /app/scryfall-import \
  --input=/app/data/scryfall-all-cards-latest.json \
  --lang=en \
  --skip-tokens=true \
  --batch=1000

# Or use the production update script
./update-cards-prod.sh --download
```

## Verification

Check card data after startup:

```bash
# Check Scryfall cards
docker compose exec postgres psql -U mage -d mage -c \
  "SELECT COUNT(*) FROM scryfall_cards WHERE lang='en';"

# Check compatibility view
docker compose exec postgres psql -U mage -d mage -c \
  "SELECT COUNT(*) FROM cards;"

# Test card lookup
docker compose exec postgres psql -U mage -d mage -c \
  "SELECT name, type_line, mana_cost FROM cards WHERE name LIKE 'Lightning%' LIMIT 5;"
```

## Troubleshooting

### Import fails on first startup
- Check logs: `docker logs mage-server`
- Verify Scryfall JSON exists: `ls -lh mage-server-go/data/`
- Check file is mounted: `docker compose exec mage-server ls -lh /app/data/`

### Server starts without card data
- The entrypoint script will warn but continue starting
- Import data manually (see Manual Import section)
- Restart container: `docker compose restart mage-server`

### Need to re-import data
```bash
# Clear existing data
docker compose exec postgres psql -U mage -d mage -c \
  "TRUNCATE scryfall_cards CASCADE;"

# Restart mage-server (will detect empty table and re-import)
docker compose restart mage-server
```

## Files Changed

1. **Dockerfile** - Added scryfall-import build + postgresql-client + entrypoint
2. **docker-entrypoint.sh** - New entrypoint script with auto-import logic
3. **docker-compose.prod.yml** - Mounts data directory + extended healthcheck
4. **docker-initdb/00_migrate_and_import.sh** - Updated comments (migrations only)

## Migration from Manual Import

If you're migrating from a setup with manually imported data:
1. The entrypoint detects existing `scryfall_cards` data
2. Skips import automatically
3. Updates compatibility view
4. No data loss or re-import required
