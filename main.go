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
		mcp.NewServerTool(
			"create_markdown_file",
			"Create a new markdown file. Parameters: name (string, required) is the file name, content (string, required) is the markdown content, path (string, optional) is a relative folder path inside the project where the file will be created. If path is omitted, the file is created in the current directory.",
			server.CreateMarkdownFile,
		),
		mcp.NewServerTool(
			"edit_markdown_file",
			"Edit an existing markdown file. Parameters: name (string, required) is the file name, content (string, required) is the new markdown content. The file must exist in the doc/ folder.",
			server.EditMarkdownFile,
		),
		mcp.NewServerTool(
			"validate_markdown_file",
			"Validate markdown content and return warnings. Parameters: content (string, required) is the markdown to validate. Does not modify any files.",
			server.ValidateMarkdownFile,
		),
		mcp.NewServerTool(
			"refactor_folder",
			"Refactor a folder by creating subdirectories and moving files. Parameters: folder_path (string, optional) is the relative path to the folder to refactor. Defaults to doc/.",
			server.RefactorFolder,
		),
	)

	if err := srv.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
