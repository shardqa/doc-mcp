package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListToolsResponse struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Tools []Tool `json:"tools"`
}

func TestListTools(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = ".."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "failed to get stdin")
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err, "failed to get stdout")
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err, "failed to get stderr")
	require.NoError(t, cmd.Start(), "failed to start server")
	defer cmd.Process.Kill()

	time.Sleep(1 * time.Second)

	req := `{"id":"1","type":"list_tools"}`
	_, err = stdin.Write([]byte(req + "\n"))
	require.NoError(t, err, "failed to write to stdin")
	stdin.Close()

	scanner := bufio.NewScanner(stdout)
	var resp ListToolsResponse
	found := false
	timeout := time.After(5 * time.Second)
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
				fmt.Printf("STDOUT: %s\n", line)
				err := json.Unmarshal([]byte(line), &resp)
				require.NoError(t, err, "invalid JSON response")
				require.Equal(t, "list_tools_response", resp.Type)
				require.NotEmpty(t, resp.Tools)
				// Check for create_markdown_file tool
				foundTool := false
				for _, tool := range resp.Tools {
					if tool.Name == "create_markdown_file" && tool.Description != "" {
						foundTool = true
						break
					}
				}
				require.True(t, foundTool, "create_markdown_file tool not found or missing description")
				found = true
			}
		}
	}
}
