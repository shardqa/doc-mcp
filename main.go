package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type ListToolsRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListToolsResponse struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Tools []Tool `json:"tools"`
}

type CreateMarkdownFileRequest struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type CreateMarkdownFileResponse struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Success bool     `json:"success"`
	Warnings []string `json:"warnings"`
}

type ToolOffering struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

func main() {
	fmt.Fprintln(os.Stderr, "MCP server started")
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for {
		var rpcReq struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      interface{}     `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		if err := decoder.Decode(&rpcReq); err != nil {
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "No more input, exiting...")
				break
			}
			fmt.Fprintln(os.Stderr, "Decode error:", err)
			continue
		}
		fmt.Fprintf(os.Stderr, "Received request: method=%s, id=%v, params=%s\n", rpcReq.Method, rpcReq.ID, string(rpcReq.Params))
		switch rpcReq.Method {
		case "initialize":
			fmt.Fprintln(os.Stderr, "Handling initialize request")
			resp := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      rpcReq.ID,
				"result": map[string]interface{}{
					"protocolVersion": "2025-06-18",
					"serverInfo": map[string]interface{}{
						"name":    "doc_mcp",
						"version": "0.1.0",
					},
					"capabilities": map[string]interface{}{
						"tools":     map[string]interface{}{},
						"prompts":   map[string]interface{}{},
						"resources": map[string]interface{}{},
						"logging":   map[string]interface{}{},
						"roots": map[string]interface{}{
							"listChanged": false,
						},
					},
				},
			}
			fmt.Fprintln(os.Stderr, "Sending initialize response")
			encoder.Encode(resp)
		case "notifications/initialized", "initialized":
			fmt.Fprintln(os.Stderr, "Ignoring notifications/initialized or initialized (no response needed)")
			continue
		case "tools/list":
			fmt.Fprintln(os.Stderr, "Handling tools/list request")
			inputSchema := map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The filename to create",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "The markdown content",
					},
				},
				"required": []string{"name", "content"},
			}
			tools := []map[string]interface{}{
				{
					"name":        "create_markdown_file",
					"description": "Create a new markdown file in the doc/ folder.",
					"inputSchema": inputSchema,
				},
			}
			resp := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      rpcReq.ID,
				"result": map[string]interface{}{
					"tools": tools,
				},
			}
			fmt.Fprintln(os.Stderr, "Sending tools/list response")
			encoder.Encode(resp)
		case "listOfferings", "list_tools":
			fmt.Fprintln(os.Stderr, "Handling listOfferings/list_tools request")
			offerings := []ToolOffering{
				{
					Name:        "create_markdown_file",
					Description: "Create a new markdown file in the doc/ folder.",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type":        "string",
								"description": "The filename to create",
							},
							"content": map[string]interface{}{
								"type":        "string",
								"description": "The markdown content",
							},
						},
						"required": []string{"name", "content"},
					},
				},
			}
			resp := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      rpcReq.ID,
				"result": map[string]interface{}{
					"offerings": offerings,
				},
			}
			fmt.Fprintln(os.Stderr, "Sending listOfferings/list_tools response")
			encoder.Encode(resp)
		case "create_markdown_file":
			fmt.Fprintln(os.Stderr, "Handling create_markdown_file request")
			var req CreateMarkdownFileRequest
			json.Unmarshal(rpcReq.Params, &req)
			warnings := []string{}
			linkRe := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
			links := linkRe.FindAllStringSubmatch(req.Content, -1)
			if len(links) < 2 {
				warnings = append(warnings, "File should have at least 2 internal links.")
			}
			lines := strings.Split(req.Content, "\n")
			if len(lines) > 100 {
				warnings = append(warnings, "File should not exceed 100 lines.")
			}
			os.MkdirAll("doc", 0755)
			f, err := os.Create("doc/" + req.Name)
			success := false
			if err == nil {
				f.WriteString(req.Content)
				f.Close()
				success = true
			}
			resp := CreateMarkdownFileResponse{
				ID:      req.ID,
				Type:    "create_markdown_file_response",
				Success: success,
				Warnings: warnings,
			}
			fmt.Fprintln(os.Stderr, "Sending create_markdown_file response")
			encoder.Encode(resp)
		default:
			fmt.Fprintf(os.Stderr, "Unknown method: %s\n", rpcReq.Method)
			errResp := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      rpcReq.ID,
				"error": map[string]interface{}{
					"code":    -32601,
					"message": "Method not found",
				},
			}
			fmt.Fprintln(os.Stderr, "Sending error response for unknown method")
			encoder.Encode(errResp)
		}
	}
} 