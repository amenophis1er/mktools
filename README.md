# mktools

A collection of powerful Makefile tools to enhance your build processes and development workflow.

## Features

- 📦 Self-contained installer with no git dependency
- 🎯 Modular target system with selective installation
- 🔄 Smart Makefile merging that preserves existing content
- 🎨 Colored output for better readability
- 📝 Comprehensive documentation generation
- 🔍 Project context dumping for documentation

## Quick Start

### Installation

Download and run the installer:

```bash
curl -sSL https://github.com/amenophis1er/mktools/releases/latest/download/install.sh | bash
```

The installer will:
1. Present available targets for installation
2. Preserve any existing Makefile content
3. Add selected mktools targets and their dependencies
4. Validate the resulting Makefile

### Basic Usage

```bash
# Show available targets
make help

# Show version information
make version

# Generate project context documentation
make dump
```

## Available Targets

### Core Targets

- `version`: Display mktools version and check for updates
- `dump`: Generate comprehensive project context documentation

### Target Groups

Basic targets (recommended for all projects):
- version

Development targets:
- dump

## Configuration

### Dump Target Configuration

Customize the dump target by setting these variables in your Makefile:

```makefile
# Add paths to exclude from dump
EXCLUDE_PATHS += \
    my/custom/path/** \
    *.custom.ext

# Add file extensions to exclude from content dump
CONTENT_EXCLUDE_EXT += \
    lock \
    conf
```

### Color Configuration

Disable colored output:

```makefile
# Disable colors
RESET :=
RED :=
GREEN :=
YELLOW :=
BLUE :=
PURPLE :=
CYAN :=
WHITE :=
```

## Advanced Usage

### Selective Target Installation

During installation, you can choose which targets to install:
```bash
# Download installer
curl -O https://github.com/amenophis1er/mktools/releases/latest/download/install.sh
chmod +x install.sh

# Run installer interactively
./install.sh
```

### Updating mktools

To update to the latest version:
1. Run `make version` to check for updates
2. Re-run the installer if an update is available

### Integration with Existing Makefiles

mktools intelligently merges with your existing Makefile:
- Preserves your existing targets and variables
- Adds mktools targets without conflicts
- Maintains existing .PHONY declarations
- Clearly marks mktools sections for easy maintenance

Example of merged Makefile:
```makefile
# Your existing comments are preserved
.PHONY: your-target mktools-target

# Your existing variables remain unchanged
YOUR_VAR := value

# BEGIN MKTOOLS VARIABLES
# mktools variables are isolated here
# END MKTOOLS VARIABLES

# Your existing targets remain unchanged
your-target:
    @echo "Your target"

# BEGIN MKTOOLS TARGETS
# mktools targets are isolated here
# END MKTOOLS TARGETS
```

## Contributing

Contributions are welcome! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
