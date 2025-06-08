#!/bin/bash

echo "Building Mark Master Sheet Consolidator..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

# Download dependencies
echo "Downloading dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "Error: Failed to download dependencies"
    exit 1
fi

# Run tests
echo "Running tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Warning: Some tests failed"
fi

# Build the applications
echo "Building CLI application..."
go build -o mark-master-sheet cmd/main.go
if [ $? -ne 0 ]; then
    echo "Error: CLI build failed"
    exit 1
fi

echo "Building GUI application..."
go build -o mark-master-sheet-gui cmd/gui/main.go
if [ $? -ne 0 ]; then
    echo "Error: GUI build failed"
    exit 1
fi

echo "Build completed successfully!"
echo "CLI Executable: mark-master-sheet"
echo "GUI Executable: mark-master-sheet-gui"
echo ""
echo "Usage:"
echo "  ./mark-master-sheet                    # CLI: Run with default config"
echo "  ./mark-master-sheet -dry-run           # CLI: Test run without changes"
echo "  ./mark-master-sheet -stats             # CLI: Show statistics"
echo "  ./mark-master-sheet -config custom.toml # CLI: Use custom config"
echo "  ./mark-master-sheet-gui                # GUI: Launch graphical interface"
echo ""
