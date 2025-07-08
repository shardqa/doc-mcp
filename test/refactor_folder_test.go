package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/shardqa/doc-mcp/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestRefactorFolder(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "refactor-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create 11 markdown files
	for i := 0; i < 5; i++ {
		content := ""
		if i == 0 {
			// Link to another file in the same group and a file in another group
			content = "[link1](group1_file1.md)\n[link2](group2_file0.md)"
		}
		fileName := fmt.Sprintf("group1_file%d.md", i)
		err := os.WriteFile(filepath.Join(tempDir, fileName), []byte(content), 0644)
		assert.NoError(t, err)
	}
	for i := 0; i < 6; i++ {
		fileName := fmt.Sprintf("group2_file%d.md", i)
		err := os.WriteFile(filepath.Join(tempDir, fileName), []byte(""), 0644)
		assert.NoError(t, err)
	}

	// Call the refactor function
	err = server.RefactorFolderLogic(tempDir)
	assert.NoError(t, err)

	// Assertions
	// Check if group1 directory is created
	group1Dir := filepath.Join(tempDir, "group1")
	assert.DirExists(t, group1Dir)

	// Check if group2 directory is created
	group2Dir := filepath.Join(tempDir, "group2")
	assert.DirExists(t, group2Dir)

	// Check if files are moved
	for i := 0; i < 5; i++ {
		fileName := fmt.Sprintf("group1_file%d.md", i)
		assert.FileExists(t, filepath.Join(group1Dir, fileName))
	}
	for i := 0; i < 6; i++ {
		fileName := fmt.Sprintf("group2_file%d.md", i)
		assert.FileExists(t, filepath.Join(group2Dir, fileName))
	}

	// Check if link is updated
	updatedContent, err := os.ReadFile(filepath.Join(group1Dir, "group1_file0.md"))
	assert.NoError(t, err)
	expectedContent := "[link1](group1_file1.md)\n[link2](../group2/group2_file0.md)\n"
	assert.Equal(t, expectedContent, string(updatedContent))
} 