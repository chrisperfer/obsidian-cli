# Tag Search Feature Design

**Date:** 2025-11-23
**Status:** Approved
**Author:** Claude Code

## Problem Statement

Finding notes by tag is currently very difficult with the CLI. Users must either:
1. Run `vault find --recursive` to list all notes
2. For each note, run `note tag list [path]` to get its tags
3. Manually filter results for the desired tag

This requires hundreds of API calls for large vaults and has no built-in command.

## Requirements

**Must Have:**
- Find all notes with a specific tag (e.g., "flooby")
- Support standard output formats (text, JSON, YAML)
- Case-insensitive tag matching
- Handle all frontmatter tag formats (string, array, with/without # prefix)

**Should Have:**
- Optional path filtering to search within specific folders
- Performance optimization (frontmatter-only requests)
- Graceful handling of partial failures

**Future Expansion (v2):**
- Multiple tag search with AND logic (`--match-all`)
- Multiple tag search with OR logic (`--match-any`)

## Solution Design

### Command Structure

Add a new subcommand under `note tag`:

```bash
obsidian-cli note tag find [tag] [flags]
```

**Flags:**
- `--path string` - Optional folder filter (e.g., "Work/")
- `--output string` - Output format: text|json|yaml (inherited from global)
- `--llm-mode` - LLM-optimized output (inherited from global)

**Examples:**
```bash
# Find all notes tagged "flooby"
obsidian-cli note tag find flooby

# Search only in Work folder
obsidian-cli note tag find flooby --path "Work/"

# JSON output
obsidian-cli note tag find urgent --output json
```

### Architecture

**Three-layer approach:**

1. **Command Layer** (`cmd/note/tag.go`):
   - Validates arguments (require tag)
   - Parses flags (--path filter)
   - Calls client method
   - Formats output using existing formatter

2. **Client Layer** (`pkg/client/notes.go`):
   - New method: `FindNotesByTag(tag string, pathFilter string) ([]string, error)`
   - Fetches all notes via `ListAllFiles()`
   - Filters by path prefix if specified
   - Fetches each note's frontmatter
   - Extracts and matches tags
   - Returns matching paths

3. **Utility Layer** (`pkg/client/tags.go`):
   - New shared function: `ExtractTags(frontmatter map[string]interface{}) []string`
   - Handles all tag format variations
   - Normalizes tags (removes # prefix, lowercase for comparison)
   - Reused by both tag management and tag search

### Data Flow

```
User: obsidian-cli note tag find flooby --path "Work/"
  ↓
tagFindCmd validates args and flags
  ↓
client.FindNotesByTag("flooby", "Work/")
  ↓
List all files → filter by "Work/" prefix
  ↓
For each path: GetNote(path) with frontmatter-only
  ↓
ExtractTags(frontmatter) → normalize → check for "flooby"
  ↓
Return ["Work/project.md", "Work/notes.md"]
  ↓
Format and display results
```

### Implementation Details

**Tag Extraction Logic:**
Reuse existing logic from `cmd/note/tag.go:46-65`:
- Handle `[]interface{}` arrays
- Handle `[]string` arrays
- Handle single `string` values
- Strip `#` prefix from all tags
- Return normalized slice of strings

**Tag Matching:**
- Normalize both search tag and note tags to lowercase
- Strip `#` from search input if present
- Match using simple equality check

**Performance Optimizations:**
1. Request only frontmatter (not full content) to reduce bandwidth
2. Pre-filter paths by folder before fetching notes (if --path specified)
3. Continue processing on individual note failures (log warnings)
4. Consider concurrent fetching in future (goroutines with rate limiting)

**Error Handling:**
- Empty results: Return empty array (not an error)
- Note fetch failures: Log warning, continue processing
- Invalid path filter: Return error before making API calls
- API errors: Propagate with context

### Output Formats

**Text (default):**
```
Work/project.md
Work/notes.md
Archive/old-tasks.md
```

**JSON:**
```json
[
  "Work/project.md",
  "Work/notes.md",
  "Archive/old-tasks.md"
]
```

**YAML:**
```yaml
- Work/project.md
- Work/notes.md
- Archive/old-tasks.md
```

## Testing Strategy

**Unit Tests:**
- `ExtractTags()` with various frontmatter formats
- Tag normalization (case, # prefix)
- Path filtering logic

**Integration Tests:**
- Find tags in test vault
- Path filtering behavior
- Output format validation

**Manual Testing:**
- Large vault performance (1000+ notes)
- Edge cases (no tags, malformed frontmatter)
- Different tag formats in wild

## Future Enhancements (v2)

**Multiple Tag Search:**
```bash
# Find notes with BOTH tags
obsidian-cli note tag find flooby urgent --match-all

# Find notes with EITHER tag
obsidian-cli note tag find flooby urgent --match-any
```

**Implementation approach:**
- Accept variadic tags: `find [tags...]`
- Add flags: `--match-all` (AND), `--match-any` (OR)
- Default to `--match-all` when multiple tags provided
- Modify matching logic in `FindNotesByTag()`

**Performance Improvements:**
- Concurrent note fetching with goroutines
- Caching frontmatter for repeated searches
- Investigate if API supports tag search natively

## Migration Path

**No breaking changes:**
- New command, existing commands unchanged
- Follows established CLI patterns
- Compatible with existing output formats

**Documentation updates:**
- Add to README command reference
- Update tag management section
- Add examples for common use cases
