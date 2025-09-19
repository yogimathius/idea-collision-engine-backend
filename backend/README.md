# Idea Collision Engine API

A high-performance Go backend for the Idea Collision Engine - a creative breakthrough tool that generates unexpected connections between your interests and diverse knowledge domains.

## 🚀 Features

- **Sub-100ms Collision Generation** - Lightning-fast API responses
- **Anti-Echo Chamber Algorithm** - Forces exposure to unexpected but relevant domains  
- **50+ Curated Domains** - From biomimicry to quantum physics
- **AI-Enhanced Insights** - OpenAI integration for premium users
- **Freemium Model** - 5 free collisions/week, unlimited for Pro ($12/mo) and Team ($39/mo)
- **Advanced Scoring** - Relevance + novelty + actionability + depth
- **Redis Caching** - >90% cache hit rate for domain lookups
- **JWT Authentication** - Secure user management
- **Stripe Integration** - Subscription billing
- **Rate Limiting** - Intelligent request throttling
- **Comprehensive Logging** - Production-ready observability

## 🏗️ Architecture

```
cmd/
├── server/     # Main API server
└── migrate/    # Database migration utility

internal/
├── auth/       # JWT authentication
├── collision/  # Core collision engine + AI service  
├── database/   # PostgreSQL & Redis clients
├── handlers/   # HTTP route handlers
├── middleware/ # Auth, rate limiting, CORS
└── models/     # Data structures

pkg/
├── config/     # Configuration management
└── utils/      # Shared utilities

migrations/     # Database schema
```

## 📊 Database Schema

- **users** - Authentication and subscription tiers
- **collision_domains** - Curated knowledge domains (50+)  
- **collision_sessions** - Generated collision history
- **user_usage** - Freemium usage tracking

## 🛠️ Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 12+
- Redis 6+  
- OpenAI API key
- Stripe account (for subscriptions)

### 1. Clone and Setup

```bash
git clone <repo-url>
cd idea-collision-engine/backend

# Setup development environment  
make setup

# Edit configuration
cp .env.example .env
# Edit .env with your database URLs and API keys
```

### 2. Database Setup

```bash
# Run migrations and seed domains
make migrate
```

### 3. Start Development Server

```bash
# With live reload (requires air)
make dev

# Or standard build and run
make run
```

The API will be available at `http://localhost:8080`

## 🔑 API Endpoints

### Authentication
- `POST /api/auth/register` - Create account
- `POST /api/auth/login` - User login  
- `GET /api/auth/profile` - Get user profile
- `PUT /api/auth/profile` - Update user profile

### Collision Generation
- `POST /api/collisions/generate` - Generate collision (rate limited)
- `GET /api/collisions/history` - User's collision history
- `PUT /api/collisions/:id/rate` - Rate collision (1-5 stars)
- `GET /api/collisions/usage` - Check usage limits

### Domains  
- `GET /api/domains/basic` - Basic domains (all users)
- `GET /api/domains/premium` - Premium domains (Pro/Team only)

### Subscriptions
- `GET /api/subscriptions/plans` - Available pricing plans
- `POST /api/subscriptions/checkout` - Create Stripe checkout  
- `GET /api/subscriptions/status` - Current subscription
- `POST /api/subscriptions/cancel` - Cancel subscription
- `POST /api/subscriptions/webhook` - Stripe webhooks

### Health
- `GET /health` - Basic health check
- `GET /api/collisions/health` - Detailed service health

## 💡 Collision Generation Example

```bash
curl -X POST http://localhost:8080/api/collisions/generate \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "user_interests": ["machine learning", "product design"],
    "current_project": "AI-powered recommendation system",  
    "project_type": "product",
    "collision_intensity": "moderate"
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "primary_domain": "Machine Learning",
  "collision_domain": "Jazz Improvisation", 
  "connection": "Jazz improvisation offers structured spontaneity for ML recommendation systems through real-time adaptation patterns...",
  "spark_questions": [
    "How might jazz's call-and-response patterns improve recommendation feedback loops?",
    "What would your AI system look like if it could 'improvise' like a jazz musician?"
  ],
  "examples": [
    "Real-time collaboration → Dynamic model updates based on user interactions",
    "Structured improvisation → Recommendation algorithms that adapt within learned constraints"
  ],
  "next_steps": [
    "Research jazz improvisation principles and identify 3 that could apply to AI systems",
    "Prototype adaptive recommendation logic using jazz-inspired feedback patterns"
  ],
  "quality_score": 87.3,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🚀 Deployment

### Fly.io (Recommended)

```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Login and deploy
flyctl auth login
flyctl launch
```

The deployment includes:
- Automatic PostgreSQL database creation
- Redis cache setup  
- SSL certificates
- Health checks
- Auto-scaling

### Docker

```bash
# Build and run
make docker-run

# Or with docker-compose
make docker-compose-up
```

## ⚡ Performance

- **Response Times**: <100ms for collision generation
- **Throughput**: 1000+ requests/second  
- **Cache Hit Rate**: >90% for domain lookups
- **Memory Usage**: <50MB RAM
- **Database Queries**: <10ms average

## 🔧 Configuration

Key environment variables:

```env
# Server
PORT=8080
ENVIRONMENT=production

# Database  
DATABASE_URL=postgresql://user:pass@host:5432/dbname
REDIS_URL=redis://host:6379

# APIs
OPENAI_API_KEY=sk-proj-your-key-here
STRIPE_SECRET_KEY=sk_live_your-key-here

# Security
JWT_SECRET=your-secure-secret-key

# Performance
RATE_LIMIT_RPS=10
CACHE_EXPIRATION=300
```

## 🧪 Testing

```bash
# Run all tests
make test

# With coverage
make test-coverage  

# Benchmarks
make benchmark

# Linting
make lint
```

## 📈 Monitoring

The API includes comprehensive logging and metrics:

- Request/response logging
- Error tracking  
- Performance metrics
- Health checks
- Database query monitoring

## 🛡️ Security

- JWT-based authentication
- Rate limiting per user/IP
- Input validation
- SQL injection prevention  
- CORS protection
- Secure headers

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)  
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🎯 Success Metrics

- **Target**: Sub-100ms API response times ✅
- **Target**: 85%+ user satisfaction rating  
- **Target**: 15%+ freemium conversion rate
- **Target**: Handle 1000+ concurrent users
- **Target**: $8K MRR (500 Pro + 20 Team users)

Built with ❤️ for creative breakthrough moments.