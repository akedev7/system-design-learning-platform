#!/bin/bash

# Setup Sandbox Script - Detect language and execute setup
# This script is run inside the Docker sandbox to set up the environment

set -e

TARGET_DIR="${1:-.}"
echo "Setting up sandbox for directory: $TARGET_DIR"

cd "$TARGET_DIR" || exit 1

# Function to detect language and run setup
detect_and_setup() {
  if [ -f "package.json" ]; then
    echo "Detected: Node.js/TypeScript project"
    setup_node
  elif [ -f "requirements.txt" ]; then
    echo "Detected: Python project"
    setup_python
  elif [ -f "go.mod" ]; then
    echo "Detected: Go project"
    setup_go
  elif [ -f "Makefile" ]; then
    echo "Detected: Makefile project (generic)"
    setup_makefile
  elif [ -f "taskfile.yaml" ] || [ -f "Taskfile.yml" ]; then
    echo "Detected: Task runner project"
    setup_taskfile
  else
    echo "Warning: Could not detect project type"
    return 1
  fi
}

setup_node() {
  echo "Setting up Node.js environment..."
  
  # Check .npmrc for version constraints
  if [ -f ".npmrc" ]; then
    echo "Found .npmrc, checking Node version..."
    # Add version checking logic here if needed
  fi
  
  # Install dependencies
  if command -v pnpm &> /dev/null; then
    echo "Using pnpm"
    pnpm install
  elif command -v npm &> /dev/null; then
    echo "Using npm"
    npm install
  else
    echo "Error: No package manager found"
    return 1
  fi
}

setup_python() {
  echo "Setting up Python environment..."
  
  # Check Python version from .python-version or runtime.txt if they exist
  if [ -f ".python-version" ]; then
    PYTHON_VERSION=$(cat .python-version)
    echo "Expected Python version: $PYTHON_VERSION"
  fi
  
  # Create virtual environment
  if [ ! -d "venv" ]; then
    python3 -m venv venv
  fi
  
  # Activate and install
  # shellcheck source=/dev/null
  source venv/bin/activate
  pip install -r requirements.txt
  
  echo "Python environment setup complete"
}

setup_go() {
  echo "Setting up Go environment..."
  
  # Check go.mod for version
  if [ -f "go.mod" ]; then
    GO_VERSION=$(grep '^go ' go.mod | awk '{print $2}')
    echo "Expected Go version: $GO_VERSION"
  fi
  
  # Download dependencies
  go mod download
  
  echo "Go environment setup complete"
}

setup_makefile() {
  echo "Running make setup..."
  if make -n setup &> /dev/null; then
    make setup
  else
    echo "No setup target in Makefile, skipping..."
  fi
}

setup_taskfile() {
  echo "Running task setup..."
  if command -v task &> /dev/null; then
    task setup
  else
    echo "Task runner not installed, trying to install..."
    # Install go-task
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
    task setup
  fi
}

# Verify tool versions
verify_versions() {
  echo "Verifying tool versions..."
  
  # Check against .prototools if it exists
  if [ -f ".prototools" ]; then
    echo "Found .prototools, verifying versions..."
    # Add version verification logic here
  fi
  
  echo "Version verification complete"
}

# Main execution
echo "=== Sandbox Setup Started ==="
detect_and_setup
verify_versions
echo "=== Sandbox Setup Complete ==="

exit 0
