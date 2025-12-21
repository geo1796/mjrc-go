# MJRC (My Jump Rope Coach) project guidelines

This document provides essential guidance for new developers joining the MJRC's backend project.

## Project Overview

MJRC is a mobile messaging app aiming to provide a nice dictionary of jump rope (freestyle) skills. This project is
MJRC's backend.

## Tech Stack

- Go
- Chi web framework
- PostgresSQL database
- Sqlc with pgx

## Architecture

## Structure

- Clean separation of concerns
- Feature-based organization placed in `/api`
- Custom libraries such as dependencies, utilities, and models are placed in `/core`
- Tiny Chi DI layer (see `/core/chix/**.go`), including error handler configured
  to return only status code without any details
- sqlc.yml file at the root of the project
- Database schema and sql queries used by `sqlc generate` are placed in `/core/postgres/migrations` &&
  `/core/postgres/queries`
- Database connection is configured in `/core/postgres/database.go`

## Feature Components

Each feature typically includes:

- one of the following: `middleware.go`, `route.go`, `group.go`
- `dto.go`: Data transfer objects (optional)
- `handler.go`: HTTP handlers/controllers
- `handler_test.go`: HTTP handlers/controllers unit tests. (optional)
- `service.go`: Business logic implementation (optional)
- `service_test.go`: Business logic unit tests. (optional)
- `repository.go`: Data access via sqlc queries (optional)
- `repository_test.go`: Data access unit tests. (optional)
- `integration_test.go`: The most important tests. For those, we use testcontainers to run a real database and mock only
  external systems. Mandatory tests

## Best Practices

- Use interfaces for dependency injection
- Maintain at least 70% test coverage
- Prefer integration tests over unit tests (JUST DON'T USE EXTERNAL API IN TESTS)
- Integration tests should use real implementations when possible
- Integration tests should not mock repositories and use testcontainers database
- Generate mocks with mockgen. If needed, you can regenerate mocks using mockgen in the command line (a simple
  mockgen.sh script is provided at the root of the project)
- Only mock external systems (e.g., mailer, s3, fcm, yoti) or repositories for unit
- MAX 100 lines per function

## JUNIE

Always look for the word 'JUNIE' in comments for guidance.

DO NOT EDIT .sql files unless I explicitly asked you to. 