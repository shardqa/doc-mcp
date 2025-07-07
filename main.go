package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shardqa/doc-mcp/internal/server"
)

func main() {
	srv := mcp.NewServer("doc_mcp", "0.1.0", nil)
	
	srv.AddTools(
		mcp.NewServerTool[server.CreateMarkdownParams, any]("create_markdown_file", "Create a new markdown file in the doc/ folder.", server.CreateMarkdownFile),
		mcp.NewServerTool[server.EditMarkdownParams, any]("edit_markdown_file", "Edit an existing markdown file in the doc/ folder.", server.EditMarkdownFile),
		mcp.NewServerTool[server.ValidateMarkdownParams, any]("validate_markdown_file", "Validate markdown content and return warnings without modifying files.", server.ValidateMarkdownFile),
	)
	
	if err := srv.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
} 