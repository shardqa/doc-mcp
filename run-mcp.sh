#!/bin/sh
SCRIPT_DIR=$(cd -- "$(dirname -- "$0")" && pwd)

# Use pre-built binary if it exists, otherwise compile and run
if [ -f "$SCRIPT_DIR/doc-mcp" ]; then
    exec "$SCRIPT_DIR/doc-mcp"
else
    exec go run "$SCRIPT_DIR"
fi 