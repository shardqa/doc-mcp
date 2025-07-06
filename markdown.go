package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// ... existing code ...
// Functions for create_markdown_file and edit_markdown_file logic, extracted from main.go
// ... existing code ...

func handleCreateMarkdownFile(params json.RawMessage, encoder *json.Encoder) {
	var req CreateMarkdownFileRequest
	json.Unmarshal(params, &req)
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
}

func handleEditMarkdownFile(params json.RawMessage, encoder *json.Encoder) {
	var req struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
	}
	json.Unmarshal(params, &req)
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
		cmd := exec.Command("markdownlint", "--fix", "doc/"+req.Name)
		cmd.Run()
		success = true
	}
	resp := map[string]interface{}{
		"id":      req.ID,
		"type":    "edit_markdown_file_response",
		"success": success,
		"warnings": warnings,
	}
	encoder.Encode(resp)
} 