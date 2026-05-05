#!/bin/bash
set -e

# Node/TypeScript setup
if [ -f "package.json" ]; then
  echo "Installing Node dependencies..."
  npm install
fi

# Java/Spring Boot setup
if [ -f "pom.xml" ]; then
  echo "Installing Maven dependencies..."
  mvn dependency:resolve -q
fi
