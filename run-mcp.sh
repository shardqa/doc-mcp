#!/bin/sh

echo "Running MCP from project directory: $PWD" > /tmp/mcp.log

# Salva o diretório atual (projeto)
PROJECT_DIR="$PWD"

# Caminho absoluto do binário
MCP_BIN="/home/richard_rosario/git/doc-mcp/doc-mcp"

# Muda para o diretório do projeto
cd "$PROJECT_DIR"

echo "Changed to project directory: $PWD" >> /tmp/mcp.log

# Executa o binário, mas a partir do diretório do projeto
exec "$MCP_BIN" 