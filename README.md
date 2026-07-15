# Sharing Vision Backend API

Backend application for the Sharing Vision fullstack engineer technical test, developed by Stefen Sutandi.

## Features
- RESTful API with standard HTTP methods
- Strict validation (minimum characters, valid statuses)
- Full test coverage for validations and business logic
- Automated database migrations using `golang-migrate`
- CI/CD pipeline using GitHub Actions validating Go tests and Docker E2E with Newman
- Configurable environment via Docker Compose

## Technology Stack
- Go 1.22
- Gin Web Framework
- GORM + MySQL
- Go Playground Validator
- `golang-migrate` for migrations
- Docker & Docker Compose
- Newman (Postman E2E tests)

## Run Locally
1. Run `docker compose up -d`
2. Wait for `article_migrate` container to finish creating tables.
3. API runs on `http://localhost:8080`

## Validation Results
- Go Unit Tests: PASSING in GitHub Actions
- Postman Newman Tests: PASSING in GitHub Actions
- Frontend Integration: PASSING
