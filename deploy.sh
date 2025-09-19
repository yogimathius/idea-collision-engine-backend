#!/bin/bash

# Idea Collision Engine Deployment Script
# This script handles Docker-based deployment for production

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="idea-collision-engine"
DOCKER_COMPOSE_FILE="docker-compose.yml"
ENV_FILE=".env.production"

echo -e "${GREEN}🚀 Starting Idea Collision Engine deployment...${NC}"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker Compose is not installed.${NC}"
    exit 1
fi

# Check for required environment variables
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${YELLOW}⚠️  Creating example environment file...${NC}"
    cat > "$ENV_FILE" << EOF
# Production Environment Configuration
NODE_ENV=production
PORT=8080

# Database Configuration
DATABASE_URL=postgres://postgres:secure_password@postgres:5432/collision_engine?sslmode=disable

# Redis Configuration
REDIS_URL=redis://redis:6379

# Security
JWT_SECRET=your-super-secure-jwt-secret-key-change-this-in-production

# External Services
OPENAI_API_KEY=your-openai-api-key
STRIPE_SECRET_KEY=your-stripe-secret-key

# CORS Origins (comma-separated)
CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
EOF
    
    echo -e "${YELLOW}⚠️  Please edit $ENV_FILE with your production values before continuing.${NC}"
    echo -e "${YELLOW}   Then run this script again.${NC}"
    exit 1
fi

# Load environment variables
set -a
source "$ENV_FILE"
set +a

# Function to check service health
check_health() {
    local service=$1
    local max_attempts=30
    local attempt=1
    
    echo -e "${YELLOW}🔍 Checking $service health...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if docker-compose ps -q "$service" | xargs docker inspect -f '{{.State.Health.Status}}' 2>/dev/null | grep -q "healthy"; then
            echo -e "${GREEN}✅ $service is healthy${NC}"
            return 0
        fi
        
        echo -e "${YELLOW}⏳ Waiting for $service to be healthy (attempt $attempt/$max_attempts)...${NC}"
        sleep 5
        attempt=$((attempt + 1))
    done
    
    echo -e "${RED}❌ $service failed to become healthy${NC}"
    return 1
}

# Pull latest images
echo -e "${YELLOW}📦 Pulling latest images...${NC}"
docker-compose pull

# Build the application
echo -e "${YELLOW}🔨 Building application...${NC}"
docker-compose build --no-cache

# Stop existing containers
echo -e "${YELLOW}🛑 Stopping existing containers...${NC}"
docker-compose down

# Start infrastructure services first
echo -e "${YELLOW}🗄️  Starting database and cache...${NC}"
docker-compose up -d postgres redis

# Wait for infrastructure to be healthy
check_health "postgres"
check_health "redis"

# Run database migrations
echo -e "${YELLOW}🔄 Running database migrations...${NC}"
docker-compose run --rm app ./migrate

# Start the application
echo -e "${YELLOW}🚀 Starting application...${NC}"
docker-compose up -d app

# Wait for application to be healthy
check_health "app"

# Show status
echo -e "${GREEN}📊 Deployment Status:${NC}"
docker-compose ps

# Show logs for verification
echo -e "${YELLOW}📝 Recent application logs:${NC}"
docker-compose logs --tail=50 app

echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo -e "${GREEN}🌐 Application is available at: http://localhost:${PORT}${NC}"
echo -e "${GREEN}🔍 Health check: http://localhost:${PORT}/health${NC}"

# Cleanup old images
echo -e "${YELLOW}🧹 Cleaning up old images...${NC}"
docker image prune -f

echo -e "${GREEN}🎉 Idea Collision Engine is now running in production mode!${NC}"