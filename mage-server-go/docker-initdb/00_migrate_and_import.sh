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

SEED_SQL="/docker-entrypoint-initdb.d/data/cards_seed.postgres.sql"
if [[ -f "$SEED_SQL" ]]; then
  echo "==> initdb: seeding cards via H2->SQL seed file: $SEED_SQL"
  # This file is expected to be PostgreSQL-compatible SQL that:
  # - creates/populates the Java-export tables (e.g. card, expansion)
  # - migrates them into the Go server's cards table
  psql_base -f "$SEED_SQL"
  echo "==> initdb: card seeding complete"
else
  echo "WARNING: seed SQL not found at $SEED_SQL (cards will NOT be seeded)"
  echo "         Generate it locally with: mage-server-go/scripts/build_postgres_seed_sql.sh"
fi

echo "==> initdb: done"
