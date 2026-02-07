# Project Backlog & Roadmap

This document consolidates pending features, ideas, and enhancements from various sources.

## 🎨 UI/UX Enhancements (from `bubblesui.md`)

*   **File Picker**:
    *   Implement a bubble tea file picker for opening files in the diff viewer or for analysis.
    *   Navigate directory tree with arrow keys.
*   **Forms**:
    *   Add structured forms for complex configuration (e.g., setting up API keys, project init).
*   **Pager Improvements**:
    *   Better scrolling and search within the output pager.
    *   Line wrapping toggle.
*   **Spinner varieties**:
    *   More context-aware loading spinners (e.g., "Reading...", "Thinking...", "Writing...").
*   **Detailed Progress Bars**:
    *   For long-running tasks like indexing or batch processing.

## 🛠️ Diff Viewer Enhancements (from `DIFF_IMPLEMENTATION.md`)

*   **Recursive Directories**: Support diffing entire directories naturally.
*   **Merge Conflict Resolution**: UI for resolving git merge conflicts (3-way merge view).
*   **Syntax Highlighting**:
    *   Enhance the diff view with language-specific syntax highlighting (currently basic).
*   **Export Diffs**: Ability to save diff output to a patch file.
*   **Git Integration**:
    *   Direct `git diff` visualization.
    *   Stage/Unstage chunks from the UI.

## 🧠 Arrakis & Logic Code

*   **Self-Correction**:
    *   If a tool fails, Arrakis should autonomously propose a fix and retry (loop limit).
*   **Memory/History**:
    *   Better persistent memory of *user preferences* across sessions (beyond just project context).
*   **Async Planning**:
    *   Allow the planner to run in the background while the user continues to interact.

## 🏗️ Architecture & Infrastructure

*   **Plugin System**:
    *   Allow users to define custom tools via scripts or WASM plugins.
    *   "Skill" loading from local markdown files (partially implemented, needs refinement).
*   **Telemetry (Optional)**:
    *   Opt-in usage stats to help improve the tool.
*   **Testing**:
    *   Establish a robust end-to-end testing framework for the TUI (using `teatest` or similar).

## 💡 Musings & Ideas (from `coding_musings.txt`)

*   "Code as a conversation": integrating voice input?
*   "Visual coding": Graph view should be editable - drag and drop nodes to refactor?
*   Real-time collaboration: Multiplayer mode for pair programming with the AI?

## 📝 Documentation
*   [ ] Add a complete API reference for the Go internal packages.
*   [ ] Create a video tutorial series (GIFs are good, video is better).
