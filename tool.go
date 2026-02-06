package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Tool represents a callable function with metadata and execution logic
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Execute     func(args map[string]interface{}) (string, error)
	Category    string                 `json:"category,omitempty"`
	Timeout     time.Duration          `json:"timeout,omitempty"`
}

// ToolManager manages tool registration, validation, and execution
type ToolManager struct {
	tools map[string]Tool
}

// NewToolManager creates a new ToolManager instance
func NewToolManager() *ToolManager {
	return &ToolManager{
		tools: make(map[string]Tool),
	}
}

// Register adds a new tool to the manager
func (tm *ToolManager) Register(tool Tool) error {
	if tool.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if tool.Execute == nil {
		return fmt.Errorf("tool %s: Execute function cannot be nil", tool.Name)
	}
	if tool.Parameters == nil {
		tool.Parameters = make(map[string]interface{})
	}
	if tool.Timeout == 0 {
		tool.Timeout = 30 * time.Second // Default timeout
	}
	
	tm.tools[tool.Name] = tool
	return nil
}

// Get retrieves a tool by name
func (tm *ToolManager) Get(name string) (Tool, bool) {
	tool, exists := tm.tools[name]
	return tool, exists
}

// List returns all registered tool names
func (tm *ToolManager) List() []string {
	names := make([]string, 0, len(tm.tools))
	for name := range tm.tools {
		names = append(names, name)
	}
	return names
}

// Execute runs a tool with the given arguments
func (tm *ToolManager) Execute(name string, argsJSON string) (string, error) {
	tool, exists := tm.Get(name)
	if !exists {
		return "", fmt.Errorf("tool '%s' not found", name)
	}

	// Parse arguments
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments JSON: %w", err)
	}

	// Validate required parameters
	if err := tm.validateParameters(tool, args); err != nil {
		return "", fmt.Errorf("parameter validation failed: %w", err)
	}

	// Execute with timeout
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := tool.Execute(args)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- result
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return "", err
	case <-time.After(tool.Timeout):
		return "", fmt.Errorf("tool execution timed out after %v", tool.Timeout)
	}
}

// validateParameters checks if all required parameters are provided
func (tm *ToolManager) validateParameters(tool Tool, args map[string]interface{}) error {
	required, ok := tool.Parameters["required"].([]string)
	if !ok {
		return nil // No required parameters
	}

	for _, param := range required {
		if _, exists := args[param]; !exists {
			return fmt.Errorf("missing required parameter: %s", param)
		}
	}
	return nil
}

