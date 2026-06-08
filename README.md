# Enterprise Core Backend

A reusable backend foundation for business systems that need authentication, users, organizations, permissions, and audit logs.

This project is meant to be the starting point for software such as ERP systems, internal admin panels, inventory systems, clinic systems, barbershop systems, restaurant back offices, and other business applications where multiple users need controlled access to the same system.

It is **not** the full business system itself. It is the core layer that solves the repeated enterprise problems first, so new projects can be built faster and with less structural mess.

## What this backend is for

Use this backend when the project needs one or more of these:

- user accounts
- login and session handling
- password reset flow
- organizations / companies
- units, branches, or locations
- roles and permissions
- internal admin access control
- audit trail for important actions
- a reusable base for future business modules

In simple terms, this backend helps when the system is not just a public website, but a real application with users, access rules, and company structure.

## Where this backend can be used

This backend is a good fit for projects such as:

- ERP-like systems
- stock and inventory systems
- clinic management systems
- barbershop management systems
- pizzeria or restaurant back-office systems
- internal dashboards for companies
- SaaS admin backends
- multi-user business tools
- management systems for small and medium businesses

The main idea is simple: if the project needs people to log in, belong to a company, have different permissions, and leave a history of important actions, this backend is useful.

## Where this backend should not be used alone

This backend is **not enough by itself** for:

- landing pages
- blogs
- portfolio websites
- brochure websites
- simple APIs without auth or business access rules
- highly specialized systems that already depend on a very different domain model

It also does not yet include business modules such as orders, products, scheduling, billing, finance, or reports. Those should be built on top of this foundation.

## What problem it solves

In many business projects, the same base work gets rebuilt over and over:

- login
- user management
- company structure
- permission checks
- session flow
- audit logging
- account lifecycle rules

This backend exists to stop rebuilding that same core every time.

Instead of starting each project from zero, the goal is to start with a clean and reliable base, then add only the business-specific modules needed for the project.

## What you gain by using it

Using this backend should help with:

- faster project starts
- cleaner project structure
- better consistency across different systems
- less repeated code
- easier long-term maintenance
- safer access control
- easier scaling of future modules
- fewer architecture mistakes early in the project

It is especially useful for freelancers, agencies, and product builders who create multiple business systems and do not want to redesign the same backend foundation every time.

## What is included

The current goal of this repository is to provide the enterprise core, especially around:

- authentication
- users
- organizations
- roles and permissions
- audit logs
- shared platform infrastructure

That means this repository focuses on the system backbone, not on domain-specific business logic.

## Project structure

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

## Example use cases

### Example 1: Clinic system

A clinic system may need receptionists, admins, and managers with different permissions. It may also need branches, protected patient-related actions, and a history of who changed what. This backend solves the base layer for that kind of structure.

### Example 2: Inventory system

An inventory system may need multiple staff users, branch-level access, controlled admin actions, and an audit trail for sensitive changes. This backend gives the access and organization layer before stock rules are added.

### Example 3: Restaurant back office

A restaurant or pizzeria system may need managers, attendants, and owners with different access levels across one or more units. This backend helps organize users, permissions, and company structure before adding menus, orders, and delivery logic.

## Development approach

This project should stay small in scope and strong in structure.

The idea is:

1. build the enterprise core well
2. avoid unnecessary abstractions
3. keep module boundaries clear
4. add business modules later without breaking the foundation

This helps the backend stay reusable instead of turning into a one-project code dump.

## Tech stack

- **Language:** Go 1.22+
- **Database:** PostgreSQL 15+
- **Database Driver / Pool:** `pgx/v5`
- **SQL Code Generation:** `sqlc`
- **HTTP Router:** `go-chi/chi/v5`
- **Migrations:** `golang-migrate`
- **Logging:** `log/slog`
- **Authentication:** JWT + password hashing

## Current status

This repository is still in the foundation stage. The goal right now is to make the core stable, reusable, and clear before expanding into more business-specific modules.

## Long-term goal

The long-term goal is to keep this repository as a professional backend base that can be reused across many business systems.

Instead of creating authentication, organization, permission, and audit logic from scratch for every new project, future systems should start here and build their own business modules on top.
