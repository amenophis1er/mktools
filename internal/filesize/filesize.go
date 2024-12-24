package filesize

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse parses a human-readable file size string (e.g., "1MB", "500KB") into bytes
func Parse(size string) (int64, error) {
	size = strings.TrimSpace(strings.ToUpper(size))

	if size == "" {
		return 0, fmt.Errorf("empty size string")
	}

	var multiplier int64 = 1
	var numeralStr string

	switch {
	case strings.HasSuffix(size, "KB"):
		multiplier = 1024
		numeralStr = strings.TrimSuffix(size, "KB")
	case strings.HasSuffix(size, "MB"):
		multiplier = 1024 * 1024
		numeralStr = strings.TrimSuffix(size, "MB")
	case strings.HasSuffix(size, "GB"):
		multiplier = 1024 * 1024 * 1024
		numeralStr = strings.TrimSuffix(size, "GB")
	case strings.HasSuffix(size, "B"):
		numeralStr = strings.TrimSuffix(size, "B")
	default:
		return 0, fmt.Errorf("invalid size suffix in %q", size)
	}

	numeral, err := strconv.ParseFloat(numeralStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size number in %q: %w", size, err)
	}

	return int64(numeral * float64(multiplier)), nil
}

// Format formats a byte count into a human-readable string
func Format(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
