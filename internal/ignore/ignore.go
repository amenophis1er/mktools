// internal/ignore/ignore.go

package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type IgnoreList struct {
	patterns []string
}

func New() *IgnoreList {
	return &IgnoreList{
		patterns: make([]string, 0),
	}
}

func (il *IgnoreList) AddPattern(pattern string) {
	if pattern = strings.TrimSpace(pattern); pattern != "" && !strings.HasPrefix(pattern, "#") {
		// Normalize pattern
		pattern = strings.ReplaceAll(pattern, "\\", "/")
		il.patterns = append(il.patterns, pattern)
	}
}

func (il *IgnoreList) AddPatterns(patterns []string) {
	for _, pattern := range patterns {
		il.AddPattern(pattern)
	}
}

func (il *IgnoreList) LoadGitignore(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		il.AddPattern(scanner.Text())
	}

	return scanner.Err()
}

func (il *IgnoreList) ShouldIgnore(path string) bool {
	// Normalize path
	path = strings.ReplaceAll(path, "\\", "/")

	for _, pattern := range il.patterns {
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}

		// Handle negation patterns
		if strings.HasPrefix(pattern, "!") {
			pattern = pattern[1:]
			if matchPath(pattern, path) {
				return false
			}
			continue
		}

		if matchPath(pattern, path) {
			return true
		}
	}
	return false
}

func matchPath(pattern, path string) bool {
	// Handle recursive globbing "**"
	if strings.Contains(pattern, "**") {
		return matchRecursiveGlob(pattern, path)
	}

	// Handle directory-only patterns ending with "/"
	if strings.HasSuffix(pattern, "/") {
		pattern = strings.TrimSuffix(pattern, "/")
		return strings.HasPrefix(path, pattern+"/") || path == pattern
	}

	// Handle direct directory matches (e.g., ".idea/")
	if strings.Contains(pattern, "/") {
		if ok, _ := filepath.Match(pattern, path); ok {
			return true
		}
		return strings.HasPrefix(path, pattern)
	}

	// Handle simple glob patterns
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if ok, _ := filepath.Match(pattern, segment); ok {
			return true
		}
	}

	return false
}

func matchRecursiveGlob(pattern, path string) bool {
	// Split pattern into parts
	parts := strings.Split(pattern, "**")

	if len(parts) != 2 {
		return false
	}

	prefix := strings.TrimRight(parts[0], "/")
	suffix := strings.TrimLeft(parts[1], "/")

	// Check prefix match
	if prefix != "" && !strings.HasPrefix(path, prefix) {
		return false
	}

	// Check suffix match
	if suffix != "" && !strings.HasSuffix(path, suffix) {
		return false
	}

	// If both prefix and suffix match, it's a match
	return true
}
