package context

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/amenophis1er/mktools/internal/config"
	"github.com/amenophis1er/mktools/internal/filesize"
	"github.com/amenophis1er/mktools/internal/metadata"
)

type ContextPlugin struct {
	config   *config.Config
	metadata *metadata.Metadata
}

type ContextOptions struct {
	OutputFile    string
	StructureOnly bool
	ContentOnly   bool
	Format        string
}

func New(cfg *config.Config) *ContextPlugin {
	return &ContextPlugin{
		config:   cfg,
		metadata: metadata.New(),
	}
}

func (p *ContextPlugin) Name() string {
	return "context"
}

func (p *ContextPlugin) Description() string {
	return "Generate context data for LLM from project files"
}

func (p *ContextPlugin) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("output", "o", "", "output file (default is ./context.md)")
	cmd.Flags().BoolP("structure-only", "s", false, "only include file structure")
	cmd.Flags().BoolP("content-only", "c", false, "only include file contents")
	cmd.Flags().StringP("format", "f", "", "output format (md or txt)")
	cmd.Flags().Int("max-files", 0, "maximum number of files to process (0 = use config value)")
	cmd.Flags().StringSlice("ignore", nil, "additional patterns to ignore")
}

func (p *ContextPlugin) Execute(ctx context.Context, cmd *cobra.Command, args []string) error {
	// Parse command options
	opts, err := p.parseFlags(cmd)
	if err != nil {
		return err
	}

	// Apply options to config
	if opts.Format != "" {
		p.config.Context.OutputFormat = opts.Format
	}
	if opts.StructureOnly {
		p.config.Context.IncludeFileContent = false
	}
	if opts.ContentOnly {
		p.config.Context.IncludeFileStructure = false
	}

	// Default to current directory if no path provided
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Check for existing context file
	existingContext := p.findExistingContext(path)
	if existingContext != "" {
		existing, err := os.ReadFile(existingContext)
		if err == nil {
			existingMeta, err := metadata.ParseFromContent(string(existing))
			if err == nil {
				changed, err := existingMeta.HasSourceChanged(path)
				if err == nil && !changed {
					fmt.Println("No changes detected in source files. Using existing context.")
					fmt.Println(string(existing))
					return nil
				}
			}
		}
	}

	// Detect project info
	projectInfo, err := detectProject(path)
	if err != nil {
		return fmt.Errorf("failed to detect project info: %w", err)
	}

	// Collect files with options
	files, err := p.collectFiles(path, opts)
	if err != nil {
		return fmt.Errorf("failed to collect files: %w", err)
	}

	// Calculate checksums for the collected files
	if err := p.metadata.CalculateSourceChecksum(files); err != nil {
		return fmt.Errorf("failed to calculate checksums: %w", err)
	}

	// Format output
	output := p.formatOutput(projectInfo, files)

	// Determine output location
	outputFile := opts.OutputFile
	if outputFile == "" {
		outputFile = p.determineOutputFile(path)
	}

	// Write output
	if outputFile != "" {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Context generated and saved to %s\n", outputFile)
	} else {
		fmt.Println(output)
	}

	return nil
}

func (p *ContextPlugin) parseFlags(cmd *cobra.Command) (*ContextOptions, error) {
	opts := &ContextOptions{}

	var err error

	opts.OutputFile, err = cmd.Flags().GetString("output")
	if err != nil {
		return nil, fmt.Errorf("error getting output flag: %w", err)
	}

	opts.StructureOnly, err = cmd.Flags().GetBool("structure-only")
	if err != nil {
		return nil, fmt.Errorf("error getting structure-only flag: %w", err)
	}

	opts.ContentOnly, err = cmd.Flags().GetBool("content-only")
	if err != nil {
		return nil, fmt.Errorf("error getting content-only flag: %w", err)
	}

	opts.Format, err = cmd.Flags().GetString("format")
	if err != nil {
		return nil, fmt.Errorf("error getting format flag: %w", err)
	}

	opts.MaxFiles, err = cmd.Flags().GetInt("max-files")
	if err != nil {
		return nil, fmt.Errorf("error getting max-files flag: %w", err)
	}

	opts.AdditionalIgnores, err = cmd.Flags().GetStringSlice("ignore")
	if err != nil {
		return nil, fmt.Errorf("error getting ignore patterns: %w", err)
	}

	// Validate flags
	if opts.StructureOnly && opts.ContentOnly {
		return nil, fmt.Errorf("cannot use both --structure-only and --content-only")
	}

	if opts.Format != "" && opts.Format != "md" && opts.Format != "txt" {
		return nil, fmt.Errorf("invalid format: %s (must be 'md' or 'txt')", opts.Format)
	}

	if opts.MaxFiles < 0 {
		return nil, fmt.Errorf("max-files must be >= 0")
	}

	return opts, nil
}

func (p *ContextPlugin) findExistingContext(path string) string {
	candidates := []string{
		filepath.Join(path, "context.md"),
		filepath.Join(path, "context.txt"),
	}

	for _, candidate := range candidates {
		if content, err := os.ReadFile(candidate); err == nil {
			if bytes.Contains(content, []byte(metadata.MetadataMarker)) {
				return candidate
			}
		}
	}

	return ""
}

