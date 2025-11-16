package client

import "time"

// Note represents an Obsidian note
type Note struct {
	Path        string                 `json:"path"`
	Content     string                 `json:"content"`
	Frontmatter map[string]interface{} `json:"frontmatter,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Stat        *FileStat              `json:"stat,omitempty"`
}

// FileStat represents file statistics
type FileStat struct {
	CTime time.Time `json:"ctime"`
	MTime time.Time `json:"mtime"`
	Size  int64     `json:"size"`
}

// VaultInfo represents vault information
type VaultInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Files int    `json:"files"`
}

// FileList represents a list of files
type FileList struct {
	Files []string `json:"files"`
}

// PeriodicNote represents a periodic note request/response
type PeriodicNote struct {
	Type string `json:"type"` // daily, weekly, monthly
	Date string `json:"date,omitempty"`
	Path string `json:"path,omitempty"`
}

// Command represents an Obsidian command
type Command struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CommandList represents a list of commands
type CommandList struct {
	Commands []Command `json:"commands"`
}

// CommandExecuteRequest represents a command execution request
type CommandExecuteRequest struct {
	CommandID string `json:"command"`
}

// SearchResult represents a search result
type SearchResult struct {
	Filename string   `json:"filename"`
	Matches  []string `json:"matches,omitempty"`
	Score    float64  `json:"score,omitempty"`
}

// SearchResults represents search results
type SearchResults struct {
	Results []SearchResult `json:"results"`
}
