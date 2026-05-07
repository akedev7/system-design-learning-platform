# Use sqlx for PostgreSQL persistence in Go backend

We chose sqlx (lightweight SQL extension for Go's database/sql) over ORMs like GORM for the Go backend's PostgreSQL access. sqlx provides struct scanning for query results, minimal boilerplate, and full control over raw SQL queries, at the cost of less built-in migration tooling compared to opinionated ORMs.

## Considered Options
- **GORM**: Popular Go ORM with automatic migrations and CRUD helpers, but adds heavy abstraction and less explicit query control.
- **database/sql**: Go standard library SQL interface, no struct scanning support, requires more boilerplate for row mapping.
- **sqlx**: Lightweight extension to database/sql with struct scanning, preserves full control over raw SQL, aligns with the migrated Go service's explicit query needs.
