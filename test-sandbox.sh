#!/bin/bash

# Test Docker Sandbox Script
# Tests the sandbox feature with resource constraints

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_NAME="docker-sandbox-test"
CONTAINER_NAME="sandbox-test"

echo "=== Building Docker Sandbox Image ==="
docker build -f Dockerfile.sandbox -t "$IMAGE_NAME" "$SCRIPT_DIR"

echo ""
echo "=== Phase 1: Setup with Network (Writable) ==="
docker run --rm \
    --name "${CONTAINER_NAME}-setup" \
    -v "$SCRIPT_DIR":/workspace:rw \
    "$IMAGE_NAME" \
    bash -c "echo '=== Environment Setup ===' && \
             echo 'Node version:' && node --version && \
             echo 'Python version:' && python3 --version && \
             echo 'Go version:' && go version && \
             echo '' && \
             echo '=== Running setup-sandbox.sh ===' && \
             setup-sandbox.sh /workspace/services/spring-boot-service"

echo ""
echo "=== Phase 2: Agent Work (Read-Only, No Network) ==="
echo "Constraints:"
echo "  - Memory: 2GB"
echo "  - CPU: 1.0 core"
echo "  - Network: Disabled"
echo "  - Root filesystem: Read-only"
echo ""

# Run container with constraints (simulating agent work)
docker run --rm \
    --name "$CONTAINER_NAME" \
    --memory="2g" \
    --memory-swap="2g" \
    --cpus="1.0" \
    --read-only \
    --tmpfs /tmp \
    --tmpfs /run \
    --network none \
    -v "$SCRIPT_DIR":/workspace:ro \
    -v "$SCRIPT_DIR/services/spring-boot-service":/workspace/services/spring-boot-service:rw \
    "$IMAGE_NAME" \
    bash -c "echo '=== Agent Sandbox ===' && \
             echo 'Node version:' && node --version && \
             echo 'Python version:' && python3 --version && \
             echo 'Go version:' && go version && \
             echo '' && \
             echo 'Testing read-only root...' && \
             touch /test-write 2>/dev/null && echo 'ERROR: Write succeeded' || echo 'OK: Root is read-only' && \
             echo '' && \
             echo 'Testing network disabled...' && \
             curl -s https://google.com 2>/dev/null && echo 'ERROR: Network works' || echo 'OK: Network is disabled'"

echo ""
echo "=== Sandbox Test Complete ==="
