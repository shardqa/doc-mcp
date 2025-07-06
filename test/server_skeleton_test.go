package test

import (
	"bufio"
	"encoding/json"
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

type Offering struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListToolsResult struct {
	Offerings []Offering `json:"offerings"`
}

type ListToolsRPCResponse struct {
	ID     string          `json:"id"`
	Result ListToolsResult `json:"result"`
}

func TestListTools(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
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

	req := `{"jsonrpc":"2.0","id":"1","method":"list_tools","params":{}}`
	_, err = stdin.Write([]byte(req + "\n"))
	require.NoError(t, err, "failed to write to stdin")
	stdin.Close()

	scanner := bufio.NewScanner(stdout)
	serr := bufio.NewScanner(stderr)
	var resp ListToolsRPCResponse
	found := false
	timeout := time.After(5 * time.Second)
	for !found {
		select {
		case <-timeout:
			for serr.Scan() {
				t.Logf("STDERR: %s", serr.Text())
			}
			require.FailNow(t, "timeout waiting for response")
		default:
			if scanner.Scan() {
				line := scanner.Text()
				t.Logf("STDOUT: %s", line)
				err := json.Unmarshal([]byte(line), &resp)
				if err != nil {
					t.Logf("JSON Unmarshal error: %v", err)
					continue
				}
				t.Logf("Unmarshaled response: %+v", resp)
				// Check for create_markdown_file tool in offerings
				foundTool := false
				for _, offering := range resp.Result.Offerings {
					if offering.Name == "create_markdown_file" && offering.Description != "" {
						foundTool = true
						break
					}
				}
				require.True(t, foundTool, "create_markdown_file tool not found or missing description")
				found = true
			}
		}
	}
	for serr.Scan() {
		t.Logf("STDERR: %s", serr.Text())
	}
}
