package main

import (
	"encoding/json"
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
				break
			}
			continue
		}
		switch rpcReq.Method {
		case "initialize":
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
			encoder.Encode(resp)
		case "notifications/initialized", "initialized":
			continue
		case "tools/list":
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
			encoder.Encode(resp)
		case "listOfferings", "list_tools":
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
			encoder.Encode(resp)
		case "create_markdown_file":
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
			encoder.Encode(resp)
		default:
			errResp := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      rpcReq.ID,
				"error": map[string]interface{}{
					"code":    -32601,
					"message": "Method not found",
				},
			}
			encoder.Encode(errResp)
		}
	}
} 