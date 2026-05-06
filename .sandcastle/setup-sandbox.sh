#!/bin/bash
set -e

# Node/TypeScript setup (only if package.json exists and has dependencies)
if [ -f "package.json" ] && grep -q '"dependencies"\|"devDependencies"' package.json; then
  echo "Installing Node dependencies..."
  npm install --prefer-offline
fi

# Java/Spring Boot setup (only if pom.xml exists)
if [ -f "pom.xml" ]; then
  echo "Installing Maven dependencies..."
  mvn dependency:resolve -q -DskipTests
fi
