package test

import (
	"strings"
	"testing"
)

func TestValidateMarkdownFile_NoWarnings(t *testing.T) {
	initReq := `{"jsonrpc":"2.0","id":"init","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	initNotif := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	toolCall := `{"jsonrpc":"2.0","id":"1","method":"tools/call","params":{"name":"validate_markdown_file","arguments":{"content":"# Title\n\nThis is a [link1](http://example.com) and another [link2](http://example.org).\n"}}}`
	
	input := initReq + "\n" + initNotif + "\n" + toolCall
	resp, _, _ := runMCP(input)
	
	if resp["jsonrpc"] != "2.0" {
		t.Fatalf("Expected jsonrpc 2.0, got %v", resp["jsonrpc"])
	}
	if resp["id"] != "1" {
		t.Fatalf("Expected id 1, got %v", resp["id"])
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result object, got %v", resp["result"])
	}
	if result["isError"] != nil && result["isError"].(bool) {
		t.Fatalf("Expected isError false, got true")
	}
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatalf("Expected non-empty content array, got %v", result["content"])
	}
	
	firstContent := content[0].(map[string]interface{})
	text := firstContent["text"].(string)
	if !strings.Contains(text, "valid") {
		t.Fatalf("Expected validation success message, got %s", text)
	}
}

func TestValidateMarkdownFile_WithWarnings(t *testing.T) {
	initReq := `{"jsonrpc":"2.0","id":"init","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	initNotif := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	toolCall := `{"jsonrpc":"2.0","id":"2","method":"tools/call","params":{"name":"validate_markdown_file","arguments":{"content":"# Title\n\nThis has only [one link](http://example.com).\n"}}}`
	
	input := initReq + "\n" + initNotif + "\n" + toolCall
	resp, _, _ := runMCP(input)
	
	result := resp["result"].(map[string]interface{})
	content := result["content"].([]interface{})
	
	if len(content) < 2 {
		t.Fatalf("Expected at least 2 content items (validation message + warnings), got %d", len(content))
	}
	
	warningsContent := content[1].(map[string]interface{})
	warningsText := warningsContent["text"].(string)
	if !strings.Contains(warningsText, "2 internal links") {
		t.Fatalf("Expected warning about links, got %s", warningsText)
	}
}

func TestValidateMarkdownFile_TooManyLines(t *testing.T) {
	lines := make([]string, 101)
	for i := range lines {
		lines[i] = "line"
	}
	lines[0] = "# Title"
	lines[1] = "This has [link1](http://example.com) and [link2](http://example.org)."
	longContent := strings.Join(lines, "\n")
	
	initReq := `{"jsonrpc":"2.0","id":"init","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	initNotif := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	toolCall := `{"jsonrpc":"2.0","id":"3","method":"tools/call","params":{"name":"validate_markdown_file","arguments":{"content":"` + strings.ReplaceAll(longContent, "\n", "\\n") + `"}}}`
	
	input := initReq + "\n" + initNotif + "\n" + toolCall
	resp, _, _ := runMCP(input)
	
	result := resp["result"].(map[string]interface{})
	content := result["content"].([]interface{})
	
	if len(content) < 2 {
		t.Fatalf("Expected at least 2 content items (validation message + warnings), got %d", len(content))
	}
	
	warningsContent := content[1].(map[string]interface{})
	warningsText := warningsContent["text"].(string)
	if !strings.Contains(warningsText, "100 lines") {
		t.Fatalf("Expected warning about lines, got %s", warningsText)
	}
} 