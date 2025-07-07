package server

import (
	"regexp"
	"strings"
)

func validateMarkdown(content string) []string {
	warnings := []string{}

	linkRe := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
	links := linkRe.FindAllStringSubmatch(content, -1)
	if len(links) < 2 {
		warnings = append(warnings, "File should have at least 2 internal links.")
	}

	lines := strings.Split(content, "\n")
	if len(lines) > 100 {
		warnings = append(warnings, "File should not exceed 100 lines.")
	}

	return warnings
} 