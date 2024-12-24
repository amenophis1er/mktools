// Example project-specific config structure for reference:
/*
# .mktools.yaml - Project-specific configuration
llm:
  provider: anthropic  # LLM provider (anthropic, openai)
  model: claude-3-sonnet  # Model to use
  api_key: ""  # Optional: Override API key

context:
  output_format: md  # Output format (md, txt)
  ignore_patterns:  # Additional patterns to ignore
    - "*.tmp"
    - "build/*"
    - "tests/fixtures/*"
  max_file_size: 1MB  # Maximum file size to include
  include_file_structure: true  # Include directory structure
  include_file_content: true   # Include file contents
  exclude_extensions:  # File extensions to exclude
    - ".exe"
    - ".dll"
    - ".so"
  max_files_to_include: 100  # Maximum number of files to process
*/

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type LLMConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	APIKey   string `yaml:"api_key"`
	Fallback *struct {
		Provider string `yaml:"provider"`
		Model    string `yaml:"model"`
		APIKey   string `yaml:"api_key"`
	} `yaml:"fallback,omitempty"`
}

type ContextConfig struct {
	OutputFormat         string   `yaml:"output_format"`
	IgnorePatterns       []string `yaml:"ignore_patterns"`
	MaxFileSize          string   `yaml:"max_file_size"`
	IncludeFileStructure bool     `yaml:"include_file_structure"`
	IncludeFileContent   bool     `yaml:"include_file_content"`
	ExcludeExtensions    []string `yaml:"exclude_extensions"`
	MaxFilesToInclude    int      `yaml:"max_files_to_include"`
}

type Config struct {
	LLM     LLMConfig     `yaml:"llm"`
	Context ContextConfig `yaml:"context"`
}

func LoadGlobal() (*Config, error) {
	config := DefaultConfig()
	globalConfig := filepath.Join(os.Getenv("HOME"), ".config", "mktools", "config.yaml")

	if err := loadFromFile(globalConfig, config); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading global config: %w", err)
		}
		// Return default config if global doesn't exist
		return config, nil
	}

	return config, nil
}

func LoadLocal() (*Config, error) {
	config := &Config{}
	if err := loadFromFile(".mktools.yaml", config); err != nil {
		return nil, err
	}
	return config, nil
}

func LoadMerged() (*Config, error) {
	// Start with global config
	config, err := LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("error loading global config: %w", err)
	}

	// Try to load and merge local config
	local, err := LoadLocal()
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // No local config, return global
		}
		return nil, fmt.Errorf("error loading local config: %w", err)
	}

	// Merge local into global
	if err := mergeConfig(config, local); err != nil {
		return nil, fmt.Errorf("error merging configs: %w", err)
	}

	return config, nil
}

func mergeConfig(dst, src *Config) error {
	if src == nil {
		return nil
	}

	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < dstValue.NumField(); i++ {
		dstField := dstValue.Field(i)
		srcField := srcValue.Field(i)

		if srcField.IsZero() {
			continue
		}

		switch dstField.Kind() {
		case reflect.Struct:
			if err := mergeStruct(dstField, srcField); err != nil {
				return err
			}
		default:
			dstField.Set(srcField)
		}
	}

	return nil
}

func mergeStruct(dst, src reflect.Value) error {
	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		srcField := src.Field(i)

		if srcField.IsZero() {
			continue
		}

		switch dstField.Kind() {
		case reflect.Slice:
			if !srcField.IsNil() {
				dstField.Set(srcField)
			}
		default:
			dstField.Set(srcField)
		}
	}

	return nil
}

func Diff(global, local *Config) (string, error) {
	var diff strings.Builder

	if local == nil {
		return "", nil
	}

	// Compare LLM config
	if d := diffLLM(&global.LLM, &local.LLM); d != "" {
		diff.WriteString("LLM Configuration:\n")
		diff.WriteString(d)
	}

	// Compare Context config
	if d := diffContext(&global.Context, &local.Context); d != "" {
		if diff.Len() > 0 {
			diff.WriteString("\n")
		}
		diff.WriteString("Context Configuration:\n")
		diff.WriteString(d)
	}

	return diff.String(), nil
}

func diffLLM(global, local *LLMConfig) string {
	var diff strings.Builder

	if local.Provider != "" && local.Provider != global.Provider {
		diff.WriteString(fmt.Sprintf("  provider: %s -> %s\n", global.Provider, local.Provider))
	}
	if local.Model != "" && local.Model != global.Model {
		diff.WriteString(fmt.Sprintf("  model: %s -> %s\n", global.Model, local.Model))
	}
	// Skip API key comparison for security

	return diff.String()
}

