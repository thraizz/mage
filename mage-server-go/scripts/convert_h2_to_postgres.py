#!/usr/bin/env python3
"""
Convert H2 SQL export to PostgreSQL-compatible SQL

H2 and PostgreSQL have different SQL dialects. This script converts:
- Data types (VARCHAR_IGNORECASE -> TEXT, etc.)
- Boolean values (keep TRUE/FALSE)
- Table/column names (lowercase for PostgreSQL)
- CREATE TABLE syntax
- INSERT statements
- Sequences
"""

import re
import sys
from pathlib import Path


def convert_data_type(h2_type):
    """Convert H2 data type to PostgreSQL equivalent"""
    h2_type_upper = h2_type.upper()
    
    type_map = {
        'VARCHAR_IGNORECASE': 'TEXT',
        'VARCHAR': 'VARCHAR',
        'CHAR': 'CHAR',
        'CLOB': 'TEXT',
        'LONGVARCHAR': 'TEXT',
        'INTEGER': 'INTEGER',
        'INT': 'INTEGER',
        'BIGINT': 'BIGINT',
        'SMALLINT': 'SMALLINT',
        'TINYINT': 'SMALLINT',
        'BOOLEAN': 'BOOLEAN',
        'BIT': 'BOOLEAN',
        'DECIMAL': 'DECIMAL',
        'DOUBLE': 'DOUBLE PRECISION',
        'FLOAT': 'REAL',
        'REAL': 'REAL',
        'TIMESTAMP': 'TIMESTAMP',
        'DATE': 'DATE',
        'TIME': 'TIME',
        'BLOB': 'BYTEA',
        'BINARY': 'BYTEA',
    }
    
    # Check for type with size (e.g., VARCHAR(255))
    match = re.match(r'(\w+)(\(\d+\))?', h2_type_upper)
    if match:
        base_type = match.group(1)
        size = match.group(2) or ''
        
        if base_type in type_map:
            pg_type = type_map[base_type]
            # Keep size for VARCHAR, CHAR, DECIMAL
            if base_type in ['VARCHAR', 'CHAR', 'DECIMAL'] and size:
                return pg_type + size
            return pg_type
    
    return h2_type  # Return as-is if no mapping found


def convert_table_name(name):
    """Convert table name to lowercase (PostgreSQL convention)"""
    return name.lower()


def convert_column_name(name):
    """Convert column name to lowercase (PostgreSQL convention)"""
    return name.lower()


def convert_create_table(statement):
    """Convert CREATE TABLE statement from H2 to PostgreSQL"""
    # Remove H2-specific clauses
    statement = re.sub(r'\s+CACHED\b', '', statement, flags=re.IGNORECASE)
    statement = re.sub(r'\s+NOT\s+PERSISTENT\b', '', statement, flags=re.IGNORECASE)
    
    # Extract table name and convert to lowercase
    match = re.search(r'CREATE\s+TABLE\s+(\w+)', statement, re.IGNORECASE)
    if match:
        old_name = match.group(1)
        new_name = convert_table_name(old_name)
        statement = statement.replace(old_name, new_name, 1)
    
    # Convert column definitions
    def replace_column_def(match):
        indent = match.group(1)
        col_name = convert_column_name(match.group(2))
        col_type = convert_data_type(match.group(3))
        rest = match.group(4)
        return f"{indent}{col_name} {col_type}{rest}"
    
    # Match column definitions: "  COLUMNNAME TYPE constraints,"
    statement = re.sub(
        r'(\s+)(\w+)\s+(\w+(?:\(\d+(?:,\s*\d+)?\))?)(.*?)(?=,|\))',
        replace_column_def,
        statement,
        flags=re.IGNORECASE
    )
    
    return statement


def convert_insert(statement):
    """Convert INSERT statement from H2 to PostgreSQL"""
    # Extract table name and convert to lowercase
    match = re.search(r'INSERT\s+INTO\s+(\w+)', statement, re.IGNORECASE)
    if match:
        old_name = match.group(1)
        new_name = convert_table_name(old_name)
        statement = statement.replace(f'INSERT INTO {old_name}', f'INSERT INTO {new_name}', 1)
    
    # PostgreSQL uses TRUE/FALSE (keep as-is)
    # Just ensure they're uppercase
    statement = re.sub(r'\btrue\b', 'TRUE', statement, flags=re.IGNORECASE)
    statement = re.sub(r'\bfalse\b', 'FALSE', statement, flags=re.IGNORECASE)
    
    return statement


