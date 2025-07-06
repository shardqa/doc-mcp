package test

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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

func TestCreateMarkdownFile(t *testing.T) {
	os.Remove("doc/test1.md")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = ".."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	require.NoError(t, cmd.Start())
	time.Sleep(1 * time.Second)

	content := "# Title\n\nThis is a [link1](file1.md) and [link2](file2.md)."
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id": "1",
		"method": "create_markdown_file",
		"params": map[string]interface{}{
			"id": "1",
			"type": "create_markdown_file",
			"name": "test1.md",
			"content": content,
		},
	}
	b, _ := json.Marshal(req)
	_, err = stdin.Write(append(b, '\n'))
	require.NoError(t, err)
	stdin.Close()

	scanner := bufio.NewScanner(stdout)
	var resp CreateMarkdownFileResponse
	timeout := time.After(5 * time.Second)
	found := false
	for !found {
		select {
		case <-timeout:
			serr := bufio.NewScanner(stderr)
			for serr.Scan() {
				t.Logf("STDERR: %s", serr.Text())
			}
			require.FailNow(t, "timeout waiting for response")
		default:
			if scanner.Scan() {
				line := scanner.Text()
				err := json.Unmarshal([]byte(line), &resp)
				require.NoError(t, err)
				require.Equal(t, "create_markdown_file_response", resp.Type)
				require.True(t, resp.Success)
				found = true
			}
		}
	}
	data, err := os.ReadFile("../doc/test1.md")
	require.NoError(t, err)
	require.Equal(t, content, strings.TrimSpace(string(data)))

	os.Remove("../doc/test1.md")

	cmd.Process.Kill()
} 