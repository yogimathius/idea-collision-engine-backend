#!/bin/bash

# Idea Collision Engine API Deployment Script
set -e

echo "🚀 Deploying Idea Collision Engine API to Fly.io..."

# Check if flyctl is installed
if ! command -v flyctl &> /dev/null; then
    echo "❌ flyctl could not be found. Please install Fly.io CLI first."
    echo "   Visit: https://fly.io/docs/flyctl/install/"
    exit 1
fi

# Check if we're logged in to Fly.io
if ! flyctl auth whoami &> /dev/null; then
    echo "🔑 Please login to Fly.io first:"
    flyctl auth login
fi

# Set up production environment variables
echo "🔧 Setting up environment variables..."

# Required environment variables
REQUIRED_VARS=(
    "DATABASE_URL"
    "OPENAI_API_KEY"  
    "STRIPE_SECRET_KEY"
    "JWT_SECRET"
)

# Check if required environment variables are set
for var in "${REQUIRED_VARS[@]}"; do
    if [[ -z "${!var}" ]]; then
        echo "❌ Required environment variable $var is not set"
        echo "   Please set it with: export $var=your_value"
        echo "   Or add it to your .env file"
        exit 1
    fi
done

# Set secrets in Fly.io
echo "🔐 Setting application secrets..."
flyctl secrets set \
    DATABASE_URL="$DATABASE_URL" \
    OPENAI_API_KEY="$OPENAI_API_KEY" \
    STRIPE_SECRET_KEY="$STRIPE_SECRET_KEY" \
    JWT_SECRET="$JWT_SECRET" \
    REDIS_URL="${REDIS_URL:-redis://redis.internal:6379}" \
    CORS_ORIGINS="${CORS_ORIGINS:-https://your-frontend-domain.com}"

# Create PostgreSQL database if it doesn't exist
echo "🗄️ Setting up PostgreSQL database..."
if ! flyctl postgres list | grep -q "collision-engine-db"; then
    echo "   Creating new PostgreSQL database..."
    flyctl postgres create collision-engine-db --region sjc --vm-size shared-cpu-1x --volume-size 10
    
    # Get database URL
    DB_URL=$(flyctl postgres connect -a collision-engine-db --command "echo \$DATABASE_URL")
    flyctl secrets set DATABASE_URL="$DB_URL"
fi

# Create Redis instance if needed  
echo "🔄 Setting up Redis cache..."
if ! flyctl redis list | grep -q "collision-engine-redis"; then
    echo "   Creating new Redis instance..."
    flyctl redis create collision-engine-redis --region sjc
    
    # Get Redis URL
    REDIS_URL=$(flyctl redis status collision-engine-redis --json | jq -r '.redis_url')
    flyctl secrets set REDIS_URL="$REDIS_URL"
fi

# Deploy the application
echo "📦 Building and deploying application..."
flyctl deploy --build-arg ENVIRONMENT=production

# Run database migrations
echo "🔄 Running database migrations..."
flyctl ssh console -C "./migrate"

# Check deployment health
echo "🏥 Checking deployment health..."
sleep 10
HEALTH_URL="https://idea-collision-engine-api.fly.dev/health"
if curl -s "$HEALTH_URL" | grep -q "healthy"; then
    echo "✅ Deployment successful! API is healthy."
    echo "🌐 API URL: https://idea-collision-engine-api.fly.dev"
    echo "📊 Health Check: $HEALTH_URL"
else
    echo "⚠️  Deployment completed but health check failed."
    echo "   Check logs with: flyctl logs"
fi

echo ""
echo "🎉 Deployment complete!"
echo ""
echo "Next steps:"
echo "1. Test your API endpoints"
echo "2. Update your frontend to use the production API URL"
echo "3. Set up monitoring and alerts"
echo "4. Configure your custom domain (optional)"
echo ""
echo "Useful commands:"
echo "  flyctl logs -a idea-collision-engine-api  # View logs"
echo "  flyctl ssh console                        # SSH into the container" 
echo "  flyctl status                             # Check app status"