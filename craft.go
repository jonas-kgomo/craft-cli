package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Tool represents a callable function
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Execute     func(args map[string]interface{}) string
}

// GroqClient — minimal HTTP implementation
type GroqClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
	tools   []Tool
}

func NewGroqClient() *GroqClient {
	return &GroqClient{
		apiKey:  os.Getenv("GROQ_API_KEY"),
		baseURL: "https://api.groq.com/openai/v1",
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

// Message types for Groq API
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []ToolDef `json:"tools,omitempty"`
}

type ToolDef struct {
	Type     string   `json:"type"`
	Function FunctionDef `json:"function"`
}

type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ChatResponse struct {
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Initialize tools
func (g *GroqClient) initTools() {
	g.tools = []Tool{
		{
			Name:        "read_file",
			Description: "Read contents of a file at the given path. Use this to examine code, configs, or documentation.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string", "description": "Absolute or relative path to file"},
				},
				"required": []string{"path"},
			},
			Execute: func(args map[string]interface{}) string {
				path := args["path"].(string)
				content, err := os.ReadFile(path)
				if err != nil {
					return fmt.Sprintf("Error reading file: %v", err)
				}
				return string(content)
			},
		},
		{
			Name:        "write_file",
			Description: "Write content to a file. Creates file if it doesn't exist, overwrites if it does.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":    map[string]string{"type": "string"},
					"content": map[string]string{"type": "string"},
				},
				"required": []string{"path", "content"},
			},
			Execute: func(args map[string]interface{}) string {
				path := args["path"].(string)
				content := args["content"].(string)
				dir := filepath.Dir(path)
				if dir != "." && dir != "/" {
					os.MkdirAll(dir, 0755)
				}
				err := os.WriteFile(path, []byte(content), 0644)
				if err != nil {
					return fmt.Sprintf("Error writing file: %v", err)
				}
				return fmt.Sprintf("Successfully wrote to %s", path)
			},
		},
		{
			Name:        "list_dir",
			Description: "List files and directories at the given path.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
			Execute: func(args map[string]interface{}) string {
				path := args["path"].(string)
				entries, err := os.ReadDir(path)
				if err != nil {
					return fmt.Sprintf("Error reading directory: %v", err)
				}
				var result []string
				for _, e := range entries {
					prefix := "📄"
					if e.IsDir() {
						prefix = "📁"
					}
					result = append(result, fmt.Sprintf("%s %s", prefix, e.Name()))
				}
				return strings.Join(result, "\n")
			},
		},
		{
			Name:        "bash",
			Description: "Execute a bash command. Use for git operations, running code, checking versions, etc.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]string{"type": "string", "description": "The bash command to execute"},
				},
				"required": []string{"command"},
			},
			Execute: func(args map[string]interface{}) string {
				command := args["command"].(string)
				// Safety: block dangerous commands
				dangerous := []string{"rm -rf /", "mkfs", ":(){ :|:& };:", "dd if=/dev/zero"}
				for _, d := range dangerous {
					if strings.Contains(command, d) {
						return "Error: Dangerous command blocked for safety"
					}
				}
				cmd := exec.Command("bash", "-c", command)
				output, err := cmd.CombinedOutput()
				if err != nil {
					return fmt.Sprintf("Error: %v\nOutput: %s", err, string(output))
				}
				return string(output)
			},
		},
		{
			Name:        "grep",
			Description: "Search for a pattern in file contents using simple string matching (case-insensitive).",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"pattern": map[string]string{"type": "string"},
					"path":    map[string]string{"type": "string", "description": "Directory to search in"},
				},
				"required": []string{"pattern", "path"},
			},
			Execute: func(args map[string]interface{}) string {
				pattern := strings.ToLower(args["pattern"].(string))
				path := args["path"].(string)
				var matches []string
				filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
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
				if len(matches) == 0 {
					return "No matches found"
				}
				return "Matches found in:\n" + strings.Join(matches, "\n")
			},
		},
	}
}

func (g *GroqClient) toToolDefs() []ToolDef {
	var defs []ToolDef
	for _, t := range g.tools {
		defs = append(defs, ToolDef{
			Type: "function",
			Function: FunctionDef{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		})
	}
	return defs
}

func (g *GroqClient) executeTool(name string, args string) string {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	for _, t := range g.tools {
		if t.Name == name {
			return t.Execute(parsed)
		}
	}
	return fmt.Sprintf("Unknown tool: %s", name)
}

func (g *GroqClient) Chat(ctx context.Context, messages []Message) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:    "llama-3.1-8b-instant",
		Messages: messages,
		Tools:    g.toToolDefs(),
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, fmt.Errorf("API error: %s", result.Error.Message)
	}

	return &result, nil
}

func getSystemPrompt() string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf(`You are CRAFT CLI, an AI coding assistant with direct access to the local filesystem.
You have access to these tools:
- read_file: Read file contents
- write_file: Write/create files  
- list_dir: List directory contents
- bash: Execute shell commands
- grep: Search file contents

Current working directory: %s
When you need to explore or modify files, use the tools directly. Always confirm successful file operations.`, cwd)
}

func main() {
	godotenv.Load()

	client := NewGroqClient()
	client.initTools()

	if client.apiKey == "" {
		fmt.Println("❌ Set GROQ_API_KEY environment variable")
		os.Exit(1)
	}

	fmt.Println("🛠️  CRAFT CLI")
	fmt.Println("Model: Model K")
	fmt.Println("Tools: read_file, write_file, list_dir, bash, grep")
	fmt.Println("Type 'exit' to quit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	var history []Message

	// Add system prompt
	history = append(history, Message{
		Role:    "system",
		Content: getSystemPrompt(),
	})

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" {
			break
		}

		history = append(history, Message{Role: "user", Content: input})

		// Agent loop: keep calling until no more tool calls
		for {
			fmt.Print("Thinking... ")
			resp, err := client.Chat(context.Background(), history)
			fmt.Print("\r") // Clear "Thinking..."

			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				break
			}

			assistantMsg := resp.Choices[0].Message
			history = append(history, assistantMsg)

			// Check if there are tool calls
			if len(assistantMsg.ToolCalls) == 0 {
				// Final response
				fmt.Println(assistantMsg.Content)
				break
			}

			// Execute tools and add results to history
			fmt.Printf("🔧 Using %d tool(s)...\n", len(assistantMsg.ToolCalls))
			for _, tc := range assistantMsg.ToolCalls {
				fmt.Printf("  → %s(%s)\n", tc.Function.Name, tc.Function.Arguments)
				result := client.executeTool(tc.Function.Name, tc.Function.Arguments)
				
				// Truncate long results for display
				display := result
				if len(display) > 200 {
					display = display[:200] + "... (truncated)"
				}
				fmt.Printf("  ← %s\n", display)

				history = append(history, Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				})
			}
		}
	}
}