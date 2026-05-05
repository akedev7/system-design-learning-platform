#!/bin/bash
# Test script for Docker Compose infrastructure
# Following TDD approach: RED -> GREEN -> REFACTOR

set -e

COMPOSE_FILE="docker-compose.yml"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "=== Docker Compose Infrastructure Tests ==="
echo ""

# Test 1: docker-compose.yml exists
echo -n "Test 1: docker-compose.yml exists... "
if [ -f "$PROJECT_ROOT/$COMPOSE_FILE" ]; then
    echo "PASS"
else
    echo "FAIL - $COMPOSE_FILE not found"
    exit 1
fi

# Test 2: docker-compose.yml is valid YAML and has required structure
echo -n "Test 2: docker-compose.yml syntax and structure... "
if command -v docker &> /dev/null && docker compose version &> /dev/null; then
    if docker compose -f "$PROJECT_ROOT/$COMPOSE_FILE" config &> /dev/null; then
        echo "PASS"
    else
        echo "FAIL - Invalid docker-compose.yml syntax"
        exit 1
    fi
else
    # If docker not available, do basic YAML check
    if python3 -c "import yaml; yaml.safe_load(open('$PROJECT_ROOT/$COMPOSE_FILE'))" 2>/dev/null; then
        echo "PASS (yaml valid, docker not available for full check)"
    else
        echo "FAIL - Invalid YAML syntax"
        exit 1
    fi
fi

# Test 3: Required services are defined (client, server, postgres)
echo -n "Test 3: Required services defined (client, server, postgres)... "
SERVICES=$(docker compose -f "$PROJECT_ROOT/$COMPOSE_FILE" config --services 2>/dev/null || \
           grep -E "^  [a-z]" "$PROJECT_ROOT/$COMPOSE_FILE" | tr -d ' :' | head -20)

if echo "$SERVICES" | grep -q "client" && \
   echo "$SERVICES" | grep -q "server" && \
   echo "$SERVICES" | grep -q "postgres"; then
    echo "PASS"
else
    echo "FAIL - Missing required services. Found: $SERVICES"
    exit 1
fi

echo ""
echo "=== All basic tests passed ==="
echo ""

# Test 4: Volume mounts for hot-reload (client)
echo -n "Test 4: Client has volume mount for hot-reload... "
if awk '/^  client:/{found=1} found && /^    volumes:/{print; exit}' "$PROJECT_ROOT/$COMPOSE_FILE" | grep -q "volumes:"; then
    echo "PASS"
else
    echo "FAIL - Client service missing volume mounts"
    exit 1
fi

# Test 5: Volume mounts for hot-reload (server)
echo -n "Test 5: Server has volume mount for hot-reload... "
if awk '/^  server:/{found=1} found && /^    volumes:/{print; exit}' "$PROJECT_ROOT/$COMPOSE_FILE" | grep -q "volumes:"; then
    echo "PASS"
else
    echo "FAIL - Server service missing volume mounts"
    exit 1
fi

# Test 6: PostgreSQL environment variables configured
echo -n "Test 6: PostgreSQL environment configured... "
if grep -A 10 "postgres:" "$PROJECT_ROOT/$COMPOSE_FILE" | grep -q "POSTGRES_DB\|POSTGRES_USER\|POSTGRES_PASSWORD"; then
    echo "PASS"
else
    echo "FAIL - PostgreSQL environment variables missing"
    exit 1
fi

# Test 7: Server can connect to PostgreSQL
echo -n "Test 7: Server datasource URL points to postgres service... "
if grep -A 20 "server:" "$PROJECT_ROOT/$COMPOSE_FILE" | grep -q "jdbc:postgresql://postgres:"; then
    echo "PASS"
else
    echo "FAIL - Server datasource not configured for postgres service"
    exit 1
fi

# Test 8: Auth0 environment variables
echo -n "Test 8: Auth0 environment variables configured... "
if grep -q "AUTH0_ISSUER_URI\|AUTH0_ISSUER_BASE_URL" "$PROJECT_ROOT/$COMPOSE_FILE"; then
    echo "PASS"
else
    echo "FAIL - Auth0 environment variables missing"
    exit 1
fi

echo ""
echo "=== All infrastructure tests passed ==="
echo ""

# Test 9: .env.example file exists
echo -n "Test 9: .env.example file exists... "
if [ -f "$PROJECT_ROOT/.env.example" ]; then
    echo "PASS"
else
    echo "FAIL - .env.example not found"
    exit 1
fi

# Test 10: README has Docker Compose instructions
echo -n "Test 10: README has Docker Compose instructions... "
if grep -qi "docker.compose" "$PROJECT_ROOT/README.md" 2>/dev/null || \
   grep -qi "docker.compose" "$PROJECT_ROOT/client/README.md" 2>/dev/null; then
    echo "PASS"
else
    echo "FAIL - README missing Docker Compose instructions"
    exit 1
fi

echo ""
echo "=== All tests passed ==="
echo ""

# Test 11: Dockerfiles exist for client and server
echo -n "Test 11: Client Dockerfile.dev exists... "
if [ -f "$PROJECT_ROOT/client/Dockerfile.dev" ]; then
    echo "PASS"
else
    echo "FAIL - client/Dockerfile.dev not found"
    exit 1
fi

echo -n "Test 12: Server Dockerfile.dev exists... "
if [ -f "$PROJECT_ROOT/services/spring-boot-service/Dockerfile.dev" ]; then
    echo "PASS"
else
    echo "FAIL - services/spring-boot-service/Dockerfile.dev not found"
    exit 1
fi

echo ""
echo "=== All tests passed ==="