def convert_create_index(statement):
    """Convert CREATE INDEX statement"""
    # Convert table name to lowercase
    match = re.search(r'ON\s+(\w+)', statement, re.IGNORECASE)
    if match:
        old_name = match.group(1)
        new_name = convert_table_name(old_name)
        statement = re.sub(
            r'(ON\s+)' + old_name,
            r'\1' + new_name,
            statement,
            flags=re.IGNORECASE
        )
    
    return statement


def should_skip_line(line):
    """Check if line should be skipped"""
    skip_patterns = [
        r'^\s*SET\s+',  # SET commands (H2 specific)
        r'^\s*CREATE\s+USER\s+',  # User creation
        r'^\s*CREATE\s+SCHEMA\s+',  # Schema creation (unless PUBLIC)
        r'^\s*GRANT\s+',  # Grant statements
    ]
    
    for pattern in skip_patterns:
        if re.match(pattern, line, re.IGNORECASE):
            return True
    
    return False


def convert_sequence(statement):
    """Convert sequence statements"""
    # H2: CREATE SEQUENCE ... START WITH X
    # PostgreSQL: CREATE SEQUENCE ... START X
    statement = re.sub(r'\bSTART\s+WITH\b', 'START', statement, flags=re.IGNORECASE)
    
    # Convert sequence name to lowercase
    match = re.search(r'CREATE\s+SEQUENCE\s+(\w+)', statement, re.IGNORECASE)
    if match:
        old_name = match.group(1)
        new_name = convert_table_name(old_name)
        statement = statement.replace(old_name, new_name, 1)
    
    return statement


def convert_h2_to_postgres(input_file, output_file):
    """Main conversion function"""
    print(f"Converting {input_file} to PostgreSQL format...")
    
    with open(input_file, 'r', encoding='utf-8') as f_in:
        content = f_in.read()
    
    # Split into statements (handle multi-line statements)
    statements = []
    current_statement = []
    in_statement = False
    
    for line in content.split('\n'):
        if not line.strip():
            continue
        
        if should_skip_line(line):
            continue
        
        current_statement.append(line)
        
        if ';' in line:
            # End of statement
            full_statement = '\n'.join(current_statement)
            statements.append(full_statement)
            current_statement = []
    
    # Convert each statement
    converted = []
    stats = {
        'create_table': 0,
        'insert': 0,
        'create_index': 0,
        'create_sequence': 0,
        'alter_sequence': 0,
        'other': 0,
    }
    
    for stmt in statements:
        stmt_upper = stmt.upper()
        
        if 'CREATE TABLE' in stmt_upper:
            converted.append(convert_create_table(stmt))
            stats['create_table'] += 1
        elif 'INSERT INTO' in stmt_upper:
            converted.append(convert_insert(stmt))
            stats['insert'] += 1
        elif 'CREATE INDEX' in stmt_upper:
            converted.append(convert_create_index(stmt))
            stats['create_index'] += 1
        elif 'CREATE SEQUENCE' in stmt_upper:
            converted.append(convert_sequence(stmt))
            stats['create_sequence'] += 1
        elif 'ALTER SEQUENCE' in stmt_upper:
            # Convert sequence name to lowercase
            match = re.search(r'ALTER\s+SEQUENCE\s+(\w+)', stmt, re.IGNORECASE)
            if match:
                old_name = match.group(1)
                new_name = convert_table_name(old_name)
                stmt = stmt.replace(old_name, new_name, 1)
            converted.append(stmt)
            stats['alter_sequence'] += 1
        else:
            converted.append(stmt)
            stats['other'] += 1
    
    # Write converted SQL
    with open(output_file, 'w', encoding='utf-8') as f_out:
        f_out.write('\n'.join(converted))
    
    print(f"✓ Conversion complete")
    print(f"  CREATE TABLE: {stats['create_table']}")
    print(f"  INSERT: {stats['insert']}")
    print(f"  CREATE INDEX: {stats['create_index']}")
    print(f"  CREATE SEQUENCE: {stats['create_sequence']}")
    print(f"  ALTER SEQUENCE: {stats['alter_sequence']}")
    print(f"  Other: {stats['other']}")
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
