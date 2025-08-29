#!/bin/bash

# Script to generate Go code from Protocol Buffer definitions

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROTO_DIR="${PROJECT_ROOT}/api/proto"
OUTPUT_DIR="${PROJECT_ROOT}/pkg/generated"

echo -e "${GREEN}Starting Protocol Buffer code generation...${NC}"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${YELLOW}protoc is not installed. Please install it first:${NC}"
    echo "  Linux: apt-get install -y protobuf-compiler"
    echo "  macOS: brew install protobuf"
    echo "  Or download from: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${YELLOW}Installing protoc-gen-go...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${YELLOW}Installing protoc-gen-go-grpc...${NC}"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directory
mkdir -p "${OUTPUT_DIR}"

# Generate Go code for common proto files
echo -e "${GREEN}Generating common proto files...${NC}"
protoc \
    --proto_path="${PROTO_DIR}" \
    --go_out="${OUTPUT_DIR}" \
    --go_opt=paths=source_relative \
    "${PROTO_DIR}"/common/*.proto

# Generate Go code for v1 user service
echo -e "${GREEN}Generating user service proto files...${NC}"
protoc \
    --proto_path="${PROTO_DIR}" \
    --go_out="${OUTPUT_DIR}" \
    --go_opt=paths=source_relative \
    --go-grpc_out="${OUTPUT_DIR}" \
    --go-grpc_opt=paths=source_relative \
    "${PROTO_DIR}"/v1/user/*.proto

# Generate Go code for v1 health service
echo -e "${GREEN}Generating health service proto files...${NC}"
protoc \
    --proto_path="${PROTO_DIR}" \
    --go_out="${OUTPUT_DIR}" \
    --go_opt=paths=source_relative \
    --go-grpc_out="${OUTPUT_DIR}" \
    --go-grpc_opt=paths=source_relative \
    "${PROTO_DIR}"/v1/health/*.proto

echo -e "${GREEN}Protocol Buffer code generation completed!${NC}"
echo -e "${GREEN}Generated files are in: ${OUTPUT_DIR}${NC}"

# List generated files
echo -e "${YELLOW}Generated files:${NC}"
find "${OUTPUT_DIR}" -name "*.pb.go" -o -name "*_grpc.pb.go" | sort