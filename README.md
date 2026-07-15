# Sharing Vision Backend API

## 1. Project Overview
Backend application for the Sharing Vision fullstack engineer technical test, developed by Stefen Sutandi.

## 2. Features
- RESTful API with standard HTTP methods
- Strict validation (minimum characters, valid statuses)
- Automated tests covering validation, service logic, handlers, pagination, and status workflows.
- Automated database migrations using `golang-migrate`
- CI/CD pipeline using GitHub Actions validating Go tests and Docker E2E with Newman
- Configurable environment via Docker Compose

## 3. Technology Stack
- Go 1.22
- Gin Web Framework
- GORM + MySQL
- Go Playground Validator
- `golang-migrate` for migrations
- Docker & Docker Compose
- Newman (Postman E2E tests)

## 4. Project Structure
- `cmd/api`: Main application entry point
- `internal/handler`: HTTP request handlers and routing
- `internal/service`: Core business logic
- `internal/repository`: Database operations and data access
- `internal/dto`: Data Transfer Objects for requests and responses
- `internal/model`: Database entities and models
- `internal/validator`: Custom validation logic and error mapping
- `migrations`: SQL migration files
- `postman`: API collection and environment files

## 5. Prerequisites
- Docker and Docker Compose
- Go 1.22 (for native setup)
- Newman (for running Postman tests natively)

## 6. Quick Start with Docker
Use this primary quick-start sequence to run the API and its dependencies via Docker Compose:

```bash
git clone https://github.com/StefenSutandi/sharing-vision-backend.git
cd sharing-vision-backend
docker compose up --build -d
docker compose ps
curl http://localhost:8080/health
```

**Startup Order:**
1. **MySQL**: Database container starts first.
2. **Migrate**: Migration container waits for MySQL to be healthy, then runs SQL migrations from the `migrations/` folder.
3. **API**: The main application waits for the migration container to complete successfully before starting.

**Inspect Logs:**
```bash
docker compose logs -f api
docker compose logs migrate
docker compose logs mysql
```

**Stop the Application:**
```bash
docker compose down
```

## 7. Native Go Setup
If you prefer running natively without Docker for development:
```bash
go mod verify
go run ./cmd/api
```
*(Requires a running MySQL database instance and manual execution of migrations)*

## 8. Environment Variables
The application uses the following environment variables (configured via Docker Compose):
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `PORT`: API port (default 8080)

## 9. Database Migration
Migrations are handled explicitly with `golang-migrate` using SQL files in the `migrations/` folder. The `migrate` service in `docker-compose.yml` runs these automatically.

## 10. API Endpoints
| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check endpoint |
| POST | `/article` | Create a new article |
| GET | `/article/:limit/:offset` | Get articles with pagination |
| GET | `/article/:limit/:offset?status=publish` | Get paginated articles by status |
| GET | `/article/:id` | Get article by ID |
| PUT | `/article/:id` | Update an existing article |
| DELETE | `/article/:id` | Soft delete an article (changes status to `thrash`) |

*Note: DELETE performs a soft-trash operation by changing status to `thrash`.*

## 11. Validation Rules
- **title**: required, 20–200 characters
- **content**: required, minimum 200 characters
- **category**: required, 3–100 characters
- **status**: `publish`, `draft`, or `thrash`

## 12. Postman Collection
The repository includes a Postman collection and environment in the `postman/` directory for testing the endpoints.

## 13. Testing
**Native Go Tests:**
```bash
go mod verify
go vet ./...
go test -v ./...
go build ./cmd/api
```

**Newman E2E Tests:**
Run the Postman collection against a live API instance:
```bash
newman run postman/Sharing-Vision-Article-API.postman_collection.json \
  -e postman/Sharing-Vision-Local.postman_environment.json
```

## 14. CI Validation
Docker-based API and Postman/Newman integration tests pass in GitHub Actions.

## 15. Database Reset
To completely reset the database and volume:
```bash
docker compose down -v
docker compose up --build -d
```

## 16. Known Limitations
- The application currently implements soft-trash rather than hard deletion.
- Paginator offsets are zero-indexed and limit/offset validation is constrained by database bounds.
