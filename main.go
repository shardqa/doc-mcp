package main

import (
	"encoding/json"
	"io"
	"os"
	"github.com/shardqa/doc-mcp/internal"
)

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
			internal.HandleToolsList(rpcReq.ID, encoder)
		case "listOfferings", "list_tools":
			internal.HandleListOfferings(rpcReq.ID, encoder)
		case "create_markdown_file":
			internal.HandleCreateMarkdownFile(rpcReq.Params, encoder)
		case "edit_markdown_file":
			internal.HandleEditMarkdownFile(rpcReq.Params, encoder)
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