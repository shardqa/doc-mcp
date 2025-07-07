package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CreateMarkdownParams struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type EditMarkdownParams struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func CreateMarkdownFile(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateMarkdownParams]) (*mcp.CallToolResultFor[any], error) {
	warnings := validateMarkdown(params.Arguments.Content)
	
	err := os.MkdirAll("doc", 0755)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Failed to create directory: " + err.Error()}},
			IsError: true,
		}, nil
	}
	
	f, err := os.Create("doc/" + params.Arguments.Name)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Failed to create file: " + err.Error()}},
			IsError: true,
		}, nil
	}
	defer f.Close()
	
	_, err = f.WriteString(params.Arguments.Content)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Failed to write file: " + err.Error()}},
			IsError: true,
		}, nil
	}
	
	cmd := exec.Command("markdownlint", "--fix", "doc/"+params.Arguments.Name)
	cmd.Run()
	
	content := []mcp.Content{&mcp.TextContent{Text: "File created successfully"}}
	if len(warnings) > 0 {
		content = append(content, &mcp.TextContent{Text: "Warnings: " + strings.Join(warnings, "; ")})
	}
	
	return &mcp.CallToolResultFor[any]{
		Content: content,
		IsError: false,
	}, nil
}

func EditMarkdownFile(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditMarkdownParams]) (*mcp.CallToolResultFor[any], error) {
	warnings := validateMarkdown(params.Arguments.Content)
	
	err := os.MkdirAll("doc", 0755)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Failed to create directory: " + err.Error()}},
			IsError: true,
		}, nil
	}
	
	f, err := os.Create("doc/" + params.Arguments.Name)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Failed to edit file: " + err.Error()}},
			IsError: true,
		}, nil
	}
	defer f.Close()
	
	_, err = f.WriteString(params.Arguments.Content)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "Failed to write file: " + err.Error()}},
			IsError: true,
		}, nil
	}
	
	cmd := exec.Command("markdownlint", "--fix", "doc/"+params.Arguments.Name)
	cmd.Run()
	
	content := []mcp.Content{&mcp.TextContent{Text: "File edited successfully"}}
	if len(warnings) > 0 {
		content = append(content, &mcp.TextContent{Text: "Warnings: " + strings.Join(warnings, "; ")})
	}
	
	return &mcp.CallToolResultFor[any]{
		Content: content,
		IsError: false,
	}, nil
}

func validateMarkdown(content string) []string {
	warnings := []string{}
	
	linkRe := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
	links := linkRe.FindAllStringSubmatch(content, -1)
	if len(links) < 2 {
		warnings = append(warnings, "File should have at least 2 internal links.")
	}
	
	lines := strings.Split(content, "\n")
	if len(lines) > 100 {
		warnings = append(warnings, "File should not exceed 100 lines.")
	}
	
	return warnings
}

func main() {
	server := mcp.NewServer("doc_mcp", "0.1.0", nil)
	
	server.AddTools(
		mcp.NewServerTool[CreateMarkdownParams, any]("create_markdown_file", "Create a new markdown file in the doc/ folder.", CreateMarkdownFile),
		mcp.NewServerTool[EditMarkdownParams, any]("edit_markdown_file", "Edit an existing markdown file in the doc/ folder.", EditMarkdownFile),
	)
	
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
} 