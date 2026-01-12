# MAGE Scripts Guide

This directory contains various scripts for working with the MAGE Go server.

## Scryfall Card Data Scripts

MAGE now uses Scryfall bulk data instead of the legacy XMage Java system.

### Download Latest Scryfall Data

**`download_scryfall_bulk.sh`** - Download the latest Scryfall "All Cards" bulk data

```bash
./scripts/download_scryfall_bulk.sh
```

Downloads to: `mage-server-go/data/scryfall-all-cards-YYYYMMDD.json`

### Migrate to Scryfall

**`migrate_to_scryfall.sh`** - Full migration from XMage to Scryfall

```bash
./scripts/migrate_to_scryfall.sh
```

This script:
1. Applies database migrations
2. Downloads latest Scryfall data (if needed)
3. Imports cards to PostgreSQL
4. Creates compatibility views
5. Verifies the migration

See `SCRYFALL_MIGRATION_GUIDE.md` for details.

### Rollback from Scryfall

**`rollback_from_scryfall.sh`** - Revert to XMage data (if needed)

```bash
./scripts/rollback_from_scryfall.sh
```

Restores the old `cards_xmage_backup` table.

## Database Migration Scripts

**`run_postgres_migrations.sh`** - Apply database migrations

```bash
./scripts/run_postgres_migrations.sh
```

Applies all pending migrations from the `migrations/` directory.

## Protobuf Scripts

**`generate_proto.sh`** - Generate Go code from .proto files

```bash
make proto  # Preferred - use the Makefile target
# OR
./scripts/generate_proto.sh
```

Generates gRPC/protobuf Go code from `.proto` files in `api/`.

## Production Scripts

See the root directory for production deployment scripts:

- `deploy.sh` - Deploy full application to production
- `update-cards-prod.sh` - Update only card data on production
- `update-cards-weekly.sh` - Automated weekly card updates

## Scryfall Importer

The Scryfall importer is a Go application located at `cmd/scryfall-import/`.

### Usage

```bash
# Import from local file
go run ./cmd/scryfall-import/main.go \
  --input=/path/to/scryfall-all-cards.json \
  --lang=en \
  --skip-tokens=true \
  --batch=1000

# Or set DATABASE_URL directly
DATABASE_URL='postgres://user:pass@localhost:5432/mage?sslmode=disable' \
  go run ./cmd/scryfall-import/main.go \
  --input=/path/to/scryfall-all-cards.json
```

### Options

- `--input` - Path to Scryfall JSON file
- `--lang` - Language filter (default: "en")
- `--skip-tokens` - Skip token cards (recommended)
- `--batch` - Batch size for inserts (default: 1000)

## Workflow

### Initial Setup

```bash
# 1. Run database migrations
./scripts/run_postgres_migrations.sh

# 2. Download and import Scryfall data
./scripts/migrate_to_scryfall.sh
```

### Update Card Data

```bash
# Local development
./scripts/download_scryfall_bulk.sh
./scripts/migrate_to_scryfall.sh

# Production
cd ../../  # Back to root
./update-cards-prod.sh --download
```

### Development Workflow

```bash
# 1. Make changes to database schema
vim migrations/NNN_description.up.sql
vim migrations/NNN_description.down.sql

# 2. Apply migrations
./scripts/run_postgres_migrations.sh

# 3. Test
go test ./...

# 4. Deploy
cd ../../
./deploy.sh
```

## Documentation

- **SCRYFALL_MIGRATION_GUIDE.md** - User guide for Scryfall data
- **DATA_MIGRATION.md** - Technical migration details
- **SCRYFALL_DFC_CARDS.md** - Double-faced card handling
- **PRODUCTION_CARD_UPDATES.md** - Production update guide
- **QUICK_REFERENCE.md** - Command reference

## Legacy XMage System (Removed)

The legacy XMage/Java-based card import system has been removed. If you need to access the old system:

```bash
# Check out the commit before Scryfall migration
git log --all --oneline --grep="Scryfall"
git checkout <commit-before-migration>
```

All Maven, Java, and XMage dependencies have been removed from the project.
