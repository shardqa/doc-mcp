package main

import "encoding/json"

// ... existing code ...
// Functions for registering tools and offerings, extracted from main.go
// ... existing code ...

func handleToolsList(rpcReqID interface{}, encoder *json.Encoder) {
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
		"id":      rpcReqID,
		"result": map[string]interface{}{
			"tools": tools,
		},
	}
	encoder.Encode(resp)
}

func handleListOfferings(rpcReqID interface{}, encoder *json.Encoder) {
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
		"id":      rpcReqID,
		"result": map[string]interface{}{
			"offerings": offerings,
		},
	}
	encoder.Encode(resp)
} 