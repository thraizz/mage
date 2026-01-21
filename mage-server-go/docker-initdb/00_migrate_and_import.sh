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

# NOTE: Scryfall card import is now handled automatically by the mage-server container
# on first startup via the docker-entrypoint.sh script.
#
# The entrypoint script will:
# 1. Check if scryfall_cards table is empty
# 2. Import Scryfall data from /app/data/scryfall-all-cards-*.json
# 3. Create the compatibility view
# 4. Start the server
#
# This happens automatically when the mage-server container starts for the first time.

echo "==> initdb: Scryfall card import"
echo "INFO: Card data will be imported automatically by mage-server on first startup"
echo "INFO: The mage-server entrypoint script will detect empty database and import data"
echo "INFO: Ensure Scryfall JSON file exists in mage-server-go/data/ directory"

echo "==> initdb: done"
