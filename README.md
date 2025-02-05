# Memecoin Management Application

## Project Overview

This is a Memecoin API built using Go and Gin as the ORM, supporting PostgreSQL as the database.

## Project Structure

```
/memecoin_homework
│
├── .env                     # Environment variables file
├── .gitignore               # Git ignore file
├── go.mod                   # Go module file
├── go.sum                   # Go dependency file
├── cmd/                     # Command files for the application
│   └── memecoin/            # Main entry point of the application
│       └── main.go          # Main application logic
├── docker-compose.yml       # Docker Compose configuration file
├── Dockerfile               # Dockerfile for building the application
├── internal/                # Contains internal logic and services
│   ├── model/               # Data models
│   │   └── memecoin.go      # Memecoin data model
│   └── service/             # Business logic services
│       └── memecoin_service.go # Memecoin related service logic
├── api/                     # API routes and handlers
│   └── memecoin_api.go      # Memecoin API routes
```

## Running the Application Locally

1. Ensure you have Go and PostgreSQL installed.
2. Clone this project to your local machine:
   ```bash
   git clone <repository-url>
   cd memecoin_homework
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Create a `.env` file in the project root with the following content:

   ```
   DB_USER=user
   DB_PASSWORD=password
   DB_NAME=memecoins
   DB_HOST=localhost

   ```

5. Start the PostgreSQL database:
   ```bash
   docker-compose up -d postgres
   ```
6. Run the application:
   ```bash
   cd ..
   go run cmd/memecoin/main.go
   ```
7. Access the API:
   - Create Memecoin: `POST /memecoins`
   - Get Memecoin: `GET /memecoins/:id`
   - Update Memecoin: `PATCH /memecoins/:id`
   - Delete Memecoin: `DELETE /memecoins/:id`
   - Increase Memecoin Popularity: `POST /memecoins/:id/poke`

## Running the Application in a Docker Container

1. Ensure you have Docker and Docker Compose installed.
2. Clone this project to your local machine:
   ```bash
   git clone <repository-url>
   cd memecoin_homework
   ```
3. Create a `.env` file in the project root with the following content:

   ```
   DB_USER=user
   DB_PASSWORD=password
   DB_NAME=memecoins
   DB_HOST=postgres

   ```

4. Start the application and database:
   ```bash
   docker-compose up --build
   ```
5. Access the API:
   - Create Memecoin: `POST /memecoins`
   - Get Memecoin: `GET /memecoins/:id`
   - Update Memecoin: `PATCH /memecoins/:id`
   - Delete Memecoin: `DELETE /memecoins/:id`
   - Increase Memecoin Popularity: `POST /memecoins/:id/poke`

## Setup and Configuration

- Database configuration is in the `docker-compose.yml` file:
  - Username: `${DB_USER}`
  - Password: `${DB_PASSWORD}`
  - Database Name: `${DB_NAME}`
- Application environment variables can be configured in the `.env` file.
- To clean up orphan containers, you can run:
  ```bash
   docker-compose down --volumes --remove-orphans
  ```
