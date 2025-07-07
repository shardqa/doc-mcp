package test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMarkdownFile(t *testing.T) {
	os.Remove("doc/test1.md")
	defer os.Remove("../doc/test1.md")

	content := "# Title\n\nThis is a [link1](file1.md) and [link2](file2.md)."
	args := map[string]interface{}{
		"name":    "test1.md",
		"content": content,
	}

	resp := RunMCPWithCommand(t, "create_markdown_file", "1", args)
	ValidateJSONRPCResponse(t, resp, "1")

	data, err := os.ReadFile("../doc/test1.md")
	require.NoError(t, err)
	require.Equal(t, content, strings.TrimSpace(string(data)))
}

func TestCreateMarkdownFile_MinimalContent(t *testing.T) {
	os.Remove("doc/test_minimal.md")
	defer os.Remove("../doc/test_minimal.md")

	content := "test"
	args := map[string]interface{}{
		"name":    "test_minimal.md",
		"content": content,
	}

	resp := RunMCPWithCommand(t, "create_markdown_file", "2", args)
	result := ValidateJSONRPCResponse(t, resp, "2")
	require.True(t, len(result.Content) >= 1)

	data, err := os.ReadFile("../doc/test_minimal.md")
	require.NoError(t, err)
	require.Equal(t, content, strings.TrimSpace(string(data)))
} 