#!/bin/bash

# Generate TypeScript types from proto files
# This script uses ts-proto to generate TypeScript code from .proto files
# It uses the SAME proto files as the Go server for consistency

set -e

# Path to server proto files (shared source of truth)
SERVER_PROTO_DIR="../mage-server-go/api/proto"
OUT_DIR="src/lib/generated"

echo "Generating TypeScript types from proto files..."
echo "Using proto files from: ${SERVER_PROTO_DIR}"

# Create output directory if it doesn't exist
mkdir -p ${OUT_DIR}

# Remove old generated files to avoid stale code
rm -rf ${OUT_DIR}/*

# Generate TypeScript code from all server proto files
# ts-proto options:
#   - outputServices=grpc-js: Generate grpc-js compatible service definitions
#   - env=browser: Generate browser-compatible code
#   - useOptionals=messages: Use optional fields for message fields
#   - exportCommonSymbols=false: Don't export common symbols
#   - esModuleInterop=true: Enable ES module interop
#   - outputIndex=true: Generate index.ts for easier imports

npx protoc \
  --plugin=./node_modules/.bin/protoc-gen-ts_proto \
  --ts_proto_out=${OUT_DIR} \
  --ts_proto_opt=outputServices=generic-definitions,useExactTypes=false,env=browser,useOptionals=messages,exportCommonSymbols=false,esModuleInterop=true,outputIndex=true \
  --proto_path=${SERVER_PROTO_DIR} \
  ${SERVER_PROTO_DIR}/mage/v1/*.proto

echo ""
echo "✓ TypeScript types generated successfully in ${OUT_DIR}/"
echo ""
echo "Generated files:"
ls -lh ${OUT_DIR}/mage/v1/
echo ""
echo "Total files: $(find ${OUT_DIR} -name '*.ts' | wc -l)"
