# System Design Learning Platform

A Next.js + Go web application for interactive system design learning, inspired by Brilliant's hands-on exercise model.

## Prerequisites

- Docker & Docker Compose
- Node.js 20+ (for local development without Docker)
- Go 1.22+ (for local development without Docker)
- pnpm (for local frontend development)

## Quick Start with Docker Compose

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd system-design-learning-platform
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and configure your Auth0 credentials:
   - `AUTH0_ISSUER_URI`: Your Auth0 tenant URL (e.g., `https://your-tenant.auth0.com/`)
   - `AUTH0_CLIENT_ID`: From your Auth0 Application settings
   - `AUTH0_CLIENT_SECRET`: From your Auth0 Application settings
   - `AUTH0_SECRET`: Generate with `openssl rand -hex 32`

3. **Start all services**
    ```bash
    docker compose up
    ```
    This starts:
    - **PostgreSQL** (port 5432) - Database
    - **Go API** (port 8080) - Backend service
    - **Next.js Client** (port 3000) - Frontend with hot-reload

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

## Hot-Reload Development

Both the frontend and backend support hot-reload during development:

- **Frontend (Next.js)**: Edit files in `client/` and changes will be reflected immediately
- **Backend (Go)**: Use `air` for hot-reload, or run `cd services/go-course-service && go run ./cmd/server`

## Services

### Client (Next.js)
- Location: `client/`
- Port: 3000
- Hot-reload: Enabled via volume mount
- Package manager: pnpm

### Server (Go)
- Location: `services/go-course-service/`
- Port: 8080
- Hot-reload: Use `air` or `fresh` for development
- Build tool: Go modules (Makefile)

### Database (PostgreSQL)
- Port: 5432
- Database: courses
- Credentials: postgres/postgres (configurable via environment variables)
- Migrations: golang-migrate (see `services/go-course-service/migrations/`)

## API Response Format

All API responses use a standardized envelope format:

```json
{
  "status": "success" | "error",
  "data": {},
  "message": ""
}
```

Example success response:
```json
{
  "status": "success",
  "data": [
    { "id": 1, "title": "System Design 101" }
  ]
}
```

Example error response:
```json
{
  "status": "error",
  "message": "course not found"
}
```

Use the `api.ts` client library in `client/src/lib/api.ts` for consistent API calls with interceptor handling.

## Environment Variables

See `.env.example` for all available configuration options.

## Auth0 Setup

1. Create an Auth0 account at https://auth0.com/
2. Create a new Application (Single Page Application or Regular Web Application)
3. Configure the following URLs in Auth0:
   - Allowed Callback URLs: `http://localhost:3000/api/auth/callback`
   - Allowed Logout URLs: `http://localhost:3000`
   - Allowed Web Origins: `http://localhost:3000`
4. Copy the Client ID and Client Secret to your `.env` file

## Without Docker (Local Development)

### Frontend
```bash
cd client
pnpm install
pnpm dev
```

### Backend
```bash
cd services/go-course-service
make run
```

### Database
Install PostgreSQL locally and create a database named `courses`.

## Project Structure

```
system-design-learning-platform/
├── client/                          # Next.js frontend
├── services/
│   └── go-course-service/          # Go (Echo + sqlx) backend API
├── docs/adr/                        # Architecture Decision Records
├── docker-compose.yml               # Local development setup
└── README.md
```

## License

[Add your license here]
