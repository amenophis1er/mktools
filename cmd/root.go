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
It can generate context for LLMs, help with commit messages, and more.`,
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mktools/config.yaml)")

	// Add context command
	contextCmd := &cobra.Command{
		Use:   "context [path]",
		Short: "Generate context for LLM",
		Long: `Generate a context file containing project structure and file contents.
The context can be used to give LLMs better understanding of your project.

Examples:
  # Generate context for current directory
  mktools context

  # Generate context for specific directory
  mktools context ./my-project

  # Generate only structure in text format
  mktools context --structure-only --format txt

  # Generate to specific file
  mktools context -o project-context.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, ok := registry.Get("context")
			if !ok {
				return fmt.Errorf("context plugin not found")
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
		Use:   "config",
		Short: "Manage mktools configuration",
	}

	configInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig()
		},
	}

	configShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showConfig()
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
			fmt.Println(err)
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
