# Use jOOQ for PostgreSQL persistence

We chose jOOQ (type-safe SQL builder) with Flyway migrations over JPA/Hibernate for Spring Boot's PostgreSQL access. jOOQ gives explicit query control, native JSONB support via custom converters, and compile-time type safety from schema-generated code, at the cost of more boilerplate than JPA for simple CRUD.

## Considered Options
- **JPA/Hibernate**: Standard Spring Boot ORM, faster initial CRUD development, but less control over JSONB column mapping and SQL query behavior.
- **jOOQ**: Type-safe SQL builder, requires Flyway for schema versioning and manual query writing, but avoids JPA's lazy loading pitfalls and magic behavior.
