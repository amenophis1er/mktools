package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amenophis1er/mktools/internal/config"
	"github.com/amenophis1er/mktools/internal/plugin"
	"github.com/amenophis1er/mktools/internal/update"
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
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mktools/config.yaml)")

	// Add version command
	rootCmd.AddCommand(newVersionCmd())

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
		Use:   "config",
		Short: "Manage mktools configuration",
		Long: `Manage mktools configuration settings.

Available Commands:
  init    Initialize configuration file
  show    Display current configuration
  diff    Show differences between global and local config`,
	}

	configInitCmd := &cobra.Command{
		Use:   "init [--local] [--minimal] [--force]",
		Short: "Initialize configuration file",
		Long: `Initialize a new configuration file with default settings.
Without flags, the configuration file will be created at $HOME/.config/mktools/config.yaml.
With --local flag, it will create .mktools.yaml in the current directory.`,
		Example: `  # Initialize global config
  mktools config init

  # Initialize local project config
  mktools config init --local

  # Initialize minimal local config (only override specific settings)
  mktools config init --local --minimal

  # Force overwrite existing config
  mktools config init --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			local, _ := cmd.Flags().GetBool("local")
			minimal, _ := cmd.Flags().GetBool("minimal")
			force, _ := cmd.Flags().GetBool("force")
			if err := initializeConfig(local, minimal, force); err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
			return nil
		},
	}

	configInitCmd.Flags().Bool("local", false, "create config file in current directory")
	configInitCmd.Flags().Bool("minimal", false, "create minimal config with only overridden settings")
	configInitCmd.Flags().Bool("force", false, "overwrite existing config file")

	configShowCmd := &cobra.Command{
		Use:   "show [--merged]",
		Short: "Show current configuration",
		Long: `Display the current active configuration settings.
By default, shows the configuration file specified.
With --merged flag, shows the effective configuration after merging global and local settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			merged, _ := cmd.Flags().GetBool("merged")
			if err := showConfig(merged); err != nil {
				return fmt.Errorf("failed to show config: %w", err)
			}
			return nil
		},
	}

	configShowCmd.Flags().Bool("merged", false, "show merged configuration")

	configDiffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Show configuration differences",
		Long: `Compare global and local configuration settings.
Shows what settings are different in the local configuration
compared to the global configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := diffConfig(); err != nil {
				return fmt.Errorf("failed to diff config: %w", err)
			}
			return nil
		},
	}

	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configDiffCmd)
	rootCmd.AddCommand(configCmd)

	// Check for updates in the background
	go func() {
		hasUpdate, newVersion, err := update.CheckForUpdate()
		if err != nil {
			return
		}
		if hasUpdate {
			fmt.Fprintf(os.Stderr, "\nNew version %s available! Run 'mktools version' for update instructions.\n", newVersion)
		}
	}()

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

func initializeConfig(local, minimal, force bool) error {
	defaultConfig := config.DefaultConfig()

	var configPath string
	if local {
		configPath = ".mktools.yaml"
	} else {
		configDir := filepath.Join(os.Getenv("HOME"), ".config", "mktools")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		configPath = filepath.Join(configDir, "config.yaml")
	}

	// Check if file exists and handle force flag
	if _, err := os.Stat(configPath); err == nil && !force {
		return fmt.Errorf("config file already exists at %s (use --force to overwrite)", configPath)
	}

	if minimal && local {
		// Create minimal config with empty structures
		defaultConfig = &config.Config{
			LLM:     config.LLMConfig{},
			Context: config.ContextConfig{},
		}
	}

	if err := defaultConfig.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Configuration initialized at %s\n", configPath)
	return nil
}

func showConfig(merged bool) error {
	if cfg == nil {
		return fmt.Errorf("no configuration loaded")
	}

	var configToShow *config.Config
	if merged {
		// Get merged config (global + local)
		var err error
		configToShow, err = config.LoadMerged()
		if err != nil {
			return fmt.Errorf("failed to load merged config: %w", err)
		}
	} else {
		configToShow = cfg
	}

	output, err := configToShow.ToString()
	if err != nil {
		return fmt.Errorf("failed to format config: %w", err)
	}

	fmt.Println(output)
	return nil
}

func diffConfig() error {
	// Load global config
	globalConfig, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Load local config
	localConfig, err := config.LoadLocal()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load local config: %w", err)
	}

	if localConfig == nil {
		fmt.Println("No local configuration found")
		return nil
	}

	// Compare and show differences
	diff, err := config.Diff(globalConfig, localConfig)
	if err != nil {
		return fmt.Errorf("failed to compare configs: %w", err)
	}

	if diff == "" {
		fmt.Println("No differences found between global and local configuration")
		return nil
	}

	fmt.Println("Configuration differences (local vs global):")
	fmt.Println(diff)
	return nil
}
