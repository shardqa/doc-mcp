package main

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