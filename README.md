# Sharing Vision Backend API

Backend API for the Sharing Vision fullstack engineer technical test, developed by Stefen Sutandi.

## Features
- Article CRUD with strict validation
- Support for `publish`, `draft`, and `thrash` statuses.
- Soft-deletion (Trash workflow) by updating status to `thrash`.
- Pagination with limit/offset and status filtering.
- Fully Dockerized with MySQL.

## Technology Stack
- Golang 1.22
- Gin HTTP Framework
- GORM (with MySQL driver)
- Go-Playground Validator v10
- MySQL 8.0

## Environment Setup
Copy the example environment file:
```bash
cp .env.example .env
```
Edit `.env` to match your database credentials.

## Local Running Instructions (Native)
If you have Go and MySQL installed locally:
1. Initialize the database: Create a database named `article`.
2. Run migrations (using raw SQL files in `migrations/`).
3. Run the application:
   ```bash
   go mod download
   go run cmd/api/main.go
   ```

## Docker Running Instructions (Recommended)
This is the easiest way to start the backend.
```bash
docker compose up --build
```
This will start both the Go API backend on `http://localhost:8080` and a MySQL 8.0 container. Migrations (table creation) are handled via raw SQL scripts. For explicit migrations, use the `migrations/` files provided.

## API Endpoint Table

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Health check endpoint |
| POST | `/article/` | Create a new article |
| GET | `/article/:limit/:offset` | List articles with pagination |
| GET | `/article/:id` | Get article by ID |
| PUT | `/article/:id` | Update an existing article |
| DELETE | `/article/:id` | Move an article to trash (soft-delete) |

## Validation Rules
- **Title**: required, 20-200 characters
- **Content**: required, min 200 characters
- **Category**: required, 3-100 characters
- **Status**: strictly one of `publish`, `draft`, or `thrash`

## Postman Collection Instructions
Import `postman/Sharing-Vision-Article-API.postman_collection.json` into Postman.
Ensure you set the `{{base_url}}` variable to `http://localhost:8080`.

## Test Instructions
To run backend unit tests:
```bash
go test -v ./...
```

## Known Limitations
- Backend and MySQL deployment to a free tier was deferred as reliable free MySQL hosting is limited. Local Docker compose is the recommended verification method.
