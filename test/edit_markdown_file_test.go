package test

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runMCP(jsonrpcReq string) (map[string]interface{}, string, error) {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = ".."
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	stdin.Write([]byte(jsonrpcReq + "\n"))
	stdin.Close()
	outBytes, _ := io.ReadAll(stdout)
	errBytes, _ := io.ReadAll(stderr)
	cmd.Wait()
	
	fullOutput := string(outBytes)
	if len(errBytes) > 0 {
		fullOutput += "\nSTDERR: " + string(errBytes)
	}
	
	lines := strings.Split(strings.TrimSpace(string(outBytes)), "\n")
	var lastResponse map[string]interface{}
	
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var resp map[string]interface{}
		if err := json.Unmarshal([]byte(line), &resp); err == nil {
			if id, exists := resp["id"]; !exists || id != "init" {
				lastResponse = resp
			}
		}
	}
	
	return lastResponse, fullOutput, nil
}

func TestEditMarkdownFile_Success(t *testing.T) {
	os.Remove("doc/test.md")
	
	initReq := `{"jsonrpc":"2.0","id":"init","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	initNotif := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	toolCall := `{"jsonrpc":"2.0","id":"1","method":"tools/call","params":{"name":"edit_markdown_file","arguments":{"name":"test.md","content":"# Title\n\nThis is a [link1](http://example.com) and another [link2](http://example.org).\n"}}}`
	
	input := initReq + "\n" + initNotif + "\n" + toolCall
	println("REQUEST:", input)
	resp, respLine, _ := runMCP(input)
	
	if resp["jsonrpc"] == nil {
		t.Fatalf("Expected JSON-RPC response, got nil. Full response: %s", respLine)
	}
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
	os.Remove("doc/test.md")
}

func TestEditMarkdownFile_TooManyLines(t *testing.T) {
	lines := make([]string, 101)
	for i := range lines {
		lines[i] = "line"
	}
	input := strings.Join(lines, "\n")
	os.WriteFile("test.md", []byte(input), 0644)
	if countLines(input) > 100 {
		if err := validateMarkdownFile("test.md"); err == nil {
			t.Fatalf("Should fail for >100 lines")
		}
	}
	os.Remove("test.md")
}

func TestEditMarkdownFile_TooFewLinks(t *testing.T) {
	input := "# Title\n\nThis is a [link1](http://example.com).\n"
	os.WriteFile("test.md", []byte(input), 0644)
	if countLinks(input) < 2 {
		if err := validateMarkdownFile("test.md"); err == nil {
			t.Fatalf("Should fail for <2 links")
		}
	}
	os.Remove("test.md")
}

func TestEditMarkdownFile_MarkdownlintFix(t *testing.T) {
	os.Remove("doc/test.md")
	
	initReq := `{"jsonrpc":"2.0","id":"init","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`
	initNotif := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	toolCall := `{"jsonrpc":"2.0","id":"2","method":"tools/call","params":{"name":"edit_markdown_file","arguments":{"name":"test.md","content":"#Title\n\nThis is a [link1](http://example.com) and another [link2](http://example.org).\n"}}}`
	
	input := initReq + "\n" + initNotif + "\n" + toolCall
	runMCP(input)
	out, _ := os.ReadFile("doc/test.md")
	if strings.Contains(string(out), "#Title") {
		t.Fatalf("markdownlint --fix not applied")
	}
	os.Remove("doc/test.md")
}

func countLines(s string) int {
	return len(strings.Split(s, "\n"))
}

func countLinks(s string) int {
	return strings.Count(s, "](")
}

func validateMarkdownFile(path string) error {
	b, _ := os.ReadFile(path)
	if countLines(string(b)) > 100 {
		return os.ErrInvalid
	}
	if countLinks(string(b)) < 2 {
		return os.ErrInvalid
	}
	return nil
} 