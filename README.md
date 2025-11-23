# Obsidian CLI

A command-line interface for interacting with your Obsidian vault through the [Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) plugin.

## Features

- **Progressive Disclosure**: Commands and subcommands reveal capabilities at each level
- **LLM-Friendly**: Optimized for use with Large Language Models
- **Multiple Output Formats**: JSON, YAML, and human-readable text
- **Comprehensive API Coverage**: Notes, vault, periodic notes, and Obsidian commands

## Prerequisites

1. [Obsidian](https://obsidian.md/) installed
2. [Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) plugin installed and configured
3. API key from the Local REST API plugin

## Installation

### From Source

```bash
git clone https://github.com/chrisperfer/obsidian-cli.git
cd obsidian-cli
go build -o obsidian-cli
```

### Using Go Install

```bash
go install github.com/chrisperfer/obsidian-cli@latest
```

## Quick Start

1. **Initialize Configuration**

```bash
obsidian-cli config init
```

You'll be prompted for:
- API URL (default: `http://localhost:27124`)
- API Key (from Local REST API plugin settings)
- TLS verification settings

2. **Verify Setup**

```bash
obsidian-cli vault info
```

## Usage

### Command Structure

```
obsidian-cli [command] [subcommand] [args] [flags]
```

### Progressive Help System

Get help at any level:

```bash
# Top-level help
obsidian-cli --help

# Category-level help
obsidian-cli note --help

# Command-level help
obsidian-cli note get --help
```

## Commands

### Configuration

#### Initialize Configuration
```bash
obsidian-cli config init
```

#### Show Current Configuration
```bash
obsidian-cli config show
```

### Vault Operations

**Note**: The Obsidian Local REST API's `/vault/` endpoint only returns top-level items. The CLI works around this limitation by using the search API for recursive operations.

#### List Notes
```bash
# List top-level items only
obsidian-cli vault list

# List items in a specific folder
obsidian-cli vault list "Food and Drink/"

# List all files recursively
obsidian-cli vault list --recursive

# List all files under a path recursively
obsidian-cli vault list "Daily/" --recursive

# List in JSON format
obsidian-cli vault list --output json
```

#### Find Notes Recursively
```bash
# Find all notes in vault
obsidian-cli vault find

# Find all notes in a folder
obsidian-cli vault find "Daily/"

# Find with pattern matching
obsidian-cli vault find --pattern "*.md"
obsidian-cli vault find --pattern "2025-*.md"

# Combine path and pattern
obsidian-cli vault find "Daily/" --pattern "2025-11-*.md"

# Find in JSON format
obsidian-cli vault find --output json
```

#### Find Recent Notes
```bash
# Notes modified in last 3 days (default)
obsidian-cli vault recent

# Custom time periods
obsidian-cli vault recent --days 7
obsidian-cli vault recent --hours 24

# Recent notes in specific path
obsidian-cli vault recent "Daily/"
obsidian-cli vault recent "Food and Drink/" --days 30

# Sort by creation time instead of modification time
obsidian-cli vault recent --sort created

# Output in JSON format
obsidian-cli vault recent --output json
```

#### Get Vault Information
```bash
obsidian-cli vault info
```

### Note Operations

#### Get a Note
```bash
# Read a note
obsidian-cli note get "Daily/2024-01-15.md"

# Get only frontmatter
obsidian-cli note get "Project.md" --frontmatter-only

# Get only content
obsidian-cli note get "Ideas.md" --content-only
```

#### Create a Note
```bash
# Create with content
obsidian-cli note create "Ideas.md" --content "# My Ideas"

# Create from stdin
echo "# Daily Notes" | obsidian-cli note create "Daily/2024-01-15.md" --stdin

# Create empty note
obsidian-cli note create "Notes/Empty.md"
```

#### Update a Note
```bash
# Update with new content
obsidian-cli note update "Ideas.md" --content "# Updated Ideas"

# Update from stdin
cat content.md | obsidian-cli note update "Notes.md" --stdin
```

#### Patch a Note (Targeted Updates)
```bash
# Append content under a heading
obsidian-cli note patch "Notes.md" --target "TODO" --target-type heading --operation append --content "- New task"

# Update a frontmatter field
obsidian-cli note patch "Document.md" --target "status" --target-type frontmatter --operation replace --content "completed"

# Add value to frontmatter array
obsidian-cli note patch "Document.md" --target "tags" --array-add "new-tag"

# Remove value from frontmatter array
obsidian-cli note patch "Document.md" --target "tags" --array-remove "old-tag"
```

#### Manage Tags
```bash
# Add tags to a note
obsidian-cli note tag add "Document.md" "important" "urgent"

# Remove tags from a note
obsidian-cli note tag remove "Document.md" "old-tag"

# List all tags on a note
obsidian-cli note tag list "Document.md"

# List tags in JSON format
obsidian-cli note tag list "Document.md" --output json
```

#### Delete a Note
```bash
# Delete with confirmation
obsidian-cli note delete "Old.md"

# Force delete without confirmation
obsidian-cli note delete "Old.md" --force
```

#### Search Notes
```bash
# Search for notes containing text
obsidian-cli note search "project"

# Search with JSON output
obsidian-cli note search "meeting" --output json
```

### Periodic Notes

#### List Periodic Notes
```bash
# List daily notes
obsidian-cli periodic list daily

# List weekly notes
obsidian-cli periodic list weekly

# List monthly notes
obsidian-cli periodic list monthly
```

#### Get Periodic Note
```bash
# Get today's daily note
obsidian-cli periodic get daily

# Get specific date
obsidian-cli periodic get daily --date 2024-01-15

# Get weekly note
obsidian-cli periodic get weekly --date 2024-01-15
```

#### Create Periodic Note
```bash
# Create today's daily note
obsidian-cli periodic create daily

# Create with content and specific date
obsidian-cli periodic create daily --date 2024-01-15 --content "# Daily Note"

# Create weekly note
obsidian-cli periodic create weekly --content "# Week 3"
```

### Obsidian Commands

#### List Available Commands
```bash
# List all Obsidian commands
obsidian-cli commands list

# List in table format
obsidian-cli commands list --output text
```

#### Execute Command
```bash
# Execute a command by ID
obsidian-cli commands execute editor:toggle-bold
```

## Global Flags

### Output Format
```bash
--output, -o    Output format: text, json, yaml (default: text)
```

### LLM Mode
```bash
--llm-mode      Enable LLM-optimized output (concise, structured)
```

### Configuration File
```bash
--config        Specify config file (default: $HOME/.obsidian-cli.yaml)
```

## Configuration File

Default location: `~/.obsidian-cli.yaml`

```yaml
api:
  url: "http://localhost:27124"
  key: "your-api-key-here"
  insecure: false
  timeout: 30s

defaults:
  output: "text"
  vault_path: ""
  confirm_deletes: true

llm:
  mode: false
  max_content_length: 10000
```

## LLM Integration

The CLI is designed for seamless LLM integration:

### Structured JSON Output
```bash
obsidian-cli note get "Ideas.md" --output json
```

Returns:
```json
{
  "path": "Ideas.md",
  "content": "# My Ideas\n\n...",
  "frontmatter": {
    "tags": ["ideas"],
    "created": "2024-01-15"
  }
}
```

### LLM Mode
```bash
obsidian-cli note get "LongNote.md" --llm-mode
```

Features:
- Truncates long content automatically
- Minimal formatting
- Predictable structure

### Error Handling
```json
{
  "success": false,
  "error": {
    "message": "Note not found: NonExistent.md"
  }
}
```

## Examples

### Daily Workflow with LLM

```bash
# Get today's daily note
obsidian-cli periodic get daily --output json

# Find what you've added recently
obsidian-cli vault recent --days 3

# Find all daily notes from November 2025
obsidian-cli vault find "Daily/" --pattern "2025-11-*.md"

# Create a new note from LLM output
echo "$LLM_OUTPUT" | obsidian-cli note create "AI/Generated.md" --stdin

# Search and process results
obsidian-cli note search "TODO" --output json | jq '.[] | .filename'
```

### Batch Operations

```bash
# Find all markdown notes and process
obsidian-cli vault find --pattern "*.md" --output json | jq -r '.[]' | while read note; do
  echo "Processing: $note"
done

# Process all notes in a specific folder
obsidian-cli vault find "Daily/" --output json | jq -r '.[]' | while read note; do
  echo "Processing: $note"
done
```

## Environment Variables

Override config with environment variables (prefix with `OBSIDIAN_`):

```bash
export OBSIDIAN_API_URL="http://localhost:27124"
export OBSIDIAN_API_KEY="your-key"
```

## Troubleshooting

### Connection Issues

1. Ensure Local REST API plugin is enabled in Obsidian
2. Verify API URL and port
3. Check API key is correct

```bash
# Test with verbose output
obsidian-cli vault info --output json
```

### TLS Certificate Errors

For local development with self-signed certificates:

```yaml
api:
  insecure: true
```

**Warning**: Only use in development environments!

## Development

### Build from Source

```bash
git clone https://github.com/chrisperfer/obsidian-cli.git
cd obsidian-cli
go mod download
go build -o obsidian-cli
```

### Run Tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Obsidian](https://obsidian.md/) - The knowledge base application
- [Local REST API Plugin](https://github.com/coddingtonbear/obsidian-local-rest-api) - API integration for Obsidian
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management

## Support

- GitHub Issues: [Report bugs or request features](https://github.com/chrisperfer/obsidian-cli/issues)
- Documentation: [Full API docs](https://coddingtonbear.github.io/obsidian-local-rest-api/)

---

**Note**: This CLI requires the Obsidian Local REST API plugin. It is not affiliated with or endorsed by Obsidian.
