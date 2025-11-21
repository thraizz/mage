#!/bin/bash

# Start Mage Web Client
# This script starts both the Go WebSocket server and Svelte dev server

set -e

echo "🚀 Starting Mage Web Client..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Check if Bun is installed
if ! command -v bun &> /dev/null; then
    echo "❌ Bun is not installed. Please install Bun first."
    exit 1
fi

# Kill any existing servers
echo "🧹 Cleaning up existing servers..."
pkill -f 'web-demo' 2>/dev/null || true
pkill -f 'vite.*5174' 2>/dev/null || true
sleep 1

# Start Go WebSocket server
echo "🔌 Starting Go WebSocket server..."
cd ../mage-server-go
nohup go run cmd/web-demo/main.go > /tmp/mage-ws-server.log 2>&1 &
GO_PID=$!
echo "   ✅ Go server started (PID: $GO_PID)"
echo "   📡 WebSocket: ws://localhost:8080/ws"

# Wait for Go server to start
sleep 2

# Start Svelte dev server
echo ""
echo "🎨 Starting Svelte dev server..."
cd ../mage-client-web
nohup bun run dev --host 0.0.0.0 --port 5174 > /tmp/mage-svelte-server.log 2>&1 &
SVELTE_PID=$!
echo "   ✅ Svelte server started (PID: $SVELTE_PID)"

# Wait for Svelte server to start
sleep 3

echo ""
echo "✅ Both servers are running!"
echo ""
echo "🎮 Mage Web Client"
echo "=================="
echo ""
echo "Go WebSocket Server:  ws://localhost:8080/ws  (PID: $GO_PID)"
echo "Svelte Dev Server:    http://localhost:5174/  (PID: $SVELTE_PID)"
echo ""
echo "📖 Open your browser to: http://localhost:5174/"
echo ""
echo "🎯 Demo features:"
echo "   - 4 creatures on battlefield"
echo "   - 2 players (Alice vs Bob)"
echo "   - Declare attackers"
echo "   - Pass priority"
echo "   - Real-time WebSocket updates"
echo ""
echo "📋 Logs:"
echo "   Go server:    tail -f /tmp/mage-ws-server.log"
echo "   Svelte server: tail -f /tmp/mage-svelte-server.log"
echo ""
echo "🛑 To stop servers:"
echo "   kill $GO_PID $SVELTE_PID"
echo "   or run: ./stop.sh"
echo ""
