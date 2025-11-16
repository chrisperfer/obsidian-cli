package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// Formatter handles output formatting
type Formatter struct {
	format  string
	llmMode bool
}

// NewFormatter creates a new output formatter
func NewFormatter(format string, llmMode bool) *Formatter {
	return &Formatter{
		format:  format,
		llmMode: llmMode,
	}
}

// Print outputs data in the configured format
func (f *Formatter) Print(data interface{}) error {
	switch f.format {
	case "json":
		return f.printJSON(data)
	case "yaml":
		return f.printYAML(data)
	case "text":
		return f.printText(data)
	default:
		return fmt.Errorf("unknown output format: %s", f.format)
	}
}

// printJSON outputs data as JSON
func (f *Formatter) printJSON(data interface{}) error {
	var output []byte
	var err error

	if f.llmMode {
		// Compact JSON for LLM mode
		output, err = json.Marshal(data)
	} else {
		// Pretty JSON for human reading
		output, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

// printYAML outputs data as YAML
func (f *Formatter) printYAML(data interface{}) error {
	output, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

// printText outputs data in human-readable text format
func (f *Formatter) printText(data interface{}) error {
	// Handle different data types
	switch v := data.(type) {
	case string:
		fmt.Println(v)
	case []string:
		for _, item := range v {
			fmt.Println(item)
		}
	default:
		// For complex types, fall back to JSON for LLM mode or pretty print
		if f.llmMode {
			return f.printJSON(data)
		}
		return f.printJSON(data)
	}
	return nil
}

// PrintSuccess prints a success message
func (f *Formatter) PrintSuccess(message string) {
	if f.format == "json" {
		result := map[string]interface{}{
			"success":   true,
			"message":   message,
			"timestamp": time.Now().Format(time.RFC3339),
		}
		f.Print(result)
	} else {
		if !f.llmMode {
			fmt.Printf("✓ %s\n", message)
		} else {
			fmt.Println(message)
		}
	}
}

// PrintError prints an error message
func (f *Formatter) PrintError(err error) {
	if f.format == "json" {
		result := map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"message": err.Error(),
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}
		f.Print(result)
	} else {
		if !f.llmMode {
			fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

// PrintTable prints data in a table format
func (f *Formatter) PrintTable(headers []string, rows [][]string) {
	if f.format == "json" {
		// Convert table to JSON structure
		data := make([]map[string]string, len(rows))
		for i, row := range rows {
			item := make(map[string]string)
			for j, header := range headers {
				if j < len(row) {
					item[header] = row[j]
				}
			}
			data[i] = item
		}
		f.Print(data)
		return
	}

	if f.format == "yaml" {
		// Convert table to YAML structure
		data := make([]map[string]string, len(rows))
		for i, row := range rows {
			item := make(map[string]string)
			for j, header := range headers {
				if j < len(row) {
					item[header] = row[j]
				}
			}
			data[i] = item
		}
		f.Print(data)
		return
	}

	// Text format - use simple text table
	table := tablewriter.NewWriter(os.Stdout)
	table.Append(headers)

	for _, row := range rows {
		table.Append(row)
	}

	table.Render()
}

// TruncateContent truncates content for LLM mode if needed
func (f *Formatter) TruncateContent(content string, maxLength int) string {
	if !f.llmMode {
		return content
	}

	if len(content) <= maxLength {
		return content
	}

	truncated := content[:maxLength]
	lines := strings.Split(truncated, "\n")
	return strings.Join(lines, "\n") + "\n\n... (truncated)"
}
