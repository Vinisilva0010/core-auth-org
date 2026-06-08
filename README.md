# Enterprise Core Backend

High-performance foundational backend for B2B enterprise applications, built with Go and PostgreSQL.

This project is designed as a modular monolith with strict domain boundaries, type-safe SQL, and minimal framework overhead. Its goal is to provide a solid, reusable core for systems such as ERP, internal admin platforms, business management software, and other organization-centric applications.

## Overview

The backend is structured around bounded contexts instead of generic technical layers. This keeps responsibilities clear, reduces coupling between modules, and makes long-term maintenance easier in growing codebases.

The current focus is the enterprise core layer: authentication, identity, organization management, permissions, and auditability. Business-specific modules should be added later on top of this foundation, not mixed into it from the start.

## Core Principles

- **Domain-first design**  
  Business rules belong to the domain and service layers, not to HTTP handlers or database adapters.

- **Modular monolith**  
  Features are separated into modules with clear boundaries, while remaining in a single deployable application for simpler operations and development flow.

- **Type-safe SQL**  
  Queries are written explicitly and mapped to Go types through `sqlc`, avoiding heavy ORM abstractions and reducing runtime surprises.

- **Minimal overhead**  
  The project favors lean, understandable building blocks over excessive abstraction, reflection-heavy tooling, or unnecessary framework complexity.

- **Operational clarity**  
  Logging, migrations, configuration, and database access are treated as first-class parts of the system, not afterthoughts.

## Tech Stack

- **Language:** Go 1.22+
- **Database:** PostgreSQL 15+
- **Database Driver / Pool:** `pgx/v5`
- **SQL Code Generation:** `sqlc`
- **HTTP Router:** `go-chi/chi/v5`
- **Migrations:** `golang-migrate`
- **Logging:** `log/slog`
- **Authentication:** JWT + password hashing
- **Architecture Style:** Modular Monolith

## Current Modules

| Module | Responsibility | Status |
|--------|----------------|--------|
| `auth` | Authentication, session flow, token lifecycle | In progress |
| `users` | User identity, registration, account lifecycle | In progress |
| `org` | Organizations, units, and tenant-related structure | Planned / WIP |
| `rbac` | Roles and permissions | Planned / WIP |
| `audit` | Security-sensitive event tracking and audit trails | Planned / WIP |

## Project Structure

```text
.
├── cmd/
│   └── api/                      # Application entrypoint
├── internal/
│   ├── platform/                 # Shared infrastructure: config, db, logger, server
│   └── modules/                  # Business modules with strict boundaries
│       ├── auth/                 # Authentication and session lifecycle
│       ├── users/                # User identity and account management
│       ├── org/                  # Organizations and related structure
│       ├── rbac/                 # Roles and permissions
│       └── audit/                # Audit trail and security event records
├── migrations/                   # Raw SQL migration files
└── sqlc.yaml                     # sqlc configuration
```

## Design Goals

This project is being built as a reusable enterprise core, not as a one-off backend. The intent is to keep the foundation stable, small, explicit, and adaptable across different business systems.

Key goals include:

- strong module boundaries
- predictable dependency flow
- explicit SQL and schema control
- low operational complexity
- long-term maintainability
- easy extension for future business modules

## What This Project Is Not

To keep the core clean, this repository does not aim to include business-specific features such as inventory, orders, scheduling, billing, or reporting at this stage. Those concerns belong in separate modules built on top of this base.

It also avoids generic repository patterns, heavy ORMs, and unnecessary abstractions that make code harder to understand and reuse over time.

## Getting Started

### Prerequisites

Before running the project, make sure the local environment includes:

- Go 1.22+
- PostgreSQL 15+
- `sqlc`
- `golang-migrate`
- a PostgreSQL database created for local development

### Environment Variables

At minimum, the application should be configured through environment variables for:

- application environment
- HTTP server port
- PostgreSQL connection string
- JWT secret
- token expiration settings
- log level

A common approach is to use a local `.env` file during development and environment-specific configuration in production.

### Initial Setup Flow

A clean setup flow should follow this order:

1. create the PostgreSQL database
2. configure environment variables
3. run database migrations
4. generate SQL code with `sqlc`
5. start the API server

This order keeps schema, generated code, and runtime configuration aligned from the beginning.

## Development Flow

The recommended workflow for this repository is module-first, not endpoint-first.

A typical development sequence should be:

1. define the module boundary
2. model the domain rules
3. create or update migrations
4. write the SQL queries
5. generate typed code with `sqlc`
6. implement repositories
7. implement services / use cases
8. expose HTTP handlers
9. add tests for critical flows
10. document the module behavior

This reduces rework and helps keep domain logic from leaking into transport or persistence layers.

## Recommended Standards

To keep the codebase clean over time, each module should follow a small and explicit internal structure. The exact naming can vary, but the responsibility split should remain stable:

- `domain` for entities, rules, and business errors
- `service` for use cases and orchestration
- `repository` for PostgreSQL access
- `transport/http` for request and response handling
- `dto` for input and output contracts when needed

The goal is not to create ceremony, but to preserve clarity as the system grows.

## Development Status

This repository is currently in the foundational stage. The architecture, module boundaries, and infrastructure setup are being defined first so that later features can be added without creating structural debt.

## Long-Term Vision

The long-term goal is to turn this repository into a professional reusable backend core for enterprise-grade systems, especially where authentication, organizational structure, access control, and auditability need to be reliable from the beginning.

Future business modules should plug into this base instead of re-implementing common enterprise concerns in every new project.
