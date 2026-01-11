#!/usr/bin/env bash
#
# Update the Go server's Postgres `cards` table from upstream magefree/mage.
#
# What this does:
# - Clones https://github.com/magefree/mage (latest master by default)
# - Starts the upstream Java server once to generate Mage.Server/db/cards.h2.mv.db
# - Exports that H2 DB to SQL (cards.sql)
# - Imports into Postgres (Java tables: card/expansion/etc), then migrates into Go table: cards
#
# Requirements:
# - git, mvn, java
# - docker (if using --docker mode; recommended)
#
# Usage:
#   ./mage-server-go/scripts/update_cards_db_from_upstream.sh --docker
#   ./mage-server-go/scripts/update_cards_db_from_upstream.sh --direct
#
# Optional env:
#   UPSTREAM_REF=master
#   DOCKER_CONTAINER=mage-postgres
#   DB_NAME=mage DB_USER=mage DB_PASSWORD=mage DB_HOST=localhost DB_PORT=5432
#
set -euo pipefail

MODE="docker"
if [[ "${1:-}" == "--direct" ]]; then
  MODE="direct"
  shift
elif [[ "${1:-}" == "--docker" ]]; then
  MODE="docker"
  shift
fi

UPSTREAM_REF="${UPSTREAM_REF:-master}"
UPSTREAM_REPO="${UPSTREAM_REPO:-https://github.com/magefree/mage.git}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
REPO_ROOT="$(cd "$GO_ROOT/.." && pwd)"

DB_NAME="${DB_NAME:-mage}"
DB_USER="${DB_USER:-mage}"
DB_PASSWORD="${DB_PASSWORD:-mage}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DOCKER_CONTAINER="${DOCKER_CONTAINER:-mage-postgres}"

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: missing required command: $cmd" >&2
    exit 1
  fi
}

require_cmd git
require_cmd mvn
require_cmd java
if [[ "$MODE" == "docker" ]]; then
  require_cmd docker
fi

echo "=== Update Cards DB From Upstream XMage ==="
echo "Repo:   $UPSTREAM_REPO ($UPSTREAM_REF)"
echo "Mode:   $MODE"
echo "Target: postgres db=$DB_NAME user=$DB_USER"
if [[ "$MODE" == "docker" ]]; then
  echo "Docker: container=$DOCKER_CONTAINER"
else
  echo "Host:   $DB_HOST:$DB_PORT"
fi
echo ""

TMP_DIR="/Users/aron/dev/opensource/mage/tmp"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

UP_DIR="$TMP_DIR/mage-upstream"
echo "Cloning upstream into: $UP_DIR"
git clone --depth=1 --branch "$UPSTREAM_REF" "$UPSTREAM_REPO" "$UP_DIR" >/dev/null
echo "✓ Clone complete"
echo ""

JAVA_SERVER_DIR="$UP_DIR/Mage.Server"
JAVA_DB_BASE="$JAVA_SERVER_DIR/db/cards.h2"
JAVA_DB_FILE="${JAVA_DB_BASE}.mv.db"

echo "Building upstream Java dependencies + Mage.Server..."
# Mage.Server depends on other upstream modules (mage, mage-common, mage-sets, etc.).
# Building only Mage.Server can fail dependency resolution, so build the reactor from repo root.
(cd "$UP_DIR" && mvn -DskipTests -pl Mage.Server -am package)
echo "✓ Build complete"
echo ""

echo "Starting upstream Java server to generate cards DB..."
(cd "$JAVA_SERVER_DIR" && mvn -DskipTests exec:java) &
MVN_PID="$!"

deadline=$((SECONDS + 1200)) # 20 minutes
while [[ ! -f "$JAVA_DB_FILE" ]]; do
  if (( SECONDS > deadline )); then
    echo "ERROR: timed out waiting for H2 DB to be generated at: $JAVA_DB_FILE" >&2
    kill -INT "$MVN_PID" >/dev/null 2>&1 || true
    wait "$MVN_PID" >/dev/null 2>&1 || true
    exit 1
  fi
  sleep 2
done

# Give it a moment to finish writing.
sleep 5
kill -INT "$MVN_PID" >/dev/null 2>&1 || true
wait "$MVN_PID" >/dev/null 2>&1 || true
echo "✓ H2 DB generated: $JAVA_DB_FILE"
echo ""

echo "Exporting H2 -> SQL (cards.sql) using existing tooling..."
(cd "$GO_ROOT" && JAVA_DB_PATH="$JAVA_DB_BASE" OUTPUT_SQL="$GO_ROOT/data/cards.sql" ./scripts/h2_to_sql.sh)
echo ""

echo "Ensuring Postgres schema is up to date..."
if [[ "$MODE" == "docker" ]]; then
  (cd "$GO_ROOT" && DOCKER_CONTAINER="$DOCKER_CONTAINER" DB_NAME="$DB_NAME" DB_USER="$DB_USER" ./scripts/run_postgres_migrations.sh)
else
  (cd "$GO_ROOT" && DB_HOST="$DB_HOST" DB_PORT="$DB_PORT" DB_NAME="$DB_NAME" DB_USER="$DB_USER" DB_PASSWORD="$DB_PASSWORD" ./scripts/run_postgres_migrations.sh --direct)
fi
echo ""

echo "Importing into Postgres and migrating to Go cards table..."
if [[ "$MODE" == "docker" ]]; then
  (cd "$GO_ROOT" && DOCKER_CONTAINER="$DOCKER_CONTAINER" DB_NAME="$DB_NAME" DB_USER="$DB_USER" ./scripts/import_to_postgres.sh --docker)
else
  (cd "$GO_ROOT" && DB_HOST="$DB_HOST" DB_PORT="$DB_PORT" DB_NAME="$DB_NAME" DB_USER="$DB_USER" DB_PASSWORD="$DB_PASSWORD" ./scripts/import_to_postgres.sh --direct)
fi

echo ""
echo "✓ Done. Your database should now include upstream's newest cards."
