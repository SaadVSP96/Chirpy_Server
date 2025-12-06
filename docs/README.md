# Chirpy API

A lightweight social media API built with Go for creating and managing chirps (short messages). Part of the [Boot.dev](https://boot.dev) backend curriculum.

## Motivation

This API was built as part of the Boot.dev backend development curriculum to demonstrate modern Go web development practices. It showcases:

-   Clean architecture with proper separation of concerns
-   JWT-based authentication with refresh token rotation
-   Database operations using SQLC for type-safe SQL
-   RESTful API design with proper HTTP status codes
-   Webhook integration for external services
-   Comprehensive error handling and logging

### Goal

The goal of Chirpy API is to provide a simple, well-documented social media backend that demonstrates production-ready Go patterns while remaining easy to understand and extend. It serves as both a learning tool and a foundation for more complex applications.

## Installation

### Prerequisites

-   Go 1.21+
-   PostgreSQL 14+
-   Git

### Setup

1. Clone the repository:

```bash
git clone https://github.com/SaadVSP96/Chirpy_Server.git
cd Chirpy_Server
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file:

```env
DB_URL=postgres://username:password@localhost/chirpy?sslmode=disable
JWT_SECRET=your-secret-key-here-minimum-32-characters
PLATFORM=local
```

4. Run database migrations:

```bash
# Use your PostgreSQL client or migration tool
psql $DB_URL -f sql/schema/*.sql
```

5. Generate SQLC code:

```bash
sqlc generate
```

6. Run the server:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Quick Start

### Create a User

```go
// POST /api/users
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

**Response:**

```json
{
    "id": "uuid-here",
    "email": "user@example.com",
    "created_at": "2025-12-06T10:00:00Z",
    "updated_at": "2025-12-06T10:00:00Z",
    "is_chirpy_red": false
}
```

### Login

```go
// POST /api/login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

**Response:**

```json
{
    "id": "uuid-here",
    "email": "user@example.com",
    "token": "jwt-token-here",
    "refresh_token": "refresh-token-here",
    "created_at": "2025-12-06T10:00:00Z",
    "updated_at": "2025-12-06T10:00:00Z",
    "is_chirpy_red": false
}
```

### Create a Chirp

```go
// POST /api/chirps
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer {jwt_token}" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello, Chirpy!"}'
```

**Response:**

```json
{
    "id": "chirp-uuid-here",
    "body": "Hello, Chirpy!",
    "user_id": "user-uuid-here",
    "created_at": "2025-12-06T10:00:00Z",
    "updated_at": "2025-12-06T10:00:00Z"
}
```

### Get Chirps

```go
// GET /api/chirps - All chirps
curl http://localhost:8080/api/chirps

// GET /api/chirps?author_id=uuid - Filtered by author
curl http://localhost:8080/api/chirps?author_id=user-uuid-here
```

## API Reference

### User Endpoints

-   `POST /api/users` - Create a new user
-   `PUT /api/users` - Update current user (requires auth)
-   `POST /api/login` - Login and get tokens

### Chirp Endpoints

-   `POST /api/chirps` - Create a new chirp (requires auth)
-   `GET /api/chirps` - List all chirps (optionally filtered by author_id)
-   `GET /api/chirps/{id}` - Get a specific chirp
-   `DELETE /api/chirps/{id}` - Delete your own chirp (requires auth)

### Authentication Endpoints

-   `POST /api/refresh` - Refresh access token
-   `POST /api/revoke` - Revoke refresh token

### Admin Endpoints

-   `POST /admin/reset` - Reset database (development only)
-   `GET /admin/metrics` - View server metrics

### Webhook Endpoints

-   `POST /api/polka/webhooks` - Process external webhooks

## Configuration

### Environment Variables

-   `DB_URL` - PostgreSQL connection string (required)
-   `JWT_SECRET` - Secret for JWT signing (required, min 32 chars)
-   `PLATFORM` - Platform identifier (optional)

### Default Settings

-   JWT tokens expire after 1 hour
-   Refresh tokens expire after 60 days
-   Chirps limited to 140 characters
-   Automatic profanity filtering enabled
-   CORS enabled for all origins (development)

## Development

### Project Structure

```
.
├── main.go              # Application entry point
├── handler_*.go         # HTTP request handlers
├── internal/
│   ├── auth/           # Authentication & JWT logic
│   └── database/       # SQLC-generated database code
├── sql/
│   ├── queries/        # SQL queries for SQLC
│   └── schema/         # Database migrations
├── cmd/                # Utility scripts
├── docs/               # Documentation
└── assets/             # Static files
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
bootdev run {test-id}
```

### Database Utilities

```bash
# Clean database (removes all users and chirps)
go run cmd/db-clean.go

# Reset database via API
curl -X POST http://localhost:8080/admin/reset
```

## Features

-   **Authentication**: Secure JWT-based auth with refresh token rotation
-   **Authorization**: Users can only delete their own chirps
-   **Filtering**: Filter chirps by author ID
-   **Webhooks**: Integration with external services for user upgrades
-   **Profanity Filter**: Automatic content moderation
