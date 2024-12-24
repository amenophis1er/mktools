package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amenophis1er/mktools/internal/config"
	"github.com/amenophis1er/mktools/internal/plugin"
	"github.com/amenophis1er/mktools/plugins/context"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	rootCmd *cobra.Command
)

func Execute() error {
	var err error
	cfg, err = config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	registry := plugin.NewRegistry()

	// Register plugins
	registry.Register(context.New(cfg))

	rootCmd = &cobra.Command{
		Use:   "mktools",
		Short: "Swiss army knife for development tasks",
		Long: `mktools is a collection of tools to help with development tasks.
It provides various utilities for:
- Generating context for LLMs
- Project structure analysis
- Development workflow automation

Use "mktools [command] --help" for more information about a command.`,
		SilenceErrors: true, // Let main() handle error output
		SilenceUsage:  true, // Don't show usage on expected errors
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mktools/config.yaml)")

	// Add context command
	contextCmd := &cobra.Command{
		Use:   "context [flags] [path]",
		Short: "Generate context for LLM",
		Long: `Generate a context file containing project structure and file contents.
The context can be used to give LLMs better understanding of your project.

The command will analyze the specified directory (or current directory if not specified)
and generate a markdown or text file containing:
- Project type detection
- Git information (if available)
- File structure
- File contents (configurable)

By default, binary files, large files, and common build artifacts are excluded.`,
		Example: `  # Generate context for current directory
  mktools context

  # Generate context for specific directory
  mktools context ./my-project

  # Generate only structure in text format
  mktools context --structure-only --format txt

  # Generate to specific file
  mktools context -o project-context.md
  
  # Generate with custom ignore patterns
  mktools context --ignore "*.tmp" --ignore "build/*"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, ok := registry.Get("context")
			if !ok {
				return fmt.Errorf("internal error: context plugin not found")
			}

			// Convert any explicit help request to cobra's help command
			if len(args) == 1 && (args[0] == "help" || args[0] == "--help") {
				return cmd.Help()
			}

			return p.Execute(cmd.Context(), cmd, args)
		},
	}

	// Let context plugin add its flags
	if p, ok := registry.Get("context"); ok {
		p.AddFlags(contextCmd)
	}

	rootCmd.AddCommand(contextCmd)

	// Add config command
	configCmd := &cobra.Command{
		Use:   "config [command]",
		Short: "Manage mktools configuration",
		Long: `Manage mktools configuration settings.

Available Commands:
  init    Initialize default configuration file
  show    Display current configuration`,
	}

	configInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		Long: `Initialize a new configuration file with default settings.
The configuration file will be created at $HOME/.config/mktools/config.yaml
if it doesn't already exist.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initializeConfig(); err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
			return nil
		},
	}

	configShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long: `Display the current active configuration settings.
This includes both default values and any overrides from:
- Global config file ($HOME/.config/mktools/config.yaml)
- Local config file (.mktools.yaml)
- Environment variables`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := showConfig(); err != nil {
				return fmt.Errorf("failed to show config: %w", err)
			}
			return nil
		},
	}

	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
	}
}

func initializeConfig() error {
	defaultConfig := config.DefaultConfig()
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "mktools")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	if err := defaultConfig.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Configuration initialized at %s\n", configPath)
	return nil
}

func showConfig() error {
	if cfg == nil {
		return fmt.Errorf("no configuration loaded")
	}

	currentConfig, err := cfg.ToString()
	if err != nil {
		return fmt.Errorf("failed to format config: %w", err)
	}

	fmt.Println(currentConfig)
	return nil
}
