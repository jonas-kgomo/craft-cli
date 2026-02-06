#!/bin/bash

# Test CRAFT CLI
echo "Testing CRAFT CLI..."
echo ""
echo "The TUI should now:"
echo "  ✅ Not flicker or blink"
echo "  ✅ Use llama-3.3-70b (better tool support)"
echo "  ✅ Not crash on nil pointer"
echo ""
echo "Try asking: 'list files in current directory'"
echo ""

./craft
