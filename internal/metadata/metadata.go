package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Metadata struct {
	GeneratedBy    string            `json:"generated_by"`
	GeneratedAt    time.Time         `json:"generated_at"`
	Version        string            `json:"version"`
	ChecksumSource string            `json:"checksum_source"`
	FileChecksums  map[string]string `json:"file_checksums"`
}

const MetadataMarker = "<!-- MKTOOLS-CONTEXT"
const MetadataEndMarker = "MKTOOLS-CONTEXT -->"

func New() *Metadata {
	return &Metadata{
		GeneratedBy:   "mktools",
		GeneratedAt:   time.Now(),
		Version:       "0.1.0", // TODO: get from version package
		FileChecksums: make(map[string]string),
	}
}

// CalculateSourceChecksum generates a checksum for all source files
func (m *Metadata) CalculateSourceChecksum(files map[string]string) error {
	// Sort files by path for consistent checksums
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Create a hash of all file contents
	h := sha256.New()
	for _, path := range paths {
		content := files[path]

		// Calculate individual file checksum
		fileHash := sha256.Sum256([]byte(content))
		m.FileChecksums[path] = hex.EncodeToString(fileHash[:])

		// Add to global checksum
		fmt.Fprintf(h, "%s:%s\n", path, content)
	}

	m.ChecksumSource = hex.EncodeToString(h.Sum(nil))
	return nil
}

// HasSourceChanged checks if source files have changed compared to stored metadata
func (m *Metadata) HasSourceChanged(root string) (bool, error) {
	for path, storedChecksum := range m.FileChecksums {
		fullPath := filepath.Join(root, path)

		// Check if file still exists
		content, err := os.ReadFile(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return true, nil // File was deleted
			}
			return false, fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Calculate current checksum
		currentHash := sha256.Sum256(content)
		currentChecksum := hex.EncodeToString(currentHash[:])

		if currentChecksum != storedChecksum {
			return true, nil // File was modified
		}
	}

	// Check for new files
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// If we find a file that's not in our checksums, source has changed
		if _, exists := m.FileChecksums[relPath]; !exists {
			return io.EOF // Use EOF as sentinel to stop walking
		}

		return nil
	})

	if err == io.EOF {
		return true, nil // New file found
	}
	if err != nil {
		return false, fmt.Errorf("error walking directory: %w", err)
	}

	return false, nil
}

// String returns a formatted string representation of the metadata
func (m *Metadata) String() string {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		// Fallback to simple format if JSON marshaling fails
		return fmt.Sprintf("%s\ngenerated_by: %s\ngenerated_at: %s\nversion: %s\nchecksum: %s\n%s",
			MetadataMarker,
			m.GeneratedBy,
			m.GeneratedAt.Format(time.RFC3339),
			m.Version,
			m.ChecksumSource,
			MetadataEndMarker)
	}

	return fmt.Sprintf("%s\n%s\n%s", MetadataMarker, string(data), MetadataEndMarker)
}

// ParseFromContent extracts metadata from content containing metadata markers
func ParseFromContent(content string) (*Metadata, error) {
	start := strings.Index(content, MetadataMarker)
	end := strings.Index(content, MetadataEndMarker)

	if start == -1 || end == -1 {
		return nil, fmt.Errorf("metadata markers not found in content")
	}

	// Extract JSON between markers
	jsonData := content[start+len(MetadataMarker) : end]
	jsonData = strings.TrimSpace(jsonData)

	var metadata Metadata
	if err := json.Unmarshal([]byte(jsonData), &metadata); err != nil {
		return nil, fmt.Errorf("error parsing metadata JSON: %w", err)
	}

	return &metadata, nil
}
