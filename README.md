# GophKeeper

GophKeeper is a secure client-server password manager built in Go. It allows users to safely store and synchronize passwords, text data, binary files, and bank card information across multiple devices.

## Features

### Server Features
- User registration, authentication, and authorization
- Secure data storage with encryption
- Data synchronization between multiple clients
- RESTful API for client communication
- JWT-based authentication
- PostgreSQL database support

### Client Features
- Cross-platform CLI application (Windows, Linux, macOS)
- Secure authentication with remote server
- Support for multiple data types:
  - Login/password pairs
  - Arbitrary text data
  - Binary data
  - Bank card information
- Local data caching and synchronization
- Version information and build date display

## Data Types

The system supports the following data types:

1. **Login/Password** - Store website credentials with optional metadata
2. **Text** - Store arbitrary text data
3. **Binary** - Store binary files and data
4. **Bank Card** - Store bank card information securely

All data types support custom metadata for additional context.

## Installation

### Prerequisites

- Go 1.23.6 or later
- PostgreSQL 12 or later (for server)
- Git

### Building from Source

1. Clone the repository:
```bash
git clone <repository-url>
cd Gophkeeper
```

2. Download dependencies:
```bash
make deps
```

3. Build the application:
```bash
make build
```

4. Build for all platforms:
```bash
make build-all
```

### Database Setup

1. Install and start PostgreSQL
2. Create a database:
```bash
createdb gophkeeper
```

3. The server will automatically create the necessary tables on first run.

## Usage

### Starting the Server

```bash
# Using make
make run-server

# Or directly
go run ./cmd/server

# With custom database settings
go run ./cmd/server -db-host=localhost -db-port=5432 -db-user=gophkeeper -db-password=password -db-name=gophkeeper
```

### Using the Client

```bash
# Register a new user
./bin/gophkeeper-client register username email@example.com password

# Login
./bin/gophkeeper-client login username password

# Add login/password data
./bin/gophkeeper-client add login_password "My Website" "username" "password" "https://example.com" "Additional notes"

# Add text data
./bin/gophkeeper-client add text "Important Note" "This is my important note"

# Add bank card data
./bin/gophkeeper-client add bank_card "My Credit Card" "1234567890123456" "12/25" "123" "John Doe" "Bank Name" "Additional notes"

# List all data
./bin/gophkeeper-client list

# Get specific data
./bin/gophkeeper-client get <data-id>

# Delete data
./bin/gophkeeper-client delete <data-id>

# Synchronize with server
./bin/gophkeeper-client sync

# Show version information
./bin/gophkeeper-client version
```

## API Endpoints

### Authentication
- `POST /api/v1/register` - Register a new user
- `POST /api/v1/login` - Authenticate user

### Data Management
- `GET /api/v1/data` - Get all user data
- `POST /api/v1/data` - Create new data
- `PUT /api/v1/data` - Update existing data
- `DELETE /api/v1/data?id=<id>` - Delete data

### Synchronization
- `POST /api/v1/sync` - Synchronize data with server

## Security

- All data is encrypted using AES-256-GCM before storage
- Passwords are hashed using SHA-256 with salt
- JWT tokens are used for authentication
- HTTPS is recommended for production deployments
- Client data is encrypted locally before transmission

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint
```

### Project Structure

```
Gophkeeper/
├── cmd/
│   ├── client/          # CLI client application
│   └── server/          # HTTP server application
├── internal/
│   ├── client/          # Client implementation
│   ├── crypto/          # Encryption and hashing
│   ├── database/        # Database operations
│   ├── models/          # Data models
│   └── server/          # Server implementation
├── scripts/
│   └── build.sh         # Cross-platform build script
├── Makefile             # Build and development commands
└── README.md
```

## Configuration

### Server Configuration

The server accepts the following command-line flags:

- `-port` - Server port (default: 8080)
- `-db-host` - Database host (default: localhost)
- `-db-port` - Database port (default: 5432)
- `-db-user` - Database user (default: gophkeeper)
- `-db-password` - Database password (default: password)
- `-db-name` - Database name (default: gophkeeper)
- `-jwt-secret` - JWT secret key (default: your-secret-key)
- `-encryption-key` - Data encryption key (default: your-encryption-key)

### Client Configuration

The client accepts the following command-line flags:

- `-server` - Server URL (default: http://localhost:8080)
- `-config` - Configuration directory (default: ~/.gophkeeper)
- `-version` - Show version information

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## Support

For issues and questions, please create an issue in the repository.
