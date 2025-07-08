package server

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/will-wow/larkdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func RefactorFolderLogic(folderPath string) error {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", folderPath, err)
	}

	markdownFiles := []os.DirEntry{}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			markdownFiles = append(markdownFiles, entry)
		}
	}

	if len(markdownFiles) <= 10 {
		return fmt.Errorf("folder %s has %d markdown files, no refactoring needed (threshold is >10)", folderPath, len(markdownFiles))
	}

	groups := make(map[string][]os.DirEntry)
	for _, file := range markdownFiles {
		filename := file.Name()
		baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

		var key string
		if strings.Contains(baseName, "_") {
			key = strings.Split(baseName, "_")[0]
		} else if strings.Contains(baseName, "-") {
			key = strings.Split(baseName, "-")[0]
		} else {
			key = "common"
		}
		groups[key] = append(groups[key], file)
	}

	movedFiles := make(map[string]string)

	for groupName, groupFiles := range groups {
		if len(groupFiles) > 1 {
			newDir := filepath.Join(folderPath, groupName)
			if err := os.MkdirAll(newDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", newDir, err)
			}

			for _, fileToMove := range groupFiles {
				oldPath := filepath.Join(folderPath, fileToMove.Name())
				newPath := filepath.Join(newDir, fileToMove.Name())
				if err := os.Rename(oldPath, newPath); err != nil {
					return fmt.Errorf("failed to move file from %s to %s: %w", oldPath, newPath, err)
				}
				absOldPath, _ := filepath.Abs(oldPath)
				absNewPath, _ := filepath.Abs(newPath)
				movedFiles[absOldPath] = absNewPath
			}
		}
	}

	if err := updateLinksLogic(movedFiles); err != nil {
		return fmt.Errorf("failed to update links: %w", err)
	}

	return nil
}

func updateLinksLogic(movedFiles map[string]string) error {
	mdParser := goldmark.New()
	mdRenderer := goldmark.New(
		goldmark.WithRenderer(larkdown.NewNodeRenderer()),
	)

	for oldPath, newPath := range movedFiles {
		source, err := os.ReadFile(newPath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", newPath, err)
		}

		doc := mdParser.Parser().Parse(text.NewReader(source))

		walkErr := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}

			if link, ok := n.(*ast.Link); ok {
				dest := string(link.Destination)
				if !strings.HasSuffix(dest, ".md") || strings.HasPrefix(dest, "http") {
					return ast.WalkContinue, nil
				}

				linkAbsPath, err := filepath.Abs(filepath.Join(filepath.Dir(oldPath), dest))
				if err != nil {
					return ast.WalkStop, err
				}

				if newDestAbs, ok := movedFiles[linkAbsPath]; ok {
					newRelDest, err := filepath.Rel(filepath.Dir(newPath), newDestAbs)
					if err != nil {
						return ast.WalkStop, err
					}
					link.Destination = []byte(newRelDest)
				}
			}

			return ast.WalkContinue, nil
		})
		if walkErr != nil {
			return fmt.Errorf("error walking ast for %s: %w", newPath, walkErr)
		}

		var buf bytes.Buffer
		if err := mdRenderer.Renderer().Render(&buf, source, doc); err != nil {
			return fmt.Errorf("failed to render markdown for %s: %w", newPath, err)
		}

		if err := os.WriteFile(newPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write updated markdown to %s: %w", newPath, err)
		}
	}

	return nil
} 