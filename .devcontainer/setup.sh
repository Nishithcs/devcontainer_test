#!/bin/bash

set -e

echo "🚀 Setting up Cluster Code Service development environment..."

# Install essential Go development tools
echo "📦 Installing Go development tools..."
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/go-delve/delve/cmd/dlv@latest

# Create golangci-lint config
mkdir -p ~/.config/golangci-lint
cat > ~/.config/golangci-lint/config.yaml << 'EOF'
run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - unused
    - misspell

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
EOF

# Create Air config for hot reload
mkdir -p ~/.air
cat > ~/.air/config.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF

# Set up Git configuration
git config --global init.defaultBranch main
git config --global pull.rebase false

echo "✅ Development environment setup complete!"
echo ""
echo "📋 Available commands:"
echo "  make help          - Show all available commands"
echo "  make dev           - Run API with hot reload"
echo "  make dev-env       - Start all services"
echo "  make test          - Run tests"
echo "  make lint          - Run linter"
echo "  make build         - Build the application"
echo ""
echo "🌐 Services will be available at:"
echo "  API Server:        http://localhost:8070"
echo "  RabbitMQ Management: http://localhost:15672"
echo "  Mongo Express:     http://localhost:8091"
echo "  Asynq Dashboard:   http://localhost:8081"
echo ""
echo "🎉 Ready to start developing!" 