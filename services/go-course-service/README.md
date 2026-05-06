# Go Course Service

Migrated from Spring Boot service to Go using Echo framework.

## Prerequisites

- Go 1.22+ (install via `brew install go` on macOS)
- PostgreSQL 16+
- Make

## Setup

1. Install Go dependencies:
   ```bash
   make setup
   ```

2. Configure environment variables (copy `.env.example` to `.env`):
   ```bash
   cp .env.example .env
   # Edit .env with your database and Auth0 credentials
   ```

3. Run database migrations:
   ```bash
   # Migrations run automatically on server start
   make run
   ```

## Running the Service

```bash
make run
```

Server starts on port 8080 by default.

## API Endpoints

| Method | Path | Description | Auth Required |
|--------|------|-------------|---------------|
| GET | `/actuator/health` | Health check | No |
| GET | `/actuator/info` | Service info | No |
| GET | `/api/v1/courses` | List courses (stub) | Yes |
| GET | `/api/v1/courses/{id}` | Get course by ID (stub) | Yes |
| GET | `/api/v1/courses/{courseId}/modules` | List modules by course | Yes |
| GET | `/api/v1/modules/{id}` | Get module by ID | Yes |
| GET | `/api/v1/modules/{moduleId}/lessons` | List lessons by module | Yes |
| GET | `/api/v1/lessons/{id}` | Get lesson by ID | Yes |

## Testing

```bash
make test
```

## Linting

```bash
make lint
```

## Building

```bash
make build
```

## Docker

Build and run with Docker:
```bash
docker build -t go-course-service .
docker run -p 8080:8080 --env-file .env go-course-service
```
