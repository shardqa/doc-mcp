package test

import (
	"bufio"
	"encoding/json"
	"io"
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

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
}

type CallToolResult struct {
	Content []map[string]interface{} `json:"content"`
	IsError bool                     `json:"isError"`
}

func sendMCPInitialization(writer io.Writer) error {
	initReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "init",
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test",
				"version": "1.0",
			},
		},
	}
	b, _ := json.Marshal(initReq)
	_, err := writer.Write(append(b, '\n'))
	if err != nil {
		return err
	}

	initNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}
	b, _ = json.Marshal(initNotif)
	_, err = writer.Write(append(b, '\n'))
	return err
}

func TestCreateMarkdownFile(t *testing.T) {
	os.Remove("doc/test1.md")
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = ".."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	require.NoError(t, cmd.Start())
	time.Sleep(1 * time.Second)

	err = sendMCPInitialization(stdin)
	require.NoError(t, err)

	content := "# Title\n\nThis is a [link1](file1.md) and [link2](file2.md)."
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "create_markdown_file",
			"arguments": map[string]interface{}{
				"name":    "test1.md",
				"content": content,
			},
		},
	}
	b, _ := json.Marshal(req)
	_, err = stdin.Write(append(b, '\n'))
	require.NoError(t, err)
	stdin.Close()

	scanner := bufio.NewScanner(stdout)
	var resp JSONRPCResponse
	timeout := time.After(5 * time.Second)
	found := false
	foundInit := false
	
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
				
				if resp.ID == "init" {
					foundInit = true
					continue
				}
				
				if !foundInit {
					continue
				}
				
				if resp.ID == "1" {
					require.Equal(t, "2.0", resp.JSONRPC)
					require.Equal(t, "1", resp.ID)
					
					resultBytes, _ := json.Marshal(resp.Result)
					var result CallToolResult
					err = json.Unmarshal(resultBytes, &result)
					require.NoError(t, err)
					require.False(t, result.IsError)
					require.NotEmpty(t, result.Content)
					found = true
				}
			}
		}
	}
	data, err := os.ReadFile("../doc/test1.md")
	require.NoError(t, err)
	require.Equal(t, content, strings.TrimSpace(string(data)))

	os.Remove("../doc/test1.md")

	cmd.Process.Kill()
}

func TestCreateMarkdownFile_MinimalContent(t *testing.T) {
	os.Remove("doc/test_minimal.md")
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = ".."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	require.NoError(t, cmd.Start())
	time.Sleep(1 * time.Second)

	err = sendMCPInitialization(stdin)
	require.NoError(t, err)

	content := "test"
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "2",
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "create_markdown_file",
			"arguments": map[string]interface{}{
				"name":    "test_minimal.md",
				"content": content,
			},
		},
	}
	b, _ := json.Marshal(req)
	_, err = stdin.Write(append(b, '\n'))
	require.NoError(t, err)
	stdin.Close()

	scanner := bufio.NewScanner(stdout)
	var resp JSONRPCResponse
	timeout := time.After(5 * time.Second)
	found := false
	foundInit := false
	
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
				
				if resp.ID == "init" {
					foundInit = true
					continue
				}
				
				if !foundInit {
					continue
				}
				
				if resp.ID == "2" {
					require.Equal(t, "2.0", resp.JSONRPC)
					require.Equal(t, "2", resp.ID)
					
					resultBytes, _ := json.Marshal(resp.Result)
					var result CallToolResult
					err = json.Unmarshal(resultBytes, &result)
					require.NoError(t, err)
					require.False(t, result.IsError)
					require.NotEmpty(t, result.Content)
					require.True(t, len(result.Content) >= 1)
					found = true
				}
			}
		}
	}
	data, err := os.ReadFile("../doc/test_minimal.md")
	require.NoError(t, err)
	require.Equal(t, content, strings.TrimSpace(string(data)))

	os.Remove("../doc/test_minimal.md")

	cmd.Process.Kill()
} 