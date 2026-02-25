# Quick setup script for development environment (Windows PowerShell)

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "🏋️ Gym Pro Backend Setup" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "✅ Go version: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Error: Go is not installed" -ForegroundColor Red
    Write-Host "Please install Go 1.24 or higher from https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}

# Check if Docker is installed
try {
    docker --version | Out-Null
    Write-Host "✅ Docker is installed" -ForegroundColor Green
} catch {
    Write-Host "❌ Warning: Docker is not installed" -ForegroundColor Yellow
    Write-Host "Docker is needed to run PostgreSQL. Install from https://www.docker.com/" -ForegroundColor Yellow
}

# Install dependencies
Write-Host ""
Write-Host "📦 Installing Go dependencies..." -ForegroundColor Yellow
go mod download
go mod tidy

# Install tools
Write-Host ""
Write-Host "🔧 Installing development tools..." -ForegroundColor Yellow

# Install golang-migrate
try {
    migrate -version | Out-Null
    Write-Host "✅ golang-migrate already installed" -ForegroundColor Green
} catch {
    Write-Host "Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
}

# Install swag
try {
    swag --version | Out-Null
    Write-Host "✅ swag already installed" -ForegroundColor Green
} catch {
    Write-Host "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
}

# Create .env if not exists
if (!(Test-Path .env)) {
    Write-Host ""
    Write-Host "📝 Creating .env file from .env.example..." -ForegroundColor Yellow
    Copy-Item .env.example .env
    Write-Host "✅ .env file created. Please update with your configuration." -ForegroundColor Green
} else {
    Write-Host "✅ .env file already exists" -ForegroundColor Green
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "✅ Setup Complete!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Start PostgreSQL: make docker-up" -ForegroundColor White
Write-Host "2. Run migrations: make migrate-up" -ForegroundColor White
Write-Host "3. (Optional) Seed data: see migrations/seed_exercises.sql" -ForegroundColor White
Write-Host "4. Start server: make run" -ForegroundColor White
Write-Host "5. Visit: http://localhost:8080/health" -ForegroundColor White
Write-Host ""
Write-Host "For more info, see GETTING_STARTED.md" -ForegroundColor Cyan
Write-Host ""
