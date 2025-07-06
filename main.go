package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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

func main() {
	fmt.Fprintln(os.Stderr, "MCP server started")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Scanner error: %v\n", err)
			}
			fmt.Fprintln(os.Stderr, "No more input, waiting...")
			// Wait for a bit before checking again
			// This prevents immediate exit if stdin closes
			// (simulate a long-running server for test)
			select {}
		}
		line := scanner.Bytes()
		fmt.Fprintf(os.Stderr, "Received: %s\n", string(line))
		var req ListToolsRequest
		err := json.Unmarshal(line, &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unmarshal error: %v\n", err)
			continue
		}
		if req.Type == "list_tools" {
			fmt.Fprintln(os.Stderr, "Processing list_tools request")
			resp := ListToolsResponse{
				ID:   req.ID,
				Type: "list_tools_response",
				Tools: []Tool{{
					Name:        "create_markdown_file",
					Description: "Create a new markdown file in the doc/ folder.",
				}},
			}
			b, _ := json.Marshal(resp)
			os.Stdout.Write(b)
			os.Stdout.Write([]byte("\n"))
			fmt.Fprintln(os.Stderr, "Sent list_tools_response")
		}
	}
} 