// GetToolDefinitions returns tools in the format expected by the API
func (tm *ToolManager) GetToolDefinitions() []ToolDef {
	var defs []ToolDef
	for _, tool := range tm.tools {
		defs = append(defs, ToolDef{
			Type: "function",
			Function: FunctionDef{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		})
	}
	return defs
}

// CreateDefaultTools creates the standard set of file system tools
func CreateDefaultTools() []Tool {
	return []Tool{
		{
			Name:        "read_file",
			Description: "Read contents of a file at the given path. Use this to examine code, configs, or documentation.",
			Category:    "filesystem",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{
						"type":        "string",
						"description": "Absolute or relative path to file",
					},
				},
				"required": []string{"path"},
			},
			Execute: func(args map[string]interface{}) (string, error) {
				path, ok := args["path"].(string)
				if !ok || path == "" {
					return "", fmt.Errorf("path must be a non-empty string")
				}
				
				content, err := os.ReadFile(path)
				if err != nil {
					return "", fmt.Errorf("failed to read file: %w", err)
				}
				return string(content), nil
			},
		},
		{
			Name:        "write_file",
			Description: "Write content to a file. Creates file if it doesn't exist, overwrites if it does.",
			Category:    "filesystem",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":    map[string]string{"type": "string"},
					"content": map[string]string{"type": "string"},
				},
				"required": []string{"path", "content"},
			},
			Execute: func(args map[string]interface{}) (string, error) {
				path, ok := args["path"].(string)
				if !ok || path == "" {
					return "", fmt.Errorf("path must be a non-empty string")
				}
				content, ok := args["content"].(string)
				if !ok {
					return "", fmt.Errorf("content must be a string")
				}
				
				dir := filepath.Dir(path)
				if dir != "." && dir != "/" {
					if err := os.MkdirAll(dir, 0755); err != nil {
						return "", fmt.Errorf("failed to create directory: %w", err)
					}
				}
				
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					return "", fmt.Errorf("failed to write file: %w", err)
				}
				return fmt.Sprintf("Successfully wrote to %s", path), nil
			},
		},
		{
			Name:        "list_dir",
			Description: "List files and directories at the given path.",
			Category:    "filesystem",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{
						"type":        "string",
						"description": "Directory path to list",
					},
				},
				"required": []string{"path"},
			},
			Execute: func(args map[string]interface{}) (string, error) {
				path, ok := args["path"].(string)
				if !ok || path == "" {
					return "", fmt.Errorf("path must be a non-empty string")
				}
				
				entries, err := os.ReadDir(path)
				if err != nil {
					return "", fmt.Errorf("failed to read directory: %w", err)
				}
				
				var result []string
				for _, e := range entries {
					prefix := "📄"
					if e.IsDir() {
						prefix = "📁"
					}
					result = append(result, fmt.Sprintf("%s %s", prefix, e.Name()))
				}
				return strings.Join(result, "\n"), nil
			},
		},
		{
			Name:        "bash",
			Description: "Execute a bash command. Use for git operations, running code, checking versions, etc.",
			Category:    "system",
			Timeout:     60 * time.Second,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]string{
						"type":        "string",
						"description": "The bash command to execute",
					},
				},
				"required": []string{"command"},
			},
			Execute: func(args map[string]interface{}) (string, error) {
				command, ok := args["command"].(string)
				if !ok || command == "" {
					return "", fmt.Errorf("command must be a non-empty string")
				}
				
				// Security: Block dangerous commands
				dangerous := []string{"rm -rf /", "mkfs", ":(){ :|:& };:", "dd if=/dev/zero"}
				for _, d := range dangerous {
					if strings.Contains(command, d) {
						return "", fmt.Errorf("dangerous command blocked for safety")
					}
				}
				
				// Additional security check
				if strings.Contains(command, "sudo") && !strings.Contains(command, "sudo -u") {
					return "", fmt.Errorf("sudo commands are restricted")
				}
				
				cmd := exec.Command("bash", "-c", command)
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Sprintf("Error: %v\nOutput: %s", err, string(output)), nil
				}
				return string(output), nil
			},
		},
		{
			Name:        "grep",
			Description: "Search for a pattern in file contents using simple string matching (case-insensitive).",
			Category:    "search",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pattern": map[string]string{
						"type":        "string",
						"description": "Search pattern",
					},
					"path": map[string]string{
						"type":        "string",
						"description": "Directory to search in",
					},
				},
				"required": []string{"pattern", "path"},
			},
			Execute: func(args map[string]interface{}) (string, error) {
				pattern, ok := args["pattern"].(string)
				if !ok || pattern == "" {
					return "", fmt.Errorf("pattern must be a non-empty string")
				}
				path, ok := args["path"].(string)
				if !ok || path == "" {
					return "", fmt.Errorf("path must be a non-empty string")
				}
				
				pattern = strings.ToLower(pattern)
				var matches []string
				
				err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
					if err != nil || info.IsDir() {
						return nil
					}
					content, err := os.ReadFile(file)
					if err != nil {
						return nil
					}
					if strings.Contains(strings.ToLower(string(content)), pattern) {
						matches = append(matches, file)
					}
					return nil
				})
				
				if err != nil {
					return "", fmt.Errorf("error walking directory: %w", err)
				}
				
				if len(matches) == 0 {
					return "No matches found", nil
				}
				
				return "Matches found in:\n" + strings.Join(matches, "\n"), nil
			},
		},
	}
}

// Helper function to safely get string from args
func getStringArg(args map[string]interface{}, key string) (string, error) {
	val, ok := args[key]
	if !ok {
		return "", fmt.Errorf("missing argument: %s", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("argument %s must be a string", key)
	}
	return str, nil
}