#!/bin/bash
# Quick setup script for development environment

set -e

echo "=========================================="
echo "🏋️ Gym Pro Backend Setup"
echo "=========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "Please install Go 1.24 or higher from https://go.dev/dl/"
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Warning: Docker is not installed"
    echo "Docker is needed to run PostgreSQL. Install from https://www.docker.com/"
else
    echo "✅ Docker is installed"
fi

# Install dependencies
echo ""
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy

# Install tools
echo ""
echo "🔧 Installing development tools..."

# Install golang-migrate
if ! command -v migrate &> /dev/null; then
    echo "Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
else
    echo "✅ golang-migrate already installed"
fi

# Install swag
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
else
    echo "✅ swag already installed"
fi

# Create .env if not exists
if [ ! -f .env ]; then
    echo ""
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please update with your configuration."
else
    echo "✅ .env file already exists"
fi

echo ""
echo "=========================================="
echo "✅ Setup Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Start PostgreSQL: make docker-up"
echo "2. Run migrations: make migrate-up"
echo "3. (Optional) Seed data: psql -U gymadmin -d gym_pro_db -f migrations/seed_exercises.sql"
echo "4. Start server: make run"
echo "5. Visit: http://localhost:8080/health"
echo ""
echo "For more info, see GETTING_STARTED.md"
echo ""
