#!/bin/bash
# Export Java H2 database to portable SQL format
# Can be imported into SQLite, PostgreSQL, MySQL, etc.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
JAVA_DB_PATH="${PROJECT_ROOT}/../Mage.Server/db/cards.h2"
OUTPUT_SQL="${PROJECT_ROOT}/data/cards.sql"

echo "=== H2 to SQL Export ==="
echo "Source: ${JAVA_DB_PATH}.mv.db"
echo "Output: $OUTPUT_SQL"
echo ""

# Check if H2 database exists
if [ ! -f "${JAVA_DB_PATH}.mv.db" ]; then
    echo "ERROR: Java H2 database not found"
    echo "Expected: ${JAVA_DB_PATH}.mv.db"
    echo ""
    echo "To generate the database:"
    echo "  cd ../Mage.Server"
    echo "  mvn clean install"
    echo "  mvn exec:java"
    echo "  (wait for server to start, then stop it)"
    exit 1
fi

# Create output directory
mkdir -p "$(dirname "$OUTPUT_SQL")"

# Download H2 JAR if not present (must match Java MAGE version: 1.4.197)
H2_JAR="${SCRIPT_DIR}/h2.jar"
H2_VERSION="1.4.197"
if [ ! -f "$H2_JAR" ] || ! java -jar "$H2_JAR" -help 2>&1 | grep -q "$H2_VERSION"; then
    echo "Downloading H2 JAR (version $H2_VERSION to match Java MAGE)..."
    rm -f "$H2_JAR"
    curl -L -o "$H2_JAR" "https://repo1.maven.org/maven2/com/h2database/h2/${H2_VERSION}/h2-${H2_VERSION}.jar"
    echo "✓ Downloaded H2 JAR"
fi

# Export to SQL using H2's Script tool
# Note: Java MAGE uses empty user/password, database name is "cards.h2" (H2 adds .mv.db)
echo "Exporting H2 database to SQL..."
java -cp "$H2_JAR" org.h2.tools.Script \
    -url "jdbc:h2:file:${JAVA_DB_PATH}" \
    -user "" \
    -password "" \
    -script "$OUTPUT_SQL"

echo "✓ Exported to SQL"

# Get file size
FILE_SIZE=$(du -h "$OUTPUT_SQL" | cut -f1)
LINE_COUNT=$(wc -l < "$OUTPUT_SQL")

echo ""
echo "=== Export Complete ==="
echo "✓ Output: $OUTPUT_SQL"
echo "✓ Size: $FILE_SIZE"
echo "✓ Lines: $LINE_COUNT"
echo ""
echo "Next steps:"
echo "  For SQLite:   ./scripts/import_to_sqlite.sh"
echo "  For PostgreSQL: ./scripts/import_to_postgres.sh"
echo "  For MySQL:    ./scripts/import_to_mysql.sh"
