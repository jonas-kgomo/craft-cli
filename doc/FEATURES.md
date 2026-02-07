# CRAFT CLI Features

This document outlines the currently implemented and working features of the CRAFT CLI.

## Core Capabilities

### 🧠 Arrakis: Deep Logic Engine
*   **Mode Switching**: Toggle between "Normal" (fast) and "Arrakis" (deep reasoning) modes using `Tab` or `/ARRAKIS` command.
*   **Factored Cognition**: Decomposes complex tasks into atomic steps using a "Planner -> Executor -> Reviewer" loop.
*   **Traceability**: Transparent logging of the agent's thought process (Plan, Thought, Observation).
*   **Persistent Context**: Maintains awareness of the project structure and history across sessions.

### ⚡ Groq Integration
*   **High-Speed Inference**: Powered by Groq's LPU inference engine for near-instant responses.
*   **Multi-Model Support**:
    *   **Llama 3.3 70B**: Primary model for high-fidelity reasoning and coding.
    *   **Mixtral 8x7B**: Fallback model for speed and efficiency.
    *   **Auto-Fallback**: Automatically switches models if rate limits are hit.
*   **Token Budgeting**: Real-time tracking of token usage with visual indicators.

### 🖥️ Lite UI (TUI)
*   **Interactive Interface**: Built with Bubble Tea for a rich terminal user experience.
*   **Slash Commands**: quick access to tools and settings (e.g., `/models`, `/vis`, `/diff`).
*   **Banners & Themes**: Dynamic ASCII art banners and color schemes (Sunset/Moonlit themes).
*   **Markdown Rendering**: Rich text rendering for agent responses, including code blocks and tables.
*   **Input Handling**: robust multi-line input with history navigation.

### 🔍 Diff Viewer & Snapshots
*   **Snapshot System**: Capture the state of files before making interactions using `/snapshot <file>`.
*   **Side-by-Side Diff**: Compare file versions using `/diff <file1> <file2>` or `/compare <file>` (against last snapshot).
*   **Visual Diffing**: Color-coded insertions and deletions for easy review of agent changes.
*   **Navigation**: Tab between panes, scroll lines/pages, and dynamic resizing.

### 🛡️ Safety & Security
*   **Sandboxed Execution**: Dangerous commands (e.g., `rm -rf`, `sudo`) are blocked by default.
*   **User Confirmation**: Critical actions require explicit user approval (unless configured otherwise).
*   **Output Truncation**: Prevents terminal flooding by truncating large file reads or command outputs (configurable).

### 📂 Context & Knowledge
*   **Context Graph**: Semantic graph allowing the agent to understand relationships between files and symbols.
*   **RAG (Retrieval Augmented Generation)**: Retrieves relevant code snippets and documentation to ground agent responses.
*   **Persistence**: Saves and loads the context index to speed up startup times.

## Command Reference

| Command | Description |
| :--- | :--- |
| `/ARRAKIS` | Activate Deep Logic mode for complex tasks |
| `/models` | Interactive menu to switch LLMs |
| `/vis` | Visualize the current context graph and Arrakis flow |
| `/diff [f1] [f2]` | Open side-by-side diff viewer for two files |
| `/snapshot [file]` | Save current state of a file for later comparison |
| `/compare [file]` | Compare current file against its last snapshot |
| `/help` | Show available commands and shortcuts |
| `/quit` | Exit the application |