func diffContext(global, local *ContextConfig) string {
	var diff strings.Builder

	if local.OutputFormat != "" && local.OutputFormat != global.OutputFormat {
		diff.WriteString(fmt.Sprintf("  output_format: %s -> %s\n", global.OutputFormat, local.OutputFormat))
	}
	if local.MaxFileSize != "" && local.MaxFileSize != global.MaxFileSize {
		diff.WriteString(fmt.Sprintf("  max_file_size: %s -> %s\n", global.MaxFileSize, local.MaxFileSize))
	}
	if local.MaxFilesToInclude != 0 && local.MaxFilesToInclude != global.MaxFilesToInclude {
		diff.WriteString(fmt.Sprintf("  max_files_to_include: %d -> %d\n", global.MaxFilesToInclude, local.MaxFilesToInclude))
	}

	// Compare slices only if they're not empty in local config
	if len(local.IgnorePatterns) > 0 {
		diff.WriteString("  ignore_patterns: (added) [\n    ")
		diff.WriteString(strings.Join(local.IgnorePatterns, "\n    "))
		diff.WriteString("\n  ]\n")
	}
	if len(local.ExcludeExtensions) > 0 {
		diff.WriteString("  exclude_extensions: (added) [\n    ")
		diff.WriteString(strings.Join(local.ExcludeExtensions, "\n    "))
		diff.WriteString("\n  ]\n")
	}

	return diff.String()
}

// defaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			Provider: "anthropic",
			Model:    "claude-3-sonnet",
		},
		Context: ContextConfig{
			OutputFormat:         "md",
			IncludeFileStructure: true,
			IncludeFileContent:   true,
			MaxFileSize:          "1MB",
			MaxFilesToInclude:    100,
			IgnorePatterns: []string{
				".git/",
				"node_modules/",
				"vendor/",
				".idea/",
				"*.pyc",
				"*.pyo",
				"*.so",
				"*.dylib",
				"*.dll",
				"*.class",
				".DS_Store",
				"Thumbs.db",
				"*.swp",
				"*.swo",
				"*~",
				".env",
				"*.log",
			},
			ExcludeExtensions: []string{
				".exe", ".bin", ".o", ".a", ".lib", ".so", ".dylib", ".dll",
				".zip", ".tar", ".gz", ".7z", ".rar",
				".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico",
				".mp3", ".mp4", ".avi", ".mov",
				".pdf", ".doc", ".docx", ".xls", ".xlsx",
			},
		},
	}
}

// Load loads the configuration from files and environment variables
func Load() (*Config, error) {
	config := DefaultConfig()

	// Load from global config file
	globalConfig := filepath.Join(os.Getenv("HOME"), ".config", "mktools", "config.yaml")
	if err := loadFromFile(globalConfig, config); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading global config: %w", err)
	}

	// Load from local config file (overrides global)
	if err := loadFromFile(".mktools.yaml", config); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading local config: %w", err)
	}

	// Environment variables override file configs
	applyEnvironmentVariables(config)

	return config, validateConfig(config)
}

func loadFromFile(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// ToString returns a string representation of the config
func (c *Config) ToString() (string, error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}
	return string(data), nil
}

func applyEnvironmentVariables(config *Config) {
	if provider := os.Getenv("MKTOOLS_LLM_PROVIDER"); provider != "" {
		config.LLM.Provider = provider
	}
	if model := os.Getenv("MKTOOLS_LLM_MODEL"); model != "" {
		config.LLM.Model = model
	}
	if apiKey := os.Getenv("MKTOOLS_API_KEY"); apiKey != "" {
		config.LLM.APIKey = apiKey
	}

	// Provider-specific API keys
	if anthropicKey := os.Getenv("ANTHROPIC_API_KEY"); anthropicKey != "" && config.LLM.Provider == "anthropic" {
		config.LLM.APIKey = anthropicKey
	}
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" && config.LLM.Provider == "openai" {
		config.LLM.APIKey = openaiKey
	}
}

func validateConfig(config *Config) error {
	if config.LLM.Provider == "" {
		return fmt.Errorf("LLM provider is required")
	}
	if config.LLM.Model == "" {
		return fmt.Errorf("LLM model is required")
	}

	// Don't validate API key during initialization
	// API key can be set later via environment variable

	// Validate output format
	switch config.Context.OutputFormat {
	case "md", "txt":
		// valid
	default:
		return fmt.Errorf("invalid output format: %s", config.Context.OutputFormat)
	}

	return nil
}

// Save saves the current configuration to a file
func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
