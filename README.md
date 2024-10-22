# mktools

A collection of reusable Makefile targets and utilities.

## Description

mktools provides a collection of reusable Makefile targets that you can easily add to your projects. Each target is self-contained with its own variables and documentation.

## Directory Structure

```
mktools/
├── README.md           # This file
├── common/            # Common variables and utilities
│   └── colors.mk      # Terminal colors for pretty output
├── install/           # Installation related files
│   └── install.sh     # Installation script
└── targets/          # Collection of available targets
└── dump/         # Example: dump target
├── dump.mk   # The target implementation
└── vars.mk   # Target specific variables
```

## Installation

You can install mktools using this command:

```bash
curl -sSL https://raw.githubusercontent.com/amenophis1er/mktools/main/install/install.sh | bash
```

The installation script will:
- Clone the repository to `~/.local/share/mktools`
- Create the `mktools` command in `~/.local/bin`
- Add `~/.local/bin` to your PATH if needed

Note: You may need to restart your shell or update your PATH after installation.

## Usage

List available targets:
```bash
mktools list
```

Install a target:
```bash
mktools install <target-name>
```

### Available Targets

#### dump
Creates a context dump of your project structure, including:
- Directory structure
- File listing
- Content of text files (excluding binary files and common formats like images, pdfs, etc)

Usage after installation:
```bash
make dump
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-target`)
3. Commit your changes (`git commit -am 'Add some amazing-target'`)
4. Push to the branch (`git push origin feature/amazing-target`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
