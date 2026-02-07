# 🚀 Quick Start

Run CRAFT CLI:

```bash
./run.sh
```

Or directly:

```bash
./craft
```

## First Run

The TUI will show:
- 🛠️ Header with current model
- Token budget bar (green/yellow/red)
- Chat history viewport
- Input box at bottom
- Status line with session stats

## Controls

- **Enter** - Send message
- **Alt+Enter** - New line in input
- **Ctrl+C** or **Esc** - Quit

## Example Session

```
> create a hello world in go

CRAFT: I'll create that for you.
  🔧 write_file({"path":"hello.go",...})
  → ✓ Wrote hello.go (85B)
CRAFT: Created hello.go with a simple program.
```

## Token Budget

Watch the bar at the top:
- 🟢 Green (0-66%): Plenty of room
- 🟡 Yellow (66-85%): Getting full
- 🔴 Red (85-100%): Near limit

History auto-evicts after 4 turns to stay under budget.

## Models

Auto-fallback on rate limits:
1. llama-3.1-8b-instant (fastest)
2. llama-3.3-70b-versatile (smarter)
3. mixtral-8x7b-32768 (fallback)

Current model shown in header.
