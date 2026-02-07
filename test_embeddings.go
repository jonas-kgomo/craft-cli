package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ctxmgr "craft-cli/internal/context"
	"craft-cli/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found")
	}

	if err := logger.Init(); err != nil {
		log.Printf("Warning: logger init failed: %v", err)
	}
	defer logger.Close()

	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatalf("GEMINI_API_KEY is not set in environment or .env")
	}

	graph := ctxmgr.NewGraph()
	targetDir := "examples/vite-hello"

	fmt.Printf("🚀 Starting test index of: %s\n", targetDir)

	count := 0
	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Basic filter
		if strings.Contains(path, "node_modules") || strings.Contains(path, "dist") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("  [!] Could not read %s: %v\n", path, err)
			return nil
		}

		fmt.Printf("  Indexing [%d] %s...\n", count+1, path)
		err = graph.AddFile(path, string(content))
		if err != nil {
			fmt.Printf("  [!] Error indexing %s: %v\n", path, err)
		} else {
			count++
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Walk failed: %v", err)
	}

	fmt.Printf("\n✅ Successfully indexed %d nodes.\n", len(graph.Nodes))
	
	indexPath := ".craft-index.json"
	if err := graph.Save(indexPath); err != nil {
		fmt.Printf(" [!] Failed to save index: %v\n", err)
	} else {
		fmt.Printf(" ✨ Index persisted to %s\n", indexPath)
	}


	// Test Search
	testQuery := "How is the vite server configured?"
	if len(os.Args) > 1 {
		testQuery = strings.Join(os.Args[1:], " ")
	}
	fmt.Printf("\n🔍 Testing Search with query: %q\n", testQuery)


	results, err := graph.Search(testQuery, 3)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	fmt.Println("\nTop Results:")
	for i, res := range results {
		contentPreview := ""
		if res.Node.Content != "" {
			runes := []rune(res.Node.Content)
			if len(runes) > 100 {
				contentPreview = string(runes[:100]) + "..."
			} else {
				contentPreview = res.Node.Content
			}
		} else {
			contentPreview = "[Path Node]"
		}
		fmt.Printf("%d. [%.4f] %s\n   Snippet: %s\n\n", i+1, res.Similarity, res.Node.Path, strings.ReplaceAll(contentPreview, "\n", " "))
	}
}
