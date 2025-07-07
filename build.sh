#!/bin/sh
# Build script for doc-mcp
echo "Building doc-mcp..."
go build -o doc-mcp .
if [ $? -eq 0 ]; then
    echo "✅ Build successful! Binary: ./doc-mcp"
else
    echo "❌ Build failed!"
    exit 1
fi 