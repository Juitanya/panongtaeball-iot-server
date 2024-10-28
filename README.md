# Project GoAgent 1.0 10/2024 (Merida)

GoAgent คือ Template โครงสร้างพื้นฐานที่ หลายๆ Service ควรจะมี

## 📋 Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Project Structure](#project-structure)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## ✨ Features

- Feature 1: Brief description
- Feature 2: Brief description
- Feature 3: Brief description

## 📌 Requirements

- Go 1.20 or higher
- PostgreSQL 14.0 or higher
- Other dependencies... (มีหละแต่คิดยังไม่ออก)

## 💻 Installation

```bash
# Clone the repository
git clone https://github.com/username/project.git

# Change directory
cd project

# Install dependencies
go mod download

# Build the project
go build -o app
```

## 🚀 Quick Start

```bash
# Set up environment variables
cp .env.example .env

# Edit your .env file with your configurations
vim .env

# Run migrations
atlas migrate apply

# Start the server
./app
```

## ⚙️ Configuration

The application can be configured using environment variables or a configuration file.

### Environment Variables

```env
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=postgres
DB_PASSWORD=secret

# Security
PASSWORD_PEPPER=your-secret-pepper
JWT_SECRET=your-jwt-secret
```

## 📖 API Documentation

### Authentication

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
}
```

## 🗄️ Database Schema

```hcl
table "users" {
    schema = schema.public
    column "id" {
        null = false
        type = bigserial
    }
    column "email" {
        null = false
        type = varchar(255)
    }
    column "password_hash" {
        null = false
        type = varchar(255)
    }
    # ... other columns
}
```

## 📁 Project Structure

```
.
├── cmd/                  # Application entrypoints
├── internal/            # Private application and library code
│   ├── api/            # API handlers
│   ├── auth/           # Authentication package
│   ├── config/         # Configuration handling
│   ├── db/             # Database operations and migrations
│   └── models/         # Data models
├── pkg/                # Library code that's safe to use by external apps
├── scripts/           # Scripts for development and deployment
├── test/             # Additional test applications and test data
└── web/              # Web application specific components
```

## 🔧 Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running Linter

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## 🚢 Deployment

### Docker

```bash
# Build the Docker image
docker build -t myapp .

# Run the container
docker run -p 8080:8080 myapp
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

## 👥 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📧 Contact

Your Name - [@yourtwitter](https://twitter.com/yourtwitter) - email@example.com

Project Link: [https://github.com/username/repo](https://github.com/username/repo)

---
Made with ❤️ by [Your Name](https://github.com/username)