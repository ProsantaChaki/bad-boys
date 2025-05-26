# Bad Boyes - Go Authentication API

A simple authentication API built with Go, Gin, and MySQL.

## Prerequisites

- Go 1.16 or higher
- MySQL 5.7 or higher

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd bad_boyes
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory with the following content:
```
DB_USER=your_mysql_username
DB_PASSWORD=your_mysql_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=bad_boyes
JWT_SECRET=your_jwt_secret_key_here
```

4. Create the database and tables:
```bash
mysql -u root -p < schema.sql
```

## Running the Application

```bash
go run main.go
```

The server will start on port 3000.

## API Endpoints

### Authentication

#### Register
- **POST** `/auth/register`
- Body:
```json
{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
}
```

#### Login
- **POST** `/auth/login`
- Body:
```json
{
    "email": "john@example.com",
    "password": "password123"
}
```

#### Get Profile (Protected)
- **GET** `/auth/profile`
- Headers:
  - `Authorization: Bearer <token>`

## Project Structure

```
.
├── internal/
│   ├── middleware/    # Authentication middleware
│   ├── models/        # Data models
│   ├── routes/        # Route handlers
│   └── services/      # Business logic
├── pkg/
│   └── database/      # Database connection
├── .env              # Environment variables
├── go.mod            # Go module file
├── main.go           # Application entry point
├── README.md         # This file
└── schema.sql        # Database schema
``` 