#!/usr/bin/env bash
#
# Run all PostgreSQL migrations (*/migrations/*_*.up.sql) in order.
#
# Why this exists:
# - The official postgres image only executes /docker-entrypoint-initdb.d/*
#   on *first* initialization of the data volume.
# - After the volume exists, new migrations won't auto-run, so we need an
#   explicit migration runner.
#
# This script is:
# - Ordered: runs migrations sorted by filename (e.g., 001_, 002_, ...)
# - Safe to re-run: records applied versions in schema_migrations
# - Docker-friendly: defaults to running against the "mage-postgres" container
#
# Usage (recommended for prod host, with docker compose):
#   cd ~/gomage
#   ./mage-server-go/scripts/run_postgres_migrations.sh
#
# Usage (direct Postgres connection, no docker exec):
#   PAGER=cat DB_HOST=localhost DB_PORT=5432 DB_USER=mage DB_PASSWORD=mage DB_NAME=mage \
#     ./mage-server-go/scripts/run_postgres_migrations.sh --direct
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-$PROJECT_ROOT/migrations}"

MODE="docker"
if [[ "${1:-}" == "--direct" ]]; then
  MODE="direct"
fi

DB_NAME="${DB_NAME:-mage}"
DB_USER="${DB_USER:-mage}"
DB_PASSWORD="${DB_PASSWORD:-mage}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DOCKER_CONTAINER="${DOCKER_CONTAINER:-mage-postgres}"

if [[ ! -d "$MIGRATIONS_DIR" ]]; then
  echo "ERROR: migrations dir not found: $MIGRATIONS_DIR" >&2
  exit 1
fi

shopt -s nullglob
MIGRATION_FILES=("$MIGRATIONS_DIR"/*_*.up.sql)
if [[ ${#MIGRATION_FILES[@]} -eq 0 ]]; then
  echo "ERROR: no migrations found at: $MIGRATIONS_DIR/*_*.up.sql" >&2
  exit 1
fi

run_psql_stdin() {
  if [[ "$MODE" == "docker" ]]; then
    # Run psql inside the postgres container
    docker exec -i "$DOCKER_CONTAINER" env PAGER=cat \
      psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME"
  else
    # Run psql directly (expects client installed + network access)
    PAGER=cat PGPASSWORD="$DB_PASSWORD" \
      psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"
  fi
}

run_psql_query() {
  local query="$1"
  if [[ "$MODE" == "docker" ]]; then
    docker exec -i "$DOCKER_CONTAINER" env PAGER=cat \
      psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
  else
    PAGER=cat PGPASSWORD="$DB_PASSWORD" \
      psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -tAc "$query"
  fi
}

echo "=== PostgreSQL Migration Runner ==="
echo "Mode: $MODE"
echo "Database: $DB_NAME"
if [[ "$MODE" == "docker" ]]; then
  echo "Container: $DOCKER_CONTAINER"
  echo "Migrations (container path): /docker-entrypoint-initdb.d"
else
  echo "Host: $DB_HOST:$DB_PORT"
  echo "Migrations (host path): $MIGRATIONS_DIR"
fi
echo ""

echo "Ensuring schema_migrations table exists..."
run_psql_query "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, name TEXT NOT NULL, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW());" >/dev/null

applied_count="$(run_psql_query "SELECT COUNT(*) FROM schema_migrations;" | xargs || true)"
echo "Already applied migrations: ${applied_count:-0}"
echo ""

applied_now=0
skipped=0

for file in "${MIGRATION_FILES[@]}"; do
  base="$(basename "$file")"
  version="${base%%_*}"
  name="${base#*_}"
  name="${name%.up.sql}"

  already="$(run_psql_query "SELECT 1 FROM schema_migrations WHERE version = '${version}' LIMIT 1;" | xargs || true)"
  if [[ "$already" == "1" ]]; then
    echo "SKIP  $version  $name"
    skipped=$((skipped + 1))
    continue
  fi

  if [[ "$MODE" == "docker" ]]; then
    migration_path="/docker-entrypoint-initdb.d/$base"
  else
    migration_path="$file"
  fi

  echo "APPLY $version  $name"

  # Run migration and record it as applied in one transaction.
  # Note: we intentionally only record on success.
  run_psql_stdin <<SQL
BEGIN;
\i '${migration_path}'
INSERT INTO schema_migrations (version, name) VALUES ('${version}', '${name}');
COMMIT;
SQL

  applied_now=$((applied_now + 1))
done

echo ""
echo "=== Done ==="
echo "Applied now: $applied_now"
echo "Skipped:     $skipped"
echo ""
echo "Verify:"
if [[ "$MODE" == "docker" ]]; then
  echo "  docker exec -i $DOCKER_CONTAINER env PAGER=cat psql -U $DB_USER -d $DB_NAME -c \"\\\\dt\""
else
  echo "  PAGER=cat psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \"\\\\dt\""
fi
