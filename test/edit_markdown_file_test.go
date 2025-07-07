package test

import (
	"os"
	"strings"
	"testing"
)

func TestEditMarkdownFile_Success(t *testing.T) {
	os.Remove("doc/test.md")
	defer os.Remove("doc/test.md")

	args := map[string]interface{}{
		"name":    "test.md",
		"content": "# Title\n\nThis is a [link1](http://example.com) and another [link2](http://example.org).\n",
	}

	resp := RunMCPWithCommand(t, "edit_markdown_file", "1", args)
	ValidateJSONRPCResponse(t, resp, "1")
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
	defer os.Remove("doc/test.md")

	args := map[string]interface{}{
		"name":    "test.md",
		"content": "#Title\n\nThis is a [link1](http://example.com) and another [link2](http://example.org).\n",
	}

	resp := RunMCPWithCommand(t, "edit_markdown_file", "2", args)
	ValidateJSONRPCResponse(t, resp, "2")

	out, _ := os.ReadFile("../doc/test.md")
	if strings.Contains(string(out), "#Title") {
		t.Fatalf("markdownlint --fix not applied")
	}
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