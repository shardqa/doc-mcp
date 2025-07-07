package server

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

func ValidateMarkdownFile(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ValidateMarkdownParams]) (*mcp.CallToolResultFor[any], error) {
	warnings := validateMarkdown(params.Arguments.Content)

	content := []mcp.Content{&mcp.TextContent{Text: "Markdown content is valid"}}
	if len(warnings) > 0 {
		content = append(content, &mcp.TextContent{Text: "Warnings: " + strings.Join(warnings, "; ")})
	}

	return &mcp.CallToolResultFor[any]{
		Content: content,
		IsError: false,
	}, nil
} 