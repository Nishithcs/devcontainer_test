# Cluster Code Service DevContainer

This directory contains the development container configuration for the Cluster Code Service project. The devcontainer provides a clean Go development environment that works with your existing `docker-compose.yml` file.

## 🚀 Quick Start

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) or Docker Engine
- [VS Code](https://code.visualstudio.com/) with the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- [DevPod CLI](https://devpod.sh/docs/getting-started/installation) (optional, for DevPod integration)

### Starting the DevContainer

1. **Using VS Code:**
   - Open the project in VS Code
   - Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on Mac)
   - Select "Dev Containers: Reopen in Container"
   - Wait for the container to build and start

2. **Using DevPod:**
   ```bash
   devpod up
   ```

## 🛠️ What's Included

### Base Image
- **Go 1.24** on Debian Bullseye
- **Docker-in-Docker** support
- **Git** and **GitHub CLI**

### Go Development Tools
- `golangci-lint` - Code linting
- `air` - Hot reload for development
- `migrate` - Database migrations
- `dlv` - Debugger
- `goimports` - Code formatting

### VS Code Extensions
- **Go** - Official Go extension
- **Docker** - Docker support
- **GitLens** - Git integration
- **YAML** - YAML support
- **Makefile Tools** - Makefile support

### Services
The devcontainer uses your existing `docker-compose.yml` file to run services. You can start them with:
```bash
make dev-env
```

This will start:
- **PostgreSQL 15** (port 5422)
- **Redis 7** (port 6379)
- **RabbitMQ 3.9** with Management UI (ports 5672, 15672)
- **MongoDB 4** (port 27017)
- **Mongo Express** (port 8091)
- **Asynq Dashboard** (port 8081)

### Port Forwarding
The following ports are automatically forwarded:
- `8070` - API Server
- `5422` - PostgreSQL
- `6379` - Redis
- `5672` - RabbitMQ
- `15672` - RabbitMQ Management UI
- `27017` - MongoDB
- `8091` - Mongo Express
- `8081` - Asynq Dashboard
- `83` - Reverse Proxy

## 📋 Available Commands

Once inside the devcontainer, you can use these commands:

```bash
# Show all available commands
make help

# Development
make dev              # Run API with hot reload
make dev-consumer     # Run consumer in development mode
make dev-env          # Start all services
make dev-env-down     # Stop all services

# Building
make build            # Build the application
make build-api        # Build API only
make build-consumer   # Build consumer only
make clean            # Clean build artifacts

# Testing & Quality
make test             # Run tests
make coverage         # Run tests with coverage
make lint             # Run linter
make fmt              # Format code
make vet              # Run go vet

# Database
make migrate          # Run migrations
make seed             # Seed database

# Dependencies
make deps             # Install development dependencies
```

## 🔧 Configuration

### Environment Variables
The devcontainer automatically sets up Go environment variables. For service configuration, use your existing `.env` file or the environment variables defined in your `docker-compose.yml`.

### Go Configuration
- **GOPATH**: `/go`
- **GOROOT**: `/usr/local/go`
- **Go modules**: Enabled
- **CGO**: Enabled

### VS Code Settings
The devcontainer includes optimized VS Code settings for Go development:
- Auto-formatting on save
- Auto-import organization
- Linting on save
- Excluded directories (bin, tmp, vendor)

## 🐳 Docker Support

The devcontainer includes full Docker support:
- Docker-in-Docker capability
- Docker Compose v2
- Access to host Docker socket
- BuildKit enabled

You can use your existing `docker-compose.yml` file directly:
```bash
docker-compose up -d
```

## 🔍 Debugging

### Go Debugging
The devcontainer includes Delve debugger. You can debug your Go applications using:
- VS Code's built-in debugging features
- Command line debugging with `dlv`

### Database Debugging
- **PostgreSQL**: Use any PostgreSQL client to connect to `localhost:5422`
- **Redis**: Use Redis CLI or GUI clients to connect to `localhost:6379`
- **RabbitMQ**: Access management UI at `http://localhost:15672`
- **MongoDB**: Use MongoDB Compass or mongo shell to connect to `localhost:27017`

## 📁 Project Structure

```
.devcontainer/
├── devcontainer.json    # Main devcontainer configuration
├── setup.sh            # Setup script for tools and configuration
└── README.md           # This file

# Key project directories
cmd/                    # Application entry points
internal/               # Private application code
├── api/               # API handlers and middleware
├── config/            # Configuration management
├── data/              # Data layer (models, repositories)
├── services/          # Business logic services
└── utils/             # Utility functions
```

## 🚨 Troubleshooting

### Common Issues

1. **Port conflicts**: If ports are already in use, the devcontainer will show warnings. You can change ports in your `docker-compose.yml`.

2. **Permission issues**: The devcontainer runs as the `vscode` user (UID 1000). If you encounter permission issues, check file ownership.

3. **Docker not accessible**: Ensure Docker Desktop is running and the Docker socket is accessible.

4. **Go tools not found**: Run `make deps` to install all Go development tools.

### Logs
- **Container logs**: Check VS Code's Dev Container output panel
- **Service logs**: Use `docker-compose logs [service-name]`
- **Application logs**: Check the application's log output

## 🔄 Updates

To update the devcontainer:
1. Modify `devcontainer.json` or `setup.sh`
2. Rebuild the container: "Dev Containers: Rebuild Container"
3. Or delete and recreate: "Dev Containers: Rebuild Container Without Cache"

## 📚 Additional Resources

- [Dev Containers Documentation](https://containers.dev/)
- [DevPod Documentation](https://devpod.sh/docs)
- [Go Development in Containers](https://golang.org/doc/install)
- [VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.go)

## 🤝 Contributing

When contributing to this project:
1. Test your changes in the devcontainer
2. Ensure all tests pass: `make test`
3. Run the linter: `make lint`
4. Update this README if you modify the devcontainer configuration 