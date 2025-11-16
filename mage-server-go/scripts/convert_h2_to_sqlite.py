#!/usr/bin/env python3
"""
Convert H2 SQL export to SQLite-compatible SQL

H2 and SQLite have different SQL dialects. This script converts:
- Data types (VARCHAR_IGNORECASE -> TEXT, etc.)
- Boolean values (TRUE/FALSE -> 1/0)
- Table/column names (case sensitivity)
- CREATE TABLE syntax
- INSERT statements
"""

import re
import sys
from pathlib import Path


def convert_data_type(h2_type):
    """Convert H2 data type to SQLite equivalent"""
    h2_type = h2_type.upper()
    
    # Remove size specifications for simplicity
    h2_type = re.sub(r'\(\d+\)', '', h2_type)
    
    type_map = {
        'VARCHAR': 'TEXT',
        'VARCHAR_IGNORECASE': 'TEXT',
        'CHAR': 'TEXT',
        'CLOB': 'TEXT',
        'LONGVARCHAR': 'TEXT',
        'INTEGER': 'INTEGER',
        'INT': 'INTEGER',
        'BIGINT': 'INTEGER',
        'SMALLINT': 'INTEGER',
        'TINYINT': 'INTEGER',
        'BOOLEAN': 'INTEGER',
        'BIT': 'INTEGER',
        'DECIMAL': 'REAL',
        'DOUBLE': 'REAL',
        'FLOAT': 'REAL',
        'REAL': 'REAL',
        'TIMESTAMP': 'INTEGER',
        'DATE': 'INTEGER',
        'TIME': 'INTEGER',
        'BLOB': 'BLOB',
        'BINARY': 'BLOB',
    }
    
    for h2, sqlite in type_map.items():
        if h2_type.startswith(h2):
            return sqlite
    
    return 'TEXT'  # Default fallback


def convert_boolean(value):
    """Convert H2 boolean to SQLite integer"""
    if value.upper() == 'TRUE':
        return '1'
    elif value.upper() == 'FALSE':
        return '0'
    return value


def convert_create_table(line):
    """Convert CREATE TABLE statement from H2 to SQLite"""
    # Remove H2-specific clauses
    line = re.sub(r'\s+CACHED\b', '', line, flags=re.IGNORECASE)
    line = re.sub(r'\s+NOT\s+PERSISTENT\b', '', line, flags=re.IGNORECASE)
    
    # Convert data types
    def replace_type(match):
        return match.group(1) + convert_data_type(match.group(2))
    
    line = re.sub(r'(\s+)(\w+(?:\(\d+\))?)\s*(?=,|\))', replace_type, line)
    
    return line


def convert_insert(line):
    """Convert INSERT statement from H2 to SQLite"""
    # Convert boolean values
    line = re.sub(r'\bTRUE\b', '1', line, flags=re.IGNORECASE)
    line = re.sub(r'\bFALSE\b', '0', line, flags=re.IGNORECASE)
    
    # Convert NULL handling
    line = re.sub(r'\bNULL\b', 'NULL', line, flags=re.IGNORECASE)
    
    return line


def should_skip_line(line):
    """Check if line should be skipped"""
    skip_patterns = [
        r'^\s*SET\s+',  # SET commands
        r'^\s*ALTER\s+SEQUENCE\s+',  # Sequence alterations
        r'^\s*CREATE\s+SEQUENCE\s+',  # Sequence creation
        r'^\s*CREATE\s+USER\s+',  # User creation
        r'^\s*CREATE\s+SCHEMA\s+',  # Schema creation
        r'^\s*GRANT\s+',  # Grant statements
        r'^\s*--',  # Comments (keep for debugging)
    ]
    
    for pattern in skip_patterns:
        if re.match(pattern, line, re.IGNORECASE):
            return True
    
    return False


def convert_h2_to_sqlite(input_file, output_file):
    """Main conversion function"""
    print(f"Converting {input_file} to SQLite format...")
    
    with open(input_file, 'r', encoding='utf-8') as f_in:
        lines = f_in.readlines()
    
    converted_lines = []
    in_create_table = False
    create_table_buffer = []
    
    stats = {
        'total_lines': len(lines),
        'converted_lines': 0,
        'skipped_lines': 0,
        'create_tables': 0,
        'inserts': 0,
    }
    
    for line in lines:
        # Skip empty lines
        if not line.strip():
            continue
        
        # Skip H2-specific commands
        if should_skip_line(line):
            stats['skipped_lines'] += 1
            continue
        
        # Handle CREATE TABLE (may span multiple lines)
        if re.match(r'^\s*CREATE\s+TABLE\s+', line, re.IGNORECASE):
            in_create_table = True
            create_table_buffer = [line]
            continue
        
        if in_create_table:
            create_table_buffer.append(line)
            if ';' in line:
                # End of CREATE TABLE
                full_statement = ''.join(create_table_buffer)
                converted = convert_create_table(full_statement)
                converted_lines.append(converted)
                stats['create_tables'] += 1
                stats['converted_lines'] += 1
                in_create_table = False
                create_table_buffer = []
            continue
        
        # Handle INSERT statements
        if re.match(r'^\s*INSERT\s+INTO\s+', line, re.IGNORECASE):
            converted = convert_insert(line)
            converted_lines.append(converted)
            stats['inserts'] += 1
            stats['converted_lines'] += 1
            continue
        
        # Handle CREATE INDEX
        if re.match(r'^\s*CREATE\s+INDEX\s+', line, re.IGNORECASE):
            converted_lines.append(line)
            stats['converted_lines'] += 1
            continue
        
        # Keep other statements as-is
        converted_lines.append(line)
        stats['converted_lines'] += 1
    
    # Write converted SQL
    with open(output_file, 'w', encoding='utf-8') as f_out:
        f_out.writelines(converted_lines)
    
    print(f"✓ Conversion complete")
    print(f"  Total lines: {stats['total_lines']}")
    print(f"  Converted: {stats['converted_lines']}")
    print(f"  Skipped: {stats['skipped_lines']}")
    print(f"  CREATE TABLE: {stats['create_tables']}")
    print(f"  INSERT: {stats['inserts']}")
    print(f"  Output: {output_file}")


def main():
    if len(sys.argv) != 3:
        print("Usage: convert_h2_to_sqlite.py <input.sql> <output.sql>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    if not Path(input_file).exists():
        print(f"ERROR: Input file not found: {input_file}")
        sys.exit(1)
    
    try:
        convert_h2_to_sqlite(input_file, output_file)
    except Exception as e:
        print(f"ERROR: Conversion failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()
