# mktools

mktools is a Swiss Army knife for development tasks, focusing on generating context for Large Language Models (LLMs) and automating common development workflows.

## Features

- **Context Generation**: Create comprehensive context files for LLM interactions
- **Project Analysis**: Automatically detect project type and structure
- **Configurable**: Support for global and project-specific configurations
- **Git Integration**: Automatic Git information inclusion
- **Smart Filtering**: Intelligent file filtering and content processing

## Installation

```bash
# Clone the repository
git clone https://github.com/amenophis1er/mktools.git

# Build and install
cd mktools
go install
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

## Configuration

mktools supports both global and project-specific configuration.

### Global Configuration

The global configuration file is located at `$HOME/.config/mktools/config.yaml`.

Initialize default configuration:
```bash
mktools config init
```

### Project-specific Configuration

Create a `.mktools.yaml` file in your project root to override global settings for that project.

Example `.mktools.yaml`:
```yaml
llm:
  provider: anthropic
  model: claude-3-sonnet

context:
  output_format: md
  ignore_patterns:
    - "*.tmp"
    - "build/*"
  max_file_size: 1MB
  include_file_structure: true
  include_file_content: true
  max_files_to_include: 100
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

### Environment Variables

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