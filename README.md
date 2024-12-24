# mktools

mktools is a Swiss Army knife for development tasks, focusing on generating context for Large Language Models (LLMs) and automating common development workflows.

## Features

- **Context Generation**: Create comprehensive context files for LLM interactions
- **Project Analysis**: Automatically detect project type and structure
- **Configurable**: Support for global and project-specific configurations
- **Git Integration**: Automatic Git information inclusion
- **Smart Filtering**: Intelligent file filtering and content processing

## Installation

### Option 1: Homebrew (macOS/Linux)

```bash
# Install
brew tap amenophis1er/mktools
brew install mktools

# Update
brew upgrade mktools
```

### Option 2: Direct Download

Download the latest binary for your platform:

```bash
# macOS (Intel)
curl -L https://github.com/amenophis1er/mktools/releases/latest/download/mktools-darwin-amd64 -o /usr/local/bin/mktools && chmod +x /usr/local/bin/mktools

# macOS (Apple Silicon)
curl -L https://github.com/amenophis1er/mktools/releases/latest/download/mktools-darwin-arm64 -o /usr/local/bin/mktools && chmod +x /usr/local/bin/mktools

# Linux (x86_64)
curl -L https://github.com/amenophis1er/mktools/releases/latest/download/mktools-linux-amd64 -o /usr/local/bin/mktools && chmod +x /usr/local/bin/mktools

# Linux (ARM64)
curl -L https://github.com/amenophis1er/mktools/releases/latest/download/mktools-linux-arm64 -o /usr/local/bin/mktools && chmod +x /usr/local/bin/mktools

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/amenophis1er/mktools/releases/latest/download/mktools-windows-amd64.exe -OutFile mktools.exe
```

### Option 3: Go Install

```bash
go install github.com/amenophis1er/mktools@latest
```

## Version Management

mktools includes built-in version management features:

```bash
# Check current version
mktools version

# The tool automatically checks for updates and notifies you when a new version is available
```

## Verification

All releases include SHA-256 checksums for verification:

```bash
# Download checksums
curl -L https://github.com/amenophis1er/mktools/releases/latest/download/checksums.txt -O

# Verify binary (adjust filename for your platform)
sha256sum -c checksums.txt --ignore-missing
```

## Quick Start

```bash
# Generate context for current project
mktools context

# Initialize configuration
mktools config init

# Show current configuration
mktools config show
```

## Commands

### context

Generate context files for LLM interactions.

```bash
# Basic usage (current directory)
mktools context

# Specify directory
mktools context ./my-project

# Generate only structure
mktools context --structure-only

# Custom output file
mktools context -o project-context.md

# Change output format
mktools context --format txt

# Custom ignore patterns
mktools context --ignore "*.tmp" --ignore "build/*"
```

### config

Manage mktools configuration.

```bash
# Initialize default configuration
mktools config init

# Show current configuration
mktools config show
```

I'll update the Configuration section in the README.md to include the new features.


## Configuration

mktools supports both global and project-specific configurations, allowing you to set defaults globally and override them per project.

### Global Configuration

The global configuration file is located at `$HOME/.config/mktools/config.yaml`.

Initialize default global configuration:
```bash
mktools config init
```

### Project-specific Configuration

Create a project-specific `.mktools.yaml` file in your project root to override global settings:

```bash
# Initialize full config with all options
mktools config init --local

# Initialize minimal config (empty template)
mktools config init --local --minimal

# Force overwrite existing config
mktools config init --local --force
```

### Managing Configurations

View current configuration:
```bash
# Show current config
mktools config show

# Show effective merged configuration (global + local)
mktools config show --merged

# Show differences between global and local configs
mktools config diff
```

### Configuration Options

#### LLM Settings

| Option | Description | Default |
|--------|-------------|---------|
| provider | LLM provider (anthropic, openai) | anthropic |
| model | Model to use | claude-3-sonnet |
| api_key | API key (optional) | - |

#### Context Settings

| Option | Description | Default |
|--------|-------------|---------|
| output_format | Output format (md, txt) | md |
| ignore_patterns | Patterns to ignore | [".git/", "node_modules/", ...] |
| max_file_size | Maximum file size | 1MB |
| include_file_structure | Include directory structure | true |
| include_file_content | Include file contents | true |
| exclude_extensions | Extensions to exclude | [".exe", ".dll", ...] |
| max_files_to_include | Maximum files to process | 100 |

### Example Configurations

Global configuration (`~/.config/mktools/config.yaml`):
```yaml
llm:
  provider: anthropic
  model: claude-3-sonnet

context:
  output_format: md
  max_file_size: 1MB
  max_files_to_include: 100
  ignore_patterns:
    - ".git/"
    - "node_modules/"
    - "vendor/"
    - ".idea/"
  exclude_extensions:
    - ".exe"
    - ".dll"
    - ".so"
```

Project-specific configuration (`.mktools.yaml`):
```yaml
# Override only needed settings
context:
  ignore_patterns:
    - "build/*"
    - "*.tmp"
  max_file_size: 2MB
```

### Environment Variables

Environment variables take precedence over both global and local configurations:

- `MKTOOLS_LLM_PROVIDER`: Override LLM provider
- `MKTOOLS_LLM_MODEL`: Override LLM model
- `MKTOOLS_API_KEY`: Set API key
- `ANTHROPIC_API_KEY`: Anthropic-specific API key
- `OPENAI_API_KEY`: OpenAI-specific API key

## Output Formats

### Markdown (Default)

```markdown
# Project Information
Type: go
Git Branch: main
Git Status: clean

# File Structure
...

# File Contents
## main.go
```

### Text

Plain text format with minimal formatting.

## File Filtering

mktools automatically excludes:

- Binary files
- Large files (configurable)
- Common build artifacts
- Version control directories
- Dependency directories
- Temporary files
- 
### Ignore Patterns

mktools uses multiple sources to determine which files to ignore:

1. Built-in patterns (common binary files, build artifacts)
2. Project-specific `.mktools.yaml` configuration
3. `.gitignore` patterns in your project
4. Command-line `--ignore` patterns
5. Project-type specific patterns (e.g., node_modules for Node.js projects)

The patterns are processed in order, with later patterns taking precedence. `.gitignore` patterns are automatically respected, meaning any files ignored by Git will also be ignored by mktools.

## Development

### Prerequisites

- Go 1.21 or higher
- Git

### Building from Source

```bash
# Clone repository
git clone https://github.com/amenophis1er/mktools.git

# Build
cd mktools
go build

# Run tests
go test ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.