#!/usr/bin/env python3
"""
Convert H2 SQL export to PostgreSQL-compatible SQL

H2 and PostgreSQL have different SQL dialects. This script converts:
- Data types (VARCHAR_IGNORECASE -> TEXT, etc.)
- Boolean values (keep TRUE/FALSE)
- Table/column names (lowercase for PostgreSQL)
- CREATE TABLE syntax (remove CACHED, PUBLIC. prefix)
- INSERT statements (remove CAST to VARCHAR_IGNORECASE)
- Sequences
"""

import re
import sys
from pathlib import Path


def _decode_h2_stringdecode_arg(h2_sql_string_literal_contents: str) -> str:
    """
    Decode the contents of an H2 SQL string literal used inside STRINGDECODE('...').

    H2 uses standard SQL single-quote escaping ('' -> ').
    STRINGDECODE additionally interprets backslash escapes like:
      - \\n, \\r, \\t
      - \\uXXXX

    We approximate H2's behavior closely enough for card text imports.
    """
    # SQL single-quote escaping
    s = h2_sql_string_literal_contents.replace("''", "'")

    # Interpret common backslash escapes (notably \\n and \\uXXXX).
    # Python's unicode_escape handles these well for our inputs.
    try:
        s = bytes(s, "utf-8").decode("unicode_escape")
    except Exception:
        # Fall back to raw string if decoding fails
        pass

    return s


def _to_pg_sql_literal(value: str) -> str:
    """Return a PostgreSQL single-quoted string literal for value."""
    return "'" + value.replace("'", "''") + "'"


def convert_h2_to_postgres(input_file, output_file):
    """Main conversion function"""
    print(f"Converting {input_file} to PostgreSQL format...")
    
    with open(input_file, 'r', encoding='utf-8') as f_in:
        content = f_in.read()
    
    lines = content.split('\n')
    converted_lines = []
    stats = {
        'create_table': 0,
        'insert': 0,
        'create_index': 0,
        'alter_table': 0,
        'skipped': 0,
    }
    
    for line in lines:
        original_line = line
        stripped = line.strip()
        
        # Skip empty lines
        if not stripped:
            continue
        
        # Skip H2-specific SET commands
        if stripped.upper().startswith('SET '):
            stats['skipped'] += 1
            continue
        
        # Skip comment lines (but keep them in output)
        if stripped.startswith('--'):
            converted_lines.append(line)
            continue
        
        # Skip standalone semicolons or user creation
        if stripped == ';' or 'CREATE USER' in stripped.upper():
            stats['skipped'] += 1
            continue
        
        # Convert CREATE CACHED TABLE to CREATE TABLE
        if 'CREATE CACHED TABLE' in line.upper():
            line = re.sub(r'CREATE\s+CACHED\s+TABLE', 'CREATE TABLE', line, flags=re.IGNORECASE)
            stats['create_table'] += 1
        
        # Remove PUBLIC. schema prefix
        line = re.sub(r'\bPUBLIC\.', '', line)
        
        # Convert VARCHAR_IGNORECASE to TEXT
        line = re.sub(r'VARCHAR_IGNORECASE(\(\d+\))?', 'TEXT', line, flags=re.IGNORECASE)

        # Convert H2 STRINGDECODE('...') to a plain PostgreSQL string literal
        # Example in H2 export:
        #   CAST(STRINGDECODE('line1\\nline2') AS VARCHAR_IGNORECASE)
        # PostgreSQL has no stringdecode() builtin.
        def _stringdecode_repl(match: re.Match) -> str:
            inner = match.group(1)
            decoded = _decode_h2_stringdecode_arg(inner)
            return _to_pg_sql_literal(decoded)

        line = re.sub(
            r"STRINGDECODE\('((?:[^']|'')*)'\)",
            _stringdecode_repl,
            line,
            flags=re.IGNORECASE,
        )
        
        # Convert CAST(... AS VARCHAR_IGNORECASE) to just the value
        # Handle: CAST('value' AS VARCHAR_IGNORECASE)
        line = re.sub(
            r"CAST\s*\(\s*'([^']*)'\s+AS\s+TEXT\s*\)",
            r"'\1'",
            line,
            flags=re.IGNORECASE
        )
        line = re.sub(
            r"CAST\s*\(\s*'([^']*)'\s+AS\s+VARCHAR_IGNORECASE\s*\)",
            r"'\1'",
            line,
            flags=re.IGNORECASE
        )
        
        # Convert NULL values in CAST
        line = re.sub(
            r"CAST\s*\(\s*NULL\s+AS\s+\w+\s*\)",
            r"NULL",
            line,
            flags=re.IGNORECASE
        )
        
        # Track INSERT statements
        if 'INSERT INTO' in line.upper():
            stats['insert'] += 1
        
        # Track CREATE INDEX statements
        if 'CREATE INDEX' in line.upper():
            stats['create_index'] += 1
        
        # Track ALTER TABLE statements
        if 'ALTER TABLE' in line.upper():
            stats['alter_table'] += 1
        
        converted_lines.append(line)
    
    # Write converted SQL
    with open(output_file, 'w', encoding='utf-8') as f_out:
        f_out.write('\n'.join(converted_lines))
    
    print("✓ Conversion complete")
    print(f"  CREATE TABLE: {stats['create_table']}")
    print(f"  INSERT: {stats['insert']}")
    print(f"  CREATE INDEX: {stats['create_index']}")
    print(f"  ALTER TABLE: {stats['alter_table']}")
    print(f"  Skipped: {stats['skipped']}")
    print(f"  Output: {output_file}")


def main():
    if len(sys.argv) != 3:
        print("Usage: convert_h2_to_postgres.py <input.sql> <output.sql>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    if not Path(input_file).exists():
        print(f"ERROR: Input file not found: {input_file}")
        sys.exit(1)
    
    try:
        convert_h2_to_postgres(input_file, output_file)
    except Exception as e:
        print(f"ERROR: Conversion failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()
