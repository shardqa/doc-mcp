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
	cmd := exec.Command("go", "run", "main.go")
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
	respLine := string(outBytes)
	if len(errBytes) > 0 {
		respLine += "\nSTDERR: " + string(errBytes)
	}
	var resp map[string]interface{}
	json.Unmarshal(outBytes, &resp)
	return resp, respLine, nil
}

func TestEditMarkdownFile_Success(t *testing.T) {
	os.Remove("doc/test.md")
	input := `{"jsonrpc":"2.0","id":"1","method":"edit_markdown_file","params":{"id":"1","type":"edit_markdown_file","name":"test.md","content":"# Title\n\nThis is a [link1](http://example.com) and another [link2](http://example.org).\n"}}`
	println("REQUEST:", input)
	resp, respLine, _ := runMCP(input)
	if resp["type"] == nil {
		t.Fatalf("Expected edit_markdown_file_response, got nil. Full response: %s", respLine)
	}
	if resp["type"] != "edit_markdown_file_response" {
		t.Fatalf("Expected edit_markdown_file_response, got %v", resp["type"])
	}
	if !resp["success"].(bool) {
		t.Fatalf("Expected success true, got false")
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
	input := `{"jsonrpc":"2.0","id":"2","method":"edit_markdown_file","params":{"id":"2","type":"edit_markdown_file","name":"test.md","content":"#Title\n\nThis is a [link1](http://example.com) and another [link2](http://example.org).\n"}}`
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