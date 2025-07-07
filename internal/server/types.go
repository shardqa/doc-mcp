package server

type CreateMarkdownParams struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type EditMarkdownParams struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ValidateMarkdownParams struct {
	Content string `json:"content"`
} 