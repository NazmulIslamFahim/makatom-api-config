#!/bin/bash

# Install Air for hot reloading
echo "Installing Air for hot reloading..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

# Install Air
echo "Installing Air..."
go install github.com/cosmtrek/air@latest

# Check if installation was successful
if command -v air &> /dev/null; then
    echo "✅ Air installed successfully!"
    echo "You can now run './run.sh' for hot reloading during development."
else
    echo "❌ Failed to install Air. Please check your Go installation."
    echo "You can still run the server with 'go run cmd/main.go'"
fi
