#!/usr/bin/env python3
"""
Script to add init() registration functions to generated_tokens.go.
This makes all tokens self-register with the token registry on package init.
"""

import re
import sys

def add_init_functions(file_path):
    """Add init() functions after each token constructor."""

    with open(file_path, 'r') as f:
        content = f.read()

    # Pattern to match token function definitions
    # Matches: func NewXXXToken() *Token {
    pattern = r'(func (New\w+Token)\(\) \*Token \{)'

    # Find all matches
    matches = list(re.finditer(pattern, content))

    if not matches:
        print("No token functions found!")
        return False

    print(f"Found {len(matches)} token functions")

    # Process in reverse order to maintain correct positions
    offset = 0
    modifications = []

    for match in matches:
        func_name = match.group(2)  # e.g., "NewATATToken"
        token_name = func_name[3:]  # Remove "New" prefix, e.g., "ATATToken"

        # Find the end of the function (next '}' at the start of a line after the function declaration)
        func_start = match.end()

        # Find the closing brace of the function
        # We need to find the '}' that closes this function
        lines_after = content[func_start:].split('\n')
        brace_count = 1  # We're already inside the function
        line_index = 0

        for i, line in enumerate(lines_after):
            if '{' in line:
                brace_count += line.count('{')
            if '}' in line:
                brace_count -= line.count('}')

            if brace_count == 0:
                line_index = i
                break

        # Find position after the closing brace
        func_end_pos = func_start + sum(len(l) + 1 for l in lines_after[:line_index+1])

        # Create the init function
        init_func = f'''
func init() {{
	Register("{token_name}", {func_name})
}}
'''

        modifications.append((func_end_pos, init_func, token_name))

    # Apply modifications in reverse order
    modifications.sort(key=lambda x: x[0], reverse=True)

    new_content = content
    for pos, init_func, token_name in modifications:
        new_content = new_content[:pos] + init_func + new_content[pos:]

    # Write back
    with open(file_path, 'w') as f:
        f.write(new_content)

    print(f"Successfully added {len(modifications)} init() functions")
    print(f"Sample tokens registered: {[m[2] for m in modifications[:5]]}")
    return True

if __name__ == '__main__':
    file_path = '/Users/aron/dev/opensource/mage/mage-server-go/internal/game/token/generated_tokens.go'

    if len(sys.argv) > 1:
        file_path = sys.argv[1]

    success = add_init_functions(file_path)
    sys.exit(0 if success else 1)
