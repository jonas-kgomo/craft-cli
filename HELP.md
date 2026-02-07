# 🛠️ CRAFT CLI - Instruction Manual

CRAFT is a highly capable AI-driven CLI for code generation, refactoring, and project context analysis.

## 🚀 Getting Started

1. Set your environment variables in `.env`:
   - `GROQ_API_KEY`: Required for main agent and tool calling.
   - `GEMINI_API_KEY`: Required for context indexing and semantic search.
2. Run CRAFT:
   - CLI Mode: `go run main.go --theme sunset`
   - TUI Mode: (Coming Soon)

## 💬 Commands

- `/theme`: Toggle between Sunset (Dune) and Moonlit (Neon) themes.
- `/models`: List available models or switch (e.g., `/models llama-3.3-70b`).
- `/context index`: Index your project for semantic search.
- `/context status`: Show current index information.
- `/context clear`: Reset the context index.
- `exit` or `quit`: Close the session.

## 🏜️ Arrakis (The Voice) - Deep Thinking

Prefix any message with `@ARRAKIS` (or type `@` and press `Tab`) to activate the Deep Thinking Suite. Arrakis will analyze your task's complexity and decide whether to **Amplify** (compact reasoning) or **Decompose** (hierarchical planning).

- **Force Amplification**: `@ARRAKIS --amp <task>`
- **Force Decomposition**: `@ARRAKIS --dec <task>`

## 📂 Context Management

CRAFT automatically retrieves relevant file snippets from your indexed project to provide better answers. If the index feels outdated, run `/context index`.

## 📝 Logging

Everything is logged in `craft.log` and `scratchpad.md`. If CRAFT seems stuck, check these files for detailed background activity.
