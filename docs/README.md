# 🚀 Idea Collision Engine API Documentation

The Idea Collision Engine API is a creative productivity service that generates unexpected but relevant idea combinations to help break out of familiar patterns and spark innovation.

## 📖 Quick Start

### Base URL
- **Development**: `http://localhost:8080`
- **Production**: `https://api.ideacollisionengine.com`

### Interactive Documentation
- **Swagger UI**: Visit `/docs/` for interactive API documentation
- **OpenAPI Spec**: Download the OpenAPI 3.0 specification at `/docs/openapi.yaml`

## 🔐 Authentication

Most API endpoints require authentication using JWT tokens.

### Getting Started
1. **Register**: `POST /api/auth/register`
2. **Login**: `POST /api/auth/login` to get your access token
3. **Use Token**: Include in Authorization header: `Authorization: Bearer <token>`

### Example Authentication Flow

```bash
# Register a new account
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123",
    "name": "John Doe"
  }'

# Login to get access token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'

# Use the token in subsequent requests
curl -X POST http://localhost:8080/api/collisions/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token-here>" \
  -d '{
    "interests": ["technology", "art", "sustainability"],
    "intensity": "moderate"
  }'
```

## 🎯 Core Features

### Collision Generation
Generate creative idea combinations based on your interests:

- **Interests**: 1-5 areas of interest (e.g., "technology", "art", "wellness")
- **Intensity Levels**: `gentle`, `moderate`, `bold`, `extreme`
- **Context**: Optional context to guide generation

### Subscription Tiers

| Feature | Free | Pro | Team |
|---------|------|-----|------|
| Collisions per week | 50 | Unlimited | Unlimited |
| Basic domains | ✅ | ✅ | ✅ |
| Premium domains | ❌ | ✅ | ✅ |
| Rate limiting | 10/min | None | None |
| API access | ✅ | ✅ | ✅ |

## 📊 Rate Limiting

The API implements rate limiting to ensure fair usage:

- **Free Users**: 10 requests per minute
- **Premium Users**: No rate limits
- **Headers**: Rate limit info in `X-RateLimit-*` headers

When rate limited, you'll receive a `429 Too Many Requests` response with retry information.

## 🔍 Health Monitoring

### Service Health
Check overall service health:
```bash
curl http://localhost:8080/health
```

Response includes status of:
- Database connectivity
- Redis cache
- AI service availability
- Collision engine status

### Collision Service Health
Check collision-specific service health:
```bash
curl http://localhost:8080/api/collisions/health
```

## 📝 API Endpoints Overview

### Authentication & User Management
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile
- `PUT /api/auth/profile` - Update user profile

### Collision Generation
- `POST /api/collisions/generate` - Generate idea collision
- `GET /api/collisions/history` - Get collision history
- `PUT /api/collisions/{id}/rate` - Rate a collision
- `GET /api/collisions/usage` - Check usage status

### Domains & Content
- `GET /api/domains/basic` - Get basic domains
- `GET /api/domains/premium` - Get premium domains (auth required)

### Subscription Management
- `GET /api/subscriptions/plans` - Get pricing plans
- `POST /api/subscriptions/checkout` - Create checkout session
- `GET /api/subscriptions/status` - Get subscription status
- `POST /api/subscriptions/cancel` - Cancel subscription
- `POST /api/subscriptions/webhook` - Stripe webhook handler

## 🚨 Error Handling

The API uses standard HTTP status codes and returns detailed error information:

```json
{
  "error": "rate_limit_exceeded",
  "message": "Rate limit exceeded. Try again in 45 seconds",
  "code": 429
}
```

### Common Error Codes
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (invalid/missing token)
- `402` - Payment Required (usage limit exceeded)
- `403` - Forbidden (premium feature required)
- `404` - Not Found
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

## 🔧 Development Setup

### Prerequisites
- Go 1.23+
- PostgreSQL 16+
- Redis 7+
- Node.js 20+ (for frontend)

### Environment Variables
```env
NODE_ENV=development
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/collision_engine
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
OPENAI_API_KEY=your-openai-key
STRIPE_SECRET_KEY=your-stripe-key
CORS_ORIGINS=http://localhost:3000
```

### Running with Docker
```bash
# Start all services
docker-compose up

# Development mode with hot reload
docker-compose --profile dev up
```

### Manual Setup
```bash
# Install dependencies
go mod download
pnpm install

# Run migrations
./migrate

# Start the API server
go run cmd/server/main.go

# Start the frontend (separate terminal)
pnpm run dev
```

## 📈 Monitoring & Observability

### Health Checks
- **Liveness**: `/health` - Service is running
- **Readiness**: `/api/collisions/health` - Service dependencies are healthy

### Metrics
The API provides built-in monitoring for:
- Request latency and throughput
- Error rates by endpoint
- Authentication success/failure rates
- Collision generation performance
- Database and Redis connectivity

### Logging
Structured logging with:
- Request/response logging
- Error tracking with stack traces
- Performance metrics
- Security events (failed auth, rate limiting)

## 🔒 Security Features

- **JWT Authentication** with secure token generation
- **Rate Limiting** to prevent abuse
- **Input Validation** on all endpoints
- **CORS Configuration** for cross-origin requests
- **SQL Injection Protection** with parameterized queries
- **Password Hashing** with bcrypt
- **Environment-based Configuration** for secrets

## 🚀 Deployment

### Production Deployment
1. Configure environment variables
2. Run database migrations
3. Deploy using Docker:
   ```bash
   ./deploy.sh
   ```

### Scaling Considerations
- **Horizontal Scaling**: Multiple API instances behind load balancer
- **Database**: PostgreSQL with read replicas
- **Caching**: Redis cluster for session management
- **AI Service**: Rate limiting and fallback strategies

## 📞 Support & Contributing

### Getting Help
- **Documentation**: Visit `/docs/` for interactive API docs
- **Health Status**: Check `/health` for service status
- **Issues**: Report bugs and feature requests on GitHub

### Contributing
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

---

**Built with ❤️ using Go, Fiber, PostgreSQL, and Redis**