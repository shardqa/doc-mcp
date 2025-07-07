#!/bin/sh
cd /home/richard_rosario/git/doc-mcp

# Use pre-built binary if it exists, otherwise compile and run
if [ -f "./doc-mcp" ]; then
    exec ./doc-mcp
else
    exec go run .
fi 