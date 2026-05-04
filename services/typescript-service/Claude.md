# TypeScript Service

## Overview
A TypeScript utility service providing string manipulation and mathematical operations. This service follows TDD principles and is designed to be used as a library.

## Service Details
- **Language**: TypeScript 5.0+ / Node.js 20+
- **Package Manager**: pnpm
- **Testing Framework**: Jest
- **Build Tool**: TypeScript compiler (tsc)
- **Task Runner**: Taskfile.yaml (go-task)

## Available Tasks
Run tasks using `task` command (requires go-task installed):

- `task setup` - Install dependencies
- `task test` - Run tests
- `task test:watch` - Run tests in watch mode
- `task lint` - Run linter
- `task build` - Compile TypeScript
- `task clean` - Remove build artifacts

## Project Structure
```
typescript-service/
├── Claude.md              # This file
├── taskfile.yaml          # Task definitions
├── package.json           # Dependencies
├── tsconfig.json          # TypeScript config
├── src/
│   ├── index.ts           # Main export
│   ├── string-utils.ts    # String utility functions
│   └── math-utils.ts      # Math utility functions
└── tests/
    ├── string-utils.test.ts
    └── math-utils.test.ts
```

## Features

### String Utilities (`string-utils.ts`)
- `capitalize(str: string): string` - Capitalizes first letter
- `reverse(str: string): string` - Reverses a string
- `truncate(str: string, maxLength: number): string` - Truncates string with ellipsis

### Math Utilities (`math-utils.ts`)
- `add(a: number, b: number): number` - Adds two numbers
- `multiply(a: number, b: number): number` - Multiplies two numbers
- `factorial(n: number): number` - Calculates factorial

## TDD Workflow
This service was built using Test-Driven Development:
1. Write failing test (Red)
2. Implement minimal code to pass (Green)
3. Refactor and lint (Refactor)

## Usage Example
```typescript
import { capitalize, add } from 'typescript-service';

console.log(capitalize('hello')); // 'Hello'
console.log(add(2, 3)); // 5
```

## Versioning
Check `package.json` for current version. Follow semantic versioning.
