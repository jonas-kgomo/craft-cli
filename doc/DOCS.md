# CRAFT CLI Documentation

## 📘 User Guide

### Introduction
CRAFT CLI (Contextual Reasoning & Agentic Framework Tool) is a next-generation terminal assistant designed to pair-program with you. It combines high-speed inference (Groq) with deep reasoning capabilities (Arrakis) wrapped in a beautiful TUI.

### Getting Started

#### Installation
```bash
make install
export GROQ_API_KEY=gsk_...
```

#### Basic Usage
Simply type your request in the prompt.
*   **Standard Mode**: Good for quick questions, shell commands, and small edits.
*   **Arrakis Mode** (`/ARRAKIS` or Tab): Use this for complex refactoring, architectural changes, or when the agent needs to "think" before acting.

---

## 🔍 Diff Viewer Guide

The Diff Viewer is a powerful tool built into CRAFT CLI to help you review changes made by the AI before they become permanent. It allows side-by-side comparison of files with snapshot capability.

### Core Concepts

*   **Snapshot**: A temporary copy of a file stored in memory/disk. You take a snapshot *before* asking the agent to make changes.
*   **Compare**: Viewing the difference between the current file state on disk and its stored snapshot.
*   **Diff**: Generic comparison between any two files (File A vs File B).

### Workflow Example

1.  **Capture State**: Before asking for a refactor, save the current state.
    ```bash
    /snapshot main.go
    ```
2.  **Instruct Agent**: "Refactor main.go to use a singleton pattern."
3.  **Review Changes**: Once the agent finishes, compare the new file against your snapshot.
    ```bash
    /compare main.go
    ```
4.  **Accept/Reject**:
    *   If satisfied, you can commit via git.
    *   If not, you can undo (if you have git integration or manual backup) or ask the agent to fix it.

### Commands

| Command | Usage | Description |
| :--- | :--- | :--- |
| `/diff` | `/diff file1.go file2.go` | Opens side-by-side diff of two distinct files. |
| `/snapshot` | `/snapshot file.go` | Saves `file.go` as the "base" version for later comparison. |
| `/compare` | `/compare file.go` | Opens diff between the *snapshot* of `file.go` and the *current* `file.go`. |

### Keyboard Shortcuts (in Diff View)

*   `Tab` / `Shift+Tab`: Switch focus between Left (Original) and Right (Modified) panes.
*   `j` / `Down`: Scroll down.
*   `k` / `Up`: Scroll up.
*   `Space` / `PageDown`: Scroll down faster.
*   `b` / `PageUp`: Scroll up faster.
*   `Esc` / `q`: Close the diff viewer and return to the chat.

### Tips
*   Always `/snapshot` critical files before complex agent operations.
*   Use `/diff` to compare generated tests against implementation files if you want to spot inconsistencies manually.
