# Use hexagonal architecture for Spring Boot backend

We chose hexagonal (ports and adapters) architecture over standard layered architecture for the Spring Boot backend. This keeps the domain layer fully isolated from infrastructure concerns (jOOQ, Spring Security), at the cost of more initial boilerplate for ports and adapter implementations.

## Considered Options
- **Layered architecture**: Standard Spring Boot structure (controller → service → repository), faster to set up but couples domain logic to infrastructure imports.
- **Hexagonal architecture**: Domain layer defines ports (interfaces), infrastructure layer implements adapters. Enforces strict separation of concerns, more upfront structure.
