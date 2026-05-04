# Python Service

## Overview
A Python data processing service providing validation utilities. This service follows TDD principles and is designed to be used as a library.

## Service Details
- **Language**: Python 3.11+
- **Package Manager**: pip
- **Testing Framework**: pytest
- **Build Tool**: setuptools
- **Task Runner**: Makefile

## Available Targets
Run targets using `make` command:

- `make setup` - Create virtual environment and install dependencies
- `make install` - Install dependencies
- `make test` - Run tests
- `make test-verbose` - Run tests with verbose output
- `make lint` - Run linter (flake8)
- `make format` - Format code (black)
- `make clean` - Remove build artifacts and cache

## Project Structure
```
python-service/
├── Claude.md              # This file
├── Makefile               # Task definitions
├── requirements.txt       # Dependencies
├── setup.py              # Package setup
├── src/
│   ├── __init__.py
│   └── validators.py     # Validation utility functions
└── tests/
    ├── __init__.py
    └── test_validators.py
```

## Features

### Validators (`validators.py`)
- `is_email(string: str) -> bool` - Validates email format
- `is_phone(string: str) -> bool` - Validates phone number format (US)
- `is_strong_password(string: str) -> bool` - Validates password strength

## TDD Workflow
This service was built using Test-Driven Development:
1. Write failing test (Red)
2. Implement minimal code to pass (Green)
3. Refactor and lint (Refactor)

## Usage Example
```python
from src.validators import is_email, is_phone, is_strong_password

print(is_email("test@example.com"))  # True
print(is_phone("(123) 456-7890"))    # True
print(is_strong_password("Abc123!@#"))  # True
```

## Versioning
Check `setup.py` for current version. Follow semantic versioning.
