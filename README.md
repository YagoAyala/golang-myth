# Technical Test - Mytheresa

This repository contains a Go application for managing products, categories, and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.

   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
   - `api/`: Common API response utilities.
   - `catalog/`: Product catalog handlers and logic.
   - `categories/`: Category management handlers.
   - `database/`: Database connection management.

3. **sql/**: Contains a very simple database migration scripts setup.
4. **models/**: Contains the data models and repositories used in the application.
5. `.env`: Environment variables file for configuration (see `env.example` for reference).

## API Documentation

**Swagger/OpenAPI Specification**: [swagger.yaml](swagger.yaml)

**Interactive API Explorer**: Open `docs/swagger-ui.html` in your browser for full interactive API documentation.

- **GET /catalog** - List products with pagination and filtering
  - Query params: `offset`, `limit`, `category`, `priceLessThan`
- **GET /catalog/{code}** - Get product details with variants
- **GET /categories** - List all categories
- **POST /categories** - Create a new category

## Architecture & Quality

This project follows Go best practices and clean architecture principles:

- ✅ **Repository Pattern**: Interfaces defined where consumed, not where implemented
- ✅ **Context Propagation**: All repository methods use `context.Context` for cancellation and timeouts
- ✅ **Structured Logging**: JSON-formatted structured logs using Go's `slog` package
- ✅ **Error Handling**: Centralized middleware for recovery and error management
- ✅ **Export Comments**: All exported types and functions have documentation comments
- ✅ **Unit Tests**: Comprehensive test coverage with mock repositories
- ✅ **CI/CD**: GitHub Actions workflow for automated testing and linting

## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Copy `env.example` to `.env` and adjust settings if needed.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run the tests.
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.

## Database Schema

The application uses the following main entities:

- **Categories**: Product categories (Clothing, Shoes, Accessories)
- **Products**: Products with code, price, and category reference
- **Product Variants**: Variants of products with optional specific pricing

## Testing

The application includes comprehensive unit tests for:
- API response utilities
- Catalog handlers (listing, filtering, pagination, product details)
- Category handlers (listing, creation)

Run tests with: `make test`

### Continuous Integration

The project uses GitHub Actions for automated testing on every push and pull request:
- Runs all tests with race detection
- Performs linting with golangci-lint
- Builds both server and seed commands
- Runs against PostgreSQL in CI environment

## Code Quality

This implementation follows:
- [Effective Go](https://go.dev/doc/effective_go) guidelines
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) conventions
- Clean Architecture principles for maintainability and scalability
