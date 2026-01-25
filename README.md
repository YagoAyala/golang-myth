# Go Hiring Challenge

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

## API Endpoints

### Catalog Endpoints

- **GET /catalog** - List all products with pagination and filtering
  - Query Parameters:
    - `offset` (optional): Pagination offset (default: 0)
    - `limit` (optional): Number of items per page (default: 10, min: 1, max: 100)
    - `category` (optional): Filter by category code (e.g., "clothing", "shoes", "accessories")
    - `priceLessThan` (optional): Filter products with price less than specified value
  - Response: `{ "products": [...], "total": <count> }`

- **GET /catalog/{code}** - Get product details by code
  - Returns product with all variants (variants inherit product price if not specified)
  - Response includes category information

### Categories Endpoints

- **GET /categories** - List all categories
  - Response: `{ "categories": [...] }`

- **POST /categories** - Create a new category
  - Request Body: `{ "code": "category-code", "name": "Category Name" }`
  - Response: The created category

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

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)
