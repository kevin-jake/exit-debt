#!/bin/bash

echo "🚀 Setting up Go Debt Tracker..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check if PostgreSQL is running
if ! pg_isready -q; then
    echo "❌ PostgreSQL is not running. Please start PostgreSQL first."
    exit 1
fi

# Install Air for hot reloading
echo "📦 Installing Air for hot reloading..."
go install github.com/cosmtrek/air@latest

# Download dependencies
echo "📥 Downloading dependencies..."
go mod download

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created. Please update it with your database credentials."
fi

# Create tmp directory for Air
mkdir -p tmp

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env file with your database credentials"
echo "2. Create a PostgreSQL database named 'debt_tracker'"
echo "3. Run 'air' to start the development server"
echo ""
echo "Happy coding! 🎉" 