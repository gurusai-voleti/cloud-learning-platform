#!/bin/bash

# Default values
VERBOSE=false
RACE=false
COVERAGE=false
WATCH=false
PKG="./..."

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -r|--race)
            RACE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -w|--watch)
            WATCH=true
            shift
            ;;
        -p|--package)
            PKG="$2"
            shift
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [-v|--verbose] [-r|--race] [-c|--coverage] [-w|--watch] [-p|--package PKG]"
            exit 1
            ;;
    esac
done

# Function to run tests
run_tests() {
    # Build the test command
    CMD="go test"
    
    if [ "$VERBOSE" = true ]; then
        CMD="$CMD -v"
    fi
    
    if [ "$RACE" = true ]; then
        CMD="$CMD -race"
    fi
    
    if [ "$COVERAGE" = true ]; then
        CMD="$CMD -coverprofile=coverage.out"
    fi
    
    # Add the package to test
    CMD="$CMD $PKG"
    
    # Run the tests
    echo "Running: $CMD"
    eval $CMD
    
    # If coverage is enabled, generate HTML report
    if [ "$COVERAGE" = true ]; then
        go tool cover -html=coverage.out -o coverage.html
        echo "Coverage report generated: coverage.html"
    fi
}

# Function to watch for changes and run tests
watch_tests() {
    echo "Watching for changes..."
    while true; do
        run_tests
        
        # Watch for file changes (requires fswatch)
        fswatch -1 .
        echo "Change detected, running tests..."
    done
}

# Main execution
if [ "$WATCH" = true ]; then
    # Check if fswatch is installed
    if ! command -v fswatch &> /dev/null; then
        echo "Error: fswatch is required for watch mode."
        echo "Install with: brew install fswatch (macOS) or apt-get install fswatch (Linux)"
        exit 1
    fi
    watch_tests
else
    run_tests
fi 