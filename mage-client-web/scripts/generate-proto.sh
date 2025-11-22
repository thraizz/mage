#!/bin/bash

# Generate TypeScript types from proto files
# This script uses ts-proto to generate TypeScript code from .proto files

set -e

echo "Generating TypeScript types from proto files..."

# Create output directory if it doesn't exist
mkdir -p src/lib/generated

# Use ts-proto via protoc to generate TypeScript files
npx protoc \
  --plugin=./node_modules/.bin/protoc-gen-ts_proto \
  --ts_proto_out=src/lib/generated \
  --ts_proto_opt=outputServices=grpc-js,env=browser,useOptionals=messages,exportCommonSymbols=false,esModuleInterop=true \
  --proto_path=proto \
  proto/*.proto

echo "✓ TypeScript types generated successfully in src/lib/generated/"
echo ""
echo "Generated files:"
ls -lh src/lib/generated/
