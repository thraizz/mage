#!/bin/bash

# Stop Mage Web Client servers

echo "🛑 Stopping Mage Web Client servers..."

# Kill Go server
pkill -f 'web-demo' && echo "   ✅ Go server stopped" || echo "   ℹ️  Go server not running"

# Kill Svelte server
pkill -f 'vite.*5174' && echo "   ✅ Svelte server stopped" || echo "   ℹ️  Svelte server not running"

echo ""
echo "✅ All servers stopped"
