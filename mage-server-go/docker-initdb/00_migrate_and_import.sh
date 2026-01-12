#!/usr/bin/env bash
#
# Postgres init script (runs ONLY when the data directory is empty).
#
# This is designed to make a fresh:
#   docker compose down -v && docker compose up -d
# bring up a working DB with schema + seeded cards.
#
# It:
# - Applies all migrations in /docker-entrypoint-initdb.d/migrations/*_*.up.sql
# - Tracks applied versions in schema_migrations (idempotent)
# - Imports cards from /docker-entrypoint-initdb.d/data/cards_export.csv (if present)
#
set -euo pipefail

echo "==> initdb: applying migrations + importing cards"

psql_base() {
  # Forward args (e.g., -c) to psql
  PAGER=cat psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "${POSTGRES_DB}" "$@"
}

echo "==> initdb: ensure schema_migrations exists"
psql_base -c "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, name TEXT NOT NULL, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW());" >/dev/null

echo "==> initdb: running migrations"
shopt -s nullglob
migrations=(/docker-entrypoint-initdb.d/migrations/*_*.up.sql)
if [[ ${#migrations[@]} -eq 0 ]]; then
  echo "WARNING: no migrations found at /docker-entrypoint-initdb.d/migrations/*_*.up.sql"
else
  for f in "${migrations[@]}"; do
    base="$(basename "$f")"
    version="${base%%_*}"
    name="${base#*_}"
    name="${name%.up.sql}"

    already="$(psql_base -tAc "SELECT 1 FROM schema_migrations WHERE version = '${version}' LIMIT 1;" | xargs || true)"
    if [[ "$already" == "1" ]]; then
      echo "SKIP  $version  $name"
      continue
    fi

    echo "APPLY $version  $name"
    psql_base <<SQL
BEGIN;
\i '${f}'
INSERT INTO schema_migrations (version, name) VALUES ('${version}', '${name}');
COMMIT;
SQL
  done
fi

# NOTE: H2 data seeding is now deprecated in favor of Scryfall import
# The old H2 seed file is no longer used. Instead, use the scryfall-import tool.
#
# To seed cards in a fresh database:
# 1. Download Scryfall bulk data:
#    ./scripts/download_scryfall_bulk.sh
#
# 2. Import Scryfall data:
#    go run ./cmd/scryfall-import/main.go \
#      --input=./data/scryfall-all-cards-latest.json \
#      --lang=en --skip-tokens=true --batch=1000
#
# 3. Create the compatibility view (if needed):
#    psql -U mage -d mage -f migrations/009_create_scryfall_tables.up.sql

echo "==> initdb: Scryfall card import"
echo "INFO: H2 seed data is no longer used"
echo "INFO: Cards will be loaded from Scryfall data via scryfall-import tool"
echo "INFO: Run the import manually after the database is initialized"
echo ""
echo "Quick start:"
echo "  1. Download: ./scripts/download_scryfall_bulk.sh"
echo "  2. Import: go run ./cmd/scryfall-import/main.go --input=./data/scryfall-all-cards-latest.json"

echo "==> initdb: done"
