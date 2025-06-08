#!/bin/bash

echo "Testing Mark Master Sheet Consolidator..."

# Check if executable exists
if [ ! -f "mark-master-sheet" ]; then
    echo "Error: mark-master-sheet not found"
    echo "Please run build.sh first"
    exit 1
fi

# Check if config file exists
if [ ! -f "config.toml" ]; then
    echo "Warning: config.toml not found, copying from sample"
    cp config.sample.toml config.toml
fi

echo ""
echo "=== Test 1: Show Version ==="
./mark-master-sheet -version

echo ""
echo "=== Test 2: Show Statistics ==="
./mark-master-sheet -stats

echo ""
echo "=== Test 3: Dry Run ==="
echo "This will scan files without making changes..."
./mark-master-sheet -dry-run

echo ""
echo "=== Test 4: Validate Configuration ==="
echo "Checking if all required directories and files exist..."

if [ ! -d "StudentFiles" ]; then
    echo "Warning: StudentFiles directory not found"
else
    echo "✓ StudentFiles directory exists"
fi

if [ ! -d "MasterSheet" ]; then
    echo "Warning: MasterSheet directory not found"
else
    echo "✓ MasterSheet directory exists"
fi

if [ ! -d "output" ]; then
    echo "Creating output directory..."
    mkdir -p output
fi

if [ ! -d "logs" ]; then
    echo "Creating logs directory..."
    mkdir -p logs
fi

if [ ! -d "backups" ]; then
    echo "Creating backups directory..."
    mkdir -p backups
fi

echo ""
echo "Testing completed!"
echo ""
echo "To run the actual processing:"
echo "  ./mark-master-sheet"
echo ""
echo "To run with custom config:"
echo "  ./mark-master-sheet -config custom.toml"
echo ""
