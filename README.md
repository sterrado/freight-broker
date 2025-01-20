# Freight Broker Application

A Go-based freight broker application for managing loads through various TMS (Transportation Management System) providers.

## Prerequisites

- Docker and Docker Compose
- Make (optional, for using Makefile commands)
- Go 1.21 or later (for local development)

## Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/sterrado/freight-broker.git
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

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
# or using yarn:
yarn install
```

3. Create a .env file for the frontend:
```bash
echo "REACT_APP_API_BASE_URL=http://localhost:8080/api" > .env
```

4. Start the development server:
```bash
npm start
# or using yarn:
yarn start
```

The frontend will be available at http://localhost:3000

## Project Structure

```
freight-broker/
├── frontend/           # React frontend application
│   ├── src/           # Frontend source code
│   ├── public/        # Static assets
│   └── package.json   # Frontend dependencies
└── backend/           # Go backend application
    ├── cmd/
    │   └── api/       # Application entrypoint
    ├── configs/       # Configuration files
    └── internal/      # Internal packages
        ├── models/    # Database models
        ├── services/  # Business logic
        ├── controllers/   # HTTP handlers
        ├── interfaces/    # Interfaces for external services
        ├── middleware/    # Middleware components
        └── dto/          # Data Transfer Objects
```

The project is organized into two main directories:
- `frontend/`: Contains the React application
- `backend/`: Contains the Go API server with all its components

## Authentication

The application uses JWT (JSON Web Token) for API authentication. To use the API:

1. Obtain a JWT token by calling the login endpoint:
```bash
POST /api/auth/login
Content-Type: application/json

{
    "username": "your-username",
    "password": "your-password"
}
```

2. Use the token in subsequent requests:
```bash
GET /api/loads
Authorization: Bearer <your-jwt-token>
```

### Protected Endpoints
All API endpoints except `/api/auth/login` require a valid JWT token in the Authorization header.

## API Documentation

### Authentication Endpoints

#### Login
```
POST /api/auth/login
Content-Type: application/json

Request:
{
    "username": "string",
    "password": "string"
}

Response:
{
    "token": "string"
}
```

### Load Management Endpoints

All load management endpoints require authentication.

#### Create Load
```
POST /api/loads
Authorization: Bearer <token>
Content-Type: application/json

Request: {
    // Load details
}
```

#### List Loads
```
GET /api/loads?page=1&size=10
Authorization: Bearer <token>
```

#### Get Load
```
GET /api/loads/:id
Authorization: Bearer <token>
```

## Environment Variables

Use .env.example to create an .env file and replace the values.
