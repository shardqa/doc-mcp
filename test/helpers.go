package test

import (
	"bufio"
	"encoding/json"
	"io"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
}

type CallToolResult struct {
	Content []map[string]interface{} `json:"content"`
	IsError bool                     `json:"isError"`
}

func SendMCPInitialization(writer io.Writer) error {
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

func RunMCPWithCommand(t *testing.T, toolName string, requestID string, args map[string]interface{}) JSONRPCResponse {
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

	err = SendMCPInitialization(stdin)
	require.NoError(t, err)

	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      requestID,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": args,
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

				if resp.ID == requestID {
					found = true
				}
			}
		}
	}

	cmd.Process.Kill()
	return resp
}

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

func ValidateJSONRPCResponse(t *testing.T, resp JSONRPCResponse, expectedID string) CallToolResult {
	require.Equal(t, "2.0", resp.JSONRPC)
	require.Equal(t, expectedID, resp.ID)

	resultBytes, _ := json.Marshal(resp.Result)
	var result CallToolResult
	err := json.Unmarshal(resultBytes, &result)
	require.NoError(t, err)
	require.False(t, result.IsError)
	require.NotEmpty(t, result.Content)

	return result
} 