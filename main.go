package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"craft-cli/internal/agent"
	ctxmgr "craft-cli/internal/context"
	"craft-cli/internal/groq"
	"craft-cli/internal/logger"
	"craft-cli/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	godotenv.Load()
	lipgloss.SetColorProfile(termenv.TrueColor)
	
	if err := logger.Init(); err != nil {
		fmt.Printf("Warning: could not initialize logger: %v\n", err)
	}
	defer logger.Close()

	themeArg := flag.String("theme", "sunset", "UI theme: sunset or moonlit")
	flag.Parse()

	if strings.ToLower(*themeArg) == "moonlit" {
		tui.ApplyTheme(tui.MoonlitTheme)
	} else {
		tui.ApplyTheme(tui.SunsetTheme)
	}

	if os.Getenv("GROQ_API_KEY") == "" {
		fmt.Println(tui.ErrorStyle.Render(" [!] Set GROQ_API_KEY"))
		os.Exit(1)
	}

	// Initialize components
	client := groq.NewClient()
	toolMgr := agent.NewToolManager(client)
	ctxGraph := ctxmgr.NewGraph()
	indexPath := ".craft-index.json"
	
	// Load existing context index
	if err := ctxGraph.Load(indexPath); err == nil {
		logger.Infof("Loaded existing context index with %d nodes", len(ctxGraph.Nodes))
	}

	// Create the Lite UI Model
	m := tui.NewLiteModel(client, toolMgr, ctxGraph)
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	m.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running CRAFT: %v\n", err)
		os.Exit(1)
	}
}
