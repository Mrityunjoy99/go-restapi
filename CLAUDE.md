# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
- `make build` - Build the application binary with static linking
- `make run` - Run the application on port 8080
- `make run-with-logs` - Run with logging to dated log files in logs/ directory
- `make watch` - Run in development mode with hot reload using air

### Testing
- `make test-run` - Run all tests with race detection
- `make test-cover` - Run tests with coverage report (generates coverage.out)
- `make html-coverage-report` - Generate HTML coverage report (coverage.html)

### Database Operations
- `make migrate` - Run all pending database migrations
- `make rollback-last` - Rollback the last migration
- `make seed-db` - Seed database with initial data
- `make generate-migration name=<migration_name>` - Generate new migration file
- `make db-reset` - Reset database schema (drops and recreates public schema)

### Development Environment
- `make dev-setup-up` - Start dev infrastructure, run migrations, and seed database
- `make dev-infra-up` - Start only the infrastructure (PostgreSQL, Prometheus via Docker)
- `make dev-setup-destroy` - Tear down development infrastructure
- `make dev-setup-reset` - Reset entire development environment

### Code Quality
- `make lint` - Run golangci-lint with 5m timeout
- `make generate` - Run go generate for all packages

### Utilities
- `make generate-admin-token` - Generate admin authentication token

## Architecture Overview

This is a Go REST API following Clean Architecture principles with clear separation of concerns:

### Layer Structure
- **cmd/**: Application entry point and command parsing
- **src/application/**: Application services (user, admin, healthcheck)
- **src/domain/**: Business logic and domain services (JWT, entities)
- **src/infrastructure/**: External dependencies (database, cache)
- **src/deployment/**: Server setup, routing, middleware
- **src/repository/**: Data access layer
- **src/common/**: Shared utilities, config, constants

### Key Components
- **Authentication**: JWT-based with role-based access (User, Admin, Manager)
- **Database**: PostgreSQL with GORM ORM and migration system
- **Web Framework**: Gin for HTTP routing and middleware
- **Configuration**: Viper for environment-based config management
- **Logging**: Custom logger with middleware for request logging

### Service Registry Pattern
The application uses a service registry (`src/domain/service/registry.go`) to manage domain services like JWT authentication, ensuring proper dependency injection and configuration validation.

### Command System
The main application supports multiple commands:
- `run` - Start the web server
- `migrate` - Database migration operations
- Command-line options are parsed via `cmd/cmdopts/`

### Testing Strategy
- Unit tests for services and domain logic
- Mock generation for interfaces (mocks/ directory)
- Coverage reporting with HTML output
- Race condition detection enabled

### Development Infrastructure
Docker Compose setup provides:
- PostgreSQL database
- Prometheus monitoring
- Development-specific configurations in `infra/dev/`