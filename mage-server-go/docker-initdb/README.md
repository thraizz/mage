# Database Initialization

This directory contains the initialization script that runs when PostgreSQL starts with an empty data volume.

## How it works

When you run `docker compose up -d` for the first time (or after `docker compose down -v`), PostgreSQL executes all scripts in `/docker-entrypoint-initdb.d/` in alphabetical order.

The `00_migrate_and_import.sh` script:
1. Creates the `schema_migrations` tracking table
2. Applies all migrations from `/docker-entrypoint-initdb.d/migrations/*.up.sql`
3. Tracks applied migrations to prevent re-running them
4. Shows instructions for importing card data

## Card Data: Scryfall (Current)

**We no longer use H2 database exports for card data.** Instead, we use Scryfall's official bulk data.

### Why Scryfall?

- **Official data source**: Direct from Scryfall, the industry-standard MTG API
- **Always up-to-date**: Easy to download latest card data
- **Better structure**: Rich metadata, proper multiface card handling
- **No Java dependency**: No need to export from XMage's H2 database

### Importing Scryfall Data

After the database is initialized, import card data:

```bash
# 1. Download latest Scryfall bulk data
./scripts/download_scryfall_bulk.sh

# 2. Apply Scryfall migrations
docker compose exec -T postgres psql -U mage -d mage < migrations/009_create_scryfall_tables.up.sql

# 3. Import Scryfall data
docker compose exec -T mage-server bash -c \
  "DATABASE_URL='postgres://mage:mage@postgres:5432/mage?sslmode=disable' \
   go run ./cmd/scryfall-import/main.go \
   --input='/app/data/scryfall-all-cards-latest.json' \
   --lang=en --skip-tokens=true --batch=1000"

# 4. Create compatibility view
docker compose exec -T postgres psql -U mage -d mage << 'SQL'
-- Backup old cards table if it exists
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'cards' AND table_type = 'BASE TABLE') THEN
        DROP TABLE IF EXISTS cards_xmage_backup CASCADE;
        ALTER TABLE cards RENAME TO cards_xmage_backup;
    END IF;
END
$$;

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
```

Or use the automated production update script:

```bash
# From the repository root
./update-cards-prod.sh
```

## Legacy: H2 Seed Data (Deprecated)

The old method using `cards_seed.postgres.sql` from H2 exports is **no longer supported**. 

If you have an existing database with H2 data, it will continue to work, but new installations should use Scryfall data exclusively.

### Migration from H2 to Scryfall

If you have an existing database with H2 data and want to migrate:

1. **Backup your current data**:
   ```bash
   docker compose exec -T postgres pg_dump -U mage -d mage -t cards > cards_backup.sql
   ```

2. **Follow the Scryfall import steps above**

3. **Verify the migration**:
   ```bash
   docker compose exec -T postgres psql -U mage -d mage -c \
     "SELECT COUNT(*) FROM scryfall_cards WHERE lang='en';"
   ```

## Troubleshooting

### Database won't initialize

If the init script fails, check logs:
```bash
docker compose logs postgres
```

### Need to re-run initialization

To completely reset and re-initialize:
```bash
docker compose down -v  # WARNING: Deletes all data!
docker compose up -d
```

### Manual migration execution

If you need to run migrations manually:
```bash
docker compose exec -T postgres psql -U mage -d mage < migrations/001_create_users_table.up.sql
# ... repeat for each migration ...
```

## See Also

- [`../migrations/`](../migrations/) - SQL migration files
- [`../cmd/scryfall-import/`](../cmd/scryfall-import/) - Scryfall import tool
- [`../scripts/download_scryfall_bulk.sh`](../scripts/download_scryfall_bulk.sh) - Download script
- [`../../update-cards-prod.sh`](../../update-cards-prod.sh) - Production update automation