func (p *ContextPlugin) determineOutputFile(path string) string {
	ext := ".md"
	if p.config.Context.OutputFormat == "txt" {
		ext = ".txt"
	}
	return filepath.Join(path, "context"+ext)
}

type ProjectInfo struct {
	Type      string
	GitBranch string
	GitStatus string
	HasGit    bool
}

func detectProject(path string) (*ProjectInfo, error) {
	info := &ProjectInfo{}

	// Detect Git
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		info.HasGit = true

		// Get git branch
		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		cmd.Dir = path
		if out, err := cmd.Output(); err == nil {
			info.GitBranch = strings.TrimSpace(string(out))
		}

		// Get git status
		cmd = exec.Command("git", "status", "--porcelain")
		cmd.Dir = path
		if out, err := cmd.Output(); err == nil {
			if len(out) == 0 {
				info.GitStatus = "clean"
			} else {
				info.GitStatus = "dirty"
			}
		}
	}

	// Project type detection
	switch {
	case fileExists(filepath.Join(path, "package.json")):
		info.Type = "nodejs"
	case fileExists(filepath.Join(path, "go.mod")):
		info.Type = "go"
	case fileExists(filepath.Join(path, "requirements.txt")):
		info.Type = "python"
	case fileExists(filepath.Join(path, "composer.json")):
		info.Type = "php"
	case fileExists(filepath.Join(path, "Cargo.toml")):
		info.Type = "rust"
	default:
		info.Type = "unknown"
	}

	return info, nil
}

func (p *ContextPlugin) collectFiles(root string, opts *ContextOptions) (map[string]string, error) {
	files := make(map[string]string)
	maxSize, err := filesize.Parse(p.config.Context.MaxFileSize)
	if err != nil {
		return nil, fmt.Errorf("invalid max file size: %w", err)
	}

	// Combine config ignore patterns with additional ones from command line
	ignorePatterns := make([]string, 0, len(p.config.Context.IgnorePatterns)+len(opts.AdditionalIgnores))
	ignorePatterns = append(ignorePatterns, p.config.Context.IgnorePatterns...)
	ignorePatterns = append(ignorePatterns, opts.AdditionalIgnores...)

	// Determine max files to process
	maxFiles := p.config.Context.MaxFilesToInclude
	if opts.MaxFiles > 0 {
		maxFiles = opts.MaxFiles
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			for _, pattern := range ignorePatterns {
				if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Skip files larger than max size
		if info.Size() > maxSize {
			return nil
		}

		// Check if path matches any ignore patterns
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		for _, pattern := range ignorePatterns {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				return nil
			}
			// Also try matching against the full relative path
			if matched, _ := filepath.Match(pattern, relPath); matched {
				return nil
			}
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(path))
		for _, excluded := range p.config.Context.ExcludeExtensions {
			if ext == excluded {
				return nil
			}
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Skip if file contains our metadata marker or is binary
		if shouldSkipFile(path, content) {
			return nil
		}

		files[relPath] = string(content)

		// Check if we've hit the max files limit
		if maxFiles > 0 && len(files) >= maxFiles {
			return io.EOF
		}

		return nil
	})

	if err == io.EOF {
		// We hit the max files limit
		fmt.Printf("Warning: Only including first %d files due to limit\n", maxFiles)
		return files, nil
	}

	return files, err
}

func shouldSkipFile(path string, content []byte) bool {
	// Skip binary files
	if !isText(content) {
		return true
	}

	// Skip if file contains our metadata marker
	if bytes.Contains(content, []byte(metadata.MetadataMarker)) {
		return true
	}

	return false
}

func (p *ContextPlugin) formatOutput(projectInfo *ProjectInfo, files map[string]string) string {
	var output strings.Builder

	// Add metadata header
	output.WriteString(p.metadata.String())
	output.WriteString("\n\n")

	// Add project info
	output.WriteString("# Project Information\n\n")
	output.WriteString(fmt.Sprintf("Type: %s\n", projectInfo.Type))
	if projectInfo.HasGit {
		output.WriteString(fmt.Sprintf("Git Branch: %s\n", projectInfo.GitBranch))
		output.WriteString(fmt.Sprintf("Git Status: %s\n", projectInfo.GitStatus))
	}
	output.WriteString("\n")

	// Add file structure
	if p.config.Context.IncludeFileStructure {
		output.WriteString("# File Structure\n\n```\n")
		paths := make([]string, 0, len(files))
		for path := range files {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			output.WriteString(path + "\n")
		}
		output.WriteString("```\n\n")
	}

	// Add file contents
	if p.config.Context.IncludeFileContent {
		output.WriteString("# File Contents\n\n")
		paths := make([]string, 0, len(files))
		for path := range files {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			ext := filepath.Ext(path)
			output.WriteString(fmt.Sprintf("## %s\n\n```%s\n%s\n```\n\n", path, ext[1:], files[path]))
		}
	}

	return output.String()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// isText reports whether a given content appears to be text data.
func isText(content []byte) bool {
	if len(content) == 0 {
		return true
	}

	// Check for null bytes which usually indicate binary data
	for _, b := range content {
		if b == 0 {
			return false
		}
	}

	return true
}
