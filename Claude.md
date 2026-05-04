# Polyglot Monorepo - Service Map

## Overview
This is a polyglot monorepo containing multiple microservices implemented in different languages (TypeScript, Python, Go). Each service resides in its own subfolder under `services/` and follows a strict TDD (Test-Driven Development) workflow with Docker AI Sandbox integration.

## Repository Structure

```
docker-sandbox-test/
├── Claude.md                    # This file - Root service map
├── .github/
│   └── workflows/               # CI/CD and agent workflows
├── services/
│   ├── typescript-service/      # TypeScript/Node.js service
│   ├── python-service/          # Python service
│   └── go-service/              # Go service
└── scripts/
    ├── recon.sh                 # File tree generation
    └── setup-sandbox.sh         # Sandbox setup helper
```

## Services

### TypeScript Service (`services/typescript-service/`)
- **Language**: TypeScript / Node.js
- **Runtime**: Node.js 20+
- **Package Manager**: pnpm
- **Testing**: Jest
- **Build Tool**: TypeScript compiler
- **Task Runner**: Taskfile.yaml (using go-task)
- **Description**: Example utility service with string manipulation and math operations
- **Port**: N/A (library service)
- **Dependencies**: None (pure implementation)

### Python Service (`services/python-service/`)
- **Language**: Python 3.11+
- **Package Manager**: pip
- **Testing**: pytest
- **Build Tool**: setuptools
- **Task Runner**: Makefile
- **Description**: Example data processing service with validation utilities
- **Port**: N/A (library service)
- **Dependencies**: None (pure implementation)

### Go Service (`services/go-service/`)
- **Language**: Go 1.21+
- **Testing**: Go built-in testing
- **Build Tool**: go build
- **Task Runner**: Makefile
- **Description**: Example CLI utility service with file operations
- **Port**: N/A (CLI tool)
- **Dependencies**: None (standard library)

## Agent Workflow

When the `agent-take` label is applied to an issue or PR, the AI agent:

1. **Discovers Context**: Reads root `Claude.md` and subfolder-specific `Claude.md`
2. **Runs Recon**: Executes `scripts/recon.sh` to generate fresh file tree
3. **Sets Up Sandbox**: Runs `scripts/setup-sandbox.sh` for language-specific setup
4. **Executes TDD**: Follows Red-Green-Refactor cycle (max 10 iterations)
5. **Rebases**: Pulls main and rebases feature branch
6. **Creates PR**: Pushes as Draft PR, marks Ready for Review when complete

## Security & Guardrails

- **Sandbox Jail**: Agent has Read-Only access to root, Read-Write to assigned subfolder
- **Network**: Disabled after initial dependency download
- **Resources**: 2GB RAM, 1.0 CPU core limit
- **Secrets**: Injected via environment variables, never logged

## Contributing

Each subfolder must contain:
- `Claude.md` or `README.md` with service-specific documentation
- `taskfile.yaml` or `Makefile` with required targets (setup, test, lint, build)
- Complete test suite following TDD principles
- Source code with proper documentation

## File Tree

*Run `scripts/recon.sh` to generate current file tree*
