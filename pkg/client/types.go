package client

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
	CTime int64 `json:"ctime"` // Unix timestamp in milliseconds
	MTime int64 `json:"mtime"` // Unix timestamp in milliseconds
	Size  int64 `json:"size"`
}

// VaultInfo represents vault information from the root API endpoint
type VaultInfo struct {
	Status        string        `json:"status"`
	Manifest      Manifest      `json:"manifest"`
	Versions      Versions      `json:"versions"`
	Service       string        `json:"service"`
	Authenticated bool          `json:"authenticated"`
	CertificateInfo CertificateInfo `json:"certificateInfo,omitempty"`
	APIExtensions []string      `json:"apiExtensions,omitempty"`
}

// Manifest represents the plugin manifest information
type Manifest struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	MinAppVersion string `json:"minAppVersion"`
	Description   string `json:"description"`
	Author        string `json:"author"`
	AuthorURL     string `json:"authorUrl"`
	IsDesktopOnly bool   `json:"isDesktopOnly"`
	Dir           string `json:"dir"`
}

// Versions represents version information
type Versions struct {
	Obsidian string `json:"obsidian"`
	Self     string `json:"self"`
}

// CertificateInfo represents SSL certificate information
type CertificateInfo struct {
	ValidityDays          float64 `json:"validityDays"`
	RegenerateRecommended bool    `json:"regenerateRecommended"`
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

// SearchMatch represents a single match within a search result
type SearchMatch struct {
	Match struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"match"`
	Context string `json:"context"`
}

// SearchResult represents a search result
type SearchResult struct {
	Filename string        `json:"filename"`
	Matches  []SearchMatch `json:"matches,omitempty"`
	Score    float64       `json:"score,omitempty"`
}
