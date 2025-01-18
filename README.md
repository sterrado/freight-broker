# Freight Broker Application

A Go-based freight broker application for managing loads through various TMS (Transportation Management System) providers.

## Prerequisites

- Docker and Docker Compose
- Make (optional, for using Makefile commands)
- Go 1.21 or later (for local development)

## Environment Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd freight-broker
```

2. Copy the example environment file and adjust the values as needed:
```bash
cp .env.example .env
```

## Running the Application

### Using Docker (Recommended)

1. Build the containers:
```bash
make build
# or without make:
docker-compose build
```

2. Start the services:
```bash
make up
# or without make:
docker-compose up -d
```

3. View the logs:
```bash
make logs
# or without make:
docker-compose logs -f
```

4. Stop the services:
```bash
make down
# or without make:
docker-compose down
```

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Ensure PostgreSQL is running and accessible with the credentials in your `.env` file

3. Run the application:
```bash
go run cmd/api/main.go
```

## Project Structure

```
freight-broker/
├── cmd/
│   └── api/            # Application entrypoint
├── configs/            # Configuration
├── internal/
│   ├── models/         # Database models
│   ├── services/       # Business logic
│   ├── controllers/    # HTTP handlers
│   ├── interfaces/     # Interfaces for external services
│   ├── dto/            # Data Transfer Objects
│   └── repository/     # Database operations
├── pkg/
│   └── tms/           # TMS integration packages
└── configs/           # Configuration files
```

## API Documentation

(TBD - Add API endpoints and usage examples)

## Development

### Running Tests
```bash
make test
# or without make:
go test -v ./...
```

### Adding New TMS Provider

1. Implement the TMSProvider interface in `internal/interfaces/tms.go`
2. Create a new provider in `pkg/tms/`
3. Update the configuration to support the new provider

## Environment Variables

- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
